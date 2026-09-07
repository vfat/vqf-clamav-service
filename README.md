# 🛡️ ClamAV Service (`vqf-clamav-service`)

> **Enterprise-Grade, High-Performance Antivirus REST API & Real-time Threat Intelligence Hub**  
> *All-in-One Docker Container • In-Memory Unix Socket Scanning • Built-in Quarantine Vault • Zero-Touch Master Key Security • Multi-Channel Alerting*

[![Go](https://img.shields.io/badge/Go-1.22-00ADD8?style=flat&logo=go)](https://go.dev)
[![Docker Hub](https://img.shields.io/badge/Docker%20Hub-vickyfatrian%2Fvqf--clamav--service-2496ED?style=flat&logo=docker)](https://hub.docker.com/r/vickyfatrian/vqf-clamav-service)
[![TDD](https://img.shields.io/badge/TDD-100%25%20Passed-10b981?style=flat)](file:///.ai-doc/tdd-overview.md)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

---

## 🌟 Key Architecture Highlights

1. **⚡ Zero-Disk-I/O In-Memory Streaming**: File payloads are piped directly to ClamAV daemon via Unix Domain Socket (`zINSTREAM\0` chunked streaming), eliminating disk overhead and providing sub-50ms verdict latencies.
2. **🔐 Zero-Touch Master Key & AES-256-GCM**: Auto-generates a cryptographic 256-bit key on first cold start and injects `ENCRYPTION_KEY` into `.env` with strict `0600` permissions.
3. **🛡️ Built-in Neutralized Quarantine Vault**: Automatically encrypts infected binaries using authenticated **AES-256-GCM** with a constant 14-byte Magic Header (`VQF_AESGCM_V1\n`), strict `0600` permissions, dual-mode restore (direct download / S3), auto-whitelisting SHA-256 hashes, and memory bomb safeguards.
4. **🚨 Multi-Channel Alerting & Flood Throttling**: Real-time notifications to Telegram Bot and Discord Webhooks with sliding-window flood protection (> 5 threats/min $\rightarrow$ batch digest).
5. **🗄️ SQLite WAL Mode Persistence**: Transactional logging with auto-purging policies (3-day audit logs, 7-day quarantine retention) and streaming CSV/JSON exports.
6. **🚀 High-Performance gRPC Streaming**: Built-in protobuf streaming service (`Port 9090`) for zero-latency inter-service RPC communications with auth interceptors.
7. **🌐 Anti-SSRF Remote URL Scanner**: Fetches remote payloads securely with DNS-rebinding guards and blocking of loopback, RFC1918, and Cloud metadata IPs (`169.254.169.254`).
8. **📥 Asynchronous Worker Queue & Spool**: Non-blocking scanning (`202 Accepted`) with spool storage (`/data/spool`), concurrent worker pools, and automated webhook callbacks.
9. **🛡️ Zero-Downtime Custom YARA Engine**: Dynamic hot-reloading of `.yara` threat signatures directly into ClamAV daemon via Unix socket commands without restarts.
10. **💻 Embedded Web Admin UI (SPA)**: Zero-dependency responsive dark glassmorphism dashboard built with Vanilla JS & CSS, embedded directly inside the Go binary.

---

## 🐳 Docker Hub Image

Official container images are available on Docker Hub:  
👉 [**`vickyfatrian/vqf-clamav-service:latest`**](https://hub.docker.com/r/vickyfatrian/vqf-clamav-service) (Release tag: [**`v1.1.0`**](https://hub.docker.com/r/vickyfatrian/vqf-clamav-service/tags))

### One-Liner Quick Run (Without cloning repository):
```bash
docker run -d \
  --name clamav-service \
  -p 8080:8080 \
  -p 9090:9090 \
  -v $(pwd)/data:/data \
  -v clamav_signatures:/var/lib/clamav \
  vickyfatrian/vqf-clamav-service:latest
```

---

## 🚀 Quick Start with Docker Compose

### 1. Clone and Prepare Configuration
```bash
git clone https://github.com/vfat/vqf-clamav-service.git
cd vqf-clamav-service
cp .env.example .env
```

### 2. Launch the All-in-One Container
```bash
docker compose up -d
```

The service will boot:
- **Web Admin UI & REST API**: `http://localhost:8080`
- **gRPC Antivirus Scanner**: `localhost:9090`
- **Health Check Probe**: `http://localhost:8080/api/v1/health`
- **Prometheus Metrics**: `http://localhost:8080/api/v1/metrics`

---

## 🔒 Security & Authorization Policies

### 1. API Authorization (`AUTH_MODE`)
Configured in `.env` via `AUTH_MODE`:
* **`none` (Default)**: Unrestricted API access. Recommended for internal VPCs / Docker bridge networks.
* **`basic`**: Requires HTTP Basic Auth (`AUTH_BASIC_USER` and `AUTH_BASIC_PASS`).
* **`bearer`**: Requires `Authorization: Bearer <AUTH_BEARER_TOKEN>` or `X-API-Key: <AUTH_BEARER_TOKEN>`.
* *Note: Healthcheck probes (`/healthz`, `/api/v1/health`) are always public/exempt to ensure continuous monitoring.*

### 2. Web Admin UI Password Protection
* Accessing `http://localhost:8080` displays a modern security lock screen prompt: *"Enter your password to access the dashboard"*.
* **Default Password**: `123456`
* **Password Change**: Click **"Change Password"** in the top navigation bar to set a custom secret. The new password is automatically salted, hashed with SHA-256, and stored persistently in SQLite (`system_settings`).

---

## 🔌 API Usage Examples

### 1. Synchronous File Scan
```bash
curl -X POST http://localhost:8080/api/v1/scan/file \
  -H "X-Consumer-Name: Billing-App" \
  -F "file=@/path/to/invoice.pdf"
```

### 2. Remote URL Scan (Protected with Anti-SSRF)
```bash
curl -X POST http://localhost:8080/api/v1/scan/url \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://secure-downloads.example.com/driver.bin",
    "expected_sha256": "optional-hash-for-integrity-verification"
  }'
```

### 3. Asynchronous Scan with Webhook Notification
```bash
curl -X POST http://localhost:8080/api/v1/scan/async \
  -F "file=@/path/to/large_archive.zip" \
  -F "callback_url=https://my-app.internal/webhooks/scan-verdict"
```
**Response (202 Accepted):**
```json
{
  "success": true,
  "job_id": "0191c984-7a1b-7000-8800-123456789abc",
  "status": "QUEUED",
  "message": "Scan job enqueued successfully"
}
```

### 4. Deploy Custom YARA Threat Signature
```bash
curl -X POST http://localhost:8080/api/v1/rules/yara \
  -H "Content-Type: application/json" \
  -d '{
    "rule_name": "detect_webshell_php",
    "author": "secops",
    "description": "Detects PHP command execution injection",
    "content": "rule detect_webshell_php {\n    strings:\n        $a = \"passthru($_GET[\x27cmd\x27])\"\n    condition:\n        $a\n}"
  }'
```

### 5. Restore Quarantined File & Whitelist Hash
```bash
curl -X POST http://localhost:8080/api/v1/quarantine/restore \
  -H "Content-Type: application/json" \
  -d '{
    "quarantine_id": "Q-20260902-8f92a10b",
    "restored_by": "secops-lead@company.com",
    "reason": "Verified false positive internal payroll document",
    "auto_whitelist": true
  }'
```

---

## 🧪 Running Unit & Integration Tests (TDD Suite)

```bash
go test -v -count=1 ./...
```

---

## 📖 Architecture & Design Documentation

Comprehensive technical documentation is maintained under [`.ai-doc/`](file:///.ai-doc/):

- 📐 [**C4 Component Diagrams**](file:///.ai-doc/C4-Component-Diagrams.md)
- 📋 [**Feature Inventory Matrix**](file:///.ai-doc/Dokumentasi-Fitur.md)
- 🧭 [**Grouped Use Case Documentation**](file:///.ai-doc/Dokumentasi-Komponen-Usecase.md)
- 🏗️ [**Design Component Documents (DCD-01 s/d DCD-04)**](file:///.ai-doc/desain-component-document/)
- 🗄️ [**Database ERD & Data Dictionaries**](file:///.ai-doc/desain-database-document/ERD-overview.md)
- 📜 [**REST API Endpoint Indeks**](file:///.ai-doc/rest-api-doc/daftar-endpoint.md)
- 📑 [**REST API Specifications Suite (API-SPEC-01 s/d 17)**](file:///.ai-doc/rest-api-doc/)
- 🖥️ [**Web UI Walkthrough Guide**](file:///.ai-doc/Web-UI-Walkthrough.md)
- 🎯 [**TDD Control Plane (17/17 Targets GREEN)**](file:///.ai-doc/tdd-overview.md)
- 📚 [**Lampiran Series (L-001 s/d L-006)**](file:///.ai-doc/lampiran/)

---

## 📄 License
MIT License. Open-source enterprise antivirus API.
