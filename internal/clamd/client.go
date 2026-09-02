package clamd

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"time"
)

const (
	defaultChunkSize = 64 * 1024 // 64 KB
)

// ScanResult represents the parsed outcome of a ClamAV scan.
type ScanResult struct {
	Verdict     string `json:"verdict"` // CLEAN, INFECTED, ERROR
	VirusName   string `json:"virus_name,omitempty"`
	RawResponse string `json:"raw_response"`
}

// IsClean returns true if no virus was detected.
func (r *ScanResult) IsClean() bool {
	return r.Verdict == "CLEAN"
}

// Client communicates with ClamAV daemon via Unix domain socket or TCP.
type Client struct {
	network string
	address string
	dialer  net.Dialer
}

// NewClient creates a new ClamAV socket client.
func NewClient(network, address string) *Client {
	if network == "" {
		network = "unix"
	}
	if address == "" {
		address = "/var/run/clamav/clamd.ctl"
	}

	return &Client{
		network: network,
		address: address,
		dialer:  net.Dialer{Timeout: 5 * time.Second},
	}
}

func (c *Client) dial(ctx context.Context) (net.Conn, error) {
	return c.dialer.DialContext(ctx, c.network, c.address)
}

// Ping sends a PING command to clamd and expects PONG.
func (c *Client) Ping(ctx context.Context) error {
	conn, err := c.dial(ctx)
	if err != nil {
		return fmt.Errorf("failed to connect to clamd at %s://%s: %w", c.network, c.address, err)
	}
	defer conn.Close()

	if _, err := conn.Write([]byte("zPING\x00")); err != nil {
		return fmt.Errorf("failed to send PING: %w", err)
	}

	resp, err := readUntilNull(conn)
	if err != nil {
		return fmt.Errorf("failed to read PING response: %w", err)
	}

	if !strings.HasPrefix(resp, "PONG") {
		return fmt.Errorf("unexpected PING response from clamd: %s", resp)
	}

	return nil
}

// Version queries clamd signature version string.
func (c *Client) Version(ctx context.Context) (string, error) {
	conn, err := c.dial(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to connect to clamd: %w", err)
	}
	defer conn.Close()

	if _, err := conn.Write([]byte("zVERSION\x00")); err != nil {
		return "", fmt.Errorf("failed to send VERSION: %w", err)
	}

	resp, err := readUntilNull(conn)
	if err != nil {
		return "", fmt.Errorf("failed to read VERSION response: %w", err)
	}

	return strings.TrimSpace(resp), nil
}

// ScanStream pipes an io.Reader stream directly to clamd using the INSTREAM protocol.
func (c *Client) ScanStream(ctx context.Context, r io.Reader) (*ScanResult, error) {
	conn, err := c.dial(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to clamd: %w", err)
	}
	defer conn.Close()

	// Send INSTREAM command
	if _, err := conn.Write([]byte("zINSTREAM\x00")); err != nil {
		return nil, fmt.Errorf("failed to initiate zINSTREAM command: %w", err)
	}

	// Stream chunks
	buf := make([]byte, defaultChunkSize)
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		n, readErr := r.Read(buf)
		if n > 0 {
			// Write 4-byte big-endian chunk length
			if err := binary.Write(conn, binary.BigEndian, uint32(n)); err != nil {
				return nil, fmt.Errorf("failed to write chunk size: %w", err)
			}
			// Write chunk payload
			if _, err := conn.Write(buf[:n]); err != nil {
				return nil, fmt.Errorf("failed to write chunk payload: %w", err)
			}
		}

		if readErr != nil {
			if errors.Is(readErr, io.EOF) {
				break
			}
			return nil, fmt.Errorf("error reading stream payload: %w", readErr)
		}
	}

	// Send 4-byte zero length to terminate stream
	if err := binary.Write(conn, binary.BigEndian, uint32(0)); err != nil {
		return nil, fmt.Errorf("failed to send stream terminator: %w", err)
	}

	// Read scan response
	resp, err := readUntilNull(conn)
	if err != nil {
		return nil, fmt.Errorf("failed to read clamd scan response: %w", err)
	}

	return parseScanResponse(resp)
}

func parseScanResponse(raw string) (*ScanResult, error) {
	trimmed := strings.TrimSpace(raw)
	res := &ScanResult{RawResponse: trimmed}

	if strings.Contains(trimmed, "OK") {
		res.Verdict = "CLEAN"
		return res, nil
	}

	if strings.Contains(trimmed, "FOUND") {
		res.Verdict = "INFECTED"
		// Response format: "stream: <VirusName> FOUND"
		parts := strings.Split(trimmed, ":")
		if len(parts) >= 2 {
			subParts := strings.TrimSpace(parts[1])
			res.VirusName = strings.TrimSuffix(subParts, " FOUND")
		}
		return res, nil
	}

	if strings.Contains(trimmed, "ERROR") {
		res.Verdict = "ERROR"
		return res, fmt.Errorf("clamd scan error: %s", trimmed)
	}

	return res, fmt.Errorf("unrecognized clamd response: %s", trimmed)
}

func readUntilNull(r io.Reader) (string, error) {
	var buf bytes.Buffer
	b := make([]byte, 1)

	for {
		n, err := r.Read(b)
		if n > 0 {
			if b[0] == 0 || b[0] == '\n' {
				break
			}
			buf.WriteByte(b[0])
		}
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return buf.String(), err
		}
	}

	return buf.String(), nil
}
