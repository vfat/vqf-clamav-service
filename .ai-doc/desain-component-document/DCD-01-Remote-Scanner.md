# DCD-01-Remote-Scanner

## 1. Metadata Komponen

| Field | Keterangan |
|---|---|
| **Komponen** | Remote Scanner (Anti-SSRF & HTTP Streaming Ingress) |
| **Kode Dokumen** | `DCD-01` |
| **Dasar Rancangan** | [`SCD-Remote-Scanner.md`](file:///home/ubuntu/workspace/plan/clamav-service/.ai-doc/plan/component/SCD-Remote-Scanner.md) |
| **Status Implementasi** | ✅ `Active / Production Ready` (`TDD-010` & `TDD-011` 100% GREEN) |
| **Target Audience** | Backend Engineers, Security Architects, DevOps |

---

## 2. Object Identification

Berdasarkan arsitektur implementasi pada [`internal/fetcher/`](file:///home/ubuntu/workspace/plan/clamav-service/internal/fetcher/) dan [`internal/api/`](file:///home/ubuntu/workspace/plan/clamav-service/internal/api/):

### Boundary
* **HTTP Ingress Router:** `http.ServeMux` pada [`internal/api/server.go`](file:///home/ubuntu/workspace/plan/clamav-service/internal/api/server.go) menangani route `POST /api/v1/scan/url`.
* **Outbound HTTP Network Socket:** `net.Dialer` dengan kustom callback `DialContext` untuk koneksi TCP ke host remote.
* **Unix Domain Socket Stream Interface:** `net.Conn` Unix socket (`/var/run/clamav/clamd.ctl`) mengalirkan data chunked ke ClamAV daemon.

### Control
* **`api.handleScanURL`:** HTTP handler utama yang mengoordinasikan validasi payload, eksekusi download remote, scanning antivirus, pencatatan audit log, dan isolasi karantina.
* **`fetcher.SafeFetcher`:** Klien HTTP terisolasi yang mengamankan proses download file dari remote URL.
* **`fetcher.IsBlockedIP`:** Logika filter pre-dial IP yang memvalidasi hasil resolusi DNS terhadap blacklist CIDR (Loopback, RFC1918, Cloud Metadata `169.254.169.254`, CGNAT `100.64.0.0/10`).
* **`fetcher.boundedReadCloser`:** Controller pembatas streaming memori yang mengontrol bytes terbaca agar tidak melebihi alokasi `MAX_SCAN_SIZE_MB`.
* **`clamd.Client.ScanStream`:** Controller pengirim chunk biner (`zINSTREAM`) ke daemon antivirus ClamAV.
* **`quarantine.Vault.QuarantineFile`:** Controller enkripsi AES-256-GCM untuk isolasi file yang terdeteksi bervirus.
* **`alert.Notifier.DispatchThreat`:** Dispatcher pengiriman alert ke Telegram/Discord/Webhook.

### Entity
* **Endpoint:** `POST /api/v1/scan/url` (REST API data access & execution contract).
* **DTO `scanURLRequest`:** Struct request JSON `{ "url": string, "expected_sha256": string }`.
* **DTO `FetchResult`:** Struct pembungkus hasil unduhan `{ Body: io.ReadCloser, StatusCode: int, ContentLength: int64, ContentType: string, ResolvedIP: string, FinalURL: string }`.
* **Domain Model `storage.ScanAuditLog`:** Record audit log di SQLite memuat histori pemindaian, SHA-256, ukuran file, durasi, dan verdict.
* **Domain Model `storage.QuarantineRecord`:** Metadata catatan file terisolasi di database SQLite.
* **Configuration `SafeFetcherConfig`:** Konfigurasi batas waktu (`Timeout: 30s`), batasan redirect (`MaxRedirects: 3`), dan kuota ukuran (`MaxBytes`).

---

## 3. Use Case List

| No | Kode Use Case | Nama Use Case | Actor | Status | Referensi Detail |
|---|---|---|---|---|---|
| 1 | `UC-URL-01` | Scan File Remote URL Valid (Clean File) | Consumer Application | Active | Section 4.1 |
| 2 | `UC-URL-02` | Blokir Percobaan Serangan SSRF Internal | Consumer / Attacker | Active | Section 4.2 |
| 3 | `UC-URL-03` | Batasi Ukuran File Remote & Cegah OOM | Consumer / Upstream | Active | Section 4.3 |
| 4 | `UC-URL-04` | Karantina Otomatis File Remote Terinfeksi | System (Autonomous) | Active | Section 4.4 |

---

## 4. Use Case Detail

### 4.1. UC-URL-01: Scan File Remote URL Valid (Clean File)
* **Deskripsi:** Mengunduh file dari remote HTTP/HTTPS URL publik atau S3 presigned URL, memindai aliran biner ke ClamAV daemon, mencatat audit log, dan mengembalikan verdict `CLEAN`.
* **Actor:** Consumer Application.
* **Precondition:** Layanan `clamav-service` aktif, daemon ClamAV siap menerima socket connection, target URL dapat diakses publik.
* **Postcondition:** Audit log tersimpan di SQLite dengan verdict `CLEAN`, respons HTTP 200 dikembalikan ke client.
* **Normal Flow:**
  1. Client mengirimkan `POST /api/v1/scan/url` dengan JSON `{ "url": "https://example.com/invoice.pdf" }`.
  2. Handler memvalidasi sintaks URL dan memastikan skema `http://` atau `https://`.
  3. `SafeFetcher` melakukan resolusi DNS ke target hostname.
  4. `IsBlockedIP` memverifikasi alamat IP publik (bukan private/loopback/metadata).
  5. `SafeFetcher` membuka koneksi langsung (*direct dial*) ke IP terverifikasi guna mencegah DNS rebinding.
  6. Aliran body response dibungkus oleh `boundedReadCloser` dan dialirkan ke memory buffer & SHA-256 hasher.
  7. Handler mengecek database whitelist lokal. Jika tidak whitelisted, bytes dialirkan via `clamd.ScanStream`.
  8. ClamAV merespons dengan status bersih (Clean).
  9. Handler mencatat audit log ke database SQLite (`InsertScanAuditLog`).
  10. Handler mengembalikan respons JSON HTTP 200 dengan metadata file dan verdict `CLEAN`.
* **Alternative Flow:**
  * **Whitelisted Hash:** Jika SHA-256 file telah terdaftar di database whitelist, pemindaian socket ClamAV di-bypass dan respons `CLEAN` langsung dikembalikan dengan flag `"whitelisted": true`.
* **Exception Flow:**
  * **Checksum Mismatch:** Jika client menyertakan `expected_sha256` dan tidak cocok dengan hash aktual file yang diunduh, sistem mengembalikan HTTP 400 (`CHECKSUM_MISMATCH`).
* **Traceability:**
  * Test: `TestHandler_ScanURL_CleanAndChecksum` di [`internal/api/handler_test.go`](file:///home/ubuntu/workspace/plan/clamav-service/internal/api/handler_test.go).
  * Implementation: `handleScanURL` di [`internal/api/server.go`](file:///home/ubuntu/workspace/plan/clamav-service/internal/api/server.go).

---

### 4.2. UC-URL-02: Blokir Percobaan Serangan SSRF Internal
* **Deskripsi:** Mencegah eksploitasi Server-Side Request Forgery (SSRF) dengan menolak target URL yang mengarah ke localhost, IP intranet (RFC1918), cloud metadata instance, atau CGNAT.
* **Actor:** Consumer Application / Security Auditor / Attacker.
* **Precondition:** Request HTTP `POST /api/v1/scan/url` diterima dengan target URL yang mencurigakan.
* **Postcondition:** Tidak ada koneksi soket TCP yang dibuat ke IP internal, percobaan serangan digagalkan, audit/error dikembalikan.
* **Normal Flow:**
  1. Client mengirim request dengan target URL internal, misalnya `http://169.254.169.254/latest/meta-data` atau `http://127.0.0.1:8080`.
  2. `fetcher.ValidateURL` memvalidasi URL.
  3. Pre-dial resolver mengevaluasi hostname target.
  4. `fetcher.IsBlockedIP` mendeteksi bahwa IP masuk dalam salah satu subnet terlarang:
     * Loopback: `127.0.0.0/8`, `::1`
     * Cloud Metadata & Link-Local: `169.254.0.0/16`
     * Private RFC1918: `10.0.0.0/8`, `172.16.0.0/12`, `192.168.0.0/16`
     * Carrier-Grade NAT: `100.64.0.0/10`
     * Wildcard: `0.0.0.0/8`, `::`
     * IPv4-mapped IPv6: `::ffff:x.x.x.x`
  5. `SafeFetcher` membatalkan dial sebelum paket SYN dikirimkan dan me-return `ErrSSRFBlocked`.
  6. Handler menangkap error `ErrSSRFBlocked` dan merespons dengan HTTP 400 Bad Request:
     ```json
     {
       "success": false,
       "error": {
         "code": "SSRF_ATTEMPT_BLOCKED",
         "message": "The provided target URL resolves to a prohibited private, loopback, or cloud-metadata IP address",
         "details": {
           "target_url": "http://169.254.169.254/latest/meta-data",
           "policy": "RFC1918, Loopback, and Link-Local IPs are strictly forbidden"
         }
       }
     }
     ```
* **Alternative Flow:**
  * **SSRF via HTTP Redirect (Hop Validation):** Jika remote server publik me-redirect client ke IP private, fungsi `CheckRedirect` pada klien HTTP mendeteksi IP redirect dan menggagalkannya dengan `ErrSSRFBlocked`.
* **Traceability:**
  * Test: `TestSSRFBlockingDirect`, `TestRedirectAntiSSRF` di [`internal/fetcher/client_test.go`](file:///home/ubuntu/workspace/plan/clamav-service/internal/fetcher/client_test.go) dan `TestHandler_ScanURL_SSRFBlocked` di [`internal/api/handler_test.go`](file:///home/ubuntu/workspace/plan/clamav-service/internal/api/handler_test.go).

---

### 4.3. UC-URL-03: Batasi Ukuran File Remote & Cegah OOM
* **Deskripsi:** Membatasi ukuran unduhan file remote agar tidak menguras kuota memori container dan melindungi daemon dari serangan payload raksasa / exhaustion.
* **Actor:** Consumer Application / Upstream Source.
* **Precondition:** Request download file remote berukuran besar diterima.
* **Postcondition:** Aliran stream dipotong seketika saat mencapai batas `MAX_SCAN_SIZE_MB`, koneksi TCP ditutup, error 413 dikembalikan.
* **Normal Flow:**
  1. Client mengirimkan target URL file yang berukuran melebihi batas (misal file 200 MB saat limit 100 MB).
  2. `SafeFetcher` memulai unduhan streaming HTTP.
  3. `boundedReadCloser` menghitung setiap byte yang dibaca ke memory buffer.
  4. Saat total byte terbaca melebihi `maxBytes`, pembacaan dihentikan seketika dan me-return sentinel error `ErrPayloadTooLarge`.
  5. Handler menangkap error `ErrPayloadTooLarge` dan menghentikan pengaliran ke ClamAV socket.
  6. Response stream ditutup rapat (`Close()`).
  7. Handler merespons ke client dengan HTTP 413 Payload Too Large:
     ```json
     {
       "success": false,
       "error": {
         "code": "FILE_TOO_LARGE",
         "message": "Remote file exceeds the configured maximum scanning size limit (100 MB)",
         "details": null
       }
     }
     ```
* **Traceability:**
  * Test: `TestFetchPayloadTooLarge` di [`internal/fetcher/client_test.go`](file:///home/ubuntu/workspace/plan/clamav-service/internal/fetcher/client_test.go) dan `TestHandler_ScanURL_FileTooLarge` di [`internal/api/handler_test.go`](file:///home/ubuntu/workspace/plan/clamav-service/internal/api/handler_test.go).

---

### 4.4. UC-URL-04: Karantina Otomatis File Remote Terinfeksi
* **Deskripsi:** Mengisolasi payload file remote bervirus ke dalam Quarantine Vault dengan enkripsi AES-256-GCM terautentikasi dan mendistribusikan notifikasi ancaman.
* **Actor:** System (Autonomous Security Engine).
* **Precondition:** File berhasil diunduh, dipindai via ClamAV socket, dan ClamAV memberikan verdict positif malware (`INFECTED`).
* **Postcondition:** File berbahaya disimpan dalam ciphertext vault `/data/quarantine` dengan magic header `VQF_AESGCM_V1\n`, log audit dicatat sebagai `INFECTED`, alert dikirimkan ke Telegram/Discord.
* **Normal Flow:**
  1. ClamAV daemon mengembalikan hasil pemindaian: `status = FOUND`, `virus_name = "Eicar-Signature"`.
  2. Handler meneruskan payload ke `Vault.QuarantineFile`.
  3. Vault membangkitkan Nonce CSPRNG 12-byte dan mengenkripsi payload file biner menggunakan cipher AES-256-GCM.
  4. Berkas terenkripsi disimpan ke disk storage karantina dengan ID unik (misal `Q-20260907-xxxx`).
  5. Catatan karantina disimpan ke tabel SQLite `quarantine_records` dengan masa retensi default 7 hari.
  6. Audit log pemindaian dicatat ke tabel `scan_audit_logs` dengan verdict `INFECTED` dan foreign key `quarantine_id`.
  7. `alert.Notifier` mendistribusikan notifikasi ancaman berformat rich Markdown ke channel Telegram/Discord terdaftar.
  8. Handler merespons client dengan HTTP 200 memuat detail ancaman dan instruksi karantina:
     ```json
     {
       "success": true,
       "verdict": "INFECTED",
       "threat": {
         "virus_name": "Eicar-Signature",
         "severity": "HIGH",
         "action_taken": "QUARANTINED",
         "quarantine_id": "Q-20260907-a1b2c3d4"
       },
       "data": {
         "source_url": "https://example.com/malware.exe",
         "file_name": "malware.exe",
         "file_size": 68,
         "file_sha256": "275a021bbfb6489e54d471899f7db9d1663fc695ec2fe2a2c4538aabf651fd0f",
         "scan_duration_ms": 45
       }
     }
     ```
* **Traceability:**
  * Test: `TestVault_AESGCM_EncryptionAndRestore` di [`internal/quarantine/vault_test.go`](file:///home/ubuntu/workspace/plan/clamav-service/internal/quarantine/vault_test.go).
  * Implementation: `handleScanURL` di [`internal/api/server.go`](file:///home/ubuntu/workspace/plan/clamav-service/internal/api/server.go).
