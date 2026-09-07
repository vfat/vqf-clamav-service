# Web UI Walkthrough & Visual Guide — `clamav-service`

## 1. Overview Dashboard & Threat Monitoring Interface

Antarmuka web mandiri (*Embedded Single Page Application*) disajikan langsung dari binary Go pada port `8080` (rute `/static/index.html` dan `/`). Mengusung tema visual *Dark Glassmorphism* dengan tipografi modern (Outfit & Inter), indikator glowing status real-time, dan dropzone interaktif.

![ClamAV Service Web UI Dashboard](/home/ubuntu/workspace/plan/clamav-service/.ai-doc/assets/clamav-web-ui-dashboard.jpg)

---

## 2. Fitur-Fitur Antarmuka Web

### 2.1. System Telemetry & Gauges
- **Total Files Scanned**: Menampilkan jumlah file yang dipindai melalui API socket `zINSTREAM\0`.
- **Threats Neutralized**: Menghitung malware yang diisolasi ke Vault Karantina.
- **Scan Latency**: Mengukur latensi pemindaian streaming secara real-time (< 45 ms).
- **Engine Security**: Menampilkan status mode enkripsi simetris master key **AES-256-GCM**.

### 2.2. Interactive Live Scan Lab
- **Drag & Drop Dropzone**: Mendukung upload file executable, dokumen, arsip (`.zip`, `.tar.gz`, `.rar`) hingga 100 MB.
- **EICAR Standard Generator**: Tombol pintas untuk menguji deteksi instan sampel uji standar malware EICAR.
- **Radar Scanning Animation**: Indikator streaming byte biner ke daemon `clamd`.
- **Verdict Cards**:
  - 🟢 **CLEAN / WHITELISTED**: Badge hijau zamrud dengan hash SHA-256 dan durasi pemindaian.
  - 🔴 **MALWARE DETECTED**: Badge merah rubi menyala dengan nama virus terdeteksi dan ID Vault isolasi (`Q-YYYYMMDD-ULID`).

### 2.3. Quarantine Vault Explorer
- Daftar file terisolasi dalam mode scrambling biner dengan permission `0600`.
- Tombol **Restore** interaktif dengan opsi otomatis mendaftarkan hash SHA-256 ke whitelist untuk mencegah re-quarantine loop.

### 2.4. Custom YARA Threat Signatures
- **Tabel Aturan YARA**: Menampilkan daftar signature khusus (`.yara`) yang aktif pada engine ClamAV beserta tanggal deploy dan author.
- **Modal Deploy Rule Baru**: Form input Rule Identifier, Author/SOC Analyst, Deskripsi, dan sintaks YARA rule. Dilengkapi validasi sintaks context-aware (mengabaikan komentar & string literal) dan auto-reload socket `zRELOAD\0` zero-downtime.
- **Modal Inspeksi Rule**: Menampilkan isi definisi sintaks YARA lengkap dalam format code-block monospace.
- **Aksi Hapus**: Menghapus file `.yara` dari disk dan memuat ulang database ClamAV secara otomatis.

### 2.5. Asynchronous Scan Queue & Webhook Monitoring
- **Enqueue Async Form**: Mengunggah file besar atau berkas latar belakang tanpa memblokir requestor (mengembalikan HTTP `202 Accepted` bersama UUID `job_id`).
- **Webhook Callback URL**: Alamat webhook untuk menerima payload JSON notifikasi hasil pemindaian. Diproteksi sistem Anti-SSRF (menolak IP privat, loopback, dan cloud metadata).
- **Scan Job Queue Table**: Memantau siklus hidup pemrosesan berkas oleh worker pool secara live (`QUEUED`, `PROCESSING`, `COMPLETED`, `FAILED`), kalkulasi hash SHA-256, ukuran berkas, dan status pengiriman webhook dengan fitur auto-refresh setiap 5 detik.

### 2.6. Audit Logs & Streaming CSV Export
- Riwayat transaksi audit lengkap dengan filter verdict (`CLEAN` / `INFECTED`).
- Tombol **Export CSV** untuk mendownload riwayat log pemindaian secara streaming via `/api/v1/audit/export`.

### 2.7. Security Lock Screen & Password Protection
- Layar kunci *Security Lock Screen* dengan latar blur saat pertama kali mengakses dashboard Web Admin.
- Penyimpanan password terenkripsi salted hash SHA-256 pada tabel SQLite `system_settings`.
- Modal penggantian kata sandi administratif langsung dari antarmuka Web UI.

### 2.8. API Authentication & Token Copier
- Visualisasi kunci integrasi aktif `clam_live_...` dan contoh integrasi cURL instan untuk microservices pengirim.
