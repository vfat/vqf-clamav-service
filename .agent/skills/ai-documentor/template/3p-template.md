# Template 3P (Progress, Plan, Pending)

## Tujuan
Template ini digunakan sebagai control plane progress tracking untuk seluruh aktivitas dokumentasi dan planning di dalam `.ai-doc/`.

Dokumen hasil template ini disimpan di:
- `.ai-doc/3p.md`

## Aturan Pembaruan
- `3p.md` wajib dibaca sebelum memulai atau melanjutkan pekerjaan dokumentasi/planning.
- `3p.md` wajib diperbarui setelah setiap langkah atau transisi bermakna.
- Format wajib memiliki 3 section utama: `Progress`, `Plan`, dan `Pending`.

---

## Struktur Dokumen

```md
# Progress, Plan, Pending (3P)

> Last Updated: <YYYY-MM-DD HH:mm>
> Active Flow: <Greenfield / Brownfield / Component Planning / Add-On / TDD>

---

## 1. Progress
Catatan langkah-langkah yang sudah selesai dikerjakan dan diverifikasi.

- [x] **<Tanggal/Aktivitas>**: <Deskripsi pencapaian/artefak yang dihasilkan>
  - *Artifact:* `<path ke artefak yang dibuat/diubah>`
  - *Status:* `<Done / Verified>`

---

## 2. Plan
Rencana langkah konkret yang akan dikerjakan berikutnya.

- [ ] **Step 1**: <Deskripsi langkah terdekat>
- [ ] **Step 2**: <Deskripsi langkah selanjutnya>

---

## 3. Pending
Daftar pertanyaan terbuka, asumsi yang perlu konfirmasi, atau blocker.

- [ ] **<Item>**: <Deskripsi hal yang menunggu konfirmasi user atau dependensi eksternal>
```
