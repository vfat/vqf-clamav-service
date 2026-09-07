# SCD-YARA-Scanner (Single Component Design)

## 1. Context

### 1.1. Masalah yang Ingin Diselesaikan
Secara default, signature database ClamAV diperbarui secara berkala oleh daemon `freshclam` dari server resmi Cisco Talos. Namun pada skenario operasional enterprise dan insiden siber aktif (misalnya serangan ransomware tertarget atau ancaman zero-day), tim keamanan siber membutuhkan kemampuan untuk menyuntikkan (*inject*) aturan deteksi khusus (**Custom YARA Rules**) secara instan tanpa menunggu rilis global dari vendor antivirus.

ClamAV engine (`libclamav`) mendukung aturan YARA secara native: berkas aturan berekstensi `.yara` atau `.yar` yang diletakkan pada direktori basis data ClamAV akan otomatis dikompilasi dan dievaluasi saat proses scanning berlangsung.

Dibutuhkan komponen **YARA Rule Injection Engine** untuk:
1. Memvalidasi sintaks dan keamanan aturan YARA sebelum disimpan ke disk.
2. Menyimpan aturan ke direktori rules ClamAV dan mencatat metadatanya ke database.
3. Memicu reload database ClamAV secara *zero-downtime* via perintah Unix Domain Socket `RELOAD`.
4. Menyediakan API untuk manajemen aturan (tambah, lihat daftar, dan hapus).

### 1.2. Posisi Komponen dalam Sistem

```
┌─────────────────────────────────────────────────────────────┐
│  clamav-service                                             │
│                                                             │
│  [Security Admin / SOC]                                     │
│         │                                                   │
│         ▼ (POST /api/v1/rules/yara)                         │
│  ┌───────────────────────────────────────────────────────┐  │
│  │ YARA Rule Injection Engine (internal/yara)            │  │
│  │  1. Syntax Validator (rule, strings, condition)       │  │
│  │  2. Rule File Persistence (/var/lib/clamav/rules/)    │  │
│  │  3. SQLite Metadata Registry (table: yara_rules)      │  │
│  │  4. clamd Socket Reload Dispatcher                    │  │
│  └──────────────────────────┬────────────────────────────┘  │
│                             │                               │
│                             ▼ (RELOAD /var/run/clamd.ctl)   │
│                 ┌───────────────────────┐                   │
│                 │ ClamAV clamd Engine   │                   │
│                 │ (Loads Custom YARA)   │                   │
│                 └───────────────────────┘                   │
└─────────────────────────────────────────────────────────────┘
```

### 1.3. Hubungan dengan Komponen Lain
* **`internal/clamd`**: Menerima perintah socket `RELOAD` untuk memuat ulang signature tanpa menghentikan proses scanning yang sedang berjalan.
* **`internal/storage`**: Menyimpan metadata aturan YARA (ID, nama aturan, deskripsi, pembuat, waktu unggah) pada tabel `yara_rules`.
* **`internal/api`**: Menyediakan endpoint REST API untuk manajemen aturan YARA.

---

## 2. Scope

### 2.1. In-Scope
1. **Validasi Sintaks YARA:**
   * Memeriksa keberadaan keyword wajib YARA (`rule <name>`, `strings:`, `condition:`).
   * Memvalidasi nama rule agar hanya memuat karakter alfanumerik dan garis bawah (`^[a-zA-Z0-9_]+$`).
   * Menolak aturan kosong, aturan duplikat, atau sintaks yang tidak valid sebelum menyentuh filesystem.
2. **Penyimpanan Berkas Aturan (`.yara`):**
   * Menyimpan berkas aturan ke direktori terisolasi (`YARA_RULES_DIR`, default `/var/lib/clamav/rules/`).
   * Mengatur izin akses berkas (`0644`) agar dapat dibaca oleh daemon ClamAV.
3. **Database Registry:**
   * Tabel SQLite `yara_rules` untuk mengelola ID unik, nama rule, konten aturan mentah, identitas pembuat, dan timestamp.
4. **Zero-Downtime Hot Reload:**
   * Mengirimkan sinyal `RELOAD` ke ClamAV Unix socket client (`clamd.Client.Reload(ctx)`).
   * Menunggu konfirmasi reload berhasil (`RELOAD OK`).
5. **REST API Endpoint Manajemen:**
   * `POST /api/v1/rules/yara`: Mendaftarkan dan mengaktifkan aturan YARA baru.
   * `GET /api/v1/rules/yara`: Menampilkan daftar aturan YARA yang aktif.
   * `DELETE /api/v1/rules/yara/{id}`: Menghapus aturan dari database dan disk, lalu memicu hot reload.

### 2.2. Out-of-Scope
* Kompilasi bytecode YARA binary eksternal CGO (menghindari ketergantungan library C native tambahan; ClamAV mengevaluasi sintaks `.yara` teks secara langsung).
* Editor aturan berbasis GUI visual (manajemen dilakukan via API dan Web Admin UI dashboard).

---

## 3. Prerequisite

1. Daemon ClamAV mendukung perintah socket `RELOAD`.
2. Migrasi skema SQLite baru untuk tabel `yara_rules` pada `internal/storage`.
3. Direktori rules ClamAV (`/var/lib/clamav/rules`) memiliki permission read untuk user `clamav`.

---

## 4. Daftar Usecase

| Kode Usecase | Nama Usecase | Deskripsi Singkat |
|---|---|---|
| **`UC-YARA-01`** | Validasi & Injeksi Aturan YARA Baru | Menerima payload teks aturan YARA, memvalidasi sintaks, menyimpan berkas `.yara`, mencatat ke database, dan memicu reload socket ClamAV. |
| **`UC-YARA-02`** | Hot-Reload Signature Daemon ClamAV | Mengirimkan perintah `RELOAD` via Unix socket sehingga ClamAV langsung memuat aturan baru tanpa downtime. |
| **`UC-YARA-03`** | Tampilkan Daftar Aturan YARA Aktif | Menampilkan seluruh custom rules yang terdaftar beserta statistik deteksi dan status keaktifan. |
| **`UC-YARA-04`** | Hapus Aturan YARA & Sinkronisasi Daemon | Menghapus berkas aturan dari filesystem, menghapus dari database, dan memicu reload ClamAV agar rule tidak lagi berlaku. |

---

## 5. Rencana Skema Database (`yara_rules`)

```sql
CREATE TABLE IF NOT EXISTS yara_rules (
    id          TEXT PRIMARY KEY,
    rule_name   TEXT NOT NULL UNIQUE,
    description TEXT,
    content     TEXT NOT NULL,
    author      TEXT,
    is_active   INTEGER NOT NULL DEFAULT 1,
    created_at  DATETIME NOT NULL,
    updated_at  DATETIME NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_yara_active ON yara_rules(is_active);
```

---

## 6. Rencana Siklus TDD

1. **`TDD-013`**: Engine Rule Manager, Syntax Validator, & Socket `RELOAD` Controller (`internal/yara`).
2. **`TDD-014`**: REST API Handlers CRUD `/api/v1/rules/yara` & Integrasi Server (`internal/api`).
