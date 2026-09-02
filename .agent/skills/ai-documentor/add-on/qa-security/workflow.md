# QA Security Add-On — Workflow

> Orchestrator / entry point untuk add-on qa-security AI Documentor.
> Panggil workflow ini ketika user meminta security screening, audit keamanan,
> atau pemetaan coverage OWASP — termasuk pada codebase brownfield yang belum
> punya artefak `.ai-doc/` lengkap.

---

## Overview

```
Input: Codebase (brownfield/greenfield) + artefak .ai-doc jika ada
        │
        ▼
┌──────────────────────────────────────────┐
│  QA Security Add-On                      │
│  (evidence-based, HALT pattern)          │
│          │                               │
│  ┌───────┴────────┐                      │
│  ▼                ▼                      │
│ Security         Findings                │
│ Coverage Matrix  Register                │
│ (.ai-doc/qa/security/                    │
│  security-coverage.md)                   │
│  + rencana mitigasi                      │
└──────────────────────────────────────────┘
```

### Output Utama

| Output | Format | Isi | Lokasi |
|---|---|---|---|
| **Security Coverage Matrix** | Markdown | Pemetaan kategori OWASP Top 10 → status/bukti/justifikasi per attack surface | `.ai-doc/qa/security/security-coverage.md` |
| **Findings Register** | Section di file yang sama | Temuan terverifikasi dengan severity, bukti, dan keputusan mitigasi/accept | section 5 di `security-coverage.md` |

### Sumber Inspirasi

Add-on ini mengadaptasi konsep `owasp-coverage.md` dari skill
[security-testing](https://github.com/petrkindlmann/qa-skills) (qa-skills) ke konvensi
AI Documentor: evidence-based, only-on-request, HALT pattern. Konfigurasi runnable
tool scan (ZAP, OSV-Scanner, Semgrep) tetap milik skill asli — add-on ini mengatur
**apa yang dites, prioritasnya, dan pelacakan buktinya**, bukan cara menjalankan tool.

### Brownfield-First

Add-on ini dirancang berguna bahkan tanpa artefak `.ai-doc/`:

- Jika artefak inti ada → dipakai sebagai sumber bukti.
- Jika tidak ada → Step 1 merekonstruksi attack surface langsung dari codebase
  (config, route, dependency manifest), mengikuti semangat `AD-BF` (Brownfield Flow).
- Rekonstruksi ditandai jelas sebagai hasil inferensi kode, bukan dokumentasi final.

---

## Flow Eksekusi

Gunakan 3 step berikut secara berurutan:

```text
Step 1: Surface Mapping
  ├── Baca .ai-doc/3p.md dan constitution.md (jika ada)
  ├── Kumpulkan bukti: artefak .ai-doc ATAU rekonstruksi dari codebase
  ├── Identifikasi auth mechanism, API surface, deployment model
  ├── Petakan attack surface per komponen/endpoint
  ├── Tandai setiap klaim hasil inferensi kode sebagai Asumsi
  └── HALT → konfirmasi scope & threat model ringkas

Step 2: Coverage Assessment
  ├── Petakan tiap attack surface ke kategori OWASP Top 10 yang relevan
  ├── Untuk tiap kategori: tentukan status Planned / Covered / N/A / Accepted Risk
  ├── Accepted Risk wajib justifikasi + persetujuan eksplisit user
  ├── Kumpulkan temuan awal dari bukti statis (config, dependency manifest)
  └── HALT → review matrix sebelum ditulis

Step 3: Findings & Mitigation Plan
  ├── Susun findings register: severity, bukti, rekomendasi
  ├── Untuk tiap temuan: usulkan Remediate / Accept / Defer
  ├── Keputusan mitigasi dikonfirmasi user satu per satu
  ├── Generate .ai-doc/qa/security/security-coverage.md dari template
  └── HALT → konfirmasi sebelum file ditulis

Wrap-Up:
  ├── Update .ai-doc/3p.md
  ├── Tawarkan tindak lanjut: item test masuk test plan (qa-strategy),
  │   eksekusi scan via qa-execution, atau skill qa-skills untuk tooling
  └── Return kontrol ke core AI Documentor
```

Aturan flow:

1. Setiap titik keputusan user diakhiri `HALT`.
2. Add-on berhenti di dokumen; menjalankan scanner adalah pekerjaan eksekusi
   (qa-execution atau skill qa-skills), bukan bagian add-on ini.
3. Tidak ada kategori OWASP yang boleh hilang diam-diam — semua harus punya status.

---

## Step 1 — Surface Mapping

Baca `.ai-doc/3p.md` dan `.ai-doc/constitution.md` bila ada.

Sumber bukti, urut prioritas:

| Prioritas | Sumber | Contoh yang diambil |
|---|---|---|
| 1 | Artefak `.ai-doc/` (`rest-api-doc`, DCD, Dokumentasi-Codebase) | Endpoint publik, use case, komponen |
| 2 | Codebase nyata | Route/handler, middleware auth, config server |
| 3 | Dependency manifest | `package.json`, `requirements.txt`, lockfile |

Pertanyaan minimum ke user (skip yang sudah terjawab di artefak):

1. Auth mechanism apa: session cookie, JWT, OAuth/OIDC, API key?
2. Data sensitif apa yang diproses (PII, pembayaran, kredensial)?
3. Deployment model: cloud mana / container / serverless?
4. Ada compliance requirement (PCI, GDPR)? — hanya dicatat sebagai konteks;
   pemetaan regulasi bukan cakupan add-on ini.

Hasil Step 1: daftar attack surface + threat model ringkas.

## Step 2 — Coverage Assessment

Setiap kategori OWASP Top 10 (2025) harus punya salah satu status:

| Status | Arti | Syarat |
|---|---|---|
| `Planned` | Perlu dites, belum ada bukti test | Terhubung ke attack surface konkret |
| `Covered` | Sudah ada test/kontrol dengan bukti | Sebut buktinya (file/test/tool output) |
| `N/A` | Attack surface-nya memang tidak ada | Wajib alasan berbasis bukti |
| `Accepted Risk` | Sadar risikonya, diputuskan menerima | Justifikasi + persetujuan eksplisit user |

Aturan keras:

- `Accepted Risk` tidak boleh diputuskan agent — selalu minta keputusan user.
- `N/A` tanpa bukti = `Perlu Dikonfirmasi`, bukan N/A.
- Kategori LLM (OWASP LLM Top 10) hanya dinilai jika produk punya fitur AI/LLM.

## Step 3 — Findings & Mitigation Plan

Temuan bisa datang dari: inspeksi config/dependency saat Step 1–2, laporan user,
atau output scan yang user tempelkan (add-on tidak menjalankan scan sendiri).

Untuk tiap temuan, usulkan salah satu:

- **Remediate** — masuk action item; bisa diturunkan jadi item test plan.
- **Accept** — jadi `Accepted Risk` dengan justifikasi tercatat.
- **Defer** — ditunda dengan alasan + target waktu.

Keputusan final selalu milik user, dikonfirmasi satu per satu.

## Wrap-Up Checklist

- Semua kategori OWASP punya status + bukti/justifikasi.
- Setiap `Accepted Risk` mencantumkan siapa yang menyetujui dan kapan.
- File tersimpan di `.ai-doc/qa/security/security-coverage.md`.
- `.ai-doc/3p.md` mencerminkan langkah terakhir.
- Kontrol kembali ke core AI Documentor.

---

## Aturan Operasional

1. **Agent = pemetakan & pendokumentasi**, bukan pengambil keputusan risiko.
2. **HALT di setiap titik keputusan** — jangan auto-pilot.
3. **Evidence-based** — status tanpa bukti adalah `Perlu Dikonfirmasi`.
4. **Only-on-request** — add-on hanya aktif lewat trigger eksplisit.
5. **Tidak menjalankan scanner** — tooling dieksekusi di luar add-on ini;
   output scan yang ditempel user boleh dijadikan bukti.
6. Integrasi opsional: qa-strategy (risiko security masuk risk matrix),
   qa-execution (eksekusi item test security), persona/method untuk sesi analisis.

---

## Output Location

```
.ai-doc/
└── qa/
  └── security/
    └── security-coverage.md
```

Satu file per project. Re-audit berikutnya memperbarui file yang sama
(riwayat di change log), bukan membuat file baru. Jika file legacy masih berada
di root `.ai-doc/`, gunakan sebagai referensi dan tawarkan pemindahan sebelum
membuat file baru.

---

## Referensi File

| File | Kegunaan |
|---|---|
| `add-on/qa-security/steps/step-01-surface.md` | Detail pemetaan attack surface |
| `add-on/qa-security/steps/step-02-assessment.md` | Detail penilaian coverage OWASP |
| `add-on/qa-security/steps/step-03-findings.md` | Detail findings register & keputusan mitigasi |
| `add-on/qa-security/template/security-coverage-template.md` | Template output utama |
