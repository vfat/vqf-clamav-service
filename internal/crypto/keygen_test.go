package crypto

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateMasterKey(t *testing.T) {
	keyStr, keyBytes, err := GenerateMasterKey()
	if err != nil {
		t.Fatalf("unexpected error generating key: %v", err)
	}

	if !strings.HasPrefix(keyStr, "clam_sec_") {
		t.Errorf("expected key to have prefix 'clam_sec_', got: %s", keyStr)
	}

	hexPart := strings.TrimPrefix(keyStr, "clam_sec_")
	if len(hexPart) != 64 {
		t.Errorf("expected 64 hex characters (32 bytes), got length: %d", len(hexPart))
	}

	if len(keyBytes) != 32 {
		t.Errorf("expected 32 raw key bytes for AES-256, got: %d", len(keyBytes))
	}
}

func TestEnsureMasterKey_AutoCreateEnv(t *testing.T) {
	tmpDir := t.TempDir()
	envPath := filepath.Join(tmpDir, ".env")

	// File .env belum ada
	keyStr, keyBytes, created, err := EnsureMasterKey(envPath)
	if err != nil {
		t.Fatalf("EnsureMasterKey failed on missing env: %v", err)
	}

	if !created {
		t.Errorf("expected created=true when .env was missing, got false")
	}

	if !strings.HasPrefix(keyStr, "clam_sec_") {
		t.Errorf("expected generated key with 'clam_sec_' prefix, got: %s", keyStr)
	}

	if len(keyBytes) != 32 {
		t.Errorf("expected 32 bytes key, got: %d", len(keyBytes))
	}

	// Verifikasi file .env benar-benar tertulis
	content, err := os.ReadFile(envPath)
	if err != nil {
		t.Fatalf("failed to read created .env file: %v", err)
	}

	if !strings.Contains(string(content), "ENCRYPTION_KEY="+keyStr) {
		t.Errorf("expected .env to contain 'ENCRYPTION_KEY=%s', got content:\n%s", keyStr, string(content))
	}

	// Panggilan kedua: harus membaca key yang sudah ada (tidak menimpa)
	secondKeyStr, secondKeyBytes, secondCreated, err := EnsureMasterKey(envPath)
	if err != nil {
		t.Fatalf("second EnsureMasterKey failed: %v", err)
	}

	if secondCreated {
		t.Errorf("expected created=false on existing env key, got true")
	}

	if secondKeyStr != keyStr {
		t.Errorf("expected identical key on second read. First: %s, Second: %s", keyStr, secondKeyStr)
	}

	if len(secondKeyBytes) != 32 {
		t.Errorf("expected 32 bytes on second key read, got: %d", len(secondKeyBytes))
	}
}

func TestParseMasterKey(t *testing.T) {
	validKey := "clam_sec_0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	rawBytes, err := ParseMasterKey(validKey)
	if err != nil {
		t.Fatalf("unexpected error parsing valid key: %v", err)
	}
	if len(rawBytes) != 32 {
		t.Errorf("expected 32 bytes, got %d", len(rawBytes))
	}

	invalidPrefix := "secret_0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	if _, err := ParseMasterKey(invalidPrefix); err == nil {
		t.Errorf("expected error for invalid prefix, got nil")
	}

	invalidLen := "clam_sec_1234"
	if _, err := ParseMasterKey(invalidLen); err == nil {
		t.Errorf("expected error for invalid length, got nil")
	}
}
