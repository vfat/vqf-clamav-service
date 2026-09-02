# Template Dokumentasi Komponen dan Use Case

Gunakan template ini untuk menyusun `.ai-doc/Dokumentasi-Komponen-Usecase.md`.
Struktur section HARUS mengikuti urutan yang sama seperti template ini agar konsisten dengan dokumen master.

## Aturan Umum

- Judul dokumen HARUS `# Dokumentasi Komponen dan Use Case`.
- Section utama HARUS tetap: `1. Ringkasan`, `2. Diagram Use Case Tergrup`, `3. Daftar Use Case per Komponen`, `4. Catatan Pengelompokan & Sinkronisasi DCD`.
- Nama komponen di section 3 HARUS sama dengan nama package di diagram dan, jika ada, sama dengan nama artefak `DCD-*` terkait.
- Nama use case di section 3 HARUS sama dengan label use case di diagram. Jangan membuat sinonim antar section.
- Ringkasan global hanya ditulis di section 1. Jangan mengulang ringkasan global di setiap sub-section komponen.
- Deskripsi komponen harus singkat dan fokus pada scope pembeda komponen itu, bukan mengulang semua butir use case. Target 2-4 kalimat pendek.
- Daftar `Use case yang terverifikasi` adalah daftar utama. Paragraf `Deskripsi` tidak boleh menyalin ulang daftar tersebut dengan kalimat yang hampir sama.
- Jika implementasi belum lengkap, tandai dengan status yang jujur seperti `Parsial`, `Draft`, `External`, atau `Belum terlihat jelas`.
- Bukti kode HARUS berupa file/class/service/viewmodel yang benar-benar dibaca, bukan inferensi dari nama folder.

## Template Dokumen

~~~md
# Dokumentasi Komponen dan Use Case

## 1. Ringkasan
Dokumen ini mengelompokkan komponen yang terlihat di dalam <Nama Sistem / Produk> berdasarkan daftar yang diberikan user, lalu menurunkan use case yang tampak dari codebase aktif. Fokus dokumen ini bukan desain ideal, tetapi use case yang punya jejak implementasi di boundary UI, service host, repository/gateway IPC, atau helper process.

Dokumen ini memakai tiga status:
- `Ada`: boundary dan alur use case terlihat cukup jelas di codebase aktif.
- `Parsial`: ada hook UI, service, atau halaman, tetapi implementasi belum lengkap atau masih stub.
- `Belum terlihat jelas`: belum ditemukan boundary aktif yang cukup untuk menyatakan fitur itu benar-benar ada.

Catatan validitas:
- Dokumen ini diposisikan sebagai peta use case praktis, bukan kontrak produk final.
- Jika sebuah komponen hanya punya hook UI, command kosong, atau page yang belum terbukti terhubung end-to-end, deskripsinya diturunkan ke `Parsial` atau `Belum terlihat jelas`.
- Use case yang sudah diperdalam di DCD lebih diutamakan daripada inferensi lama dari nama file.

## 2. Diagram Use Case Tergrup
~~~plantuml
@startuml
left to right direction

actor User
actor "Operator / Installer" as Operator
actor "Backend API / Socket.IO" as Backend
actor "System / Host" as SystemActor

package "Nama Komponen 1" {
  usecase "Nama Use Case 1" as UC_Component1_UC1
  usecase "Nama Use Case 2" as UC_Component1_UC2
}

package "Nama Komponen 2" {
  usecase "Nama Use Case 1" as UC_Component2_UC1
}

User --> UC_Component1_UC1
Operator --> UC_Component1_UC2
Backend --> UC_Component2_UC1
SystemActor --> UC_Component1_UC1

UC_Component1_UC1 ..> UC_Component2_UC1 : <<include>>
UC_Component1_UC2 ..> UC_Component1_UC1 : <<extend>>
@enduml
~~~

## 3. Daftar Use Case per Komponen
### 3.1 Nama Komponen
**Status**: ✅ **Active**
**Reference**: `DCD-XX-Nama-Komponen.md` (opsional jika memang ada)

Deskripsi:
Paragraf singkat 2-4 kalimat yang menjelaskan scope komponen, pembeda utamanya, dan tingkat kematangan implementasi. Jangan menyalin ulang daftar use case satu per satu dalam bentuk paragraf panjang.

Use case yang terverifikasi:
- ✅ `Nama Use Case 1` (UC-1, **Active**)
- ⏳ `Nama Use Case 2` (UC-2, **Parsial/Draft**)

Bukti kode:
- `path/file-1`
- `path/file-2`
- `path/file-3`

Implementasi:
- Catatan implementasi yang benar-benar berbeda dari daftar use case, misalnya event, polling, IPC, persistence, scheduler, atau batas validasi.
- Hindari bullet yang hanya mengulang nama use case dengan kata lain.

Catatan:
- Opsional. Isi hanya jika ada batas validitas, placeholder, atau area yang sengaja belum diklaim penuh.

### 3.2 Nama Komponen Lain
**Status**: ⏳ **Parsial**
**Reference**: `DCD-YY-Komponen-Lain.md`

Deskripsi:
Paragraf singkat yang fokus pada scope, bukan pengulangan daftar use case.

Use case yang terverifikasi:
- ✅ `Nama Use Case 1` (UC-1, **Active**)
- ⏳ `Nama Use Case 2` (UC-2, **Draft**)

Bukti kode:
- `path/file-1`

Implementasi:
- Catatan implementasi ringkas.

## 4. Catatan Pengelompokan & Sinkronisasi DCD

### Synchronization Status
- ✅ **Synchronized (X/Y)**: daftar DCD yang sudah sinkron
- ⏳ **Pending**: daftar komponen yang belum sinkron, jika ada

### Key Separation Notes
- Komponen A vs Komponen B: jelaskan alasan pemisahan bila rawan tumpang tindih.
- Komponen background vs komponen UI: jelaskan jika ada komponen yang bukan halaman UI tunggal.
~~~

## Checklist Review Sebelum Final

- Apakah empat section utama ada dan urutannya benar?
- Apakah nama komponen konsisten antara section 3, diagram, dan DCD?
- Apakah daftar use case tidak diulang lagi di paragraf `Deskripsi`?
- Apakah bullet `Implementasi` berisi detail runtime/teknis yang berbeda dari daftar use case?
- Apakah semua klaim status punya bukti kode yang sudah dibaca?