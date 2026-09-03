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
3. **🛡️ Built-in Neutralized Quarantine Vault**: Automatically scrambles infected binaries (XOR mask) and stores them under `/data/quarantine/` with permissions `0600`, supporting dual-mode restore (direct download / S3) and auto-whitelisting SHA-256 hashes.
4. **🚨 Multi-Channel Alerting & Flood Throttling**: Real-time notifications to Telegram Bot and Discord Webhooks with sliding-window flood protection (> 5 threats/min $\rightarrow$ batch digest).
5. **🗄️ SQLite WAL Mode Persistence**: Transactional logging with auto-purging policies (3-day audit logs, 7-day quarantine retention) and streaming CSV/JSON exports.
6. **💻 Embedded Web Admin UI (SPA)**: Zero-dependency responsive dark glassmorphism dashboard built with Vanilla JS & CSS, embedded directly inside the Go binary.

---

## 🐳 Docker Hub Image

Image resmi tersedia di Docker Hub:  
👉 [**`vickyfatrian/vqf-clamav-service:latest`**](https://hub.docker.com/r/vickyfatrian/vqf-clamav-service) (atau tag rilis [**`v1.0.0`**](https://hub.docker.com/r/vickyfatrian/vqf-clamav-service/tags))

### One-Liner Quick Run (Tanpa Clone Source Code):
```bash
docker run -d \
  --name clamav-service \
  -p 8080:8080 \
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
- **Web UI & REST API**: `http://localhost:8080`
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

**Response (Clean File):**
```json
{
  "success": true,
  "verdict": "CLEAN",
  "data": {
    "file_name": "invoice.pdf",
    "file_size": 2048576,
    "file_sha256": "275a021bbfb6489e54d471899f7db9d1663fc695ec2fe2a2c4538aabf651fd0f",
    "scan_duration_ms": 42,
    "scanned_at": "2026-09-02T19:30:00Z"
  }
}
```

**Response (Malware Detected):**
```json
{
  "success": true,
  "verdict": "INFECTED",
  "threat": {
    "virus_name": "Win.Trojan.Agent-1234",
    "severity": "HIGH",
    "action_taken": "QUARANTINED",
    "quarantine_id": "Q-20260902-8f92a10b"
  },
  "data": {
    "file_name": "invoice.pdf.exe",
    "file_size": 1048576,
    "file_sha256": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
    "scan_duration_ms": 38
  }
}
```

### 2. Restore Quarantined File & Whitelist Hash
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

## 🧪 Running Unit Tests (TDD Suite)

```bash
go test -v ./...
```

---

## 📖 Architecture & Design Documentation

Comprehensive technical documentation is maintained under [`.ai-doc/`](file:///.ai-doc/):

- 📐 [**C4 Component Diagrams**](file:///.ai-doc/C4-Component-Diagrams.md)
- 📋 [**Feature Inventory Matrix**](file:///.ai-doc/Dokumentasi-Fitur.md)
- 🧭 [**Grouped Use Case Documentation**](file:///.ai-doc/Dokumentasi-Komponen-Usecase.md)
- 📜 [**REST API Endpoint List**](file:///.ai-doc/rest-api-doc/daftar-endpoint.md)
- 📑 [**REST API Specification & Swimlanes**](file:///.ai-doc/rest-api-doc/)
- 🎯 [**TDD Control Plane**](file:///.ai-doc/tdd-overview.md)
- 📚 [**Lampiran Series (L-001 s/d L-006)**](file:///.ai-doc/lampiran/)

---

## 📄 License
MIT License. Open-source enterprise antivirus API.
