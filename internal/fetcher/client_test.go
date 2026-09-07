package fetcher

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestIsBlockedIP(t *testing.T) {
	tests := []struct {
		ip      string
		blocked bool
	}{
		// IPv4 Loopback
		{"127.0.0.1", true},
		{"127.0.0.2", true},
		{"127.255.255.255", true},

		// IPv6 Loopback
		{"::1", true},

		// Wildcard / Unspecified
		{"0.0.0.0", true},
		{"::", true},

		// Link-Local / Cloud Metadata
		{"169.254.169.254", true},
		{"169.254.1.1", true},
		{"fe80::1", true},

		// RFC1918 Private IPv4
		{"10.0.0.1", true},
		{"10.255.255.254", true},
		{"172.16.0.1", true},
		{"172.31.255.254", true},
		{"192.168.1.1", true},
		{"192.168.254.254", true},

		// Carrier-Grade NAT (RFC6598: 100.64.0.0/10)
		{"100.64.0.1", true},
		{"100.127.255.254", true},

		// IPv6 Unique Local Address (fc00::/7)
		{"fc00::1", true},
		{"fd12:3456::1", true},

		// IPv4-mapped IPv6
		{"::ffff:127.0.0.1", true},
		{"::ffff:169.254.169.254", true},
		{"::ffff:10.0.0.1", true},
		{"::ffff:192.168.1.1", true},

		// Allowed Public IPs
		{"8.8.8.8", false},
		{"1.1.1.1", false},
		{"52.216.1.1", false},
		{"142.250.190.46", false},
		{"100.63.255.255", false}, // Outside CGNAT
		{"100.128.0.1", false},    // Outside CGNAT
		{"172.15.255.255", false}, // Outside RFC1918
		{"172.32.0.1", false},     // Outside RFC1918
		{"2607:f8b0:4005:809::200e", false}, // Google Public IPv6
	}

	for _, tt := range tests {
		ip := net.ParseIP(tt.ip)
		if ip == nil {
			t.Fatalf("failed to parse IP: %s", tt.ip)
		}
		got := IsBlockedIP(ip)
		if got != tt.blocked {
			t.Errorf("IsBlockedIP(%s) = %v; want %v", tt.ip, got, tt.blocked)
		}
	}
}

func TestValidateURLScheme(t *testing.T) {
	tests := []struct {
		url     string
		wantErr bool
	}{
		{"http://example.com/file.pdf", false},
		{"https://s3.amazonaws.com/bucket/doc.exe", false},
		{"ftp://example.com/file.pdf", true},
		{"file:///etc/passwd", true},
		{"gopher://evil.com", true},
		{"data:text/plain;base64,SGVsbG8=", true},
		{"", true},
		{"://invalid-url", true},
	}

	for _, tt := range tests {
		err := ValidateURL(tt.url)
		if (err != nil) != tt.wantErr {
			t.Errorf("ValidateURL(%q) err = %v; wantErr = %v", tt.url, err, tt.wantErr)
		}
	}
}

func TestSSRFBlockingDirect(t *testing.T) {
	fetcher := NewSafeFetcher(SafeFetcherConfig{
		Timeout: 5 * time.Second,
	})

	blockedTargets := []string{
		"http://127.0.0.1:8080/secret",
		"http://localhost:8080/admin",
		"http://169.254.169.254/latest/meta-data/",
		"http://10.0.0.1/config",
		"http://172.16.0.1/internal",
		"http://192.168.1.1/router",
		"http://100.64.0.1/cgnat",
		"http://[::1]:8080/ipv6-local",
	}

	ctx := context.Background()
	for _, target := range blockedTargets {
		res, err := fetcher.Fetch(ctx, target)
		if err == nil {
			if res != nil && res.Body != nil {
				res.Body.Close()
			}
			t.Fatalf("expected SSRF error for target %s, got success", target)
		}
		if !errors.Is(err, ErrSSRFBlocked) {
			t.Errorf("expected ErrSSRFBlocked for %s, got: %v", target, err)
		}
	}
}

func TestRedirectAntiSSRF(t *testing.T) {
	// Server that redirects to an internal IP
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "http://169.254.169.254/latest/meta-data", http.StatusFound)
	}))
	defer ts.Close()

	// Configure fetcher allowing the test server loopback IP initially,
	// but redirect must still trigger the anti-SSRF check on the new destination.
	fetcher := NewSafeFetcher(SafeFetcherConfig{
		Timeout:                 5 * time.Second,
		AllowLoopbackForTesting: true, // initial request reaches httptest server
	})

	ctx := context.Background()
	res, err := fetcher.Fetch(ctx, ts.URL)
	if err == nil {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
		t.Fatal("expected redirect to 169.254.169.254 to be blocked by SSRF check, got nil error")
	}

	if !errors.Is(err, ErrSSRFBlocked) {
		t.Errorf("expected ErrSSRFBlocked on redirect, got %v", err)
	}
}

func TestRedirectMaxHops(t *testing.T) {
	var hopCount int
	var ts *httptest.Server
	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hopCount++
		http.Redirect(w, r, fmt.Sprintf("%s/hop/%d", ts.URL, hopCount), http.StatusFound)
	}))
	defer ts.Close()

	fetcher := NewSafeFetcher(SafeFetcherConfig{
		Timeout:                 5 * time.Second,
		MaxRedirects:            3,
		AllowLoopbackForTesting: true,
	})

	ctx := context.Background()
	res, err := fetcher.Fetch(ctx, ts.URL)
	if err == nil {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
		t.Fatal("expected error due to redirect loop/limit, got nil")
	}

	if !errors.Is(err, ErrTooManyRedirects) {
		t.Errorf("expected ErrTooManyRedirects, got: %v", err)
	}
}

func TestFetchSuccessAndChecksum(t *testing.T) {
	payload := []byte("This is safe remote file content designed for clean antivirus scanning!")
	hasher := sha256.New()
	hasher.Write(payload)
	expectedSHA256 := hex.EncodeToString(hasher.Sum(nil))

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/octet-stream")
		w.WriteHeader(http.StatusOK)
		w.Write(payload)
	}))
	defer ts.Close()

	fetcher := NewSafeFetcher(SafeFetcherConfig{
		Timeout:                 5 * time.Second,
		MaxBytes:                1024 * 1024,
		AllowLoopbackForTesting: true,
	})

	ctx := context.Background()
	res, err := fetcher.Fetch(ctx, ts.URL)
	if err != nil {
		t.Fatalf("unexpected fetch error: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", res.StatusCode)
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("failed reading body: %v", err)
	}

	if string(bodyBytes) != string(payload) {
		t.Errorf("body mismatch: got %q, want %q", string(bodyBytes), string(payload))
	}

	h := sha256.New()
	h.Write(bodyBytes)
	actualSHA256 := hex.EncodeToString(h.Sum(nil))
	if actualSHA256 != expectedSHA256 {
		t.Errorf("sha256 mismatch: got %s, want %s", actualSHA256, expectedSHA256)
	}
}

func TestFetchPayloadTooLarge(t *testing.T) {
	largeData := make([]byte, 5000)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(largeData)
	}))
	defer ts.Close()

	fetcher := NewSafeFetcher(SafeFetcherConfig{
		Timeout:                 5 * time.Second,
		MaxBytes:                1024, // 1 KB max
		AllowLoopbackForTesting: true,
	})

	ctx := context.Background()
	res, err := fetcher.Fetch(ctx, ts.URL)
	if err != nil {
		t.Fatalf("unexpected fetch start error: %v", err)
	}
	defer res.Body.Close()

	buf := make([]byte, 2048)
	_, readErr := io.ReadAll(res.Body)
	_ = buf
	if readErr == nil {
		t.Fatal("expected ErrPayloadTooLarge when exceeding MaxBytes, got nil")
	}

	if !errors.Is(readErr, ErrPayloadTooLarge) {
		t.Errorf("expected ErrPayloadTooLarge, got %v", readErr)
	}
}
