# SCD-gRPC-Scanner (Single Component Design)

## 1. Context

### 1.1. Masalah yang Ingin Diselesaikan
Saat ini `clamav-service` mengekspos pemindaian malware melalui HTTP/1.1 REST API:
1. `multipart/form-data` upload (`POST /api/v1/scan/file`)
2. Raw binary HTTP body stream (`POST /api/v1/scan/stream`)
3. Remote URL stream fetcher (`POST /api/v1/scan/url`)

Pada infrastruktur microservice skala enterprise dan Kubernetes mesh, komunikasi antar-layanan (inter-service communication) memerlukan throughput tinggi dengan latency minimal. Penggunaan HTTP/1.1 REST memiliki overhead berupa text header framing, base64/multipart parsing, dan pembentukan koneksi TCP berulang.

Dibutuhkan antarmuka **gRPC (HTTP/2 binary framing)** dengan mekanisme **client-streaming RPC** sehingga upstream service dapat mengalirkan bongkahan biner (*chunked bytes*) langsung ke pipeline pemindaian antivirus tanpa overhead parsing multipart dan mendukung multiplexing dalam satu koneksi TCP persisten.

### 1.2. Posisi Komponen dalam Sistem
Komponen `gRPC-Scanner` bertindak sebagai *High-Throughput Binary Ingress Engine* yang berdampingan dengan HTTP REST API server:

```
┌─────────────────────────────────────────────────────────────┐
│  clamav-service                                             │
│                                                             │
│  ┌───────────────────────────┐ ┌──────────────────────────┐ │
│  │ HTTP REST API (Port 8080) │ │ gRPC Server (Port 9090)  │ │
│  └─────────────┬─────────────┘ └────────────┬─────────────┘ │
│                │                            │               │
│                │ Auth / Rate Limit          │ Interceptors  │
│                │                            │ (Auth, Limit) │
│                ▼                            ▼               │
│  ┌────────────────────────────────────────────────────────┐ │
│  │ ClamAV Scanner Core / Pipeline Interface               │ │
│  │  - Chunked Streamer (io.Pipe / io.LimitReader)          │ │
│  │  - SHA-256 Hasher & Whitelist Checker                  │ │
│  │  - ClamAV Unix Socket Client (/var/run/clamd)          │ │
│  │  - Quarantine Vault (AES-256-GCM)                      │ │
│  │  - SQLite Audit Logger & Alert Dispatcher              │ │
│  └────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

### 1.3. Hubungan dengan Komponen Lain
* **`internal/clamd`**: Menerima chunked biner langsung dari gRPC stream via Unix Domain Socket (`zINSTREAM`).
* **`internal/quarantine`**: Mengisolasi payload file terinfeksi ke Vault (AES-256-GCM terenkripsi).
* **`internal/storage`**: Menyimpan riwayat audit log pemindaian gRPC dan mencocokkan hash whitelist.
* **`internal/alert`**: Mengirimkan notifikasi insiden jika file terindikasi virus.
* **`internal/ratelimit`**: Membatasi laju request per consumer via gRPC Unary/Stream Interceptor.

---

## 2. Scope

### 2.1. In-Scope
1. **Definisi Protocol Buffers (`proto3`)**:
   * Service `clamav.v1.ScannerService`.
   * RPC `ScanStream(stream ScanChunkRequest) returns (ScanVerdictResponse)`: Client-streaming chunked biner.
   * RPC `HealthCheck(HealthCheckRequest) returns (HealthCheckResponse)`: Standar liveness/readiness probe untuk gRPC ecosystem.
2. **Pola Chunked Client Streaming**:
   * Chunk pertama (`metadata`): Memuat metadata file (`file_name`, `consumer_name`, `expected_sha256`).
   * Chunk berikutnya (`chunk_data`): Potongan biner (misal per 32 KB atau 64 KB).
3. **Batas Ukuran & Proteksi Alokasi Memori**:
   * Penegakan batas alokasi memori selaras dengan `MAX_SCAN_SIZE_MB` via `io.LimitReader`.
   * Menolak stream jika total ukuran melebihi batas dengan gRPC status code `RESOURCE_EXHAUSTED`.
4. **Verifikasi Integritas & Whitelist**:
   * Perhitungan SHA-256 secara on-the-fly saat stream dibaca.
   * Validasi terhadap `expected_sha256` (menghasilkan gRPC status `INVALID_ARGUMENT` jika tidak cocok).
   * Pengecekan hash whitelist lokal SQLite (verdict `CLEAN` instan).
5. **Autentikasi & Otorisasi Interceptor**:
   * Pemeriksaan metadata gRPC (`authorization` Bearer token atau `x-api-key`).
   * Mengembalikan gRPC status `UNAUTHENTICATED` jika token tidak valid dan `AUTH_MODE != none`.
6. **Graceful Shutdown**:
   * gRPC server mendukung `GracefulStop()` terkoordinasi dengan sinyal shutdown supervisor OS (`SIGINT`, `SIGTERM`).

### 2.2. Out-of-Scope
* Bidirectional live stream feedback per byte (karena ClamAV mengevaluasi verdict setelah transfer stream selesai / zero-byte chunk).
* Web UI gRPC-Web proxy (Web UI tetap berkomunikasi melalui HTTP REST API Port 8080).

---

## 3. Prerequisite

1. Dependency Go `google.golang.org/grpc` dan `google.golang.org/protobuf`.
2. Generator tool `protoc` dan plugin `protoc-gen-go`, `protoc-gen-go-grpc` atau pre-generated code yang dikompilasi secara deterministik.
3. Subkomponen `internal/clamd`, `internal/quarantine`, `internal/storage`, dan `internal/crypto` stabil dan berstatus 100% GREEN.

---

## 4. Daftar Usecase

| Kode Usecase | Nama Usecase | Deskripsi Singkat |
|---|---|---|
| **`UC-GRPC-01`** | Client-Streaming Scan File Bersih | Mengalirkan biner file bersih via gRPC stream chunked, memindai ke ClamAV daemon, dan mengembalikan `ScanVerdictResponse` dengan status `CLEAN`. |
| **`UC-GRPC-02`** | Deteksi Virus & Karantina Otomatis | Memindai file terinfeksi via gRPC stream, mendeteksi virus, otomatis mengisolasi ke Vault AES-256-GCM, dan mengembalikan status `INFECTED` beserta `quarantine_id`. |
| **`UC-GRPC-03`** | Verifikasi SHA-256 & Proteksi Payload Limit | Menolak stream jika `expected_sha256` tidak cocok (`INVALID_ARGUMENT`) atau total bytes melebihi `MAX_SCAN_SIZE_MB` (`RESOURCE_EXHAUSTED`). |
| **`UC-GRPC-04`** | gRPC Health Probe & Auth Interceptor | Menyediakan health check endpoint `HealthCheck` serta penegakan autentikasi token metadata via gRPC Interceptor. |

---

## 5. Spesifikasi Protocol Buffers (`scanner.proto`)

```protobuf
syntax = "proto3";

package clamav.v1;

option go_package = "github.com/vfat/vqf-clamav-service/proto/clamav/v1;clamavv1";

service ScannerService {
  // Client-streaming RPC untuk memindai biner file
  rpc ScanStream(stream ScanChunkRequest) returns (ScanVerdictResponse);

  // Health check RPC untuk observability & k8s probes
  rpc HealthCheck(HealthCheckRequest) returns (HealthCheckResponse);
}

message FileMetadata {
  string file_name = 1;
  string consumer_name = 2;
  string expected_sha256 = 3;
}

message ScanChunkRequest {
  oneof payload {
    FileMetadata metadata = 1;
    bytes chunk_data = 2;
  }
}

enum VerdictStatus {
  VERDICT_UNSPECIFIED = 0;
  CLEAN = 1;
  INFECTED = 2;
}

message ThreatInfo {
  string virus_name = 1;
  string severity = 2;
  string action_taken = 3;
  string quarantine_id = 4;
}

message ScanData {
  string file_name = 1;
  int64 file_size = 2;
  string file_sha256 = 3;
  int64 scan_duration_ms = 4;
  string scanned_at = 5;
  bool whitelisted = 6;
}

message ScanVerdictResponse {
  bool success = 1;
  VerdictStatus verdict = 2;
  ThreatInfo threat = 3;
  ScanData data = 4;
}

message HealthCheckRequest {
  string service = 1;
}

message HealthCheckResponse {
  enum ServingStatus {
    UNKNOWN = 0;
    SERVING = 1;
    NOT_SERVING = 2;
  }
  ServingStatus status = 1;
  string version = 2;
  string engine = 3;
}
```

---

## 6. Asumsi, Risiko, dan Mitigasi

| Kategori | Asumsi / Risiko | Mitigasi Desain |
|---|---|---|
| **Port Collision** | Port gRPC bentrok dengan port REST API (8080). | gRPC server berjalan pada port terpisah yang dapat dikonfigurasi via env `GRPC_PORT` (default: `9090`). |
| **Stream Chunk Flooding** | Client mengirim ribuan chunk kosong atau chunk terlalu besar tanpa henti. | Interceptor memeriksa ukuran maksimal per message (`MaxReceiveMessageSize: 4MB`) dan total stream dibatasi selaras dengan `MAX_SCAN_SIZE_MB`. |
| **Auth Metadata Missing** | Service upstream tidak mengirimkan token di metadata header gRPC. | Interceptor mengevaluasi header metadata `authorization` atau `x-api-key`, mengembalikan error `codes.Unauthenticated` jika gagal. |
| **ClamAV Socket Timeout** | ClamAV lambat merespons stream gRPC besar. | Context timeout 30 detik diterapkan pada level call handler; koneksi ditutup rapi jika waktu habis. |

---

## 7. Rencana Siklus TDD

1. **Target `TDD-012`**:
   - Pembuatan proto definition dan struktur kode package `internal/grpcserver`.
   - Unit test suite gRPC client & in-memory buffer (`bufconn`):
     - Test scan file bersih (`CLEAN`).
     - Test scan file terinfeksi (`INFECTED`) dengan auto-quarantine.
     - Test proteksi ukuran file melebihi limit (`RESOURCE_EXHAUSTED`).
     - Test validasi checksum tidak cocok (`INVALID_ARGUMENT`).
     - Test autentikasi metadata interceptor (`UNAUTHENTICATED`).
     - Test gRPC health check.
2. **Wiring di `cmd/server/main.go`**:
   - Menjalankan gRPC server secara konkuren dengan HTTP REST API server.
   - Integrasi graceful shutdown untuk kedua server secara bersamaan.
