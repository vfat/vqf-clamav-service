# AI Documentor — Skill Help

Gunakan saat mendokumentasikan codebase yang sudah ada atau merencanakan sistem baru melalui workflow `.ai-doc/` berbasis bukti. Mencakup greenfield project overview, brownfield codebase docs, feature docs, grouped use case docs, DCD, database design docs, REST API docs, synchronization pass, dan documentation review.

---

## Daftar Sub-Skill

| # | Menu Code | Display Name | Deskripsi | Phase | Required | Output Location | Outputs |
|---|-----------|-------------|-----------|-------|----------|-----------------|---------|
| — | `AD` | AI Documentor | Gunakan saat mendokumentasikan codebase yang sudah ada atau merencanakan sistem baru melalui workflow `.ai-doc/` berbasis bukti. Mencakup greenfield project overview, brownfield codebase docs, feature docs, grouped use case docs, DCD, database design docs, REST API docs, synchronization pass, dan documentation review. | anytime | false | `.ai-doc/` | berbagai artefak di bawah `.ai-doc/` |
| 1 | `AD-FA` | First Actions | Bootstrap workspace `.ai-doc/`, baca file control plane (`3p.md`, `constitution.md`), tentukan artefak target, audit bukti kode, dan terapkan aturan only-on-request sebelum menulis dokumentasi apa pun. | bootstrap | true | `.ai-doc/3p.md` | `3p.md` diupdate |
| 2 | `AD-AM` | Artifact Map | Pilih dokumen target dan template yang sesuai berdasarkan request user. Memetakan jenis request ke file output dan template (contoh: `project-overview.md`, DCD, ERD, REST API doc). | anytime | true | in-memory selection | artefak target + template terpilih |
| 3 | `AD-GF` | Greenfield Flow | Rencanakan project atau sistem baru dari nol. Membuat `.ai-doc/project-overview.md` melalui tanya jawab terstruktur dengan user. Mencakup problem statement, stakeholder, scope, constraint, dan arah solusi high-level. | planning | false | `.ai-doc/project-overview.md` | `project-overview.md` |
| 4 | `AD-BF` | Brownfield Flow | Dokumentasikan codebase yang sudah ada. Rekonstruksi kondisi saat ini dari source code nyata (entry point, config, service, schema). Pilih artefak paling sempit yang tepat dan tulis secara konservatif dengan pelacakan bukti. | documentation | false | `.ai-doc/Dokumentasi-Codebase.md` | `Dokumentasi-Codebase.md` + dokumen terkait |
| 5 | `AD-CP` | Component Planning Flow | Rencanakan komponen baru. Membuat Spec Component Document (SCD) di `.ai-doc/plan/component/` melalui tanya jawab. Menangkap context, scope, prerequisite, dan daftar use case awal. | planning | false | `.ai-doc/plan/component/SCD-<nama>.md` | `SCD-<nama>.md` |
| 6 | `AD-CON` | General Constitution | Terapkan konstitusi dokumentasi berbasis bukti default ke sebuah project. Mendefinisikan aturan workspace, progress tracking, aturan bukti, klasifikasi arsitektur, dan aturan identifikasi objek. | bootstrap | true | `.ai-doc/constitution.md` | `constitution.md` |
| 7 | `AD-DIAG` | Diagram Rule | Aturan pembuatan diagram C4 component dan contoh baseline Mermaid. Mencakup aturan penamaan komponen, relationship evidence-based, dan format diagram. | anytime | false | `references/diagram-rule.md` | — |
| 8 | `AD-GUC` | Grouped Use Case Rule | Aturan penulisan dokumen grouped use case. Mencakup struktur section standar, konsistensi nama komponen/use case, dan aturan paragraf deskripsi. | anytime | false | `references/grouped-usecase-rule.md` | — |
| 9 | `AD-REV` | Review Mode | Alur verifikasi dokumen: buka dokumen target, bandingkan claim dengan kode, turunkan wording bila evidence lemah, sinkronkan dokumen tetangga. | anytime | false | `references/review-mode.md` | — |
| 10 | `AD-OD` | Output Discipline | Aturan output: pilih artefak terkecil, jangan buat dokumen turunan tanpa permintaan eksplisit, prefer exact names dari repo. | anytime | false | `references/output-discipline.md` | — |
| 11 | `AD-AO-BRAIN` | Brainstorming Add-On | Memfasilitasi sesi brainstorming terstruktur dengan output MoM + Discussion Summary. Hanya aktif jika diminta eksplisit ("brainstorm", "diskusi ide"). | add-on | false | `.ai-doc/brainstorming/` | `mom-{date}-{topic}.md`, `discussion-{subtopic}-{date}.md` |
| 12 | `AD-AO-PERSONA` | Persona Add-On | Menyiapkan template persona yang bisa dipakai berkali-kali. Pilih blueprint (Problem Solver, Creative Visionary, Technical Architect, Data Specialist, Senior Backend Engineer, Custom), beri nama, generate ke `.ai-doc/personas/`. | add-on | false | `.ai-doc/personas/<nama-persona>/` | `persona.md`, `customize.toml` |
| 13 | `AD-AO-METHOD` | Method Add-On | Koleksi 100 teknik/metode/framework: 45 Brain Methods (ideation), 25 Solving Methods (problem-solving), 30 Innovation Frameworks (strategi inovasi). Dilengkapi techniques-index dengan persona mapping dan quick-select matrix. | add-on | false | in-memory (opsional: `.ai-doc/method/`) | output sesuai metode |
| 14 | `AD-AO-TDD` | TDD Add-On | Add-on mandiri untuk greenfield: keputusan eksplisit di constitution, test gagal sebelum production code, dan tracking RED/GREEN/REFACTOR di `.ai-doc/tdd-overview.md`. | add-on | false | `.ai-doc/tdd-overview.md` | `tdd-overview.md` |
| 15 | `AD-AO-QASTRAT` | QA Strategy Add-On | Menyusun test strategy (risk matrix, test pyramid, quality gates, entry/exit criteria) dan opsional test plan per sprint/release — semuanya berbasis bukti dari artefak `.ai-doc/` yang ada. Hanya aktif jika diminta eksplisit. | add-on | false | `.ai-doc/qa/strategy/` | `test-strategy.md`, `plans/test-plan-{slug}-{date}.md` |
| 16 | `AD-AO-QASEC` | QA Security Add-On | Security screening & audit berbasis bukti: attack surface mapping (brownfield-friendly), coverage matrix OWASP Top 10 2025, findings register dengan keputusan Remediate/Accept/Defer. Tidak menjalankan scanner — hanya memetakan & mendokumentasikan. Hanya aktif jika diminta eksplisit. | add-on | false | `.ai-doc/qa/security/` | `security-coverage.md` |
| 17 | `AD-AO-QAEXEC` | QA Execution Add-On | Mengeksekusi item test plan (manual/automated/hybrid) dengan status berbasis bukti nyata — command + exit code + output. Control plane `.ai-doc/qa/execution/qa-overview.md`. Gate keras: butuh test plan valid. Hanya aktif jika diminta eksplisit. | add-on | false | `.ai-doc/qa/execution/` | `qa-overview.md`, update §7 test plan |

---

## Alur Eksekusi

```text
AD-FA (First Actions)  →  AD-CON (General Constitution)
        ↓
AD-AM (Artifact Map)  →  Pilih sub-skill sesuai request:
                           ├── AD-GF (Greenfield)
                           ├── AD-BF (Brownfield)
                           ├── AD-CP (Component Planning)
                           └── AD-REV (Review Mode)
```

---

## Referensi Aturan

| Menu Code | File Referensi |
|-----------|---------------|
| `AD-DIAG` | `references/diagram-rule.md` |
| `AD-GUC` | `references/grouped-usecase-rule.md` |
| `AD-REV` | `references/review-mode.md` |
| `AD-OD` | `references/output-discipline.md` |
| `AD-AO-BRAIN` | `add-on/addon-help.md` |
| `AD-AO-PERSONA` | `add-on/addon-help.md` |
| `AD-AO-METHOD` | `add-on/addon-help.md` |
| `AD-AO-METHOD` | `add-on/method/techniques-index.md` |
| `AD-AO-TDD` | `add-on/addon-help.md` |
| `AD-AO-QASTRAT` | `add-on/addon-help.md` |
| `AD-AO-QASTRAT` | `add-on/qa-strategy/workflow.md` |
| `AD-AO-QASEC` | `add-on/addon-help.md` |
| `AD-AO-QASEC` | `add-on/qa-security/workflow.md` |
| `AD-AO-QAEXEC` | `add-on/addon-help.md` |
| `AD-AO-QAEXEC` | `add-on/qa-execution/workflow.md` |
