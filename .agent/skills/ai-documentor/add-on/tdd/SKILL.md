---
name: tdd
description: Use when a greenfield project explicitly opts into test-driven development, or when implementing a planned behavior under an enabled TDD policy.
---

# TDD Add-On

Add-on ini menjalankan pengembangan greenfield dengan siklus **RED → GREEN → REFACTOR**. Add-on hanya aktif setelah user menyetujui penggunaannya dan keputusan tersebut dicatat di `.ai-doc/constitution.md`.

## Kapan Digunakan

Gunakan saat:

- bootstrap greenfield project dan user memilih TDD;
- use case atau behavior baru akan diimplementasikan;
- `.ai-doc/constitution.md` menyatakan `TDD: Enabled`;
- user meminta tracking TDD secara eksplisit.

Jangan aktifkan diam-diam. Untuk dokumentasi atau planning tanpa implementasi, TDD tidak diperlukan.

## Aturan Utama

1. **Tidak ada production code untuk behavior baru sebelum test yang sesuai ditulis.**
2. Jalankan test baru dan buktikan kegagalan yang diharapkan (**RED**).
3. Tulis implementasi paling kecil yang membuat test lulus (**GREEN**).
4. Setelah GREEN, rapikan implementasi tanpa menambah behavior (**REFACTOR**).
5. Jalankan kembali test terkait dan regression test setelah refactor.
6. Catat setiap target dan bukti transisinya di `.ai-doc/tdd-overview.md`.
7. Jika test tidak dapat dijalankan, gunakan `BLOCKED` atau `Perlu Dikonfirmasi`; jangan mengklaim RED/GREEN tanpa output nyata.

## Urutan Operasional

1. Baca `.ai-doc/3p.md` dan `.ai-doc/constitution.md`.
2. Konfirmasi policy TDD bila belum tercatat.
3. Setelah TDD enabled, buat `.ai-doc/tdd-overview.md` dari `template/tdd-overview-template.md`.
4. Turunkan scope/use case dan acceptance criteria menjadi target test yang dapat diamati.
5. Untuk setiap target: `PLANNED` → `RED` → `GREEN` → `REFACTORING` → `REFACTORED`.
6. Update overview dan `3p.md` pada setiap langkah bermakna.
7. Hentikan atau minta keputusan user bila target `BLOCKED` atau memerlukan `EXCEPTION`.

## Referensi Lokal

- `workflow.md` — alur activation, siklus, dan wrap-up.
- `writing-good-tests.md` — aturan kualitas test dan mock.
- `template/tdd-overview-template.md` — template pusat kontrol project.

Add-on ini berdiri sendiri dan tidak memerlukan skill atau file eksternal.
