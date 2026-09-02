package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

// DB wraps SQLite database connection.
type DB struct {
	conn *sql.DB
}

// ScanAuditLog represents a row in scan_audit_logs.
type ScanAuditLog struct {
	ID             string    `json:"id"`
	Timestamp      time.Time `json:"timestamp"`
	ConsumerName   string    `json:"consumer_name"`
	ClientIP       string    `json:"client_ip"`
	FileName       string    `json:"file_name"`
	FileSizeBytes  int64     `json:"file_size_bytes"`
	FileSHA256     string    `json:"file_sha256"`
	Verdict        string    `json:"verdict"`
	VirusName      string    `json:"virus_name,omitempty"`
	ScanDurationMs int64     `json:"scan_duration_ms"`
	QuarantineID   string    `json:"quarantine_id,omitempty"`
}

// AuditFilter contains filter parameters for listing audit logs.
type AuditFilter struct {
	Verdict      string
	ConsumerName string
	From         *time.Time
	To           *time.Time
	Limit        int
	Offset       int
}

// QuarantineRecord represents a row in quarantine_records.
type QuarantineRecord struct {
	ID               string     `json:"id"`
	OriginalFilename string     `json:"original_filename"`
	FileSizeBytes    int64      `json:"file_size_bytes"`
	FileSHA256       string     `json:"file_sha256"`
	VirusName        string     `json:"virus_name"`
	SourceConsumer   string     `json:"source_consumer"`
	StoredPath       string     `json:"stored_path"`
	Status           string     `json:"status"` // QUARANTINED, RESTORED, DELETED
	CreatedAt        time.Time  `json:"created_at"`
	ExpiresAt        time.Time  `json:"expires_at"`
	RestoredAt       *time.Time `json:"restored_at,omitempty"`
	RestoredBy       string     `json:"restored_by,omitempty"`
	RestoreReason    string     `json:"restore_reason,omitempty"`
}

// NewDB initializes SQLite connection with WAL mode and runs table migrations.
func NewDB(dbPath string) (*DB, error) {
	if dbPath == "" {
		dbPath = "/data/clamav-service.db"
	}

	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory %s: %w", dir, err)
	}

	dsn := fmt.Sprintf("file:%s?_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)&_pragma=busy_timeout(5000)&_pragma=foreign_keys(ON)", dbPath)
	conn, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite database: %w", err)
	}

	conn.SetMaxOpenConns(1) // SQLite single writer safe concurrency

	db := &DB{conn: conn}
	if err := db.migrate(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return db, nil
}

// Close closes the underlying SQLite database connection.
func (db *DB) Close() error {
	if db.conn != nil {
		return db.conn.Close()
	}
	return nil
}

func (db *DB) migrate() error {
	schema := `
	CREATE TABLE IF NOT EXISTS scan_audit_logs (
		id               TEXT PRIMARY KEY,
		timestamp        DATETIME NOT NULL,
		consumer_name    TEXT,
		client_ip        TEXT,
		file_name        TEXT,
		file_size_bytes  INTEGER,
		file_sha256      TEXT NOT NULL,
		verdict          TEXT NOT NULL,
		virus_name       TEXT,
		scan_duration_ms INTEGER NOT NULL,
		quarantine_id    TEXT
	);
	CREATE INDEX IF NOT EXISTS idx_audit_timestamp ON scan_audit_logs(timestamp);
	CREATE INDEX IF NOT EXISTS idx_audit_verdict ON scan_audit_logs(verdict);
	CREATE INDEX IF NOT EXISTS idx_audit_sha256 ON scan_audit_logs(file_sha256);

	CREATE TABLE IF NOT EXISTS quarantine_records (
		id                 TEXT PRIMARY KEY,
		original_filename  TEXT NOT NULL,
		file_size_bytes    INTEGER NOT NULL,
		file_sha256        TEXT NOT NULL,
		virus_name         TEXT NOT NULL,
		source_consumer    TEXT,
		stored_path        TEXT NOT NULL,
		status             TEXT NOT NULL,
		created_at         DATETIME NOT NULL,
		expires_at         DATETIME NOT NULL,
		restored_at        DATETIME,
		restored_by        TEXT,
		restore_reason     TEXT
	);
	CREATE INDEX IF NOT EXISTS idx_quar_status ON quarantine_records(status);
	CREATE INDEX IF NOT EXISTS idx_quar_expires ON quarantine_records(expires_at);

	CREATE TABLE IF NOT EXISTS whitelist_signatures (
		sha256_hash TEXT PRIMARY KEY,
		description TEXT,
		added_by    TEXT,
		created_at  DATETIME NOT NULL
	);

	CREATE TABLE IF NOT EXISTS api_keys (
		id          TEXT PRIMARY KEY,
		key_hash    TEXT NOT NULL UNIQUE,
		name        TEXT NOT NULL,
		permissions TEXT NOT NULL,
		is_active   INTEGER NOT NULL DEFAULT 1,
		created_at  DATETIME NOT NULL
	);

	CREATE TABLE IF NOT EXISTS system_settings (
		key             TEXT PRIMARY KEY,
		value_encrypted TEXT NOT NULL,
		updated_at      DATETIME NOT NULL
	);
	`
	_, err := db.conn.Exec(schema)
	return err
}

// InsertScanAuditLog inserts a new log record.
func (db *DB) InsertScanAuditLog(l ScanAuditLog) error {
	query := `
	INSERT INTO scan_audit_logs (
		id, timestamp, consumer_name, client_ip, file_name,
		file_size_bytes, file_sha256, verdict, virus_name,
		scan_duration_ms, quarantine_id
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := db.conn.Exec(query,
		l.ID, l.Timestamp.Format(time.RFC3339), l.ConsumerName, l.ClientIP, l.FileName,
		l.FileSizeBytes, l.FileSHA256, l.Verdict, l.VirusName,
		l.ScanDurationMs, l.QuarantineID,
	)
	return err
}

// ListScanAuditLogs queries scan logs based on filter with pagination.
func (db *DB) ListScanAuditLogs(f AuditFilter) ([]ScanAuditLog, int, error) {
	if f.Limit <= 0 {
		f.Limit = 50
	}

	whereClause := "WHERE 1=1"
	args := []interface{}{}

	if f.Verdict != "" {
		whereClause += " AND verdict = ?"
		args = append(args, f.Verdict)
	}
	if f.ConsumerName != "" {
		whereClause += " AND consumer_name = ?"
		args = append(args, f.ConsumerName)
	}
	if f.From != nil {
		whereClause += " AND timestamp >= ?"
		args = append(args, f.From.Format(time.RFC3339))
	}
	if f.To != nil {
		whereClause += " AND timestamp <= ?"
		args = append(args, f.To.Format(time.RFC3339))
	}

	// Count total
	countQuery := "SELECT COUNT(*) FROM scan_audit_logs " + whereClause
	var total int
	if err := db.conn.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Query data
	dataQuery := fmt.Sprintf(`
		SELECT id, timestamp, consumer_name, client_ip, file_name,
		       file_size_bytes, file_sha256, verdict, virus_name,
		       scan_duration_ms, quarantine_id
		FROM scan_audit_logs
		%s
		ORDER BY timestamp DESC
		LIMIT ? OFFSET ?
	`, whereClause)

	dataArgs := append(args, f.Limit, f.Offset)
	rows, err := db.conn.Query(dataQuery, dataArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var logs []ScanAuditLog
	for rows.Next() {
		var l ScanAuditLog
		var tsStr string
		var vName, qID sql.NullString

		err := rows.Scan(
			&l.ID, &tsStr, &l.ConsumerName, &l.ClientIP, &l.FileName,
			&l.FileSizeBytes, &l.FileSHA256, &l.Verdict, &vName,
			&l.ScanDurationMs, &qID,
		)
		if err != nil {
			return nil, 0, err
		}

		l.Timestamp, _ = time.Parse(time.RFC3339, tsStr)
		if vName.Valid {
			l.VirusName = vName.String
		}
		if qID.Valid {
			l.QuarantineID = qID.String
		}
		logs = append(logs, l)
	}

	return logs, total, nil
}

// PurgeAuditLogs deletes logs older than retentionDays.
func (db *DB) PurgeAuditLogs(retentionDays int) (int64, error) {
	threshold := time.Now().UTC().AddDate(0, 0, -retentionDays).Format(time.RFC3339)
	res, err := db.conn.Exec("DELETE FROM scan_audit_logs WHERE timestamp < ?", threshold)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// InsertQuarantineRecord inserts a quarantine record.
func (db *DB) InsertQuarantineRecord(q QuarantineRecord) error {
	query := `
	INSERT INTO quarantine_records (
		id, original_filename, file_size_bytes, file_sha256, virus_name,
		source_consumer, stored_path, status, created_at, expires_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := db.conn.Exec(query,
		q.ID, q.OriginalFilename, q.FileSizeBytes, q.FileSHA256, q.VirusName,
		q.SourceConsumer, q.StoredPath, q.Status,
		q.CreatedAt.Format(time.RFC3339), q.ExpiresAt.Format(time.RFC3339),
	)
	return err
}

// GetQuarantineRecord retrieves a single record by ID.
func (db *DB) GetQuarantineRecord(id string) (*QuarantineRecord, error) {
	query := `
	SELECT id, original_filename, file_size_bytes, file_sha256, virus_name,
	       source_consumer, stored_path, status, created_at, expires_at,
	       restored_at, restored_by, restore_reason
	FROM quarantine_records
	WHERE id = ?
	`
	var q QuarantineRecord
	var createdStr, expiresStr string
	var restoredStr, restoredBy, restoreReason sql.NullString

	err := db.conn.QueryRow(query, id).Scan(
		&q.ID, &q.OriginalFilename, &q.FileSizeBytes, &q.FileSHA256, &q.VirusName,
		&q.SourceConsumer, &q.StoredPath, &q.Status, &createdStr, &expiresStr,
		&restoredStr, &restoredBy, &restoreReason,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("quarantine record not found")
	}
	if err != nil {
		return nil, err
	}

	q.CreatedAt, _ = time.Parse(time.RFC3339, createdStr)
	q.ExpiresAt, _ = time.Parse(time.RFC3339, expiresStr)
	if restoredStr.Valid {
		t, _ := time.Parse(time.RFC3339, restoredStr.String)
		q.RestoredAt = &t
	}
	if restoredBy.Valid {
		q.RestoredBy = restoredBy.String
	}
	if restoreReason.Valid {
		q.RestoreReason = restoreReason.String
	}

	return &q, nil
}

// RestoreQuarantineRecord updates record status to RESTORED and optionally adds hash to whitelist.
func (db *DB) RestoreQuarantineRecord(id, restoredBy, reason string, autoWhitelist bool) error {
	rec, err := db.GetQuarantineRecord(id)
	if err != nil {
		return err
	}

	now := time.Now().UTC().Format(time.RFC3339)
	query := `
	UPDATE quarantine_records
	SET status = 'RESTORED', restored_at = ?, restored_by = ?, restore_reason = ?
	WHERE id = ?
	`
	if _, err := db.conn.Exec(query, now, restoredBy, reason, id); err != nil {
		return err
	}

	if autoWhitelist && rec.FileSHA256 != "" {
		return db.AddWhitelist(rec.FileSHA256, "Auto-whitelisted on restore: "+reason, restoredBy)
	}

	return nil
}

// IsWhitelisted checks if a SHA256 hash is in whitelist_signatures.
func (db *DB) IsWhitelisted(sha256Hash string) (bool, error) {
	var count int
	err := db.conn.QueryRow("SELECT COUNT(*) FROM whitelist_signatures WHERE sha256_hash = ?", sha256Hash).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// AddWhitelist inserts a SHA256 hash into whitelist_signatures.
func (db *DB) AddWhitelist(sha256Hash, description, addedBy string) error {
	query := `
	INSERT OR REPLACE INTO whitelist_signatures (sha256_hash, description, added_by, created_at)
	VALUES (?, ?, ?, ?)
	`
	_, err := db.conn.Exec(query, sha256Hash, description, addedBy, time.Now().UTC().Format(time.RFC3339))
	return err
}

// RemoveWhitelist removes a SHA256 hash from whitelist_signatures.
func (db *DB) RemoveWhitelist(sha256Hash string) error {
	_, err := db.conn.Exec("DELETE FROM whitelist_signatures WHERE sha256_hash = ?", sha256Hash)
	return err
}
