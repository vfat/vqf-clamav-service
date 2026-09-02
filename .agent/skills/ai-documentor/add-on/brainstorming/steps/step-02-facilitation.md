# Step 2: Interactive Facilitation

> Loop per sub-topik. Agent memfasilitasi, persona merespons (jika aktif), user memutuskan.
> **HALT** di setiap titik keputusan.

---

## Alur per Sub-Topik

```
Untuk setiap sub-topik (urut dari step-01-setup):
  │
  ├─ 1. Tampilkan sub-topik + kategori + teknik yang dipakai
  ├─ 2. Agent lempar prompt sesuai teknik
  ├─ 3. Persona respons dengan ide (jika persona-driven)
  │     └─ Jika tidak persona-driven → agent jawab sendiri
  ├─ 4. User tanggapi:
  │     ├─ Approve → ide masuk ke list (status: approved)
  │     ├─ Reject → ide di-skip (status: rejected)
  │     └─ Modifikasi → ide di-update (status: modified)
  ├─ 5. Catat ide (append-only + status):
  │     [Ide #X]: Judul Singkat
  │     Detail: ...
  │     Dari: user / joni / jono / agent
  │     Status: approved / rejected / modified
  ├─ 6. HALT → "Lanjut / Ganti sub-topik / Selesai?"
  │     ├─ Lanjut → ulangi facilitation (poin 2-5)
  │     ├─ Ganti sub-topik → pilih sub-topik berikutnya
  │     └─ Selesai → Step 3 (wrap-up)
  └─ 7. Loop sampai user bilang "selesai" atau semua sub-topik selesai
```

---

## Aturan Facilitation

### 1. Agent = Facilitator

- Bukan generator ide utama — kecuali tidak ada persona aktif.
- Tugas agent: melempar prompt, mencatat, menjaga fokus.
- Jangan mendominasi diskusi dengan ide sendiri jika ada persona aktif.

### 2. Satu Elemen Teknik Per Interaksi

Satu prompt = satu pertanyaan/perspektif. Jangan lempar 3 pertanyaan sekaligus.

**Contoh baik:**
> "Oke, kita pakai Five Whys. Kenapa kira-kira login sering timeout?"

**Contoh buruk:**
> "Dari Five Whys, Fishbone, dan SCAMPER, mana yang mau dipakai? Atau mungkin kita langsung ke First Principles?"

### 3. HALT di Setiap Titik Keputusan

Setelah setiap ide direspons user, tanyakan:
> "Lanjut eksplorasi sub-topik ini? Ganti sub-topik? Atau selesai?"

**JANGAN auto-proceed.**

### 4. Anti-Bias

Setiap ~10 ide, lakukan pivot domain:
> "Kita sudah dapat 10 ide. Coba lihat dari sudut pandang lain — bagaimana kalau dari sisi [domain berbeda]?"

### 5. Persona Responses

Jika persona aktif, persona merespons **sesuai karakteristiknya**:

| Persona | Karakter Respons |
|---|---|
| **Joni** (Eksploratif) | Memunculkan banyak kemungkinan, tidak cepat mengerucut, suka analogi dan skenario "what if" |
| **Jono** (Analitis) | Menantang asumsi, bertanya sistematis, mencari bukti dan hubungan sebab-akibat |

Karakteristik detail ada di `.ai-doc/personas/list.md`.

### 6. Idea Status Tracking

Setiap ide dicatat dengan format:

```
[Ide #X]: Judul Singkat
  Detail: ...
  Dari: user / joni / jono / agent
  Status: approved / rejected / modified
```

- **approved**: ide diterima, masuk ke output.
- **rejected**: ide ditolak, tetap dicatat untuk audit trail.
- **modified**: ide diubah dari versi awal.

### 7. Dokumentasi Selama Sesi

- Semua ide dicatat **append-only** per sub-topik.
- Setiap sub-topik punya section sendiri.
- Keputusan/catatan khusus ditandai dengan `[KEPUTUSAN]` atau `[CATATAN]`.

---

## Format Sesi

### Tampilan per Iterasi

```
────────────────────────────────────────
Sub-Topik: [nama sub-topik]
Kategori: [diagnosis/ideation/analysis/perspective/solution/planning/innovation]
Teknik: [nama teknik]
────────────────────────────────────────

[Prompt dari agent sesuai teknik]

[Dari Joni]: [ide]
[Dari Jono]: [tanggapan]

User: [approve/reject/modify]

[Ide #1]: [judul]
  Detail: ...
  Dari: joni
  Status: approved

────────────────────────────────────────
Lanjut / Ganti sub-topik / Selesai?
```

---

## Transisi ke Step 3

Trigger ke Step 3:
- User bilang "selesai" atau "cukup"
- Semua sub-topik sudah selesai difasilitasi
- User memilih "Selesai" di prompt HALT

→ Lanjut ke `step-03-wrap-up.md`
