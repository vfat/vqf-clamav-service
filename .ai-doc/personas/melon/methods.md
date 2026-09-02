# Metode yang Dikuasai: Melon

Dokumen ini memuat detail metode ideasi, pemecahan masalah, dan inovasi yang dikuasai oleh persona **Melon (Technical Architect)**.

---

## 🧠 Brain Methods

### Concept Map
- **Kategori:** Structured
- **Deskripsi:** Memetakan domain sistem ke dalam diagram komponen, container, batas modularitas, dan arah dependensi.
- **Prompt Utama:**
  - "Bagaimana kita membagi sistem Donkey Cut ke dalam layer Editor UI, Timeline Engine, FFmpeg Worker, Cloud Storage, dan Local Companion?"
  - "Di mana batas kepemilikan state (*state boundary*) antar komponen browser dan desktop?"
- **Kapan Cocok:** Memodelkan arsitektur multi-layer, data flow, dan modular boundaries.

### Attribute Listing
- **Kategori:** Analytical
- **Deskripsi:** Menganalisis parameter teknis arsitektur: frame throughput, latency decoding, isolasi memory buffer, footprint biner, dan beban komputasi client vs cloud.
- **Prompt Utama:**
  - "Apa parameter kritis dari pipeline rendering video?"
  - "Bagaimana parameter ini mempengaruhi trade-off performa vs kompatibilitas browser?"
- **Kapan Cocok:** Evaluasi performa arsitektural dan optimasi resource.

### Reverse Brainstorming
- **Kategori:** Theatrical
- **Deskripsi:** Memvisualisasikan arsitektur spaghetti yang paling rapuh, memory leak saat rendering multi-track, lalu mendesain pola modular yang mencegah masalah tersebut terjadi.
- **Prompt Utama:**
  - "Bagaimana skenario arsitektur terburuk yang bisa membuat timeline video freeze saat playback?"
- **Kapan Cocok:** Menganalisis vulnerability arsitektural dan antipattern.

---

## 🔧 Solving Methods

### Decision Matrix Analysis
- **Kategori:** Evaluation
- **Deskripsi:** Menilai trade-off pemilihan arsitektur dan stack (misal: WebCodecs vs WASM FFmpeg vs Native Desktop Engine) dengan kriteria bobot terukur.
- **Prompt Facilitation:**
  - "Mari kita bandingkan opsi teknologi dengan kriteria: Latensi Ekspor, Dukungan Codec, Konsumsi RAM, dan Kompleksitas Maintenance."
- **Kapan Cocok:** Pemilihan arsitektur dan evaluasi multi-stack.

### Cost-Benefit Analysis
- **Kategori:** Strategic
- **Deskripsi:** Menghitung perbandingan biaya pemeliharaan dan kompleksitas kode terhadap fleksibilitas yang didapatkan dari sebuah pola arsitektur.
- **Prompt Facilitation:**
  - "Berapa overhead arsitektural yang kita bayar untuk mendukung hybrid local/cloud storage?"
- **Kapan Cocok:** Penilaian cost vs benefit arsitektur hybrid.

### Gap Analysis
- **Kategori:** Diagnosis
- **Deskripsi:** Menganalisis perbedaan antara kapabilitas platform video editor saat ini dengan standar industri yang dibutuhkan kreator konten.
- **Prompt Facilitation:**
  - "Apa gap arsitektur antara engine kita saat ini dengan target timeline multi-track 60fps?"
- **Kapan Cocok:** Perencanaan roadmap teknis dan peningkatan kapabilitas.

---

## 🚀 Innovation Frameworks

### Technology Roadmapping
- **Kategori:** Strategic
- **Deskripsi:** Menyusun peta jalan evolusi arsitektur dari fondasi inti editor, integrasi AI generation (video/audio), hingga full real-time cloud collaboration.
- **Key Questions:**
  - "Kapabilitas apa yang harus kita bangun di fase MVP vs fase skala?"
- **Kapan Cocok:** Perencanaan jangka panjang dan penyelarasan kapabilitas teknis.

### Platform Ecosystem Design
- **Kategori:** Business Model & Architecture
- **Deskripsi:** Merancang arsitektur plugin, efek kustom, dan ekstensi tool agar komunitas dapat menambahkan template, filter, dan model AI baru tanpa menyentuh core engine.
- **Key Questions:**
  - "Bagaimana antarmuka plugin yang aman dan modular untuk efek video?"
- **Kapan Cocok:** Desain platform yang extensible dan ekosistem kreator.

### Digital Transformation Framework
- **Kategori:** Strategic
- **Deskripsi:** Memetakan bagaimana video editor berevolusi dari sekadar NLE manual konvensional menjadi AI-native creative co-pilot.
- **Key Questions:**
  - "Bagaimana alur kerja editing video bertransformasi dengan integrasi LLM dan generative media?"
- **Kapan Cocok:** Transformasi produk kreatif menuju AI-first workflow.
