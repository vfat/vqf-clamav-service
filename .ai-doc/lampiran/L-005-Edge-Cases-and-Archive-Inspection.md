# Lampiran L-005: Penanganan Kasus Khusus, Proteksi Zip-Bomb, & Arsip Terenkripsi

| Metadata | Nilai |
|---|---|
| **Kode Lampiran** | `L-005` |
| **Nama Lampiran** | Penanganan Kasus Khusus, Proteksi Zip-Bomb, & Arsip Terenkripsi |
| **Target Service** | `clamav-service` |
| **Status** | Approved |
| **Versi** | 1.0 |
| **Tanggal Pembuatan** | 2026-09-02 |
| **Dokumen Utama** | [`.ai-doc/project-overview.md`](file:///home/ubuntu/workspace/plan/clamav-service/.ai-doc/project-overview.md) |

---

## 1. Ringkasan
Dokumen ini mendefinisikan batasan ambang batas (*limits*), algoritma proteksi dari serangan *decompression bomb* (Zip-Bomb), mekanisme penanganan file arsip ber-password (ala VirusTotal / Malware Labs), serta kebijakan penanganan batas waktu (*timeout*).

---

## 2. Parameter Batasan Ekstraksi & Proteksi Zip-Bomb

Untuk mencegah kehabisan memori (*OOM*) dan kelelahan CPU akibat dekompresi bertingkat:

| Parameter | Nilai Default | Konfigurasi Environment | Deskripsi |
|---|---|---|---|
| **Max Scan Size** | `100 MB` | `MAX_SCAN_SIZE_MB=100` | Batas maksimum ukuran file yang diterima gateway. |
| **Max Recursion Depth** | `5 Level` | `ARCHIVE_MAX_RECURSION=5` | Kedalaman nesting arsip di dalam arsip (zip in zip). |
| **Max Files in Archive** | `1.000 File` | `ARCHIVE_MAX_FILES=1000` | Jumlah maksimal file yang diizinkan diekstrak per arsip. |
| **Max Extraction Memory** | `250 MB` | `ARCHIVE_MAX_EXTRACT_MB=250`| Batas total ukuran uncompressed bytes di memori. |
| **Context Scan Timeout** | `30 Detik` | `SCAN_TIMEOUT_SECONDS=30` | Batas waktu eksekusi pemindaian sebelum koneksi di-abort. |

Jika ada arsip yang melanggar batasan di atas, pemindaian dihentikan seketika dan mengembalikan response error HTTP `400 Bad Request` dengan error code `ZIP_BOMB_DETECTED`.

---

## 3. Penanganan Arsip Terenkripsi Password (*Encrypted Archive Workflow*)

Arsip yang diproteksi password (seperti ZIP, RAR, 7z terenkripsi) adalah teknik klasik untuk menyembunyikan malware dari deteksi antivirus.

```
[Menerima File Arsip Terenkripsi]
               │
               ▼
   Parameter 'archive_password' ada?
   ├── Ya ──> Ekstrak in-memory dengan password tersebut ──> Scan Isi File
   └── Tidak ──────────────────────────────────────────────┐
                                                           ▼
                                      Coba Kamus Password Populer
                                      ('infected', 'password', '123456', dll)
                                                           │
                                   ┌───────────────────────┴───────────────────────┐
                                   ▼                                               ▼
                             Password Cocok                                Password Gagal
                             Ekstrak & Scan                                Kembalikan Status:
                                                                           UNSCANNABLE / SUSPICIOUS
```

### 3.1. Kebijakan Verdict: `UNSCANNABLE` (Bukan `CLEAN`)
Sistem **tidak pernah** mengembalikan status `CLEAN` untuk file ber-password yang tidak berhasil diekstrak.
* **Response Payload:**
  ```json
  {
    "success": true,
    "verdict": "UNSCANNABLE",
    "reason": "PASSWORD_PROTECTED_ARCHIVE",
    "is_safe": false,
    "risk_level": "SUSPICIOUS",
    "message": "Archive is encrypted. Provide 'archive_password' in request body to inspect content."
  }
  ```

### 3.2. Dukungan Parameter User-Supplied Password
API scan menerima field password opsional:
```http
POST /api/v1/scan/file
Content-Type: multipart/form-data

file: [binary payload]
archive_password: "mySecretPassword123"
```

---

## 4. Penanganan Timeout & Executable Rusak

1. **Timeout Handling:**
   * Setiap request scan diikat dengan `context.WithTimeout(ctx, 30*time.Second)`.
   * Jika timeout terjadi, socket ke `clamd` diputus secara elegan dan mengembalikan HTTP `504 Gateway Timeout` dengan error code `SCAN_TIMEOUT`.
2. **Broken Executables:**
   * File biner yang korup atau memiliki PE/ELF header tidak valid ditandai sebagai `verdict: SUSPICIOUS` dengan catatan `reason: BROKEN_EXECUTABLE`.
