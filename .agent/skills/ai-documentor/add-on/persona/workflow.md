# Persona Add-On — Workflow

> Orchestrator / entry point untuk add-on persona AI Documentor.
> Panggil workflow ini ketika user ingin menyiapkan atau membuat persona.

---

## Overview

```
User meminta "add on persona" atau "buatkan persona"
        │
        ▼
┌───────────────────────────────────┐
│  Persona Add-On                   │
│  (pilih template → generate)      │
│          │                        │
│  ┌───────┴───────┐                │
│  ▼               ▼                │
│ Pilih           Generated         │
│ Blueprint       Persona Files     │
│ (templates)     (.ai-doc/personas │
│  di add-on)     /<nama-persona>/) │
└───────────────────────────────────┘
```

### Output

| Output | Lokasi | Isi |
|---|---|---|
| **Persona Directory** | `.ai-doc/personas/<nama-persona>/` | Folder persona hasil generate |
| **Persona File** | `.ai-doc/personas/<nama-persona>/persona.md` | Detail persona: karakteristik, cara berpikir, komunikasi; bagian Load Config hanya merujuk ke `customize.toml` |
| **Customize File** | `.ai-doc/personas/<nama-persona>/customize.toml` | Source of truth konfigurasi persona untuk agent (metadata, identity, nama, role, prinsip) |
| **Methods File** | `.ai-doc/personas/<nama-persona>/methods.md` | Detail lengkap metode yang dikuasai (Brain, Solving, Innovation) |

---

## Flow Eksekusi

Gunakan 2 step berikut secara berurutan:

```text
Step 1: Select Blueprint
  ├── Agent tampilkan daftar blueprint yang tersedia
  ├── Agent jelaskan perbedaan masing-masing blueprint
  ├── User pilih blueprint
  ├── User tentukan nama persona
  └── HALT → konfirmasi pilihan

Step 2: Generate Persona
  ├── Buat folder .ai-doc/personas/<nama-persona>/
  ├── Generate customize.toml dari blueprint sebagai source of truth
  ├── Generate persona.md dari blueprint + referensi Load Config ke customize.toml
  ├── Generate methods.md dari referensi add-on method
  ├── Update .ai-doc/personas/list.md
  └── Selesai → persona siap dipakai
```

---

## Daftar Blueprint

| # | Blueprint | Archetype | Cocok Untuk |
|---|---|---|---|
| 1 | **Problem Solver** 🔬 | The Diagnostician | Debugging, root cause analysis, TRIZ, problem-solving sistematis |
| 2 | **Creative Visionary** ✨ | The Innovator | Ideation, design thinking, blue ocean, product discovery |
| 3 | **Technical Architect** 🏗️ | The System Builder | Architecture review, platform design, technology roadmap |
| 4 | **Data Specialist** 📊 | The Insight Miner | Data exploration, KPI review, experiment analysis, insight extraction |
| 5 | **Senior Backend Engineer** ⚙️ | The Reliability Builder | API design, production debugging, service reliability, backend scaling |
| 6 | **Custom** 📋 | — | Anda tentukan sendiri dari awal |

Lihat detail masing-masing blueprint di `add-on/persona/blueprints/`.

---

## Cara Pemilihan

1. Agent tanya kebutuhan user: "Untuk keperluan apa persona ini?"
2. Agent rekomendasikan blueprint yang paling sesuai berdasarkan kebutuhan.
3. User pilih blueprint (boleh juga Custom).
4. User tentukan nama untuk persona-nya.
5. Agent konfirmasi sebelum generate.

---

## Integrasi dengan Add-On Lain

Persona yang sudah dibuat bisa digunakan oleh add-on lain:

- **Brainstorming Add-On** — persona bisa di-load sebagai peserta sesi brainstorming
- **Add-On lainnya** — persona sebagai karakter fasilitator/kolaborator

Saat add-on lain membutuhkan persona, ia akan membaca:
1. `.ai-doc/personas/list.md` — daftar persona yang tersedia
2. `.ai-doc/personas/<nama>/customize.toml` — source of truth konfigurasi agent/persona
3. `.ai-doc/personas/<nama>/persona.md` — detail karakter & komunikasi
4. `.ai-doc/personas/<nama>/methods.md` — detail metode yang dikuasai

---

## Referensi File

| File | Kegunaan |
|---|---|
| `add-on/persona/blueprints/problem-solver/` | Blueprint Problem Solver (Dr. Quinn) |
| `add-on/persona/blueprints/creative-visionary/` | Blueprint Creative Visionary (Nova) |
| `add-on/persona/blueprints/technical-architect/` | Blueprint Technical Architect (Arch) |
| `add-on/persona/blueprints/data-specialist/` | Blueprint Data Specialist (Mira) |
| `add-on/persona/blueprints/senior-backend-engineer/` | Blueprint Senior Backend Engineer (Raka) |
| `add-on/persona/blueprints/custom/` | Blueprint kosong untuk persona kustom |
| `add-on/persona/steps/step-01-select.md` | Detail langkah pemilihan blueprint |
| `add-on/persona/steps/step-02-generate.md` | Detail langkah generate persona |
