# Discussion Summary — Arsitektur & Spesifikasi Lengkap `clamav-service`

| Field | Value |
|-------|-------|
| **Sesi** | Perancangan Arsitektur, Fitur, & Spesifikasi `clamav-service` |
| **Area** | Feature / Greenfield Architecture Design |
| **Sub-Topik** | Core Scanning Engine, Quarantine System, Containerization, Alerts, Rate Limiting, Audit, & Process Supervision |
| **Tanggal** | 2026-09-02 |
| **Kategori** | Ideation, Solution Architecture, Reliability, & Security Engineering |
| **Teknik** | First Principles, Possibility Mapping, Comparative Trade-off, Lifecycle State Modeling, Error Contract Standardization |
| **Peserta** | User, Sherin ✋, Melon 🏗️, Sultan ⚙️, Nindi 🔬, Lugi 📊, Bernadya ✨ |
| **Jumlah Ide** | 28 ide terpetakan & disepakati |

---

## Konteks
Sesi brainstorming greenfield komprehensif bersama Full Team untuk merancang arsitektur, spesifikasi, dan detail teknis `clamav-service`. Diskusi mencakup seluruh aspek mulai dari engine scanning streaming, proteksi memori VPS, sistem karantina, antarmuka Web UI, database SQLite WAL, sistem notifikasi Telegram/Discord/Email, otomatisasi enkripsi key, pembatasan rate limit, audit logging & auto-purge retensi, pengawasan proses container Go, hingga kontrak response JSON.

---

## Rincian Ide yang Teridentifikasi

### 1. Unified Security Inspection Layer & Blast Radius Isolation
Memisahkan pemindaian antivirus ke dalam container independen agar aplikasi bisnis utama terhindar dari crash, OOM, atau serangan eksploitasi file.
> **Sumber:** Melon 🏗️ & Nindi 🔬 | **Status:** `approved` | **Kategori:** Architecture

### 2. High-Throughput Streaming Antivirus Gateway
Menerima raw binary / multipart stream via HTTP dan meneruskannya via Unix Domain Socket (`/tmp/clamd.sock`) langsung ke RAM daemon `clamd` tanpa disk I/O perantara.
> **Sumber:** Sultan ⚙️ | **Status:** `approved` | **Kategori:** Backend Performance

### 3. Pluggable Security-as-a-Service (ala Stripe)
Integrasi instan dengan model Sync Direct Scan (verdict dalam hitungan milidetik) dan Async Cloud Scan (kirim URL + webhook callback).
> **Sumber:** Bernadya ✨ | **Status:** `approved` | **Kategori:** Product/DX

### 4. Agentless Threat Intelligence & Audit Hub
Pencatatan telemetri scan (SHA256, virus signature, latensi ms, status) tanpa perlu menginstal agent pada host/laptop.
> **Sumber:** Lugi 📊 | **Status:** `approved` | **Kategori:** Data/Observability

### 5. Built-in Quarantine Vault dengan Safe Storage Pattern
Penyimpanan file terinfeksi di volume lokal `/data/quarantine/` dengan nama netral (`.quarantine`), izin ketat `0600`, dan scrambled/encrypted.
> **Sumber:** Sultan ⚙️ & Nindi 🔬 | **Status:** `approved` | **Kategori:** Security

### 6. Dual Mode Restore & Anti Re-Quarantine Hash Whitelisting
Rilis file false positive via Direct Download atau Push balik ke S3 / Webhook, disertai auto-whitelist SHA256 agar file tidak terkena karantina ulang.
> **Sumber:** Nindi 🔬 & Melon 🏗️ | **Status:** `approved` | **Kategori:** Operational Workflow

### 7. All-in-One Single Container Deployment
Mengemas Go API wrapper, `clamd`, `freshclam`, embedded UI, dan SQLite dalam 1 Docker container dengan volume persisten `/data`.
> **Sumber:** Melon 🏗️ & Sultan ⚙️ | **Status:** `approved` | **Kategori:** Infrastructure

### 8. In-Process SQLite (WAL Mode) Database
Penyimpanan metadata karantina, log audit, API keys, dan settings dengan SQLite WAL mode. Zero extra RAM overhead, single file backup di `/data/clamav-service.db`.
> **Sumber:** Lugi 📊 | **Status:** `approved` | **Kategori:** Database

### 9. AES-256-GCM Field-Level Encryption & Key Hashing
Enkripsi simetris untuk kredensial sensitif di database menggunakan master key dari Environment Variable, dan salted hash untuk API Keys.
> **Sumber:** Nindi 🔬 | **Status:** `approved` | **Kategori:** Security

### 10. Zip-Bomb Limits & Archive Inspection Safeguards
Membatasi rekursi unzip maks 5 level, maks 1.000 file, maks 250 MB ekstraksi memori, dan batas timeout 30 detik.
> **Sumber:** Nindi 🔬 | **Status:** `approved` | **Kategori:** Resilience

### 11. Password-Protected Archive Handling (ala VirusTotal)
Return status `UNSCANNABLE (PASSWORD_PROTECTED)`. Mendukung input `archive_password` di API dan auto-try kamus password malware umum.
> **Sumber:** Melon 🏗️ & Lugi 📊 | **Status:** `approved` | **Kategori:** Security

### 12. Multi-Channel Alert Manager (Telegram, Discord, Email)
Pengiriman notifikasi otomatis saat malware terdeteksi, dengan auto-detect kredensial dan toggle switch ON/OFF di Web UI.
> **Sumber:** Bernadya ✨ & Sultan ⚙️ | **Status:** `approved` | **Kategori:** Alerting

### 13. Alert Throttling & Flood Control
Proteksi anti-spam pesan saat lonjakan serangan malware massal dengan mengirimkan batch summary digest jika > 5 malware/menit.
> **Sumber:** Nindi 🔬 | **Status:** `approved` | **Kategori:** Alerting

### 14. Zero-Touch Encryption Key Generation
Sistem otomatis men-generate key dan langsung menuliskan `ENCRYPTION_KEY=...` ke file `.env` saat boot/init pertama. User menyalin key hanya untuk keperluan backup darurat.
> **Sumber:** User, Sultan ⚙️, & Melon 🏗️ | **Status:** `approved` | **Kategori:** Security/DX

### 15. 2-Tier Rate Limiting & API Key Policy
Proteksi server global di `.env` (misal 50 RPS) + aturan spesifik per API Key di UI (RPM, daily quota, IP whitelist, permissions) dengan in-memory token bucket.
> **Sumber:** Sultan ⚙️ & Lugi 📊 | **Status:** `approved` | **Kategori:** Traffic Management

### 16. Short-Cycle Audit & Quarantine Retention
Default retensi log audit scan = 3 hari, default retensi file karantina = 7 hari, dilengkapi auto-vacuum SQLite mingguan.
> **Sumber:** User, Lugi 📊, & Sultan ⚙️ | **Status:** `approved` | **Kategori:** Data Governance

### 17. Streaming CSV & JSON Audit Exporter
Download laporan riwayat pemindaian dengan filter fleksibel via Web UI dan REST API secara streaming tanpa lonjakan RAM.
> **Sumber:** Bernadya ✨ & Lugi 📊 | **Status:** `approved` | **Kategori:** Compliance

### 18. Go Native Process Supervisor (PID 1)
Binary Go bertindak sebagai master entrypoint yang men-spawn, memantau, dan me-restart `clamd` dan `freshclam`, serta menangani graceful shutdown OS signal.
> **Sumber:** Sultan ⚙️ & Melon 🏗️ | **Status:** `approved` | **Kategori:** Process Management

### 19. Standarisasi JSON Response & Error Contracts
Struktur konsisten untuk balasan `CLEAN`, `INFECTED`, `SUSPICIOUS`, `UNSCANNABLE`, `ACCEPTED`, dan inventaris error code standar.
> **Sumber:** Bernadya ✨ & Sultan ⚙️ | **Status:** `approved` | **Kategori:** API Design

---

## Keputusan Final

- ✅ **Bahasa & Runtime:** **Go (Golang)** sebagai API Wrapper & Process Supervisor (RAM ~15–30 MB).
- ✅ **Engine Antivirus:** ClamAV Daemon (`clamd`) via Unix Domain Socket + Freshclam background updater.
- ✅ **Packaging:** All-in-One Single Docker Container (Alpine-based, volume `/data`).
- ✅ **Database:** SQLite (WAL Mode) di `/data/clamav-service.db`.
- ✅ **Antarmuka:** REST API (Port 8080) + Embedded Web Admin UI (SPA).
- ✅ **Keamanan:** AES-256-GCM field encryption, salted hash API keys, Zero-Touch key generation ke `.env`.
- ✅ **Karantina & Restore:** Built-in Vault `/data/quarantine/`, dual restore (Download/Push), SHA256 whitelisting.
- ✅ **Alerting:** Telegram, Discord, Email dengan Throttling Anti-Spam.
- ✅ **Retensi:** 3 Hari Log Audit, 7 Hari Karantina.

---
*Generated by AI Documentor — Brainstorming Add-On (Session Completed)*
