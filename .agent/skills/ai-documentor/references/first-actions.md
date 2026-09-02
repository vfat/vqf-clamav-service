# First Actions

Gunakan checklist ini sebelum menulis atau mengubah dokumentasi.

## 1. Bootstrap Workspace Minimum

Jika `.ai-doc/` belum ada:
- buat `.ai-doc/`
- buat `.ai-doc/3p.md`
- buat `.ai-doc/constitution.md` dari `references/constitution-general.md`

Untuk workflow ini, skill default yang dipakai adalah `ai-documentor`.
Skill tambahan boleh dipakai kemudian bila user memang membutuhkan.

Skill ini tidak perlu menyalin template ke `.ai-doc/template/`.
Gunakan template dari folder `template/` skill ini, kecuali project memang sudah punya template
lokal yang sengaja dipelihara sendiri.

## 2. Baca Control Plane

Wajib baca:
- `.ai-doc/3p.md`
- `.ai-doc/constitution.md`

Tujuannya:
- memahami progress terakhir
- melihat aturan lokal project
- menghindari membuat artefak baru yang tidak diminta user

## 3. Tentukan Artefak Target

Pilih artefak terkecil yang tepat:
- project overview greenfield
- dokumentasi codebase
- planning komponen baru
- dokumentasi fitur
- grouped use case
- DCD
- C4/runtime diagram
- desain database document
- REST API doc
- review/sinkronisasi

Jika request bersifat greenfield dan akan berlanjut ke implementasi, tanyakan secara eksplisit
apakah user ingin mengaktifkan TDD add-on. Jangan mengaktifkan TDD hanya karena project
tergolong greenfield.

## 4. Audit Bukti Kode Dulu

Sebelum menulis:
- baca entry point
- baca flow startup/runtime
- baca config dan dependency
- baca service/repository/viewmodel/controller/handler yang relevan
- baca source of truth untuk artifact yang diminta

Jangan infer hanya dari nama file.

## 5. Apply Only-On-Request Rule

Artefak turunan berikut hanya dibuat bila user meminta eksplisit:
- `Dokumentasi-Fitur.md`
- `Dokumentasi-Komponen-Usecase.md`
- `C4-Component-Diagrams.md`
- `desain-component-document/`
- `desain-database-document/`
- `rest-api-doc/`
