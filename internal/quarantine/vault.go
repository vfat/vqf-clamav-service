package quarantine

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/vfat/vqf-clamav-service/internal/storage"
)

// Vault manages safe quarantined file isolation, scrambling, and restoration.
type Vault struct {
	vaultDir string
	db       *storage.DB
}

// NewVault initializes the quarantine vault manager.
func NewVault(vaultDir string, db *storage.DB) *Vault {
	if vaultDir == "" {
		vaultDir = "/data/quarantine"
	}
	os.MkdirAll(vaultDir, 0700)

	return &Vault{
		vaultDir: vaultDir,
		db:       db,
	}
}

// QuarantineFile isolates an infected file into the vault, calculates its SHA256,
// scrambles its binary content, sets 0600 permissions, and logs to the database.
func (v *Vault) QuarantineFile(ctx context.Context, originalName, consumer, virusName string, r io.Reader, retentionDays int) (*storage.QuarantineRecord, error) {
	if retentionDays <= 0 {
		retentionDays = 7 // Default 7 days
	}

	rawBytes, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read payload for quarantine: %w", err)
	}

	// Calculate SHA-256 hash
	hasher := sha256.New()
	hasher.Write(rawBytes)
	hashStr := hex.EncodeToString(hasher.Sum(nil))

	// Generate Vault ID
	now := time.Now().UTC()
	randSuffix := make([]byte, 8)
	rand.Read(randSuffix)
	vaultID := fmt.Sprintf("Q-%s-%s", now.Format("20060102"), hex.EncodeToString(randSuffix))

	fileName := vaultID + ".quarantine"
	filePath := filepath.Join(v.vaultDir, fileName)

	// Scramble and write with 0600 permissions
	scrambled := scrambleBytes(rawBytes)
	if err := os.WriteFile(filePath, scrambled, 0600); err != nil {
		return nil, fmt.Errorf("failed to write scrambled quarantine file %s: %w", filePath, err)
	}

	record := storage.QuarantineRecord{
		ID:               vaultID,
		OriginalFilename: originalName,
		FileSizeBytes:    int64(len(rawBytes)),
		FileSHA256:       hashStr,
		VirusName:        virusName,
		SourceConsumer:   consumer,
		StoredPath:       filePath,
		Status:           "QUARANTINED",
		CreatedAt:        now,
		ExpiresAt:        now.Add(time.Duration(retentionDays) * 24 * time.Hour),
	}

	if err := v.db.InsertQuarantineRecord(record); err != nil {
		return nil, fmt.Errorf("failed to insert quarantine DB record: %w", err)
	}

	return &record, nil
}

// RestoreFile de-scrambles a quarantined file, updates its database status,
// and optionally adds its SHA256 hash to the whitelist.
func (v *Vault) RestoreFile(ctx context.Context, id, restoredBy, reason string, autoWhitelist bool) (io.ReadCloser, *storage.QuarantineRecord, error) {
	record, err := v.db.GetQuarantineRecord(id)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to find quarantine record: %w", err)
	}

	if record.Status != "QUARANTINED" {
		return nil, nil, fmt.Errorf("cannot restore file with status '%s'", record.Status)
	}

	// Read and de-scramble file
	scrambledBytes, err := os.ReadFile(record.StoredPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read quarantine file on disk: %w", err)
	}

	plainBytes := descrambleBytes(scrambledBytes)

	// Update record in database
	if err := v.db.RestoreQuarantineRecord(id, restoredBy, reason, autoWhitelist); err != nil {
		return nil, nil, fmt.Errorf("failed to update quarantine database record: %w", err)
	}

	updatedRec, err := v.db.GetQuarantineRecord(id)
	if err != nil {
		return nil, nil, err
	}

	return io.NopCloser(bytes.NewReader(plainBytes)), updatedRec, nil
}

// DownloadFile de-scrambles and returns the plain reader for a quarantined payload.
func (v *Vault) DownloadFile(ctx context.Context, id string) (io.ReadCloser, *storage.QuarantineRecord, error) {
	record, err := v.db.GetQuarantineRecord(id)
	if err != nil {
		return nil, nil, fmt.Errorf("quarantine record not found: %w", err)
	}

	scrambledBytes, err := os.ReadFile(record.StoredPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read quarantine file on disk: %w", err)
	}

	plainBytes := descrambleBytes(scrambledBytes)
	return io.NopCloser(bytes.NewReader(plainBytes)), record, nil
}

// PurgeExpired deletes physical quarantine files that have exceeded their TTL.
func (v *Vault) PurgeExpired(ctx context.Context) (int, error) {
	// Find and remove files older than expires_at
	entries, err := os.ReadDir(v.vaultDir)
	if err != nil {
		return 0, err
	}

	purged := 0
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".quarantine" {
			info, err := entry.Info()
			if err == nil && time.Since(info.ModTime()) > 7*24*time.Hour {
				if err := os.Remove(filepath.Join(v.vaultDir, entry.Name())); err == nil {
					purged++
				}
			}
		}
	}

	return purged, nil
}

// DeleteFile permanently removes the physical scrambled payload and its database record.
func (v *Vault) DeleteFile(ctx context.Context, id string) error {
	record, err := v.db.GetQuarantineRecord(id)
	if err != nil {
		return fmt.Errorf("quarantine record not found: %w", err)
	}

	if record.StoredPath != "" {
		_ = os.Remove(record.StoredPath)
	}

	return v.db.DeleteQuarantineRecord(id)
}

const xorMask = 0xA5

func scrambleBytes(data []byte) []byte {
	out := make([]byte, len(data))
	for i, b := range data {
		// XOR mask + bit inversion
		out[i] = b ^ xorMask
	}
	return out
}

func descrambleBytes(data []byte) []byte {
	// Symmetric XOR reversal
	out := make([]byte, len(data))
	for i, b := range data {
		out[i] = b ^ xorMask
	}
	return out
}
