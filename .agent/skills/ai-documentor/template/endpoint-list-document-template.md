# Template Endpoint List Document

## Tujuan
Template ini dipakai untuk mendokumentasikan daftar endpoint per komponen sebelum file
`API-SPEC-*` dibuat satu per satu.

Template ini cocok dipakai saat:
- membuat inventaris endpoint awal dari codebase
- mengelompokkan endpoint per komponen atau domain
- menandai status lifecycle endpoint sebelum dibuat spec detail

## Penamaan Dokumen

Untuk indeks default workspace `rest-api-doc`, gunakan:

```text
daftar-endpoint.md
```

Jika user meminta daftar endpoint per domain terpisah, nama file MAY memakai pola:

```text
daftar-endpoint-<nama-komponen>.md
```

Contoh:
- `daftar-endpoint-user.md`
- `daftar-endpoint-auth.md`
- `daftar-endpoint-metatrader.md`

## Struktur Standar

```md
# Daftar Endpoint

## 1. Ringkasan
## 2. Daftar Endpoint per Komponen
## 3. Catatan Validitas
## 4. Asumsi, Risiko, dan Hal yang Perlu Dikonfirmasi
```

## 1. Ringkasan

Tuliskan ringkasan singkat:
- sumber bukti endpoint
- scope komponen
- apakah daftar ini lengkap atau masih parsial

## 2. Daftar Endpoint per Komponen

Gunakan satu tabel per komponen.

### Format Tabel

| Method | Endpoint | Description | Status |
|--------|----------|-------------|--------|
| GET | /api/v1/example | Penjelasan singkat endpoint | Published |

### Aturan Kolom
- `Method`: metode HTTP seperti `GET`, `POST`, `PUT`, `PATCH`, `DELETE`
- `Endpoint`: path endpoint lengkap
- `Description`: fungsi endpoint secara singkat dan teknis
- `Status`: gunakan salah satu nilai berikut:
  - `Draft`
  - `Published`
  - `Deprecated`
  - `Removed`

### Contoh Pengelompokan

#### Komponen Auth

| Method | Endpoint | Description | Status |
|--------|----------|-------------|--------|
| POST | /api/v1/auth/login | Login user dengan credential | Published |
| GET | /api/v1/auth/google/callback | Callback login Google | Published |

#### Komponen User

| Method | Endpoint | Description | Status |
|--------|----------|-------------|--------|
| GET | /api/v1/user/profile | Ambil profil user aktif | Published |
| PATCH | /api/v1/user/profile | Ubah profil user aktif | Draft |

## 3. Catatan Validitas

Tuliskan bila ada hal seperti:
- endpoint hanya terlihat di router tetapi handler belum ditemukan
- endpoint terlihat di OpenAPI tetapi implementasi runtime belum terverifikasi
- endpoint legacy masih ada di code tetapi kemungkinan tidak dipakai

## 4. Asumsi, Risiko, dan Hal yang Perlu Dikonfirmasi

Pisahkan bagian yang belum pasti, misalnya:
- grouping komponen masih inferensi
- status `Published` belum tervalidasi runtime
- endpoint callback berasal dari config, bukan route registration langsung
