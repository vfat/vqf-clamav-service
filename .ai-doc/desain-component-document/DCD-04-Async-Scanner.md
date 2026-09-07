# DCD-04-Async-Scanner

## 1. Metadata Komponen

| Field | Keterangan |
|---|---|
| **Komponen** | Asynchronous Scan Job Queue & Safe Webhook Dispatcher |
| **Kode Dokumen** | `DCD-04` |
| **Dasar Rancangan** | [`SCD-Async-Scanner.md`](file:///home/ubuntu/workspace/plan/clamav-service/.ai-doc/plan/component/SCD-Async-Scanner.md) |
| **Status Implementasi** | 🟡 `In-Progress / TDD Phase` (`TDD-015` & `TDD-016`) |
| **Target Audience** | Backend Engineers, Security Architects, Integration Partners |

---

## 2. Object Identification

Berdasarkan rancangan arsitektur pada [`internal/asyncscan/`](file:///home/ubuntu/workspace/plan/clamav-service/internal/asyncscan/), [`internal/api/`](file:///home/ubuntu/workspace/plan/clamav-service/internal/api/), dan [`internal/storage/`](file:///home/ubuntu/workspace/plan/clamav-service/internal/storage/):

### Boundary
* **HTTP Ingress Router:** `http.ServeMux` pada [`internal/api/server.go`](file:///home/ubuntu/workspace/plan/clamav-service/internal/api/server.go) menangani route:
  - `POST /api/v1/scan/async`
  - `GET /api/v1/scan/jobs/{id}`
* **Spool Storage Filesystem:** Direktori penyimpanan berkas sementara (`SPOOL_DIR`, default `/data/spool/`).
* **Outbound Webhook HTTP Transport:** `fetcher.SafeFetcher` dengan Anti-SSRF protection untuk mengirim HTTP POST payload ke `callback_url`.
* **Unix Domain Socket Stream Interface:** `net.Conn` Unix socket ClamAV daemon (`/var/run/clamav/clamd.ctl`).

### Control
* **`api.handleScanAsync`:** HTTP handler penerima upload multipart dan callback URL, memvalidasi parameter, membuat spool file, mendaftarkan job ke DB dan mem-push ke queue channel, lalu mengembalikan HTTP 202 Accepted.
* **`api.handleScanJobStatus`:** HTTP handler yang melayani query polling status job berdasarkan `job_id`.
* **`asyncscan.Service`:** Orkestrator utama yang mengelola queue channel, worker pool goroutines, dan siklus hidup job.
* **`asyncscan.Worker`:** Background goroutine yang membaca job dari channel, membaca berkas spool, memindai ke `clamd`, mengisolasi ke vault jika infected, mencatat audit log, memperbarui DB status, dan memicu webhook dispatcher.
* **`asyncscan.WebhookDispatcher`:** Controller pengirim notifikasi hasil scan ke `callback_url` dengan validasi skema, Anti-SSRF filtering, timeout 10 detik, dan mekanisme retry.

### Entity
* **Endpoint:**
  - `POST /api/v1/scan/async`
  - `GET /api/v1/scan/jobs/{id}`
* **DTO `scanAsyncResponse`:** `{ "success": true, "status": "ACCEPTED", "job_id": string, "message": string }`.
* **Domain Model `storage.ScanJob`:** Record tabel `scan_jobs` SQLite (`id`, `file_name`, `file_size`, `file_sha256`, `callback_url`, `consumer`, `status`, `verdict`, `virus_name`, `error_msg`, `created_at`, `updated_at`).
* **Database Table `scan_jobs`:** Skema tabel SQLite dengan indeks `idx_jobs_status`.

---

## 3. Use Case List

| No | Kode Use Case | Nama Use Case | Actor | Status | Referensi Detail |
|---|---|---|---|---|---|
| 1 | `UC-ASYNC-01` | Pengajuan Job Pemindaian Asinkron | Consumer Application | In-Progress | Section 4.1 |
| 2 | `UC-ASYNC-02` | Eksekusi Pemindaian Latar Belakang | Background Worker Pool | In-Progress | Section 4.2 |
| 3 | `UC-ASYNC-03` | Pengiriman Hasil via Webhook Anti-SSRF | Webhook Dispatcher | In-Progress | Section 4.3 |
| 4 | `UC-ASYNC-04` | Query Status Scan Job | Consumer Application | In-Progress | Section 4.4 |

---

## 4. Use Case Detail

### 4.1. UC-ASYNC-01: Pengajuan Job Pemindaian Asinkron
* **Deskripsi:** Menerima unggahan berkas multipart dan URL webhook target, memvalidasi URL callback (skema http/https), menyimpan berkas sementara ke spool storage, mendaftarkan job ke tabel `scan_jobs` dengan status `QUEUED`, memasukkan ke antrean worker, dan mengembalikan HTTP 202 Accepted beserta `job_id`.
* **Actor:** Consumer Application.
* **Precondition:** Layanan `clamav-service` aktif, parameter `file` dan `callback_url` disertakan.
* **Postcondition:** Berkas spool tersimpan, record `scan_jobs` berstatus `QUEUED`, respons HTTP 202 dikembalikan.

### 4.2. UC-ASYNC-02: Eksekusi Pemindaian Latar Belakang
* **Deskripsi:** Worker membaca job dari antrean, memperbarui status menjadi `PROCESSING`, membaca berkas spool, menghitung hash SHA-256, mengalirkan biner ke daemon ClamAV, mengisolasi file ke karantina jika terdeteksi malware, mencatat audit log, menghapus berkas spool dari disk, dan memperbarui status job menjadi `COMPLETED` (atau `FAILED`).
* **Actor:** Background Worker.
* **Precondition:** Worker aktif dan menerima job dari queue channel.
* **Postcondition:** Pemindaian selesai, record audit log tercatat di SQLite, berkas spool dihapus.

### 4.3. UC-ASYNC-03: Pengiriman Hasil via Webhook Anti-SSRF
* **Deskripsi:** Setelah pemindaian selesai, Webhook Dispatcher menyusun payload hasil pemindaian dan mengirimkannya via HTTP POST ke `callback_url`. Validasi Anti-SSRF dijalankan untuk menolak target IP privat atau loopback. Jika pengiriman gagal, dilakukan retry hingga 3 kali.
* **Actor:** Webhook Dispatcher.
* **Precondition:** Job pemindaian telah mencapai verdict final (`CLEAN`, `INFECTED`, atau `ERROR`).
* **Postcondition:** Webhook terkirim ke endpoint consumer.

### 4.4. UC-ASYNC-04: Query Status Scan Job
* **Deskripsi:** Mengambil detail informasi status pemrosesan scan job dan hasilnya berdasarkan `job_id`. Berguna sebagai mekanisme polling bagi consumer yang tidak dapat menerima webhook masuk.
* **Actor:** Consumer Application.
* **Precondition:** `job_id` terdaftar di database.
* **Postcondition:** Detail status job dan verdict dikembalikan dalam format JSON (HTTP 200 OK).
