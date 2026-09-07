# API Specification: `DELETE /api/v1/rules/yara/{id}`

## 1. Ringkasan Endpoint

| Item | Deskripsi |
|---|---|
| **Metode HTTP** | `DELETE` |
| **Path** | `/api/v1/rules/yara/{id}` atau `/api/v1/rules/yara?id={id}` |
| **Deskripsi** | Menghapus aturan YARA dari filesystem disk, menghapus record dari basis data SQLite, dan memicu reload clamd. |
| **Autentikasi** | Sesuai `AUTH_MODE` (`none`, `basic`, `bearer`) |
| **Rate Limit** | Berlaku (default 100 RPM) |
| **Status** | `Published` (Active / Production Ready) |

---

## 2. Response Contract

### 2.1. Sukses (200 OK)
```json
{
  "success": true,
  "message": "YARA rule removed and ClamAV reloaded successfully",
  "rule_id": "e37cbfa6-1e66-419b-8d75-471a921d7b30"
}
```

### 2.2. Error (404 Not Found)
```json
{
  "success": false,
  "error": {
    "code": "RULE_NOT_FOUND",
    "message": "YARA rule with id 'e37cbfa6-1e66-419b-8d75-471a921d7b30' not found",
    "details": null
  }
}
```
