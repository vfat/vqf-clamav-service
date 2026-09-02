# Dokumentasi Fitur — `clamav-service`

> **Status Dokumen:** Active (Evidence-Backed)  
> **Last Updated:** 2026-09-02 19:38  
> **Scope:** Greenfield Feature Inventory & Live Codebase Implementation

---

## 1. Ringkasan Fitur

`clamav-service` adalah layanan pemindaian malware dan antivirus berbasis HTTP REST API berkinerja tinggi yang dikemas dalam satu container mandiri (*All-in-One Docker Container*). Layanan ini dirancang dengan pendekatan *Security-as-a-Service (ala Stripe)* untuk memudahkan integrasi keamanan file ke berbagai aplikasi pengirim (*consumer applications*).

---

## 2. Matriks Inventaris Fitur

| Kode Fitur | Nama Fitur | Kategori | Status | Komponen Penanggung Jawab | Bukti Implementasi & Test |
|---|---|---|---|---|---|
| **FEAT-01** | Core Synchronous File Scan | Scanning | ✅ **ACTIVE** | `api.Server`, `clamd.Client` | [`internal/api/server.go`](file:///home/ubuntu/workspace/plan/clamav-service/internal/api/server.go#L108-L230), [`internal/api/handler_test.go`](file:///home/ubuntu/workspace/plan/clamav-service/internal/api/handler_test.go) |
| **FEAT-02** | Raw Binary Chunked Stream Scan | Scanning | ✅ **ACTIVE** | `api.Server`, `clamd.Client` | [`internal/clamd/client.go`](file:///home/ubuntu/workspace/plan/clamav-service/internal/clamd/client.go#L104-L162), [`internal/clamd/client_test.go`](file:///home/ubuntu/workspace/plan/clamav-service/internal/clamd/client_test.go) |
| **FEAT-03** | Remote URL / Cloud Object Scan | Scanning | ⏳ `PLANNED` | `api.Server`, `clamd.Client` | Backlog enhancement |
| **FEAT-04** | Async Scan Job & Webhook Callback | Async Queue | ⏳ `PLANNED` | `api.Server`, `alert.Notifier` | Backlog enhancement |
| **FEAT-05** | Password-Protected Archive Inspection | Resilience | ⏳ `PLANNED` | `scanner.ArchiveInspector` | Backlog enhancement (`TDD-009`) |
| **FEAT-06** | Zip-Bomb & Decompression Limiter | Resilience | ⏳ `PLANNED` | `scanner.ArchiveInspector` | Backlog enhancement (`TDD-009`) |
| **FEAT-07** | Built-in Quarantine Vault | Quarantine | ✅ **ACTIVE** | `quarantine.Vault` | [`internal/quarantine/vault.go`](file:///home/ubuntu/workspace/plan/clamav-service/internal/quarantine/vault.go#L35-L84), [`internal/quarantine/vault_test.go`](file:///home/ubuntu/workspace/plan/clamav-service/internal/quarantine/vault_test.go) |
| **FEAT-08** | Dual Mode File Restore & Whitelisting | Quarantine | ✅ **ACTIVE** | `quarantine.Vault`, `storage.DB` | [`internal/quarantine/vault.go`](file:///home/ubuntu/workspace/plan/clamav-service/internal/quarantine/vault.go#L86-L119), [`internal/storage/sqlite.go`](file:///home/ubuntu/workspace/plan/clamav-service/internal/storage/sqlite.go#L254-L290) |
| **FEAT-09** | Short-Cycle Auto-Purge & Retention | Data Lifecycle | ✅ **ACTIVE** | `storage.DB`, `cmd/server/main.go` | [`internal/storage/sqlite.go`](file:///home/ubuntu/workspace/plan/clamav-service/internal/storage/sqlite.go#L210-L218), [`cmd/server/main.go`](file:///home/ubuntu/workspace/plan/clamav-service/cmd/server/main.go#L119-L135) |
| **FEAT-10** | Multi-Channel Alerting (Telegram/Discord/Email) | Alerting | ✅ **ACTIVE** | `alert.Notifier` | [`internal/alert/notifier.go`](file:///home/ubuntu/workspace/plan/clamav-service/internal/alert/notifier.go#L53-L162), [`internal/alert/notifier_test.go`](file:///home/ubuntu/workspace/plan/clamav-service/internal/alert/notifier_test.go) |
| **FEAT-11** | Alert Flood Throttling & Anti-Spam | Alerting | ✅ **ACTIVE** | `alert.Notifier` | [`internal/alert/notifier.go`](file:///home/ubuntu/workspace/plan/clamav-service/internal/alert/notifier.go#L70-L93), [`internal/alert/notifier_test.go`](file:///home/ubuntu/workspace/plan/clamav-service/internal/alert/notifier_test.go#L59-L90) |
| **FEAT-12** | 2-Tier Rate Limiting & Server Protection | Traffic Control | ✅ **ACTIVE** | `ratelimit.Limiter`, `api.Server` | [`internal/ratelimit/limiter.go`](file:///home/ubuntu/workspace/plan/clamav-service/internal/ratelimit/limiter.go#L27-L86), [`internal/ratelimit/limiter_test.go`](file:///home/ubuntu/workspace/plan/clamav-service/internal/ratelimit/limiter_test.go) |
| **FEAT-13** | Zero-Touch Master Key & AES-256-GCM | Security | ✅ **ACTIVE** | `crypto.Keygen`, `crypto.AES` | [`internal/crypto/keygen.go`](file:///home/ubuntu/workspace/plan/clamav-service/internal/crypto/keygen.go), [`internal/crypto/aes.go`](file:///home/ubuntu/workspace/plan/clamav-service/internal/crypto/aes.go) |
| **FEAT-14** | Embedded Web Admin Dashboard (SPA) | Management | ⏳ `IN-PROGRESS` | `web/ui`, `api.Server` | Next target step |
| **FEAT-15** | Streaming Audit Log Exporter (CSV/JSON) | Compliance | ✅ **ACTIVE** | `api.Server`, `storage.DB` | [`internal/api/server.go`](file:///home/ubuntu/workspace/plan/clamav-service/internal/api/server.go#L248-L253), [`internal/storage/sqlite.go`](file:///home/ubuntu/workspace/plan/clamav-service/internal/storage/sqlite.go#L149-L208) |
| **FEAT-16** | Go Native Process Supervisor (PID 1) | Operations | ✅ **ACTIVE** | `cmd/server/main.go` | [`cmd/server/main.go`](file:///home/ubuntu/workspace/plan/clamav-service/cmd/server/main.go#L101-L117) |
| **FEAT-17** | Health & Readiness Probe / Metrics | Observability | ✅ **ACTIVE** | `api.Server` | [`internal/api/server.go`](file:///home/ubuntu/workspace/plan/clamav-service/internal/api/server.go#L95-L106) |

---

## 3. Rincian Implementasi Fitur Aktif

### 3.1. Core Synchronous & Stream Scanning (`FEAT-01`, `FEAT-02`)
* Menerima payload file multipart melalui `POST /api/v1/scan/file` atau stream biner langsung melalui `POST /api/v1/scan/stream`.
* Mengalirkan byte stream via Unix Domain Socket `/var/run/clamav/clamd.ctl` menggunakan protokol ClamAV `zINSTREAM\0` secara chunked (64 KB).
* Mengembalikan struktur JSON standar: `verdict` (`CLEAN` / `INFECTED`), `threat` details, `file_sha256`, dan `scan_duration_ms`.

### 3.2. Quarantine Vault & Safe Isolation (`FEAT-07`, `FEAT-08`, `FEAT-09`)
* File malware otomatis diisolasi ke `/data/quarantine/` dengan format nama netral `Q-YYYYMMDD-ULID.quarantine`.
* File mengalami XOR binary scrambling dan disimpan dengan permission `0600`.
* Restore file didukung via `POST /api/v1/quarantine/restore`, secara otomatis mendaftarkan hash SHA-256 ke tabel `whitelist_signatures`.
* Background worker membersihkan log audit lama (> 3 hari) dan file karantina kadaluarsa (> 7 hari).

### 3.3. Multi-Channel Alerting & Throttling (`FEAT-10`, `FEAT-11`)
* Notifikasi ancaman seketika terkirim ke bot Telegram dan webhook Discord saat virus ditemukan.
* Dilengkapi algoritma sliding-window 60 detik (*Anti-Spam Flood Throttling*) yang mencegah spam saat terjadi lonjakan deteksi malware (> 5 malware/menit).

### 3.4. Rate Limiting, Master Key, & Security (`FEAT-12`, `FEAT-13`)
* Algoritma *token bucket* membatasi request per client IP / API key dan mengembalikan header `X-RateLimit-*`.
* First boot otomatis men-generate master encryption key (`clam_sec_...`) dan menginjeksi ke `.env`.
* Kredensial sensitif dienkripsi menggunakan cipher simetris **AES-256-GCM** dengan nonce 96-bit.
