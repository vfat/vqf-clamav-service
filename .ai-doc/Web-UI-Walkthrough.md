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

### 2.4. Audit Logs & Streaming CSV Export
- Riwayat transaksi audit lengkap dengan filter verdict.
- Tombol **Export CSV** untuk mendownload ribuan log secara streaming via `/api/v1/audit/export`.

### 2.5. API Authentication & Token Copier
- Visualisasi kunci integrasi aktif `clam_live_...` dan contoh integrasi cURL instan untuk microservices pengirim.
