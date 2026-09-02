# Dokumentasi Fitur — `clamav-service`

> **Status Dokumen:** Approved Feature Map  
> **Tanggal Pembuatan:** 2026-09-02  
> **Scope:** Greenfield Feature Inventory & Mapping

---

## 1. Ringkasan Fitur

`clamav-service` adalah layanan pemindaian malware dan antivirus berbasis HTTP REST API berkinerja tinggi yang dikemas dalam satu container mandiri (*All-in-One Docker Container*). Layanan ini dirancang dengan pendekatan *Security-as-a-Service (ala Stripe)* untuk memudahkan integrasi keamanan file ke berbagai aplikasi pengirim (*consumer applications*).

---

## 2. Matriks Inventaris Fitur

| Kode Fitur | Nama Fitur | Kategori | Status | Komponen Penanggung Jawab | Deskripsi Singkat |
|---|---|---|---|---|---|
| **FEAT-01** | Core Synchronous File Scan | Scanning | `PLANNED` (MVP) | `ScanController`, `ClamdBridge` | Pemindaian langsung via multipart file upload dengan instant verdict (`CLEAN`, `INFECTED`). |
| **FEAT-02** | Raw Binary Chunked Stream Scan | Scanning | `PLANNED` (MVP) | `ScanController`, `ClamdBridge` | Pemindaian raw binary stream tanpa multipart overhead via protokol socket `INSTREAM`. |
| **FEAT-03** | Remote URL / Cloud Object Scan | Scanning | `PLANNED` (MVP) | `ScanController`, `ClamdBridge` | Mengunduh stream dari URL publik / presigned S3 URL dan memindainya secara on-the-fly. |
| **FEAT-04** | Async Scan Job & Webhook Callback | Async Queue | `PLANNED` (MVP) | `ScanController`, `WebhookDispatcher` | Menerima job scan (return `202 Accepted` + `job_id`) dan menembakkan hasil ke URL callback. |
| **FEAT-05** | Password-Protected Archive Inspection | Resilience | `PLANNED` (MVP) | `ScanController`, `ArchiveInspector` | Inspeksi arsip terenkripsi dengan parameter `archive_password`, kamus password, atau verdict `UNSCANNABLE`. |
| **FEAT-06** | Zip-Bomb & Decompression Limiter | Resilience | `PLANNED` (MVP) | `ScanController`, `ArchiveInspector` | Proteksi DoS dengan membatasi rekursi (5 level), jumlah file (1.000), dan batas ekstrak (250 MB). |
| **FEAT-07** | Built-in Quarantine Vault | Quarantine | `PLANNED` (MVP) | `QuarantineManager` | Isolasi file terinfeksi ke folder `/data/quarantine/` dengan nama netral, izin `0600`, dan scrambled biner. |
| **FEAT-08** | Dual Mode File Restore & Whitelisting | Quarantine | `PLANNED` (MVP) | `QuarantineManager` | Pemulihan file false-positive (Direct Admin Download / S3 Push) dan auto-whitelist SHA256. |
| **FEAT-09** | Short-Cycle Auto-Purge & Retention | Data Lifecycle | `PLANNED` (MVP) | `SQLiteManager`, `CleanerWorker` | Pembersihan otomatis file karantina (7 hari) dan log audit (3 hari) dengan auto-vacuum SQLite. |
| **FEAT-10** | Multi-Channel Alerting (Telegram/Discord/Email) | Alerting | `PLANNED` (MVP) | `AlertDispatcher` | Pengiriman notifikasi seketika saat terdeteksi malware ke channel Telegram, Discord, dan Email. |
| **FEAT-11** | Alert Flood Throttling & Anti-Spam | Alerting | `PLANNED` (MVP) | `AlertDispatcher` | Penggabungan notifikasi ke Batch Digest Summary saat lonjakan serangan > 5 deteksi/menit. |
| **FEAT-12** | 2-Tier Rate Limiting & API Key Policy | Traffic Control | `PLANNED` (MVP) | `AuthRateLimiter` | Proteksi server global (`.env`) + aturan granular RPM, kuota, dan IP whitelist per API Key. |
| **FEAT-13** | Zero-Touch Master Key & AES-256-GCM | Security | `PLANNED` (MVP) | `CryptoKeygen`, `CryptoAES` | Otomatisasi generate key ke `.env` saat first boot dan enkripsi field kredensial di SQLite. |
| **FEAT-14** | Embedded Web Admin Dashboard (SPA) | Management | `PLANNED` (MVP) | `AdminWebUI`, `APIGateway` | Web dashboard mandiri untuk monitor daemon, drag-and-drop test scan, inspeksi vault, dan export data. |
| **FEAT-15** | Streaming Audit Log Exporter (CSV/JSON) | Compliance | `PLANNED` (MVP) | `SQLiteManager`, `APIGateway` | Download riwayat log audit dengan filter fleksibel secara streaming tanpa membebani memori server. |
| **FEAT-16** | Go Native Process Supervisor (PID 1) | Operations | `PLANNED` (MVP) | `ProcessSupervisor` | Pengawasan proses `clamd` & `freshclam`, auto-restart, graceful shutdown, dan zero-downtime signature reload. |
| **FEAT-17** | Health & Readiness Probe / Metrics | Observability | `PLANNED` (MVP) | `APIGateway`, `MetricsExporter` | Endpoint `/healthz` dan `/metrics` untuk liveness/readiness probe Kubernetes dan Prometheus. |

---

## 3. Rincian Fitur Utama

### 3.1. Fitur Core Scanning & Resilience (FEAT-01 s/d FEAT-06)
* **Kebutuhan Bisnis:** Memberikan verdict pemindaian file yang cepat (<100ms) tanpa risiko crash engine akibat serangan *decompression bomb* atau arsip terenkripsi.
* **Integrasi Teknis:**
  - Mengalirkan byte stream via Unix Domain Socket `/tmp/clamd.sock` menggunakan perintah ClamAV `zINSTREAM\0`.
  - Mengembalikan struktur data JSON standar dengan hash SHA-256, latensi scan, dan versi signature.
  - Memverifikasi status arsip ber-password dan mengembalikan status `UNSCANNABLE (PASSWORD_PROTECTED)` jika arsip terkunci.

### 3.2. Fitur Quarantine Vault & Restore (FEAT-07 s/d FEAT-09)
* **Kebutuhan Bisnis:** Mencegah penyebaran file malware ke sistem lain tanpa menghilangkan bukti atau dokumen penting jika terjadi *false positive*.
* **Integrasi Teknis:**
  - File malware dinetralkan (format `Q-YYYYMMDD-ULID.quarantine`) dengan izin akses `0600`.
  - Restore manual: Admin mendownload file asli yang telah di-descramble.
  - Restore otomatis: Service mengunggah file ke target S3 dan menembak webhook event.
  - Hash SHA-256 file didaftarkan ke tabel `whitelist_signatures` agar tidak terkena karantina ulang.
  - Pembersihan otomatis berbasis TTL: 7 hari untuk file karantina dan 3 hari untuk log audit.

### 3.3. Fitur Alerting, Rate Limiting, & Security (FEAT-10 s/d FEAT-13)
* **Kebutuhan Bisnis:** Observabilitas real-time terhadap ancaman keamanan dan proteksi server dari beban berlebih (fair usage).
* **Integrasi Teknis:**
  - Alert instan via bot Telegram, Discord webhook, dan email SMTP dengan proteksi anti-spam flood throttling.
  - In-memory token bucket rate limiting mengembalikan header `X-RateLimit-*` dan HTTP `429 Too Many Requests`.
  - Zero-touch master key generation otomatis menulis `ENCRYPTION_KEY` ke `.env` pada booting pertama.
  - Kredensial webhook dan token bot tersimpan terenkripsi dengan AES-256-GCM di SQLite.

### 3.4. Fitur Admin Dashboard, Supervisi, & Ops (FEAT-14 s/d FEAT-17)
* **Kebutuhan Bisnis:** Kemudahan pengoperasian bagi DevOps & SecOps tanpa ketergantungan software supervisor eksternal.
* **Integrasi Teknis:**
  - Binary Go bertindak sebagai Master Process (PID 1) yang men-spawn subprocess `clamd` dan `freshclam`.
  - Dashboard Web UI (SPA) disajikan langsung dari binary Go di port 8080.
  - Export data riwayat scan ke format CSV dan JSON secara streaming.
  - Endpoint health probe `/healthz` dan Prometheus `/metrics`.
