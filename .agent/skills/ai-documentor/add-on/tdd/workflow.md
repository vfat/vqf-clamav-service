# TDD Add-On — Workflow

> Workflow mandiri untuk TDD pada greenfield project.

## Tujuan

Membuat keputusan TDD eksplisit, memulai dari test yang gagal, dan menjaga `.ai-doc/tdd-overview.md` sebagai pusat status setiap behavior.

## Prasyarat

Sebelum memulai:

- project dikategorikan greenfield;
- scope, use case, dan acceptance criteria awal sudah cukup jelas;
- `.ai-doc/3p.md` dan `.ai-doc/constitution.md` sudah dibaca;
- test runner, lokasi test, atau command test diketahui. Jika belum, tandai `Perlu Dikonfirmasi`.

## Flow

```text
Step 1: Activation
  ├── Tanya user apakah TDD digunakan
  ├── Catat Enabled atau Disabled di constitution
  ├── Jika Enabled, buat .ai-doc/tdd-overview.md
  └── HALT → konfirmasi policy dan target awal

Step 2: RED / GREEN / REFACTOR loop
  ├── PLANNED: daftarkan behavior dan acceptance criteria
  ├── RED: tulis test sebelum production code
  ├── Jalankan test dan verifikasi gagal karena behavior belum ada
  ├── Update overview + 3p.md
  ├── GREEN: tulis implementasi minimal
  ├── Jalankan test dan verifikasi lulus
  ├── REFACTORING: rapikan tanpa menambah behavior
  ├── Jalankan test terkait dan regression test
  ├── Tandai REFACTORED atau BLOCKED
  └── HALT → lanjut target berikutnya atau wrap-up

Step 3: Wrap-Up
  ├── Periksa target yang belum selesai
  ├── Ringkas bukti dan blocker
  ├── Update overview + 3p.md
  └── Kembalikan kontrol ke AI Documentor
```

## Step 1 — Activation

Pertanyaan minimum:

> Apakah project ini ingin dikembangkan dengan TDD: menulis test untuk behavior yang belum ada, memastikan test gagal terlebih dahulu, lalu implementasi minimal dan refactor?

Jika user menjawab ya:

- tulis `TDD: Enabled` dan scope-nya di `.ai-doc/constitution.md`;
- buat `.ai-doc/tdd-overview.md` dari template;
- jangan membuat status `RED` sebelum test benar-benar dijalankan.

Jika user menjawab tidak:

- tulis `TDD: Disabled` di constitution;
- jangan jalankan TDD diam-diam.

## Step 2 — Cycle Rules

Setiap target memiliki ID stabil, misalnya `TDD-001`.

### RED

- Test menguji satu behavior yang dapat diamati.
- Test ditulis sebelum production implementation.
- Jalankan command test yang sebenarnya.
- Catat command, test target, exit status, dan ringkasan output.
- Kegagalan harus disebabkan behavior yang belum tersedia, bukan typo atau setup rusak.

### GREEN

- Tulis implementasi minimal.
- Jangan menambah behavior yang belum dituntut test.
- Jalankan test target dan test terkait.
- Catat output pass yang nyata.

### REFACTOR

- Ubah struktur, nama, atau duplikasi tanpa mengubah behavior.
- Tandai `REFACTORING` saat pekerjaan berlangsung.
- Jalankan ulang test setelah refactor.
- Gunakan `REFACTORED` hanya jika test tetap lulus.

## Blocked dan Exception

Gunakan `BLOCKED` jika runner, dependency, environment, atau requirement menghalangi verifikasi. Gunakan `EXCEPTION` hanya setelah user menyetujui pengecualian dan alasannya dicatat di overview serta constitution bila berlaku project-wide.

## Wrap-Up Checklist

- Semua target memiliki status terbaru.
- Setiap `RED` memiliki failure evidence.
- Setiap `GREEN` memiliki passing evidence.
- Setiap `REFACTORED` memiliki regression evidence.
- Blocker dan exception tidak disamarkan sebagai selesai.
- `.ai-doc/3p.md` sudah mencerminkan langkah terakhir.
