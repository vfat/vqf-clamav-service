package asyncscan

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/vfat/vqf-clamav-service/internal/alert"
	"github.com/vfat/vqf-clamav-service/internal/clamd"
	"github.com/vfat/vqf-clamav-service/internal/fetcher"
	"github.com/vfat/vqf-clamav-service/internal/quarantine"
	"github.com/vfat/vqf-clamav-service/internal/storage"
)

type mockClamdScanner struct {
	clean       bool
	virusName   string
	scanStreamErr error
}

func (m *mockClamdScanner) ScanStream(ctx context.Context, r io.Reader) (*clamd.ScanResult, error) {
	if m.scanStreamErr != nil {
		return nil, m.scanStreamErr
	}
	// Consume stream
	_, _ = io.ReadAll(r)
	if m.clean {
		return &clamd.ScanResult{Verdict: "CLEAN", RawResponse: "stream: OK"}, nil
	}
	return &clamd.ScanResult{
		Verdict:     "INFECTED",
		RawResponse: fmt.Sprintf("stream: %s FOUND", m.virusName),
		VirusName:   m.virusName,
	}, nil

}

func setupTestQueue(t *testing.T, scanner clamdScanner) (*Queue, *storage.DB, *quarantine.Vault, string) {
	t.Helper()
	tmpDir := t.TempDir()
	spoolDir := filepath.Join(tmpDir, "spool")
	vaultDir := filepath.Join(tmpDir, "vault")
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := storage.NewDB(dbPath)
	if err != nil {
		t.Fatalf("failed to init DB: %v", err)
	}

	vault := quarantine.NewVault(vaultDir, db)
	notifier := alert.NewNotifier(alert.Config{})

	safeFetcher := fetcher.NewSafeFetcher(fetcher.SafeFetcherConfig{
		Timeout:                 5 * time.Second,
		AllowLoopbackForTesting: true,
	})

	cfg := Config{
		SpoolDir:      spoolDir,
		Workers:       2,
		QueueCapacity: 50,
		QuarRetention: 7,
	}

	queue := NewQueue(cfg, db, scanner, vault, notifier, safeFetcher)
	return queue, db, vault, spoolDir
}

func TestQueue_SubmitAndProcess_CleanFile(t *testing.T) {
	mockScanner := &mockClamdScanner{clean: true}
	queue, db, _, spoolDir := setupTestQueue(t, mockScanner)
	defer db.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	queue.Start(ctx)
	defer queue.Stop()

	// Setup mock webhook receiver
	webhookReceived := make(chan WebhookPayload, 1)
	webhookServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload WebhookPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		webhookReceived <- payload
		w.WriteHeader(http.StatusOK)
	}))
	defer webhookServer.Close()

	fileData := []byte("clean file contents for async queue processing")
	job, err := queue.SubmitJob(ctx, "report.pdf", bytes.NewReader(fileData), int64(len(fileData)), webhookServer.URL, "client-app")
	if err != nil {
		t.Fatalf("SubmitJob failed: %v", err)
	}

	if job.ID == "" || job.Status != "QUEUED" {
		t.Fatalf("unexpected job response: %+v", job)
	}

	// Wait for webhook callback
	select {
	case payload := <-webhookReceived:
		if payload.JobID != job.ID {
			t.Errorf("expected job ID '%s', got '%s'", job.ID, payload.JobID)
		}
		if payload.Verdict != "CLEAN" {
			t.Errorf("expected verdict 'CLEAN', got '%s'", payload.Verdict)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("timed out waiting for webhook callback")
	}

	// Verify DB status updated to COMPLETED
	updatedJob, err := db.GetScanJob(job.ID)
	if err != nil {
		t.Fatalf("failed to fetch updated job: %v", err)
	}
	if updatedJob.Status != "COMPLETED" || updatedJob.Verdict != "CLEAN" {
		t.Errorf("expected status COMPLETED, got '%s'", updatedJob.Status)
	}

	// Verify spool file was deleted
	spoolFilePath := filepath.Join(spoolDir, job.ID+".spool")
	if _, err := os.Stat(spoolFilePath); !os.IsNotExist(err) {
		t.Errorf("expected spool file to be deleted from disk, but it still exists")
	}
}

func TestQueue_SubmitAndProcess_InfectedFile(t *testing.T) {
	mockScanner := &mockClamdScanner{clean: false, virusName: "Eicar-Test-Signature"}
	queue, db, vault, _ := setupTestQueue(t, mockScanner)
	defer db.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	queue.Start(ctx)
	defer queue.Stop()

	webhookReceived := make(chan WebhookPayload, 1)
	webhookServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload WebhookPayload
		_ = json.NewDecoder(r.Body).Decode(&payload)
		webhookReceived <- payload
		w.WriteHeader(http.StatusOK)
	}))
	defer webhookServer.Close()

	fileData := []byte("X5O!P%@AP[4\\PZX54(P^)7CC)7}$EICAR-STANDARD-ANTIVIRUS-TEST-FILE!$H+H*")
	job, err := queue.SubmitJob(ctx, "eicar.com", bytes.NewReader(fileData), int64(len(fileData)), webhookServer.URL, "client-app")
	if err != nil {
		t.Fatalf("SubmitJob failed: %v", err)
	}

	select {
	case payload := <-webhookReceived:
		if payload.JobID != job.ID {
			t.Errorf("expected job ID '%s', got '%s'", job.ID, payload.JobID)
		}
		if payload.Verdict != "INFECTED" {
			t.Errorf("expected verdict 'INFECTED', got '%s'", payload.Verdict)
		}

		if payload.Threat == nil || payload.Threat.VirusName != "Eicar-Test-Signature" {
			t.Errorf("expected threat Eicar-Test-Signature, got %+v", payload.Threat)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("timed out waiting for webhook callback for infected file")
	}

	// Verify quarantine record created
	records, total, _ := db.ListQuarantineRecords(10, 0, "")
	if total < 1 || len(records) < 1 {
		t.Errorf("expected quarantine record to be created in DB, got total %d", total)
	}
	_ = vault
}

func TestQueue_AntiSSRF_WebhookBlocked(t *testing.T) {
	mockScanner := &mockClamdScanner{clean: true}
	tmpDir := t.TempDir()
	db, _ := storage.NewDB(filepath.Join(tmpDir, "test.db"))
	defer db.Close()

	// Use standard SafeFetcher that disallows loopback & metadata IPs
	strictFetcher := fetcher.NewSafeFetcher(fetcher.SafeFetcherConfig{
		Timeout:                 2 * time.Second,
		AllowLoopbackForTesting: false,
	})

	queue := NewQueue(Config{
		SpoolDir:      filepath.Join(tmpDir, "spool"),
		Workers:       1,
		QueueCapacity: 10,
	}, db, mockScanner, nil, nil, strictFetcher)

	ctx := context.Background()

	// Submitting a private / cloud metadata webhook URL should be rejected immediately or fail safely
	_, err := queue.SubmitJob(ctx, "test.txt", bytes.NewReader([]byte("test")), 4, "http://169.254.169.254/secret", "attacker")
	if err == nil {
		t.Fatal("expected error submitting job with cloud metadata SSRF callback URL")
	}
}
