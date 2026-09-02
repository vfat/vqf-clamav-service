# Template Spec Component Document

## Tujuan
Template ini dipakai untuk planning komponen baru sebelum implementasi detail dimulai.

Dokumen hasil template ini disimpan di:
- `.ai-doc/plan/component/`

Format penamaan file:

```text
SCD-<nama komponen baru>.md
```

Contoh:
- `SCD-User-Profile.md`
- `SCD-Notification-Center.md`
- `SCD-MetaTrader-Sync.md`

## Struktur Standar

```md
# SCD-<nama komponen baru>

## 1. Context
## 2. Scope
## 3. Prerequisite
## 4. Daftar Usecase
## 5. Catatan Diskusi
## 6. Asumsi, Risiko, dan Hal yang Perlu Dikonfirmasi
```

## 1. Context

Tuliskan:
- masalah yang ingin diselesaikan
- posisi komponen dalam sistem
- hubungan komponen dengan komponen lain yang sudah ada

## 2. Scope

Tuliskan:
- apa yang termasuk dalam komponen
- apa yang di luar scope
- boundary awal komponen

## 3. Prerequisite

Tuliskan dependency perencanaan atau implementasi awal, misalnya:
- endpoint/backend yang harus tersedia
- entity/schema yang perlu ada
- dependency ke runtime atau komponen lain
- prerequisite UI/UX atau flow bisnis

## 4. Daftar Usecase

Gunakan tabel berikut.

| Kode Usecase | Nama Usecase | Deskripsi Singkat |
|--------------|--------------|-------------------|
| UC-01 | Nama usecase | Ringkasan singkat usecase |

Aturan:
- hanya isi kode usecase, nama usecase, dan deskripsi singkat
- detail flow use case belum masuk di tahap SCD
- gunakan kode yang stabil dan mudah diturunkan ke DCD di tahap berikutnya

## 5. Catatan Diskusi

Bagian ini dipakai untuk menangkap hasil tanya jawab dengan user:
- keputusan naming
- alternatif scope
- saran agent
- keputusan prioritas

## 6. Asumsi, Risiko, dan Hal yang Perlu Dikonfirmasi

Pisahkan:
- asumsi agent
- risiko desain bila requirement masih kabur
- pertanyaan lanjutan yang perlu dikonfirmasi sebelum implementasi atau DCD
