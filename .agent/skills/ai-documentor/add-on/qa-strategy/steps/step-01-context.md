# Step 1: Context Gathering

> Kumpulkan bukti dari artefak `.ai-doc/` yang sudah ada sebelum menyusun strategi apa pun.
> **HALT** — tunggu konfirmasi user sebelum lanjut ke Step 2.

---

## Input dari User

- **Cakupan strategi** — seluruh sistem / satu komponen / satu release (wajib)
- **Target sprint atau release** — opsional di step ini, wajib jika lanjut ke Step 3
- **Konteks tambahan** — incident record, keluhan user, area yang sering rusak

---

## Yang Dilakukan Agent

### 1. Baca Control Plane

Baca dua file berikut sebelum hal lain:

- `.ai-doc/3p.md` — status project dan langkah terakhir
- `.ai-doc/constitution.md` — aturan lokal yang mungkin membatasi strategi
  (misal aturan bahasa, klasifikasi arsitektur, kebijakan TDD)

### 2. Inventarisasi Artefak Sumber

Cek artefak berikut dan catat mana yang ada:

| Artefak | Yang diekstrak | Jika tidak ada |
|---|---|---|
| `project-overview.md` | Critical journeys, constraint | Tandai `Perlu Dikonfirmasi`, tanya user |
| `Dokumentasi-Codebase.md` | Komponen nyata, entry point | Boleh lanjut tanpa ini, catat gap |
| `DCD-*.md` | Use case per komponen | Coverage jadi kasar, tandai `Partial` |
| ERD / data dictionary | Entitas kritikal | Risiko data tidak terpetakan |
| `rest-api-doc` | Endpoint publik | API test tidak bisa diprioritaskan |

### 3. Ekstraksi Bahan Strategi

Susun daftar kerja (in-memory, belum ditulis ke file):

```yaml
evidence_inventory:
  components:
    - name: "<nama komponen>"
      source: "<artefakt asal>"
      use_cases: ["<UC-001>", "..."]
  public_endpoints:
    - path: "<endpoint>"
      source: "rest-api-doc"
  known_risks:
    - description: "<risiko>"
      source: "<bukti: incident/churn/integrasi>"
      confidence: "evidence" | "asumsi"
  gaps:
    - "<area tanpa bukti>"
```

### 4. Presentasikan Ringkasan & Konfirmasi Cakupan

> Dari artefak yang ada, saya menemukan **[N] komponen**, **[M] endpoint publik**,
> dan **[K] risiko tercatat**. Area tanpa bukti: [daftar gap].
>
> Usulan cakupan strategi: **[usulan]** karena [alasan].
> - ✅ **Lanjut dengan cakupan ini**
> - ✏️ **Ubah cakupan** — misal fokus ke satu komponen saja
> - ➕ **Lengkapi bukti dulu** — buat artefak yang hilang via core AI Documentor

**HALT** — tunggu keputusan user.

---

## Output Step 1

Session state (hold in memory):

```yaml
session_state:
  scope: "system" | "component" | "release"
  scope_target: "<nama komponen/release jika relevan>"
  evidence_inventory:
    components_count: 0
    endpoints_count: 0
    risks: [...]
    gaps: [...]
  artifacts_available:
    project_overview: true/false
    codebase_doc: true/false
    dcd: true/false
    erd: true/false
    rest_api_doc: true/false
```

---

## Aturan

- **JANGAN menyusun strategi dari nama file saja** — baca isi artefaknya.
- **JANGAN mengisi gap dengan asumsi diam-diam** — tandai dan laporkan.
- **REKOMENDASIKAN cakupan** — saran dengan alasan, keputusan milik user.
- Jika artefak inti belum ada sama sekali, sarankan kembali ke core flow
  (`AD-GF`/`AD-BF`) sebelum add-on ini berguna.
