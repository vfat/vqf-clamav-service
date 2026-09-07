# API Specification: `POST /api/v1/rules/yara`

## 1. Ringkasan Endpoint

| Item | Deskripsi |
|---|---|
| **Metode HTTP** | `POST` |
| **Path** | `/api/v1/rules/yara` |
| **Deskripsi** | Mendaftarkan aturan YARA kustom baru, menyimpan berkas `.yara` ke disk direktori rules, menyimpan metadata ke tabel `yara_rules`, dan memicu zero-downtime hot-reload ClamAV daemon. |
| **Autentikasi** | Sesuai `AUTH_MODE` (`none`, `basic`, `bearer`) |
| **Rate Limit** | Berlaku (default 100 RPM) |
| **Status** | `Published` (Active / Production Ready) |

---

## 2. Request Contract

### 2.1. Headers
* `Content-Type: application/json` (Wajib)
* `Authorization: Bearer <token>` atau `Basic <base64>` (Opsional jika `AUTH_MODE != none`)
* `X-Consumer-Name: <nama-client>` (Opsional, dicatat sebagai `author`)

### 2.2. Request Body Schema
```json
{
  "rule_name": "detect_webshell_php",
  "description": "Deteksi indikasi webshell PHP berdasar fungsi eksekusi shell",
  "content": "rule detect_webshell_php { strings: $a = \"passthru($_GET['cmd'])\" condition: $a }",
  "author": "sec-ops"
}
```

| Field | Tipe | Wajib | Keterangan |
|---|---|---|---|
| `rule_name` | `string` | **Ya** | Nama pengenal unik aturan (alfanumerik & garis bawah, max 64 karakter). |
| `description` | `string` | Tidak | Deskripsi ancaman atau referensi CVE/IoC. |
| `content` | `string` | **Ya** | Konten sintaks aturan YARA lengkap (`rule <name> { ... condition: ... }`). |
| `author` | `string` | Tidak | Identitas analis pembuat aturan. |

---

## 3. Response Contract

### 3.1. Sukses (201 Created)
```json
{
  "success": true,
  "message": "YARA rule deployed and ClamAV reloaded successfully",
  "data": {
    "id": "e37cbfa6-1e66-419b-8d75-471a921d7b30",
    "rule_name": "detect_webshell_php",
    "description": "Deteksi indikasi webshell PHP berdasar fungsi eksekusi shell",
    "content": "rule detect_webshell_php { strings: $a = \"passthru($_GET['cmd'])\" condition: $a }",
    "author": "sec-ops",
    "is_active": true,
    "created_at": "2026-09-07T07:28:40Z",
    "updated_at": "2026-09-07T07:28:40Z"
  }
}
```

### 3.2. Error (400 Bad Request — Invalid Syntax / Duplicate Rule)
```json
{
  "success": false,
  "error": {
    "code": "INVALID_YARA_RULE",
    "message": "yara validation failed: unbalanced braces '{' and '}' in YARA rule",
    "details": null
  }
}
```
