# Minutes of Meeting — Brainstorming Session (Final Wrap-Up)

| Field | Value |
|-------|-------|
| **Tanggal** | 2026-09-02 |
| **Waktu** | 14:28 - 17:15 WIB |
| **Topik Utama** | Perancangan Arsitektur, Fitur, & Spesifikasi Lengkap `clamav-service` |
| **Area** | Feature / Greenfield Architecture Design |
| **Peserta** | User, Sherin ✋ (Moderator), Melon 🏗️ (Technical Architect), Sultan ⚙️ (Senior Backend Engineer), Nindi 🔬 (Problem Solver), Lugi 📊 (Data Specialist), Bernadya ✨ (Creative Visionary) |
| **Teknik** | First Principles, Possibility Mapping, Comparative Trade-off, Lifecycle State Modeling, Error Contract Standardization |

---

## Agenda
1. Identitas, Hakikat, & Arah Pengembangan `clamav-service` (*Security-as-a-Service ala Stripe*)
2. Protokol Komunikasi (HTTP/REST Streaming + Async Webhook) & Kelompok Endpoint
3. Arsitektur Karantina (*Quarantine Vault*), Safe Storage Pattern, & Dual Mode Restore
4. Profil Beban VPS/RAM, Containerization Pattern (All-in-One Single Container)
5. User Interface (Embedded Web Admin SPA), Database Engine (SQLite WAL), & Config Layering
6. Edge Cases Handling (Zip-Bomb, File Size Limit, Password-Protected Archive Inspection)
7. Sistem Notifikasi & Alerting (Telegram, Discord, Email, Anti-Spam Flood Throttling)
8. Zero-Touch Master Key Generation & Auto-Write ke `.env`
9. Rate Limiting, Quota, & Granular API Key Policy (In-Memory Token Bucket)
10. Audit Logging, Retention Defaults (3 Hari Log / 7 Hari Karantina), & Streaming Export CSV/JSON
11. Manajemen Proses Container (Go Native Process Supervisor)
12. Standarisasi Format JSON Response & Error Codes Inventory

---

## Ringkasan Keputusan per Topik

### 1. Hakikat & Visi Sistem
- ✅ **Agentless Architecture**: Service bertindak sebagai *Unified Security Inspection Layer* independen tanpa butuh host agent atau sandbox VM berat di fase 1.
- ✅ **Security-as-a-Service**: Menyediakan integrasi instan *Sync Direct Scan* dan *Async Job + Webhook*.

### 2. Protokol & Endpoint
- ✅ **HTTP REST API**: Format multipart & raw binary stream via Unix Socket (`/tmp/clamd.sock`).
- ✅ **4 Kelompok Endpoint**: Core Scan, Async Webhook Jobs, Quarantine Vault Management, dan Health/Metrics/Ops.

### 3. Karantina & Restore
- ✅ **Built-in Vault**: Disimpan di volume terisolasi `/data/quarantine/` dengan nama netral (`.quarantine`), izin `0600`, dan scrambled/encrypted.
- ✅ **Dual Mode Restore**: Rilis file via Direct Admin Download atau Push balik ke S3 / Webhook.
- ✅ **Auto-Whitelist SHA256**: Otomatis mendaftarkan hash file saat di-restore agar tidak kena karantina berulang.

### 4. Sizing VPS & Containerization
- ✅ **All-in-One Container**: Mengemas API Go, `clamd`, `freshclam`, SQLite, dan Web UI dalam 1 image.
- ✅ **Sizing Hardware**: Minimum 2 GB RAM + 2 GB Swap (Direkomendasikan: 2 vCPU / 4 GB RAM). MinIO tidak dijalankan di VPS yang sama.

### 5. UI, Database, Config, & Enkripsi
- ✅ **Embedded Admin UI**: Single Page App ringan disajikan langsung dari binary Go.
- ✅ **SQLite (WAL Mode)**: Disimpan di `/data/clamav-service.db` (zero extra RAM overhead).
- ✅ **Hybrid Config**: File `.env` (bootstrap) + SQLite `system_settings` (live dynamic UI overrides).
- ✅ **AES-256-GCM Encryption**: Enkripsi field-level untuk data kredensial di SQLite, serta salted one-way hashing untuk API Keys.

### 6. Edge Cases & Password-Protected Files
- ✅ **Zip-Bomb Limits**: Maks 5 level rekursi, maks 1.000 file dalam arsip, maks 250 MB ekstraksi in-memory.
- ✅ **Encrypted Archives**: Return status `UNSCANNABLE (PASSWORD_PROTECTED)`. Mendukung input `archive_password` di API dan auto-try kamus password umum malware.

### 7. Alerting & Notifikasi
- ✅ **Multi-Channel**: Mendukung Telegram Bot, Discord Webhook, dan Email (SMTP).
- ✅ **Auto-Detect & UI Toggle**: Otomatis aktif jika kredensial di `.env` ada, serta dapat di-mute/toggle ON/OFF via Web UI.
- ✅ **Anti-Spam Flood Control**: Throttling lonjakan notifikasi dengan batch digest jika terdapat > 5 malware/menit.

### 8. Key Lifecycle & Inisialisasi
- ✅ **Zero-Touch Key Gen**: Sistem otomatis men-generate key dan langsung menuliskan `ENCRYPTION_KEY=...` ke file `.env` saat boot/init pertama. User menyalin key hanya untuk backup darurat/disaster recovery.

### 9. Rate Limiting & API Keys
- ✅ **2-Tier Policy**: Global server protection (`.env`) + Granular rate limit / IP whitelist / quota per API Key (Web UI).
- ✅ **In-Memory Token Bucket**: Performa tinggi (< 1 MB RAM) dengan HTTP 429 & standar header rate limit.

### 10. Audit Logging & Retensi
- ✅ **Retensi Log Audit (`scan_audit_logs`)**: **3 Hari** (Default).
- ✅ **Retensi File Karantina**: **7 Hari** (Default).
- ✅ **Export**: Streaming download CSV & JSON dengan filter tanggal, status, dan consumer.

### 11. Manajemen Proses Container
- ✅ **Go Native Supervisor**: Binary Go sebagai PID 1 yang men-spawn dan memonitor `clamd` dan `freshclam`, menangani graceful shutdown, dan auto-restart jika clamd crash.

### 12. Kontrak JSON Response
- ✅ Format balasan terstruktur untuk `CLEAN`, `INFECTED`, `SUSPICIOUS`, `UNSCANNABLE`, dan `ACCEPTED` (async).
- ✅ Inventaris kode error standar: `UNAUTHORIZED`, `FORBIDDEN`, `RATE_LIMIT_EXCEEDED`, `FILE_TOO_LARGE`, `ZIP_BOMB_DETECTED`, `SCAN_TIMEOUT`, `ENGINE_UNAVAILABLE`, `QUARANTINE_NOT_FOUND`.

---

## Action Items

| # | Item | Owner | Deadline | Status |
|---|------|-------|----------|--------|
| ACT-01 | Sinkronkan seluruh 12 pilar keputusan ke `.ai-doc/project-overview.md` | Agent / Melon 🏗️ | Today | 🟢 DONE |
| ACT-02 | Perbarui control plane `.ai-doc/3p.md` dan constitution | Agent / Sherin ✋ | Today | 🟢 DONE |
| ACT-03 | Push seluruh pembaruan artefak ke GitHub remote repo | Agent | Today | 🟢 DONE |
| ACT-04 | Transisi ke fase implementasi greenfield / pembuatan SCD & TDD | User & Team | Next Session | ⬜ PLANNED |

---
*Generated by AI Documentor — Brainstorming Add-On (Session Completed)*
