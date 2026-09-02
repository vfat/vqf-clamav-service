# Brainstorming Add-On — Workflow

> Orchestrator / entry point untuk add-on brainstorming AI Documentor.
> Panggil workflow ini ketika user meminta sesi brainstorming.

---

## Overview

```
Input: Topik + Konteks + Teknik (opsional)
        │
        ▼
┌───────────────────────────────┐
│    Sesi Brainstorming         │
│  (interaktif, facilitator-    │
│   driven, HALT pattern)       │
│          │                    │
│  ┌───────┴───────┐            │
│  ▼               ▼            │
│ MoM          Discussion       │
│ (ringkasan)   Summary per     │
│               Sub-Topik       │
└───────────────────────────────┘
```

### 2 Output Utama

| Output | Format | Isi | Kegunaan |
|---|---|---|---|
| **MoM** | Markdown | Metadata sesi, daftar topik, ringkasan keputusan, action items | Dokumentasi meeting, share ke stakeholder |
| **Discussion Summary** | Markdown per topik | Konteks, ide-ide yang muncul, keputusan/pembahasan, next steps per topik | Deep-dive per topik, referensi follow-up |

---

## Flow Eksekusi

Gunakan 3 step berikut secara berurutan:

```text
Step 1: Setup
  ├── Tentukan mode moderator (Eksploratif / Analitis-Evidence / Pengambil Keputusan Aplikatif)
  ├── Cek .ai-doc/personas/list.md (jika ada)
  ├── Tentukan area (Troubleshooting / Feature / Release)
  ├── Breakdown sub-topik + rekomendasi teknik
  ├── Konfirmasi peserta (user-only / +Joni / +Jono / keduanya)
  ├── Mode moderator final: [mode terpilih] — tidak berubah selama sesi
  └── HALT → konfirmasi setup

Step 2: Facilitation (Loop per Sub-Topik)
  ├── Tampilkan sub-topik + teknik
  ├── Lempar prompt sesuai teknik
  ├── Persona respons / agent jawab sendiri
  ├── User tanggapi (approve / reject / modify)
  ├── Catat ide (append-only + status)
  └── HALT → "Lanjut / Ganti sub-topik / Selesai?"

Step 3: Wrap-Up
  ├── Review semua ide
  ├── Konfirmasi keputusan & action items
  ├── Generate MoM → `.ai-doc/brainstorming/mom-{date}-{topic-slug}.md`
  └── Generate Discussion Summary → `.ai-doc/brainstorming/discussion-{subtopic-slug}-{date}.md`
```

---

## Mode Moderator

Mode moderator ditentukan di **Step 1 (Setup)** dan **tidak berubah selama sesi berlangsung**.

### Tiga Mode dengan Karakter Fasilitator

| Mode | Nama | Fokus | Cocok Untuk |
|---|---|---|---|
| **Eksploratif** | Sherin ✋ | Membuka ruang ide, memperbanyak kemungkinan | Greenfield, project awal, fitur baru, konteks belum jelas |
| **Analitis-Evidence** | Manda 🔍 | Memvalidasi ide, melacak bukti | Brownfield, troubleshooting, root cause, banyak klaim |
| **Pengambil Keputusan Aplikatif** | Dinda 🎯 | Mengerucutkan ke opsi realistis | Ide sudah banyak, perlu prioritas & aksi nyata |

### Cara Pemilihan

1. Agent menganalisis topik user dan **merekomendasikan** mode yang paling sesuai dengan alasan.
2. Agent jelaskan 3 mode ke user beserta fokus dan cocok untuk situasi apa — serta perkenalkan nama fasilitator (Sherin, Manda, Dinda).
3. User pilih salah satu.
4. Agent tampilkan **Greet the User** sesuai mode yang dipilih (sapaan hangat dari karakter fasilitator).
5. Mode dicatat di session state dan dipakai sampai sesi selesai.

Lihat detail masing-masing mode + greet dialog di `add-on/brainstorming/mode-moderator.md`.

---

## Aturan Facilitation

1. **Agent = facilitator**, bukan generator ide (kecuali tidak ada persona aktif).
2. **Satu elemen teknik per interaksi** — jangan lempar banyak prompt sekaligus.
3. **HALT di setiap titik keputusan** — jangan auto-pilot.
4. **Anti-bias** — pivot domain setiap ~10 ide.
5. **Default: keep exploring** — user yang tentukan kapan berhenti.
6. **Persona respons sesuai karakteristiknya** (lihat `.ai-doc/personas/list.md`).
7. **Idea status tracking** — setiap ide punya status: `approved` / `rejected` / `modified`.
8. **Ide yang di-reject tetap dicatat** untuk audit trail.

---

## 2 Area Fokus

| Area | Deskripsi | Kategori Teknik |
|---|---|---|
| **Troubleshooting / Bug Fix / Hotfix / Improvement** | Memperbaiki yang sudah ada — dari diagnosis sampai solusi inovatif | Diagnosis → Analysis → Solution → Innovation |
| **Feature / Release** | Membuat yang baru — dari ideasi sampai perencanaan release | Ideation → Perspective → Planning → Innovation |

Lihat `add-on/brainstorming/techniques.md` untuk detail teknik per area.

---

## Integrasi Persona

Jika folder `.ai-doc/personas/` ada dan berisi `list.md`, agent harus:

1. Load daftar persona yang tersedia.
2. Tawarkan ke user pilihan peserta:
   - User only
   - User + Joni (persona)
   - User + Jono (persona)
   - User + Joni + Jono
3. Persona yang aktif akan merespons sesuai karakteristiknya selama sesi.

Jika tidak ada persona, agent bertindak sebagai fasilitator dan generator ide.

---

## Output Location

Semua output disimpan di:

```
.ai-doc/brainstorming/
├── mom-{YYYY-MM-DD}-{topic-slug}.md
└── discussion-{subtopic-slug}-{YYYY-MM-DD}.md
```

Format slug: lowercase, gunakan `-` sebagai separator, maksimal 5 kata.

---

## Referensi File

| File | Kegunaan |
|---|---|
| `add-on/brainstorming/techniques.md` | Daftar teknik brainstorming per area & kategori |
| `add-on/brainstorming/mode-moderator.md` | Detail 3 mode moderator (Eksploratif, Analitis-Evidence, Pengambil Keputusan Aplikatif) |
| `add-on/brainstorming/template/mom-template.md` | Template output MoM |
| `add-on/brainstorming/template/discussion-template.md` | Template output Discussion Summary |
| `add-on/brainstorming/steps/step-01-setup.md` | Detail langkah setup sesi |
| `add-on/brainstorming/steps/step-02-facilitation.md` | Detail langkah facilitation loop |
| `add-on/brainstorming/steps/step-03-wrap-up.md` | Detail langkah wrap-up & generate output |
| `tes/01-brainstrom-addon.md` | Dokumen desain awal add-on brainstorming |
