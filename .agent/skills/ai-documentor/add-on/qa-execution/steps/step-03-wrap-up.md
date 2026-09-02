# Step 3: Wrap-Up & Reporting

> Rekap hasil sesi, soroti gap terhadap exit criteria, dan tutup sesi
> dengan kontrol kembali ke core AI Documentor.
> **HALT** — konfirmasi rekap sebelum sesi ditutup.

---

## Input dari User

- **Konfirmasi rekap** — setuju ditutup / ada item yang mau diulang
- **Keputusan tindak lanjut** — perbaiki FAILED / batch berikutnya / cukup

---

## Yang Dilakukan Agent

### 1. Susun Rekap Sesi

```
=== QA EXECUTION — REKAP SESI ===
Plan    : .ai-doc/qa/strategy/plans/test-plan-auth-2026-08-25.md
Batch   : 5 item (P0)

Status akhir sesi ini:
  PASSED   : 3  (QA-001, QA-002, QA-003)
  FAILED   : 1  (QA-005 — hipotesis: error handler lempar exception)
  BLOCKED  : 1  (QA-004 — DB staging tidak dapat diakses)

Posisi plan keseluruhan:
  PASSED 5/12 · PLANNED 6 · FAILED 1 · BLOCKED 1 · DEFERRED 0

Exit criteria (dari plan):
  ❌ Belum terpenuhi — 1 item P0 masih FAILED
```

### 2. Soroti Gap terhadap Exit Criteria

Bandingkan posisi plan dengan exit criteria di test plan:

- Item P0 yang belum PASSED → sebutkan satu per satu.
- Quality gates dari test strategy yang terancam (misal coverage threshold).
- Jika semua P0 lolos → nyatakan exit criteria terpenuhi untuk batch ini.

### 3. Tawarkan Tindak Lanjut

> Mau lanjut ke mana?
> - 🔧 **Perbaiki QA-005** — keluar add-on, masuk mode kerja normal (perbaikan kode)
> - ▶️ **Batch berikutnya** — 4 item P1
> - 🔁 **Re-run QA-004** — setelah environment siap
> - 🛡️ **Audit area terdampak** — trigger add-on qa-security bila FAILED menyentuh surface keamanan
> - ✅ **Cukup untuk sesi ini**

### 4. Tutup Sesi

Pastikan sebelum return kontrol:

- `qa-overview.md` konsisten dengan §7 semua test plan yang disentuh.
- Semua FAILED punya hipotesis; semua BLOCKED punya alasan konkret.
- `.ai-doc/3p.md` mencerminkan langkah terakhir.
- Change Log qa-overview mendapat baris baru untuk sesi ini.

---

## Output Step 3

```yaml
session_state:
  session_summary:
    passed: 3
    failed: 1
    blocked: 1
    deferred: 0
    remaining_planned: 6
  exit_criteria_met: false
  follow_up_offered: true
  session_closed: true
```

---

## Aturan

- **JANGAN tutup sesi dengan control plane basah** — semua update harus
  sudah tertulis di qa-overview.md dan test plan.
- **JANGAN klaim exit criteria terpenuhi** kalau masih ada P0 belum PASSED.
- **REKOMENDASIKAN prioritas perbaikan FAILED P0** — tapi keputusan milik user.
- Return kontrol ke core AI Documentor setelah konfirmasi rekap.
