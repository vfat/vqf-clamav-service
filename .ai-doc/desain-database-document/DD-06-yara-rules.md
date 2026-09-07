# Data Dictionary: `yara_rules`

## 1. Ringkasan Tabel

| Metadata | Keterangan |
|---|---|
| **Nama Tabel** | `yara_rules` |
| **Deskripsi** | Menyimpan registrasi metadata custom threat detection rules (YARA) yang diinjeksikan ke engine ClamAV |
| **Primary Key** | `id` (TEXT) |
| **Model Struct Go** | [`storage.YARARule`](file:///home/ubuntu/workspace/plan/clamav-service/internal/storage/sqlite.go) |
| **Estimasi Volume** | Puluhan hingga ratusan aturan aktif per kluster keamanan |

---

## 2. Struktur Kolom & Tipe Data

| Nama Kolom | Tipe SQLite | Nullable | Default | Deskripsi Teknis |
|---|---|---|---|---|
| `id` | `TEXT` | **NO** | - | UUID v4 unik pengenal aturan YARA. |
| `rule_name` | `TEXT` | **NO** | - | Nama unik aturan YARA (alfanumerik & garis bawah, max 64 karakter). Contoh: `detect_webshell_php`. |
| `description` | `TEXT` | YES | `NULL` | Penjelasan deskripsi ancaman atau indikator kompromi (IoC). |
| `content` | `TEXT` | **NO** | - | Payload teks mentah isi aturan YARA lengkap (`rule ... { ... condition: ... }`). |
| `author` | `TEXT` | YES | `NULL` | Identitas analis SOC / consumer pembuat aturan. |
| `is_active` | `INTEGER` | **NO** | `1` | Status keaktifan aturan (1: aktif dimuat ke ClamAV, 0: nonaktif). |
| `created_at` | `DATETIME` | **NO** | - | Timestamp UTC pembuatan aturan (format RFC3339). |
| `updated_at` | `DATETIME` | **NO** | - | Timestamp UTC pembaruan aturan terakhir (format RFC3339). |

---

## 3. Batasan (*Constraints*) & Indeks

* **Primary Key:** `id`
* **Unique Constraint:** `rule_name UNIQUE` (Mencegah konflik nama berkas `.yara` pada disk).
* **Indeks:**
  * `idx_yara_active`: Indeks pada kolom `is_active` untuk memfilter aturan yang aktif dioperasikan.
