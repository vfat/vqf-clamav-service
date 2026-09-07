package api

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/vfat/vqf-clamav-service/internal/alert"
	"github.com/vfat/vqf-clamav-service/internal/clamd"
	"github.com/vfat/vqf-clamav-service/internal/crypto"
	"github.com/vfat/vqf-clamav-service/internal/fetcher"
	"github.com/vfat/vqf-clamav-service/internal/quarantine"
	"github.com/vfat/vqf-clamav-service/internal/ratelimit"
	"github.com/vfat/vqf-clamav-service/internal/asyncscan"
	"github.com/vfat/vqf-clamav-service/internal/storage"
	"github.com/vfat/vqf-clamav-service/internal/yara"
	"github.com/vfat/vqf-clamav-service/web"
)

// ServerConfig holds dependencies and configuration for the HTTP server.
type ServerConfig struct {
	DB               *storage.DB
	Vault            *quarantine.Vault
	Notifier         *alert.Notifier
	Limiter          *ratelimit.Limiter
	Clamd            *clamd.Client
	Fetcher          *fetcher.SafeFetcher
	YARAManager      *yara.Manager
	AsyncQueue       *asyncscan.Queue
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
	if cfg.Fetcher == nil {
		cfg.Fetcher = fetcher.NewSafeFetcher(fetcher.SafeFetcherConfig{
			Timeout:  30 * time.Second,
			MaxBytes: cfg.MaxScanSizeMB * 1024 * 1024,
		})
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
	s.mux.HandleFunc("GET /api/v1/stats", s.handleStats)

	// Web UI Password Auth
	s.mux.HandleFunc("GET /api/v1/auth/ui-status", s.handleUIStatus)
	s.mux.HandleFunc("POST /api/v1/auth/ui-login", s.handleUILogin)
	s.mux.HandleFunc("POST /api/v1/auth/ui-password", s.handleUIPassword)

	// Scanning
	s.mux.HandleFunc("POST /api/v1/scan/file", s.handleScanFile)
	s.mux.HandleFunc("POST /api/v1/scan/stream", s.handleScanStream)
	s.mux.HandleFunc("POST /api/v1/scan/url", s.handleScanURL)
	s.mux.HandleFunc("POST /api/v1/scan/async", s.handleScanAsync)
	s.mux.HandleFunc("GET /api/v1/scan/jobs/{id}", s.handleScanJobStatus)
	s.mux.HandleFunc("GET /api/v1/scan/jobs", s.handleScanJobList)

	// Quarantine
	s.mux.HandleFunc("GET /api/v1/quarantine", s.handleQuarantineList)
	s.mux.HandleFunc("GET /api/v1/quarantine/download", s.handleQuarantineDownload)
	s.mux.HandleFunc("POST /api/v1/quarantine/restore", s.handleQuarantineRestore)
	s.mux.HandleFunc("POST /api/v1/quarantine/delete", s.handleQuarantineDelete)
	s.mux.HandleFunc("DELETE /api/v1/quarantine", s.handleQuarantineDelete)

	// Audit Logs
	s.mux.HandleFunc("GET /api/v1/audit/export", s.handleAuditExport)

	// YARA Custom Rules Management
	s.mux.HandleFunc("POST /api/v1/rules/yara", s.handleYARAAdd)
	s.mux.HandleFunc("GET /api/v1/rules/yara", s.handleYARAList)
	s.mux.HandleFunc("DELETE /api/v1/rules/yara/{id}", s.handleYARADelete)
	s.mux.HandleFunc("DELETE /api/v1/rules/yara", s.handleYARADelete)
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
			path == "/api/v1/stats" || path == "/api/v1/auth/ui-login" || path == "/api/v1/auth/ui-status" ||
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

type scanURLRequest struct {
	URL            string `json:"url"`
	ExpectedSHA256 string `json:"expected_sha256"`
}

func (s *Server) handleScanURL(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	if r.Body == nil {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST_PAYLOAD", "Missing request body", nil)
		return
	}

	var reqBody scanURLRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST_PAYLOAD", "Invalid JSON request body", nil)
		return
	}

	reqBody.URL = strings.TrimSpace(reqBody.URL)
	if reqBody.URL == "" {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST_PAYLOAD", "Field 'url' is required", nil)
		return
	}

	if len(reqBody.URL) > 2048 {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST_PAYLOAD", "Field 'url' exceeds maximum length of 2048 characters", nil)
		return
	}

	if err := fetcher.ValidateURL(reqBody.URL); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST_PAYLOAD", err.Error(), nil)
		return
	}

	// Fetch remote target with Anti-SSRF protection
	fetchCtx, cancel := context.WithTimeout(r.Context(), 35*time.Second)
	defer cancel()

	fetchRes, err := s.config.Fetcher.Fetch(fetchCtx, reqBody.URL)
	if err != nil {
		if errors.Is(err, fetcher.ErrSSRFBlocked) {
			respondError(w, http.StatusBadRequest, "SSRF_ATTEMPT_BLOCKED", "The provided target URL resolves to a prohibited private, loopback, or cloud-metadata IP address", map[string]interface{}{
				"target_url": reqBody.URL,
				"policy":     "RFC1918, Loopback, and Link-Local IPs are strictly forbidden",
			})
			return
		}
		if errors.Is(err, fetcher.ErrInvalidScheme) {
			respondError(w, http.StatusBadRequest, "INVALID_REQUEST_PAYLOAD", err.Error(), nil)
			return
		}
		if errors.Is(err, fetcher.ErrTooManyRedirects) {
			respondError(w, http.StatusBadRequest, "TOO_MANY_REDIRECTS", "Redirect limit exceeded", nil)
			return
		}
		if errors.Is(err, context.DeadlineExceeded) {
			respondError(w, http.StatusGatewayTimeout, "DOWNLOAD_TIMEOUT", "Failed to stream remote file within the allowed window", nil)
			return
		}
		respondError(w, http.StatusBadGateway, "DOWNLOAD_FAILED", fmt.Sprintf("Failed to download remote file: %v", err), nil)
		return
	}
	defer fetchRes.Body.Close()

	if fetchRes.StatusCode < 200 || fetchRes.StatusCode >= 300 {
		respondError(w, http.StatusBadGateway, "DOWNLOAD_FAILED", fmt.Sprintf("Remote server responded with status code %d", fetchRes.StatusCode), nil)
		return
	}

	payload, err := io.ReadAll(fetchRes.Body)
	if err != nil {
		if errors.Is(err, fetcher.ErrPayloadTooLarge) {
			respondError(w, http.StatusRequestEntityTooLarge, "FILE_TOO_LARGE", fmt.Sprintf("Remote file exceeds the configured maximum scanning size limit (%d MB)", s.config.MaxScanSizeMB), nil)
			return
		}
		respondError(w, http.StatusBadGateway, "DOWNLOAD_FAILED", fmt.Sprintf("Failed reading remote file stream: %v", err), nil)
		return
	}

	hasher := sha256.New()
	hasher.Write(payload)
	fileHash := hex.EncodeToString(hasher.Sum(nil))

	// Verify expected SHA256 checksum if provided
	if reqBody.ExpectedSHA256 != "" && !strings.EqualFold(reqBody.ExpectedSHA256, fileHash) {
		respondError(w, http.StatusBadRequest, "CHECKSUM_MISMATCH", fmt.Sprintf("Downloaded file SHA256 (%s) does not match expected checksum (%s)", fileHash, reqBody.ExpectedSHA256), map[string]interface{}{
			"computed_sha256": fileHash,
			"expected_sha256": reqBody.ExpectedSHA256,
		})
		return
	}

	// Determine file name from URL
	fileName := "remote_file"
	if parsed, err := url.Parse(reqBody.URL); err == nil {
		base := filepath.Base(parsed.Path)
		if base != "" && base != "." && base != "/" {
			fileName = base
		}
	}

	// Check Whitelist
	if isWhitelisted, _ := s.config.DB.IsWhitelisted(fileHash); isWhitelisted {
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"success":     true,
			"verdict":     "CLEAN",
			"whitelisted": true,
			"data": map[string]interface{}{
				"source_url":       reqBody.URL,
				"file_name":        fileName,
				"file_size":        len(payload),
				"file_sha256":      fileHash,
				"scan_duration_ms": time.Since(startTime).Milliseconds(),
				"scanned_at":       time.Now().UTC().Format(time.RFC3339),
			},
		})
		return
	}

	// Scan with Clamd daemon
	scanCtx, scanCancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer scanCancel()

	scanRes, err := s.config.Clamd.ScanStream(scanCtx, bytes.NewReader(payload))
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
			FileName:       fileName,
			FileSizeBytes:  int64(len(payload)),
			FileSHA256:     fileHash,
			Verdict:        "CLEAN",
			ScanDurationMs: durationMs,
		})

		respondJSON(w, http.StatusOK, map[string]interface{}{
			"success": true,
			"verdict": "CLEAN",
			"data": map[string]interface{}{
				"source_url":       reqBody.URL,
				"file_name":        fileName,
				"file_size":        len(payload),
				"file_sha256":      fileHash,
				"scan_duration_ms": durationMs,
				"scanned_at":       time.Now().UTC().Format(time.RFC3339),
			},
		})
		return
	}

	// INFECTED: Quarantine & Alert
	quarRec, _ := s.config.Vault.QuarantineFile(r.Context(), fileName, consumerName, scanRes.VirusName, bytes.NewReader(payload), s.config.QuarRetention)
	quarID := ""
	if quarRec != nil {
		quarID = quarRec.ID
	}

	_ = s.config.DB.InsertScanAuditLog(storage.ScanAuditLog{
		ID:             fmt.Sprintf("audit_%d", time.Now().UnixNano()),
		Timestamp:      time.Now().UTC(),
		ConsumerName:   consumerName,
		ClientIP:       extractIP(r),
		FileName:       fileName,
		FileSizeBytes:  int64(len(payload)),
		FileSHA256:     fileHash,
		Verdict:        "INFECTED",
		VirusName:      scanRes.VirusName,
		ScanDurationMs: durationMs,
		QuarantineID:   quarID,
	})

	if s.config.Notifier != nil {
		_ = s.config.Notifier.DispatchThreat(r.Context(), alert.ThreatAlert{
			VirusName:      scanRes.VirusName,
			FileName:       fileName,
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
			"source_url":       reqBody.URL,
			"file_name":        fileName,
			"file_size":        len(payload),
			"file_sha256":      fileHash,
			"scan_duration_ms": durationMs,
		},
	})
}

func (s *Server) handleQuarantineList(w http.ResponseWriter, r *http.Request) {
	limit := 50
	offset := 0
	status := r.URL.Query().Get("status")
	if l := r.URL.Query().Get("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil && val > 0 {
			limit = val
		}
	}
	if o := r.URL.Query().Get("offset"); o != "" {
		if val, err := strconv.Atoi(o); err == nil && val >= 0 {
			offset = val
		}
	}

	records, total, err := s.config.DB.ListQuarantineRecords(limit, offset, status)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "Failed to retrieve quarantine records", nil)
		return
	}

	type quarantineItem struct {
		ID            string `json:"id"`
		FileName      string `json:"file_name"`
		FileSize      int64  `json:"file_size"`
		FileSHA256    string `json:"file_sha256"`
		VirusName     string `json:"virus_name"`
		ConsumerName  string `json:"consumer_name"`
		Status        string `json:"status"`
		QuarantinedAt string `json:"quarantined_at"`
		ExpiresAt     string `json:"expires_at"`
	}

	items := make([]quarantineItem, 0)
	for _, rec := range records {
		items = append(items, quarantineItem{
			ID:            rec.ID,
			FileName:      rec.OriginalFilename,
			FileSize:      rec.FileSizeBytes,
			FileSHA256:    rec.FileSHA256,
			VirusName:     rec.VirusName,
			ConsumerName:  rec.SourceConsumer,
			Status:        rec.Status,
			QuarantinedAt: rec.CreatedAt.Format(time.RFC3339),
			ExpiresAt:     rec.ExpiresAt.Format(time.RFC3339),
		})
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"total":   total,
		"items":   items,
	})
}

func (s *Server) handleQuarantineRestore(w http.ResponseWriter, r *http.Request) {
	var body struct {
		QuarantineID  string `json:"quarantine_id"`
		RestoredBy    string `json:"restored_by"`
		Reason        string `json:"reason"`
		AutoWhitelist bool   `json:"auto_whitelist"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST_PAYLOAD", "Malformed JSON payload", nil)
		return
	}

	if body.QuarantineID == "" {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST_PAYLOAD", "Missing quarantine_id parameter", nil)
		return
	}

	if body.RestoredBy == "" {
		body.RestoredBy = "admin"
	}
	if body.Reason == "" {
		body.Reason = "Restored via Web Admin UI"
	}

	rec, err := s.config.DB.GetQuarantineRecord(body.QuarantineID)
	if err != nil {
		respondError(w, http.StatusNotFound, "NOT_FOUND", "Quarantine record not found", nil)
		return
	}

	if err := s.config.DB.RestoreQuarantineRecord(body.QuarantineID, body.RestoredBy, body.Reason, body.AutoWhitelist); err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "Failed to restore quarantine record: "+err.Error(), nil)
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"status":  "RESTORED",
		"message": "File restored successfully and hash whitelisted.",
		"data": map[string]interface{}{
			"quarantine_id": rec.ID,
			"status":        "RESTORED",
			"file_sha256":   rec.FileSHA256,
			"whitelisted":   body.AutoWhitelist,
			"restored_at":   time.Now().UTC().Format(time.RFC3339),
		},
	})
}

func (s *Server) handleQuarantineDelete(w http.ResponseWriter, r *http.Request) {
	var body struct {
		QuarantineID string `json:"quarantine_id"`
	}

	if id := r.URL.Query().Get("id"); id != "" {
		body.QuarantineID = id
	} else if r.Body != nil {
		_ = json.NewDecoder(r.Body).Decode(&body)
	}

	if body.QuarantineID == "" {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST_PAYLOAD", "Missing quarantine_id parameter", nil)
		return
	}

	if err := s.config.Vault.DeleteFile(r.Context(), body.QuarantineID); err != nil {
		respondError(w, http.StatusNotFound, "NOT_FOUND", "Quarantine item not found or deletion failed: "+err.Error(), nil)
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success":       true,
		"status":        "DELETED",
		"message":       "Quarantined malware payload permanently destroyed.",
		"quarantine_id": body.QuarantineID,
	})
}

func (s *Server) handleQuarantineDownload(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST_PAYLOAD", "Missing id parameter", nil)
		return
	}

	reader, record, err := s.config.Vault.DownloadFile(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, "NOT_FOUND", "Quarantine item not found: "+err.Error(), nil)
		return
	}
	defer reader.Close()

	filename := record.OriginalFilename
	if filename == "" {
		filename = record.ID + ".bin"
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	w.Header().Set("Content-Length", strconv.FormatInt(record.FileSizeBytes, 10))
	w.WriteHeader(http.StatusOK)

	_, _ = io.Copy(w, reader)
}

func (s *Server) handleStats(w http.ResponseWriter, r *http.Request) {
	stats, err := s.config.DB.GetSystemStats()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "Failed to fetch stats", nil)
		return
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    stats,
	})
}

func (s *Server) handleAuditExport(w http.ResponseWriter, r *http.Request) {
	format := r.URL.Query().Get("format")
	logs, _, err := s.config.DB.ListScanAuditLogs(storage.AuditFilter{Limit: 100})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "Failed to query audit logs", nil)
		return
	}

	if format == "json" || strings.Contains(r.Header.Get("Accept"), "application/json") {
		respondJSON(w, http.StatusOK, map[string]interface{}{
			"success": true,
			"items":   logs,
		})
		return
	}

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=\"scan_audit_logs.csv\"")
	fmt.Fprintf(w, "id,timestamp,consumer,file_name,verdict,virus_name,duration_ms\n")
	for _, l := range logs {
		fmt.Fprintf(w, "%s,%s,%s,%s,%s,%s,%d\n",
			l.ID, l.Timestamp.Format(time.RFC3339), l.ConsumerName, l.FileName,
			l.Verdict, l.VirusName, l.ScanDurationMs,
		)
	}
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

func (s *Server) handleYARAAdd(w http.ResponseWriter, r *http.Request) {
	if s.config.YARAManager == nil {
		respondError(w, http.StatusServiceUnavailable, "YARA_NOT_AVAILABLE", "YARA rule injection engine is not configured", nil)
		return
	}

	var req struct {
		RuleName    string `json:"rule_name"`
		Description string `json:"description"`
		Content     string `json:"content"`
		Author      string `json:"author"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST_PAYLOAD", "Malformed JSON payload", nil)
		return
	}

	req.RuleName = strings.TrimSpace(req.RuleName)
	if req.RuleName == "" {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST_PAYLOAD", "Field 'rule_name' is required", nil)
		return
	}

	req.Content = strings.TrimSpace(req.Content)
	if req.Content == "" {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST_PAYLOAD", "Field 'content' is required", nil)
		return
	}

	if req.Author == "" {
		req.Author = extractConsumer(r)
	}

	rule, err := s.config.YARAManager.AddRule(r.Context(), req.RuleName, req.Description, req.Content, req.Author)
	if err != nil {
		errStr := err.Error()
		if strings.Contains(errStr, "yara validation failed") || strings.Contains(errStr, "already exists") {
			respondError(w, http.StatusBadRequest, "INVALID_YARA_RULE", errStr, nil)
			return
		}
		if strings.Contains(errStr, "failed to reload clamd") {
			respondError(w, http.StatusInternalServerError, "RELOAD_FAILED", "Failed to reload ClamAV with new rule: "+errStr, nil)
			return
		}
		respondError(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "Failed to save YARA rule: "+errStr, nil)
		return
	}

	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"message": "YARA rule deployed and ClamAV reloaded successfully",
		"data":    rule,
	})
}

func (s *Server) handleYARAList(w http.ResponseWriter, r *http.Request) {
	if s.config.YARAManager == nil {
		respondError(w, http.StatusServiceUnavailable, "YARA_NOT_AVAILABLE", "YARA rule injection engine is not configured", nil)
		return
	}

	rules, err := s.config.YARAManager.ListRules(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "Failed to list YARA rules: "+err.Error(), nil)
		return
	}

	if rules == nil {
		rules = []storage.YARARule{}
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"total":   len(rules),
		"items":   rules,
	})
}

func (s *Server) handleYARADelete(w http.ResponseWriter, r *http.Request) {
	if s.config.YARAManager == nil {
		respondError(w, http.StatusServiceUnavailable, "YARA_NOT_AVAILABLE", "YARA rule injection engine is not configured", nil)
		return
	}

	id := r.PathValue("id")
	if id == "" {
		id = r.URL.Query().Get("id")
	}
	if id == "" && r.Body != nil {
		var body struct {
			ID string `json:"id"`
		}
		_ = json.NewDecoder(r.Body).Decode(&body)
		id = body.ID
	}

	if id == "" {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST_PAYLOAD", "Missing rule id parameter", nil)
		return
	}

	if err := s.config.YARAManager.DeleteRule(r.Context(), id); err != nil {
		if strings.Contains(err.Error(), "not found") {
			respondError(w, http.StatusNotFound, "RULE_NOT_FOUND", fmt.Sprintf("YARA rule with id '%s' not found", id), nil)
			return
		}
		respondError(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "Failed to delete YARA rule: "+err.Error(), nil)
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "YARA rule removed and ClamAV reloaded successfully",
		"rule_id": id,
	})
}

func (s *Server) handleScanAsync(w http.ResponseWriter, r *http.Request) {
	if s.config.AsyncQueue == nil {
		respondError(w, http.StatusServiceUnavailable, "ASYNC_SCAN_UNAVAILABLE", "Async scan queue engine is not configured", nil)
		return
	}

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

	callbackURL := strings.TrimSpace(r.FormValue("callback_url"))
	if callbackURL == "" {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST_PAYLOAD", "Missing 'callback_url' form field", nil)
		return
	}

	consumerName := extractConsumer(r)

	job, err := s.config.AsyncQueue.SubmitJob(r.Context(), header.Filename, file, header.Size, callbackURL, consumerName)
	if err != nil {
		errStr := err.Error()
		if strings.Contains(errStr, "invalid callback_url") || strings.Contains(errStr, "prohibited address") || strings.Contains(errStr, "SSRF") {
			respondError(w, http.StatusBadRequest, "SSRF_ATTEMPT_BLOCKED", errStr, nil)
			return
		}
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST_PAYLOAD", errStr, nil)
		return
	}

	respondJSON(w, http.StatusAccepted, map[string]interface{}{
		"success": true,
		"status":  "ACCEPTED",
		"job_id":  job.ID,
		"message": "Scan job queued. Verdict will be posted to callback_url.",
		"data": map[string]interface{}{
			"job_id":       job.ID,
			"file_name":    job.FileName,
			"file_size":    job.FileSize,
			"callback_url": job.CallbackURL,
			"status":       job.Status,
			"created_at":   job.CreatedAt.Format(time.RFC3339),
		},
	})
}

func (s *Server) handleScanJobStatus(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		id = r.URL.Query().Get("id")
	}

	if id == "" {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST_PAYLOAD", "Missing job id parameter", nil)
		return
	}

	job, err := s.config.DB.GetScanJob(id)
	if err != nil {
		respondError(w, http.StatusNotFound, "JOB_NOT_FOUND", fmt.Sprintf("Scan job with id '%s' not found", id), nil)
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"job_id":       job.ID,
			"file_name":    job.FileName,
			"file_size":    job.FileSize,
			"file_sha256":  job.FileSHA256,
			"callback_url": job.CallbackURL,
			"consumer":     job.Consumer,
			"status":       job.Status,
			"verdict":      job.Verdict,
			"virus_name":   job.VirusName,
			"error_msg":    job.ErrorMsg,
			"created_at":   job.CreatedAt.Format(time.RFC3339),
			"updated_at":   job.UpdatedAt.Format(time.RFC3339),
		},
	})
}

func (s *Server) handleScanJobList(w http.ResponseWriter, r *http.Request) {
	if s.config.DB == nil {
		respondError(w, http.StatusServiceUnavailable, "STORAGE_UNAVAILABLE", "Database not configured", nil)
		return
	}

	limit := 50
	offset := 0
	if l := r.URL.Query().Get("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil && val > 0 {
			limit = val
		}
	}
	if o := r.URL.Query().Get("offset"); o != "" {
		if val, err := strconv.Atoi(o); err == nil && val >= 0 {
			offset = val
		}
	}

	jobs, total, err := s.config.DB.ListScanJobs(limit, offset)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "Failed to query scan jobs: "+err.Error(), nil)
		return
	}
	if jobs == nil {
		jobs = []storage.ScanJob{}
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"total":   total,
		"items":   jobs,
	})
}


