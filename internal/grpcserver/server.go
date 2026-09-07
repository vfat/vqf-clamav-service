package grpcserver

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	"github.com/vfat/vqf-clamav-service/internal/alert"
	"github.com/vfat/vqf-clamav-service/internal/clamd"
	"github.com/vfat/vqf-clamav-service/internal/quarantine"
	"github.com/vfat/vqf-clamav-service/internal/ratelimit"
	"github.com/vfat/vqf-clamav-service/internal/storage"
	clamavv1 "github.com/vfat/vqf-clamav-service/proto/clamav/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Config defines dependencies and tuning parameters for the gRPC server.
type Config struct {
	DB            *storage.DB
	Vault         *quarantine.Vault
	Notifier      *alert.Notifier
	Limiter       *ratelimit.Limiter
	Clamd         *clamd.Client
	AuthMode      string // "none", "bearer"
	BearerToken   string
	MaxScanSizeMB int64
	QuarRetention int
}

// Server implements the gRPC ScannerService server.
type Server struct {
	clamavv1.UnimplementedScannerServiceServer
	config     Config
	grpcServer *grpc.Server
}

// NewServer initializes and configures a new gRPC Server instance with interceptors.
func NewServer(cfg Config) *Server {
	if cfg.MaxScanSizeMB <= 0 {
		cfg.MaxScanSizeMB = 100
	}
	if cfg.QuarRetention <= 0 {
		cfg.QuarRetention = 7
	}

	s := &Server{
		config: cfg,
	}

	unaryInterceptor := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if err := s.authenticate(ctx); err != nil {
			return nil, err
		}
		return handler(ctx, req)
	}

	streamInterceptor := func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if err := s.authenticate(ss.Context()); err != nil {
			return err
		}
		return handler(srv, ss)
	}

	s.grpcServer = grpc.NewServer(
		grpc.UnaryInterceptor(unaryInterceptor),
		grpc.StreamInterceptor(streamInterceptor),
	)

	clamavv1.RegisterScannerServiceServer(s.grpcServer, s)
	return s
}

// Serve binds and serves the gRPC server on the given listener.
func (s *Server) Serve(lis net.Listener) error {
	return s.grpcServer.Serve(lis)
}

// GracefulStop gracefully shuts down the gRPC server.
func (s *Server) GracefulStop() {
	s.grpcServer.GracefulStop()
}

// Stop abruptly shuts down the gRPC server.
func (s *Server) Stop() {
	s.grpcServer.Stop()
}

func (s *Server) authenticate(ctx context.Context) error {
	if s.config.AuthMode == "" || s.config.AuthMode == "none" {
		return nil
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Error(codes.Unauthenticated, "missing metadata")
	}

	authHeaders := md.Get("authorization")
	if len(authHeaders) == 0 {
		return status.Error(codes.Unauthenticated, "missing authorization header")
	}

	token := authHeaders[0]
	if strings.HasPrefix(strings.ToLower(token), "bearer ") {
		token = strings.TrimSpace(token[7:])
	}

	if s.config.AuthMode == "bearer" {
		if s.config.BearerToken != "" && token != s.config.BearerToken {
			return status.Error(codes.Unauthenticated, "invalid bearer token")
		}
	}

	return nil
}

// HealthCheck returns the server readiness and engine metadata.
func (s *Server) HealthCheck(ctx context.Context, req *clamavv1.HealthCheckRequest) (*clamavv1.HealthCheckResponse, error) {
	return &clamavv1.HealthCheckResponse{
		Status:  clamavv1.HealthCheckResponse_SERVING,
		Version: "1.2.0",
		Engine:  "ClamAV",
	}, nil
}

// ScanStream accepts a client stream of chunks and returns a final scan verdict.
func (s *Server) ScanStream(stream clamavv1.ScannerService_ScanStreamServer) error {
	startTime := time.Now()
	maxBytes := s.config.MaxScanSizeMB * 1024 * 1024

	var meta *clamavv1.FileMetadata
	var buffer bytes.Buffer
	hasher := sha256.New()
	var totalBytes int64

	for {
		chunkReq, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return status.Errorf(codes.Internal, "failed receiving chunk: %v", err)
		}

		if m := chunkReq.GetMetadata(); m != nil {
			meta = m
			continue
		}

		chunkData := chunkReq.GetChunkData()
		if len(chunkData) > 0 {
			totalBytes += int64(len(chunkData))
			if totalBytes > maxBytes {
				return status.Errorf(codes.ResourceExhausted, "payload exceeds configured limit of %d MB", s.config.MaxScanSizeMB)
			}
			buffer.Write(chunkData)
			hasher.Write(chunkData)
		}
	}

	fileName := "stream_file"
	consumerName := "grpc_client"
	expectedSHA256 := ""
	if meta != nil {
		if meta.FileName != "" {
			fileName = meta.FileName
		}
		if meta.ConsumerName != "" {
			consumerName = meta.ConsumerName
		}
		expectedSHA256 = meta.ExpectedSha256
	}

	fileHash := hex.EncodeToString(hasher.Sum(nil))

	// Verify checksum if expected_sha256 was provided
	if expectedSHA256 != "" && !strings.EqualFold(expectedSHA256, fileHash) {
		return status.Errorf(codes.InvalidArgument, "file checksum mismatch: computed %s, expected %s", fileHash, expectedSHA256)
	}

	durationMs := time.Since(startTime).Milliseconds()

	// Check SQLite Whitelist
	if isWhitelisted, _ := s.config.DB.IsWhitelisted(fileHash); isWhitelisted {
		return stream.SendAndClose(&clamavv1.ScanVerdictResponse{
			Success: true,
			Verdict: clamavv1.VerdictStatus_CLEAN,
			Data: &clamavv1.ScanData{
				FileName:       fileName,
				FileSize:       totalBytes,
				FileSha256:     fileHash,
				ScanDurationMs: durationMs,
				ScannedAt:      time.Now().UTC().Format(time.RFC3339),
				Whitelisted:    true,
			},
		})
	}

	// Scan via ClamAV daemon
	scanRes, err := s.config.Clamd.ScanStream(stream.Context(), bytes.NewReader(buffer.Bytes()))
	if err != nil {
		return status.Errorf(codes.Unavailable, "clamav daemon unavailable: %v", err)
	}

	if scanRes.IsClean() {
		_ = s.config.DB.InsertScanAuditLog(storage.ScanAuditLog{
			ID:             fmt.Sprintf("audit_grpc_%d", time.Now().UnixNano()),
			Timestamp:      time.Now().UTC(),
			ConsumerName:   consumerName,
			ClientIP:       "grpc",
			FileName:       fileName,
			FileSizeBytes:  totalBytes,
			FileSHA256:     fileHash,
			Verdict:        "CLEAN",
			ScanDurationMs: durationMs,
		})

		return stream.SendAndClose(&clamavv1.ScanVerdictResponse{
			Success: true,
			Verdict: clamavv1.VerdictStatus_CLEAN,
			Data: &clamavv1.ScanData{
				FileName:       fileName,
				FileSize:       totalBytes,
				FileSha256:     fileHash,
				ScanDurationMs: durationMs,
				ScannedAt:      time.Now().UTC().Format(time.RFC3339),
				Whitelisted:    false,
			},
		})
	}

	// INFECTED: Quarantine & Alert
	quarRec, _ := s.config.Vault.QuarantineFile(stream.Context(), fileName, consumerName, scanRes.VirusName, bytes.NewReader(buffer.Bytes()), s.config.QuarRetention)
	quarID := ""
	if quarRec != nil {
		quarID = quarRec.ID
	}

	_ = s.config.DB.InsertScanAuditLog(storage.ScanAuditLog{
		ID:             fmt.Sprintf("audit_grpc_%d", time.Now().UnixNano()),
		Timestamp:      time.Now().UTC(),
		ConsumerName:   consumerName,
		ClientIP:       "grpc",
		FileName:       fileName,
		FileSizeBytes:  totalBytes,
		FileSHA256:     fileHash,
		Verdict:        "INFECTED",
		VirusName:      scanRes.VirusName,
		ScanDurationMs: durationMs,
		QuarantineID:   quarID,
	})

	if s.config.Notifier != nil {
		_ = s.config.Notifier.DispatchThreat(stream.Context(), alert.ThreatAlert{
			VirusName:      scanRes.VirusName,
			FileName:       fileName,
			FileSizeBytes:  totalBytes,
			FileSHA256:     fileHash,
			QuarantineID:   quarID,
			SourceConsumer: consumerName,
			DetectedAt:     time.Now().UTC(),
		})
	}

	return stream.SendAndClose(&clamavv1.ScanVerdictResponse{
		Success: true,
		Verdict: clamavv1.VerdictStatus_INFECTED,
		Threat: &clamavv1.ThreatInfo{
			VirusName:    scanRes.VirusName,
			Severity:     "HIGH",
			ActionTaken:  "QUARANTINED",
			QuarantineId: quarID,
		},
		Data: &clamavv1.ScanData{
			FileName:       fileName,
			FileSize:       totalBytes,
			FileSha256:     fileHash,
			ScanDurationMs: durationMs,
			ScannedAt:      time.Now().UTC().Format(time.RFC3339),
		},
	})
}
