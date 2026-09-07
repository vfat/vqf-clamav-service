# DCD-03-YARA-Scanner

## 1. Metadata Komponen

| Field | Keterangan |
|---|---|
| **Komponen** | YARA Custom Signature Rule Injection Engine |
| **Kode Dokumen** | `DCD-03` |
| **Dasar Rancangan** | [`SCD-YARA-Scanner.md`](file:///home/ubuntu/workspace/plan/clamav-service/.ai-doc/plan/component/SCD-YARA-Scanner.md) |
| **Status Implementasi** | ✅ `Active / Production Ready` (`TDD-013` & `TDD-014` 100% GREEN) |
| **Target Audience** | Backend Engineers, Security Architects, SOC Analysts |

---

## 2. Object Identification

Berdasarkan implementasi arsitektur pada [`internal/yara/`](file:///home/ubuntu/workspace/plan/clamav-service/internal/yara/), [`internal/api/`](file:///home/ubuntu/workspace/plan/clamav-service/internal/api/), dan [`internal/storage/`](file:///home/ubuntu/workspace/plan/clamav-service/internal/storage/):

### Boundary
* **HTTP REST Ingress Router:** `http.ServeMux` pada [`internal/api/server.go`](file:///home/ubuntu/workspace/plan/clamav-service/internal/api/server.go) menangani route:
  - `POST /api/v1/rules/yara`
  - `GET /api/v1/rules/yara`
  - `DELETE /api/v1/rules/yara/{id}` dan `DELETE /api/v1/rules/yara`
* **Local Filesystem Rules Storage:** Direktori aturan berekstensi `.yara` (`YARA_RULES_DIR`, default `/var/lib/clamav/rules/`).
* **Unix Domain Socket Interface:** `net.Conn` Unix socket ClamAV daemon (`/var/run/clamav/clamd.ctl`) untuk transmisi perintah zero-downtime hot-reload `zRELOAD\x00`.

### Control
* **`api.handleYARAAdd`:** HTTP handler validasi input payload REST API, pemanggilan engine rule registration, dan penanganan HTTP response.
* **`api.handleYARAList`:** HTTP handler untuk query seluruh daftar aturan YARA yang aktif.
* **`api.handleYARADelete`:** HTTP handler untuk penghapusan aturan YARA berdasarkan ID.
* **`yara.ValidateRule`:** Syntax validator yang memvalidasi format nama rule (`^[a-zA-Z0-9_]{1,64}$`), keberadaan keyword deklarasi `rule <name>`, `condition:`, serta keseimbangan kurung kurawal `{}` dengan pengabaian string literal dan komentar.
* **`yara.Manager.AddRule`:** Orkestrator penambahan aturan: memvalidasi sintaks, menulis file aturan `.yara` berizin `0644`, menyisipkan metadata ke SQLite, dan memicu perintah `Reload` ClamAV secara atomik dengan mekanisme rollback.
* **`yara.Manager.DeleteRule`:** Controller pembersihan file aturan dari filesystem disk, penghapusan record dari database SQLite, dan penembakan sinyal `Reload` ke clamd.
* **`clamd.Client.Reload`:** Controller pengirim perintah ClamAV socket `zRELOAD\x00` untuk memuat ulang signature database secara instan tanpa downtime.

### Entity
* **Endpoint:**
  - `POST /api/v1/rules/yara`
  - `GET /api/v1/rules/yara`
  - `DELETE /api/v1/rules/yara/{id}`
* **DTO `yaraAddRequest`:** Struct JSON `{ "rule_name": string, "description": string, "content": string, "author": string }`.
* **Domain Model `storage.YARARule`:** Record database tabel `yara_rules` SQLite (`id`, `rule_name`, `description`, `content`, `author`, `is_active`, `created_at`, `updated_at`).
* **Database Table `yara_rules`:** Skema tabel SQLite dengan indeks `idx_yara_active`.

---

## 3. Use Case List

| No | Kode Use Case | Nama Use Case | Actor | Status | Referensi Detail |
|---|---|---|---|---|---|
| 1 | `UC-YARA-01` | Validasi & Injeksi Aturan YARA Baru | Security Admin / SOC | Active | Section 4.1 |
| 2 | `UC-YARA-02` | Hot-Reload Signature Daemon ClamAV | System / clamd | Active | Section 4.2 |
| 3 | `UC-YARA-03` | Tampilkan Daftar Aturan YARA Aktif | Security Admin / Dashboard | Active | Section 4.3 |
| 4 | `UC-YARA-04` | Hapus Aturan YARA & Sinkronisasi Daemon | Security Admin / SOC | Active | Section 4.4 |

---

## 4. Use Case Detail

### 4.1. UC-YARA-01: Validasi & Injeksi Aturan YARA Baru
* **Deskripsi:** Menerima payload aturan deteksi kustom YARA dari administrator keamanan, memvalidasi integritas sintaks, menulis berkas `.yara` ke disk sistem ClamAV, mendaftarkan metadata ke basis data SQLite, dan memicu reload socket clamd.
* **Actor:** Security Admin / SOC Analyst.
* **Precondition:** Layanan aktif, daemon clamd berjalan dan menerima koneksi Unix socket, direktori rules memiliki izin tulis.
* **Postcondition:** Berkas `${rulesDir}/${ruleName}.yara` tersimpan, tercatat di tabel `yara_rules`, clamd memuat database baru, dan response HTTP 201 Created dikembalikan.

### 4.2. UC-YARA-02: Hot-Reload Signature Daemon ClamAV
* **Deskripsi:** Mengirimkan sinyal biner `zRELOAD\x00` melalui Unix domain socket ke daemon clamd agar ClamAV membaca ulang seluruh signature dasar dan aturan `.yara` kustom secara zero-downtime tanpa merestart proses daemon.
* **Actor:** System (Autonomous Internal Controller).
* **Precondition:** Koneksi socket clamd tersedia.
* **Postcondition:** Daemon clamd merespons `RELOADING` dan memuat aturan baru tanpa memutus request pemindaian yang sedang berjalan. Jika gagal, operasi penambahan aturan di-rollback secara otomatis.

### 4.3. UC-YARA-03: Tampilkan Daftar Aturan YARA Aktif
* **Deskripsi:** Menampilkan seluruh daftar custom signature YARA yang terdaftar pada sistem beserta deskripsi, pembuat (*author*), status keaktifan, dan timestamp pendaftaran.
* **Actor:** Security Admin / Dashboard UI.
* **Precondition:** Database SQLite dapat diakses.
* **Postcondition:** Daftar array aturan YARA dikembalikan dalam format JSON (HTTP 200 OK).

### 4.4. UC-YARA-04: Hapus Aturan YARA & Sinkronisasi Daemon
* **Deskripsi:** Menghapus aturan YARA dari sistem berdasarkan ID aturan. Berkas `.yara` dihapus dari filesystem disk, record dihapus dari tabel `yara_rules`, dan daemon clamd dipicu `RELOAD` agar threat detection rule tersebut tidak lagi aktif.
* **Actor:** Security Admin / SOC Analyst.
* **Precondition:** ID aturan YARA terdaftar di database.
* **Postcondition:** Berkas fisik terhapus, data SQLite terhapus, clamd di-reload, dan response HTTP 200 OK dikembalikan.
