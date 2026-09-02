# Step 1: Select Blueprint

> User memilih blueprint persona dan menentukan nama.
> **HALT** — tunggu konfirmasi user sebelum generate.

---

## Input dari User

- **Kebutuhan** — untuk apa persona ini? (brainstorming, problem-solving, architecture review, dll.)
- **Nama persona** — user memberi nama untuk persona-nya

---

## Yang Dilakukan Agent

### 1. Tanya Kebutuhan

Tanya user untuk apa persona akan digunakan. Berdasarkan jawaban, rekomendasikan blueprint:

| Jika User Butuh | Rekomendasi Blueprint | Alasan |
|---|---|---|
| Memecahkan masalah kompleks, debugging, root cause analysis | **Problem Solver** 🔬 | Dibekali TRIZ, Systems Thinking, Failure Mode Analysis |
| Ideation, product discovery, inovasi, eksplorasi pasar | **Creative Visionary** ✨ | Dibekali Design Thinking, Blue Ocean, SCAMPER |
| Architecture review, technology roadmap, platform design | **Technical Architect** 🏗️ | Dibekali Decision Matrix, Tech Roadmapping, Platform Design |
| Analisis data, KPI review, experiment analysis, insight extraction | **Data Specialist** 📊 | Dibekali metric framing, gap analysis, dan market/data insight |
| API design, backend troubleshooting, service reliability, scaling | **Senior Backend Engineer** ⚙️ | Dibekali failure analysis, systems thinking, dan roadmap teknologi |
| Tidak cocok dengan 5 di atas / ingin dari awal | **Custom** 📋 | Blueprint kosong, isi sendiri |

**Cara penyampaian:**

> Baik, untuk keperluan **[kebutuhan user]** , saya rekomendasikan blueprint **[nama blueprint]** karena **[alasan]**. Tapi Anda bisa pilih yang lain:
>
> 1. **Problem Solver** 🔬 — untuk problem-solving sistematis
> 2. **Creative Visionary** ✨ — untuk ideation dan inovasi
> 3. **Technical Architect** 🏗️ — untuk perancangan sistem
> 4. **Data Specialist** 📊 — untuk analisis data dan insight
> 5. **Senior Backend Engineer** ⚙️ — untuk engineering backend produksi
> 6. **Custom** 📋 — isi sendiri dari awal
>
> Blueprint mana yang Anda pilih?

**HALT** — tunggu user memilih blueprint.

### 2. Tentukan Nama Persona

Setelah blueprint dipilih, minta user memberi nama persona:

> Nama untuk persona ini? Misalnya: "Dr. Quinn", "Joni", "Sinta", atau nama lain yang Anda suka.

**Catatan:**
- Nama akan jadi folder name → gunakan format slug: lowercase, tanpa spasi, gunakan `-` sebagai separator.
- Contoh: `dr-quinn`, `joni`, `sinta-the-architect`.

### 3. Tampilkan Pratinjau

Tampilkan pratinjau blueprint yang dipilih:

```yaml
blueprint: "problem-solver"
persona_name: "dr-quinn"
title: "Master Problem Solver"
icon: "🔬"
role: "Crack complex challenges with systematic problem-solving methodologies..."
description: "Apply systematic problem-solving methodologies to a hard challenge"
methods:
  brain: ["First Principles Thinking", "Solution Matrix", "Role Playing"]
  solving: ["Systems Thinking", "Failure Mode Analysis", "TRIZ Contradiction Matrix", "Feasibility Study"]
  innovation: ["Business Model Patterns", "Jobs to be Done"]
```

### 4. Konfirmasi

**HALT** — tanya user:

> Siap saya generate persona **{{nama}}** dengan blueprint **{{blueprint}}**?
> - ✅ **Generate** — lanjut ke Step 2
> - ✏️ **Ubah nama** — ganti nama persona
> - 🔄 **Ganti blueprint** — pilih blueprint lain
> - ❌ **Batal** — batalkan pembuatan persona

---

## Output Step 1

```yaml
persona:
  blueprint: "problem-solver" | "creative-visionary" | "technical-architect" | "data-specialist" | "senior-backend-engineer" | "custom"
  persona_name: "dr-quinn"        # slug, jadi folder name
  persona_label: "Dr. Quinn"      # display name
  kebutuhan: "..."                # catatan kebutuhan user
```

---

## Aturan

- **JANGAN auto-generate** sebelum user konfirmasi.
- **JANGAN paksakan blueprint** — user boleh pilih blueprint yang tidak sesuai rekomendasi.
- **JANGAN gunakan nama yang sudah ada** di `.ai-doc/personas/list.md` — beri tahu user bahwa nama sudah dipakai.
- **REKOMENDASIKAN** blueprint — saran diterima user yang putuskan.
- Nama folder harus **slug** (lowercase, `-` sebagai separator).
