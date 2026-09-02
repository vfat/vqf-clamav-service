# QA Execution Add-On — Workflow

> Orchestrator / entry point untuk add-on qa-execution AI Documentor.
> Panggil workflow ini ketika user meminta eksekusi test plan yang sudah ada —
> manual, automated, atau hybrid — dengan pelacakan status berbasis bukti nyata.

---

## Overview

```
Input: .ai-doc/qa/strategy/plans/test-plan-{slug}-{date}.md (valid) + codebase
        │
        ▼
┌──────────────────────────────────────────────┐
│  QA Execution Add-On                         │
│  (evidence-based, HALT pattern)              │
│          │                                   │
│  ┌───────┴────────┐                          │
│  ▼                ▼                          │
│ Eksekusi Item    Control Plane               │
│ (manual/auto/    .ai-doc/qa/execution/       │
│  hybrid)         qa-overview.md               │
│  hybrid)         (status + bukti semua item) │
└──────────────────────────────────────────────┘
```

### Output Utama

| Output | Format | Isi | Lokasi |
|---|---|---|---|
| **QA Overview** | Markdown | Control plane: status + bukti terakhir untuk setiap item test plan | `.ai-doc/qa/execution/qa-overview.md` |
| **Hasil Eksekusi** | Section di test plan | Kolom Status & Bukti di §7 test plan ikut diperbarui | `test-plan-{slug}-{date}.md` |

### Sumber Inspirasi

Add-on ini mengadaptasi pola `test-execution` dan `test-reporting` dari repo
qa-skills ke konvensi AI Documentor. Disiplin bukti disalin dari add-on tdd:
status tanpa bukti nyata tidak sah.

### Posisi dalam Rantai QA

```
qa-strategy → test-plan-*.md → qa-execution → laporan status
   (rencana)      (item test)     (eksekusi)     (bukti nyata)
```

- **Manual** — user menjalankan langkah, agent memandu & mencatat hasilnya.
- **Automated** — agent menjalankan command test di terminal (pytest,
  Playwright, Jest, go test, dll.) dan menangkap exit code + output sebagai bukti.
- **Hybrid** — setup manual + verifikasi otomatis.

Boundary: add-on ini mengorkestrasi dan melacak. Setup CI/CD pipeline,
scheduling, dan pembangunan framework automation dari nol di luar mandat.

---

## Flow Eksekusi

Gunakan 3 step berikut secara berurutan:

```text
Step 1: Activation
   ├── Cari test plan valid di .ai-doc/qa/strategy/plans/ (glob test-plan-*.md)
  ├── Jika tidak ada / tidak valid → HALT → tawarkan jalankan qa-strategy dulu
  ├── Baca item test: ID, jenis, level, prioritas, status saat ini
  ├── Deteksi runner yang tersedia di project (package.json, Makefile, dll.)
  ├── Usulkan urutan eksekusi (P0 dulu, lalu P1, P2)
  └── HALT → konfirmasi batch item yang akan dieksekusi

Step 2: Execution Loop
  ├── Untuk tiap item: tentukan mode (manual / automated / hybrid)
  ├── Automated: jalankan command nyata → tangkap exit code + output
  ├── Manual: pandu user langkah demi langkah → catat konfirmasi hasil
  ├── Tentukan status baru dari bukti (bukan dari perasaan):
  │   PASSED / FAILED / BLOCKED / DEFERRED
  ├── FAILED → catat output gagal + hipotesis penyebab
  ├── Update qa-overview.md + kolom §7 test plan per item selesai
  └── HALT → jika ada FAILED/BLOCKED, konfirmasi lanjut atau berhenti

Step 3: Wrap-Up & Reporting
  ├── Rekap sesi: jumlah PASSED/FAILED/BLOCKED/DEFERRED + sisa PLANNED
  ├── Soroti item P0 yang belum lolos (exit criteria test plan)
  ├── Update .ai-doc/3p.md
  └── HALT → konfirmasi rekap sebelum kontrol dikembalikan

Return kontrol ke core AI Documentor
```

Aturan flow:

1. Setiap titik keputusan user diakhiri `HALT`.
2. Status hanya boleh berubah kalau ada bukti yang menyertainya.
3. Satu sesi boleh mengeksekusi banyak item, tapi satu item = satu transisi
   status yang tercatat.

---

## Step 1 — Activation

Gate aktivasi (keras):

- Harus ada minimal satu file `test-plan-{slug}-{date}.md` di `.ai-doc/qa/strategy/plans/`.
- File harus punya section Test Items dengan ID + status.
- Jika tidak ada:

> Saya tidak menemukan test plan di `.ai-doc/qa/strategy/plans/`. qa-execution butuh test plan
> sebagai sumber item. Mau jalankan add-on **qa-strategy** dulu untuk menyusunnya?

**HALT** — jangan melanjutkan tanpa plan valid.

Deteksi runner (urutan prioritas):

| Sumber | Contoh |
|---|---|
| `package.json` scripts | `npm test`, `npx playwright test` |
| `Makefile` / `justfile` | `make test`, `just e2e` |
| Config framework | `pytest.ini`, `playwright.config.ts`, `go.mod` |
| Constitution / 3p.md | Catatan tooling dari sesi qa-strategy |

Usulkan urutan eksekusi: P0 → P1 → P2; dalam level yang sama, effort S dulu.
User bisa override kapan saja.

## Step 2 — Execution Loop

Untuk tiap item, agent memilih mode berdasar jenis test di plan:

| Jenis di Plan | Mode Default |
|---|---|
| Unit / Integration | Automated (biasanya sudah ada test file) |
| E2E / API | Automated bila suite ada; manual bila belum |
| UX / Visual / Eksploratif | Manual |

Status model & syarat bukti (disalin dari disiplin TDD):

| Status Baru | Syarat Bukti Wajib |
|---|---|
| `PASSED` | Command nyata + exit code 0 + output lolos ATAU konfirmasi eksplisit user (manual) |
| `FAILED` | Command/output kegagalan + hipotesis penyebab |
| `BLOCKED` | Alasan environment/data yang konkret (misal DB tidak bisa diakses) |
| `DEFERRED` | Alasan + target waktu |

Interaksi TDD: jika constitution menyatakan `TDD: Enabled` dan item berada di
level unit, eksekusinya lewat siklus RED→GREEN add-on tdd — qa-execution hanya
mencatat referensi buktinya, tidak menduplikasi siklusnya.

Update incremental: setiap item selesai langsung dicatat ke `qa-overview.md`
dan kolom §7 test plan — jangan menumpuk update di akhir sesi.

## Step 3 — Wrap-Up & Reporting

Rekap akhir sesi mencakup:

- Tabel status: PASSED / FAILED / BLOCKED / DEFERRED / sisa PLANNED.
- Item P0 yang belum lolos vs exit criteria di test plan.
- Daftar hipotesis untuk item FAILED (bahan perbaikan).
- Saran tindak lanjut: perbaiki FAILED, jalankan batch berikutnya, atau
  audit security via qa-security untuk area yang terdampak.

---

## Aturan Operasional

1. **Tanpa test plan valid, add-on tidak aktif** — arahkan ke qa-strategy.
2. **Tidak ada klaim PASSED tanpa bukti nyata** — sama kerasnya dengan aturan TDD.
3. **Agent menjalankan command test secara nyata** di terminal — bukan simulasi;
   exit code dan output wajib ditangkap apa adanya.
4. **FAILED bukan kegagalan agent** — catat jujur, susun hipotesis, biarkan
   user memutuskan tindak lanjut.
5. **Only-on-request** — aktif hanya lewat trigger eksplisit.
6. Integrasi: qa-strategy (sumber plan), tdd (unit-level), qa-security
   (temuan security jadi item test).

---

## Output Location

```
.ai-doc/
└── qa/
   ├── execution/
   │   └── qa-overview.md              # control plane (dibuat sekali, terus diperbarui)
   └── strategy/
      └── plans/
         └── test-plan-{slug}-{YYYY-MM-DD}.md  # kolom §7 Hasil Eksekusi diperbarui
```

Satu `qa-overview.md` per project, merujuk semua test plan yang pernah dieksekusi.
Jika hanya plan legacy yang ditemukan di root `.ai-doc/`, gunakan sebagai
referensi dan tawarkan pemindahan ke `qa/strategy/plans/` sebelum eksekusi.

---

## Referensi File

| File | Kegunaan |
|---|---|
| `add-on/qa-execution/steps/step-01-activation.md` | Detail gate aktivasi & deteksi runner |
| `add-on/qa-execution/steps/step-02-execution.md` | Detail loop eksekusi & transisi status |
| `add-on/qa-execution/steps/step-03-wrap-up.md` | Detail rekap & pelaporan |
| `add-on/qa-execution/template/qa-overview-template.md` | Template control plane |
