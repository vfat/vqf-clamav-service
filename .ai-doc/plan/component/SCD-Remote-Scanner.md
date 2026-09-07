# SCD-Remote-Scanner (Single Component Design)

## 1. Context

### 1.1. Masalah yang Ingin Diselesaikan
Saat ini pemindaian malware pada `clamav-service` hanya mendukung dua metode input lokal:
1. `multipart/form-data` upload langsung via `POST /api/v1/scan/file`.
2. Raw binary stream via `POST /api/v1/scan/stream`.

Pada arsitektur cloud modern (microservices, serverless, presigned S3/GCS uploads), client sering kali tidak mengunggah file langsung ke service antivirus, melainkan mengunggah ke object storage (AWS S3, Google Cloud Storage, MinIO). Layanan antivirus membutuhkan kemampuan untuk mengunduh dan memeriksa objek dari remote URL secara aman tanpa membebani storage lokal atau memicu risiko keamanan jaringan.

### 1.2. Posisi Komponen dalam Sistem
`Remote-Scanner` bertindak sebagai *Ingress Boundary & Streaming Fetcher* yang berada di antara HTTP API Server dan ClamAV UNIX Domain Socket:

```
[Consumer Application]
       │
       ▼ (POST /api/v1/scan/url)
┌──────────────────────────────────────────────────────────┐
│  clamav-service                                          │
│                                                          │
│  ┌────────────────────────────────────────────────────┐  │
│  │ Remote-Scanner Component                           │  │
│  │  1. URL Parser & Protocol Whitelist (http/https)   │  │
│  │  2. DNS Resolver & Anti-SSRF IP Sanitizer          │  │
│  │  3. Hardened HTTP Client (Redirect & Timeout Caps) │  │
│  │  4. Streaming Pipe (io.LimitReader -> ClamAV)      │  │
│  └────────────────────────────────────────────────────┘  │
│                           │                              │
│                           ▼ (zINSTREAM /var/run/clamd)   │
│                 ┌────────────────────┐                   │
│                 │ ClamAV clamd Engine│                   │
│                 └────────────────────┘                   │
└──────────────────────────────────────────────────────────┘
```

### 1.3. Hubungan dengan Komponen Lain
* **`api.Server`**: Menerima request `POST /api/v1/scan/url`, melakukan autentikasi (`AUTH_MODE`), dan meneruskan payload ke `Remote-Scanner`.
* **`clamd.Client`**: Menerima chunked stream dari hasil unduhan remote URL via Unix Domain Socket.
* **`quarantine.Vault`**: Menerima stream payload jika hasil verdict adalah `INFECTED`.
* **`storage.DB`**: Mencatat audit log hasil pemindaian remote URL.
* **`alert.Notifier`**: Mengirimkan alert jika file remote terinfeksi.

---

## 2. Scope

### 2.1. In-Scope
1. **Validasi Skema URL:** Hanya mengizinkan protokol `http://` dan `https://`.
2. **🛡️ Strict Anti-SSRF Defense (Network Boundary):**
   * Pengecekan alamat IP hasil resolusi DNS sebelum koneksi TCP dibuat (*pre-dial IP validation*).
   * Pemblokiran total terhadap:
     - IPv4 Loopback (`127.0.0.0/8`).
     - IPv6 Loopback (`::1`).
     - Private IP RFC1918 (`10.0.0.0/8`, `172.16.0.0/12`, `192.168.0.0/16`).
     - Link-Local & Cloud Metadata endpoint (`169.254.0.0/16`, termasuk `169.254.169.254`).
     - Carrier-Grade NAT (`100.64.0.0/10`).
   * Pencegahan *DNS Rebinding Attack* dengan mengunci koneksi hanya ke IP yang telah tervalidasi saat handshaking.
3. **Zero-Disk Streaming Pipe:** Mengalirkan bytes langsung dari response body remote HTTP ke Unix socket ClamAV tanpa menulis file perantara ke disk container.
4. **Alokasi Memori Terbatas (`io.LimitReader`):** Membatasi maksimal ukuran unduhan remote selaras dengan `MAX_SCAN_SIZE_MB`.
5. **Download Timeout & Redirect Ceiling:** Maksimal 30 detik total unduhan dan maksimal 3 HTTP redirects (setiap redirect URL divalidasi ulang terhadap aturan Anti-SSRF).

### 2.2. Out-of-Scope
* Protokol non-HTTP (misalnya `ftp://`, `gopher://`, `file://`, `sftp://`) ditolak keras pada tahap validasi URL.
* Autentikasi OAuth2 internal khusus provider cloud (user wajib menyediakan presigned URL yang sudah memiliki otorisasi akses unduhan mandiri).

---

## 3. Prerequisite

1. Helper enkripsi dan vault karantina `internal/quarantine` telah stabil dan berstatus 100% GREEN (`TDD-006B`).
2. Daemon `clamd` Unix Domain Socket `/var/run/clamav/clamd.ctl` siap menerima protokol chunked `zINSTREAM`.
3. Standar error code `SSRF_ATTEMPT_BLOCKED` dan `DOWNLOAD_FAILED` terdaftar pada kontrak JSON API.

---

## 4. Daftar Usecase

| Kode Usecase | Nama Usecase | Deskripsi Singkat |
|---|---|---|
| **`UC-URL-01`** | Scan File Remote URL Valid | Mengunduh file dari URL publik atau S3 presigned URL, memindai via ClamAV, dan mengembalikan verdict `CLEAN` atau `INFECTED`. |
| **`UC-URL-02`** | Blokir Serangan SSRF Internal | Menolak request pemindaian jika target URL mengarah ke IP localhost (`127.0.0.1`), metadata server cloud (`169.254.169.254`), atau intranet privat. |
| **`UC-URL-03`** | Cegah Eksploitasi File Raksasa / OOM | Membatasi streaming download remote menggunakan `io.LimitReader` hingga batas `MAX_SCAN_SIZE_MB` (default 100 MB) dan mengembalikan HTTP 413 jika melebihi batas. |
| **`UC-URL-04`** | Karantina Otomatis File Remote Terinfeksi | Jika file remote terdeteksi virus, payload dialirkan ke `quarantine.Vault` (AES-256-GCM) dan notifikasi dikirimkan ke webhook/Telegram. |

---

## 5. Catatan Diskusi & Keputusan Arsitektur

* **Keputusan Custom Transport (Sultan ⚙️ & Nindi 🔬):**  
  Alih-alih menggunakan `http.DefaultClient`, komponen ini wajib menginstansiasi `http.Client` dengan `net.Dialer.Control` kustom. Sebelum socket TCP `connect()` dijalankan oleh OS, callback `Control` memeriksa IP target. Jika IP masuk ke *disallowed CIDR*, dialer langsung me-return `errors.New("connection to private/link-local address forbidden")`.
* **DNS Rebinding Mitigation:** Menggunakan IP yang telah lolos validasi untuk dial langsung, bukan melakukan resolve ulang setelah pemeriksaan selesai.
* **Redirect Policy (Melon 🏗️):** Implementasi `CheckRedirect` pada `http.Client` untuk membatasi maksimal 3 hops dan melakukan validasi Anti-SSRF pada setiap hop tujuan.

---

## 6. Asumsi, Risiko, dan Mitigasi

| Kategori | Asumsi / Risiko | Mitigasi Desain |
|---|---|---|
| **Asumsi** | Target URL remote dapat diakses via port 80 (HTTP) atau 443 (HTTPS) dengan response status `200 OK`. | Kembalikan error `DOWNLOAD_FAILED` dengan HTTP status remote jika target mengembalikan 404, 403, atau 500. |
| **Risiko** | Slowloris download attack (remote server sengaja mengalirkan 1 byte per detik untuk menghabiskan worker). | Bungkus request dengan `context.WithTimeout(ctx, 30*time.Second)` dan per-read deadline. |
| **Risiko** | SSRF bypass melalui redirect ke URL internal (misal dari `https://evil.com` me-redirect ke `http://169.254.169.254`). | Handler `CheckRedirect` memvalidasi alamat tujuan setiap kali HTTP 301/302/307/308 terjadi. |

---
*Generated by AI Documentor — Component Planning Flow*
