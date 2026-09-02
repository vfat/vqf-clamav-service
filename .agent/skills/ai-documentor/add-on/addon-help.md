# Add-On Help — AI Documentor

> Index deklaratif semua add-on yang tersedia di AI Documentor.
> Add-on hanya aktif jika diminta **eksplisit** oleh user.

---

## Daftar Add-On

| Add-On | Status | Trigger User | Tujuan | Entry File | Supports Persona | Output Location |
|---|---|---|---|---|---|---|
| brainstorming | `active` | "brainstorm", "brainstrom", "diskusi ide", "ideation", "sesi brainstorming" | Memfasilitasi sesi brainstorming terstruktur dengan 2 output: MoM + Discussion Summary | `add-on/brainstorming/workflow.md` | single / multi | `.ai-doc/brainstorming/` |
| persona | `active` | "add on persona", "buatkan persona", "template persona", "persona" | Menyiapkan template persona yang bisa dipakai berkali-kali oleh user — pilih blueprint, generate ke `.ai-doc/personas/` | `add-on/persona/workflow.md` | — | `.ai-doc/personas/<nama-persona>/` |
| method | `active` | "brain methods", "solving method", "innovation framework", "teknik brainstorming", "metode problem-solving", "framework inovasi", "cari metode" | Koleksi 100 teknik/metode/framework: 45 Brain Methods (ideation), 25 Solving Methods (problem-solving), 30 Innovation Frameworks (strategi bisnis). | `add-on/method/techniques-index.md` | — | in-memory (opsional: `.ai-doc/method/`) |
| tdd | `active` | "tdd", "test-driven", "gunakan TDD", "aktifkan TDD" | Menjalankan siklus RED → GREEN → REFACTOR secara eksplisit untuk greenfield dan melacak status target di `.ai-doc/tdd-overview.md`. | `add-on/tdd/SKILL.md` | — | `.ai-doc/tdd-overview.md` |
| qa-strategy | `active` | "test strategy", "qa strategy", "rencana testing", "qa plan", "test plan" | Menyusun test strategy (risk matrix, pyramid, quality gates) dan opsional test plan per sprint/release, semuanya berbasis bukti dari artefak `.ai-doc/`. | `add-on/qa-strategy/workflow.md` | single / multi | `.ai-doc/qa/strategy/` |
| qa-security | `active` | "security screening", "audit keamanan", "owasp", "security review", "coverage keamanan" | Security screening & audit berbasis bukti: attack surface mapping (brownfield-friendly), coverage matrix OWASP Top 10 2025, findings register dengan keputusan Remediate/Accept/Defer. Tidak menjalankan scanner — tooling tetap di qa-skills. | `add-on/qa-security/workflow.md` | single / multi | `.ai-doc/qa/security/` |
| qa-execution | `active` | "jalankan test plan", "eksekusi testing", "qa run", "run test plan", "eksekusi test plan" | Mengeksekusi item test plan (manual/automated/hybrid) dengan status berbasis bukti nyata: command + exit code + output. Control plane `.ai-doc/qa/execution/qa-overview.md`. Gate keras: butuh test plan valid dari qa-strategy. | `add-on/qa-execution/workflow.md` | single / multi | `.ai-doc/qa/execution/` |

---

## Detail Add-On

### brainstorming

| Atribut | Nilai |
|---|---|
| **Status** | `active` — siap digunakan |
| **Tujuan** | Memfasilitasi sesi brainstorming interaktif yang menghasilkan output terstruktur — bukan cuma ide mengalir, tapi terdokumentasi rapi dalam format yang langsung dipakai |
| **Kapan digunakan** | User ingin brainstorming topik tertentu (troubleshooting, feature, release) dengan output MoM dan/atau Discussion Summary |
| **Kapan tidak digunakan** | User hanya ingin diskusi ringan tanpa output terdokumentasi; user hanya ingin dokumentasi codebase biasa (gunakan core AD) |
| **Persona yang cocok** | Eksploratif (eksplorasi ide), Analitis-Evidence (validasi), Pengambil Keputusan Aplikatif (pengerucutan) |
| **Output minimum** | MoM di `.ai-doc/brainstorming/mom-{date}-{topic-slug}.md` |
| **Output tambahan** | Discussion Summary per sub-topik di `.ai-doc/brainstorming/discussion-{subtopic-slug}-{date}.md` |
| **Teknik** | 18 teknik dalam 2 area × 4 kategori — lihat `add-on/brainstorming/techniques.md` |
| **Template output** | `add-on/brainstorming/template/mom-template.md` dan `discussion-template.md` |
| **Steps** | 3 step: Setup → Facilitation (loop) → Wrap-Up |
| **HALT pattern** | Ya — agent tidak auto-pilot, tunggu konfirmasi user di setiap titik keputusan |
| **Integrasi persona** | Load dari `.ai-doc/personas/list.md` jika ada; user pilih peserta (user-only / +Joni / +Jono / keduanya) |
| **Idea tracking** | Setiap ide tercatat dengan status: `approved` / `rejected` / `modified` |
| **Mode moderator** | 3 mode: Eksploratif, Analitis-Evidence, Pengambil Keputusan Aplikatif — lihat `add-on/brainstorming/mode-moderator.md` |
| **Dokumen desain** | `tes/01-brainstrom-addon.md` |

---

### persona

| Atribut | Nilai |
|---|---|
| **Status** | `active` — siap digunakan |
| **Tujuan** | Menyiapkan template persona yang bisa dipakai berkali-kali oleh user. User memilih blueprint, memberi nama, dan agent generate folder persona di `.ai-doc/personas/` |
| **Kapan digunakan** | User ingin membuat persona untuk dipakai di sesi AI Documentor — baik untuk brainstorming, problem-solving, atau kebutuhan lainnya |
| **Kapan tidak digunakan** | User hanya ingin dokumentasi codebase biasa (gunakan core AD); user ingin brainstorming biasa tanpa persona khusus |
| **Blueprint tersedia** | Problem Solver 🔬, Creative Visionary ✨, Technical Architect 🏗️, Data Specialist 📊, Senior Backend Engineer ⚙️, Custom 📋 |
| **Output** | `.ai-doc/personas/<nama-persona>/persona.md` + `customize.toml` |
| **Entry file** | `add-on/persona/workflow.md` |
| **Steps** | 2 step: Select Blueprint → Generate Persona |
| **HALT pattern** | Ya — agent tidak auto-pilot, tunggu konfirmasi user di setiap titik keputusan |
| **Integrasi** | Persona yang sudah dibuat bisa dipakai oleh add-on brainstorming sebagai peserta sesi |
| **Persistensi** | Persona tersimpan di `.ai-doc/personas/` dan bisa dipakai kapan saja |

---

### tdd

| Atribut | Nilai |
|---|---|
| **Status** | `active` — mandiri dan hanya aktif setelah user memilih TDD |
| **Tujuan** | Memulai implementasi greenfield dari test yang gagal dan melacak RED/GREEN/REFACTOR sebagai control plane project |
| **Kapan digunakan** | Saat bootstrap greenfield atau implementasi behavior baru ketika user menyetujui TDD |
| **Kapan tidak digunakan** | Dokumentasi murni, brownfield tanpa opt-in, atau ketika constitution menyatakan TDD disabled |
| **Entry file** | `add-on/tdd/SKILL.md` |
| **Workflow** | `add-on/tdd/workflow.md` |
| **Test quality rules** | `add-on/tdd/writing-good-tests.md` |
| **Template** | `add-on/tdd/template/tdd-overview-template.md` |
| **Output** | `.ai-doc/tdd-overview.md` |
| **Status model** | `PLANNED` → `RED` → `GREEN` → `REFACTORING` → `REFACTORED`, dengan `BLOCKED`/`EXCEPTION` bila perlu |
| **Konstitusi** | Keputusan `TDD: Enabled` atau `TDD: Disabled` wajib dicatat di `.ai-doc/constitution.md` |

---

### qa-strategy

| Atribut | Nilai |
|---|---|
| **Status** | `active` — siap digunakan |
| **Tujuan** | Menyusun strategi testing evidence-based dari artefak `.ai-doc/` yang ada: risk matrix (Impact × Likelihood), test pyramid, quality gates, entry/exit criteria — lalu opsional menurunkan test plan per sprint/release |
| **Kapan digunakan** | Artefak inti sudah ada (project-overview / DCD / REST API doc) dan user ingin rencana testing yang terlacak ke bukti |
| **Kapan tidak digunakan** | Belum ada artefak `.ai-doc/` sama sekali (jalankan core flow dulu); eksekusi test di CI/tool (di luar mandat AI Documentor) |
| **Sumber inspirasi** | Skill `test-strategy` + `test-planning` dari repo qa-skills, diadaptasi ke konvensi AI Documentor |
| **Entry file** | `add-on/qa-strategy/workflow.md` |
| **Steps** | 3 step: Context Gathering → Strategy Building → Test Planning (opsional) |
| **HALT pattern** | Ya — agent tidak auto-pilot; file hanya ditulis setelah konfirmasi user |
| **Integrasi persona** | Persona bisa memfasilitasi sesi strategy; method menyediakan teknik prioritisasi |
| **Template output** | `add-on/qa-strategy/template/test-strategy-template.md` dan `test-plan-template.md` |
| **Output minimum** | `.ai-doc/qa/strategy/test-strategy.md` |
| **Output tambahan** | `.ai-doc/qa/strategy/plans/test-plan-{slug}-{date}.md` per sprint/release |

---

### qa-security

| Atribut | Nilai |
|---|---|
| **Status** | `active` — siap digunakan |
| **Tujuan** | Security screening & audit evidence-based: petakan attack surface, nilai coverage OWASP Top 10 (2025) tanpa kategori hilang diam-diam, dan kelola findings register dengan keputusan Remediate/Accept/Defer |
| **Kapan digunakan** | User minta security screening/audit; juga berguna untuk brownfield yang belum punya artefak `.ai-doc/` lengkap — Step 1 merekonstruksi attack surface dari codebase nyata |
| **Kapan tidak digunakan** | Eksekusi scanner/tool security di CI (di luar mandat AI Documentor — gunakan skill qa-skills); pemetaan regulasi compliance penuh |
| **Sumber inspirasi** | Skill `security-testing` dari repo qa-skills (konsep owasp-coverage Done When), diadaptasi ke konvensi AI Documentor |
| **Entry file** | `add-on/qa-security/workflow.md` |
| **Steps** | 3 step: Surface Mapping → Coverage Assessment → Findings & Mitigation Plan |
| **HALT pattern** | Ya — agent tidak auto-pilot; `Accepted Risk` wajib persetujuan eksplisit user dengan nama approver + tanggal |
| **Integrasi persona** | Persona bisa memfasilitasi sesi threat modeling ringan; method menyediakan teknik analisis risiko |
| **Template output** | `add-on/qa-security/template/security-coverage-template.md` |
| **Output minimum** | `.ai-doc/qa/security/security-coverage.md` (matrix + findings register dalam satu file) |
| **Output tambahan** | Re-audit log di file yang sama; action items bisa diturunkan ke test plan via add-on qa-strategy |

---

### qa-execution

| Atribut | Nilai |
|---|---|
| **Status** | `active` — siap digunakan |
| **Tujuan** | Mengeksekusi item test plan secara terlacak: mode manual / automated / hybrid, status PLANNED → IN PROGRESS → PASSED/FAILED/BLOCKED/DEFERRED, semuanya berbasis bukti eksekusi nyata |
| **Kapan digunakan** | Test plan sudah ada (hasil qa-strategy) dan user ingin menjalankannya — sebagian (batch P0) atau seluruhnya, dengan pelacakan status di control plane |
| **Kapan tidak digunakan** | Belum ada test plan valid (jalankan qa-strategy dulu); setup CI/CD pipeline atau scheduling; pembangunan framework automation dari nol |
| **Sumber inspirasi** | Skill `test-execution` + `test-reporting` dari repo qa-skills; disiplin bukti disalin dari add-on tdd |
| **Entry file** | `add-on/qa-execution/workflow.md` |
| **Steps** | 3 step: Activation → Execution Loop → Wrap-Up & Reporting |
| **HALT pattern** | Ya — konfirmasi batch sebelum mulai; setiap FAILED/BLOCKED dikonfirmasi sebelum lanjut; rekap dikonfirmasi sebelum sesi ditutup |
| **Mode eksekusi** | Manual (user jalankan, agent catat), Automated (agent jalankan command nyata + tangkap exit code/output), Hybrid (setup manual + verifikasi otomatis) |
| **Integrasi TDD** | Saat constitution menyatakan `TDD: Enabled`, item unit-level dieksekusi via siklus RED→GREEN add-on tdd — qa-execution hanya mencatat referensinya |
| **Template output** | `add-on/qa-execution/template/qa-overview-template.md` |
| **Output minimum** | `.ai-doc/qa/execution/qa-overview.md` + pembaruan kolom §7 test plan |
| **Output tambahan** | Riwayat eksekusi per item, daftar hipotesis FAILED, rekap gap terhadap exit criteria |

---

### method

| Atribut | Nilai |
|---|---|
| **Status** | `active` — siap digunakan |
| **Tujuan** | Menyediakan koleksi 125 teknik/metode/framework/story lintas 4 kategori utama: Brain Methods (45 teknik ideation), Solving Methods (25 metode problem-solving), Innovation Frameworks (30 framework strategi inovasi), Story Methods (25 storytelling framework). |
| **Kapan digunakan** | User ingin mencari/menggunakan teknik brainstorming, metode problem-solving, framework inovasi, atau storytelling — baik standalone maupun dalam sesi persona |
| **Kapan tidak digunakan** | User hanya ingin dokumentasi codebase biasa; user hanya ingin brainstorming biasa (gunakan add-on brainstorming) |
| **Koleksi** | 45 Brain Methods (10 kategori: Collaborative, Creative, Deep, Introspective, Structured, Theatrical, Wild, Biomimetic, Quantum, Cultural) |
| | 25 Solving Methods (5 fase: Diagnosis, Analysis, Synthesis, Evaluation, Implementation) |
| | 30 Innovation Frameworks (6 kategori: Disruption, Business Model, Market Analysis, Strategic, Value Chain, Technology) |
| | 25 Story Methods (6 kategori: Transformation, Strategic, Persuasive, Analytical, Emotional, Conversational) |
| **Persona yang cocok** | Problem Solver 🔬, Creative Visionary ✨, Technical Architect 🏗️, Data Specialist 📊, Senior Backend Engineer ⚙️ — setiap persona punya rekomendasi metode spesifik di techniques-index |
| **Entry file** | `add-on/method/techniques-index.md` |
| **File konten** | `add-on/method/brain/README.md` (45 teknik), `add-on/method/solving/README.md` (25 metode), `add-on/method/innovation/README.md` (30 framework), `add-on/method/story/README.md` (25 story type) |
| **Template output** | `add-on/method/template/` (opsional) |

---

## Aturan Penggunaan

1. **Add-on hanya aktif jika diminta eksplisit.** Jangan auto-run.
2. **Satu sesi = satu add-on aktif.** Tidak boleh menjalankan dua add-on bersamaan.
3. **Persona boleh menyertai add-on** (misal: "jadikan aku moderator pengambil keputusan" + "brainstorming").
4. **Output add-on tetap mengikuti aturan AI Documentor:**
   - Berbasis bukti (`evidence-based`)
   - Tidak over-infer
   - Hanya buat artefak jika diminta
5. **Setelah selesai**, return kontrol ke core AI Documentor.
