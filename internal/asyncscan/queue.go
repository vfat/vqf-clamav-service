package asyncscan

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/vfat/vqf-clamav-service/internal/alert"
	"github.com/vfat/vqf-clamav-service/internal/clamd"
	"github.com/vfat/vqf-clamav-service/internal/fetcher"
	"github.com/vfat/vqf-clamav-service/internal/quarantine"
	"github.com/vfat/vqf-clamav-service/internal/storage"
)

// clamdScanner defines the interface to scan stream with clamd.
type clamdScanner interface {
	ScanStream(ctx context.Context, r io.Reader) (*clamd.ScanResult, error)
}

// Config holds configuration parameters for the async queue.
type Config struct {
	SpoolDir      string
	Workers       int
	QueueCapacity int
	QuarRetention int
}

// WebhookPayload represents the standard JSON payload sent to callback_url.
type WebhookPayload struct {
	Event   string      `json:"event"`
	JobID   string      `json:"job_id"`
	Status  string      `json:"status"`
	Verdict string      `json:"verdict"`
	Threat  *ThreatInfo `json:"threat,omitempty"`
	Data    WebhookData `json:"data"`
}

// ThreatInfo details detected malware.
type ThreatInfo struct {
	VirusName    string `json:"virus_name"`
	ActionTaken  string `json:"action_taken"`
	QuarantineID string `json:"quarantine_id,omitempty"`
}

// WebhookData holds file and scan metadata.
type WebhookData struct {
	FileName       string `json:"file_name"`
	FileSize       int64  `json:"file_size"`
	FileSHA256     string `json:"file_sha256"`
	ScanDurationMs int64  `json:"scan_duration_ms"`
	CompletedAt    string `json:"completed_at"`
}

// Queue manages background scanning jobs, spool storage, and webhook delivery.
type Queue struct {
	cfg      Config
	db       *storage.DB
	scanner  clamdScanner
	vault    *quarantine.Vault
	notifier *alert.Notifier
	fetcher  *fetcher.SafeFetcher
	jobsChan chan *storage.ScanJob
	wg       sync.WaitGroup
	ctx      context.Context
	cancel   context.CancelFunc
}

// NewQueue constructs a new Queue instance.
func NewQueue(cfg Config, db *storage.DB, scanner clamdScanner, vault *quarantine.Vault, notifier *alert.Notifier, f *fetcher.SafeFetcher) *Queue {
	if cfg.SpoolDir == "" {
		cfg.SpoolDir = "/data/spool"
	}
	if cfg.Workers <= 0 {
		cfg.Workers = 2
	}
	if cfg.QueueCapacity <= 0 {
		cfg.QueueCapacity = 100
	}
	if cfg.QuarRetention <= 0 {
		cfg.QuarRetention = 7
	}
	if f == nil {
		f = fetcher.NewSafeFetcher(fetcher.SafeFetcherConfig{
			Timeout: 10 * time.Second,
		})
	}

	return &Queue{
		cfg:      cfg,
		db:       db,
		scanner:  scanner,
		vault:    vault,
		notifier: notifier,
		fetcher:  f,
		jobsChan: make(chan *storage.ScanJob, cfg.QueueCapacity),
	}
}

// Start launches worker goroutines.
func (q *Queue) Start(ctx context.Context) {
	q.ctx, q.cancel = context.WithCancel(ctx)

	if err := os.MkdirAll(q.cfg.SpoolDir, 0700); err != nil {
		log.Printf("[ASYNC] Failed to create spool dir %s: %v", q.cfg.SpoolDir, err)
	}

	for i := 0; i < q.cfg.Workers; i++ {
		q.wg.Add(1)
		go q.worker(i)
	}
	log.Printf("[ASYNC] Started %d scan worker goroutines", q.cfg.Workers)
}

// Stop terminates worker goroutines and waits for pending jobs to finish.
func (q *Queue) Stop() {
	if q.cancel != nil {
		q.cancel()
	}
	close(q.jobsChan)
	q.wg.Wait()
	log.Printf("[ASYNC] All scan workers stopped gracefully")
}

// SubmitJob writes file payload to spool directory, records job in SQLite, and enqueues to channel.
func (q *Queue) SubmitJob(ctx context.Context, fileName string, fileStream io.Reader, fileSize int64, callbackURL, consumer string) (*storage.ScanJob, error) {
	callbackURL = strings.TrimSpace(callbackURL)
	if callbackURL == "" {
		return nil, errors.New("callback_url is required")
	}

	// Validate callback URL with Anti-SSRF
	if err := fetcher.ValidateURL(callbackURL); err != nil {
		return nil, fmt.Errorf("invalid callback_url: %w", err)
	}

	// Pre-validate callback host IP to reject SSRF early if host is raw IP or resolvable
	if err := q.validateWebhookURL(ctx, callbackURL); err != nil {
		return nil, err
	}

	if err := os.MkdirAll(q.cfg.SpoolDir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create spool directory: %w", err)
	}

	jobID := fmt.Sprintf("job_%d", time.Now().UnixNano())
	spoolPath := filepath.Join(q.cfg.SpoolDir, jobID+".spool")

	spoolFile, err := os.OpenFile(spoolPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return nil, fmt.Errorf("failed to create spool file: %w", err)
	}

	written, err := io.Copy(spoolFile, fileStream)
	_ = spoolFile.Close()
	if err != nil {
		_ = os.Remove(spoolPath)
		return nil, fmt.Errorf("failed writing to spool file: %w", err)
	}

	if fileSize <= 0 {
		fileSize = written
	}

	now := time.Now().UTC()
	job := &storage.ScanJob{
		ID:          jobID,
		FileName:    fileName,
		FileSize:    fileSize,
		CallbackURL: callbackURL,
		Consumer:    consumer,
		Status:      "QUEUED",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := q.db.InsertScanJob(*job); err != nil {
		_ = os.Remove(spoolPath)
		return nil, fmt.Errorf("failed to record scan job in database: %w", err)
	}

	select {
	case q.jobsChan <- job:
	default:
		// Queue buffer full: job is still in DB, will return error or wait
		log.Printf("[ASYNC WARN] Job channel full for job %s", jobID)
	}

	return job, nil
}

func (q *Queue) validateWebhookURL(ctx context.Context, rawURL string) error {
	if q.fetcher == nil {
		return fetcher.ValidateURL(rawURL)
	}
	return q.fetcher.ValidateTargetIP(ctx, rawURL)
}


func (q *Queue) worker(workerID int) {
	defer q.wg.Done()

	for job := range q.jobsChan {
		q.processJob(job)
	}
}

func (q *Queue) processJob(job *storage.ScanJob) {
	spoolPath := filepath.Join(q.cfg.SpoolDir, job.ID+".spool")
	defer func() {
		_ = os.Remove(spoolPath)
	}()

	_ = q.db.UpdateScanJobStatus(job.ID, "PROCESSING", "", "", "", "")

	file, err := os.Open(spoolPath)
	if err != nil {
		_ = q.db.UpdateScanJobStatus(job.ID, "FAILED", "ERROR", "", "Failed to open spool file: "+err.Error(), "")
		q.dispatchWebhook(job, WebhookPayload{
			Event:   "scan.failed",
			JobID:   job.ID,
			Status:  "FAILED",
			Verdict: "ERROR",
			Data: WebhookData{
				FileName: job.FileName,
				FileSize: job.FileSize,
			},
		})
		return
	}
	defer file.Close()

	// Compute SHA256 while streaming to clamd
	hasher := sha256.New()
	teeReader := io.TeeReader(file, hasher)

	startTime := time.Now()
	scanCtx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	scanRes, err := q.scanner.ScanStream(scanCtx, teeReader)
	if err != nil {
		_ = q.db.UpdateScanJobStatus(job.ID, "FAILED", "ERROR", "", "Daemon scan error: "+err.Error(), "")
		return
	}

	durationMs := time.Since(startTime).Milliseconds()
	fileSHA256 := hex.EncodeToString(hasher.Sum(nil))

	if scanRes.IsClean() {
		_ = file.Close()
		_ = os.Remove(spoolPath)

		_ = q.db.UpdateScanJobStatus(job.ID, "COMPLETED", "CLEAN", "", "", fileSHA256)
		_ = q.db.InsertScanAuditLog(storage.ScanAuditLog{
			ID:             fmt.Sprintf("audit_%d", time.Now().UnixNano()),
			Timestamp:      time.Now().UTC(),
			ConsumerName:   job.Consumer,
			FileName:       job.FileName,
			FileSizeBytes:  job.FileSize,
			FileSHA256:     fileSHA256,
			Verdict:        "CLEAN",
			ScanDurationMs: durationMs,
		})

		q.dispatchWebhook(job, WebhookPayload{
			Event:   "scan.completed",
			JobID:   job.ID,
			Status:  "COMPLETED",
			Verdict: "CLEAN",
			Data: WebhookData{
				FileName:       job.FileName,
				FileSize:       job.FileSize,
				FileSHA256:     fileSHA256,
				ScanDurationMs: durationMs,
				CompletedAt:    time.Now().UTC().Format(time.RFC3339),
			},
		})
		return
	}

	// INFECTED
	quarID := ""
	if q.vault != nil {
		_, _ = file.Seek(0, io.SeekStart)
		quarRec, _ := q.vault.QuarantineFile(context.Background(), job.FileName, job.Consumer, scanRes.VirusName, file, q.cfg.QuarRetention)
		if quarRec != nil {
			quarID = quarRec.ID
		}
	}

	_ = file.Close()
	_ = os.Remove(spoolPath)

	_ = q.db.UpdateScanJobStatus(job.ID, "COMPLETED", "INFECTED", scanRes.VirusName, "", fileSHA256)
	_ = q.db.InsertScanAuditLog(storage.ScanAuditLog{
		ID:             fmt.Sprintf("audit_%d", time.Now().UnixNano()),
		Timestamp:      time.Now().UTC(),
		ConsumerName:   job.Consumer,
		FileName:       job.FileName,
		FileSizeBytes:  job.FileSize,
		FileSHA256:     fileSHA256,
		Verdict:        "INFECTED",
		VirusName:      scanRes.VirusName,
		ScanDurationMs: durationMs,
		QuarantineID:   quarID,
	})

	if q.notifier != nil {
		_ = q.notifier.DispatchThreat(context.Background(), alert.ThreatAlert{
			VirusName:      scanRes.VirusName,
			FileName:       job.FileName,
			FileSizeBytes:  job.FileSize,
			FileSHA256:     fileSHA256,
			QuarantineID:   quarID,
			SourceConsumer: job.Consumer,
			DetectedAt:     time.Now().UTC(),
		})
	}

	q.dispatchWebhook(job, WebhookPayload{
		Event:   "scan.completed",
		JobID:   job.ID,
		Status:  "COMPLETED",
		Verdict: "INFECTED",
		Threat: &ThreatInfo{
			VirusName:    scanRes.VirusName,
			ActionTaken:  "QUARANTINED",
			QuarantineID: quarID,
		},
		Data: WebhookData{
			FileName:       job.FileName,
			FileSize:       job.FileSize,
			FileSHA256:     fileSHA256,
			ScanDurationMs: durationMs,
			CompletedAt:    time.Now().UTC().Format(time.RFC3339),
		},
	})
}


func (q *Queue) dispatchWebhook(job *storage.ScanJob, payload WebhookPayload) {
	if job.CallbackURL == "" {
		return
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return
	}

	// Safe webhook POST with Anti-SSRF
	client := q.fetcher.Client()

	for attempt := 1; attempt <= 3; attempt++ {
		reqCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		req, err := http.NewRequestWithContext(reqCtx, http.MethodPost, job.CallbackURL, bytes.NewReader(bodyBytes))
		if err != nil {
			cancel()
			return
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", "ClamAV-Service-Webhook/1.0")

		resp, err := client.Do(req)
		if err == nil {
			_ = resp.Body.Close()
			cancel()
			if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				return
			}
		} else {
			cancel()
			if errors.Is(err, fetcher.ErrSSRFBlocked) {
				log.Printf("[ASYNC SSRF BLOCKED] Webhook to %s blocked by Anti-SSRF", job.CallbackURL)
				return
			}
		}

		time.Sleep(100 * time.Millisecond)
	}
}
