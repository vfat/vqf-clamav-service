package api

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/vfat/vqf-clamav-service/internal/alert"
	"github.com/vfat/vqf-clamav-service/internal/clamd"
	"github.com/vfat/vqf-clamav-service/internal/quarantine"
	"github.com/vfat/vqf-clamav-service/internal/ratelimit"
	"github.com/vfat/vqf-clamav-service/internal/storage"
)

func createServerWithAuth(t *testing.T, authMode, basicUser, basicPass, bearerToken string) (*Server, *storage.DB) {
	tmpDir := t.TempDir()
	db, err := storage.NewDB(tmpDir + "/test_auth.db")
	if err != nil {
		t.Fatalf("failed to init db: %v", err)
	}

	server := NewServer(ServerConfig{
		DB:            db,
		Vault:         quarantine.NewVault(tmpDir+"/vault", db),
		Notifier:      alert.NewNotifier(alert.Config{}),
		Limiter:       ratelimit.NewLimiter(),
		Clamd:         clamd.NewClient("unix", "/tmp/mock.sock"),
		AuthMode:      authMode,
		BasicUser:     basicUser,
		BasicPass:     basicPass,
		BearerToken:   bearerToken,
		UIPassword:    "123456",
		MaxScanSizeMB: 100,
	})

	return server, db
}

func TestAuth_ModeNone(t *testing.T) {
	server, db := createServerWithAuth(t, "none", "", "", "")
	defer db.Close()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/quarantine", nil)
	w := httptest.NewRecorder()

	server.Router().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200 for AuthMode=none without credentials, got %d", w.Code)
	}
}

func TestAuth_ModeBasic(t *testing.T) {
	server, db := createServerWithAuth(t, "basic", "admin", "secretPass", "")
	defer db.Close()

	// 1. Unauthenticated request -> 401
	req := httptest.NewRequest(http.MethodGet, "/api/v1/quarantine", nil)
	w := httptest.NewRecorder()
	server.Router().ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401 for missing basic auth, got %d", w.Code)
	}
	if w.Header().Get("WWW-Authenticate") == "" {
		t.Errorf("expected WWW-Authenticate header")
	}

	// 2. Wrong credentials -> 401
	req = httptest.NewRequest(http.MethodGet, "/api/v1/quarantine", nil)
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("admin:wrong")))
	w = httptest.NewRecorder()
	server.Router().ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401 for wrong basic auth, got %d", w.Code)
	}

	// 3. Valid credentials -> 200
	req = httptest.NewRequest(http.MethodGet, "/api/v1/quarantine", nil)
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("admin:secretPass")))
	w = httptest.NewRecorder()
	server.Router().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200 for valid basic auth, got %d", w.Code)
	}
}

func TestAuth_ModeBearer(t *testing.T) {
	server, db := createServerWithAuth(t, "bearer", "", "", "my_secret_token_123")
	defer db.Close()

	// 1. Unauthenticated request -> 401
	req := httptest.NewRequest(http.MethodGet, "/api/v1/quarantine", nil)
	w := httptest.NewRecorder()
	server.Router().ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401 for missing bearer token, got %d", w.Code)
	}

	// 2. Invalid bearer token -> 401
	req = httptest.NewRequest(http.MethodGet, "/api/v1/quarantine", nil)
	req.Header.Set("Authorization", "Bearer invalid_token")
	w = httptest.NewRecorder()
	server.Router().ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401 for invalid bearer token, got %d", w.Code)
	}

	// 3. Valid bearer token -> 200
	req = httptest.NewRequest(http.MethodGet, "/api/v1/quarantine", nil)
	req.Header.Set("Authorization", "Bearer my_secret_token_123")
	w = httptest.NewRecorder()
	server.Router().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200 for valid bearer token, got %d", w.Code)
	}

	// 4. Valid via X-API-Key header -> 200
	req = httptest.NewRequest(http.MethodGet, "/api/v1/quarantine", nil)
	req.Header.Set("X-API-Key", "my_secret_token_123")
	w = httptest.NewRecorder()
	server.Router().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200 for valid X-API-Key, got %d", w.Code)
	}
}

func TestAuth_HealthExemptions(t *testing.T) {
	server, db := createServerWithAuth(t, "basic", "admin", "secretPass", "")
	defer db.Close()

	// Both /healthz and /api/v1/health should pass without any credentials
	for _, path := range []string{"/healthz", "/api/v1/health"} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		w := httptest.NewRecorder()
		server.Router().ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected 200 OK for exempted health probe '%s', got %d", path, w.Code)
		}
	}
}

func TestUIAuth_LoginAndChangePassword(t *testing.T) {
	server, db := createServerWithAuth(t, "none", "", "", "")
	defer db.Close()

	// 1. Initial login with wrong password -> 401
	body, _ := json.Marshal(map[string]string{"password": "wrong"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/ui-login", bytes.NewReader(body))
	w := httptest.NewRecorder()
	server.Router().ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401 for wrong UI password, got %d", w.Code)
	}

	// 2. Initial login with default password "123456" -> 200
	body, _ = json.Marshal(map[string]string{"password": "123456"})
	req = httptest.NewRequest(http.MethodPost, "/api/v1/auth/ui-login", bytes.NewReader(body))
	w = httptest.NewRecorder()
	server.Router().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200 for default UI password, got %d", w.Code)
	}

	// 3. Change password from "123456" to "customPass2026"
	body, _ = json.Marshal(map[string]string{
		"current_password": "123456",
		"new_password":     "customPass2026",
	})
	req = httptest.NewRequest(http.MethodPost, "/api/v1/auth/ui-password", bytes.NewReader(body))
	w = httptest.NewRecorder()
	server.Router().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200 for changing password, got %d", w.Code)
	}

	// 4. Old password "123456" now fails -> 401
	body, _ = json.Marshal(map[string]string{"password": "123456"})
	req = httptest.NewRequest(http.MethodPost, "/api/v1/auth/ui-login", bytes.NewReader(body))
	w = httptest.NewRecorder()
	server.Router().ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected old password to fail with 401, got %d", w.Code)
	}

	// 5. New password "customPass2026" succeeds -> 200
	body, _ = json.Marshal(map[string]string{"password": "customPass2026"})
	req = httptest.NewRequest(http.MethodPost, "/api/v1/auth/ui-login", bytes.NewReader(body))
	w = httptest.NewRecorder()
	server.Router().ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected new password to succeed with 200, got %d", w.Code)
	}
}
