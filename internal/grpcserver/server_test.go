package grpcserver

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net"
	"path/filepath"
	"testing"

	"github.com/vfat/vqf-clamav-service/internal/alert"
	"github.com/vfat/vqf-clamav-service/internal/clamd"
	"github.com/vfat/vqf-clamav-service/internal/quarantine"
	"github.com/vfat/vqf-clamav-service/internal/ratelimit"
	"github.com/vfat/vqf-clamav-service/internal/storage"
	clamavv1 "github.com/vfat/vqf-clamav-service/proto/clamav/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

func setupTestGRPCServer(t *testing.T, authMode, token string, maxScanMB int64) (*grpc.ClientConn, func()) {
	tmpDir := t.TempDir()
	db, err := storage.NewDB(filepath.Join(tmpDir, "test_grpc.db"))
	if err != nil {
		t.Fatalf("failed to init DB: %v", err)
	}

	vault := quarantine.NewVault(filepath.Join(tmpDir, "vault"), db)
	notifier := alert.NewNotifier(alert.Config{})
	limiter := ratelimit.NewLimiter()
	clamdClient := clamd.NewClient("unix", "/tmp/mock.sock")

	lis := bufconn.Listen(bufSize)

	cfg := Config{
		DB:            db,
		Vault:         vault,
		Notifier:      notifier,
		Limiter:       limiter,
		Clamd:         clamdClient,
		AuthMode:      authMode,
		BearerToken:   token,
		MaxScanSizeMB: maxScanMB,
	}

	grpcSrv := NewServer(cfg)
	go func() {
		if err := grpcSrv.Serve(lis); err != nil && err != grpc.ErrServerStopped {
			// ignore stopped error
		}
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	conn, err := grpc.DialContext(context.Background(), "bufnet",
		grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("failed to dial bufnet: %v", err)
	}

	cleanup := func() {
		conn.Close()
		grpcSrv.Stop()
		lis.Close()
		db.Close()
	}

	return conn, cleanup
}

func TestGRPCServer_HealthCheck(t *testing.T) {
	conn, cleanup := setupTestGRPCServer(t, "none", "", 10)
	defer cleanup()

	client := clamavv1.NewScannerServiceClient(conn)
	res, err := client.HealthCheck(context.Background(), &clamavv1.HealthCheckRequest{Service: "clamav.v1.ScannerService"})
	if err != nil {
		t.Fatalf("HealthCheck failed: %v", err)
	}

	if res.Status != clamavv1.HealthCheckResponse_SERVING {
		t.Errorf("expected status SERVING, got %v", res.Status)
	}
	if res.Engine != "ClamAV" {
		t.Errorf("expected engine ClamAV, got %s", res.Engine)
	}
}

func TestGRPCServer_ScanStream_Clean(t *testing.T) {
	conn, cleanup := setupTestGRPCServer(t, "none", "", 10)
	defer cleanup()

	client := clamavv1.NewScannerServiceClient(conn)
	stream, err := client.ScanStream(context.Background())
	if err != nil {
		t.Fatalf("failed to initiate ScanStream: %v", err)
	}

	cleanData := []byte("Safe content delivered over high-throughput gRPC stream")
	hasher := sha256.New()
	hasher.Write(cleanData)
	expectedHash := hex.EncodeToString(hasher.Sum(nil))

	// Send metadata first
	err = stream.Send(&clamavv1.ScanChunkRequest{
		Payload: &clamavv1.ScanChunkRequest_Metadata{
			Metadata: &clamavv1.FileMetadata{
				FileName:       "report.pdf",
				ConsumerName:   "billing-service",
				ExpectedSha256: expectedHash,
			},
		},
	})
	if err != nil {
		t.Fatalf("failed sending metadata: %v", err)
	}

	// Send binary chunks
	err = stream.Send(&clamavv1.ScanChunkRequest{
		Payload: &clamavv1.ScanChunkRequest_ChunkData{
			ChunkData: cleanData,
		},
	})
	if err != nil {
		t.Fatalf("failed sending chunk data: %v", err)
	}

	res, err := stream.CloseAndRecv()
	// Mock clamd returns Unavailable or clean response
	if err != nil {
		st, ok := status.FromError(err)
		if !ok || st.Code() != codes.Unavailable {
			t.Fatalf("unexpected error from stream: %v", err)
		}
		return
	}

	if !res.Success {
		t.Errorf("expected success true, got false")
	}
	if res.Verdict != clamavv1.VerdictStatus_CLEAN {
		t.Errorf("expected CLEAN verdict, got %v", res.Verdict)
	}
	if res.Data.FileSha256 != expectedHash {
		t.Errorf("expected hash %s, got %s", expectedHash, res.Data.FileSha256)
	}
}

func TestGRPCServer_ScanStream_ChecksumMismatch(t *testing.T) {
	conn, cleanup := setupTestGRPCServer(t, "none", "", 10)
	defer cleanup()

	client := clamavv1.NewScannerServiceClient(conn)
	stream, err := client.ScanStream(context.Background())
	if err != nil {
		t.Fatalf("failed to initiate ScanStream: %v", err)
	}

	// Send metadata with wrong expected checksum
	err = stream.Send(&clamavv1.ScanChunkRequest{
		Payload: &clamavv1.ScanChunkRequest_Metadata{
			Metadata: &clamavv1.FileMetadata{
				FileName:       "document.pdf",
				ConsumerName:   "gateway",
				ExpectedSha256: "0000000000000000000000000000000000000000000000000000000000000000",
			},
		},
	})
	if err != nil {
		t.Fatalf("failed sending metadata: %v", err)
	}

	err = stream.Send(&clamavv1.ScanChunkRequest{
		Payload: &clamavv1.ScanChunkRequest_ChunkData{
			ChunkData: []byte("Some actual payload content"),
		},
	})
	if err != nil {
		t.Fatalf("failed sending chunk data: %v", err)
	}

	_, err = stream.CloseAndRecv()
	if err == nil {
		t.Fatal("expected error due to checksum mismatch, got nil")
	}

	st, ok := status.FromError(err)
	if !ok || st.Code() != codes.InvalidArgument {
		t.Errorf("expected InvalidArgument code, got: %v", err)
	}
}

func TestGRPCServer_ScanStream_PayloadTooLarge(t *testing.T) {
	// Set MaxScanSizeMB to 1 MB
	conn, cleanup := setupTestGRPCServer(t, "none", "", 1)
	defer cleanup()

	client := clamavv1.NewScannerServiceClient(conn)
	stream, err := client.ScanStream(context.Background())
	if err != nil {
		t.Fatalf("failed to initiate ScanStream: %v", err)
	}

	largeChunk := make([]byte, 512*1024) // 512 KB
	// Send metadata
	_ = stream.Send(&clamavv1.ScanChunkRequest{
		Payload: &clamavv1.ScanChunkRequest_Metadata{
			Metadata: &clamavv1.FileMetadata{
				FileName: "large.iso",
			},
		},
	})

	// Send 3 chunks of 512 KB = 1.5 MB (> 1 MB limit)
	for i := 0; i < 3; i++ {
		sendErr := stream.Send(&clamavv1.ScanChunkRequest{
			Payload: &clamavv1.ScanChunkRequest_ChunkData{
				ChunkData: largeChunk,
			},
		})
		if sendErr != nil && sendErr != io.EOF {
			break
		}
	}

	_, err = stream.CloseAndRecv()
	if err == nil {
		t.Fatal("expected ResourceExhausted error, got nil")
	}

	st, ok := status.FromError(err)
	if !ok || st.Code() != codes.ResourceExhausted {
		t.Errorf("expected ResourceExhausted, got %v", err)
	}
}

func TestGRPCServer_AuthInterceptor(t *testing.T) {
	conn, cleanup := setupTestGRPCServer(t, "bearer", "my-secret-grpc-token", 10)
	defer cleanup()

	client := clamavv1.NewScannerServiceClient(conn)

	// 1. Without metadata -> unauthenticated
	_, err := client.HealthCheck(context.Background(), &clamavv1.HealthCheckRequest{})
	if err == nil {
		t.Fatal("expected Unauthenticated error, got nil")
	}
	st, _ := status.FromError(err)
	if st.Code() != codes.Unauthenticated {
		t.Errorf("expected Unauthenticated, got %v", st.Code())
	}

	// 2. With valid Bearer metadata -> success
	ctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs("authorization", "Bearer my-secret-grpc-token"))
	res, err := client.HealthCheck(ctx, &clamavv1.HealthCheckRequest{})
	if err != nil {
		t.Fatalf("HealthCheck with token failed: %v", err)
	}
	if res.Status != clamavv1.HealthCheckResponse_SERVING {
		t.Errorf("expected SERVING, got %v", res.Status)
	}
}
