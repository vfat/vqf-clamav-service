# AI Documentor — Add-On

> Referensi entry point untuk mengaktifkan **add-on** secara eksplisit.
> Add-on tidak aktif secara default — hanya berjalan jika user minta.
> Format ini mengikuti struktur yang sama dengan `skill-help.md` untuk konsistensi.

---

## Daftar Add-On

| # | Menu Code | Display Name | Deskripsi | Trigger User | Entry File | Output Location | Outputs |
|---|-----------|-------------|-----------|-------------|-----------|-----------------|---------|
| 1 | `AD-AO-BRAIN` | Brainstorming | Memfasilitasi sesi brainstorming terstruktur dengan output MoM + Discussion Summary. Mencakup 2 area (Troubleshooting & Feature/Release), 18 teknik, idea status tracking, dan integrasi persona. | "brainstorm", "brainstrom", "diskusi ide", "ideation", "sesi brainstorming" | `add-on/brainstorming/workflow.md` | `.ai-doc/brainstorming/` | `mom-{date}-{topic}.md`, `discussion-{subtopic}-{date}.md` |
| 2 | `AD-AO-PERSONA` | Persona | Menyiapkan template persona yang bisa dipakai berkali-kali. Pilih blueprint (Problem Solver, Creative Visionary, Technical Architect, Data Specialist, Senior Backend Engineer, Custom), beri nama, dan generate ke `.ai-doc/personas/`. | "add on persona", "buatkan persona", "template persona", "persona" | `add-on/persona/workflow.md` | `.ai-doc/personas/<nama-persona>/` | `persona.md`, `customize.toml` |
| 3 | `AD-AO-METHOD` | Method | Koleksi 100 teknik/metode/framework lintas 3 kategori: Brain Methods (45 teknik ideation), Solving Methods (25 metode problem-solving), Innovation Frameworks (30 framework strategi inovasi). Dilengkapi techniques-index untuk navigasi. | "brain methods", "solving method", "innovation framework", "teknik brainstorming", "cari metode", "metode problem-solving", "framework inovasi" | `add-on/method/techniques-index.md` | in-memory (opsional: `.ai-doc/method/`) | output sesuai metode |
| 4 | `AD-AO-QASTRAT` | QA Strategy | Menyusun test strategy dan test plan berbasis bukti dari artefak `.ai-doc/` yang ada: risk matrix, test pyramid, quality gates, entry/exit criteria, lalu turunan plan per sprint/release. Adaptasi pola qa-skills ke konvensi AI Documentor. | "test strategy", "qa strategy", "rencana testing", "qa plan", "test plan" | `add-on/qa-strategy/workflow.md` | `.ai-doc/qa/strategy/` | `test-strategy.md`, `plans/test-plan-{slug}-{date}.md` |
| 5 | `AD-AO-QASEC` | QA Security | Security screening & audit berbasis bukti: pemetaan attack surface (brownfield-friendly — rekonstruksi dari codebase bila artefak belum ada), coverage matrix OWASP Top 10 2025 tanpa kategori hilang diam-diam, findings register dengan keputusan Remediate/Accept/Defer. Add-on tidak menjalankan scanner; tooling tetap di qa-skills. | "security screening", "audit keamanan", "owasp", "security review", "coverage keamanan" | `add-on/qa-security/workflow.md` | `.ai-doc/qa/security/` | `security-coverage.md` |
| 6 | `AD-AO-QAEXEC` | QA Execution | Mengeksekusi item dari test plan yang valid (wajib ada, gate keras) dalam mode manual / automated / hybrid — agent menjalankan command test nyata (pytest, Playwright, Jest, go test) dan mencatat exit code + output sebagai bukti. Status model PLANNED → IN PROGRESS → PASSED/FAILED/BLOCKED/DEFERRED dengan disiplin bukti ala TDD; control plane `.ai-doc/qa/execution/qa-overview.md`. Item unit-level saat TDD Enabled dieksekusi via add-on tdd. | "jalankan test plan", "eksekusi testing", "qa run", "run test plan", "eksekusi test plan" | `add-on/qa-execution/workflow.md` | `.ai-doc/qa/execution/` | `qa-overview.md`, update §7 test plan |

---

## Alur Aktivasi

```text
User menyebut trigger
        │
        ▼
Deteksi trigger add-on
        │
        ▼
Buka add-on/addon-help.md → jalankan workflow
        │
        ▼
Jalankan sesi (dengan HALT pattern)
        │
        ▼
Selesai → return kontrol ke core AI Documentor
```

Saat trigger add-on terdeteksi:

1. Buka `add-on/addon-help.md` untuk melihat index dan detail add-on.
2. Jalankan workflow sesuai entry file di tabel.
3. Output mengikuti aturan AI Documentor (evidence-based, no over-infer).

---

## Aturan

- **Jangan auto-run add-on.** Hanya jalankan jika user minta eksplisit.
- **Setelah add-on selesai**, return kontrol ke core AI Documentor.
- **Jika user tidak menyebut add-on**, jalankan AI Documentor normal (`SKILL.md`).
- `SKILL.md` tetap bersih — tidak perlu diedit setiap kali add-on ditambah.
- Cukup update `skill-help.md` dan/atau `references/add-on.md` untuk menambahkan menu baru.

---

## Referensi Terkait

| Menu Code | File Referensi |
|-----------|---------------|
| `AD-AO-BRAIN` | `add-on/addon-help.md` |
| `AD-AO-BRAIN` | `add-on/brainstorming/workflow.md` |
| `AD-AO-BRAIN` | `add-on/brainstorming/techniques.md` |
| `AD-AO-BRAIN` | `add-on/brainstorming/mode-moderator.md` |
| `AD-AO-PERSONA` | `add-on/addon-help.md` |
| `AD-AO-PERSONA` | `add-on/persona/workflow.md` |
| `AD-AO-PERSONA` | `add-on/persona/blueprints/problem-solver/persona.md` |
| `AD-AO-PERSONA` | `add-on/persona/blueprints/creative-visionary/persona.md` |
| `AD-AO-PERSONA` | `add-on/persona/blueprints/technical-architect/persona.md` |
| `AD-AO-PERSONA` | `add-on/persona/blueprints/data-specialist/persona.md` |
| `AD-AO-PERSONA` | `add-on/persona/blueprints/senior-backend-engineer/persona.md` |
| `AD-AO-PERSONA` | `add-on/persona/blueprints/custom/persona.md` |
| `AD-AO-METHOD` | `add-on/addon-help.md` |
| `AD-AO-METHOD` | `add-on/method/techniques-index.md` |
| `AD-AO-METHOD` | `add-on/method/brain/README.md` |
| `AD-AO-METHOD` | `add-on/method/solving/README.md` |
| `AD-AO-METHOD` | `add-on/method/innovation/README.md` |
| `AD-AO-QASTRAT` | `add-on/addon-help.md` |
| `AD-AO-QASTRAT` | `add-on/qa-strategy/workflow.md` |
| `AD-AO-QASTRAT` | `add-on/qa-strategy/steps/step-01-context.md` |
| `AD-AO-QASTRAT` | `add-on/qa-strategy/steps/step-02-strategy.md` |
| `AD-AO-QASTRAT` | `add-on/qa-strategy/steps/step-03-plan.md` |
| `AD-AO-QASTRAT` | `add-on/qa-strategy/template/test-strategy-template.md` |
| `AD-AO-QASTRAT` | `add-on/qa-strategy/template/test-plan-template.md` |
| `AD-AO-QASEC` | `add-on/addon-help.md` |
| `AD-AO-QASEC` | `add-on/qa-security/workflow.md` |
| `AD-AO-QASEC` | `add-on/qa-security/steps/step-01-surface.md` |
| `AD-AO-QASEC` | `add-on/qa-security/steps/step-02-assessment.md` |
| `AD-AO-QASEC` | `add-on/qa-security/steps/step-03-findings.md` |
| `AD-AO-QASEC` | `add-on/qa-security/template/security-coverage-template.md` |
| `AD-AO-QAEXEC` | `add-on/addon-help.md` |
| `AD-AO-QAEXEC` | `add-on/qa-execution/workflow.md` |
| `AD-AO-QAEXEC` | `add-on/qa-execution/steps/step-01-activation.md` |
| `AD-AO-QAEXEC` | `add-on/qa-execution/steps/step-02-execution.md` |
| `AD-AO-QAEXEC` | `add-on/qa-execution/steps/step-03-wrap-up.md` |
| `AD-AO-QAEXEC` | `add-on/qa-execution/template/qa-overview-template.md` |
| — | `tes/07-add-on.md` (arsitektur add-on) |
