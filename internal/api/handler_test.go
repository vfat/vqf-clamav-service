package api

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/vfat/vqf-clamav-service/internal/alert"
	"github.com/vfat/vqf-clamav-service/internal/clamd"
	"github.com/vfat/vqf-clamav-service/internal/fetcher"
	"github.com/vfat/vqf-clamav-service/internal/quarantine"
	"github.com/vfat/vqf-clamav-service/internal/ratelimit"
	"github.com/vfat/vqf-clamav-service/internal/storage"
	"github.com/vfat/vqf-clamav-service/internal/yara"
	"github.com/vfat/vqf-clamav-service/internal/asyncscan"
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

func TestHandler_ServeStaticWebUI(t *testing.T) {
	server, db := setupTestServer(t)
	defer db.Close()

	req := httptest.NewRequest(http.MethodGet, "/static/app.css", nil)
	w := httptest.NewRecorder()

	server.Router().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200 for /static/app.css, got %d", w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "CLAMAV-SERVICE") {
		t.Errorf("expected embedded CSS to contain 'CLAMAV-SERVICE'")
	}
}

func TestHandler_ScanURL_SSRFBlocked(t *testing.T) {
	server, db := setupTestServer(t)
	defer db.Close()

	payload := map[string]string{
		"url": "http://169.254.169.254/latest/meta-data",
	}
	bodyBytes, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/scan/url", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	server.Router().ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400 Bad Request for SSRF target, got %d. Body: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	errObj, ok := resp["error"].(map[string]interface{})
	if !ok || errObj["code"] != "SSRF_ATTEMPT_BLOCKED" {
		t.Errorf("expected error code SSRF_ATTEMPT_BLOCKED, got: %v", resp)
	}
}

func TestHandler_ScanURL_InvalidPayload(t *testing.T) {
	server, db := setupTestServer(t)
	defer db.Close()

	// Invalid URL scheme
	payload := map[string]string{
		"url": "ftp://example.com/malware.exe",
	}
	bodyBytes, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/scan/url", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	server.Router().ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400 for ftp:// scheme, got %d", w.Code)
	}

	var resp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	errObj, _ := resp["error"].(map[string]interface{})
	if errObj["code"] != "INVALID_REQUEST_PAYLOAD" {
		t.Errorf("expected INVALID_REQUEST_PAYLOAD, got %v", errObj["code"])
	}
}

func TestHandler_ScanURL_CleanAndChecksum(t *testing.T) {
	tmpDir := t.TempDir()
	db, err := storage.NewDB(filepath.Join(tmpDir, "test.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	vault := quarantine.NewVault(filepath.Join(tmpDir, "vault"), db)
	notifier := alert.NewNotifier(alert.Config{})
	limiter := ratelimit.NewLimiter()
	clamdClient := clamd.NewClient("unix", "/tmp/mock.sock")

	fileContent := []byte("Safe content from remote S3 bucket")
	hasher := sha256.New()
	hasher.Write(fileContent)
	expectedHash := hex.EncodeToString(hasher.Sum(nil))

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(fileContent)
	}))
	defer mockServer.Close()

	safeFetcher := fetcher.NewSafeFetcher(fetcher.SafeFetcherConfig{
		Timeout:                 5 * time.Second,
		AllowLoopbackForTesting: true,
	})

	server := NewServer(ServerConfig{
		DB:            db,
		Vault:         vault,
		Notifier:      notifier,
		Limiter:       limiter,
		Clamd:         clamdClient,
		Fetcher:       safeFetcher,
		MaxScanSizeMB: 100,
	})

	// Test with matching checksum
	payload := map[string]string{
		"url":             mockServer.URL + "/report.pdf",
		"expected_sha256": expectedHash,
	}
	bodyBytes, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/scan/url", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	server.Router().ServeHTTP(w, req)

	// Mock clamd returns 503 or 200
	if w.Code != http.StatusOK && w.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected status 200 or 503, got %d. Body: %s", w.Code, w.Body.String())
	}

	// Test with MISMATCHING checksum -> must return 400 CHECKSUM_MISMATCH before reaching clamd
	mismatchPayload := map[string]string{
		"url":             mockServer.URL + "/report.pdf",
		"expected_sha256": "0000000000000000000000000000000000000000000000000000000000000000",
	}
	mismatchBytes, _ := json.Marshal(mismatchPayload)
	reqMismatch := httptest.NewRequest(http.MethodPost, "/api/v1/scan/url", bytes.NewReader(mismatchBytes))
	reqMismatch.Header.Set("Content-Type", "application/json")
	wMismatch := httptest.NewRecorder()

	server.Router().ServeHTTP(wMismatch, reqMismatch)

	if wMismatch.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400 for checksum mismatch, got %d. Body: %s", wMismatch.Code, wMismatch.Body.String())
	}

	var mismatchResp map[string]interface{}
	_ = json.Unmarshal(wMismatch.Body.Bytes(), &mismatchResp)
	errObj, _ := mismatchResp["error"].(map[string]interface{})
	if errObj["code"] != "CHECKSUM_MISMATCH" {
		t.Errorf("expected CHECKSUM_MISMATCH, got %v", errObj["code"])
	}
}

func TestHandler_ScanURL_FileTooLarge(t *testing.T) {
	tmpDir := t.TempDir()
	db, _ := storage.NewDB(filepath.Join(tmpDir, "test.db"))
	defer db.Close()

	vault := quarantine.NewVault(filepath.Join(tmpDir, "vault"), db)
	notifier := alert.NewNotifier(alert.Config{})
	limiter := ratelimit.NewLimiter()
	clamdClient := clamd.NewClient("unix", "/tmp/mock.sock")

	largeContent := make([]byte, 5000)
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(largeContent)
	}))
	defer mockServer.Close()

	safeFetcher := fetcher.NewSafeFetcher(fetcher.SafeFetcherConfig{
		Timeout:                 5 * time.Second,
		MaxBytes:                1024, // 1KB limit
		AllowLoopbackForTesting: true,
	})

	server := NewServer(ServerConfig{
		DB:            db,
		Vault:         vault,
		Notifier:      notifier,
		Limiter:       limiter,
		Clamd:         clamdClient,
		Fetcher:       safeFetcher,
		MaxScanSizeMB: 1,
	})

	payload := map[string]string{
		"url": mockServer.URL + "/large.bin",
	}
	bodyBytes, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/scan/url", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	server.Router().ServeHTTP(w, req)

	if w.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("expected status 413 Payload Too Large, got %d. Body: %s", w.Code, w.Body.String())
	}
}

func TestHandler_YARARules_CRUD(t *testing.T) {
	tmpDir := t.TempDir()
	db, err := storage.NewDB(filepath.Join(tmpDir, "test.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	rulesDir := filepath.Join(tmpDir, "yara_rules")
	mockClamd := clamd.NewClient("unix", "/tmp/nonexistent.sock")
	yaraMgr := yara.NewManager(rulesDir, db, nil)

	server := NewServer(ServerConfig{
		DB:          db,
		YARAManager: yaraMgr,
		Clamd:       mockClamd,
	})

	// 1. Add invalid YARA rule
	invalidPayload := map[string]string{
		"rule_name": "bad_rule",
		"content":   "rule bad_rule { invalid }",
	}
	bodyBytes, _ := json.Marshal(invalidPayload)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/rules/yara", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	server.Router().ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400 for invalid YARA rule, got %d. Body: %s", w.Code, w.Body.String())
	}

	// 2. Add valid YARA rule
	validPayload := map[string]string{
		"rule_name":   "custom_webshell",
		"description": "Detect webshell injection",
		"content":     "rule custom_webshell { condition: true }",
		"author":      "sec-admin",
	}
	bodyBytes, _ = json.Marshal(validPayload)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/rules/yara", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	server.Router().ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201 Created, got %d. Body: %s", w.Code, w.Body.String())
	}

	var createdResp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &createdResp)
	ruleData, ok := createdResp["data"].(map[string]interface{})
	if !ok || ruleData["id"] == "" {
		t.Fatalf("expected created rule data with id, got %v", createdResp)
	}
	ruleID := ruleData["id"].(string)

	// 3. List YARA rules
	req = httptest.NewRequest(http.MethodGet, "/api/v1/rules/yara", nil)
	w = httptest.NewRecorder()

	server.Router().ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200 OK, got %d", w.Code)
	}

	var listResp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &listResp)
	items, ok := listResp["items"].([]interface{})
	if !ok || len(items) != 1 {
		t.Fatalf("expected 1 item in list, got %v", listResp)
	}

	// 4. Delete YARA rule
	req = httptest.NewRequest(http.MethodDelete, "/api/v1/rules/yara/"+ruleID, nil)
	w = httptest.NewRecorder()

	server.Router().ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200 OK on delete, got %d. Body: %s", w.Code, w.Body.String())
	}

	// 5. Delete non-existent rule
	req = httptest.NewRequest(http.MethodDelete, "/api/v1/rules/yara/nonexistent-id", nil)
	w = httptest.NewRecorder()

	server.Router().ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status 404 for non-existent rule, got %d. Body: %s", w.Code, w.Body.String())
	}

	// 6. List after delete should be empty
	req = httptest.NewRequest(http.MethodGet, "/api/v1/rules/yara", nil)
	w = httptest.NewRecorder()

	server.Router().ServeHTTP(w, req)
	var emptyListResp map[string]interface{}
	_ = json.Unmarshal(w.Body.Bytes(), &emptyListResp)
	emptyItems, _ := emptyListResp["items"].([]interface{})
	if len(emptyItems) != 0 {
		t.Errorf("expected 0 items after delete, got %d", len(emptyItems))
	}
}

func TestHandler_ScanAsync_SubmitAndPoll(t *testing.T) {
	tmpDir := t.TempDir()
	db, err := storage.NewDB(filepath.Join(tmpDir, "test.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	spoolDir := filepath.Join(tmpDir, "spool")
	mockClamd := clamd.NewClient("unix", "/tmp/nonexistent.sock")
	safeFetcher := fetcher.NewSafeFetcher(fetcher.SafeFetcherConfig{
		Timeout:                 2 * time.Second,
		AllowLoopbackForTesting: true,
	})

	asyncQueue := asyncscan.NewQueue(asyncscan.Config{
		SpoolDir:      spoolDir,
		Workers:       1,
		QueueCapacity: 10,
	}, db, mockClamd, nil, nil, safeFetcher)

	server := NewServer(ServerConfig{
		DB:         db,
		AsyncQueue: asyncQueue,
		Fetcher:    safeFetcher,
	})

	// 1. Submit async scan job missing callback_url -> 400 Bad Request
	bodyMissing := &bytes.Buffer{}
	writerMissing := multipart.NewWriter(bodyMissing)
	part, _ := writerMissing.CreateFormFile("file", "test.zip")
	part.Write([]byte("dummy content"))
	writerMissing.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/scan/async", bodyMissing)
	req.Header.Set("Content-Type", writerMissing.FormDataContentType())
	w := httptest.NewRecorder()
	server.Router().ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400 for missing callback_url, got %d", w.Code)
	}

	mockWebhookServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer mockWebhookServer.Close()

	// 2. Submit valid async scan job -> 202 Accepted
	bodyValid := &bytes.Buffer{}
	writerValid := multipart.NewWriter(bodyValid)
	partValid, _ := writerValid.CreateFormFile("file", "archive.zip")
	partValid.Write([]byte("safe archive file bytes"))
	writerValid.WriteField("callback_url", mockWebhookServer.URL+"/webhook")
	writerValid.Close()


	reqValid := httptest.NewRequest(http.MethodPost, "/api/v1/scan/async", bodyValid)
	reqValid.Header.Set("Content-Type", writerValid.FormDataContentType())
	wValid := httptest.NewRecorder()
	server.Router().ServeHTTP(wValid, reqValid)

	if wValid.Code != http.StatusAccepted {
		t.Fatalf("expected status 202 Accepted, got %d. Body: %s", wValid.Code, wValid.Body.String())
	}

	var acceptedResp map[string]interface{}
	_ = json.Unmarshal(wValid.Body.Bytes(), &acceptedResp)
	jobID, _ := acceptedResp["job_id"].(string)
	if jobID == "" {
		t.Fatalf("expected job_id in response, got %v", acceptedResp)
	}

	// 3. Poll job status -> 200 OK
	reqPoll := httptest.NewRequest(http.MethodGet, "/api/v1/scan/jobs/"+jobID, nil)
	wPoll := httptest.NewRecorder()
	server.Router().ServeHTTP(wPoll, reqPoll)

	if wPoll.Code != http.StatusOK {
		t.Fatalf("expected status 200 OK polling job status, got %d", wPoll.Code)
	}

	var pollResp map[string]interface{}
	_ = json.Unmarshal(wPoll.Body.Bytes(), &pollResp)
	jobData, ok := pollResp["data"].(map[string]interface{})
	if !ok || jobData["job_id"] != jobID {
		t.Fatalf("unexpected poll data: %v", pollResp)
	}

	// 4. Poll non-existent job -> 404 Not Found
	reqNotFound := httptest.NewRequest(http.MethodGet, "/api/v1/scan/jobs/nonexistent_job", nil)
	wNotFound := httptest.NewRecorder()
	server.Router().ServeHTTP(wNotFound, reqNotFound)

	if wNotFound.Code != http.StatusNotFound {
		t.Fatalf("expected status 404 for non-existent job, got %d", wNotFound.Code)
	}
}




