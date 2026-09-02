# Step 2: Generate Persona

> Membuat folder dan file persona berdasarkan blueprint yang dipilih.
> Output disimpan di `.ai-doc/personas/<nama-persona>/`.

---

## Input dari Step 1

```yaml
persona:
  blueprint: "problem-solver"      # blueprint yang dipilih
  persona_name: "dr-quinn"         # slug untuk folder name
  persona_label: "Dr. Quinn"       # display name
  kebutuhan: "..."                 # catatan kebutuhan user
```

---

## Yang Dilakukan Agent

### 1. Cek Folder `.ai-doc/personas/`

- Jika belum ada: buat folder `.ai-doc/personas/`
- Jika sudah ada: lanjut

### 2. Generate Folder Persona

Buat struktur folder:

```
.ai-doc/personas/<nama-persona>/
├── persona.md
├── customize.toml
└── methods.md
```

### 3. Generate `customize.toml`

Baca blueprint `add-on/persona/blueprints/<blueprint>/customize.toml` dan sesuaikan. **File ini adalah source of truth konfigurasi persona**.

- Isi `[persona]`:
  - `slug` → `persona_name`
  - `blueprint` → blueprint yang dipilih
  - `archetype` → archetype blueprint
- Isi `[identity]`:
  - `user_name` → gunakan `{user_name}` sebagai placeholder (akan diisi saat greeting)
  - `greeting_template` → sesuaikan dengan archetype dan gaya komunikasi persona
- Isi `[agent]`:
  - `name` → `persona_label`
  - Biarkan `title`, `icon`, `role`, prinsip, dan metode lain sesuai blueprint
- Untuk blueprint **Custom**: generate file dengan semua field kosong/komentar agar user mengisi sendiri

### 4. Generate `persona.md`

Baca blueprint `add-on/persona/blueprints/<blueprint>/persona.md` dan salin dengan penyesuaian:

- Ubah judul jadi "Persona: {{persona_label}}"
- Ubah nama default di tabel jika ada dengan nama yang dipilih user
- Untuk blueprint **Custom**: kosongkan semua field, user isi sendiri nanti
- Bagian **Load Config** jangan berisi YAML duplikat. Isi dengan referensi ke `customize.toml` sebagai source of truth:
  - `persona.slug` → `[persona].slug`
  - `persona.blueprint` → `[persona].blueprint`
  - `persona.archetype` → `[persona].archetype`
  - `agent.name` → `[agent].name`
  - `identity.user_name` → `[identity].user_name`
  - `identity.greeting_template` → `[identity].greeting_template`
- Isi **Metode yang Dikuasai** — referensi ke `methods.md`:
  - Tulis ringkasan daftar metode per kategori (nama saja, tanpa detail)
  - Tambahkan link referensi ke `methods.md` untuk detail lengkap
  - Format:

    ```markdown
    ## Metode yang Dikuasai

    Detail metode lengkap ada di [methods.md](methods.md).

    | Kategori | Metode |
    |---|---|
    | **🧠 Brain Method** | First Principles Thinking, Concept Map, ... |
    | **🔧 Solving Method** | Systems Thinking, Gap Analysis, ... |
    | **🚀 Innovation Method** | Technology Roadmapping, ... |
    ```

### 5. Generate `methods.md`

Buat file `methods.md` di folder persona yang berisi detail lengkap semua metode yang dikuasai:

- Baca daftar metode dari blueprint `persona.md` section "Metode yang Dikuasai"
- Untuk setiap metode, buka file referensi di `add-on/method/`:
  - Brain Methods → `add-on/method/brain/README.md`
  - Solving Methods → `add-on/method/solving/README.md`
  - Innovation Frameworks → `add-on/method/innovation/README.md`
- Untuk setiap metode, tulis detail lengkap:
  - **Deskripsi** — penjelasan singkat metode
  - **Prompt Utama / Prompt Facilitation / Key Questions** — pertanyaan panduan
  - **Kapan Cocok** — situasi terbaik penggunaan
  - **Catatan** — tips tambahan
- Susun dalam 3 sub-section: 🧠 Brain Methods, 🔧 Solving Methods, 🚀 Innovation Frameworks
- Format contoh untuk Brain Method:

  ```markdown
  **First Principles Thinking**
  - **Kategori:** Deep
  - **Deskripsi:** Buang semua asumsi, bangun ulang dari kebenaran fundamental.
  - **Prompt Utama:**
    - "Apa yang kita tahu dengan pasti?"
    - "Apa kebenaran fundamentalnya?"
    - "Kalau kita mulai dari nol?"
  - **Kapan Cocok:** Inovasi radikal, tantangan status quo.
  - **Catatan:** Teknik favorit Elon Musk. Butuh waktu tapi hasilnya powerful.
  ```

- Format contoh untuk Solving Method:

  ```markdown
  **Five Whys Root Cause**
  - **Kategori:** Diagnosis
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
  ```

- Format contoh untuk Innovation Framework:

  ```markdown
  **Disruptive Innovation Theory**
  - **Kategori:** Disruption
  - **Deskripsi:** Identifikasi bagaimana new entrants menggunakan solusi lebih sederhana dan murah untuk mengalahkan incumbents.
  - **Key Questions:**
    - "Siapa non-consumers-nya?"
    - "Apa yang 'cukup baik' untuk mereka?"
    - "Kelemahan incumbent apa yang ada?"
  - **Kapan Cocok:** Evaluasi potensi disruptif, analisis kompetitor, startup strategy.
  - **Output:** Disruption opportunity assessment.
  - **Catatan:** Disruption ≠ teknologi baru. Tentang business model yang lebih sederhana dan accessible.
  ```

### 6. Update atau Buat `list.md`

Cek apakah `.ai-doc/personas/list.md` sudah ada:

- Jika belum ada: buat dengan format berikut
- Jika sudah ada: tambahkan baris baru ke tabel

**Format `list.md`:**

```markdown
# Daftar Persona

Berikut persona yang tersedia untuk sesi AI Documentor:

| Nama | Blueprint | Archetype | Kegunaan | Metode |
|---|---|---|---|---|
| **Dr. Quinn** | problem-solver | The Diagnostician | Root cause analysis, debugging | [methods.md](dr-quinn/methods.md) |

## Cara Pakai

Persona bisa dipanggil di sesi brainstorming atau add-on lain dengan format:
- `+Dr. Quinn` — untuk debugging
```

### 7. Konfirmasi ke User

Setelah semua file tergenerate, tampilkan ringkasan:

```
✅ Persona "{{persona_label}}" berhasil dibuat!

  📁 .ai-doc/personas/{{persona_name}}/
  ├── 📄 persona.md
  ├── 📄 customize.toml
  └── 📄 methods.md

Blueprint: {{blueprint}}
Archetype: {{archetype}}
Metode: {{daftar metode}}

Persona siap digunakan! Anda bisa:
  - Panggil dengan: "Gunakan persona {{persona_name}}"
  - Edit sendiri file customize.toml untuk ubah prinsip/metode
  - Lihat methods.md untuk detail lengkap metode yang dikuasai
  - Gunakan di add-on brainstorming sebagai peserta
```

---

## Output Step 2

```yaml
generated:
  persona_path: ".ai-doc/personas/{{persona_name}}/"
  files:
    - "persona.md"
    - "customize.toml"
    - "methods.md"
  list_updated: true
  status: "ready"
```

---

## Aturan

- **JANGAN ubah blueprint asli** di `add-on/persona/blueprints/`. Blueprint adalah template master — tidak boleh dimodifikasi.
- **JANGAN overwrite** file yang sudah ada — jika folder `.ai-doc/personas/<nama>/` sudah ada, beri tahu user dan tawarkan:
  - Pakai nama lain
  - Overwrite (hanya jika user setuju)
- **PASTIKAN** `list.md` selalu terupdate setelah generate.
- **SIMPAN** semua output di `.ai-doc/personas/` — jangan di tempat lain.
- **KEMBALIKAN** kontrol ke core AI Documentor setelah selesai.
