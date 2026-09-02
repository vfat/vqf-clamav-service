# 🔧 Solving Methods — 25 Metode Problem-Solving

> **AD-AO-METHOD-SOLVING**
> Koleksi lengkap 25 metode problem-solving yang terbagi dalam 5 fase berurutan.
> Dirancang untuk menganalisis, mengevaluasi, dan mengimplementasikan solusi.
> Format: Pipeline Prompts + Output + Panduan Penggunaan.

---

## Executive Summary

| Item | Detail |
|---|---|
| **Total Metode** | 25 |
| **Fase** | 5 berurutan: Diagnosis → Analysis → Synthesis → Evaluation → Implementation |
| **Fokus** | Problem-solving — menyelesaikan masalah secara sistematis |
| **Alur** | Sequential (harus berurutan untuk hasil optimal) |
| **Output** | Solusi tervalidasi + rencana eksekusi |
| **Cocok Untuk** | Semua persona, especially Problem Solver, Architect, Backend Engineer |

### Quick Persona → Solving Methods

| Persona | Fase Terkuat | Metode Andalan |
|---|---|---|
| 🔬 Problem Solver | Diagnosis, Synthesis | Five Whys, TRIZ, FMEA, Feasibility Study |
| ✨ Creative Visionary | Synthesis | Lateral Thinking, Biomimicry, Synectics |
| 🏗️ Technical Architect | Analysis, Evaluation | Systems Thinking, Decision Matrix, Cost Benefit |
| 📊 Data Specialist | Analysis | Gap Analysis, Pareto, Is/Is Not, Monitoring Dashboard |
| ⚙️ Backend Engineer | Evaluation, Implementation | FMEA, Risk Assessment, PDCA, Gantt Chart |

---

## Daftar Isi

1. [Kategori Overview](#kategori-overview)
2. [Diagnosis (5 metode)](#1-diagnosis)
3. [Analysis (5 metode)](#2-analysis)
4. [Synthesis (5 metode)](#3-synthesis)
5. [Evaluation (5 metode)](#4-evaluation)
6. [Implementation (5 metode)](#5-implementation)
7. [Perbandingan dengan Brain Methods](#perbandingan-dengan-brain-methods)
8. [Matriks Pemilihan](#matriks-pemilihan)

---

## Format Metode

Setiap metode dijelaskan dengan format:

```
**Nama Metode**
- Kategori: ...
- Deskripsi: ...
- Prompt Facilitation: pertanyaan pipeline (dipisah |)
- Kapan Cocok: situasi terbaik
- Output: hasil yang diharapkan
- Catatan: tips tambahan
```

---

## Kategori Overview

Berbeda dengan Brain Methods yang fokus pada **ideation**, Solving Methods fokus pada **problem-solving** — dari diagnosis sampai implementasi. Alurnya mengikuti siklus:

```
┌────────────┐    ┌────────────┐    ┌────────────┐    ┌────────────┐    ┌────────────┐
│  DIAGNOSIS │───▶│  ANALYSIS  │───▶│  SYNTHESIS │───▶│ EVALUATION │───▶│IMPLEMENTATION│
│ (5 metode) │    │ (5 metode) │    │ (5 metode) │    │ (5 metode) │    │ (5 metode)  │
└────────────┘    └────────────┘    └────────────┘    └────────────┘    └────────────┘
     │                 │                 │                 │                 │
  Apa masalahnya?  Mengapa?         Solusi apa?      Mana yang terbaik?  Bagaimana
  Seberapa parah?  Apa penyebabnya? Bagaimana         Risiko apa?         menjalankannya?
                                    caranya?
```

| # | Kategori | Jumlah | Fase | Untuk Apa |
|---|----------|--------|------|-----------|
| 1 | **Diagnosis** | 5 | Awal | Memahami & mendefinisikan masalah |
| 2 | **Analysis** | 5 | Kedalaman | Menganalisis penyebab & constraint |
| 3 | **Synthesis** | 5 | Solusi | Menghasilkan solusi kreatif |
| 4 | **Evaluation** | 5 | Keputusan | Memilih & memvalidasi solusi |
| 5 | **Implementation** | 5 | Eksekusi | Menjalankan & memantau solusi |

---

## 1. Diagnosis

> Fokus: **memahami masalah** sebelum menyelesaikannya. "Solving the right problem."

### 1.1 Five Whys Root Cause

- **Deskripsi:** Gali lapisan gejala untuk menemukan root cause dengan bertanya "kenapa" lima kali.
- **Prompt Facilitation:**
  - "Kenapa ini terjadi?"
  - "Kenapa itu terjadi?"
  - "Kenapa itu terjadi?"
  - "Apa yang ada di bawahnya?"
  - "Apa root cause-nya?"
- **Kapan Cocok:** Masalah berulang, troubleshooting, saat gejala jelas tapi penyebab tidak.
- **Output:** Root cause statement yang actionable.
- **Catatan:** Jangan berhenti di gejala. Terus tanya "kenapa" sampai ke fondasi.

### 1.2 Fishbone Diagram (Ishikawa)

- **Deskripsi:** Map semua potensi penyebab di kategori: people, process, materials, equipment, environment.
- **Prompt Facilitation:**
  - "Faktor people apa yang berkontribusi?"
  - "Issue process apa?"
  - "Masalah material apa?"
  - "Faktor equipment apa?"
  - "Kondisi environment apa?"
- **Kapan Cocok:** Masalah kompleks dengan banyak potensi penyebab, quality control, manufacturing.
- **Output:** Visual diagram semua potensi penyebab terorganisasi per kategori.
- **Catatan:** Jangan hanya fokus satu kategori. Semua "tulang" fishbone harus dieksplorasi.

### 1.3 Problem Statement Refinement

- **Deskripsi:** Transformasi keluhan vague menjadi problem statement yang precise dan actionable.
- **Prompt Facilitation:**
  - "Apa yang salah secara spesifik?"
  - "Siapa yang terdampak dan bagaimana?"
  - "Kapan dan di mana ini terjadi?"
  - "Apa gap antara kondisi sekarang dan yang diinginkan?"
  - "Apa yang membuat ini jadi masalah?"
- **Kapan Cocok:** Awal proyek, saat masalah belum jelas, sebelum solusi.
- **Output:** Problem statement 1-2 kalimat yang jelas dan terukur.
- **Catatan:** Problem statement yang baik = setengah solusi ditemukan.

### 1.4 Is/Is Not Analysis

- **Deskripsi:** Definisikan batasan masalah dengan membandingkan di mana masalah ada vs tidak ada.
- **Prompt Facilitation:**
  - "Di mana masalah terjadi?"
  - "Di mana tidak?"
  - "Kapan terjadi?"
  - "Kapan tidak?"
  - "Siapa yang mengalami?"
  - "Siapa yang tidak?"
  - "Pola apa yang muncul?"
- **Kapan Cocok:** Masalah yang sporadis, narrowing investigation scope, debugging.
- **Output:** Boundary definition yang jelas — apa yang IN dan apa yang OUT.
- **Catatan:** Kontras antara "is" dan "is not" sering mengungkap pola yang tersembunyi.

### 1.5 Systems Thinking

- **Deskripsi:** Map elemen sistem yang saling terkait, feedback loops, dan leverage points.
- **Prompt Facilitation:**
  - "Apa saja komponen sistem?"
  - "Hubungan apa yang ada?"
  - "Feedback loop apa yang ada?"
  - "Delay apa yang terjadi?"
  - "Di mana leverage points?"
- **Kapan Cocok:** Masalah sistemik, organisasi kompleks, masalah yang "selesaikan satu muncul yang lain".
- **Output:** System map dengan komponen, hubungan, dan leverage points.
- **Catatan:** Masalah sistemik tidak punya single cause. Fokus pada pola, bukan event.

---

## 2. Analysis

> Fokus: **menganalisis masalah secara mendalam** — penyebab, constraint, gap, dan risiko.

### 2.1 Force Field Analysis

- **Deskripsi:** Identifikasi driving forces (mendorong solusi) dan restraining forces (menghambat progress).
- **Prompt Facilitation:**
  - "Gaya apa yang mendorong solusi?"
  - "Gaya apa yang menghambat perubahan?"
  - "Mana yang paling kuat?"
  - "Mana yang bisa kita pengaruhi?"
  - "Apa strateginya?"
- **Kapan Cocok:** Change management, saat ada resistance, sebelum implementasi.
- **Output:** Map gaya pendorong vs penghambat + strategi intervensi.
- **Catatan:** Sering lebih efektif melemahkan restraining forces daripada menambah driving forces.

### 2.2 Pareto Analysis (80/20)

- **Deskripsi:** Aplikasikan aturan 80/20 untuk mengidentifikasi vital few causes yang menciptakan mayoritas dampak.
- **Prompt Facilitation:**
  - "Penyebab apa saja yang ada?"
  - "Berapa frekuensi atau dampak masing-masing?"
  - "Berapa cumulative impact?"
  - "Vital few apa yang mendorong 80%?"
  - "Fokus di mana?"
- **Kapan Cocok:** Banyak penyebab potensial, resource terbatas, prioritization.
- **Output:** Ranked list penyebab + fokus area vital few.
- **Catatan:** 20% penyebab = 80% masalah. Fokus di sana dulu.

### 2.3 Gap Analysis

- **Deskripsi:** Bandingkan kondisi saat ini dengan kondisi yang diinginkan di berbagai dimensi.
- **Prompt Facilitation:**
  - "Apa kondisi saat ini?"
  - "Apa kondisi yang diinginkan?"
  - "Gap apa yang ada?"
  - "Seberapa besar gap?"
  - "Apa penyebab gap?"
  - "Fokus prioritas?"
- **Kapan Cocok:** Performance improvement, capability assessment, strategic planning.
- **Output:** Gap map dengan ukuran dan prioritas per dimensi.
- **Catatan:** Gap yang terukur = gap yang bisa ditutup.

### 2.4 Constraint Identification (Theory of Constraints)

- **Deskripsi:** Temukan bottleneck yang membatasi performa sistem.
- **Prompt Facilitation:**
  - "Apa constraint-nya?"
  - "Apa yang membatasi throughput?"
  - "Apa yang harus kita optimasi?"
  - "Apa yang terjadi kalau constraint diangkat?"
  - "Constraint berikutnya apa?"
- **Kapan Cocok:** Optimasi proses, bottleneck identification, throughput improvement.
- **Output:** Identifikasi bottleneck + rencana eliminasi.
- **Catatan:** Optimasi di luar bottleneck = waste. Fokus pada constraint utama.

### 2.5 Failure Mode Analysis (FMEA)

- **Deskripsi:** Antisipasi bagaimana solusi bisa gagal dan engineer preventions sebelum masalah terjadi.
- **Prompt Facilitation:**
  - "Apa yang bisa salah?"
  - "Berapa likelihood?"
  - "Berapa impact?"
  - "Bagaimana mencegah?"
  - "Bagaimana deteksi dini?"
  - "Apa mitigasinya?"
- **Kapan Cocok:** Risk assessment, safety-critical systems, sebelum implementasi besar.
- **Output:** Risk register dengan likelihood, impact, prevention, dan mitigation.
- **Catatan:** Lebih murah mencegah daripada memperbaiki. Lakukan ini SEBELUM implementasi.

---

## 3. Synthesis

> Fokus: **menghasilkan solusi** — dari analisis menjadi ide solusi konkret.

### 3.1 TRIZ Contradiction Matrix

- **Deskripsi:** Resolusi kontradiksi teknis menggunakan 40 inventive principles dari analisis pattern paten.
- **Prompt Facilitation:**
  - "Apa yang membaik?"
  - "Apa yang memburuk?"
  - "Apa kontradiksinya?"
  - "Prinsip apa yang berlaku?"
  - "Bagaimana resolusinya?"
- **Kapan Cocok:** Engineering problems, technical trade-offs, optimasi yang saling bertentangan.
- **Output:** Solusi berdasarkan inventive principles TRIZ.
- **Catatan:** TRIZ = "Theory of Inventive Problem Solving". Sangat powerful untuk masalah teknis.

### 3.2 Lateral Thinking Techniques

- **Deskripsi:** Gunakan provocative operations dan random entry untuk memecah pattern-thinking.
- **Prompt Facilitation:**
  - "Buat provokasi"
  - "Tantang asumsi"
  - "Gunakan stimulus random"
  - "Lari dari ide dominan"
  - "Generate alternatif"
- **Kapan Cocok:** Solusi konvensional gagal, butuh breakthrough, stuck di local optimum.
- **Output:** Solusi non-obvious yang tidak terpikir sebelumnya.
- **Catatan:** Lateral thinking ≠ logical thinking. Tentang memecah pola, bukan mengikuti pola.

### 3.3 Morphological Analysis

- **Deskripsi:** Eksplorasi sistematis semua kombinasi parameter solusi untuk menemukan konfigurasi optimal.
- **Prompt Facilitation:**
  - "Apa parameter kunci?"
  - "Opsi apa untuk masing-masing?"
  - "Coba kombinasi berbeda"
  - "Pola apa yang muncul?"
  - "Apa yang optimal?"
- **Kapan Cocok:** Desain sistem kompleks, architecture decisions, multi-variable optimization.
- **Output:** Matrix parameter × opsi + kombinasi optimal.
- **Catatan:** Sangat sistematis tapi butuh waktu. Cocok untuk architect persona.

### 3.4 Biomimicry Problem Solving

- **Deskripsi:** Belajar dari 3.8 miliar tahun R&D alam untuk menemukan solusi elegan.
- **Prompt Facilitation:**
  - "Bagaimana alam menyelesaikan ini?"
  - "Analogi biologis apa?"
  - "Prinsip apa yang bisa ditransfer?"
  - "Bagaimana mengadaptasi?"
- **Kapan Cocok:** Engineering challenges, sustainability, design optimization.
- **Output:** Solusi terinspirasi biomimicry + prinsip transfer.
- **Catatan:** Alam sudah "mensolve" banyak masalah yang kita hadapi. Tinggal mengadaptasi.

### 3.5 Synectics Method

- **Deskripsi:** Buat yang asing jadi familiar dan yang familiar jadi asing melalui analogi.
- **Prompt Facilitation:**
  - "Ini seperti apa?"
  - "Bagaimana mereka serupa?"
  - "Metafora apa yang cocok?"
  - "Apa yang disarankan?"
  - "Insight apa yang muncul?"
- **Kapan Cocok:** Masalah yang terlalu familiar (butuh defamiliarisasi), creative breakthrough.
- **Output:** Insight dari analogi + solusi yang terinspirasi.
- **Catatan:** "Make the strange familiar and the familiar strange" — inti dari Synectics.

---

## 4. Evaluation

> Fokus: **memilih & memvalidasi solusi** — dari banyak opsi menjadi satu keputusan.

### 4.1 Decision Matrix

- **Deskripsi:** Evaluasi opsi solusi secara sistematis terhadap kriteria berbobot untuk seleksi objektif.
- **Prompt Facilitation:**
  - "Apa saja opsi?"
  - "Kriteria apa yang penting?"
  - "Berapa bobot masing-masing?"
  - "Rating setiap opsi"
  - "Hitung skor"
  - "Apa yang menang?"
- **Kapan Cocok:** Multi-option decision, saat ada banyak kriteria, objektif selection.
- **Output:** Ranked options dengan skor per kriteria.
- **Catatan:** Bobot kriteria harus ditentukan SEBELUM evaluasi. Hindari bias post-hoc.

### 4.2 Cost Benefit Analysis

- **Deskripsi:** Kuantifikasi cost dan benefit yang diharapkan untuk mendukung keputusan investasi rasional.
- **Prompt Facilitation:**
  - "Apa saja cost-nya?"
  - "Apa saja benefit-nya?"
  - "Kuantifikasi masing-masing"
  - "Berapa payback period?"
  - "Berapa ROI?"
  - "Apa rekomendasinya?"
- **Kapan Cocok:** Investment decisions, resource allocation, business case.
- **Output:** Cost-benefit breakdown + ROI + rekomendasi.
- **Catatan:** Quantify what you can, estimate what you can't. Angka kasar lebih baik tanpa angka.

### 4.3 Risk Assessment Matrix

- **Deskripsi:** Evaluasi risiko solusi di dimensi likelihood dan impact untuk prioritisasi mitigasi.
- **Prompt Facilitation:**
  - "Apa yang bisa salah?"
  - "Berapa probabilitasnya?"
  - "Berapa dampaknya?"
  - "Plot di matrix"
  - "Berapa risk score?"
  - "Rencana mitigasi?"
- **Kapan Cocok:** Sebelum implementasi besar, risk-sensitive projects, compliance.
- **Output:** Risk matrix + risk scores + mitigation plans.
- **Catatan:** Fokus pada high-likelihood + high-impact dulu. Sisanya monitor.

### 4.4 Pilot Testing Protocol

- **Deskripsi:** Desain eksperimen skala kecil untuk validasi solusi sebelum full implementation.
- **Prompt Facilitation:**
  - "Apa yang akan kita test?"
  - "Apa success criteria?"
  - "Apa test plan-nya?"
  - "Data apa yang dikumpulkan?"
  - "Apa yang kita pelajari?"
  - "Scale atau pivot?"
- **Kapan Cocok:** Solusi baru yang belum terbukti, high-stakes decisions, innovation projects.
- **Output:** Pilot plan + success criteria + learnings + go/no-go decision.
- **Catatan:** Pilot = belajar murah sebelum investasi mahal. Selalu pilot dulu jika memungkinkan.

### 4.5 Feasibility Study

- **Deskripsi:** Assess kelayakan teknis, operasional, finansial, dan jadwal dari opsi solusi.
- **Prompt Facilitation:**
  - "Apakah technically possible?"
  - "Operationally viable?"
  - "Financially sound?"
  - "Schedule realistic?"
  - "Overall feasibility?"
- **Kapan Cocok:** Sebelum commit ke solusi besar, project approval, stakeholder buy-in.
- **Output:** Feasibility assessment per dimensi + overall recommendation.
- **Catatan:** Satu dimensi "tidak layak" bisa kill entire project. Cek semua dimensi.

---

## 5. Implementation

> Fokus: **menjalankan solusi** — dari rencana menjadi eksekusi nyata.

### 5.1 PDCA Cycle (Plan-Do-Check-Act)

- **Deskripsi:** Siklus iteratif: Plan, Do, Check, Act — implementasi dengan continuous learning.
- **Prompt Facilitation:**
  - "Apa rencananya?"
  - "Eksekusi rencana"
  - "Cek hasil"
  - "Apa yang berhasil?"
  - "Apa yang tidak?"
  - "Adjust dan ulangi"
- **Kapan Cocok:** Continuous improvement, iterative implementation, quality management.
- **Output:** PDCA cycle documentation + learnings + adjustments.
- **Catatan:** PDCA tidak pernah berhenti. Setiap cycle = improvement.

### 5.2 Gantt Chart Planning

- **Deskripsi:** Visualisasi timeline proyek dengan task, dependencies, dan milestones.
- **Prompt Facilitation:**
  - "Apa saja task-nya?"
  - "Urutan apa?"
  - "Dependencies apa?"
  - "Timeline-nya apa?"
  - "Siapa responsible?"
  - "Milestones apa?"
- **Kapan Cocok:** Project planning, resource allocation, stakeholder communication.
- **Output:** Gantt chart dengan timeline, dependencies, dan milestones.
- **Catatan:** Dependencies adalah kunci. Task tanpa dependency bisa di-parallel.

### 5.3 Stakeholder Mapping

- **Deskripsi:** Identifikasi semua pihak yang terdampak dan rencana engagement strategy.
- **Prompt Facilitation:**
  - "Siapa yang terdampak?"
  - "Apa interest-nya?"
  - "Apa influence-nya?"
  - "Apa engagement strategy-nya?"
  - "Bagaimana komunikasinya?"
- **Kapan Cocok:** Change management, project kickoff, saat ada resistance.
- **Output:** Stakeholder map + engagement strategy per stakeholder.
- **Catatan:** High influence + high interest = manage closely. Jangan diabaikan.

### 5.4 Change Management Protocol

- **Deskripsi:** Kelola dimensi organisasi dan manusia dari implementasi solusi secara sistematis.
- **Prompt Facilitation:**
  - "Apa yang berubah?"
  - "Siapa yang terdampak?"
  - "Resistance apa yang diharapkan?"
  - "Bagaimana komunikasinya?"
  - "Bagaimana support transition?"
  - "Bagaimana sustain?"
- **Kapan Cocok:** Organizational change, new system adoption, culture shift.
- **Output:** Change management plan + communication plan + support plan.
- **Catatan:** 70% perubahan gagal karena faktor manusia, bukan teknis. Kelola manusia dulu.

### 5.5 Monitoring Dashboard

- **Deskripsi:** Buat sistem tracking visual untuk key metrics agar solusi deliver hasil yang diharapkan.
- **Prompt Facilitation:**
  - "Metrics apa yang penting?"
  - "Target apa?"
  - "Bagaimana mengukur?"
  - "Bagaimana memvisualisasikan?"
  - "Apa yang trigger action?"
  - "Review frequency?"
- **Kapan Cocok:** Post-implementation, ongoing operations, KPI tracking.
- **Output:** Dashboard design + metrics + targets + alert thresholds.
- **Catatan:** Yang tidak diukur = tidak dikelola. Dashboard = visibility.

---

## Perbandingan dengan Brain Methods

| Aspek | Brain Methods (45 teknik) | Solving Methods (25 metode) |
|-------|---------------------------|---------------------------|
| **Fokus** | Ideation — menghasilkan ide | Problem-solving — menyelesaikan masalah |
| **Fase** | Divergent thinking | Full problem-solving lifecycle |
| **Alur** | Tidak berurutan, bisa random | Berurutan: Diagnosis → Analysis → Synthesis → Evaluation → Implementation |
| **Prompt** | Deskripsi saja | Pipeline prompts (dipisah `\|`) |
| **Output** | Ide-ide kreatif | Solusi yang tervalidasi + rencana eksekusi |
| **Persona fit** | Semua persona | Lebih cocok Problem Solver, Architect, Backend Engineer |
| **Kapan pakai** | Butuh ide segar | Masalah sudah jelas, butuh solusi |

### Teknik yang Ada di Kedua Sistem

| Teknik | Brain Methods | Solving Methods |
|--------|---------------|-----------------|
| Five Whys | ✅ (Deep) | ✅ (Diagnosis) |
| Morphological Analysis | ✅ (Deep) | ✅ (Synthesis) |
| Biomimicry | ✅ (Biomimetic) | ✅ (Synthesis) |

### Kombinasi Ideal

```
Brain Methods (ideasi)          Solving Methods (eksekusi)
        │                                  │
        ▼                                  ▼
   Generate ideas ──────────────▶ Validate & implement
   "Apa yang bisa?"              "Bagaimana caranya?"
        │                                  │
        └──────────┬───────────────────────┘
                   │
                   ▼
            Complete Workflow
        Ideasi → Analisis → Solusi → Evaluasi → Eksekusi
```

---

## Matriks Pemilihan

### Berdasarkan Fase Proyek

| Fase | Metode Rekomendasi |
|---|---|
| **Masalah belum jelas** | Problem Statement Refinement, Is/Is Not, Five Whys |
| **Butuh paham root cause** | Five Whys, Fishbone, Systems Thinking |
| **Butuh analisis mendalam** | Force Field, Pareto, Gap Analysis, Constraint ID |
| **Butuh solusi kreatif** | TRIZ, Lateral Thinking, Biomimicry, Synectics |
| **Butuh pilih solusi** | Decision Matrix, Cost Benefit, Risk Assessment |
| **Butuh validasi** | Pilot Testing, Feasibility Study |
| **Butuh eksekusi** | PDCA, Gantt Chart, Change Management |
| **Butuh monitoring** | Monitoring Dashboard |

### Berdasarkan Tipe Masalah

| Tipe Masalah | Metode |
|---|---|
| **Masalah berulang** | Five Whys, Fishbone, Pareto |
| **Masalah kompleks/sistemik** | Systems Thinking, Force Field |
| **Masalah teknis** | TRIZ, Biomimicry, Morphological Analysis |
| **Masalah organisasi** | Change Management, Stakeholder Mapping |
| **Masalah dengan banyak opsi** | Decision Matrix, Cost Benefit |
| **Masalah high-risk** | Risk Assessment, Failure Mode, Pilot Testing |
| **Masalah resource terbatas** | Pareto, Constraint Identification |

### Berdasarkan Persona

| Persona | Metode yang Cocok |
|---|---|
| 🔬 **Problem Solver** | Five Whys, Fishbone, TRIZ, FMEA, Feasibility Study |
| ✨ **Creative Visionary** | Lateral Thinking, Biomimicry, Synectics, Pilot Testing |
| 🏗️ **Technical Architect** | Systems Thinking, Morphological Analysis, TRIZ, Decision Matrix, Feasibility Study |
| 📊 **Data Specialist** | Problem Statement Refinement, Gap Analysis, Pareto, Is/Is Not, Monitoring Dashboard |
| ⚙️ **Backend Engineer** | Constraint Identification, FMEA, Risk Assessment, PDCA, Gantt Chart, Monitoring Dashboard |
