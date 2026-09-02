# Brownfield Flow

Workflow ini dipakai saat mendokumentasikan codebase yang sudah ada.

## 1. Reconstruct the Current Truth

Audit source nyata:
- entry points
- startup/lifecycle
- runtime communication
- config/dependency files
- services/repositories
- UI/ViewModel/API surfaces
- schema/contracts bila relevan

Untuk review dokumen, jangan percaya dokumen lama. Re-check code.

## 2. Select the Smallest Correct Artifact

Pilih artefak yang paling sempit yang menjawab permintaan user:
- overview doc
- feature doc
- grouped use case
- DCD
- database design doc
- REST API doc
- sync/review pass

## 3. Write Conservatively

Aturan penulisan:
- prefer wording yang konkret
- jangan membesar-besarkan shell/stub menjadi flow lengkap
- kalau evidence lemah, tandai sebagai:
  - `Draft`
  - `Partial`
  - `Placeholder`
  - `Asumsi`
  - `Perlu Dikonfirmasi`

## 4. Sync Neighbor Documents

Contoh sync:
- perubahan DCD -> review `Dokumentasi-Komponen-Usecase.md`
- perubahan status fitur -> review `Dokumentasi-Fitur.md`
- perubahan interpretasi runtime -> review `Dokumentasi-Codebase.md`
- perubahan workflow/aturan -> review `.ai-doc/constitution.md`

## 5. Database Design Branch

Jika user meminta `desain-database-document`:
1. buat `.ai-doc/desain-database-document/` bila belum ada
2. baca schema dulu
3. buat ERD dulu
4. lanjut DD per tabel
5. bedakan relasi logical vs physical bila FK fisik tidak terlihat

## 6. REST API Branch

Jika user meminta `rest-api-doc` atau REST API spec:
1. buat `.ai-doc/rest-api-doc/` bila belum ada
2. baca router/controller/middleware/DTO/OpenAPI lebih dulu
3. buat `daftar-endpoint.md` terlebih dahulu memakai `template/endpoint-list-document-template.md`
4. lanjut file `API-SPEC-*` satu per satu
5. pakai swimlane bila flow lintas layer memang terlihat dari code

## 7. Record State

Setelah langkah bermakna selesai:
- update `.ai-doc/3p.md`
- isi `Progress`
- isi `Plan`
- isi `Pending`
