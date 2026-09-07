package fetcher

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	ErrSSRFBlocked      = errors.New("SSRF attempt blocked: target IP or domain is prohibited")
	ErrInvalidScheme    = errors.New("invalid URL scheme: only http and https are allowed")
	ErrTooManyRedirects = errors.New("too many redirects: maximum hops exceeded")
	ErrPayloadTooLarge  = errors.New("remote file exceeds maximum scanning size limit")
)

var blockedCIDRs []*net.IPNet

func init() {
	cidrs := []string{
		// IPv4 Private (RFC 1918)
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		// IPv4 Loopback (RFC 1122)
		"127.0.0.0/8",
		// IPv4 Link-Local / Cloud Metadata (RFC 3927)
		"169.254.0.0/16",
		// Carrier-Grade NAT (RFC 6598)
		"100.64.0.0/10",
		// Broadcast & Wildcard
		"0.0.0.0/8",
		"240.0.0.0/4",
		"255.255.255.255/32",
		// Multicast
		"224.0.0.0/4",
		// Benchmark
		"198.18.0.0/15",

		// IPv6 Unique Local Address (RFC 4193)
		"fc00::/7",
		// IPv6 Link-Local Unicast (RFC 4291)
		"fe80::/10",
		// IPv6 Loopback
		"::1/128",
		// IPv6 Unspecified
		"::/128",
		// IPv6 Multicast
		"ff00::/8",
	}

	for _, c := range cidrs {
		_, ipNet, err := net.ParseCIDR(c)
		if err == nil {
			blockedCIDRs = append(blockedCIDRs, ipNet)
		}
	}
}

// IsBlockedIP checks whether an IP address belongs to loopback, private,
// link-local, cloud metadata, or CGNAT networks.
func IsBlockedIP(ip net.IP) bool {
	if ip == nil {
		return true
	}

	// Unwrap IPv4-mapped IPv6 addresses (e.g. ::ffff:127.0.0.1)
	if v4 := ip.To4(); v4 != nil {
		ip = v4
	}

	// Native Go IP checks
	if ip.IsLoopback() || ip.IsUnspecified() || ip.IsLinkLocalUnicast() ||
		ip.IsLinkLocalMulticast() || ip.IsMulticast() || ip.IsPrivate() {
		return true
	}

	// Explicit CIDR check (covers CGNAT 100.64.0.0/10 and Cloud Metadata 169.254.0.0/16)
	for _, network := range blockedCIDRs {
		if network.Contains(ip) {
			return true
		}
	}

	return false
}

// ValidateURL ensures the target URL uses http or https schemes and has a valid host.
func ValidateURL(rawURL string) error {
	if strings.TrimSpace(rawURL) == "" {
		return fmt.Errorf("%w: empty URL", ErrInvalidScheme)
	}

	parsed, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidScheme, err)
	}

	scheme := strings.ToLower(parsed.Scheme)
	if scheme != "http" && scheme != "https" {
		return fmt.Errorf("%w: %s", ErrInvalidScheme, scheme)
	}

	if parsed.Host == "" {
		return fmt.Errorf("%w: missing host", ErrInvalidScheme)
	}

	return nil
}

// SafeFetcherConfig defines limits and behaviors for safe remote file downloads.
type SafeFetcherConfig struct {
	Timeout                 time.Duration
	MaxRedirects            int
	MaxBytes                int64
	AllowLoopbackForTesting bool
}

// SafeFetcher manages secure HTTP fetching with Anti-SSRF and streaming safeguards.
type SafeFetcher struct {
	config SafeFetcherConfig
	client *http.Client
}

// FetchResult contains the streaming reader and metadata for the retrieved remote file.
type FetchResult struct {
	Body          io.ReadCloser
	StatusCode    int
	ContentLength int64
	ContentType   string
	ResolvedIP    string
	FinalURL      string
}

// NewSafeFetcher creates a hardened HTTP client with strict Anti-SSRF controls.
func NewSafeFetcher(cfg SafeFetcherConfig) *SafeFetcher {
	if cfg.Timeout <= 0 {
		cfg.Timeout = 30 * time.Second
	}
	if cfg.MaxRedirects <= 0 {
		cfg.MaxRedirects = 3
	}
	if cfg.MaxBytes <= 0 {
		cfg.MaxBytes = 100 * 1024 * 1024 // 100 MB default
	}

	fetcher := &SafeFetcher{
		config: cfg,
	}

	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			host, port, err := net.SplitHostPort(addr)
			if err != nil {
				return nil, err
			}

			// Validate IP addresses before connecting
			targetIPs, err := net.DefaultResolver.LookupIP(ctx, "ip", host)
			if err != nil {
				// If host is an IP literal, try parsing directly
				parsedIP := net.ParseIP(host)
				if parsedIP == nil {
					return nil, err
				}
				targetIPs = []net.IP{parsedIP}
			}

			var dialIP net.IP
			for _, ip := range targetIPs {
				isBlocked := IsBlockedIP(ip)
				if cfg.AllowLoopbackForTesting && ip.IsLoopback() {
					isBlocked = false
				}

				if isBlocked {
					return nil, fmt.Errorf("%w: %s resolves to blocked address %s", ErrSSRFBlocked, host, ip.String())
				}
				if dialIP == nil {
					dialIP = ip
				}
			}

			if dialIP == nil {
				return nil, fmt.Errorf("%w: no valid resolved IP for %s", ErrSSRFBlocked, host)
			}

			dialer := &net.Dialer{
				Timeout:   cfg.Timeout,
				KeepAlive: 30 * time.Second,
			}

			// Connect directly to the verified IP to prevent DNS rebinding
			return dialer.DialContext(ctx, network, net.JoinHostPort(dialIP.String(), port))
		},
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   cfg.Timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= cfg.MaxRedirects {
				return fmt.Errorf("%w: limit %d", ErrTooManyRedirects, cfg.MaxRedirects)
			}

			// Validate redirect URL scheme
			if err := ValidateURL(req.URL.String()); err != nil {
				return err
			}

			// Validate target IP on redirect hop
			host := req.URL.Hostname()
			targetIPs, err := net.DefaultResolver.LookupIP(req.Context(), "ip", host)
			if err != nil {
				parsedIP := net.ParseIP(host)
				if parsedIP == nil {
					return err
				}
				targetIPs = []net.IP{parsedIP}
			}

			for _, ip := range targetIPs {
				isBlocked := IsBlockedIP(ip)
				if cfg.AllowLoopbackForTesting && ip.IsLoopback() {
					isBlocked = false
				}
				if isBlocked {
					return fmt.Errorf("%w: redirect target %s resolves to blocked IP %s", ErrSSRFBlocked, host, ip.String())
				}
			}

			return nil
		},
	}

	fetcher.client = client
	return fetcher
}

// Fetch executes a hardened GET request, enforcing Anti-SSRF, timeouts, and byte limits.
func (f *SafeFetcher) Fetch(ctx context.Context, targetURL string) (*FetchResult, error) {
	if err := ValidateURL(targetURL); err != nil {
		return nil, err
	}

	parsed, err := url.Parse(targetURL)
	if err != nil {
		return nil, err
	}

	// Pre-dial verification of initial host
	host := parsed.Hostname()
	targetIPs, err := net.DefaultResolver.LookupIP(ctx, "ip", host)
	if err != nil {
		parsedIP := net.ParseIP(host)
		if parsedIP == nil {
			return nil, fmt.Errorf("%w: failed resolving host %s: %v", ErrSSRFBlocked, host, err)
		}
		targetIPs = []net.IP{parsedIP}
	}

	var firstAllowedIP string
	for _, ip := range targetIPs {
		isBlocked := IsBlockedIP(ip)
		if f.config.AllowLoopbackForTesting && ip.IsLoopback() {
			isBlocked = false
		}
		if isBlocked {
			return nil, fmt.Errorf("%w: host %s resolves to prohibited address %s", ErrSSRFBlocked, host, ip.String())
		}
		if firstAllowedIP == "" {
			firstAllowedIP = ip.String()
		}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, targetURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "ClamAV-Remote-Scanner/1.1 (Anti-SSRF Protected)")

	resp, err := f.client.Do(req)
	if err != nil {
		if errors.Is(err, ErrSSRFBlocked) || strings.Contains(err.Error(), ErrSSRFBlocked.Error()) {
			return nil, ErrSSRFBlocked
		}
		if errors.Is(err, ErrTooManyRedirects) || strings.Contains(err.Error(), ErrTooManyRedirects.Error()) {
			return nil, ErrTooManyRedirects
		}
		return nil, err
	}

	// Wrap body in safe bounded reader
	boundedBody := &boundedReadCloser{
		rc:       resp.Body,
		maxBytes: f.config.MaxBytes,
	}

	return &FetchResult{
		Body:          boundedBody,
		StatusCode:    resp.StatusCode,
		ContentLength: resp.ContentLength,
		ContentType:   resp.Header.Get("Content-Type"),
		ResolvedIP:    firstAllowedIP,
		FinalURL:      resp.Request.URL.String(),
	}, nil
}

// boundedReadCloser enforces MaxBytes and translates overflows to ErrPayloadTooLarge.
type boundedReadCloser struct {
	rc       io.ReadCloser
	maxBytes int64
	readSoFar int64
}

func (b *boundedReadCloser) Read(p []byte) (int, error) {
	if b.readSoFar >= b.maxBytes {
		return 0, ErrPayloadTooLarge
	}

	// Calculate remaining allowed quota
	remaining := b.maxBytes - b.readSoFar
	toRead := p
	if int64(len(toRead)) > remaining {
		toRead = toRead[:remaining]
	}

	n, err := b.rc.Read(toRead)
	b.readSoFar += int64(n)

	if err != nil {
		return n, err
	}

	if b.readSoFar >= b.maxBytes {
		// Try to read 1 extra byte to check if stream truly continues beyond maxBytes
		extra := make([]byte, 1)
		extraN, extraErr := b.rc.Read(extra)
		if extraN > 0 {
			return n, ErrPayloadTooLarge
		}
		if extraErr != nil && !errors.Is(extraErr, io.EOF) {
			return n, extraErr
		}
	}

	return n, nil
}

func (b *boundedReadCloser) Close() error {
	return b.rc.Close()
}

// Client returns the underlying hardened http.Client configured with Anti-SSRF transport.
func (f *SafeFetcher) Client() *http.Client {
	return f.client
}

// ValidateTargetIP verifies that rawURL parses cleanly, uses http/https, and does not resolve to a blocked IP.
func (f *SafeFetcher) ValidateTargetIP(ctx context.Context, rawURL string) error {
	if err := ValidateURL(rawURL); err != nil {
		return err
	}

	parsed, err := url.Parse(rawURL)
	if err != nil {
		return err
	}

	host := parsed.Hostname()
	targetIPs, err := net.DefaultResolver.LookupIP(ctx, "ip", host)
	if err != nil {
		parsedIP := net.ParseIP(host)
		if parsedIP == nil {
			return fmt.Errorf("%w: failed resolving host %s: %v", ErrSSRFBlocked, host, err)
		}
		targetIPs = []net.IP{parsedIP}
	}

	for _, ip := range targetIPs {
		isBlocked := IsBlockedIP(ip)
		if f.config.AllowLoopbackForTesting && ip.IsLoopback() {
			isBlocked = false
		}
		if isBlocked {
			return fmt.Errorf("%w: host %s resolves to prohibited address %s", ErrSSRFBlocked, host, ip.String())
		}
	}

	return nil
}

