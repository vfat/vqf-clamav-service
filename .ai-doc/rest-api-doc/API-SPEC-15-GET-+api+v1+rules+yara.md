# API Specification: `GET /api/v1/rules/yara`

## 1. Ringkasan Endpoint

| Item | Deskripsi |
|---|---|
| **Metode HTTP** | `GET` |
| **Path** | `/api/v1/rules/yara` |
| **Deskripsi** | Menampilkan daftar seluruh custom threat detection rules (YARA) yang terdaftar pada sistem. |
| **Autentikasi** | Sesuai `AUTH_MODE` (`none`, `basic`, `bearer`) |
| **Rate Limit** | Berlaku (default 100 RPM) |
| **Status** | `Published` (Active / Production Ready) |

---

## 2. Response Contract

### 2.1. Sukses (200 OK)
```json
{
  "success": true,
  "total": 1,
  "items": [
    {
      "id": "e37cbfa6-1e66-419b-8d75-471a921d7b30",
      "rule_name": "detect_webshell_php",
      "description": "Deteksi indikasi webshell PHP berdasar fungsi eksekusi shell",
      "content": "rule detect_webshell_php { strings: $a = \"passthru($_GET['cmd'])\" condition: $a }",
      "author": "sec-ops",
      "is_active": true,
      "created_at": "2026-09-07T07:28:40Z",
      "updated_at": "2026-09-07T07:28:40Z"
    }
  ]
}
```
