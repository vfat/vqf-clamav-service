package clamd

import (
	"bytes"
	"context"
	"encoding/binary"
	"io"
	"net"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// mockClamdServer sets up a temporary unix domain socket server simulating clamd responses
func mockClamdServer(t *testing.T, sockPath string) net.Listener {
	listener, err := net.Listen("unix", sockPath)
	if err != nil {
		t.Fatalf("failed to start mock unix socket listener: %v", err)
	}

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			go handleMockConn(conn)
		}
	}()

	return listener
}

func handleMockConn(conn net.Conn) {
	defer conn.Close()

	// Read command until null byte
	var cmdBuf bytes.Buffer
	b := make([]byte, 1)
	for {
		n, err := conn.Read(b)
		if n > 0 {
			if b[0] == 0 || b[0] == '\n' {
				break
			}
			cmdBuf.WriteByte(b[0])
		}
		if err != nil {
			return
		}
	}

	cmd := cmdBuf.String()

	switch {
	case strings.HasPrefix(cmd, "zPING") || strings.HasPrefix(cmd, "PING"):
		conn.Write([]byte("PONG\x00"))

	case strings.HasPrefix(cmd, "zVERSION") || strings.HasPrefix(cmd, "VERSION"):
		conn.Write([]byte("ClamAV 1.4.0/27400/Wed Sep  2 17:00:00 2026\x00"))

	case strings.HasPrefix(cmd, "zINSTREAM") || strings.HasPrefix(cmd, "nINSTREAM") || strings.HasPrefix(cmd, "INSTREAM"):
		// Read chunks until 4-byte 0 length chunk
		var payload bytes.Buffer
		for {
			var length uint32
			if err := binary.Read(conn, binary.BigEndian, &length); err != nil {
				break
			}
			if length == 0 {
				break // End of stream
			}
			chunk := make([]byte, length)
			if _, err := io.ReadFull(conn, chunk); err != nil {
				break
			}
			payload.Write(chunk)
		}

		if strings.Contains(payload.String(), "EICAR-STANDARD-ANTIVIRUS-TEST-FILE") {
			conn.Write([]byte("stream: Eicar-Test-Signature FOUND\x00"))
		} else {
			conn.Write([]byte("stream: OK\x00"))
		}
	}
}

func TestClamdClient_PingAndVersion(t *testing.T) {
	tmpDir := t.TempDir()
	sockPath := filepath.Join(tmpDir, "mock_clamd.sock")

	listener := mockClamdServer(t, sockPath)
	defer listener.Close()

	client := NewClient("unix", sockPath)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Test Ping
	if err := client.Ping(ctx); err != nil {
		t.Fatalf("Ping failed: %v", err)
	}

	// Test Version
	version, err := client.Version(ctx)
	if err != nil {
		t.Fatalf("Version failed: %v", err)
	}

	if !strings.Contains(version, "ClamAV 1.4.0") {
		t.Errorf("expected version to contain 'ClamAV 1.4.0', got: %s", version)
	}
}

func TestClamdClient_ScanStream_Clean(t *testing.T) {
	tmpDir := t.TempDir()
	sockPath := filepath.Join(tmpDir, "mock_clamd.sock")

	listener := mockClamdServer(t, sockPath)
	defer listener.Close()

	client := NewClient("unix", sockPath)

	cleanData := strings.NewReader("Hello world, this is a clean document.")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	result, err := client.ScanStream(ctx, cleanData)
	if err != nil {
		t.Fatalf("ScanStream failed on clean data: %v", err)
	}

	if !result.IsClean() {
		t.Errorf("expected clean result, got verdict: %s, virus: %s", result.Verdict, result.VirusName)
	}
}

func TestClamdClient_ScanStream_Infected(t *testing.T) {
	tmpDir := t.TempDir()
	sockPath := filepath.Join(tmpDir, "mock_clamd.sock")

	listener := mockClamdServer(t, sockPath)
	defer listener.Close()

	client := NewClient("unix", sockPath)

	eicarData := strings.NewReader("X5O!P%@AP[4\\PZX54(P^)7CC)7}$EICAR-STANDARD-ANTIVIRUS-TEST-FILE!$H+H*")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	result, err := client.ScanStream(ctx, eicarData)
	if err != nil {
		t.Fatalf("ScanStream failed on infected data: %v", err)
	}

	if result.IsClean() {
		t.Fatalf("expected infected result, got clean")
	}

	if result.Verdict != "INFECTED" {
		t.Errorf("expected verdict 'INFECTED', got '%s'", result.Verdict)
	}

	if result.VirusName != "Eicar-Test-Signature" {
		t.Errorf("expected virus name 'Eicar-Test-Signature', got '%s'", result.VirusName)
	}
}
