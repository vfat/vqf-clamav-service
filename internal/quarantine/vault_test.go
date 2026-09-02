package quarantine

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"path/filepath"
	"strings"
	"testing"

	"github.com/vfat/vqf-clamav-service/internal/storage"
)

func TestVault_QuarantineAndRestore(t *testing.T) {
	tmpDir := t.TempDir()
	vaultDir := filepath.Join(tmpDir, "vault")
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := storage.NewDB(dbPath)
	if err != nil {
		t.Fatalf("failed to init test db: %v", err)
	}
	defer db.Close()

	vault := NewVault(vaultDir, db)

	// File berbahaya untuk dikarantina
	maliciousContent := "X5O!P%@AP[4\\PZX54(P^)7CC)7}$EICAR-STANDARD-ANTIVIRUS-TEST-FILE!$H+H*"
	reader := strings.NewReader(maliciousContent)

	hasher := sha256.New()
	hasher.Write([]byte(maliciousContent))
	expectedHash := hex.EncodeToString(hasher.Sum(nil))

	ctx := context.Background()

	// 1. Karantina file
	record, err := vault.QuarantineFile(ctx, "malicious_script.bat", "Test-Consumer", "Eicar-Test-Signature", reader, 7)
	if err != nil {
		t.Fatalf("QuarantineFile failed: %v", err)
	}

	if record.OriginalFilename != "malicious_script.bat" {
		t.Errorf("expected filename 'malicious_script.bat', got '%s'", record.OriginalFilename)
	}
	if record.FileSHA256 != expectedHash {
		t.Errorf("expected hash '%s', got '%s'", expectedHash, record.FileSHA256)
	}
	if record.Status != "QUARANTINED" {
		t.Errorf("expected status 'QUARANTINED', got '%s'", record.Status)
	}

	// 2. Restore file
	restoredReader, updatedRec, err := vault.RestoreFile(ctx, record.ID, "security-admin@corp.internal", "False positive verified", true)
	if err != nil {
		t.Fatalf("RestoreFile failed: %v", err)
	}
	defer restoredReader.Close()

	if updatedRec.Status != "RESTORED" {
		t.Errorf("expected updated status 'RESTORED', got '%s'", updatedRec.Status)
	}

	restoredBytes, err := io.ReadAll(restoredReader)
	if err != nil {
		t.Fatalf("failed to read restored payload: %v", err)
	}

	if string(restoredBytes) != maliciousContent {
		t.Errorf("restored content mismatch. Expected '%s', got '%s'", maliciousContent, string(restoredBytes))
	}

	// 3. Verifikasi auto-whitelist
	whitelisted, err := db.IsWhitelisted(expectedHash)
	if err != nil || !whitelisted {
		t.Errorf("expected hash to be whitelisted in DB, got whitelisted=%v, err=%v", whitelisted, err)
	}
}

func TestVault_ScrambleIntegrity(t *testing.T) {
	data := []byte("Sensitive malicious binary payload with specific bytes 1234567890")
	scrambled := scrambleBytes(data)

	if bytes.Equal(data, scrambled) {
		t.Fatalf("scrambled payload must not be identical to raw plaintext data")
	}

	descrambled := descrambleBytes(scrambled)
	if !bytes.Equal(data, descrambled) {
		t.Fatalf("descramble failed to reconstruct original payload")
	}
}
