# Minutes of Meeting — Brainstorming Session

| Field | Value |
|-------|-------|
| **Tanggal** | 2026-09-02 |
| **Waktu** | 14:28 - 14:50 WIB |
| **Topik Utama** | Perancangan Arsitektur, Fitur, & Spesifikasi `clamav-service` |
| **Area** | Feature / Greenfield Architecture Design |
| **Peserta** | User, Sherin ✋ (Moderator), Melon 🏗️ (Technical Architect), Sultan ⚙️ (Senior Backend Engineer), Nindi 🔬 (Problem Solver), Lugi 📊 (Data Specialist), Bernadya ✨ (Creative Visionary) |
| **Teknik** | First Principles, Possibility Mapping, Comparative Trade-off, Lifecycle State Modeling |

---

## Agenda
1. Identitas, Hakikat, & Arah Pengembangan `clamav-service` (*What is it and where can it go?*)
2. Protokol Komunikasi (gRPC vs HTTP/REST) & Pola Integrasi Security-as-a-Service (ala Stripe)
3. Arsitektur Karantina (*Quarantine Vault*), Mekanisme Restore, & Whitelisting
4. Profil Beban VPS / RAM, Containerization Pattern (All-in-One vs Terpisah)
5. User Interface (Embedded Web Admin), Database Engine (SQLite), Konfigurasi, & Enkripsi Kredensial

---

## Ringkasan Diskusi

### 1. Hakikat & Visi `clamav-service`
**Jumlah ide:** 6

**Ide utama:**
- Unified Security Inspection Layer & Blast Radius Isolation untuk ekosistem aplikasi.
- High-Throughput Streaming Antivirus Gateway non-blocking.
- Security-as-a-Service yang pluggable (Sync Instant Verdict & Async Webhook/S3 presigned scan).
- Threat Intelligence & Audit Hub untuk pencatatan telemetry hash, signature virus, dan scan metrics.

**Keputusan:**
- ✅ Disepakati fokus pada model **Agentless File Inspection Service** (tidak butuh agent host EDR dan tidak perlu sandbox VM berat di fase awal).

**Next steps:**
- [ ] Tuangkan konsep ini ke dalam `project-overview.md`.

---

### 2. Protokol Komunikasi & Endpoint Design
**Jumlah ide:** 8

**Ide utama:**
- HTTP/REST dengan streaming socket bridge dipilih untuk adopsi universal dan kemudahan integrasi.
- Endpoint dibagi menjadi 4 kelompok: Core Scanning (File/Stream/URL), Async Jobs & Webhook, Quarantine Management, dan Health/Metrics/Ops.

**Keputusan:**
- ✅ Menggunakan protokol utama **HTTP REST API (Multipart + Binary Streaming)** dengan arsitektur modular yang siap ditambah gRPC adapter kelak.
- ✅ Mendukung endpoint Sync (Direct scan) dan Async (Job ID + Webhook callback).

---

### 3. Sistem Karantina & Restore
**Jumlah ide:** 5

**Ide utama:**
- Built-in Quarantine Vault menggunakan Local Volume isolasi (`/data/quarantine`) atau S3/MinIO eksternal secara pluggable.
- File terinfeksi dinetralkan: ekstensi diubah (`.quarantine`), permission `0600`, opsional XOR-scramble.
- Restore mendukung 2 mode: Direct Admin Download & S3 Push / Webhook notification.
- Fitur **Auto-Whitelist Hash (SHA256)** saat restore untuk mencegah *re-quarantine loop*.

**Keputusan:**
- ✅ Implementasi karantina built-in (tanpa vendor SaaS 3rd party).
- ✅ Mengaktifkan retensi TTL (auto-purge file karantina kadaluarsa) dan audit log aksi restore.

---

### 4. Containerization & Sizing VPS
**Jumlah ide:** 4

**Ide utama:**
- ClamAV memuat ~8.5 juta signature ke RAM: footprint idle `clamd` ~1.1–1.4 GB RAM, spike hingga ~2.0 GB saat freshclam database reload.
- Pola deployment All-in-One Single Container menghemat RAM dan menggunakan komunikasi Unix Domain Socket (`/tmp/clamd.sock`) yang berlatensi ultra-rendah.

**Keputusan:**
- ✅ Menggunakan **All-in-One Single Container** (`clamav-service` + `clamd` + `freshclam` + SQLite + Local Quarantine Volume).
- ✅ Rekomendasi VPS: Minimum 2 GB RAM + 2 GB Swap (Ideal: 2 vCPU / 4 GB RAM).
- ✅ MinIO tidak dijalankan di VPS yang sama untuk menghemat alokasi RAM.

---

### 5. UI, Database, Config, & Enkripsi Kredensial
**Jumlah ide:** 5

**Ide utama:**
- **UI:** Embedded Single Page Application (SPA) ringan yang disajikan dari backend binary (Dashboard status clamd, drag-and-drop test scan, quarantine explorer, API keys manager, audit logs).
- **Database:** **SQLite (WAL Mode)** di `/data/clamav-service.db` — zero extra RAM, performa ribuan ops/detik, single-file backup.
- **Config:** Hybrid 12-factor Environment Variables (`.env`) + Dynamic Database Overrides.
- **Enkripsi:** One-way salted hash untuk API Keys, dan **AES-256-GCM** (Field-level encryption) untuk data sensitif/kredensial di SQLite.

**Keputusan:**
- ✅ Arsitektur penyimpanan data menggunakan SQLite WAL mode.
- ✅ Embedded Web Admin UI disertakan langsung dalam container.
- ✅ Enkripsi AES-256-GCM wajib untuk semua kredensial yang tersimpan di DB.

---

## Action Items

| # | Item | Owner | Deadline | Status |
|---|------|-------|----------|--------|
| ACT-01 | Susun artefak `project-overview.md` mengintegrasikan seluruh konsensus brainstorming | Agent / Melon 🏗️ | Next Session | ⬜ |
| ACT-02 | Konfirmasi kesiapan transisi dari tahap Brainstorming ke Greenfield Planning Document | User / Sherin ✋ | After Break | ⬜ |

---

## Catatan Tambahan
Sesi dijeda sementara (break). Seluruh poin konsensus telah terdokumentasi dan siap diturunkan ke Project Overview & spesifikasi teknis setelah user kembali.

---
*Generated by AI Documentor — Brainstorming Add-On*
