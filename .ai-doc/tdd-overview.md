# TDD Overview — `clamav-service`

> **Pusat Kontrol TDD Greenfield Project.**  
> Seluruh implementasi behavior baru wajib mengikuti siklus: `RED` (Test gagal terbukti) $\rightarrow$ `GREEN` (Implementasi minimal lolos) $\rightarrow$ `REFACTOR` (Pembersihan struktural tanpa merusak test).

---

## 1. Metadata

- **Project:** `clamav-service`
- **TDD Policy:** `Enabled`
- **Scope:** Greenfield Go Backend Service & Supervisor
- **Activated at:** 2026-09-02
- **Constitution:** [`.ai-doc/constitution.md`](file:///home/ubuntu/workspace/plan/clamav-service/.ai-doc/constitution.md)
- **Last updated:** 2026-09-02 17:25
- **Overall status:** `Active`

---

## 2. Progress Summary

| Metric | Count |
|---|---:|
| Total targets | 8 |
| PLANNED | 8 |
| RED | 0 |
| GREEN | 0 |
| REFACTORING | 0 |
| REFACTORED | 0 |
| BLOCKED | 0 |
| EXCEPTION | 0 |

---

## 3. TDD Registry

| ID | Component | Use Case / Behavior | Acceptance Criteria | Test File | Current Status | Last Evidence | Notes |
|---|---|---|---|---|---|---|---|
| **TDD-001** | `crypto/keygen` | Zero-Touch Master Key Auto-Gen & Env Inject | Menghasilkan key 256-bit dan menulis `ENCRYPTION_KEY=...` ke file `.env` jika belum ada | `internal/crypto/keygen_test.go` | `PLANNED` | — | Target 1 |
| **TDD-002** | `crypto/aes` | Field-Level AES-256-GCM Encrypt & Decrypt | Enkripsi payload dengan nonce acak 96-bit dan autentikasi tag; dekripsi menghasilkan plaintext asli | `internal/crypto/aes_test.go` | `PLANNED` | — | Target 2 |
| **TDD-003** | `storage/sqlite` | SQLite Schema Migrations & CRUD Operations | Auto-create tabel `scan_audit_logs`, `quarantine_records`, `api_keys`, `system_settings` di WAL mode | `internal/storage/sqlite_test.go` | `PLANNED` | — | Target 3 |
| **TDD-004** | `ratelimit` | In-Memory Token Bucket Rate Limiter | Mengizinkan request sesuai RPM, menolak dengan HTTP 429 saat kuota habis | `internal/ratelimit/limiter_test.go` | `PLANNED` | — | Target 4 |
| **TDD-005** | `clamd` | Unix Domain Socket Client & Stream Chunking | Menghubungkan ke Unix socket clamd, kirim chunk data biner, parse verdict `OK` / `FOUND` | `internal/clamd/client_test.go` | `PLANNED` | — | Target 5 |
| **TDD-006** | `quarantine` | Quarantine Vault Neutralization & Restore | Simpan file biner sebagai `.quarantine` permission `0600`, de-scramble saat restore, auto-whitelist SHA256 | `internal/quarantine/vault_test.go` | `PLANNED` | — | Target 6 |
| **TDD-007** | `alert` | Multi-Channel Notifier & Anti-Spam Throttling | Kirim payload Telegram/Discord/Email dan gabungkan pesan ke batch digest saat lonjakan > 5 alert/menit | `internal/alert/notifier_test.go` | `PLANNED` | — | Target 7 |
| **TDD-008** | `api` | REST API Gateway & Standard JSON Responses | Routing endpoint scan/quarantine/health dengan header standar, error codes, dan verdict JSON | `internal/api/handler_test.go` | `PLANNED` | — | Target 8 |

---

## 4. Cycle Detail

### TDD-001 — Zero-Touch Master Key Auto-Gen & Env Inject
- **Component:** `crypto/keygen`
- **Use case source:** [Lampiran L-003](file:///home/ubuntu/workspace/plan/clamav-service/.ai-doc/lampiran/L-003-Security-Encryption-and-Key-Lifecycle.md)
- **Acceptance criteria:** Jika `.env` tidak ada atau `ENCRYPTION_KEY` kosong, fungsi men-generate key 256-bit (prefix `clam_sec_`) dan menuliskannya ke `.env`. Jika sudah ada, membaca key eksisting tanpa menimpa.
- **Current status:** `PLANNED`

#### RED
- **Test file:** `internal/crypto/keygen_test.go`
- **Test name/target:** `TestEnsureMasterKey`
- **Command:** `—`
- **Exit status:** `—`
- **Failure evidence:** `—`
- **Verified at:** `—`

#### GREEN
- **Implementation file(s):** `internal/crypto/keygen.go`
- **Minimal change:** `—`
- **Command:** `—`
- **Exit status:** `—`
- **Passing evidence:** `—`
- **Verified at:** `—`

#### REFACTOR
- **Status:** `—`
- **Changes:** `—`
- **Regression command:** `—`
- **Exit status:** `—`
- **Regression evidence:** `—`
- **Verified at:** `—`

---

## 5. Blockers and Exceptions

*Tidak ada blocker saat ini.*

---

## 6. Change Log

| Date | Target | Phase | Change | Evidence / Reference |
|---|---|---|---|---|
| 2026-09-02 | ALL | Activation | TDD policy `Enabled` diaktifkan secara eksplisit oleh user | `.ai-doc/constitution.md` |

---

## 7. Operating Rules

1. Test unit/integration **MUST** ditulis sebelum kode implementasi.
2. Status `RED` wajib dibuktikan dengan menjalankan command `go test` dan mencatat failure output nyata.
3. Status `GREEN` dicapai hanya setelah implementasi minimal berhasil meloloskan test.
4. Refactoring hanya boleh dilakukan ketika test tetap `GREEN`.
5. Update file ini dan `.ai-doc/3p.md` setelah setiap transisi bermakna.
