# DCD-02-gRPC-Scanner

## 1. Metadata Komponen

| Field | Keterangan |
|---|---|
| **Komponen** | gRPC Scanner Service (High-Throughput Binary Streaming Ingress) |
| **Kode Dokumen** | `DCD-02` |
| **Dasar Rancangan** | [`SCD-gRPC-Scanner.md`](file:///home/ubuntu/workspace/plan/clamav-service/.ai-doc/plan/component/SCD-gRPC-Scanner.md) |
| **Status Implementasi** | ✅ `Active / Production Ready` (`TDD-012` 100% GREEN) |
| **Target Audience** | Microservice Developers, Backend Engineers, DevOps |

---

## 2. Object Identification

Berdasarkan arsitektur implementasi pada [`internal/grpcserver/`](file:///home/ubuntu/workspace/plan/clamav-service/internal/grpcserver/) dan [`proto/clamav/v1/`](file:///home/ubuntu/workspace/plan/clamav-service/proto/clamav/v1/):

### Boundary
* **TCP Socket Listener:** `net.Listener` membuka socket TCP pada port `9090` (`GRPC_PORT`) untuk koneksi HTTP/2 gRPC.
* **gRPC Server Engine:** `grpc.Server` membungkus multiplexer streaming RPC dan mengikat *unary & stream interceptors*.
* **Unix Domain Socket Stream Interface:** `net.Conn` Unix socket (`/var/run/clamav/clamd.ctl`) mengalirkan data chunked langsung ke engine ClamAV.

### Control
* **`grpcserver.Server.ScanStream`:** Controller client-streaming RPC yang membaca pesan biner bertahap, memvalidasi kuota memori, menghitung SHA-256 secara live, dan menyalurkan payload ke socket ClamAV.
* **`grpcserver.Server.authenticate`:** Interceptor otentikasi metadata gRPC yang memverifikasi token Bearer pada setiap panggilan RPC.
* **`grpcserver.Server.HealthCheck`:** Controller pelaporan kesehatan server yang mengembalikan status kesiapan engine (`SERVING`).
* **`clamd.Client.ScanStream`:** Controller pengirim chunk biner (`zINSTREAM`) ke daemon antivirus ClamAV.
* **`quarantine.Vault.QuarantineFile`:** Controller enkripsi AES-256-GCM untuk isolasi file bervirus.
* **`alert.Notifier.DispatchThreat`:** Dispatcher pengiriman alert ke Telegram/Discord/Webhook.

### Entity
* **Service RPC Endpoints:**
  * `clamav.v1.ScannerService/ScanStream`
  * `clamav.v1.ScannerService/HealthCheck`
* **Protobuf Message Entities:**
  * `clamavv1.ScanChunkRequest`: Envelop pembungkus oneof payload (`metadata` atau `chunk_data`).
  * `clamavv1.FileMetadata`: Metadata file (`file_name`, `consumer_name`, `expected_sha256`).
  * `clamavv1.ScanVerdictResponse`: Kontrak hasil pemindaian (`success`, `verdict`, `threat`, `data`).
  * `clamavv1.ThreatInfo`: Informasi ancaman (`virus_name`, `severity`, `action_taken`, `quarantine_id`).
  * `clamavv1.ScanData`: Detail statistik pemindaian (`file_name`, `file_size`, `file_sha256`, `scan_duration_ms`).
  * `clamavv1.HealthCheckResponse`: Status kesiapan daemon (`SERVING`, `version`, `engine`).
* **Domain Model `storage.ScanAuditLog`:** Catatan audit log di SQLite dengan client IP `"grpc"`.
* **Configuration `grpcserver.Config`:** Konfigurasi server (`MaxScanSizeMB: 100`, `QuarRetention: 7`, `AuthMode: "bearer"`, `BearerToken`).

---

## 3. Use Case List

| No | Kode Use Case | Nama Use Case | Actor | Status | Referensi Detail |
|---|---|---|---|---|---|
| 1 | `UC-GRPC-01` | Client-Streaming Scan File Bersih | Inter-Service Consumer | Active | Section 4.1 |
| 2 | `UC-GRPC-02` | Deteksi Virus & Karantina Otomatis | System (Autonomous) | Active | Section 4.2 |
| 3 | `UC-GRPC-03` | Verifikasi SHA-256 & Proteksi Payload Limit | Upstream Client | Active | Section 4.3 |
| 4 | `UC-GRPC-04` | gRPC Health Probe & Auth Interceptor | Kubernetes / API Gateway | Active | Section 4.4 |

---

## 4. Use Case Detail

### 4.1. UC-GRPC-01: Client-Streaming Scan File Bersih
* **Deskripsi:** Menerima aliran potongan biner (*chunked bytes*) dari klien mikroservis melalui client-streaming RPC, menyalurkannya langsung ke ClamAV Unix socket, dan mengembalikan status `CLEAN`.
* **Actor:** Inter-Service Consumer (e.g., File Upload Service, Document Management Service).
* **Precondition:** Koneksi gRPC HTTP/2 terbentuk, ClamAV socket siap menerima perintah `zINSTREAM`.
* **Postcondition:** File dinyatakan `CLEAN`, audit log tercatat di SQLite, respons `ScanVerdictResponse` dikembalikan.
* **Normal Flow:**
  1. Client memanggil RPC `ScanStream`.
  2. Client mengirim chunk pertama berupa `ScanChunkRequest_Metadata` berisi `file_name = "annual_report.pdf"` dan `consumer_name = "document-svc"`.
  3. Client mengirim rentetan `ScanChunkRequest_ChunkData` (misal per 64 KB) secara berulang.
  4. Server membaca chunk via `stream.Recv()`, mengumpulkan bytes ke buffer in-memory, dan memperbarui perhitungan SHA-256 secara live.
  5. Setelah selesai, client menutup arah pengiriman via `stream.CloseAndRecv()`.
  6. Server memeriksa apakah hash SHA-256 telah terdaftar pada database whitelist SQLite. Jika tidak:
  7. Server mengalirkan payload ke ClamAV daemon via `clamd.ScanStream`.
  8. ClamAV mengonfirmasi stream bebas virus (status clean).
  9. Server mencatat histori audit log ke tabel `scan_audit_logs`.
  10. Server mengirimkan respons gRPC final `ScanVerdictResponse` dengan `verdict = VerdictStatus_CLEAN` dan `success = true`.
* **Traceability:**
  * Test: `TestGRPCServer_ScanStream_Clean` di [`internal/grpcserver/server_test.go`](file:///home/ubuntu/workspace/plan/clamav-service/internal/grpcserver/server_test.go).
  * Implementation: `ScanStream` di [`internal/grpcserver/server.go`](file:///home/ubuntu/workspace/plan/clamav-service/internal/grpcserver/server.go).

---

### 4.2. UC-GRPC-02: Deteksi Virus & Karantina Otomatis
* **Deskripsi:** Memproses aliran stream biner yang mengandung malware, mendeteksi jenis virus, otomatis mengenkripsi dan mengisolasi ke Quarantine Vault (AES-256-GCM), serta mengembalikan status `INFECTED`.
* **Actor:** System (Autonomous Security Engine).
* **Precondition:** Stream file berbahaya selesai diterima oleh server gRPC.
* **Postcondition:** Payload diisolasi di disk karantina `/data/quarantine` dengan header `VQF_AESGCM_V1\n`, log audit dicatat sebagai `INFECTED`, ancaman dinotifikasikan, respons `INFECTED` dikembalikan dengan `quarantine_id`.
* **Normal Flow:**
  1. Aliran stream biner selesai diterima oleh server.
  2. Server mengirimkan bytes ke socket ClamAV (`zINSTREAM`).
  3. ClamAV merespons dengan deteksi virus: `status = FOUND`, `virus_name = "Win.Trojan.Generic-10294"`.
  4. Server memanggil `Vault.QuarantineFile`:
     * Membangkitkan CSPRNG nonce 12-byte.
     * Mengenkripsi biner utuh menggunakan cipher AES-256-GCM.
     * Menulis berkas terenkripsi ke `/data/quarantine/Q-20260907-xxxx`.
     * Mencatat record ke tabel SQLite `quarantine_records`.
  5. Server mencatat log pemindaian ke tabel SQLite `scan_audit_logs` dengan status `INFECTED` dan foreign key `quarantine_id`.
  6. `alert.Notifier` mendistribusikan pesan peringatan bahaya ke channel webhook/Telegram.
  7. Server mengembalikan respons gRPC `ScanVerdictResponse` memuat:
     * `verdict = VerdictStatus_INFECTED`
     * `threat.virus_name = "Win.Trojan.Generic-10294"`
     * `threat.action_taken = "QUARANTINED"`
     * `threat.quarantine_id = "Q-20260907-xxxx"`
* **Traceability:**
  * Test: `TestVault_AESGCM_EncryptionAndRestore` di [`internal/quarantine/vault_test.go`](file:///home/ubuntu/workspace/plan/clamav-service/internal/quarantine/vault_test.go).
  * Implementation: `ScanStream` di [`internal/grpcserver/server.go`](file:///home/ubuntu/workspace/plan/clamav-service/internal/grpcserver/server.go).

---

### 4.3. UC-GRPC-03: Verifikasi SHA-256 & Proteksi Payload Limit
* **Deskripsi:** Menjaga integritas data dan stabilitas memori gRPC server dengan membatasi ukuran stream maksimal (`MAX_SCAN_SIZE_MB`) serta memvalidasi kesesuaian hash yang diharapkan.
* **Actor:** Upstream Client / Ingress Gateway.
* **Precondition:** Klien memulai streaming biner ke `ScanStream`.
* **Postcondition:** Stream yang melanggar batas atau tidak cocok hash-nya ditolak dengan gRPC error status code spesifik.
* **Normal Flow (Payload Limit Guard):**
  1. Klien mengalirkan bongkahan data besar secara terus menerus.
  2. Server menghitung akumulasi total bytes yang diterima (`totalBytes += len(chunkData)`).
  3. Ketika `totalBytes > maxBytes` (misal melebihi batas 100 MB), server seketika menghentikan pembacaan stream.
  4. Server mengembalikan gRPC error dengan status `codes.ResourceExhausted`:
     > *"payload exceeds configured limit of 100 MB"*
* **Alternative Flow (Checksum Mismatch Guard):**
  1. Klien menyertakan `expected_sha256` pada metadata awal.
  2. Seluruh chunk selesai diterima dan di-hash oleh server.
  3. Server membandingkan hash kalkulasi dengan hash expected klien.
  4. Jika hash tidak cocok, server menolak melanjutkan ke ClamAV engine dan mengembalikan gRPC error status `codes.InvalidArgument`:
     > *"file checksum mismatch: computed <hash>, expected <hash>"*
* **Traceability:**
  * Test: `TestGRPCServer_ScanStream_PayloadTooLarge`, `TestGRPCServer_ScanStream_ChecksumMismatch` di [`internal/grpcserver/server_test.go`](file:///home/ubuntu/workspace/plan/clamav-service/internal/grpcserver/server_test.go).

---

### 4.4. UC-GRPC-04: gRPC Health Probe & Auth Interceptor
* **Deskripsi:** Memeriksa kesiapan dan kesehatan layanan antivirus melalui RPC standar `HealthCheck`, serta memvalidasi hak akses klien via interceptor metadata gRPC.
* **Actor:** Kubernetes Liveness/Readiness Probe / API Gateway.
* **Precondition:** Port gRPC aktif dan siap menerima panggilan RPC.
* **Postcondition:** Panggilan tanpa kredensial valid ditolak dengan `codes.Unauthenticated`, panggilan valid diproses.
* **Normal Flow (Health Check):**
  1. Kubernetes probe memanggil RPC `HealthCheck` dengan `HealthCheckRequest{Service: "clamav.v1.ScannerService"}`.
  2. Interceptor memverifikasi header auth (jika `AUTH_MODE=bearer`).
  3. Handler mengembalikan `HealthCheckResponse` dengan `status = SERVING`, `version = "1.2.0"`, dan `engine = "ClamAV"`.
* **Alternative Flow (Auth Enforcement):**
  1. Jika `AUTH_MODE=bearer` aktif dan klien tidak mengirimkan metadata `authorization` atau token salah:
  2. Interceptor mencegat eksekusi handler dan langsung me-return gRPC error `codes.Unauthenticated`:
     > *"invalid or missing authentication token"*
* **Traceability:**
  * Test: `TestGRPCServer_HealthCheck`, `TestGRPCServer_AuthInterceptor` di [`internal/grpcserver/server_test.go`](file:///home/ubuntu/workspace/plan/clamav-service/internal/grpcserver/server_test.go).
