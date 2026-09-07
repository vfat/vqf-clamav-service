package quarantine

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

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

func TestVault_AESGCM_EncryptionAndRestore(t *testing.T) {
	tmpDir := t.TempDir()
	vaultDir := filepath.Join(tmpDir, "vault")
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := storage.NewDB(dbPath)
	if err != nil {
		t.Fatalf("failed to init test db: %v", err)
	}
	defer db.Close()

	// 32-byte master key
	masterKey := []byte("01234567890123456789012345678901")
	vault := NewVault(vaultDir, db, masterKey)

	plainContent := "MALWARE-PAYLOAD-SAMPLE-FOR-AES-GCM-TEST"
	ctx := context.Background()

	record, err := vault.QuarantineFile(ctx, "sample.exe", "SecOps", "Win.Trojan.Generic", strings.NewReader(plainContent), 7)
	if err != nil {
		t.Fatalf("QuarantineFile failed: %v", err)
	}

	// 1. Verify file on disk has Magic Header
	diskBytes, err := os.ReadFile(record.StoredPath)
	if err != nil {
		t.Fatalf("failed to read file on disk: %v", err)
	}

	expectedMagic := []byte("VQF_AESGCM_V1\n")
	if len(diskBytes) < len(expectedMagic) || !bytes.Equal(diskBytes[:len(expectedMagic)], expectedMagic) {
		t.Fatalf("expected stored file to start with magic header 'VQF_AESGCM_V1\\n', got %q", diskBytes[:min(len(diskBytes), 20)])
	}

	// 2. Verify disk data is NOT static XOR 0xA5
	xorReversed := make([]byte, len(diskBytes)-len(expectedMagic))
	for i, b := range diskBytes[len(expectedMagic):] {
		xorReversed[i] = b ^ 0xA5
	}
	if string(xorReversed) == plainContent {
		t.Fatalf("disk payload is still scrambled with static XOR, expected authenticated AES-256-GCM encryption!")
	}

	// 3. Verify RestoreFile decrypts AES-GCM correctly
	restoredR, _, err := vault.RestoreFile(ctx, record.ID, "admin", "verified clean", false)
	if err != nil {
		t.Fatalf("RestoreFile failed: %v", err)
	}
	defer restoredR.Close()

	restoredBytes, err := io.ReadAll(restoredR)
	if err != nil {
		t.Fatalf("failed to read restored content: %v", err)
	}
	if string(restoredBytes) != plainContent {
		t.Fatalf("restored content mismatch. Got %q, expected %q", string(restoredBytes), plainContent)
	}

	// 4. Verify DownloadFile also decrypts AES-GCM correctly
	downR, _, err := vault.DownloadFile(ctx, record.ID)
	if err != nil {
		t.Fatalf("DownloadFile failed: %v", err)
	}
	defer downR.Close()

	downBytes, _ := io.ReadAll(downR)
	if string(downBytes) != plainContent {
		t.Fatalf("downloaded content mismatch. Got %q, expected %q", string(downBytes), plainContent)
	}
}

func TestVault_LegacyXOR_BackwardCompatibility(t *testing.T) {
	tmpDir := t.TempDir()
	vaultDir := filepath.Join(tmpDir, "vault")
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := storage.NewDB(dbPath)
	if err != nil {
		t.Fatalf("failed to init test db: %v", err)
	}
	defer db.Close()

	masterKey := []byte("01234567890123456789012345678901")
	vault := NewVault(vaultDir, db, masterKey)

	// Craft a legacy file scrambled with XOR 0xA5 without any magic header
	legacyPlain := "LEGACY-QUARANTINED-VIRUS-CONTENT-BEFORE-MIGRATION"
	legacyScrambled := make([]byte, len(legacyPlain))
	for i, b := range []byte(legacyPlain) {
		legacyScrambled[i] = b ^ 0xA5
	}

	legacyID := "Q-LEGACY-001"
	legacyPath := filepath.Join(vaultDir, legacyID+".quarantine")
	_ = os.MkdirAll(vaultDir, 0700)
	if err := os.WriteFile(legacyPath, legacyScrambled, 0600); err != nil {
		t.Fatalf("failed to write legacy file: %v", err)
	}

	// Insert DB record for legacy file
	now := time.Now().UTC()
	err = db.InsertQuarantineRecord(storage.QuarantineRecord{
		ID:               legacyID,
		OriginalFilename: "old_trojan.exe",
		FileSizeBytes:    int64(len(legacyPlain)),
		FileSHA256:       "legacy_sha256",
		VirusName:        "Legacy.Trojan.Old",
		StoredPath:       legacyPath,
		Status:           "QUARANTINED",
		CreatedAt:        now,
		ExpiresAt:        now.Add(7 * 24 * time.Hour),
	})
	if err != nil {
		t.Fatalf("failed to insert legacy record: %v", err)
	}

	ctx := context.Background()

	// Verify RestoreFile falls back to XOR for legacy file
	restoredR, _, err := vault.RestoreFile(ctx, legacyID, "admin", "false positive on legacy", false)
	if err != nil {
		t.Fatalf("RestoreFile on legacy XOR file failed: %v", err)
	}
	defer restoredR.Close()

	restoredBytes, _ := io.ReadAll(restoredR)
	if string(restoredBytes) != legacyPlain {
		t.Fatalf("legacy fallback restore mismatch. Expected %q, got %q", legacyPlain, string(restoredBytes))
	}

	// Verify DownloadFile also falls back to XOR for legacy file
	downR, _, err := vault.DownloadFile(ctx, legacyID)
	if err != nil {
		t.Fatalf("DownloadFile on legacy XOR file failed: %v", err)
	}
	defer downR.Close()

	downBytes, _ := io.ReadAll(downR)
	if string(downBytes) != legacyPlain {
		t.Fatalf("legacy fallback download mismatch. Expected %q, got %q", legacyPlain, string(downBytes))
	}
}

func TestVault_AntiMemoryBomb_LimitReader(t *testing.T) {
	tmpDir := t.TempDir()
	vaultDir := filepath.Join(tmpDir, "vault")
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := storage.NewDB(dbPath)
	if err != nil {
		t.Fatalf("failed to init test db: %v", err)
	}
	defer db.Close()

	masterKey := []byte("01234567890123456789012345678901")
	vault := NewVault(vaultDir, db, masterKey)
	vault.SetMaxQuarantineBytes(1024) // Set small 1KB limit for testing

	// Attempt to quarantine payload exceeding limit (2KB)
	oversizedContent := strings.Repeat("A", 2048)
	ctx := context.Background()

	_, err = vault.QuarantineFile(ctx, "bomb.bin", "attacker", "Heuristics.ZipBomb", strings.NewReader(oversizedContent), 7)
	if err == nil {
		t.Fatalf("expected QuarantineFile to reject oversized payload exceeding MaxQuarantineBytes limit, got nil error")
	}
}
