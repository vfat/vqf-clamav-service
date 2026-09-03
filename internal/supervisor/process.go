package supervisor

import (
	"context"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"
)

// Supervisor manages clamd and freshclam child processes.
type Supervisor struct {
	clamdCmd     *exec.Cmd
	freshclamCmd *exec.Cmd
	stopChan     chan struct{}
}

// NewSupervisor creates a process supervisor.
func NewSupervisor() *Supervisor {
	return &Supervisor{
		stopChan: make(chan struct{}),
	}
}

// Start boots freshclam seed (if needed), clamd daemon, and freshclam daemon.
func (s *Supervisor) Start(ctx context.Context, clamdConf, freshclamConf, dbDir, socketPath string) error {
	_ = os.MkdirAll(filepath.Dir(socketPath), 0755)
	_ = os.MkdirAll(dbDir, 0755)

	// 1. Check if virus signatures exist
	if !hasSignatures(dbDir) {
		log.Println("[SUPERVISOR] First run: Virus signatures missing. Running initial freshclam download (this may take 1-2 minutes)...")
		cmd := exec.CommandContext(ctx, "freshclam", "--config-file="+freshclamConf)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Printf("[SUPERVISOR WARNING] Initial freshclam run: %v", err)
		} else {
			log.Println("[SUPERVISOR] Initial virus signatures downloaded successfully.")
		}
	}

	// 2. Start clamd daemon
	go s.runClamd(clamdConf, socketPath)

	// 3. Start freshclam background updater
	go s.runFreshclam(freshclamConf)

	return nil
}

func hasSignatures(dir string) bool {
	files, err := os.ReadDir(dir)
	if err != nil {
		return false
	}
	for _, f := range files {
		ext := filepath.Ext(f.Name())
		if ext == ".cvd" || ext == ".cld" {
			return true
		}
	}
	return false
}

func (s *Supervisor) runClamd(confPath, socketPath string) {
	for {
		select {
		case <-s.stopChan:
			return
		default:
		}

		_ = os.Remove(socketPath)

		log.Printf("[SUPERVISOR] Starting clamd daemon with config %s...", confPath)
		cmd := exec.Command("clamd", "-c", confPath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		s.clamdCmd = cmd

		if err := cmd.Start(); err != nil {
			log.Printf("[SUPERVISOR ERROR] Failed to start clamd: %v. Retrying in 5s...", err)
			time.Sleep(5 * time.Second)
			continue
		}

		err := cmd.Wait()
		select {
		case <-s.stopChan:
			return
		default:
			log.Printf("[SUPERVISOR WARNING] clamd exited (%v). Restarting in 3s...", err)
			time.Sleep(3 * time.Second)
		}
	}
}

func (s *Supervisor) runFreshclam(confPath string) {
	// Delay initial daemon start slightly so clamd gets CPU priority
	time.Sleep(10 * time.Second)

	for {
		select {
		case <-s.stopChan:
			return
		default:
		}

		log.Printf("[SUPERVISOR] Starting freshclam daemon updater with config %s...", confPath)
		cmd := exec.Command("freshclam", "-d", "-c", confPath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		s.freshclamCmd = cmd

		if err := cmd.Start(); err != nil {
			log.Printf("[SUPERVISOR ERROR] Failed to start freshclam: %v. Retrying in 10s...", err)
			time.Sleep(10 * time.Second)
			continue
		}

		err := cmd.Wait()
		select {
		case <-s.stopChan:
			return
		default:
			log.Printf("[SUPERVISOR WARNING] freshclam exited (%v). Restarting in 10s...", err)
			time.Sleep(10 * time.Second)
		}
	}
}

// Stop terminates clamd and freshclam child processes.
func (s *Supervisor) Stop() {
	close(s.stopChan)
	if s.clamdCmd != nil && s.clamdCmd.Process != nil {
		_ = s.clamdCmd.Process.Signal(syscall.SIGTERM)
	}
	if s.freshclamCmd != nil && s.freshclamCmd.Process != nil {
		_ = s.freshclamCmd.Process.Signal(syscall.SIGTERM)
	}
}
