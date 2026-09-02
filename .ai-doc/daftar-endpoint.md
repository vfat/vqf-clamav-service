# Daftar Endpoint — `clamav-service`

## 1. Ringkasan
Dokumen ini mendokumentasikan daftar lengkap endpoint HTTP REST API yang tersedia pada `clamav-service`. Seluruh endpoint pada daftar ini telah diimplementasikan di layer router [`internal/api/server.go`](file:///home/ubuntu/workspace/plan/clamav-service/internal/api/server.go) dan diverifikasi melalui unit test suite [`internal/api/handler_test.go`](file:///home/ubuntu/workspace/plan/clamav-service/internal/api/handler_test.go).

---

## 2. Daftar Endpoint per Komponen

### 2.1. Komponen Antivirus Scanning

| Method | Endpoint | Description | Status |
|---|---|---|---|
| `POST` | `/api/v1/scan/file` | Pemindaian file sinkron via multipart file upload dengan instant verdict (`CLEAN` / `INFECTED`). | `Published` |
| `POST` | `/api/v1/scan/stream` | Pemindaian raw binary stream chunked tanpa overhead multipart. | `Published` |
| `POST` | `/api/v1/scan/url` | Pemindaian file remote melalui unduhan stream URL publik / S3. | `Draft` |
| `POST` | `/api/v1/scan/async` | Pengajuan job pemindaian asinkron dengan webhook callback (`202 Accepted`). | `Draft` |

### 2.2. Komponen Quarantine Vault

| Method | Endpoint | Description | Status |
|---|---|---|---|
| `GET` | `/api/v1/quarantine` | Mengambil daftar file terisolasi di karantina dengan pagination dan filter status. | `Published` |
| `POST` | `/api/v1/quarantine/restore` | Memulihkan (*restore*) file karantina ke sistem dan mendaftarkan hash ke whitelist. | `Published` |

### 2.3. Komponen Audit & Compliance

| Method | Endpoint | Description | Status |
|---|---|---|---|
| `GET` | `/api/v1/audit/export` | Mengunduh riwayat log audit pemindaian dalam format CSV atau JSON streaming. | `Published` |

### 2.4. Komponen Health, Observability & Ops

| Method | Endpoint | Description | Status |
|---|---|---|---|
| `GET` | `/api/v1/health` | Healthcheck & readiness probe status database, socket ClamAV, dan service uptime. | `Published` |
| `GET` | `/healthz` | Kubernetes lightweight liveness/readiness probe. | `Published` |
| `GET` | `/api/v1/metrics` | Endpoint Prometheus metrics (`text/plain; version=0.0.4`). | `Published` |

---

## 3. Catatan Validitas
- Seluruh endpoint dengan status `Published` telah terhubung langsung ke handler aktif dan database SQLite.
- Endpoint dengan status `Draft` (`/scan/url`, `/scan/async`) terdaftar pada blueprint arsitektur untuk rilis berikutnya.

---

## 4. Asumsi, Keamanan, dan Kebijakan Akses
- **Authentication**: Mendukung validasi token via header `Authorization: Bearer <API_KEY>` jika `REQUIRE_API_KEY=true`.
- **Rate Limiting**: Header `X-RateLimit-Limit`, `X-RateLimit-Remaining`, dan `X-RateLimit-Reset` disertakan pada setiap response.
