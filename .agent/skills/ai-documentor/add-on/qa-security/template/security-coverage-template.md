# Security Coverage — <Nama Project>

> Output add-on **qa-security** (AI Documentor). Dibuat: `<YYYY-MM-DD>` ·
> Re-audit terakhir: `<YYYY-MM-DD>` · Bootstrap: `ai-doc` / `codebase-reconstruction`

---

## 1. Metadata

| Field | Nilai |
|---|---|
| Project | `<nama>` |
| Scope audit | system / component / surface-group |
| Auth mechanism | `<jwt/session/oauth/api-key/Perlu Dikonfirmasi>` |
| Data sensitif | `<PII/payment/kredensial/...>` |
| Deployment | `<cloud/container/serverless>` |
| API style | rest / graphql / grpc / mixed |
| LLM features | ya / tidak |
| Compliance context | `<opsional, hanya konteks>` |

## 2. Attack Surface

| ID | Surface | Jenis | Auth? | Data Sensitif | Sumber Bukti |
|---|---|---|---|---|---|
| S-001 | `<endpoint/komponen>` | API publik / admin / webhook / upload | Ya/Tidak | `<jenis>` | `<artefak/file kode>` |

## 3. Coverage Matrix OWASP Top 10 (2025)

Semua kategori wajib muncul. Status: `Covered` / `Planned` / `N/A` / `Accepted Risk`.

| Kategori | Surface Terkait | Status | Bukti / Alasan |
|---|---|---|---|
| A01 Broken Access Control | S-00x | Planned | `<bukti/alasan>` |
| A02 Security Misconfiguration | S-00x | Planned | `<bukti/alasan>` |
| A03 Supply Chain Failures | S-00x | Planned | `<bukti/alasan>` |
| A04 Cryptographic Failures | S-00x | Planned | `<bukti/alasan>` |
| A05 Injection | S-00x | Planned | `<bukti/alasan>` |
| A06 Insecure Design | S-00x | Planned | `<bukti/alasan>` |
| A07 Authentication Failures | S-00x | Planned | `<bukti/alasan>` |
| A08 Data Integrity Failures | S-00x | N/A | `<alasan berbasis bukti>` |
| A09 Logging & Alerting Failures | S-00x | Planned | `<bukti/alasan>` |
| A10 Mishandling Exceptional Conditions | S-00x | Planned | `<bukti/alasan>` |

### 3a. OWASP LLM Top 10 (hanya jika ada fitur LLM)

| Kategori | Surface Terkait | Status | Bukti / Alasan |
|---|---|---|---|
| LLM01 Prompt Injection | S-00x | Planned | `<bukti/alasan>` |
| LLM02–LLM10 | ... | ... | ... |

## 4. Accepted Risk Register

Hanya diisi dengan persetujuan eksplisit user.

| Kategori/Temuan | Justifikasi | Disetujui Oleh | Tanggal |
|---|---|---|---|
| `<A0x / F-00x>` | `<mengapa diterima>` | `<nama user>` | `<YYYY-MM-DD>` |

## 5. Findings Register

Severity: Critical / High / Medium / Low. Keputusan: Remediate / Accept / Defer.

| ID | Temuan | Kategori | Surface | Severity | Sumber Bukti | Keputusan | Target |
|---|---|---|---|---|---|---|---|
| F-001 | `<deskripsi>` | A0x | S-00x | High | `<file/config/scan>` | Remediate | `<milestone/tanggal>` |
| F-002 | `<deskripsi>` | A0x | S-00x | Medium | `<...>` | Defer | `<target wajib>` |

### Action Items (hasil keputusan Remediate)

| Item | Temuan Terkait | Bisa Masuk Test Plan? |
|---|---|---|
| `<aksi konkret>` | F-001 | Ya → qa-strategy / Tidak (perubahan kode) |

## 6. Re-Audit Log

Diisi pada setiap siklus audit ulang (re-scan → triage → accept/remit).

| Tanggal | Pemicu | Perubahan Utama | Temuan Baru/Ditutup |
|---|---|---|---|
| `<YYYY-MM-DD>` | initial audit | matrix awal dibuat | F-001..F-00n dibuka |
| `<YYYY-MM-DD>` | `<re-scan/perubahan kode>` | `<kategori status berubah>` | `<F-xxx ditutup>` |

## 7. Unknowns & Asumsi

| Item | Jenis | Rencana Konfirmasi |
|---|---|---|
| `<klaim hasil rekonstruksi>` | Asumsi | `<cara verifikasi>` |
| `<auth mechanism belum pasti>` | Perlu Dikonfirmasi | `<tanya owner / cek config>` |

---

## Aturan Pemakaian File Ini

1. Semua kategori OWASP harus punya status — tidak boleh ada baris kosong.
2. `Accepted Risk` tanpa nama approver + tanggal tidak sah.
3. `Defer` tanpa target tidak sah.
4. Re-audit memperbarui file ini (tambah baris di §6), bukan membuat file baru.
5. Klaim tanpa bukti = `Perlu Dikonfirmasi`, bukan fakta.
