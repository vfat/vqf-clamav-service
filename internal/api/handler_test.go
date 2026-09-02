package api

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/vfat/vqf-clamav-service/internal/alert"
	"github.com/vfat/vqf-clamav-service/internal/clamd"
	"github.com/vfat/vqf-clamav-service/internal/quarantine"
	"github.com/vfat/vqf-clamav-service/internal/ratelimit"
	"github.com/vfat/vqf-clamav-service/internal/storage"
)

func setupTestServer(t *testing.T) (*Server, *storage.DB) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test_api.db")
	vaultDir := filepath.Join(tmpDir, "vault")

	db, err := storage.NewDB(dbPath)
	if err != nil {
		t.Fatalf("failed to init test DB: %v", err)
	}

	vault := quarantine.NewVault(vaultDir, db)
	notifier := alert.NewNotifier(alert.Config{})
	limiter := ratelimit.NewLimiter()

	// In unit test, mock clamd client pointing to non-existent or mock address
	clamdClient := clamd.NewClient("unix", "/tmp/mock.sock")

	server := NewServer(ServerConfig{
		DB:             db,
		Vault:          vault,
		Notifier:       notifier,
		Limiter:        limiter,
		Clamd:          clamdClient,
		RequireAPIKey:  false,
		MaxScanSizeMB:  100,
		LogRetention:   3,
		QuarRetention:  7,
	})

	return server, db
}

func TestHandler_HealthCheck(t *testing.T) {
	server, db := setupTestServer(t)
	defer db.Close()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	w := httptest.NewRecorder()

	server.Router().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse JSON response: %v", err)
	}

	if resp["status"] != "healthy" {
		t.Errorf("expected status 'healthy', got '%v'", resp["status"])
	}
}

func TestHandler_RateLimitingMiddleware(t *testing.T) {
	server, db := setupTestServer(t)
	defer db.Close()

	// Hit /api/v1/health 7 times rapidly with limited RPM = 5
	server.config.RateLimitRPM = 5
	server.config.RateLimitEnabled = true

	var lastCode int
	for i := 0; i < 7; i++ {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
		req.RemoteAddr = "192.168.1.100:1234"
		w := httptest.NewRecorder()
		server.Router().ServeHTTP(w, req)
		lastCode = w.Code
	}

	if lastCode != http.StatusTooManyRequests {
		t.Fatalf("expected status 429 Too Many Requests on exceeded rate limit, got %d", lastCode)
	}
}

func TestHandler_ScanMultipartClean(t *testing.T) {
	server, db := setupTestServer(t)
	defer db.Close()

	// Buat multipart form data
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "clean_document.pdf")
	if err != nil {
		t.Fatalf("failed to create form file: %v", err)
	}
	part.Write([]byte("Clean PDF document content without viruses"))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/scan/file", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	server.Router().ServeHTTP(w, req)

	// Since socket is mock in unit test, it should handle gracefully or return verdict
	if w.Code != http.StatusOK && w.Code != http.StatusServiceUnavailable {
		t.Errorf("expected status 200 or 503 (mock socket), got %d. Body: %s", w.Code, w.Body.String())
	}
}
