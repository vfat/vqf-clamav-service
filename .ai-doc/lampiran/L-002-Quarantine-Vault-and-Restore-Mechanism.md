# Lampiran L-002: Arsitektur Quarantine Vault & Mekanisme File Restore

| Metadata | Nilai |
|---|---|
| **Kode Lampiran** | `L-002` |
| **Nama Lampiran** | Arsitektur Quarantine Vault & Mekanisme File Restore |
| **Target Service** | `clamav-service` |
| **Status** | Approved |
| **Versi** | 1.0 |
| **Tanggal Pembuatan** | 2026-09-02 |
| **Dokumen Utama** | [`.ai-doc/project-overview.md`](file:///home/ubuntu/workspace/plan/clamav-service/.ai-doc/project-overview.md) |

---

## 1. Ringkasan
Dokumen ini merinci arsitektur ruang penyimpanan terisolasi (*Quarantine Vault*), format netralisasi file terinfeksi, mekanisme retensi kadaluarsa otomatis, serta dua model pemulihan (*restore*) file jika terjadi kasus *false-positive*, termasuk pencegahan siklus karantina berulang (*Anti Re-Quarantine Hash Whitelisting*).

---

## 2. Struktur Penyimpanan Terisolasi (*Quarantine Vault*)

### 2.1. Lokasi & Izin Akses File
* **Direktori Default:** `/data/quarantine/` pada volume persisten lokal.
* **Hak Akses File (Permissions):** `0600` (`-rw-------`), hanya proses user service yang memiliki hak baca/tulis. Dilarang memiliki izin eksekusi (`no-exec`).
* **Format Nama File (Safe Neutralization):**
  File yang terinfeksi **tidak boleh disimpan dengan nama atau ekstensi aslinya** untuk mencegah eksekusi tidak sengaja oleh OS atau tool lain.
  * *Pola Penamaan:* `Q-{YYYYMMDD}-{ULID}.quarantine`
  * *Contoh:* `Q-20260902-01J7XYZ894K2019.quarantine`
* **Scrambling / Enkripsi:** Konten biner di-scramble (XOR mask / AES-GCM ringan) sebelum ditulis ke disk.

---

## 3. Metadata Karantina di SQLite

Setiap file yang masuk ke karantina dicatat dalam tabel `quarantine_records`:

```sql
CREATE TABLE quarantine_records (
    id                  TEXT PRIMARY KEY,          -- Format: Q-YYYYMMDD-ULID
    original_filename   TEXT NOT NULL,             -- Nama asli: misal 'invoice.pdf.exe'
    file_size_bytes     INTEGER NOT NULL,          -- Ukuran file asli
    file_sha256         TEXT NOT NULL,             -- SHA256 checksum asli
    virus_name          TEXT NOT NULL,             -- Nama malware yang terdeteksi
    source_consumer     TEXT,                      -- ID API Key atau Consumer Name
    stored_path         TEXT NOT NULL,             -- Path di /data/quarantine/
    status              TEXT NOT NULL,             -- 'QUARANTINED', 'RESTORED', 'DELETED'
    created_at          DATETIME NOT NULL,         -- Waktu karantina
    expires_at          DATETIME NOT NULL,         -- Waktu kadaluarsa TTL (default 7 hari)
    restored_at         DATETIME,                  -- Waktu di-restore (nullable)
    restored_by         TEXT,                      -- Akun/Admin yang me-restore
    restore_reason      TEXT                       -- Alasan/justifikasi restore
);
```

---

## 4. Siklus Hidup File & Retensi Otomatis (Default 7 Hari)

1. **Auto-Purge Background Worker:**  
   Service menjalankan cron cleaner internal setiap 24 jam.
2. **Kriteria Pembersihan:**
   ```sql
   SELECT * FROM quarantine_records 
   WHERE status = 'QUARANTINED' AND expires_at < datetime('now');
   ```
3. **Aksi:**  
   File biner `.quarantine` terkait dihapus permanen dari disk, status record di database diubah menjadi `DELETED` (metadata audit tetap dipertahankan untuk histori pelacakan).

---

## 5. Mekanisme Pemulihan (*File Restore Workflow*)

Jika file dinyatakan sebagai *false-positive* oleh tim SecOps, file dapat dipulihkan melalui 2 mode:

```
[File Terkarantina (Status: QUARANTINED)]
                 │
                 ├── 1. Verifikasi Checksum SHA256 Integritas
                 ├── 2. Auto-Daftarkan SHA256 ke Whitelist (Opsional)
                 │
         ┌───────┴────────────────────────┐
         ▼                                ▼
[Mode A: Direct Download]        [Mode B: S3 / Webhook Push]
De-scramble stream langsung       Upload kembali ke S3 target &
ke browser Security Admin        tembak webhook event ke consumer
         │                                │
         └───────────────┬────────────────┘
                         ▼
             Ubah Status: RESTORED
             Catat Audit (restored_by, reason)
```

### 5.1. Mode A: Direct Admin Download (Manual Recovery)
* **Endpoint:** `POST /api/v1/quarantine/:id/restore?mode=download`
* **Cara Kerja:** Service membaca file `.quarantine`, melakukan de-scramble, dan mengembalikan file biner asli via HTTP stream dengan header `Content-Disposition: attachment; filename="invoice.pdf"`.

### 5.2. Mode B: Target Callback / S3 Push (Automated Re-Dispatch)
* **Endpoint:** `POST /api/v1/quarantine/:id/restore`
* **Request Body:**
  ```json
  {
    "mode": "dispatch",
    "target_s3_url": "s3://app-storage/uploads/2026/invoice.pdf",
    "notify_webhook_url": "https://app.com/api/quarantine-callback",
    "restored_by": "security-officer@internal.corp",
    "reason": "False-positive verified by vendor SecOps",
    "auto_whitelist": true
  }
  ```
* **Cara Kerja:** Service mengunggah file yang telah dipulihkan ke target S3 dan mengirimkan webhook event `quarantine.file_restored`.

---

## 6. Fitur Anti Re-Quarantine Hash Whitelisting

* **Masalah:** File yang sudah di-restore sering kali di-scan ulang oleh cron background atau consumer lain, menyebabkan file masuk karantina berulang kali.
* **Solusi:** Saat parameter `auto_whitelist: true` disertakan dalam aksi restore, hash SHA256 file tersebut otomatis disimpan ke tabel `whitelist_signatures`:
  ```sql
  INSERT INTO whitelist_signatures (sha256_hash, description, added_by, created_at)
  VALUES ('275a021b...', 'Whitelisted false-positive invoice', 'security-officer', datetime('now'));
  ```
* **Efek:** Pada pemindaian berikutnya, jika hash file cocok dengan tabel whitelist, service langsung mengembalikan status `verdict: CLEAN` dengan catatan `whitelisted: true`.
