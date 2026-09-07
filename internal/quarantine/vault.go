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

	"github.com/vfat/vqf-clamav-service/internal/crypto"
	"github.com/vfat/vqf-clamav-service/internal/storage"
)

// MagicHeader is the constant 14-byte prefix identifying AES-256-GCM encrypted quarantine files.
const MagicHeader = "VQF_AESGCM_V1\n"

// Vault manages safe quarantined file isolation, authenticated encryption, and restoration.
type Vault struct {
	vaultDir          string
	db                *storage.DB
	masterKey         []byte
	maxQuarantineSize int64
}

// NewVault initializes the quarantine vault manager.
// Accepts optional 32-byte masterKey for AES-256-GCM authenticated encryption.
func NewVault(vaultDir string, db *storage.DB, masterKey ...[]byte) *Vault {
	if vaultDir == "" {
		vaultDir = "/data/quarantine"
	}
	_ = os.MkdirAll(vaultDir, 0700)

	var key []byte
	if len(masterKey) > 0 && len(masterKey[0]) == 32 {
		key = masterKey[0]
	} else {
		// Auto-ensure fallback master key
		_, k, _, _ := crypto.EnsureMasterKey("")
		key = k
	}

	return &Vault{
		vaultDir:          vaultDir,
		db:                db,
		masterKey:         key,
		maxQuarantineSize: 100 * 1024 * 1024, // 100MB default limit
	}
}

// SetMaxQuarantineBytes configures the maximum file size acceptable by QuarantineFile to prevent memory exhaustion.
func (v *Vault) SetMaxQuarantineBytes(limit int64) {
	v.maxQuarantineSize = limit
}

// QuarantineFile isolates an infected file into the vault, calculates its SHA256,
// encrypts its binary content with AES-256-GCM, sets 0600 permissions, and logs to the database.
func (v *Vault) QuarantineFile(ctx context.Context, originalName, consumer, virusName string, r io.Reader, retentionDays int) (*storage.QuarantineRecord, error) {
	if retentionDays <= 0 {
		retentionDays = 7 // Default 7 days
	}

	// Memory bomb safeguard: enforce strict reader ceiling
	if v.maxQuarantineSize > 0 {
		r = io.LimitReader(r, v.maxQuarantineSize+1)
	}

	rawBytes, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read payload for quarantine: %w", err)
	}

	if v.maxQuarantineSize > 0 && int64(len(rawBytes)) > v.maxQuarantineSize {
		return nil, fmt.Errorf("payload exceeds maximum quarantine limit of %d bytes", v.maxQuarantineSize)
	}

	// Calculate SHA-256 hash
	hasher := sha256.New()
	hasher.Write(rawBytes)
	hashStr := hex.EncodeToString(hasher.Sum(nil))

	// Generate Vault ID
	now := time.Now().UTC()
	randSuffix := make([]byte, 8)
	_, _ = rand.Read(randSuffix)
	vaultID := fmt.Sprintf("Q-%s-%s", now.Format("20060102"), hex.EncodeToString(randSuffix))

	fileName := vaultID + ".quarantine"
	filePath := filepath.Join(v.vaultDir, fileName)

	// Encrypt with AES-256-GCM + Magic Header
	var storedData []byte
	if len(v.masterKey) == 32 {
		encrypted, err := crypto.EncryptAESGCMBytes(rawBytes, v.masterKey)
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt quarantine payload: %w", err)
		}
		storedData = append([]byte(MagicHeader), encrypted...)
	} else {
		// Fallback to legacy XOR if no valid key is present
		storedData = scrambleBytes(rawBytes)
	}

	if err := os.WriteFile(filePath, storedData, 0600); err != nil {
		return nil, fmt.Errorf("failed to write quarantine file %s: %w", filePath, err)
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

// decryptPayload inspects the magic header to either decrypt AES-256-GCM or fallback to legacy XOR descrambling.
func (v *Vault) decryptPayload(storedBytes []byte) ([]byte, error) {
	magic := []byte(MagicHeader)
	if len(storedBytes) >= len(magic) && bytes.Equal(storedBytes[:len(magic)], magic) {
		cipherData := storedBytes[len(magic):]
		return crypto.DecryptAESGCMBytes(cipherData, v.masterKey)
	}

	// Dynamic fallback for legacy XOR 0xA5 scrambled files
	return descrambleBytes(storedBytes), nil
}

// RestoreFile de-scrambles/decrypts a quarantined file, updates its database status,
// and optionally adds its SHA256 hash to the whitelist.
func (v *Vault) RestoreFile(ctx context.Context, id, restoredBy, reason string, autoWhitelist bool) (io.ReadCloser, *storage.QuarantineRecord, error) {
	record, err := v.db.GetQuarantineRecord(id)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to find quarantine record: %w", err)
	}

	if record.Status != "QUARANTINED" {
		return nil, nil, fmt.Errorf("cannot restore file with status '%s'", record.Status)
	}

	storedBytes, err := os.ReadFile(record.StoredPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read quarantine file on disk: %w", err)
	}

	plainBytes, err := v.decryptPayload(storedBytes)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decrypt quarantine payload: %w", err)
	}

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

// DownloadFile decrypts/de-scrambles and returns the plain reader for a quarantined payload.
func (v *Vault) DownloadFile(ctx context.Context, id string) (io.ReadCloser, *storage.QuarantineRecord, error) {
	record, err := v.db.GetQuarantineRecord(id)
	if err != nil {
		return nil, nil, fmt.Errorf("quarantine record not found: %w", err)
	}

	storedBytes, err := os.ReadFile(record.StoredPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read quarantine file on disk: %w", err)
	}

	plainBytes, err := v.decryptPayload(storedBytes)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decrypt quarantine payload: %w", err)
	}

	return io.NopCloser(bytes.NewReader(plainBytes)), record, nil
}

// PurgeExpired deletes physical quarantine files that have exceeded their TTL.
func (v *Vault) PurgeExpired(ctx context.Context) (int, error) {
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
		out[i] = b ^ xorMask
	}
	return out
}

func descrambleBytes(data []byte) []byte {
	out := make([]byte, len(data))
	for i, b := range data {
		out[i] = b ^ xorMask
	}
	return out
}

