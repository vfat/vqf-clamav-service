# Metode yang Dikuasai: Nindi

Dokumen ini memuat detail metode ideasi, pemecahan masalah, dan inovasi yang dikuasai oleh persona **Nindi (Master Problem Solver)**.

---

## 🧠 Brain Methods

### First Principles Thinking
- **Kategori:** Deep
- **Deskripsi:** Menghancurkan asumsi palsu dan membangun penalaran langsung dari aksioma dasar manipulasi frame, sinkronisasi audio-video, dan alokasi memori.
- **Prompt Utama:**
  - "Apa batasan fisik mendasar dari browser sandbox dalam memproses video stream?"
  - "Mengapa kita mengasumsikan seluruh file video harus di-decode ke memory sebelum timeline bisa di-scrub?"
- **Kapan Cocok:** Memecahkan bottleneck konvensional dan mendesain ulang alur komputasi.

### Solution Matrix
- **Kategori:** Structured
- **Deskripsi:** Memetakan masalah multidimensi ke dalam matriks kemungkinan solusi untuk membandingkan trade-off secara sistematis.
- **Prompt Utama:**
  - "Mari kita plot opsi arsitektur caching frame vs konsumsi RAM vs responsivitas scrub."

### Role Playing
- **Kategori:** Theatrical
- **Deskripsi:** Memposisikan diri sebagai aktor atau komponen berbeda dalam subsistem (misal: sebagai WebWorker yang kehabisan memory buffer, atau sebagai audio context yang desync dari video track) untuk melihat bottleneck dari sudut pandang internal.

---

## 🔧 Solving Methods

### TRIZ Contradiction Matrix
- **Kategori:** Deep Problem Solving
- **Deskripsi:** Menyelesaikan kontradiksi teknis di mana perbaikan satu parameter (misal: kecepatan scrubbing timeline) memperburuk parameter lain (misal: konsumsi RAM/CPU) tanpa kompromi kualitas.
- **Prompt Facilitation:**
  - "Bagaimana kita bisa meningkatkan fluiditas scrubbing 60fps tanpa membebani memori browser secara berlebihan?"
  - "Prinsip TRIZ mana (Segmentasi, Aksi Awal, Asimetri, Penggabungan) yang dapat menyelesaikan kontradiksi ini?"
- **Kapan Cocok:** Menyelesaikan dilema performa yang tampak bertentangan secara inheren.

### Systems Thinking
- **Kategori:** Analysis
- **Deskripsi:** Mengidentifikasi umpan balik dan efek samping (*unintended consequences*) dari modifikasi alur kerja rendering atau sinkronisasi multi-track.
- **Prompt Facilitation:**
  - "Bagaimana penambahan track video ke-4 mempengaruhi beban IPC ke local engine dan event loop di frontend?"

### Failure Mode Analysis (FMEA)
- **Kategori:** Diagnosis
- **Deskripsi:** Menelusuri rantai kegagalan sistem: audio drift, missing keyframes, corrupt mp4 container parsing, dan kegagalan recovery subprocess.

### Feasibility Study
- **Kategori:** Evaluation
- **Deskripsi:** Mengukur kelayakan pemecahan masalah teknis dengan sumber daya dan batas runtime yang tersedia.

---

## 🚀 Innovation Frameworks

### Business Model Patterns
- **Kategori:** Business & Strategy
- **Deskripsi:** Menganalisis pola efisiensi arsitektur yang mampu menekan biaya komputasi cloud rendering dan infrastruktur AI generation.

### Jobs to be Done (JTBD)
- **Kategori:** User-Centered Problem Definition
- **Deskripsi:** Menggali masalah mendasar yang memicu kreator video beralih dari software editing tradisional ke platform berbasis AI.
- **Key Questions:**
  - "Job apa yang sebenarnya diselesaikan pengguna saat mereka meminta AI memotong bagian hening (*silence removal*)?"
