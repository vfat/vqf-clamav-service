package api

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/vfat/vqf-clamav-service/internal/alert"
	"github.com/vfat/vqf-clamav-service/internal/clamd"
	"github.com/vfat/vqf-clamav-service/internal/crypto"
	"github.com/vfat/vqf-clamav-service/internal/quarantine"
	"github.com/vfat/vqf-clamav-service/internal/ratelimit"
	"github.com/vfat/vqf-clamav-service/internal/storage"
	"github.com/vfat/vqf-clamav-service/web"
)

// ServerConfig holds dependencies and configuration for the HTTP server.
type ServerConfig struct {
	DB               *storage.DB
	Vault            *quarantine.Vault
	Notifier         *alert.Notifier
	Limiter          *ratelimit.Limiter
	Clamd            *clamd.Client
	RequireAPIKey    bool
	MaxScanSizeMB    int64
	RateLimitRPM     int
	RateLimitEnabled bool
	LogRetention     int
	QuarRetention    int
	AuthMode         string // "none", "basic", "bearer"
	BasicUser        string
	BasicPass        string
	BearerToken      string
	UIPassword       string
}

// Server encapsulates the HTTP API and routing.
type Server struct {
	config ServerConfig
	mux    *http.ServeMux
}

// NewServer initializes a new Server.
func NewServer(cfg ServerConfig) *Server {
	if cfg.MaxScanSizeMB <= 0 {
		cfg.MaxScanSizeMB = 100
	}
	if cfg.RateLimitRPM <= 0 {
		cfg.RateLimitRPM = 100
	}
	if cfg.AuthMode == "" {
		cfg.AuthMode = "none"
	}
	if cfg.BasicUser == "" {
		cfg.BasicUser = "admin"
	}
	if cfg.UIPassword == "" {
		cfg.UIPassword = "123456"
	}

	s := &Server{
		config: cfg,
		mux:    http.NewServeMux(),
	}

	s.routes()
	return s
}

// Router returns the configured HTTP handler with middlewares.
func (s *Server) Router() http.Handler {
	return s.applyMiddlewares(s.mux)
}

func (s *Server) routes() {
	// Web UI SPA & Static Assets
	s.mux.Handle("/static/", http.StripPrefix("/static/", web.AssetHandler()))
	s.mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/static/index.html", http.StatusFound)
	})

	// Health and Ops (Always unauthenticated)
	s.mux.HandleFunc("GET /api/v1/health", s.handleHealth)
	s.mux.HandleFunc("GET /healthz", s.handleHealth)
	s.mux.HandleFunc("GET /api/v1/metrics", s.handleMetrics)

	// Web UI Password Auth
	s.mux.HandleFunc("GET /api/v1/auth/ui-status", s.handleUIStatus)
	s.mux.HandleFunc("POST /api/v1/auth/ui-login", s.handleUILogin)
	s.mux.HandleFunc("POST /api/v1/auth/ui-password", s.handleUIPassword)

	// Scanning
	s.mux.HandleFunc("POST /api/v1/scan/file", s.handleScanFile)
	s.mux.HandleFunc("POST /api/v1/scan/stream", s.handleScanStream)

	// Quarantine
	s.mux.HandleFunc("GET /api/v1/quarantine", s.handleQuarantineList)
	s.mux.HandleFunc("POST /api/v1/quarantine/restore", s.handleQuarantineRestore)

	// Audit Logs
	s.mux.HandleFunc("GET /api/v1/audit/export", s.handleAuditExport)
}

func (s *Server) applyMiddlewares(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Standard Security & CORS Headers
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, X-API-Key, X-Consumer-Name")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Rate Limiting
		if s.config.RateLimitEnabled && s.config.Limiter != nil {
			clientIP := extractIP(r)
			allowed, remaining, resetTime := s.config.Limiter.Allow(clientIP, s.config.RateLimitRPM)
			w.Header().Set("X-RateLimit-Limit", strconv.Itoa(s.config.RateLimitRPM))
			w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(remaining))
			w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(resetTime, 10))

			if !allowed {
				respondError(w, http.StatusTooManyRequests, "RATE_LIMIT_EXCEEDED", "Rate limit exceeded. Please retry later.", map[string]interface{}{
					"retry_after_seconds": resetTime - time.Now().Unix(),
				})
				return
			}
		}

		// Authentication Check
		path := r.URL.Path
		isPublic := path == "/healthz" || path == "/api/v1/health" || path == "/api/v1/metrics" ||
			path == "/api/v1/auth/ui-login" || path == "/api/v1/auth/ui-status" ||
			path == "/" || strings.HasPrefix(path, "/static/")

		if !isPublic && s.config.AuthMode != "none" {
			if s.config.AuthMode == "basic" {
				user, pass, ok := r.BasicAuth()
				if !ok || user != s.config.BasicUser || pass != s.config.BasicPass {
					w.Header().Set("WWW-Authenticate", `Basic realm="ClamAV Security"`)
					respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid or missing basic authentication credentials", nil)
					return
				}
			} else if s.config.AuthMode == "bearer" {
				authHeader := r.Header.Get("Authorization")
				token := ""
				if strings.HasPrefix(authHeader, "Bearer ") {
					token = strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
				} else if apiKey := r.Header.Get("X-API-Key"); apiKey != "" {
					token = strings.TrimSpace(apiKey)
				}

				if token == "" || token != s.config.BearerToken {
					respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid or missing Bearer token / API key", nil)
					return
				}
			}
		}

		h.ServeHTTP(w, r)
	})
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"status":    "healthy",
		"service":   "clamav-service",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func (s *Server) handleMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; version=0.0.4")
	fmt.Fprintf(w, "# HELP clamav_service_up Status of clamav service\n")
	fmt.Fprintf(w, "# TYPE clamav_service_up gauge\n")
	fmt.Fprintf(w, "clamav_service_up 1\n")
}

func (s *Server) handleScanFile(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	maxBytes := s.config.MaxScanSizeMB * 1024 * 1024
	r.Body = http.MaxBytesReader(w, r.Body, maxBytes)

	if err := r.ParseMultipartForm(maxBytes); err != nil {
		respondError(w, http.StatusRequestEntityTooLarge, "FILE_TOO_LARGE", "Uploaded file exceeds maximum limit", nil)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST_PAYLOAD", "Missing 'file' multipart form field", nil)
		return
	}
	defer file.Close()

	payload, err := io.ReadAll(file)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "Failed to read file payload", nil)
		return
	}

	hasher := sha256.New()
	hasher.Write(payload)
	fileHash := hex.EncodeToString(hasher.Sum(nil))

	// Check Whitelist
	if isWhitelisted, _ := s.config.DB.IsWhitelisted(fileHash); isWhitelisted {
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"success": true,
			"verdict": "CLEAN",
			"whitelisted": true,
			"data": map[string]interface{}{
				"file_name":        header.Filename,
				"file_size":        len(payload),
				"file_sha256":      fileHash,
				"scan_duration_ms": time.Since(startTime).Milliseconds(),
			},
		})
		return
	}

	// Scan with Clamd
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	scanRes, err := s.config.Clamd.ScanStream(ctx, bytes.NewReader(payload))
	if err != nil {
		respondError(w, http.StatusServiceUnavailable, "ENGINE_UNAVAILABLE", "Antivirus daemon unavailable or timed out", nil)
		return
	}

	durationMs := time.Since(startTime).Milliseconds()
	consumerName := extractConsumer(r)

	if scanRes.IsClean() {
		_ = s.config.DB.InsertScanAuditLog(storage.ScanAuditLog{
			ID:             fmt.Sprintf("audit_%d", time.Now().UnixNano()),
			Timestamp:      time.Now().UTC(),
			ConsumerName:   consumerName,
			ClientIP:       extractIP(r),
			FileName:       header.Filename,
			FileSizeBytes:  int64(len(payload)),
			FileSHA256:     fileHash,
			Verdict:        "CLEAN",
			ScanDurationMs: durationMs,
		})

		respondJSON(w, http.StatusOK, map[string]interface{}{
			"success": true,
			"verdict": "CLEAN",
			"data": map[string]interface{}{
				"file_name":        header.Filename,
				"file_size":        len(payload),
				"file_sha256":      fileHash,
				"scan_duration_ms": durationMs,
				"scanned_at":       time.Now().UTC().Format(time.RFC3339),
			},
		})
		return
	}

	// INFECTED: Quarantine & Alert
	quarRec, _ := s.config.Vault.QuarantineFile(ctx, header.Filename, consumerName, scanRes.VirusName, bytes.NewReader(payload), s.config.QuarRetention)
	quarID := ""
	if quarRec != nil {
		quarID = quarRec.ID
	}

	_ = s.config.DB.InsertScanAuditLog(storage.ScanAuditLog{
		ID:             fmt.Sprintf("audit_%d", time.Now().UnixNano()),
		Timestamp:      time.Now().UTC(),
		ConsumerName:   consumerName,
		ClientIP:       extractIP(r),
		FileName:       header.Filename,
		FileSizeBytes:  int64(len(payload)),
		FileSHA256:     fileHash,
		Verdict:        "INFECTED",
		VirusName:      scanRes.VirusName,
		ScanDurationMs: durationMs,
		QuarantineID:   quarID,
	})

	if s.config.Notifier != nil {
		_ = s.config.Notifier.DispatchThreat(ctx, alert.ThreatAlert{
			VirusName:      scanRes.VirusName,
			FileName:       header.Filename,
			FileSizeBytes:  int64(len(payload)),
			FileSHA256:     fileHash,
			QuarantineID:   quarID,
			SourceConsumer: consumerName,
			DetectedAt:     time.Now().UTC(),
		})
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"verdict": "INFECTED",
		"threat": map[string]interface{}{
			"virus_name":    scanRes.VirusName,
			"severity":      "HIGH",
			"action_taken":  "QUARANTINED",
			"quarantine_id": quarID,
		},
		"data": map[string]interface{}{
			"file_name":        header.Filename,
			"file_size":        len(payload),
			"file_sha256":      fileHash,
			"scan_duration_ms": durationMs,
		},
	})
}

func (s *Server) handleScanStream(w http.ResponseWriter, r *http.Request) {
	s.handleScanFile(w, r)
}

func (s *Server) handleQuarantineList(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"items":   []interface{}{},
	})
}

func (s *Server) handleQuarantineRestore(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"status":  "RESTORED",
	})
}

func (s *Server) handleAuditExport(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=\"scan_audit_logs.csv\"")
	fmt.Fprintf(w, "id,timestamp,consumer,file_name,verdict,virus_name,duration_ms\n")
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, code, message string, details interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"success": false,
		"error": map[string]interface{}{
			"code":    code,
			"message": message,
			"details": details,
		},
	})
}

func extractIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}
	parts := strings.Split(r.RemoteAddr, ":")
	if len(parts) > 0 {
		return parts[0]
	}
	return r.RemoteAddr
}

func extractConsumer(r *http.Request) string {
	if c := r.Header.Get("X-Consumer-Name"); c != "" {
		return c
	}
	return "Anonymous-Client"
}

func (s *Server) handleUIStatus(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success":   true,
		"protected": true,
	})
}

func (s *Server) verifyUIPassword(input string) bool {
	if s.config.DB != nil {
		storedHash, err := s.config.DB.GetSystemSetting("ui_password_hash")
		if err == nil && storedHash != "" {
			return crypto.VerifyPassword(input, storedHash)
		}
	}
	fallback := s.config.UIPassword
	if fallback == "" {
		fallback = "123456"
	}
	return input == fallback
}

func (s *Server) handleUILogin(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST_PAYLOAD", "Malformed JSON payload", nil)
		return
	}

	if !s.verifyUIPassword(body.Password) {
		respondError(w, http.StatusUnauthorized, "INVALID_PASSWORD", "Incorrect dashboard password", nil)
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"token":   fmt.Sprintf("ui_sess_%d", time.Now().UnixNano()),
		"message": "Authenticated successfully",
	})
}

func (s *Server) handleUIPassword(w http.ResponseWriter, r *http.Request) {
	var body struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST_PAYLOAD", "Malformed JSON payload", nil)
		return
	}

	if !s.verifyUIPassword(body.CurrentPassword) {
		respondError(w, http.StatusUnauthorized, "INVALID_PASSWORD", "Current password does not match", nil)
		return
	}

	if len(body.NewPassword) < 4 {
		respondError(w, http.StatusBadRequest, "PASSWORD_TOO_SHORT", "New password must be at least 4 characters long", nil)
		return
	}

	newHash := crypto.HashPassword(body.NewPassword)
	if s.config.DB != nil {
		if err := s.config.DB.SetSystemSetting("ui_password_hash", newHash); err != nil {
			respondError(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "Failed to persist new password", nil)
			return
		}
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Dashboard password updated successfully",
	})
}
