# Brainstorming Add-On — Techniques

> Daftar teknik brainstorming untuk 2 area fokus: Troubleshooting dan Feature/Release.
> Setiap area memiliki 4 kategori: Diagnosis/Ideation → Analysis/Perspective → Solution/Planning → Innovation.
> Ini adalah versi Markdown deklaratif — LLM-friendly, tanpa CSV.

---

## Mapping Area × Kategori

| Area | Diagnosis / Ideation | Analysis / Perspective | Solution / Planning | Innovation |
|---|---|---|---|---|
| **Troubleshooting** | Cari akar masalah | Analisis mendalam | Solusi perbaikan | Cara baru yang radikal |
| **Feature / Release** | Hasilkan ide baru | Evaluasi multi-sisi | Perencanaan release | Terobosan produk |

---

## Area 1: Troubleshooting / Bug Fixing / Hotfix / Improvement

> Fokus: **memperbaiki yang sudah ada** — dari diagnosis sampai solusi inovatif.

### Diagnosis

| # | Teknik | Kapan Cocok | Prompt Utama |
|---|--------|-------------|--------------|
| 1 | **Five Whys** | Root cause tidak jelas, bug berulang | "Kenapa ini terjadi?" (ulangi 5×) |
| 2 | **Fishbone Diagram** | Bug kompleks dengan banyak potensi penyebab | "Faktor apa saja yang berkontribusi?" |

### Analysis

| # | Teknik | Kapan Cocok | Prompt Utama |
|---|--------|-------------|--------------|
| 3 | **Failure Analysis** | Post-mortem, belajar dari bug sebelumnya | "Apa yang salah? Kenapa gagal?" |
| 4 | **Assumption Reversal** | Asumsi tersembunyi yang salah | "Asumsi apa yang kita buat? Kalau sebaliknya?" |
| 5 | **Constraint Mapping** | Improvement dengan resource terbatas | "Batasan apa yang real vs imagined?" |

### Solution

| # | Teknik | Kapan Cocok | Prompt Utama |
|---|--------|-------------|--------------|
| 6 | **Reversal Inversion** | Solusi konvensional tidak cukup | "Bagaimana cara membuat ini lebih buruk?" |
| 7 | **SCAMPER** | Improvement fitur/proses existing | "Apa yang bisa di-eliminate/modify/combine?" |
| 8 | **Pareto Analysis** | Banyak bug, prioritas mana dulu | "20% penyebab = 80% masalah — fokus di mana?" |

### Innovation

| # | Teknik | Kapan Cocok | Prompt Utama |
|---|--------|-------------|--------------|
| 9 | **First Principles** | Perlu pendekatan radikal, bukan patch | "Kalau mulai dari nol, bagaimana?" |
| 10 | **Analogical Thinking** | Inspirasi dari sistem/domain lain | "Sistem lain menyelesaikan ini bagaimana?" |

---

## Area 2: Feature / Release

> Fokus: **membuat yang baru** — dari ideasi sampai perencanaan release.

### Ideation

| # | Teknik | Kapan Cocok | Prompt Utama |
|---|--------|-------------|--------------|
| 11 | **What If Scenarios** | Butuh ide fitur out-of-the-box | "Bagaimana kalau user bisa...?" |
| 12 | **Yes And Building** | Rapid ideation, bangun di atas ide sebelumnya | "Yes, dan kita juga bisa..." |

### Perspective

| # | Teknik | Kapan Cocok | Prompt Utama |
|---|--------|-------------|--------------|
| 13 | **Role Playing** | Multi-stakeholder perspective | "Kalau kamu [user/admin/dev], apa yang kamu mau?" |
| 14 | **Six Thinking Hats** | Evaluasi fitur dari semua sisi | "Fakta? Emosi? Manfaat? Risiko? Alternatif?" |

### Planning

| # | Teknik | Kapan Cocok | Prompt Utama |
|---|--------|-------------|--------------|
| 15 | **Mind Mapping** | Organisasi ide fitur, release planning | "Fitur utama → sub-fitur → task" |
| 16 | **Resource Constraints** | MVP, fitur dengan resource terbatas | "Kalau cuma 1 minggu, apa yang paling penting?" |

### Innovation

| # | Teknik | Kapan Cocok | Prompt Utama |
|---|--------|-------------|--------------|
| 17 | **First Principles** | Fitur yang "selalu dilakukan begini" | "Kalau mulai dari nol, bagaimana?" |
| 18 | **Analogical Thinking** | Inspirasi dari produk/domain lain | "Produk lain menyelesaikan ini bagaimana?" |

---

## Cara Penggunaan

1. **Tentukan area** terlebih dahulu (Troubleshooting atau Feature/Release).
2. **Pilih 1 teknik per kategori** — AI rekomendasikan, user konfirmasi.
3. **Gunakan prompt utama** sebagai pembuka sesi untuk teknik tersebut.
4. **Setelah selesai** dengan satu teknik, lanjut ke teknik di kategori berikutnya.
5. **Innovation** opsional — hanya jika user ingin eksplorasi radikal.

### Contoh Sesi Troubleshooting

```
Topik: "Login sering timeout di jam sibuk"
Teknik: Five Whys (diagnosis) → Fishbone (analysis) → First Principles (innovation)
Output MoM: Root cause = DB connection pool terlalu kecil
Output Discussion: 3 opsi fix — patch vs redesign → dipilih redesign
```

### Contoh Sesi Feature/Release

```
Topik: "Fitur notifikasi untuk v2.0"
Teknik: What If (ideation) → Role Playing (perspective) → Mind Mapping (planning) → First Principles (innovation)
Output MoM: 12 ide notifikasi, dipilih 5 untuk v2.0
Output Discussion: Detail per fitur, prioritas, dan estimasi
```

---

## Catatan

- Teknik bersifat **rekomendasi**. User boleh menyesuaikan atau mengganti.
- Boleh menggunakan teknik dari area berbeda jika relevan.
- Innovation hanya dipakai jika user ingin eksplorasi di luar perbaikan/penambahan biasa.
- Teknik bisa dikombinasikan secara berurutan (tidak paralel).
