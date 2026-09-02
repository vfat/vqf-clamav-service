# Greenfield Flow

Workflow ini dipakai saat user ingin merencanakan project atau sistem baru yang belum ada
codebase-nya atau belum punya bentuk solusi yang stabil.

## Tujuan

Artefak pertama untuk greenfield planning adalah:
- `.ai-doc/project-overview.md`

Dokumen ini memakai:
- `template/project-overview-template.md`

## Trigger

Gunakan flow ini saat user meminta hal seperti:
- planning project baru
- product/system overview untuk project baru
- project overview
- discovery awal sebelum desain komponen atau implementasi

## Urutan Kerja

1. Pastikan request memang bersifat greenfield planning, bukan dokumentasi brownfield.
2. Buat `.ai-doc/` bila belum ada.
3. Buat `project-overview.md` sebagai artefak planning pertama.
4. Bangun isi dokumen melalui tanya jawab agent dengan user.
5. Jika user belum jelas, agent SHOULD memberi saran tentang:
   - perumusan problem statement
   - target users/stakeholders
   - scope awal
   - arah solusi high level
   - constraint utama
   - prerequisite awal
6. Catat hasil klarifikasi, asumsi, dan risiko secara eksplisit.
7. Tanyakan keputusan TDD secara eksplisit bila greenfield akan berlanjut ke implementasi:
   - jika user memilih `Enabled`, catat policy dan scope di `.ai-doc/constitution.md`
   - inisialisasi `.ai-doc/tdd-overview.md` dari `add-on/tdd/template/tdd-overview-template.md`
   - baca dan jalankan `add-on/tdd/workflow.md`
   - jika user memilih `Disabled`, catat keputusan tersebut dan jangan menjalankan TDD diam-diam
8. Setelah TDD enabled dan use case/acceptance criteria cukup jelas, turunkan behavior menjadi target
   test. Tulis test lebih dahulu, jalankan sampai failure evidence valid (`RED`), lalu baru implementasi.
9. Update `.ai-doc/3p.md`.

## TDD Decision Gate

TDD bukan default greenfield dan tidak boleh diaktifkan berdasarkan asumsi agent. Pertanyaan minimum:

> Apakah project ini ingin dikembangkan dengan TDD: menulis test untuk behavior yang belum ada,
> memastikan test gagal terlebih dahulu, lalu membuat implementasi minimal dan melakukan refactor?

Jika test runner, lokasi test, atau command belum diketahui, tandai `Perlu Dikonfirmasi` atau
`BLOCKED`; jangan membuat klaim RED/GREEN tanpa hasil command nyata. TDD hanya mengontrol implementasi
behavior yang dipilih dan tidak membuat artefak turunan lain tanpa permintaan eksplisit user.

## Aturan Only-On-Request

Di fase greenfield overview, artefak turunan berikut MUST NOT dibuat kecuali user meminta
secara eksplisit:
- `rest-api-doc/`
- `desain-database-document/`
- `desain-component-document/`
- `C4-Component-Diagrams.md`
- C4 Container Diagram
- C4 Context Diagram

Urutan default greenfield tetap:
1. mulai dari `project-overview.md`
2. turun ke artefak turunan hanya jika diminta user atau sudah diputuskan eksplisit bersama user

## Aturan Tanya Jawab

Pertanyaan agent SHOULD menggali:
- masalah apa yang ingin diselesaikan
- siapa pengguna dan stakeholder utama
- apa target outcome yang diinginkan
- apa yang masuk dan di luar scope
- constraint apa yang sudah diketahui
- dependency awal apa yang harus tersedia

Hindari:
- langsung membuat DCD
- langsung membuat ERD atau API spec tanpa overview
- langsung membuat C4 diagram tanpa permintaan eksplisit user
- langsung mengunci arsitektur detail sebelum context dan scope jelas

## Struktur Project Overview

`project-overview.md` SHOULD memuat:
- `Problem Statement`
- `Target Users & Stakeholders`
- `Assumptions`
- `Goals & Objectives`
- `Scope`
- `High-Level System Direction`
- `Key Constraints`
- `Prerequisite`

Bagian tambahan MAY memuat:
- `Catatan Diskusi`
- `Risiko, Asumsi, dan Hal yang Perlu Dikonfirmasi`

## Handoff ke Tahap Berikutnya

Jika planning berlanjut:
- `project-overview.md` dapat diturunkan ke SCD per komponen
- dapat menjadi dasar untuk API planning, DB planning, atau runtime planning
- dapat menjadi dasar untuk backlog awal implementasi
