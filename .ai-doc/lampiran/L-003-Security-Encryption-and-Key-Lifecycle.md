# Lampiran L-003: Keamanan, Enkripsi, & Manajemen Siklus Kunci

| Metadata | Nilai |
|---|---|
| **Kode Lampiran** | `L-003` |
| **Nama Lampiran** | Keamanan, Enkripsi, & Manajemen Siklus Kunci |
| **Target Service** | `clamav-service` |
| **Status** | Approved |
| **Versi** | 1.0 |
| **Tanggal Pembuatan** | 2026-09-02 |
| **Dokumen Utama** | [`.ai-doc/project-overview.md`](file:///home/ubuntu/workspace/plan/clamav-service/.ai-doc/project-overview.md) |

---

## 1. Ringkasan
Dokumen ini mendefinisikan standar kriptografi, mekanisme inisialisasi otomatis kunci master (*Zero-Touch Key Generation*), enkripsi field-level pada database SQLite, pengelolaan salted hash untuk API Keys, serta prosedur pencadangan darurat (*Disaster Recovery*).

---

## 2. Inisialisasi Otomatis Master Key (*Zero-Touch Auto-Generation*)

Untuk memberikan pengalaman pengguna yang mulus tanpa mengorbankan keamanan, sistem menerapkan mekanisme otomatisasi 100%:

```
[Container Boot / CLI init]
             │
             ▼
     Cek file .env di host?
             │
     ┌───────┴────────────────────────┐
     ▼                                ▼
[File .env Belum Ada / Kosong]   [File .env Sudah Ada]
1. Generate AES-256 Key kuat     Gunakan ENCRYPTION_KEY
2. Buat file .env otomatis       yang sudah terpasang
3. Tulis ENCRYPTION_KEY=...
             │
             ▼
[Tampilkan Banner Terminal & UI]
(Sebagai catatan backup darurat user)
```

### 2.1. Banner Terminal Saat Inisialisasi Pertama
```text
╔══════════════════════════════════════════════════════════════════════════╗
║                     🔑 CLAMAV-SERVICE MASTER KEY                        ║
╠══════════════════════════════════════════════════════════════════════════╣
║                                                                          ║
║  clam_sec_89fa72bc19401e9a7c3b2f8109d4e56a78bc34a8e23f091c784...         ║
║                                                                          ║
╠══════════════════════════════════════════════════════════════════════════╣
║  ✅ Berhasil digenerate dan otomatis ditulis ke file .env                ║
║  ⚠️ Salin dan simpan key ini untuk keperluan Disaster Recovery / Backup.  ║
╚══════════════════════════════════════════════════════════════════════════╝
```

---

## 3. Spesifikasi Enkripsi Database (Field-Level AES-256-GCM)

Seluruh kredensial sensitif di dalam database SQLite (`clamav-service.db`) dienkripsi pada level kolom sebelum disimpan ke disk:

* **Algoritma:** `AES-256-GCM` (Galois/Counter Mode) dengan authenticated encryption.
* **Panjang Kunci:** 256-bit (32 bytes).
* **Nonce/IV:** 96-bit (12 bytes) unik yang di-generate secara acak (`crypto/rand`) untuk setiap operasi enkripsi.
* **Payload Format:** `base64( nonce [12B] + ciphertext + auth_tag [16B] )`.
* **Field yang Wajib Dienkripsi:**
  1. `alert_configs.bot_token` (Token Bot Telegram)
  2. `alert_configs.webhook_url` (URL Webhook Discord)
  3. `alert_configs.smtp_password` (Password SMTP Email)
  4. `storage_configs.s3_secret_key` (Kredensial S3 Eksternal)

---

## 4. Keamanan API Keys (Salted One-Way Hashing)

* **Format Token Publik:** `clam_live_{32_karakter_acak}` (Contoh: `clam_live_8fbc2190ad4e...`).
* **Prinsip Tampilan Tunggal:** Token plaintext hanya ditampilkan **satu kali** di Web Admin UI saat pertama kali dibuat.
* **Penyimpanan di Database:**
  * Disimpan dalam bentuk **Argon2id** atau **SHA-256 salted hash** di tabel `api_keys`.
  * Sistem tidak pernah menyimpan API Key dalam bentuk plaintext.
* **Verifikasi Request:**
  1. Service mengekstrak token dari header `Authorization: Bearer clam_live_...`.
  2. Menghitung hash token dan mencocokkannya dengan database.

---

## 5. Prosedur Disaster Recovery & Migrasi Server

Jika server VPS mengalami kerusakan dan database SQLite dipindahkan ke host baru:
1. Pindahkan folder `/data/` ke server baru.
2. Salin nilai `ENCRYPTION_KEY` lama ke file `.env` di server baru.
3. Jalankan container di server baru. Seluruh data kredensial dan konfigurasi akan otomatis terbaca dan terdekripsi secara normal.
