# REST API Spec — `GET /api/v1/scan/jobs/{id}`

## 1. Metadata

| Method | Endpoint | Description | Status |
|---|---|---|---|
| `GET` | `/api/v1/scan/jobs/{id}` | Query polling status dan hasil pemrosesan scan job asinkron. | `Published` (Active / Production Ready) |

---

## 2. API Spec

### 2.1 Authentication
- Mengikuti `AUTH_MODE` (`none`, `basic`, `bearer`).

### 2.2 Path Parameters
- `id` (string, wajib): ID job pemindaian (misal: `job_1788350219`).

### 2.3 Response (200 OK — Completed Job)
```json
{
  "success": true,
  "data": {
    "job_id": "job_1788350219",
    "file_name": "large_archive.zip",
    "file_size": 104857600,
    "file_sha256": "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08",
    "callback_url": "https://consumer.corp.internal/webhooks/antivirus",
    "status": "COMPLETED",
    "verdict": "CLEAN",
    "virus_name": "",
    "error_msg": "",
    "created_at": "2026-09-07T16:00:00Z",
    "updated_at": "2026-09-07T16:00:15Z"
  }
}
```

### 2.4 Response (404 Not Found)
```json
{
  "success": false,
  "error": {
    "code": "JOB_NOT_FOUND",
    "message": "Scan job with id 'job_1788350219' not found",
    "details": null
  }
}
```
