# QA Strategy Add-On — Workflow

> Orchestrator / entry point untuk add-on qa-strategy AI Documentor.
> Panggil workflow ini ketika user meminta penyusunan strategi testing atau
> rencana test berbasis artefak `.ai-doc/` yang sudah ada.

---

## Overview

```
Input: Artefak .ai-doc existing (project-overview, DCD, REST API doc)
        │
        ▼
┌─────────────────────────────────────┐
│  QA Strategy Add-On                 │
│  (evidence-based, HALT pattern)     │
│          │                          │
│  ┌───────┴────────┐                 │
│  ▼                ▼                 │
│ Test Strategy   Test Plan           │
│ (pyramid, gates,(per sprint/release)│
│  risk matrix)   │                   │
│                 ▼                   │
│ .ai-doc/qa/strategy/plans/          │
│          test-plan-*.md             │
└─────────────────────────────────────┘
```

### Output Utama

| Output | Format | Isi | Lokasi |
|---|---|---|---|
| **Test Strategy** | Markdown | Risk matrix, test pyramid, quality gates, entry/exit criteria | `.ai-doc/qa/strategy/test-strategy.md` |
| **Test Plan** | Markdown per sprint/release | Feature → test case mapping, prioritas, effort, scope boundary | `.ai-doc/qa/strategy/plans/test-plan-{slug}-{YYYY-MM-DD}.md` |

### Sumber Inspirasi

Add-on ini mengadaptasi pola `test-strategy` dan `test-planning` dari repo
[qa-skills](https://github.com/petrkindlmann/qa-skills) ke konvensi AI Documentor:
evidence-based, only-on-request, dan output di bawah `.ai-doc/`.

---

## Flow Eksekusi

Gunakan 3 step berikut secara berurutan:

```text
Step 1: Context Gathering
  ├── Baca .ai-doc/3p.md dan constitution.md
  ├── Inventarisasi artefak sumber bukti (overview, DCD, ERD, REST API doc)
  ├── Ekstraksi daftar komponen, use case, endpoint, dan risiko yang tercatat
  ├── Tandai area tanpa bukti sebagai Perlu Dikonfirmasi
  └── HALT → konfirmasi cakupan strategi

Step 2: Strategy Building
  ├── Susun risk matrix (Impact × Likelihood) dari bukti Step 1
  ├── Tentukan test pyramid per level (unit / integration / E2E)
  ├── Tetapkan quality gates untuk CI dan release
  ├── Definisikan entry/exit criteria per level test
  ├── Pilih tool dengan justifikasi (opsional, hanya jika user minta)
  └── HALT → review strategy sebelum ditulis ke file

Step 3: Test Planning
  ├── Turunkan scope dari strategy + input sprint/release dari user
  ├── Mapping feature/use case → jenis test → prioritas
  ├── Estimasi effort relatif (S/M/L) dan tandai yang ditunda
  ├── Generate test plan dari template
  └── HALT → konfirmasi sebelum file ditulis

Wrap-Up:
  ├── Update .ai-doc/3p.md
  ├── Tawarkan sinkronisasi dokumen tetangga bila ada klaim baru
  └── Return kontrol ke core AI Documentor
```

Aturan flow:

1. Setiap titik keputusan user diakhiri `HALT`.
2. Step 3 hanya dijalankan jika user memintanya — strategy saja adalah output yang sah.
3. Jangan membuat status "selesai" untuk section yang buktinya belum ada.

---

## Step 1 — Context Gathering

Baca `.ai-doc/3p.md` dan `.ai-doc/constitution.md` lebih dulu.

Inventarisasi artefak sumber:

| Artefak | Yang diekstrak |
|---|---|
| `project-overview.md` | Critical user journeys, stakeholder, constraint |
| `Dokumentasi-Codebase.md` | Komponen nyata, entry point, dependency |
| `DCD-*.md` | Use case per komponen — kandidat coverage |
| `ERD` / data dictionary | Entitas kritikal, aturan integritas data |
| `rest-api-doc` | Endpoint publik — kandidat API test |

Area yang tidak punya bukti ditandai `Perlu Dikonfirmasi`, bukan ditebak.

Pertanyaan minimum ke user:

> Strategi ini untuk cakupan apa: seluruh sistem, satu komponen, atau satu release?

## Step 2 — Strategy Building

Semua keputusan strategi harus bisa dilacak ke bukti Step 1.

- **Risk matrix** — setiap risiko punya sumber bukti (incident record, churn tinggi,
  endpoint publik, integrasi pihak ketiga). Skor Impact × Likelihood.
- **Test pyramid** — proporsi unit/integration/E2E disesuaikan karakter sistem
  (API-heavy ≠ UI-heavy). Berikan alasan, bukan angka generik.
- **Quality gates** — apa yang memblokir merge/deploy, beserta ambangnya.
- **Entry/exit criteria** — per level test, objektif dan dapat dicek.

Tool hanya dibahas jika user memintanya; satu rekomendasi default + alasannya,
alternatif disebut sekilas.

## Step 3 — Test Planning

Plan adalah turunan langsung dari strategy — jangan menambah scope baru.

- Input wajib: target sprint/release dan daftar fitur dari user.
- Prioritas mengikuti risk matrix, bukan selera.
- Item yang tidak sempat ditandai `Deferred` dengan alasannya, bukan dihapus diam-diam.

## Wrap-Up Checklist

- Strategy dan/atau plan tersimpan di lokasi output yang benar.
- Setiap section tanpa bukti ditandai `Perlu Dikonfirmasi` atau `Asumsi`.
- `.ai-doc/3p.md` mencerminkan langkah terakhir.
- Kontrol kembali ke core AI Documentor.

---

## Aturan Operasional

1. **Agent = fasilitator**, keputusan strategi tetap milik user.
2. **HALT di setiap titik keputusan** — jangan auto-pilot.
3. **Evidence-based** — setiap klaim risiko/prioritas menunjuk sumbernya.
4. **Only-on-request** — add-on hanya aktif lewat trigger eksplisit.
5. Integrasi opsional: persona (fasilitator sesi), method (teknik prioritisasi).

---

## Output Location

```
.ai-doc/
└── qa/
  └── strategy/
    ├── test-strategy.md
    └── plans/
      └── test-plan-{sprint-atau-release-slug}-{YYYY-MM-DD}.md
```

Format slug: lowercase, gunakan `-` sebagai separator, maksimal 5 kata.
Jika artefak legacy masih berada di root `.ai-doc/`, baca sebagai referensi dan
tawarkan pemindahan sebelum membuat artefak baru agar tidak terjadi duplikasi.

---

## Referensi File

| File | Kegunaan |
|---|---|
| `add-on/qa-strategy/steps/step-01-context.md` | Detail pengumpulan konteks & bukti |
| `add-on/qa-strategy/steps/step-02-strategy.md` | Detail penyusunan strategy |
| `add-on/qa-strategy/steps/step-03-plan.md` | Detail penurunan test plan |
| `add-on/qa-strategy/template/test-strategy-template.md` | Template output strategy |
| `add-on/qa-strategy/template/test-plan-template.md` | Template output test plan |
