# Data Dictionary: `scan_jobs`

## 1. Ringkasan Tabel

| Metadata | Keterangan |
|---|---|
| **Nama Tabel** | `scan_jobs` |
| **Deskripsi** | Menyimpan antrean dan histori status job pemindaian malware asinkron beserta hasil verdict dan metadata webhook |
| **Primary Key** | `id` (TEXT) |
| **Model Struct Go** | [`storage.ScanJob`](file:///home/ubuntu/workspace/plan/clamav-service/internal/storage/sqlite.go) |
| **Estimasi Volume** | Ratusan hingga ribuan job per hari |

---

## 2. Struktur Kolom & Tipe Data

| Nama Kolom | Tipe SQLite | Nullable | Default | Deskripsi Teknis |
|---|---|---|---|---|
| `id` | `TEXT` | **NO** | - | ID unik pemindaian asinkron (`job_<timestamp_nano>` atau UUID). |
| `file_name` | `TEXT` | **NO** | - | Nama asli berkas yang diunggah. |
| `file_size` | `INTEGER` | **NO** | - | Ukuran berkas dalam satuan bytes. |
| `file_sha256` | `TEXT` | YES | `NULL` | Hash SHA-256 berkas hasil komputasi saat pemindaian. |
| `callback_url`| `TEXT` | **NO** | - | URL tujuan webhook pengiriman hasil pemindaian. |
| `consumer`    | `TEXT` | YES | `NULL` | Identitas consumer/client yang mengajukan job. |
| `status`      | `TEXT` | **NO** | - | Status pemrosesan job: `QUEUED`, `PROCESSING`, `COMPLETED`, `FAILED`. |
| `verdict`     | `TEXT` | YES | `NULL` | Verdict pemindaian: `CLEAN`, `INFECTED`, `ERROR`. |
| `virus_name`  | `TEXT` | YES | `NULL` | Nama virus/signature jika terinfeksi malware. |
| `error_msg`   | `TEXT` | YES | `NULL` | Pesan error jika job gagal diproses. |
| `created_at`  | `DATETIME` | **NO** | - | Waktu penerimaan job (UTC RFC3339). |
| `updated_at`  | `DATETIME` | **NO** | - | Waktu pembaruan status job terakhir (UTC RFC3339). |

---

## 3. Batasan (*Constraints*) & Indeks

* **Primary Key:** `id`
* **Indeks:**
  * `idx_jobs_status`: Indeks pada kolom `status` untuk query worker dan filter pemantauan antrean.
