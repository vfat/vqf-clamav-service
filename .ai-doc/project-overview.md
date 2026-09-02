# Project Overview — `clamav-service`

> **Status:** Draft (Greenfield Planning)  
> **Created Date:** 2026-09-02  
> **Author:** Antigravity Team & User  
> **Reference:** [MoM Brainstorming](file:///home/ubuntu/workspace/plan/clamav-service/.ai-doc/brainstorming/mom-2026-09-02-clamav-service-architecture.md), [Discussion Summary](file:///home/ubuntu/workspace/plan/clamav-service/.ai-doc/brainstorming/discussion-clamav-service-core-architecture-2026-09-02.md)

---

## 1. Problem Statement

1. **Kompleksitas Integrasi Malware Scanner:**  
   Banyak aplikasi backend dan microservices membutuhkan kemampuan pemindaian file dari ancaman virus/malware sebelum file disimpan ke storage. Namun, mengintegrasikan ClamAV secara langsung (via process execution atau custom socket) di setiap service menciptakan duplikasi kode, kerumitan manajemen daemon, dan pemborosan resource.
2. **Risiko Blast Radius & Keamanan Runtime:**  
   Proses pemindaian file berisiko mengalami serangan *denial-of-service* (seperti *zip-bomb* atau eksploitasi engine). Tanpa pemisahan service yang terisolasi, kegagalan engine scanner dapat melumpuhkan seluruh aplikasi bisnis utama.
3. **Penanganan False Positive & File Terinfeksi yang Kaku:**  
   Banyak sistem file scanner langsung menghapus file yang terdeteksi virus tanpa opsi karantina yang aman. Hal ini berisiko menyebabkan hilangnya dokumen penting bisnis jika terjadi *false positive*.

---

## 2. Target Users & Stakeholders

* **Consumer Applications & Microservices (Primary Users):**  
  Aplikasi web, upload gateway, worker background, atau service lain yang mengirimkan file/stream untuk dipindai secara instan atau asinkron.
* **Security Administrators & SecOps:**  
  Tim keamanan yang memantau log audit pemindaian, mengevaluasi file di *Quarantine Vault*, melakukan rilis/whitelist file *false-positive*, dan mengelola API Keys.
* **DevOps & System Engineers:**  
  Tim infrastruktur yang mengoperasikan, memonitor resource (CPU/RAM/Socket), dan mengelola update signature virus secara otomatis.
* **End Users Aplikasi:**  
  Pengguna akhir yang terlindungi dari unggahan file berbahaya di seluruh ekosistem aplikasi.

---

## 3. Assumptions

1. Service akan di-deploy menggunakan **All-in-One Docker Container** di VPS/Server dengan akses internet keluar untuk pembaruan database signature via `freshclam`.
2. Host VPS memiliki spesifikasi minimum **2 GB RAM + 2 GB Swap** (Direkomendasikan **4 GB RAM**) untuk menjamin daemon `clamd` berjalan stabil tanpa risiko *OOM Killer*.
3. Penyimpanan persisten default untuk database SQLite dan file karantina menggunakan **Docker Volume Lokal** (`/data`), dengan fleksibilitas adapter untuk S3/MinIO eksternal.
4. Komunikasi antara API wrapper dan daemon ClamAV menggunakan **Unix Domain Socket** lokal (`/tmp/clamd.sock` atau `/var/run/clamav/clamd.ctl`) untuk zero-latency in-memory streaming.

---

## 4. Goals & Objectives

* **High-Throughput Streaming Engine:** Menyediakan endpoint pemindaian file langsung (chunked streaming / in-memory) dengan latensi rendah (<100ms untuk file umum) tanpa overhead disk I/O sementara.
* **Developer Experience yang Sangat Mudah (Security-as-a-Service ala Stripe):** Menyediakan API intuitif: *Sync Direct Scan* untuk verdict instan dan *Async Job + Webhook* untuk pemindaian file besar/cloud object.
* **Safe Quarantine Vault & Dual Restore:** Menyimpan file terinfeksi ke ruang terisolasi (ekstensi `.quarantine`, izin akses `0600`, scrambled/encrypted), dilengkapi fitur *Restore* (Direct Download / S3 Push) dan *Auto-Whitelist SHA256* untuk mencegah siklus karantina berulang.
* **Zero-Overhead Embedded Architecture:** Mengemas API, ClamAV daemon, Freshclam, Embedded Web Admin UI, dan SQLite (WAL mode) dalam satu container tunggal yang hemat resource.
* **Auditability & Credential Security:** Menyediakan audit log komprehensif (SHA256, signature, latency) serta enkripsi field-level **AES-256-GCM** untuk kredensial sensitif.

---

## 5. Scope

### In-Scope (Fase 1 / MVP)
1. **Core Scanning Engine:**
   - Pemindaian multipart upload file (`POST /api/v1/scan/file`).
   - Pemindaian raw binary stream (`POST /api/v1/scan/stream`).
   - Pemindaian remote URL / Presigned S3 Link (`POST /api/v1/scan/url`).
2. **Asynchronous & Webhook Engine:**
   - Pemindaian asinkron dengan return `202 Accepted` + `job_id` (`POST /api/v1/scan/async`).
   - Polling status job (`GET /api/v1/scan/jobs/:job_id`).
   - Pengiriman hasil otomatis via HTTP Webhook Callback.
3. **Built-in Quarantine Vault:**
   - Penyimpanan terisolasi di Local Volume (`/data/quarantine`).
   - Manajemen REST API: List, Detail Metadata, Restore (Download / Push), dan Permanent Delete.
   - Retensi TTL otomatis (auto-purge file kadaluarsa, default 30 hari).
   - Fitur Hash Whitelisting (SHA256) saat file di-restore.
4. **Embedded Web Admin Dashboard (SPA):**
   - Overview status daemon ClamAV, versi signature, rasio infected/clean, throughput grafik.
   - Interactive Drag-and-Drop Test Scan.
   - Quarantine Explorer & Audit Log Viewer.
   - API Key Generator & Manager.
5. **Data & Config:**
   - SQLite (WAL Mode) di `/data/clamav-service.db`.
   - Hybrid Config (12-factor Environment Variables + Database overrides).
   - Enkripsi AES-256-GCM untuk data rahasia & salted hashing untuk API Keys.
6. **Operations & Health:**
   - Health check probe (`GET /api/v1/health` dan `/healthz`).
   - Prometheus metrics endpoint (`GET /api/v1/metrics`).
   - Daemon reload trigger (`POST /api/v1/engine/reload`).

### Out-of-Scope (Fase 1)
- Host/endpoint EDR agent yang dipasang di level sistem operasi laptop/server.
- Dynamic VM sandbox execution (seperti Cuckoo Sandbox / interactive detonation).
- Multi-tenant SaaS billing system.
- Distributed multi-node clustered clamd farm (cukup single all-in-one instance).

---

## 6. High-Level System Direction

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    clamav-service (All-in-One Container)                    │
│                                                                             │
│   ┌───────────────────────────┐         ┌───────────────────────────────┐   │
│   │  Embedded Web Admin UI    │ <─────> │      HTTP REST API Gateway    │   │
│   │  (SPA Dashboard / Tester) │         │ (Auth, Rate Limiter, Stream)  │   │
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
│       │ (Safe Vault / Neutralizer)    │                 │ (Auto-Update│ │   │
│       └───────────────┬───────────────┘                 └─────────────┘ │   │
│                       │                                                 │   │
│                       ▼                                                 │   │
│       ┌───────────────────────────────┐                                 │   │
│       │    SQLite Storage (WAL)       │                                 │   │
│       │ (Audit Log, Keys, Metadata)   │                                 │   │
│       └───────────────┬───────────────┘                                 │   │
│                       │                                                 │   │
│   Persistent Volume ──┼─────────────────────────────────────────────────┼───┘
│   (/data)             ▼                                                 ▼
│             [clamav-service.db]                              [/data/quarantine/]
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 7. Key Constraints

1. **Kebutuhan RAM ClamAV:** Daemon `clamd` memuat ~8.5 juta signature virus ke RAM (~1.1–1.4 GB baseline). Host wajib memiliki kapasitas RAM memadai + swap untuk mencegah crash.
2. **Neutralization Constraint:** File dalam karantina tidak boleh disimpan dalam format executable asli, wajib berizin `0600` dan diisolasi dari path web server.
3. **Single Port Exposure:** Seluruh layanan (REST API, Web Admin UI, Health, Metrics) disajikan melalui satu port HTTP utama (default `8080`).

---

## 8. Prerequisite

1. Docker Engine / Docker Compose terpasang pada environment target.
2. Port 8080 (atau yang ditentukan) dapat diakses oleh consumer service.
3. Akses keluar (outbound HTTPS) ke server database ClamAV (`database.clamav.net`) untuk download signature virus pertama kali.

---

## 9. Catatan Diskusi

* **Protokol:** Opsi gRPC dipertimbangkan, namun diputuskan menggunakan HTTP/REST dengan streaming socket sebagai prioritas utama demi kemudahan adopsi, dengan abstraksi core engine yang tetap modular untuk adapter gRPC kelak.
* **Storage Karantina:** Menggunakan Local Persistent Volume sebagai default mengeliminasi kebutuhan dependensi MinIO/S3 terpisah pada VPS kecil, namun tetap mempertahankan S3 driver untuk arsitektur multi-server.
* **Database:** SQLite WAL mode dipilih karena zero-RAM overhead, zero maintenance, dan performa tinggi untuk kebutuhan single instance.

---

## 10. Risiko, Asumsi, dan Hal yang Perlu Dikonfirmasi

* [ ] **Pilihan Bahasa Backend API:** Pilihan stack untuk API wrapper (misal: **Go** vs **Node.js/TypeScript** vs **Python**).  
  *Catatan:* Go sangat direkomendasikan karena konsumsi RAM sangat minim (~15-30 MB) dan menghasilkan binary mandiri yang sangat cocok untuk container all-in-one.
* [ ] **TDD Decision:** Konfirmasi keputusan penerapan Test-Driven Development (TDD) untuk fase implementasi.
