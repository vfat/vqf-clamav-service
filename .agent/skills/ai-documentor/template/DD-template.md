# Template Data Dictionary (DD)

## Tujuan
Template ini mendefinisikan struktur standar untuk mendokumentasikan **Data Dictionary (DD)** bagi tabel database atau storage entity.

Data Dictionary dipakai untuk menjelaskan:
- metadata tabel
- struktur kolom
- index
- foreign key

Template ini cocok dipakai saat:
- mendokumentasikan schema database
- melengkapi ERD
- menjelaskan storage model dalam dokumentasi fitur
- melengkapi DCD yang bergantung pada tabel atau entitas persistence

## Penamaan File

Format penamaan file Data Dictionary:

```text
DD-<nomor urut>-<nama tabel>.md
```

Contoh:
- `DD-01-users.md`
- `DD-02-accounts.md`
- `DD-03-trades.md`

Panduan:
- `DD` adalah prefix tetap untuk Data Dictionary.
- `nomor urut` menggunakan dua digit atau lebih bila diperlukan.
- `nama tabel` mengikuti nama tabel utama yang didokumentasikan.
- Jika satu dokumen memuat beberapa tabel dalam satu domain, gunakan nama tabel/kelompok paling representatif.
- Gunakan huruf kecil dan pisahkan kata dengan tanda `-` bila nama tabel terdiri dari beberapa kata.

## Prinsip
- Data Dictionary MUST berbasis bukti dari codebase atau schema nyata:
  - `DbContext`
  - entity/model class
  - migration
  - SQL schema
  - repository query shape
- Jangan mengarang kolom, index, atau foreign key yang tidak terlihat dari bukti.
- Jika sebagian detail belum jelas, dokumentasi SHOULD menambahkan catatan asumsi atau hal yang perlu dikonfirmasi.
- Satu dokumen DD MAY memuat:
  - satu tabel
  - satu kelompok tabel
  - seluruh schema kecil

## Struktur Standar

```md
# Data Dictionary

## 1. Metadata Table Database

| Table Name | Prefix | Description |
|------------|--------|-------------|
| users | usr | Menyimpan data user aplikasi |

## 2. Struktur Kolom

### 2.1 users

| Column Name | Type | Keterangan |
|-------------|------|------------|
| id | string | Primary key user |
| email | string | Email login user |

## 3. Indexes

### 3.1 users

| Tipe Index | Value |
|------------|-------|
| Primary | id |
| Unique | email |
| Index | created_at |

## 4. Foreign Keys

### 4.1 users

| Source | Target | On Delete | On Update |
|--------|--------|-----------|-----------|
| profile_id | profiles.id | SET NULL | CASCADE |
```

## 1. Metadata Table Database

Bagian ini mendokumentasikan daftar tabel beserta metadata singkatnya.

### Format

| Table Name | Prefix | Description |
|------------|--------|-------------|
| `<nama_tabel>` | `<prefix>` | `<deskripsi_tabel>` |

### Panduan
- `Table Name`: nama tabel fisik atau nama entity storage yang dipetakan
- `Prefix`: prefix internal atau konvensi penamaan tabel jika ada
- `Description`: fungsi tabel dalam bahasa yang singkat dan teknis

### Catatan
- Jika project tidak memakai prefix tabel, isi `-`
- Jika tabel berasal dari entity class dan nama tabel fisik belum eksplisit, tulis nama entity/table yang paling dapat dipertanggungjawabkan dan beri catatan bila perlu

## 2. Struktur Kolom

Bagian ini mendokumentasikan kolom per tabel.

### Format

| Column Name | Type | Keterangan |
|-------------|------|------------|
| `<nama_kolom>` | `<tipe>` | `<fungsi / arti kolom>` |

### Panduan
- `Column Name`: nama kolom fisik atau properti entity yang dipetakan
- `Type`: tipe data yang terlihat dari schema, migration, atau entity
- `Keterangan`: arti bisnis atau teknis dari kolom

### Catatan
- Jika nullable penting untuk dipahami, tulis di kolom `Keterangan`
- Jika tipe database tidak terlihat pasti, gunakan tipe yang paling defensible dari bukti code dan tambahkan catatan di dokumen pemakai

## 3. Indexes

Bagian ini mendokumentasikan constraint atau index penting per tabel.

### Format

| Tipe Index | Value |
|------------|-------|
| Primary | `<kolom / kombinasi kolom>` |
| Unique | `<kolom / kombinasi kolom>` |
| Index | `<kolom / kombinasi kolom>` |

### Panduan
- `Tipe Index` hanya memakai:
  - `Primary`
  - `Unique`
  - `Index`
- `Value` dapat berupa:
  - satu kolom
  - kombinasi beberapa kolom

### Contoh

| Tipe Index | Value |
|------------|-------|
| Primary | id |
| Unique | email |
| Index | user_id, created_at |

## 4. Foreign Keys

Bagian ini mendokumentasikan relasi foreign key per tabel.

### Format

| Source | Target | On Delete | On Update |
|--------|--------|-----------|-----------|
| `<kolom_sumber>` | `<tabel_tujuan>.<kolom_tujuan>` | `<aksi>` | `<aksi>` |

### Panduan
- `Source`: kolom FK di tabel sumber
- `Target`: referensi target dalam format `table.column`
- `On Delete`: aksi saat record target dihapus
- `On Update`: aksi saat key target berubah

### Value Umum
- `CASCADE`
- `SET NULL`
- `RESTRICT`
- `NO ACTION`
- `-` bila tidak terlihat dari bukti

## Template per Tabel

Gunakan pola ini untuk setiap tabel:

```md
## 1. Metadata Table Database

| Table Name | Prefix | Description |
|------------|--------|-------------|
| users | usr | Menyimpan data user aplikasi |

## 2. Struktur Kolom

### 2.1 users

| Column Name | Type | Keterangan |
|-------------|------|------------|
| id | string | Primary key user |
| email | string | Email login user |
| created_at | datetime | Waktu pembuatan record |

## 3. Indexes

### 3.1 users

| Tipe Index | Value |
|------------|-------|
| Primary | id |
| Unique | email |
| Index | created_at |

## 4. Foreign Keys

### 4.1 users

| Source | Target | On Delete | On Update |
|--------|--------|-----------|-----------|
| - | - | - | - |
```

## Template Multi-Tabel

Gunakan pola ini jika satu dokumen memuat beberapa tabel:

```md
# Data Dictionary

## 1. Metadata Table Database

| Table Name | Prefix | Description |
|------------|--------|-------------|
| users | usr | Menyimpan data user aplikasi |
| accounts | acc | Menyimpan akun trading |
| trades | trd | Menyimpan transaksi trade |

## 2. Struktur Kolom

### 2.1 users
| Column Name | Type | Keterangan |
|-------------|------|------------|
| id | string | Primary key user |

### 2.2 accounts
| Column Name | Type | Keterangan |
|-------------|------|------------|
| id | string | Primary key account |
| user_id | string | Relasi ke user |

### 2.3 trades
| Column Name | Type | Keterangan |
|-------------|------|------------|
| id | string | Primary key trade |
| account_id | string | Relasi ke account |

## 3. Indexes

### 3.1 users
| Tipe Index | Value |
|------------|-------|
| Primary | id |

### 3.2 accounts
| Tipe Index | Value |
|------------|-------|
| Primary | id |
| Index | user_id |

### 3.3 trades
| Tipe Index | Value |
|------------|-------|
| Primary | id |
| Index | account_id, created_at |

## 4. Foreign Keys

### 4.1 users
| Source | Target | On Delete | On Update |
|--------|--------|-----------|-----------|
| - | - | - | - |

### 4.2 accounts
| Source | Target | On Delete | On Update |
|--------|--------|-----------|-----------|
| user_id | users.id | CASCADE | CASCADE |

### 4.3 trades
| Source | Target | On Delete | On Update |
|--------|--------|-----------|-----------|
| account_id | accounts.id | CASCADE | CASCADE |
```

## Panduan Pembuatan

### Step 1: Identifikasi Sumber Kebenaran

Cari dari:
- `DbContext`
- entity/model
- migration
- SQL DDL
- schema snapshot

### Step 2: Tentukan Scope

Pilih salah satu:
- satu tabel
- beberapa tabel dalam satu domain
- seluruh schema kecil

### Step 3: Isi Metadata Table Database

Daftar semua tabel yang termasuk scope.

### Step 4: Isi Struktur Kolom

Tulis hanya kolom yang benar-benar terlihat dari bukti code/schema.

### Step 5: Isi Indexes

Gunakan index yang terlihat dari:
- migration
- fluent config
- annotation
- DDL

### Step 6: Isi Foreign Keys

Gunakan relasi yang memang terlihat dari:
- FK column
- migration relation
- ORM mapping
- DDL constraint

### Step 7: Tambahkan Catatan di Dokumen Pemakai

Jika ada yang belum pasti, tambahkan:
- `Asumsi`
- `Perlu Dikonfirmasi`
- `Risiko Dokumentasi Tidak Lengkap`

## Kapan Template Ini Dipakai

Gunakan template ini saat:
- mendokumentasikan tabel database
- melengkapi ERD
- melengkapi DCD yang bergantung pada storage
- menjelaskan persistence model untuk onboarding developer baru

## Kapan Tidak Dipakai

Jangan gunakan template ini untuk:
- alur runtime
- use case actor
- dependency antar module
- lifecycle service/process

Untuk itu, gunakan template lain yang sesuai.
