# Step 3: Findings & Mitigation Plan

> Susun findings register dari temuan statis (+ output scan yang ditempel user),
> lalu konfirmasikan keputusan mitigasi satu per satu sebelum file digenerate.
> **HALT** — konfirmasi akhir sebelum `.ai-doc/qa/security/security-coverage.md` ditulis.

---

## Input dari User

- **Output scan (opsional)** — hasil ZAP/OSV-Scanner/Semgrep yang ditempel user;
  dipakai sebagai bukti temuan dinamis
- **Keputusan per temuan** — Remediate / Accept / Defer
- **Target waktu** — untuk item Defer

---

## Yang Dilakukan Agent

### 1. Konsolidasi Temuan

Gabungkan temuan statis (Step 2) + temuan dinamis (output scan user):

| ID | Temuan | Kategori OWASP | Surface | Severity | Sumber Bukti |
|---|---|---|---|---|---|
| F-001 | `<deskripsi>` | `<A0x>` | `<S-00x>` | Critical/High/Medium/Low | `<file/config/scan output>` |

Aturan severity:

- Gunakan severity dari tool bila ada (misal output OSV-Scanner).
- Temuan statis pakai draft severity agent + alasan singkat.
- Duplikat antara statis & dinamis digabung, sumber dicatat keduanya.

### 2. Usulkan Keputusan per Temuan

| Opsi | Arti | Konsekuensi dokumentasi |
|---|---|---|
| **Remediate** | Diperbaiki | Masuk action item; bisa diturunkan jadi item test plan (qa-strategy) |
| **Accept** | Diterima sebagai risiko | Jadi `Accepted Risk`: justifikasi + persetujuan user tercatat |
| **Defer** | Ditunda | Alasan + target waktu wajib |

Presentasikan satu per satu:

> **F-001** [High] `<temuan>` — bukti: `<sumber>`.
> Usulan saya: **Remediate** karena [alasan].
> Keputusan Anda? (Remediate / Accept / Defer / tanya detail dulu)

**HALT per temuan** — jangan batch semua keputusan dalam satu pertanyaan.

### 3. Generate File Output

Setelah semua keputusan terkumpul, tampilkan ringkasan akhir:

```
=== SECURITY COVERAGE — RINGKASAN ===
Kategori : Covered=x, Planned=y, N/A=z, Accepted Risk=w
Temuan   : Total=n → Remediate=a, Accept=b, Defer=c
File     : .ai-doc/qa/security/security-coverage.md
```

> Siap saya tulis ke `.ai-doc/qa/security/security-coverage.md`?
> - ✅ **Tulis** — generate dari template
> - ✏️ **Revisi** — sebutkan bagiannya
> - ❌ **Batal**

**HALT** — tunggu konfirmasi. Setelah disetujui, generate file dari
`template/security-coverage-template.md`.

### 4. Tawarkan Tindak Lanjut

Setelah file tersimpan, tawarkan (jangan jalankan sendiri):

- Item `Remediate` → turunkan jadi test plan via add-on qa-strategy
- Item `Planned` → dieksekusi via add-on qa-execution saat plan tersusun
- Butuh konfigurasi runnable scanner (ZAP/OSV/Semgrep) → skill qa-skills
  (`security-testing`) bisa di-install terpisah untuk tooling-nya

---

## Output Step 3

```yaml
session_state:
  findings:
    - id: "F-001"
      decision: "remediate" | "accept" | "defer"
      decided_by: "user"
      defer_target: null | "<tanggal/milestone>"
  output_file: ".ai-doc/qa/security/security-coverage.md"
  follow_up_offered: true/false
```

---

## Aturan

- **JANGAN menulis file sebelum konfirmasi akhir user.**
- **Keputusan risiko milik user** — agent hanya mengusulkan dengan alasan.
- **Accept tanpa justifikasi tidak sah** — tolak dan minta alasannya.
- **Defer tanpa target waktu tidak sah** — minimal milestone, bukan "nanti".
- Update `.ai-doc/3p.md` setelah file ditulis.
