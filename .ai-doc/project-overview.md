# Project Overview — `clamav-service`

> **Status:** Approved Greenfield Blueprint  
> **Created Date:** 2026-09-02  
> **Author:** Antigravity Full Team & User  
> **Reference:** [Final MoM Brainstorming](file:///home/ubuntu/workspace/plan/clamav-service/.ai-doc/brainstorming/mom-2026-09-02-clamav-service-architecture.md), [Discussion Summary](file:///home/ubuntu/workspace/plan/clamav-service/.ai-doc/brainstorming/discussion-clamav-service-core-architecture-2026-09-02.md)

---

## 1. Problem Statement

1. **Kompleksitas & Duplikasi Integrasi Malware Scanner:**  
   Banyak aplikasi backend dan microservices membutuhkan kemampuan pemindaian file dari ancaman virus/malware sebelum file disimpan ke storage. Mengintegrasikan ClamAV secara langsung (via process execution atau custom socket) di setiap service menciptakan duplikasi kode, kerumitan manajemen daemon, dan pemborosan resource.
2. **Risiko Blast Radius & Keamanan Runtime:**  
   Proses pemindaian file rentan terhadap serangan *denial-of-service* (seperti *zip-bomb*, broken executable, memory leak, atau eksploitasi engine). Tanpa pemisahan service yang terisolasi, kegagalan engine scanner dapat melumpuhkan seluruh aplikasi bisnis utama.
3. **Penanganan False Positive & File Terinfeksi yang Kaku:**  
   Banyak sistem file scanner langsung menghapus file yang terdeteksi virus tanpa opsi karantina yang aman. Hal ini berisiko menyebabkan hilangnya dokumen penting bisnis jika terjadi *false positive*.
4. **Ketiadaan Observabilitas, Alerting & Audit Terpadu:**  
   Tim keamanan (SecOps) dan DevOps kesulitan melacak riwayat pemindaian file, telemetri malware, serta membutuhkan notifikasi instan (Telegram/Discord/Email) saat ada file berbahaya yang masuk.

---

## 2. Target Users & Stakeholders

* **Consumer Applications & Microservices (Primary Users):**  
  Aplikasi web, upload gateway, worker background, atau service lain yang mengirimkan file/stream untuk dipindai secara instan (sync) atau asinkron (async job + webhook).
* **Security Administrators & SecOps:**  
  Tim keamanan yang memantau log audit pemindaian, mengevaluasi file di *Quarantine Vault*, melakukan rilis/whitelist file *false-positive*, mengonfigurasi notifikasi alert, dan mengelola API Keys.
* **DevOps & System Engineers:**  
  Tim infrastruktur yang mengoperasikan, memonitor resource (CPU/RAM/Socket), dan mengelola update signature virus secara otomatis.
* **End Users Aplikasi:**  
  Pengguna akhir yang terlindungi dari unggahan file berbahaya di seluruh ekosistem aplikasi.

---

## 3. Assumptions

1. Service berjalan sebagai **All-in-One Docker Container** mandiri berbasis Alpine Linux di VPS/Server dengan akses internet keluar untuk pembaruan database signature via `freshclam`.
2. Host VPS memiliki spesifikasi minimum **2 GB RAM + 2 GB Swap** (Direkomendasikan **2 vCPU / 4 GB RAM**) untuk menjamin daemon `clamd` berjalan stabil tanpa risiko *OOM Killer*.
3. Penyimpanan persisten default untuk database SQLite dan file karantina menggunakan **Docker Volume Lokal** (`/data`), dengan fleksibilitas adapter untuk S3/MinIO eksternal.
4. Komunikasi antara API wrapper dan daemon ClamAV menggunakan **Unix Domain Socket** lokal (`/tmp/clamd.sock`) untuk zero-latency in-memory streaming.
5. Runtime API wrapper menggunakan **Go (Golang)** yang bertindak sebagai **PID 1 Process Supervisor** untuk mengelola lifecycle `clamd` dan `freshclam`.

---

## 4. Goals & Objectives

* **High-Throughput Streaming Engine:** Menyediakan endpoint pemindaian file langsung (chunked streaming / in-memory) dengan latensi rendah (<100ms untuk file umum) tanpa disk I/O perantara.
* **Developer Experience yang Sangat Mudah (Security-as-a-Service ala Stripe):** Menyediakan API intuitif: *Sync Direct Scan* untuk verdict instan dan *Async Job + Webhook* untuk pemindaian file besar/cloud object.
* **Safe Quarantine Vault & Dual Restore:** Menyimpan file terinfeksi ke ruang terisolasi (ekstensi `.quarantine`, izin akses `0600`, scrambled/encrypted), dilengkapi fitur *Restore* (Direct Download / S3 Push) dan *Auto-Whitelist SHA256* untuk mencegah siklus karantina berulang.
* **Zero-Touch Automation & Bulletproof Encryption:** Otomatis men-generate dan meng-inject `ENCRYPTION_KEY` ke `.env` saat boot pertama, serta mengenkripsi data sensitif di SQLite menggunakan **AES-256-GCM**.
* **Proactive Multi-Channel Alerting:** Mengirimkan alert instan ke **Telegram, Discord, dan Email** saat malware terdeteksi, dilengkapi fitur *Anti-Spam Flood Throttling*.
* **Granular Traffic Control:** Menyediakan proteksi *2-Tier Rate Limiting* (global server protection + granular limit & IP whitelist per API Key) menggunakan in-memory token bucket.
* **Efficient Audit & Auto-Purge:** Menyediakan log audit komprehensif dengan retensi default hemat disk (**3 hari log audit, 7 hari karantina**) dan fitur streaming export CSV/JSON.

---

## 5. Scope

### In-Scope (Fase 1 / MVP)
1. **Core Scanning Engine:**
   - Pemindaian multipart upload file (`POST /api/v1/scan/file`).
   - Pemindaian raw binary chunked stream (`POST /api/v1/scan/stream`).
   - Pemindaian remote URL / Presigned S3 Link (`POST /api/v1/scan/url`).
2. **Asynchronous & Webhook Engine:**
   - Pemindaian asinkron dengan return `202 Accepted` + `job_id` (`POST /api/v1/scan/async`).
   - Polling status job (`GET /api/v1/scan/jobs/:job_id`).
   - Pengiriman hasil otomatis via HTTP Webhook Callback.
3. **Edge Case & Password-Protected Handling:**
   - Proteksi Zip-Bomb (Maks 5 level rekursi, maks 1.000 file, maks 250 MB ekstraksi).
   - Penanganan file terenkripsi password (`UNSCANNABLE / PASSWORD_PROTECTED`), dukungan parameter `archive_password`, dan auto-try kamus password umum.
   - Strict context timeout (30 detik).
4. **Built-in Quarantine Vault:**
   - Penyimpanan terisolasi di Local Volume (`/data/quarantine/`).
   - Manajemen REST API: List, Detail Metadata, Restore (Download / Push), dan Permanent Delete.
   - Retensi TTL otomatis (auto-purge file kadaluarsa, default **7 hari**).
   - Fitur Hash Whitelisting (SHA256) saat file di-restore.
5. **Multi-Channel Alerting & Notifikasi:**
   - Channel: Telegram Bot, Discord Webhook, Email (SMTP).
   - Mode aktivasi: Auto-detect kredensial + UI Toggle switch ON/OFF.
   - Proteksi anti-spam flood throttling (batch digest jika > 5 alert/menit).
6. **Embedded Web Admin Dashboard (SPA):**
   - Overview status daemon ClamAV, versi signature, rasio infected/clean, throughput grafik.
   - Interactive Drag-and-Drop Test Scan.
   - Quarantine Vault Explorer & Audit Log Viewer.
   - API Key Generator & Granular Policy Manager.
   - Alert Settings & Test Alert Trigger.
7. **Rate Limiting & API Keys:**
   - In-memory token bucket rate limiting.
   - 2-Tier: Global protection (`.env`) + Granular RPM/quota/IP whitelist per API Key (Web UI).
8. **Data, Audit & Retensi:**
   - SQLite (WAL Mode) di `/data/clamav-service.db`.
   - Retensi log audit scan default **3 hari** dengan pembersihan harian & auto-vacuum.
   - Fitur streaming export laporan audit ke **CSV & JSON**.
9. **Process Supervision & Lifecycle:**
   - Go Native Process Supervisor (PID 1) yang mengelola `clamd` dan `freshclam`, auto-restart, dan zero-downtime signature reload.
   - Zero-Touch Encryption Key auto-generation ke `.env`.
   - Health check probe (`GET /api/v1/health` dan `/healthz`) & Prometheus metrics (`GET /api/v1/metrics`).

### Out-of-Scope (Fase 1)
- Host/endpoint EDR agent yang dipasang di level sistem operasi laptop/server.
- Dynamic VM sandbox execution (seperti Cuckoo Sandbox / interactive detonation).
- Multi-tenant SaaS billing system.
- Distributed multi-node clustered clamd farm.

---

## 6. High-Level System Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    clamav-service (All-in-One Container)                    │
│                                                                             │
│   ┌───────────────────────────┐         ┌───────────────────────────────┐   │
│   │  Embedded Web Admin UI    │ <─────> │      HTTP REST API Gateway    │   │
│   │  (SPA Dashboard / Vault)  │         │ (Auth, Rate Limiter, Stream)  │   │
│   └───────────────────────────┘         └───────┬──────────────┬────────┘   │
│                                                 │              │            │
│                       ┌─────────────────────────┘              └────────┐   │
│                       ▼                                                 ▼   │
│       ┌───────────────────────────────┐                 ┌─────────────┐ │   │
│       │      Core Scanning Bridge     │ <──Unix Sock──> │ clamd       │ │   │
│       │ (Stream chunking / Clam Socket)                 │ (In-Memory  │ │   │
│       └───────────────┬───────────────┘                 │  Signatures)│ │   │
│                       │                                 └──────▲──────┘ │   │
│                       ▼                                        │        │   │
│       ┌───────────────────────────────┐                 ┌──────┴──────┐ │   │
│       │      Quarantine Manager       │                 │ freshclam   │ │   │
│       │ (Safe Vault / Dual Restore)   │                 │ (Auto-Update│ │   │
│       └───────────────┬───────────────┘                 └─────────────┘ │   │
│                       │                                                 │   │
│                       ▼                                                 │   │
│       ┌───────────────────────────────┐                                 │   │
│       │    SQLite Storage (WAL)       │ <── AES-256-GCM Encryption      │   │
│       │ (Audit Log, Keys, Settings)   │                                 │   │
│       └───────────────┬───────────────┘                                 │   │
│                       │                                                 │   │
│   Persistent Volume ──┼─────────────────────────────────────────────────┼───┘
│   (/data)             ▼                                                 ▼
│             [clamav-service.db]                              [/data/quarantine/]
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 7. Key Constraints

1. **Kebutuhan RAM ClamAV:** Daemon `clamd` memuat ~8.5 juta signature virus ke RAM (~1.1–1.4 GB baseline, spike ~2.0 GB saat reload). Host wajib memiliki kapasitas RAM minimal 2 GB + 2 GB swap.
2. **Neutralization Constraint:** File dalam karantina tidak boleh disimpan dalam format executable asli, wajib berizin `0600` dan diisolasi dari path web server.
3. **Single Port Exposure:** Seluruh layanan (REST API, Web Admin UI, Health, Metrics) disajikan melalui satu port HTTP utama (default `8080`).

---

## 8. Prerequisite

1. Docker Engine / Docker Compose terpasang pada environment target.
2. Port 8080 (atau yang ditentukan) dapat diakses oleh consumer service.
3. Akses keluar (outbound HTTPS) ke server database ClamAV (`database.clamav.net`) untuk download signature virus pertama kali.

---

## 9. Catatan Konsensus Brainstorming

* **Protokol:** Diputuskan menggunakan HTTP/REST dengan streaming socket sebagai prioritas utama demi kemudahan adopsi, dengan abstraksi core engine yang tetap modular.
* **Storage Karantina:** Menggunakan Local Persistent Volume sebagai default mengeliminasi kebutuhan dependensi MinIO/S3 terpisah pada VPS kecil.
* **Database & Retensi:** SQLite WAL mode dengan retensi default 3 hari log audit dan 7 hari karantina untuk menjaga footprint disk tetap minimal.
* **Supervisi Proses:** Menggunakan Go Native Supervisor (PID 1) tanpa supervisor eksternal (Supervisord) untuk menghemat RAM dan menjaga image tetap kecil.
* **Keamanan Kredensial:** Otomatisasi generate key pada first-boot langsung ke `.env`, enkripsi AES-256-GCM pada SQLite, dan salted hash pada API Keys.

---

## 10. Keputusan & Status

* [x] **Pilihan Bahasa Backend API:** **Go (Golang)** — binary tunggal, super hemat RAM (~15-30 MB), native goroutine streaming.
* [x] **Architecture Blueprint:** Approved & Locked.
* [ ] **TDD Decision:** Konfirmasi keputusan penerapan Test-Driven Development (TDD) untuk fase implementasi.
