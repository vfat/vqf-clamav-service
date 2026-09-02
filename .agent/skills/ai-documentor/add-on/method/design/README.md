# 🎨 Design Methods — 30 Metode Design Thinking

> **AD-AO-METHOD-DESIGN**
> Koleksi lengkap 30 metode design thinking yang terbagi dalam 6 fase berurutan.
> Dirancang untuk memahami pengguna, mendefinisikan masalah, menghasilkan ide, memprototipe, menguji, dan mengimplementasikan solusi.
> Format: Pipeline Prompts + Output + Panduan Penggunaan.

---

## Executive Summary

| Item | Detail |
|---|---|
| **Total Metode** | 30 |
| **Fase** | 6 berurutan: Empathize → Define → Ideate → Prototype → Test → Implement |
| **Fokus** | Design thinking — human-centered innovation |
| **Alur** | Sequential (idealnya iteratif antar fase) |
| **Output** | Solusi teruji + rencana implementasi |
| **Cocok Untuk** | Semua persona, terutama Problem Solver, Creative Visionary, Architect, Backend Engineer |
| **Sumber CSV** | `design/design-methods.csv` |

### Quick Persona → Design Methods

| Persona | Fase Terkuat | Metode Andalan |
|---|---|---|
| 🔬 Problem Solver | Empathize, Define, Test | User Interviews, How Might We, Usability Testing |
| ✨ Creative Visionary | Ideate, Prototype | SCAMPER Design, Crazy 8s, Provotype Sketching |
| 🏗️ Technical Architect | Prototype, Implement | Wizard of Oz, Service Blueprinting, Design System Creation |
| 📊 Data Specialist | Empathize, Test | Journey Mapping, Feedback Capture Grid, A/B Testing |
| ⚙️ Backend Engineer | Implement | Pilot Programs, Measurement Framework, Design System Creation |

---

## Daftar Isi

1. [Kategori Overview](#kategori-overview)
2. [Empathize (5 metode)](#1-empathize)
3. [Define (5 metode)](#2-define)
4. [Ideate (5 metode)](#3-ideate)
5. [Prototype (5 metode)](#4-prototype)
6. [Test (5 metode)](#5-test)
7. [Implement (5 metode)](#6-implement)
8. [Perbandingan dengan Solving Methods](#perbandingan-dengan-solving-methods)
9. [Matriks Pemilihan](#matriks-pemilihan)
10. [Data CSV](#data-csv)

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

Berbeda dengan Solving Methods yang fokus pada **problem-solving secara umum**, Design Methods fokus pada **human-centered design** — menempatkan pengguna sebagai pusat dari setiap keputusan. Alurnya mengikuti siklus Design Thinking klasik:

```
┌────────────┐    ┌────────────┐    ┌────────────┐    ┌────────────┐    ┌────────────┐    ┌────────────┐
│  EMPATHIZE  │───▶│   DEFINE   │───▶│   IDEATE   │───▶│ PROTOTYPE  │───▶│    TEST    │───▶│  IMPLEMENT │
│ (5 metode)  │    │ (5 metode) │    │ (5 metode) │    │ (5 metode) │    │ (5 metode) │    │ (5 metode) │
└────────────┘    └────────────┘    └────────────┘    └────────────┘    └────────────┘    └────────────┘
     │                │                │                │                │                │
  Siapa user?     Apa masalah      Bagaimana        Seperti apa       Berfungsi?       Bagaimana
  Apa yang        sebenarnya?      kita menyelesaikannya? bentuk solusi?   Apa yang perlu   menjalankannya
  mereka rasakan?                                                          diubah?         di dunia nyata?
                                                                                            
                          ◀───── iterate / loop ─────▶
```

| # | Kategori | Jumlah | Fase | Untuk Apa |
|---|----------|--------|------|-----------|
| 1 | **Empathize** | 5 | Awal | Memahami pengguna secara mendalam |
| 2 | **Define** | 5 | Framing | Merumuskan masalah & opportunity |
| 3 | **Ideate** | 5 | Eksplorasi | Menghasilkan beragam solusi |
| 4 | **Prototype** | 5 | Bentuk | Mewujudkan ide menjadi tangible |
| 5 | **Test** | 5 | Validasi | Menguji solusi dengan pengguna nyata |
| 6 | **Implement** | 5 | Eksekusi | Meluncurkan & memantau dampak |

---

## 1. Empathize

> Fokus: **memahami pengguna** — pain points, perilaku, dan kebutuhan yang belum terucap.

### 1.1 User Interviews

- **Deskripsi:** Conduct deep conversations to understand user needs, experiences, and pain points through active listening.
- **Prompt Facilitation:**
  - "What brings you here today?"
  - "Walk me through a recent experience"
  - "What frustrates you most?"
  - "What would make this easier?"
  - "Tell me more about that"
- **Kapan Cocok:** Riset awal produk, eksplorasi masalah baru, sebelum redesign.
- **Output:** Rich qualitative insights tentang kebutuhan, frustrasi, dan motivasi pengguna.
- **Catatan:** Active listening > asking. Jangan jual solusi, dengarkan tanpa menghakimi.

### 1.2 Empathy Mapping

- **Deskripsi:** Create visual representation of what users say, think, do, and feel to build deep understanding.
- **Prompt Facilitation:**
  - "What did they say?"
  - "What might they be thinking?"
  - "What actions did they take?"
  - "What emotions surfaced?"
- **Kapan Cocok:** Sintesis riset kualitatif, alignment tim, ekspos pain points.
- **Output:** Peta empati 4 kuadran (say, think, do, feel) yang bisa digunakan sebagai kompas desain.
- **Catatan:** Akui gap antara say dan do — apa yang mereka katakan sering berbeda dari apa yang mereka lakukan.

### 1.3 Shadowing

- **Deskripsi:** Observe users in their natural environment to see unspoken behaviors and contextual factors.
- **Prompt Facilitation:**
  - "Watch without interrupting"
  - "Note their workarounds"
  - "What patterns emerge?"
  - "What do they not say?"
- **Kapan Cocok:** Konteks fisik/kompleks, perilaku yang tidak sadar, riset lapangan.
- **Output:** Catatan observasi perilaku natural beserta workarounds yang ditemukan.
- **Catatan:** Tanyakan izin sebelum shadowing. Behavior yang terlihat sering kali lebih jujur dari interview.

### 1.4 Journey Mapping

- **Deskripsi:** Document complete user experience across touchpoints to identify pain points and opportunities.
- **Prompt Facilitation:**
  - "What's their starting point?"
  - "What steps do they take?"
  - "Where do they struggle?"
  - "What delights them?"
  - "What's the emotional arc?"
- **Kapan Cocok:** Memetakan end-to-end experience, identifikasi touchpoint kritis, alignment lintas tim.
- **Output:** Journey map visual dengan emotional arc, pain points, dan opportunity areas.
- **Catatan:** Bikin journey dari perspektif pengguna, bukan dari perspektif sistem atau organisasi.

### 1.5 Diary Studies

- **Deskripsi:** Have users document experiences over time to capture authentic moments and evolving needs.
- **Prompt Facilitation:**
  - "What did you experience today?"
  - "How did you feel?"
  - "What worked or didn't?"
  - "What surprised you?"
- **Kapan Cocok:** Perjalanan panjang, pola musiman, perilaku yang muncul perlahan.
- **Output:** Catatan longitudinal dari pengalaman nyata pengguna dalam periode waktu tertentu.
- **Catatan:** Lebih ringan daripada interview; ideal untuk riset remote atau saat interview sulit dilakukan.

---

## 2. Define

> Fokus: **merumuskan masalah** — translate insights menjadi opportunity area yang jelas.

### 2.1 Problem Framing

- **Deskripsi:** Transform observations into clear, actionable problem statements that inspire solution generation.
- **Prompt Facilitation:**
  - "What's the real problem?"
  - "Who experiences this?"
  - "Why does it matter?"
  - "What would success look like?"
- **Kapan Cocok:** Setelah riset empathize, sebelum mulai ideation, alignment tim.
- **Output:** Problem statement yang menginspirasi — bukan kalimat pasif, melainkan tantangan yang bisa diselesaikan.
- **Catatan:** Frame sebagai opportunity, bukan sekadar daftar masalah. Bahasa memengaruhi solusi yang muncul.

### 2.2 How Might We

- **Deskripsi:** Reframe problems as opportunity questions that open solution space without prescribing answers.
- **Prompt Facilitation:**
  - "How might we help users...?"
  - "How might we make it easier to...?"
  - "How might we reduce the friction of...?"
- **Kapan Cocok:** Transisi dari define ke ideate, brainstorming kick-off, ekspansi solusi space.
- **Output:** Daftar HMW questions yang siap membawa tim ke sesi ideation.
- **Catatan:** HMW yang terlalu lebar → tidak fokus. HMW yang terlalu sempit → membatasi kreativitas. Cari sweet spot-nya.

### 2.3 Point of View Statement

- **Deskripsi:** Create specific user-centered problem statements that capture who, what, and why.
- **Prompt Facilitation:**
  - "User type needs what because insight"
  - "What's driving this need?"
  - "Why does it matter to them?"
- **Kapan Cocok:** Mendokumentasikan insight spesifik, handover ke tim desain/produk.
- **Output:** POV statement terstruktur: `[user] needs [need] because [insight]`.
- **Catatan:** Jangan masukkan solusi di POV. Fokus pada user-need-insight, biarkan ideation mencari solusi.

### 2.4 Affinity Clustering

- **Deskripsi:** Group related observations and insights to reveal patterns and opportunity themes.
- **Prompt Facilitation:**
  - "What connects these?"
  - "What themes emerge?"
  - "Group similar items"
  - "Name each cluster"
  - "What story do they tell?"
- **Kapan Cocok:** Data riset banyak (sticker notes, transcripts), menemukan tema, menyintesis insights.
- **Output:** Klaster tema dengan label yang merepresentasikan insight utama.
- **Catatan:** Beri nama cluster dengan makna, bukan kategori literal. Nama yang baik mengundang aksi.

### 2.5 Jobs to be Done

- **Deskripsi:** Identify functional, emotional, and social jobs users are hiring solutions to accomplish.
- **Prompt Facilitation:**
  - "What job are they trying to do?"
  - "What progress do they want?"
  - "What are they really hiring this for?"
  - "What alternatives exist?"
- **Kapan Cocok:** Memahami motivasi di balik keputusan, inovasi produk, segmentasi.
- **Output:** Daftar jobs (functional/emotional/social) yang sedang "ditugaskan" ke solusi lain maupun produk kita.
- **Catatan:** Pengguna tidak membeli produk; mereka "memperkerjakan" produk untuk menyelesaikan job.

---

## 3. Ideate

> Fokus: **menghasilkan volume ide** — divergent thinking sebelum konvergen.

### 3.1 Brainstorming

- **Deskripsi:** Generate large quantity of diverse ideas without judgment to explore solution space fully.
- **Prompt Facilitation:**
  - "No bad ideas"
  - "Build on others"
  - "Go for quantity"
  - "Be visual"
  - "Stay on topic"
  - "Defer judgment"
- **Kapan Cocok:** Kick-off ideation, memperluas solusi space, keterlibatan tim.
- **Output:** Volume ide beragam yang akan di-cluster dan disaring pada langkah berikutnya.
- **Catatan:** Kuantitas awal, kualitas belakangan. Defer judgment untuk menjaga momentum kreatif.

### 3.2 Crazy 8s

- **Deskripsi:** Rapidly sketch eight solution variations in eight minutes to force quick creative thinking.
- **Prompt Facilitation:**
  - "Fold paper in 8"
  - "1 minute per sketch"
  - "No overthinking"
  - "Quantity over quality"
  - "Push past obvious"
- **Kapan Cocok:** Individu yang butuh dorongan kreatif, sesi pendek, eksplorasi cepat.
- **Output:** 8 sketsa kasar dalam 8 menit — memaksa otak melampaui solusi pertama.
- **Catatan:** Matikan inner critic. Sketsa jelek lebih berharga daripada ide sempurna yang tidak pernah keluar.

### 3.3 SCAMPER Design

- **Deskripsi:** Apply seven design lenses to existing solutions: Substitute, Combine, Adapt, Modify, Purposes, Eliminate, Reverse.
- **Prompt Facilitation:**
  - "What could we substitute?"
  - "How could we combine elements?"
  - "What could we adapt?"
  - "How could we modify it?"
  - "Other purposes?"
  - "What to eliminate?"
  - "What if reversed?"
- **Kapan Cocok:** Iterasi desain, modifikasi konsep existing, ideation terstruktur.
- **Output:** Tujuh transformasi alternatif dari solusi existing.
- **Catatan:** Setiap lensa membuka perspektif berbeda — pastikan semuanya dieksplorasi.

### 3.4 Provotype Sketching

- **Deskripsi:** Create deliberately provocative or extreme prototypes to spark breakthrough thinking.
- **Prompt Facilitation:**
  - "What's the most extreme version?"
  - "Make it ridiculous"
  - "Push boundaries"
  - "What useful insights emerge?"
- **Kapan Cocok:** Saat tim stuck di solusi mainstream, butuh breakthrough, eksplorasi avant-garde.
- **Output:** Prototipe "provokatif" yang sengaja berlebihan — tujuannya memicu insight, bukan jadi produk akhir.
- **Catatan:** Provotype bukan produk. Ekstrak insight-nya, lalu moderasi untuk jadi solusi realistis.

### 3.5 Analogous Inspiration

- **Deskripsi:** Find inspiration from completely different domains to spark innovative connections.
- **Prompt Facilitation:**
  - "What other field solves this?"
  - "How does nature handle this?"
  - "What's an analogous problem?"
  - "What can we borrow?"
- **Kapan Cocok:** Stuck di local maxima, butuh perspektif baru, inovasi cross-domain.
- **Output:** Analog dari domain lain yang bisa diadopsi atau diadaptasi untuk masalah kita.
- **Catatan:** Jangan terjebak di industri sendiri — museum, biologi, olahraga, militer adalah gudang analog yang kaya.

---

## 4. Prototype

> Fokus: **mewujudkan ide** — buat tangible untuk bisa diskusi dan uji.

### 4.1 Paper Prototyping

- **Deskripsi:** Create quick low-fidelity sketches and mockups to make ideas tangible for testing.
- **Prompt Facilitation:**
  - "Sketch it out"
  - "Make it rough"
  - "Focus on core concept"
  - "Test assumptions"
  - "Learn fast"
- **Kapan Cocok:** Sangat awal ideation, test flow, eksplorasi layout, usability preliminer.
- **Output:** Mockup kertas / wireframe kasar yang siap diuji dengan pengguna.
- **Catatan:** Fidelity rendah = gagal murah, iterasi cepat. Jangan over-detail di fase ini.

### 4.2 Role Playing

- **Deskripsi:** Act out user scenarios and service interactions to test experience flow and pain points.
- **Prompt Facilitation:**
  - "Play the user"
  - "Act out the scenario"
  - "What feels awkward?"
  - "Where does it break?"
  - "What works?"
- **Kapan Cocok:** Service design, multi-stakeholder experience, training tim, dry-run experience.
- **Output:** Identifikasi pain point pada flow dan alignment emosional anggota tim terhadap skenario.
- **Catatan:** Pura-pura pertama, authentic insight akan mengikuti setelah players terhanyut.

### 4.3 Wizard of Oz

- **Deskripsi:** Simulate complex functionality manually behind the scenes to test concept before building.
- **Prompt Facilitation:**
  - "Fake the backend"
  - "Focus on experience"
  - "What do they think is happening?"
  - "Does the concept work?"
- **Kapan Cocok:** Konsep radikal, biaya bangun tinggi, butuh validasi sebelum commit engineering.
- **Output:** Pengalaman pengguna yang "utuh" dengan backend yang sebenarnya masih manual.
- **Catatan:** Coba instrumentasi seminimal mungkin — yang penting concept validated, bukan sistem nyata.

### 4.4 Storyboarding

- **Deskripsi:** Visualize user experience across time and touchpoints as sequential illustrated narrative.
- **Prompt Facilitation:**
  - "What's scene 1?"
  - "How does it progress?"
  - "What's the emotional journey?"
  - "Where's the climax?"
  - "How does it resolve?"
- **Kapan Cocok:** Multi-touchpoint experience, alignment lintas fungsi, pitch internal.
- **Output:** Board visual cerita dari state awal → tengah → akhir, lengkap dengan emosinya.
- **Catatan:** Sketsa sederhana cukup. Yang penting narasi dan emosi terlihat — bukan seni.

### 4.5 Physical Mockups

- **Deskripsi:** Build tangible artifacts users can touch and interact with to test form and function.
- **Prompt Facilitation:**
  - "Make it 3D"
  - "Use basic materials"
  - "Make it interactive"
  - "Test ergonomics"
  - "Gather reactions"
- **Kapan Cocok:** Produk fisik, form factor decisions, ergonomic test, konsep retail/spatial.
- **Output:** Mockup 3D dari berbagai material yang bisa disentuh dan diuji langsung.
- **Catatan:** Kardus, clay, foam — material murah sudah cukup untuk dapat insight yang signifikan.

---

## 5. Test

> Fokus: **memvalidasi dengan pengguna** — feedback loop sebelum peluncuran.

### 5.1 Usability Testing

- **Deskripsi:** Watch users attempt tasks with prototype to identify friction points and opportunities.
- **Prompt Facilitation:**
  - "Try to accomplish X"
  - "Think aloud please"
  - "Don't help them"
  - "Where do they struggle?"
  - "What surprises them?"
- **Kapan Cocok:** Validasi flow, identifikasi bug UX, sebelum code, sebelum launch.
- **Output:** Daftar friction points, observed confusions, dan rekomendasi perbaikan.
- **Catatan:** Moderator harus menahan diri untuk tidak membantu — di sinilah insight terbesar muncul.

### 5.2 Feedback Capture Grid

- **Deskripsi:** Organize user feedback across likes, questions, ideas, and changes for actionable insights.
- **Prompt Facilitation:**
  - "What did they like?"
  - "What questions arose?"
  - "What ideas did they have?"
  - "What needs changing?"
- **Kapan Cocok:** Sesi feedback yang multi-peserta, retrospective iterasi desain.
- **Output:** Kuadran 4 (likes / questions / ideas / changes) yang siap diprioritaskan.
- **Catatan:** Metode ini juga berguna untuk menilai positioning proposal/arahan baru sebelum development.

### 5.3 A/B Testing

- **Deskripsi:** Compare two variations to understand which approach better serves user needs.
- **Prompt Facilitation:**
  - "Show version A"
  - "Show version B"
  - "Which works better?"
  - "Why the difference?"
  - "What does data show?"
- **Kapan Cocok:** Hipotesis terukur, optimasi on-page, decision yang perlu data objektif.
- **Output:** Statistik signifikan dari performa dua varian + rekomendasi versi lanjutan.
- **Catatan:** Hanya ukur satu variabel per test untuk isolasi sebab-akibat.

### 5.4 Assumption Testing

- **Deskripsi:** Identify and validate critical assumptions underlying your solution to reduce risk.
- **Prompt Facilitation:**
  - "What are we assuming?"
  - "How can we test this?"
  - "What would prove us wrong?"
  - "What's the riskiest assumption?"
- **Kapan Cocok:** Pre-build, pre-launch, eksperimen murah untuk validasi asumsi.
- **Output:** Map of key assumptions + status (validated / invalid / untested) + tes berikutnya.
- **Catatan:** Mulai dari asumsi paling berisiko — bukan yang paling mudah diuji.

### 5.5 Iterate and Refine

- **Deskripsi:** Use test insights to improve prototype through rapid cycles of refinement and re-testing.
- **Prompt Facilitation:**
  - "What did we learn?"
  - "What needs fixing?"
  - "What stays?"
  - "Make changes quickly"
  - "Test again"
- **Kapan Cocok:** Paska usability test, iterasi berkala, continuous improvement.
- **Output:** Prototype versi baru yang siap diuji ulang dengan perubahan terarah.
- **Catatan:** Setiap iterasi harus punya satu pertanyaan desain jelas. Jangan iterasi tanpa arah.

---

## 6. Implement

> Fokus: **membawa ke dunia nyata** — launched, monitored, dan disempurnakan.

### 6.1 Pilot Programs

- **Deskripsi:** Launch small-scale real-world implementation to learn before full rollout.
- **Prompt Facilitation:**
  - "Start small"
  - "Real users"
  - "Real context"
  - "What breaks?"
  - "What works?"
  - "Scale lessons learned"
- **Kapan Cocok:** Validasi real-world, mitigasi risiko launch, market entry bertahap.
- **Output:** Lessons-learned dari pilot + go/no-go criteria untuk scale-up.
- **Catatan:** Pilih pilot audience yang merepresentasikan target pasar — bukan yang paling mudah dicapai.

### 6.2 Service Blueprinting

- **Deskripsi:** Map all service components, interactions, and touchpoints to guide implementation.
- **Prompt Facilitation:**
  - "What's visible to users?"
  - "What happens backstage?"
  - "What systems are needed?"
  - "Where are handoffs?"
- **Kapan Cocok:** Implementasi layanan, alignment antar tim, ekspos failure points.
- **Output:** Blueprint visual garis depan (visible) vs garis belakang (backstage) + supporting processes.
- **Catatan:** Handoff adalah titik failure paling umum — bingkai secara eksplisit siapa pegang apa kapan.

### 6.3 Design System Creation

- **Deskripsi:** Build consistent patterns, components, and guidelines for scalable implementation.
- **Prompt Facilitation:**
  - "What patterns repeat?"
  - "Create reusable components"
  - "Document standards"
  - "Enable consistency"
- **Kapan Cocok:** Skala tim, multi-product, konsistensi lintas touchpoint, efisiensi design-dev.
- **Output:** Library komponen + token + guidelines yang dipakai lintas tim dan produk.
- **Catatan:** Design system yang hidup > design system yang "launched and forgotten". Maintain terus.

### 6.4 Stakeholder Alignment

- **Deskripsi:** Bring team and stakeholders along journey to build shared understanding and commitment.
- **Prompt Facilitation:**
  - "Show the research"
  - "Walk through prototypes"
  - "Share user stories"
  - "Build empathy"
  - "Get buy-in"
- **Kapan Cocok:** Sebelum implementasi besar, proyek lintas fungsi, banyak decision makers.
- **Output:** Tim yang tidak hanya tahu, tapi juga *mengalami* dan *membeli* visi desain.
- **Catatan:** Stakeholder yang diajak empather dulu akan menjadi sekutu — bukan penghalang.

### 6.5 Measurement Framework

- **Deskripsi:** Define success metrics and feedback loops to track impact and inform future iterations.
- **Prompt Facilitation:**
  - "How will we measure success?"
  - "What are key metrics?"
  - "How do we gather feedback?"
  - "When do we revisit?"
- **Kapan Cocok:** Pre-launch, post-launch review, continuous learning loop.
- **Output:** Dashboard metrics + review cadence + decision criteria untuk iterasi berikutnya.
- **Catatan:** Pilih metrics yang actionable — bukan vanity metrics. Definisi "success" harus jelas di awal.

---

## Perbandingan dengan Solving Methods

| Aspek | Design Methods | Solving Methods |
|---|---|---|
| **Jumlah** | 30 metode | 25 metode |
| **Fase** | 6 (Empathize, Define, Ideate, Prototype, Test, Implement) | 5 (Diagnosis, Analysis, Synthesis, Evaluation, Implementation) |
| **Titik awal** | Empati terhadap pengguna | Definisi masalah |
| **Fokus utama** | User-centered, experientially grounded | Problem-centric, analitis |
| **Peran pengguna** | Sentral — di setiap fase | Subjek verifikasi |
| **Output khas** | Solusi yang delight + feasible | Solusi yang valid + terukur |
| **Cocok untuk** | Inovasi produk/layanan, eksplorasi | Optimasi sistem, troubleshooting, strategi |
| **Iterasi** | Antar fase (prototype ↔ test) | Biasanya linear; revisit hanya saat ada anomali |

### Kapan Memilih Design vs Solving

| Situasi | Pilih |
|---|---|
| Belum jelas siapa user atau apa masalah mereka | **Design Methods** (mulai dari Empathize) |
| Masalah sudah terdefinisi dengan baik, perlu solusi optimal | **Solving Methods** (mulai dari Analysis) |
| Produk baru / fitur baru / market baru | **Design Methods** |
| Optimasi sistem / layanan yang sudah jalan | **Solving Methods** |
| Kombinasi keduanya | **Gunakan keduanya**: Design Methods untuk explore, Solving Methods untuk optimize |

---

## Matriks Pemilihan

| Fase | Outcome Utama | Pilih Metode Jika… |
|---|---|---|
| **Empathize** | Pemahaman user yang kaya | Butuh insight mendalam; riset kualitatif |
| **Define** | Problem statement yang tajam | Butuh framing; alignment tim |
| **Ideate** | Volume ide yang beragam | Butuh divergent thinking; eksplorasi |
| **Prototype** | Ide yang tangible | Butuh bentuk untuk diskusi & test |
| **Test** | Validasi empiris | Butuh feedback loop sebelum commit |
| **Implement** | Real-world impact | Butuh peluncuran yang aman & terukur |

---

## Data CSV

Data lengkap 30 metode tersedia dalam format CSV di:
```
add-on/method/design/design-methods.csv
```

Kolom CSV:
- `phase` — fase design thinking (empathize / define / ideate / prototype / test / implement)
- `method_name` — nama metode
- `description` — deskripsi singkat
- `facilitation_prompts` — daftar pertanyaan / prompt (dipisah dengan `|`)

CSV ini bisa dipakai langsung oleh AI agent untuk melakukan facilitation otomatis pada tiap fase.

---

*Methods ini mengikuti alur Design Thinking klasik — dari empati hingga implementasi — dengan iterate loop antar fase. Gunakan untuk inovasi human-centered, apakah produk digital, layanan, atau sistem.*
