# SCD-Async-Scanner (Single Component Design)

## 1. Context

### 1.1. Masalah yang Ingin Diselesaikan
Pada pemindaian sinkron (`POST /api/v1/scan/file`, `POST /api/v1/scan/url`), koneksi HTTP client dipertahankan tetap terbuka (*blocking*) sampai proses pemindaian ClamAV selesai. Pada pemrosesan berkas berukuran besar (misalnya file arsip 50-500 MB) atau saat beban lalu lintas tinggi, penahanan koneksi ini dapat menyebabkan:
1. HTTP gateway/reverse proxy timeout (misalnya batas 60 detik Cloudflare/Nginx/AWS ALB).
2. *Connection exhaustion* pada thread pool aplikasi pengirim (*consumer*).
3. Kegagalan proses akibat fluktuasi koneksi jaringan client.

Dibutuhkan arsitektur **Asynchronous Scan Job Queue & Webhook Delivery Engine** (`Async-Scanner`) untuk memproses pemindaian file di latar belakang tanpa menahan koneksi HTTP client.

### 1.2. Posisi Komponen dalam Sistem

```
[Consumer Application]
       │
       │ 1. POST /api/v1/scan/async (Multipart + callback_url)
       ▼
┌──────────────────────────────────────────────────────────────┐
│  clamav-service                                              │
│                                                              │
│  ┌────────────────────────────────────────────────────────┐  │
│  │ Async-Scanner Component (internal/asyncscan)           │  │
│  │  1. Spool Manager (/data/spool/)                       │  │
│  │  2. SQLite Job Registry (table: scan_jobs)             │  │
│  │  3. Worker Pool (Concurrent background scanners)       │  │
│  │  4. Safe Webhook Dispatcher (Anti-SSRF Protected)      │  │
│  └────────────────────────────────────────────────────────┘  │
│         │                                                    │
│         ├───────────────► 202 Accepted { "job_id": "..." }   │
│         │                 (dikembalikan seketika)            │
│         ▼ (Worker streams to socket)                         │
│  ┌───────────────────────┐                                   │
│  │ ClamAV clamd Engine   │                                   │
│  └───────────────────────┘                                   │
│         │                                                    │
│         ▼ 2. POST callback_url (Verdict payload)             │
│  [Consumer Webhook Endpoint]                                 │
└──────────────────────────────────────────────────────────────┘
```

### 1.3. Hubungan dengan Komponen Lain
* **`api.Server`**: Menyediakan endpoint ingress `POST /api/v1/scan/async` dan polling status `GET /api/v1/scan/jobs/{id}`.
* **`clamd.Client`**: Menerima streaming biner berkas dari background worker via Unix domain socket.
* **`fetcher.SafeFetcher`**: Mengirimkan payload webhook hasil pemindaian ke `callback_url` dengan proteksi **Anti-SSRF** ketat (mencegah penyerang menyalahgunakan webhook callback untuk membidik IP internal/cloud metadata).
* **`quarantine.Vault`**: Mengisolasi file malware jika verdict bernilai `INFECTED`.
* **`storage.DB`**: Menyimpan status job pada tabel `scan_jobs` dan mencatat histori pemindaian pada `scan_audit_logs`.
* **`alert.Notifier`**: Mengirimkan notifikasi instan jika file terdeteksi malware.

---

## 2. Scope

### 2.1. In-Scope
1. **Spool File Management (`/data/spool`):**
   - Menyimpan payload file unggahan sementara ke direktori spool container secara aman (`0600`).
   - Membersihkan berkas spool secara otomatis setelah proses scan dan pengiriman webhook selesai.
2. **Persistent Job Registry (`scan_jobs`):**
   - Mencatat ID unik job, nama file, ukuran, hash SHA-256, URL callback, status (`QUEUED`, `PROCESSING`, `COMPLETED`, `FAILED`), verdict, dan error message.
3. **Concurrent Worker Pool:**
   - Channel-based worker pool yang membatasi konkurensi pemindaian background (default: 2 worker).
   - Mendukung graceful termination (menyelesaikan job aktif sebelum proses mati).
4. **Anti-SSRF Protected Webhook Dispatcher:**
   - Memvalidasi skema `http://` / `https://`.
   - Menggunakan `fetcher.SafeFetcher` untuk menolak callback ke private IP RFC1918, loopback, dan cloud metadata `169.254.169.254`.
   - Melakukan retry pengiriman webhook (hingga 3 kali) dengan timeout terkendali.
5. **REST API Endpoints:**
   - `POST /api/v1/scan/async`: Mengunggah file + callback URL, mengembalikan `202 Accepted` dengan `job_id`.
   - `GET /api/v1/scan/jobs/{id}`: Menampilkan status terbaru dari scan job (polling alternative).

### 2.2. Out-of-Scope
* Integrasi message broker eksternal (Kafka/RabbitMQ/Redis) — implementasi menggunakan arsitektur pure Go goroutines & SQLite WAL queue untuk menjaga prinsip zero-external-dependency dan memori minimal (~25 MB RAM).

---

## 3. Prerequisite

1. Direktori `/data/spool` tersedia untuk penyimpanan file sementara.
2. Skema SQLite baru untuk tabel `scan_jobs` pada `internal/storage`.
3. Komponen `fetcher.SafeFetcher` siap digunakan untuk pengiriman webhook dengan Anti-SSRF.

---

## 4. Daftar Usecase

| Kode Usecase | Nama Usecase | Deskripsi Singkat |
|---|---|---|
| **`UC-ASYNC-01`** | Pengajuan Job Pemindaian Asinkron | Menerima file multipart dan callback URL, menyimpan ke spool, mencatat ke DB, dan mengembalikan HTTP 202 Accepted dengan `job_id`. |
| **`UC-ASYNC-02`** | Eksekusi Pemindaian Latar Belakang | Worker mengambil job dari queue, memindai aliran berkas ke clamd, mengisolasi file jika infected, mencatat audit log, dan menghapus berkas spool. |
| **`UC-ASYNC-03`** | Pengiriman Hasil via Webhook Anti-SSRF | Mengirimkan payload JSON hasil pemindaian ke callback URL menggunakan transport Anti-SSRF dengan mekanisme retry. |
| **`UC-ASYNC-04`** | Query Status Scan Job | Menyediakan endpoint polling untuk memeriksa status dan hasil pemindaian job berdasarkan `job_id`. |

---

## 5. Rencana Skema Database (`scan_jobs`)

```sql
CREATE TABLE IF NOT EXISTS scan_jobs (
    id           TEXT PRIMARY KEY,
    file_name    TEXT NOT NULL,
    file_size    INTEGER NOT NULL,
    file_sha256  TEXT,
    callback_url TEXT NOT NULL,
    consumer     TEXT,
    status       TEXT NOT NULL, -- QUEUED, PROCESSING, COMPLETED, FAILED
    verdict      TEXT,          -- CLEAN, INFECTED, ERROR
    virus_name   TEXT,
    error_msg    TEXT,
    created_at   DATETIME NOT NULL,
    updated_at   DATETIME NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_jobs_status ON scan_jobs(status);
```

---

## 6. Rencana Siklus TDD

1. **`TDD-015`**: Core Async Queue, Spool Management, Worker Pool & Safe Webhook Dispatcher (`internal/asyncscan`).
2. **`TDD-016`**: REST API Ingress Handler `POST /api/v1/scan/async` & Status Poller `GET /api/v1/scan/jobs/{id}` (`internal/api`).
