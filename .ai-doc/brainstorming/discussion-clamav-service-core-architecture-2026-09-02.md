# Discussion Summary — Arsitektur & Spesifikasi Inti `clamav-service`

| Field | Value |
|-------|-------|
| **Sesi** | Perancangan Arsitektur, Fitur, & Spesifikasi `clamav-service` |
| **Area** | Feature / Greenfield Architecture Design |
| **Sub-Topik** | Core Scanning Engine, Quarantine System, Containerization, & Administration |
| **Tanggal** | 2026-09-02 |
| **Kategori** | Ideation, Solution Architecture, & Reliability Planning |
| **Teknik** | First Principles, Possibility Mapping, Comparative Trade-off, Lifecycle State Modeling |
| **Peserta** | User, Sherin ✋, Melon 🏗️, Sultan ⚙️, Nindi 🔬, Lugi 📊, Bernadya ✨ |
| **Jumlah Ide** | 28 ide terpetakan |

---

## Konteks
Diskusi awal greenfield planning untuk membangun `clamav-service` sebagai layanan antivirus/malware scanner yang modern, ringan, berkinerja tinggi, dan mudah diintegrasikan oleh berbagai aplikasi pengirim file. Diskusi berfokus pada penentuan identitas sistem, protokol komunikasi, mekanisme penanganan file terinfeksi (karantina & restore), efisiensi resource VPS, serta komponen pendukung (Web UI, SQLite DB, Config, dan Enkripsi).

---

## Ide yang Teridentifikasi

### 1. Unified Security Inspection Layer & Blast Radius Isolation
Memisahkan proses pemindaian malware dari backend aplikasi utama. Kerentanan scanning (seperti zip-bomb, memori overhead, crash clamd) terisolasi sepenuhnya di service ini.
> **Sumber:** Melon 🏗️ & Nindi 🔬 | **Status:** `approved` | **Kategori:** Architecture

### 2. High-Throughput Streaming Antivirus Gateway
Service menerima payload/stream dari HTTP body secara non-blocking dan meneruskannya via Unix Domain Socket langsung ke memory daemon `clamd` tanpa disk I/O perantara.
> **Sumber:** Sultan ⚙️ | **Status:** `approved` | **Kategori:** Backend Performance

### 3. Pluggable Security-as-a-Service (ala Stripe)
Pengalaman integrasi pengembang (DX) yang super sederhana: Sync Direct Scan (instant verdict), Async Cloud Scan (kirim URL/presigned link + webhook callback), dan SDK/API 1-baris.
> **Sumber:** Bernadya ✨ | **Status:** `approved` | **Kategori:** Product/DX

### 4. Agentless Threat Intelligence & Audit Hub
Pencatatan telemetri scan (SHA256 fingerprint, virus signature name, duration latency, status) tanpa perlu memasang agent di host/endpoint.
> **Sumber:** Lugi 📊 | **Status:** `approved` | **Kategori:** Data/Observability

### 5. Built-in Quarantine Vault dengan Safe Storage Pattern
File terinfeksi tidak langsung dihapus melainkan disimpan di storage terisolasi (`/data/quarantine` atau S3) dengan ekstensi `.quarantine`, permission `0600`, dan scrambled/encrypted.
> **Sumber:** Sultan ⚙️ & Nindi 🔬 | **Status:** `approved` | **Kategori:** Security

### 6. Dual Mode Restore & Anti Re-Quarantine Hash Whitelisting
Mekanisme rilis file false-positive: bisa via Direct Admin Download atau Push balik ke S3 / Webhook, disertai pendaftaran SHA256 ke daftar whitelist lokal agar tidak kena karantina ulang.
> **Sumber:** Nindi 🔬 & Melon 🏗️ | **Status:** `approved` | **Kategori:** Operational Workflow

### 7. All-in-One Single Container Deployment
Mengemas API wrapper, `clamd`, `freshclam`, embedded UI, dan SQLite dalam 1 Docker container dengan volume persisten `/data`. Menghemat RAM, latensi Unix socket 0 ms, dan setup 1-line run.
> **Sumber:** Melon 🏗️ & Sultan ⚙️ | **Status:** `approved` | **Kategori:** Infrastructure

### 8. In-Process SQLite (WAL Mode) Database
Penyimpanan metadata karantina, audit log, API keys, dan settings menggunakan SQLite. Zero extra RAM overhead, single file backup di volume `/data`.
> **Sumber:** Lugi 📊 | **Status:** `approved` | **Kategori:** Database

### 9. AES-256-GCM Field-Level Encryption & Key Hashing
Enkripsi simetris untuk kredensial sensitif di database menggunakan master key dari Environment Variable, dan one-way hashing untuk API Keys.
> **Sumber:** Nindi 🔬 | **Status:** `approved` | **Kategori:** Security

---

## Pembahasan & Keputusan

### Pembahasan 1: gRPC vs HTTP/REST
- **Analisis:** gRPC sangat cepat untuk internal binary RPC, tetapi HTTP/REST lebih universal untuk multipart upload dan mudah dikonsumsi oleh cURL, frontend, webhook, dan berbagai framework.
- **Keputusan:** Menggunakan **HTTP REST API** dengan streaming handler sebagai antarmuka utama fase 1, dirancang dengan core engine modular yang siap menerima gRPC adapter di masa depan.

### Pembahasan 2: Karantina & Ketergantungan Storage
- **Analisis:** Menggunakan vendor SaaS 3rd party atau mewajibkan MinIO terpisah akan membebani biaya atau RAM VPS.
- **Keputusan:** Menerapkan **Pluggable Storage Driver**. Default menggunakan **Local Volume Terisolasi** (zero extra dependency/RAM), dengan opsi konfigurasi environment variable ke S3/MinIO eksternal bila diperlukan.

### Pembahasan 3: Sizing VPS & Dampak Memori
- **Analisis:** ClamAV memuat ~8.5 juta signature ke RAM, membutuhkan footprint ~1.1–1.4 GB RAM idle dan spike hingga ~2.0 GB saat reload freshclam harian.
- **Keputusan:** VPS minimum 2 GB RAM + 2 GB Swap (Direkomendasikan: 2 vCPU / 4 GB RAM). Hindari menjalankan server MinIO di VPS yang sama.

### Pembahasan 4: Embedded Admin Dashboard & Autentikasi
- **Analisis:** Antarmuka web mandiri dibutuhkan untuk monitoring status daemon, uji coba scan drag-and-drop, manajemen karantina, dan API keys.
- **Keputusan:** Menyertakan **Embedded SPA Admin UI** yang disajikan langsung oleh backend binary, diproteksi Session/JWT Cookie dengan Rate Limiter.

---

## Keputusan Final

- ✅ **Protokol:** HTTP/REST API (Sync, Chunked Stream, Async Webhook).
- ✅ **Deployment:** All-in-One Single Docker Container (ClamAV daemon + Freshclam + API + Web UI + SQLite).
- ✅ **Karantina:** Built-in Quarantine Vault di Local Volume `/data/quarantine/` dengan TTL auto-purge dan fitur Restore + SHA256 Whitelisting.
- ✅ **Database:** SQLite (WAL Mode) di `/data/clamav-service.db`.
- ✅ **Keamanan:** AES-256-GCM field encryption, salted hash API keys, role-based access.

---

## Next Steps

- [ ] **Penyusunan Project Overview**: Turunkan hasil kesepakatan ini ke dokumen resmi `.ai-doc/project-overview.md`.
- [ ] **Review Arsitektur**: Finalisasi use case dan acceptance criteria bersama user setelah sesi break.

---

## Referensi & Context
- Skill: [ai-documentor](file:///home/ubuntu/workspace/plan/clamav-service/.agent/skills/ai-documentor/SKILL.md)
- Minutes of Meeting: [mom-2026-09-02-clamav-service-architecture.md](file:///home/ubuntu/workspace/plan/clamav-service/.ai-doc/brainstorming/mom-2026-09-02-clamav-service-architecture.md)
- Control Plane: [.ai-doc/3p.md](file:///home/ubuntu/workspace/plan/clamav-service/.ai-doc/3p.md)

---
*Generated by AI Documentor — Brainstorming Add-On*
