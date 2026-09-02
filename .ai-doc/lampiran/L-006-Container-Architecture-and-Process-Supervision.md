# Lampiran L-006: Arsitektur Container & Supervisi Proses Go Native

| Metadata | Nilai |
|---|---|
| **Kode Lampiran** | `L-006` |
| **Nama Lampiran** | Arsitektur Container & Supervisi Proses Go Native |
| **Target Service** | `clamav-service` |
| **Status** | Approved |
| **Versi** | 1.0 |
| **Tanggal Pembuatan** | 2026-09-02 |
| **Dokumen Utama** | [`.ai-doc/project-overview.md`](file:///home/ubuntu/workspace/plan/clamav-service/.ai-doc/project-overview.md) |

---

## 1. Ringkasan
Dokumen ini mendefinisikan arsitektur container tunggal (*All-in-One Docker Image*), peran binary Go sebagai **Master Process (PID 1)** yang mengawasi daemon ClamAV (`clamd`) dan pembaruan berkala (`freshclam`), siklus urutan startup (*Boot Sequence*), mekanisme *Zero-Downtime Signature Reload*, serta profil memori RAM VPS.

---

## 2. Topologi All-in-One Container

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                 All-in-One Docker Container (Alpine Linux)                  │
│                                                                             │
│   ┌─────────────────────────────────────────────────────────────────────┐   │
│   │               Go API Binary (Master Process / PID 1)                │   │
│   │   ├── REST API Server (Port 8080) & Embedded Web UI (SPA)           │   │
│   │   ├── In-Memory Token Bucket Rate Limiter                           │   │
│   │   ├── SQLite Database Manager (WAL Mode)                            │   │
│   │   └── Process Supervisor (Lifecycle, OS Signal & Health Monitoring) │   │
│   └─────────┬─────────────────────────────────────────────┬─────────────┘   │
│             │                                             │                 │
│      Spawn & Supervise                             Spawn & Supervise        │
│             ▼                                             ▼                 │
│   ┌───────────────────┐                         ┌───────────────────┐       │
│   │ clamd Daemon      │ <───── Unix Socket ─────│ freshclam Daemon  │       │
│   │ (In-Memory Scans) │  (/tmp/clamd.sock)      │ (Update Signature)│       │
│   └───────────────────┘                         └───────────────────┘       │
│             │                                             │                 │
│             └──────────────────────┬──────────────────────┘                 │
│                                    ▼                                        │
│   Persistent Volume (/data/) ──> [clamav-service.db] + [/data/quarantine/]  │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 3. Siklus Urutan Booting (*Startup Boot Sequence*)

Untuk mencegah kegagalan socket connection saat container dinyalakan:

1. **Step 1: Security Init:**  
   Binary Go mengecek keberadaan `ENCRYPTION_KEY` di `.env` / volume. Jika belum ada, otomatis men-generate key dan menuliskannya ke `.env`.
2. **Step 2: Database Signature Verification:**  
   Mengecek apakah database virus di direktori `/var/lib/clamav/` sudah tersedia. Jika container baru (fresh install), menjalankan `freshclam` satu kali secara synchronous untuk mengunduh database awal.
3. **Step 3: Spawn `clamd` Daemon:**  
   Menjalankan subprocess `clamd --config-file=/etc/clamav/clamd.conf`.
4. **Step 4: Socket Readiness Polling:**  
   Binary Go melakukan polling pengecekan ketersediaan Unix Domain Socket `/tmp/clamd.sock` (mengirim perintah ping `PING` $\rightarrow$ menunggu `PONG`). Proses memuat signature ke RAM memakan waktu ~10–15 detik.
5. **Step 5: Spawn `freshclam` Daemon:**  
   Menjalankan background daemon `freshclam -d -c 12` untuk mengecek update signature virus setiap 2 jam.
6. **Step 6: Start HTTP Server & Web UI:**  
   Membuka listener HTTP di port `8080`. Endpoint `/api/v1/health` kini mengembalikan status `200 OK (Ready)`.

---

## 4. Penanganan Sinyal OS & Graceful Shutdown

* **PID 1 Signal Handler:**  
  Binary Go menangkap sinyal `SIGTERM` dan `SIGINT` dari Docker:
  1. Menutup HTTP listener agar tidak menerima request scan baru.
  2. Menyelesaikan request pemindaian yang sedang berjalan (grace period 15 detik).
  3. Mengirim sinyal shutdown ke subprocess `clamd` dan `freshclam`.
  4. Menutup koneksi database SQLite (`PRAGMA optimize` dan close WAL).
  5. Keluar dengan exit code `0` tanpa meninggalkan *zombie process*.

---

## 5. Mekanisme Zero-Downtime Signature Reload

Ketika `freshclam` berhasil mengunduh file update signature terbaru dari `database.clamav.net`:
1. `freshclam` mengirimkan perintah socket `RELOAD` ke daemon `clamd`.
2. `clamd` memuat database baru ke RAM secara paralel tanpa memutus koneksi socket aktif yang sedang melayani request scan.
3. Setelah database baru aktif, memori database lama dilepaskan.
4. Binary Go memperbarui versi signature pada audit log dan dashboard UI.

---

## 6. Profil Memori RAM & Rekomendasi Hardware

| Status Operasional | Penggunaan RAM (Perkiraan) | Catatan |
|---|---|---|
| **Engine `clamd` (Idle / Standby)** | **~1.1 GB – 1.4 GB** | Menampung ~8.5 juta signature virus di RAM. |
| **Saat Pemindaian Aktif** | **+100 MB – 300 MB** | Buffer stream in-memory dan parsing file. |
| **Saat `freshclam` Reload** | **Spike ~1.8 – 2.0 GB** | Selama ~5 detik saat swapping database baru di RAM. |
| **Go API Binary & SQLite** | **~15 MB – 30 MB** | Sangat ringan, in-process WAL mode. |
| **Total Baseline Container** | **~1.3 GB – 1.6 GB** | **Wajib VPS minimal 2 GB RAM + 2 GB Swap (Direkomendasikan: 2 vCPU / 4 GB RAM).** |
