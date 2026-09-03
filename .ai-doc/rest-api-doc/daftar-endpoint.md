# Daftar Endpoint — `clamav-service`

## 1. Ringkasan
Dokumen ini mendokumentasikan daftar lengkap endpoint HTTP REST API yang tersedia pada `clamav-service`. Seluruh endpoint pada daftar ini telah diimplementasikan di layer router [`internal/api/server.go`](file:///home/ubuntu/workspace/plan/clamav-service/internal/api/server.go) dan diverifikasi melalui unit test suite [`internal/api/handler_test.go`](file:///home/ubuntu/workspace/plan/clamav-service/internal/api/handler_test.go) serta container runtime live.

Setiap endpoint memiliki dokumen spesifikasi teknis detail tersendiri di dalam folder ini (`API-SPEC-*`).

---

## 2. Daftar Endpoint per Komponen

### 2.1. Komponen Antivirus Scanning

| Method | Endpoint | Description | Status | Spesifikasi Teknis |
|---|---|---|---|---|
| `POST` | `/api/v1/scan/file` | Pemindaian file sinkron via multipart file upload dengan instant verdict (`CLEAN` / `INFECTED`). | `Published` | [API-SPEC-01](./API-SPEC-01-POST-+api+v1+scan+file.md) |
| `POST` | `/api/v1/scan/stream` | Pemindaian raw binary stream chunked tanpa overhead multipart. | `Published` | [API-SPEC-02](./API-SPEC-02-POST-+api+v1+scan+stream.md) |
| `POST` | `/api/v1/scan/url` | Pemindaian file remote melalui unduhan stream URL publik / S3. | `Draft` | [API-SPEC-12](./API-SPEC-12-POST-+api+v1+scan+url.md) |
| `POST` | `/api/v1/scan/async` | Pengajuan job pemindaian asinkron dengan webhook callback (`202 Accepted`). | `Draft` | [API-SPEC-13](./API-SPEC-13-POST-+api+v1+scan+async.md) |

### 2.2. Komponen Quarantine Vault

| Method | Endpoint | Description | Status | Spesifikasi Teknis |
|---|---|---|---|---|
| `GET` | `/api/v1/quarantine` | Mengambil daftar file terisolasi di karantina dengan pagination dan filter status. | `Published` | [API-SPEC-03](./API-SPEC-03-GET-+api+v1+quarantine.md) |
| `POST` | `/api/v1/quarantine/restore` | Memulihkan (*restore*) file karantina ke sistem dan mendaftarkan hash ke whitelist. | `Published` | [API-SPEC-04](./API-SPEC-04-POST-+api+v1+quarantine+restore.md) |

### 2.3. Komponen Audit & Compliance

| Method | Endpoint | Description | Status | Spesifikasi Teknis |
|---|---|---|---|---|
| `GET` | `/api/v1/audit/export` | Mengunduh riwayat log audit pemindaian dalam format CSV atau JSON streaming. | `Published` | [API-SPEC-05](./API-SPEC-05-GET-+api+v1+audit+export.md) |

### 2.4. Komponen Web UI & Dashboard Authentication

| Method | Endpoint | Description | Status | Spesifikasi Teknis |
|---|---|---|---|---|
| `GET` | `/api/v1/auth/ui-status` | Memeriksa status proteksi keamanan antarmuka Web Admin UI. | `Published` | [API-SPEC-06](./API-SPEC-06-GET-+api+v1+auth+ui-status.md) |
| `POST` | `/api/v1/auth/ui-login` | Autentikasi sesi dashboard Web Admin UI (default password `123456`). | `Published` | [API-SPEC-07](./API-SPEC-07-POST-+api+v1+auth+ui-login.md) |
| `POST` | `/api/v1/auth/ui-password` | Mengubah password dashboard Web Admin UI dan menyimpannya secara persisten ke SQLite. | `Published` | [API-SPEC-08](./API-SPEC-08-POST-+api+v1+auth+ui-password.md) |

### 2.5. Komponen Health, Observability & Ops

| Method | Endpoint | Description | Status | Spesifikasi Teknis |
|---|---|---|---|---|
| `GET` | `/api/v1/health` | Healthcheck & readiness probe (Selalu public / bypass auth). | `Published` | [API-SPEC-09](./API-SPEC-09-GET-+api+v1+health.md) |
| `GET` | `/healthz` | Kubernetes lightweight liveness/readiness probe (Selalu public / bypass auth). | `Published` | [API-SPEC-10](./API-SPEC-10-GET-+healthz.md) |
| `GET` | `/api/v1/metrics` | Endpoint Prometheus metrics (`text/plain; version=0.0.4`). | `Published` | [API-SPEC-11](./API-SPEC-11-GET-+api+v1+metrics.md) |

---

## 3. Catatan Validitas
- Seluruh endpoint dengan status `Published` telah terhubung langsung ke handler aktif di router [`internal/api/server.go`](file:///home/ubuntu/workspace/plan/clamav-service/internal/api/server.go) dan database SQLite.
- Endpoint dengan status `Draft` (`/scan/url`, `/scan/async`) terdaftar pada blueprint arsitektur untuk rilis berikutnya.

---

## 4. Asumsi, Keamanan, dan Kebijakan Akses
- **API Authorization Policy (`AUTH_MODE`)**:
  - `none` (*Default*): Bebas autentikasi untuk kemudahan integrasi di internal VPC/Docker network.
  - `basic`: Mewajibkan header `Authorization: Basic <base64>` menggunakan kredensial `AUTH_BASIC_USER` & `AUTH_BASIC_PASS`.
  - `bearer`: Mewajibkan header `Authorization: Bearer <token>` atau `X-API-Key: <token>`.
  - *Exemption*: Endpoint `/healthz` dan `/api/v1/health` selalu di-bypass agar container healthcheck tidak terblokir.
- **Web Admin UI Security**: Dilengkapi *Security Lock Screen* dengan default password `123456`, dapat diubah sewaktu-waktu dan disimpan dengan hash SHA-256 salted di tabel SQLite `system_settings`.
- **Rate Limiting**: Header `X-RateLimit-Limit`, `X-RateLimit-Remaining`, dan `X-RateLimit-Reset` disertakan pada setiap response.
