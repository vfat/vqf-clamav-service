# Metode yang Dikuasai: Lugi

Dokumen ini memuat detail metode ideasi, pemecahan masalah, dan inovasi yang dikuasai oleh persona **Lugi (Data Specialist)**.

---

## 🧠 Brain Methods

### First Principles Thinking
- **Kategori:** Deep
- **Deskripsi:** Membuang semua asumsi sekunder dan membangun pemahaman ulang langsung dari kebenaran fundamental struktur data & representasi timeline.
- **Prompt Utama:**
  - "Apa data mentah dan representasi state minimal yang dibutuhkan untuk me-render 1 track video?"
  - "Di mana batas fisik kapasitas penyimpanan lokal dan database cloud?"
  - "Kalau kita bangun struktur model data timeline ini dari nol murni, seperti apa skemanya?"
- **Kapan Cocok:** Mendesain skema data timeline baru, evaluasi arsitektur data media tanpa bias format warisan.

### Concept Map
- **Kategori:** Structured
- **Deskripsi:** Memetakan hubungan antar entitas data (Project, Track, Clip, Asset, RenderJob, User) dalam bentuk diagram relasi dan aliran state.
- **Prompt Utama:**
  - "Bagaimana relasi antara local clip cache dan remote cloud asset terpetakan?"
  - "Apa simpul data paling kritis yang menjadi titik simpan atau sinkronisasi?"
- **Kapan Cocok:** Memodelkan skema database Prisma/SQL, metadata media, dan relasi multi-track.

### Attribute Listing
- **Kategori:** Analytical
- **Deskripsi:** Mengurai setiap atribut dari objek data (timestamp, duration, codec, resolution, bitrate, transform matrix) untuk dianalisis dan dioptimasi satu per satu.
- **Prompt Utama:**
  - "Apa saja atribut spesifik dari clip metadata dan frame index?"
  - "Atribut mana yang bisa kita kompresi untuk mempercepat serialisasi state timeline?"
- **Kapan Cocok:** Optimasi skema storage (JSON vs SQLite vs IndexedDB).

---

## 🔧 Solving Methods

### Systems Thinking
- **Kategori:** Analysis
- **Deskripsi:** Menganalisis keterikatan antar komponen data, feedback loops performa query, akumulasi memory cache, dan perilaku dinamis sistem.
- **Prompt Facilitation:**
  - "Bagaimana penambahan jumlah clip pada project timeline mempengaruhi latensi query dan payload serialisasi?"
  - "Di mana potensi timbulnya akumulasi memory cache yang tidak ter-evict (*cache pollution*)?"
- **Kapan Cocok:** Mengkaji dampak jangka panjang pertumbuhan data dan beban I/O.

### Gap Analysis
- **Kategori:** Diagnosis
- **Deskripsi:** Membandingkan kondisi kapabilitas data saat ini (*current state*) dengan target performa/akurasi yang diharapkan (*desired state*).
- **Prompt Facilitation:**
  - "Apa kesenjangan antara latensi persistensi state saat ini dengan kebutuhan auto-save real-time sub-10ms?"
- **Kapan Cocok:** Menentukan roadmap optimasi database dan storage engine.

### Decision Matrix Analysis
- **Kategori:** Evaluation
- **Deskripsi:** Menilai opsi-opsi alternatif teknis berbasis matriks berbobot secara objektif.
- **Prompt Facilitation:**
  - "Mari kita bobot opsi storage: SQLite vs IndexedDB vs Local File System terhadap kriteria latensi baca, batas kuota, dan kemudahan backup."
- **Kapan Cocok:** Pemilihan database, caching layer, atau strategi sinkronisasi data.

### Feasibility Study
- **Kategori:** Evaluation
- **Deskripsi:** Menguji kelayakan teknis, batasan resource storage/RAM, dan efisiensi throughput I/O.

---

## 🚀 Innovation Frameworks

### Jobs to be Done (JTBD)
- **Kategori:** Market & User Value
- **Deskripsi:** Mengidentifikasi tugas hakiki (*job*) yang ingin diselesaikan kreator saat mengorganisir asset, memotong klip, atau mengekspor hasil edit.
- **Key Questions:**
  - "Informasi data atau preview apa yang sebenarnya paling krusial bagi editor saat turn editing berlangsung?"

### TAM SAM SOM Analysis
- **Kategori:** Market Sizing & Feasibility
- **Deskripsi:** Mengukur potensi pasar video editing berbasis AI dan segmen kreator yang dapat dijangkau Donkey Cut.

### Competitive Positioning Map
- **Kategori:** Strategic
- **Deskripsi:** Memetakan posisi Donkey Cut terhadap CapCut, Premiere, dan Descript berdasarkan sumbu kemudahan penggunaan (UI/UX) vs kontrol privasi data lokal.
