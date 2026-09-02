# Component Planning Flow

Workflow ini dipakai saat user ingin merencanakan komponen baru yang belum ada atau belum
jelas bentuknya.

## Tujuan

Hasil akhirnya adalah **Spec Component Document** di:
- `.ai-doc/plan/component/`

Dengan format nama:
- `SCD-<nama komponen baru>.md`

## Trigger

Gunakan flow ini saat user meminta hal seperti:
- planning komponen baru
- spec komponen baru
- rancangan awal komponen
- pemetaan use case komponen yang belum dibangun

## Urutan Kerja

1. Pastikan request memang bersifat planning, bukan dokumentasi brownfield komponen yang sudah ada.
2. Buat folder `.ai-doc/plan/component/` bila belum ada.
3. Mulai tanya jawab singkat dengan user untuk menggali:
   - masalah yang mau diselesaikan
   - posisi komponen dalam sistem
   - scope utama
   - dependency atau prerequisite
   - daftar use case awal
4. Jika user bingung, agent SHOULD memberi saran:
   - usulan nama komponen
   - usulan scope yang lebih sempit
   - usulan daftar use case awal
   - usulan komponen existing yang terkait
5. Tulis hasilnya ke file `SCD-<nama komponen baru>.md` memakai `template/SCD-template.md`.
6. Catat asumsi dan hal yang masih perlu dikonfirmasi.
7. Update `.ai-doc/3p.md`.

## Aturan Tanya Jawab

Pertanyaan agent SHOULD berfokus pada:
- apa tujuan komponen
- siapa aktor atau pengguna utamanya
- apa input dan output utamanya
- komponen existing apa yang berinteraksi
- use case minimum yang benar-benar perlu

Hindari:
- langsung membuat DCD penuh
- langsung membuat detail normal/alternative flow
- langsung mengarang dependency teknis yang belum dibahas

## Struktur SCD

SCD hanya memuat:
- `Context`
- `Scope`
- `Prerequisite`
- `Daftar Usecase`
- `Catatan Diskusi`
- `Asumsi, Risiko, dan Hal yang Perlu Dikonfirmasi`

Di tahap ini, `Daftar Usecase` hanya berisi:
- kode usecase
- nama usecase
- deskripsi singkat

## Handoff ke Tahap Berikutnya

Jika nanti user ingin melanjutkan:
- SCD bisa diturunkan menjadi DCD
- daftar use case bisa dipakai untuk grouped use case doc
- prerequisite bisa dipakai untuk planning API, DB, atau dependency komponen
