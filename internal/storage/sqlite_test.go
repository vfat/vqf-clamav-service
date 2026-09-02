package storage

import (
	"path/filepath"
	"testing"
	"time"
)

func TestNewDB_InitAndMigrations(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test_clamav.db")

	db, err := NewDB(dbPath)
	if err != nil {
		t.Fatalf("failed to initialize SQLite DB: %v", err)
	}
	defer db.Close()

	// Verifikasi tabel berhasil dibuat
	tables := []string{
		"scan_audit_logs",
		"quarantine_records",
		"whitelist_signatures",
		"api_keys",
		"system_settings",
	}

	for _, table := range tables {
		var name string
		err := db.conn.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name=?", table).Scan(&name)
		if err != nil || name != table {
			t.Errorf("expected table '%s' to exist in database, got error: %v", table, err)
		}
	}
}

func TestScanAuditLogs_InsertQueryAndPurge(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test_audit.db")

	db, err := NewDB(dbPath)
	if err != nil {
		t.Fatalf("failed to init DB: %v", err)
	}
	defer db.Close()

	logEntry := ScanAuditLog{
		ID:               "audit_01j7xyz894",
		Timestamp:        time.Now().UTC(),
		ConsumerName:     "Billing-App",
		ClientIP:         "192.168.1.50",
		FileName:         "invoice.pdf.exe",
		FileSizeBytes:    2048576,
		FileSHA256:       "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		Verdict:          "INFECTED",
		VirusName:        "Win.Trojan.Agent-123",
		ScanDurationMs:   45,
		QuarantineID:     "Q-20260902-01",
	}

	if err := db.InsertScanAuditLog(logEntry); err != nil {
		t.Fatalf("failed to insert scan audit log: %v", err)
	}

	// Query by verdict
	logs, total, err := db.ListScanAuditLogs(AuditFilter{Verdict: "INFECTED", Limit: 10, Offset: 0})
	if err != nil {
		t.Fatalf("failed to list scan audit logs: %v", err)
	}

	if total != 1 || len(logs) != 1 {
		t.Fatalf("expected 1 audit log, got total=%d, len=%d", total, len(logs))
	}

	if logs[0].VirusName != "Win.Trojan.Agent-123" {
		t.Errorf("expected virus name 'Win.Trojan.Agent-123', got '%s'", logs[0].VirusName)
	}

	// Test Purge: insert old log 5 days ago
	oldLog := logEntry
	oldLog.ID = "audit_old_01"
	oldLog.Timestamp = time.Now().UTC().AddDate(0, 0, -5)
	if err := db.InsertScanAuditLog(oldLog); err != nil {
		t.Fatalf("failed to insert old log: %v", err)
	}

	// Purge logs older than 3 days (default policy)
	purgedCount, err := db.PurgeAuditLogs(3)
	if err != nil {
		t.Fatalf("failed to purge old audit logs: %v", err)
	}

	if purgedCount != 1 {
		t.Errorf("expected 1 old log purged, got %d", purgedCount)
	}
}

func TestQuarantineRecords_InsertRestoreAndWhitelist(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test_quar.db")

	db, err := NewDB(dbPath)
	if err != nil {
		t.Fatalf("failed to init DB: %v", err)
	}
	defer db.Close()

	quar := QuarantineRecord{
		ID:               "Q-20260902-01J7XYZ",
		OriginalFilename: "document.pdf.exe",
		FileSizeBytes:    1048576,
		FileSHA256:       "275a021bbfb6489e54d471899f7db9d1663fc695ec2fe2a2c4538aabf651fd0f",
		VirusName:        "Eicar-Test-Signature",
		SourceConsumer:   "Upload-Gateway",
		StoredPath:       "/data/quarantine/Q-20260902-01J7XYZ.quarantine",
		Status:           "QUARANTINED",
		CreatedAt:        time.Now().UTC(),
		ExpiresAt:        time.Now().UTC().Add(7 * 24 * time.Hour),
	}

	if err := db.InsertQuarantineRecord(quar); err != nil {
		t.Fatalf("failed to insert quarantine record: %v", err)
	}

	// Get by ID
	fetched, err := db.GetQuarantineRecord(quar.ID)
	if err != nil {
		t.Fatalf("failed to get quarantine record: %v", err)
	}

	if fetched.Status != "QUARANTINED" {
		t.Errorf("expected status 'QUARANTINED', got '%s'", fetched.Status)
	}

	// Restore record & whitelist hash
	err = db.RestoreQuarantineRecord(quar.ID, "security-admin@corp.internal", "False positive verified", true)
	if err != nil {
		t.Fatalf("failed to restore quarantine record: %v", err)
	}

	updated, _ := db.GetQuarantineRecord(quar.ID)
	if updated.Status != "RESTORED" {
		t.Errorf("expected status 'RESTORED', got '%s'", updated.Status)
	}

	// Check whitelist
	whitelisted, err := db.IsWhitelisted(quar.FileSHA256)
	if err != nil {
		t.Fatalf("failed to check whitelist: %v", err)
	}
	if !whitelisted {
		t.Errorf("expected SHA256 '%s' to be whitelisted", quar.FileSHA256)
	}
}
