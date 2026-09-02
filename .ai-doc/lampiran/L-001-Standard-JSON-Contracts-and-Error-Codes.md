# Lampiran L-001: Standarisasi Kontrak JSON Response & Daftar Kode Error

| Metadata | Nilai |
|---|---|
| **Kode Lampiran** | `L-001` |
| **Nama Lampiran** | Standarisasi Kontrak JSON Response & Daftar Kode Error |
| **Target Service** | `clamav-service` |
| **Status** | Approved |
| **Versi** | 1.0 |
| **Tanggal Pembuatan** | 2026-09-02 |
| **Dokumen Utama** | [`.ai-doc/project-overview.md`](file:///home/ubuntu/workspace/plan/clamav-service/.ai-doc/project-overview.md) |

---

## 1. Ringkasan
Dokumen ini mendefinisikan format baku payload response JSON untuk semua endpoint pemindaian, operasional, dan penanganan error pada `clamav-service`. Semua consumer service wajib mengacu pada spesifikasi ini untuk melakukan parsing hasil pemindaian dan penanganan error.

---

## 2. Header Standar HTTP

Setiap response dari `clamav-service` menyertakan header standar berikut:
* `Content-Type: application/json; charset=utf-8`
* `X-Request-ID: req_01j7xyz984k...` (Tracking ID unik per request)
* `X-RateLimit-Limit: 100` (Batas kuota request)
* `X-RateLimit-Remaining: 92` (Sisa kuota request)
* `X-RateLimit-Reset: 1725268900` (Epoch timestamp reset kuota)

---

## 3. Format Response Sukses

### 3.1. File Bersih (`verdict: CLEAN`)
Dikembalikan saat file tidak mengandung malware/virus apa pun.
* **HTTP Status:** `200 OK`
```json
{
  "success": true,
  "verdict": "CLEAN",
  "data": {
    "file_name": "document_contract.pdf",
    "file_size": 1048576,
    "file_sha256": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
    "mime_type": "application/pdf",
    "scan_duration_ms": 42,
    "signature_version": "ClamAV 1.4.0 / Daily:27400",
    "scanned_at": "2026-09-02T17:15:00Z"
  }
}
```

### 3.2. File Terinfeksi Malware (`verdict: INFECTED`)
Dikembalikan saat file terdeteksi mengandung signature virus/malware dan telah otomatis dipindahkan ke Quarantine Vault.
* **HTTP Status:** `200 OK`
```json
{
  "success": true,
  "verdict": "INFECTED",
  "threat": {
    "virus_name": "Win.Trojan.Agent-9812",
    "severity": "HIGH",
    "action_taken": "QUARANTINED",
    "quarantine_id": "Q-2026-9812",
    "quarantine_expires_at": "2026-09-09T17:15:00Z"
  },
  "data": {
    "file_name": "invoice_payment.pdf.exe",
    "file_size": 2415082,
    "file_sha256": "275a021bbfb6489e54d471899f7db9d1663fc695ec2fe2a2c4538aabf651fd0f",
    "scan_duration_ms": 68,
    "scanned_at": "2026-09-02T17:15:00Z"
  }
}
```

### 3.3. File Mencurigakan / Terenkripsi (`verdict: SUSPICIOUS` / `UNSCANNABLE`)
Dikembalikan saat isi arsip tidak dapat dipindai karena terkunci password atau rasio kompresi mencurigakan.
* **HTTP Status:** `200 OK`
```json
{
  "success": true,
  "verdict": "UNSCANNABLE",
  "reason": "PASSWORD_PROTECTED_ARCHIVE",
  "is_safe": false,
  "risk_level": "SUSPICIOUS",
  "message": "Archive is encrypted. Provide 'archive_password' in request body to inspect content.",
  "data": {
    "file_name": "confidential_payroll.zip",
    "file_size": 512000,
    "file_sha256": "4b227777d4dd1fc61c6f884f48641d02b4d121d3fd328cb08b5531fcacdabf8a",
    "scan_duration_ms": 12,
    "scanned_at": "2026-09-02T17:15:00Z"
  }
}
```

### 3.4. Async Scan Job Diterima (`status: ACCEPTED`)
Dikembalikan saat consumer mengirim job pemindaian asinkron (`POST /api/v1/scan/async`).
* **HTTP Status:** `202 Accepted`
```json
{
  "success": true,
  "status": "ACCEPTED",
  "job_id": "job_01j7xyz894k2019",
  "message": "Scan job submitted to background queue. Result will be dispatched to webhook_url.",
  "check_status_url": "/api/v1/scan/jobs/job_01j7xyz894k2019"
}
```

---

## 4. Format Standard Error Response

Semua respon error menggunakan format terstruktur seragam dengan field `code`, `message`, dan opsional `details`:

```json
{
  "success": false,
  "error": {
    "code": "FILE_TOO_LARGE",
    "message": "Uploaded file exceeds maximum allowed scan size of 100 MB.",
    "details": {
      "max_allowed_bytes": 104857600,
      "received_bytes": 157286400
    }
  }
}
```

---

## 5. Daftar Kode Error Standar (*Error Code Inventory*)

| Error Code | HTTP Status | Deskripsi Penyebab | Penanganan di Consumer |
|---|---|---|---|
| `UNAUTHORIZED` | 401 | Header `Authorization: Bearer <key>` tidak disertakan atau token invalid. | Periksa API Key di pengaturan aplikasi. |
| `FORBIDDEN` | 403 | API Key tidak memiliki izin untuk resource ini atau IP tidak ada dalam whitelist. | Hubungi SecOps untuk memperbarui IP Whitelist / Scope. |
| `RATE_LIMIT_EXCEEDED` | 429 | Request melebihi batas RPS/RPM per-key atau global. Response menyertakan field `retry_after_seconds`. | Lakukan backoff retry sesuai `retry_after_seconds`. |
| `FILE_TOO_LARGE` | 413 | Ukuran file melebihi batas konfigurasi `MAX_SCAN_SIZE_MB` (default 100 MB). | Tolak upload di level gateway atau gunakan streaming URL. |
| `ZIP_BOMB_DETECTED` | 400 | Arsip terdeteksi sebagai decompression bomb (kedalaman > 5 level atau file > 1.000). | Tolak file sebagai potensi serangan DoS. |
| `SCAN_TIMEOUT` | 504 | Pemindaian memakan waktu melebihi batas `SCAN_TIMEOUT_SECONDS` (30 detik). | Ulangi scan atau periksa beban CPU/RAM server. |
| `INVALID_REQUEST_PAYLOAD`| 400 | Format multipart form, body JSON, atau query parameter tidak valid. | Periksa dokumentasi API spec request body. |
| `QUARANTINE_NOT_FOUND` | 404 | Record ID karantina tidak ditemukan atau sudah terhapus oleh auto-purge TTL. | File sudah kadaluarsa atau ID salah. |
| `ENGINE_UNAVAILABLE` | 503 | Daemon `clamd` sedang reload signature atau sedang restart. | Retry request setelah 5-10 detik. |
| `INTERNAL_SERVER_ERROR` | 500 | Terjadi kesalahan internal yang tidak tertangani pada service. | Laporkan ke administrator dan periksa log sistem. |
