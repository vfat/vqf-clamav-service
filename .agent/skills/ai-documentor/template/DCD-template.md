# Template Desain Component Document (DCD)

## Tujuan
Template ini mendefinisikan struktur untuk membuat Desain Component Document (DCD) yang comprehensive.
Setiap komponen produk akan memiliki satu DCD yang mencakup object identification, use case listing, dan detail setiap use case.

## Penamaan File
Format: `DCD-<nomor urut>-<nama komponen>.md`

**Contoh**:
- `DCD-01-Login.md`
- `DCD-02-Sign-Up.md`
- `DCD-03-Dashboard.md`
- `DCD-04-MetaTrader.md`
- `DCD-05-Chart.md`
- dst.

---

## Struktur Dokumen

### 1. Judul Dokumen

```markdown
# DCD-<nomor urut>-<nama komponen>
```

**Contoh**:
```markdown
# DCD-01-Login
```

Catatan: Gunakan format yang konsisten. Nomor urut sesuai dengan urutan komponen produk yang sudah didaftar.

---

### 2. Section Object Identification

Gunakan format section dengan list dari `.ai-doc/template/object-identification-section-template.md`.

**Struktur**:
```markdown
## Object Identification

### Boundary
- [List item boundary untuk komponen ini]
- ...

### Control
- [List item control untuk komponen ini]
- ...

### Entity
- [List item entity untuk komponen ini]
- ...
```

**Penjelasan**:
- **Boundary**: Antarmuka/UI yang user lihat (screen, form, modal, button, API endpoint)
- **Control**: Proses/logika di balik layar (service, handler, validator, manager)
- **Entity**: Data yang diakses (database table, DTO, API response, configuration)

**Referensi**: Lihat `.ai-doc/template/object-identification-section-template.md` untuk detail lengkap.

---

### 3. Tabel List Use Case

Tabel berisi daftar semua use case untuk komponen ini dengan status dan referensi.

**Struktur**:
```markdown
## Use Case List

| No | Use Case Name | Actor | Status | Detail |
|----|---------------|-------|--------|--------|
| 1 | [UC Name 1] | [Actor] | [Active/Draft/TODO] | [Link atau referensi] |
| 2 | [UC Name 2] | [Actor] | [Active/Draft/TODO] | [Link atau referensi] |
| ... | ... | ... | ... | ... |
```

**Kolom**:
- **No**: Nomor urut use case dalam komponen ini (1, 2, 3, dst)
- **Use Case Name**: Nama use case yang jelas
- **Actor**: Siapa yang melakukan use case (User, System, Operator)
- **Status**: Active (sudah implemented), Draft (sedang dikerjakan), TODO (belum dikerjakan)
- **Detail**: Link ke section detail atau penjelasan singkat

**Contoh untuk Login**:
```markdown
| No | Use Case Name | Actor | Status | Detail |
|----|---------------|-------|--------|--------|
| 1 | Login Via Form | User | Active | See section 4.1 |
| 2 | Login Via SSO Google | User | Draft | See section 4.2 |
| 3 | Forgot Password | User | TODO | See section 4.3 |
```

---

### 4. Use Case Detail

Untuk setiap use case, buatkan dokumentasi detail menggunakan template dari `.ai-doc/template/usecase-detail-template.md`.

**Struktur per Use Case**:
```markdown
## 4. Use Case Details

### 4.1 [Use Case Name 1]

| Kategori | Keterangan |
|----------|-----------|
| Use Case Name | [nama] |
| Actor | [actor] |
| Deskripsi | [deskripsi] |
| Pre-Condition | [pre-condition] |
| Post-Condition | [post-condition] |
| Normal Flow | [flow langkah demi langkah] |
| Alternative Flow | [flow alternatif] |
| Exception Flow | [flow exception] |
| Related Use Case | [UC yang terkait] |
| Business Rule | [aturan bisnis] |

### 4.2 [Use Case Name 2]

| Kategori | Keterangan |
|----------|-----------|
| Use Case Name | [nama] |
| ... | ... |
```

**Penjelasan**:
- Setiap subsection adalah satu use case dengan format table 2 kolom.
- Ikuti template dari `.ai-doc/template/usecase-detail-template.md`.
- Nomor subsection harus match dengan nomor di tabel Use Case List.

**Referensi**: Lihat `.ai-doc/template/usecase-detail-template.md` untuk detail lengkap tentang setiap field.

---

## Template Lengkap (Skeleton)

```markdown
# DCD-<nomor urut>-<nama komponen>

## Overview
[Penjelasan singkat tentang komponen ini, apa fungsinya, dan siapa target usernya]

## Object Identification

### Boundary
- [Boundary items]

### Control
- [Control items]

### Entity
- [Entity items]

## Use Case List

| No | Use Case Name | Actor | Status | Detail |
|----|---------------|-------|--------|--------|
| 1 | [UC Name] | [Actor] | [Status] | See section 4.1 |

## Use Case Details

### 4.1 [Use Case Name]

| Kategori | Keterangan |
|----------|-----------|
| Use Case Name | [nama] |
| Actor | [actor] |
| Deskripsi | [deskripsi] |
| Pre-Condition | [pre-condition] |
| Post-Condition | [post-condition] |
| Normal Flow | [flow] |
| Alternative Flow | [flow] |
| Exception Flow | [flow] |
| Related Use Case | [UC terkait] |
| Business Rule | [aturan bisnis] |
```

---

## Daftar Komponen Produk yang Akan Memiliki DCD

Berdasarkan dokumentasi komponen, berikut adalah 15 komponen produk yang akan memiliki DCD masing-masing:

| No | Komponen | File | Status |
|----|----------|------|--------|
| 1 | Apps Installer | DCD-01-Apps-Installer.md | TODO |
| 2 | Login | DCD-02-Login.md | TODO |
| 3 | Sign Up | DCD-03-Sign-Up.md | TODO |
| 4 | Dashboard | DCD-04-Dashboard.md | TODO |
| 5 | MetaTrader | DCD-05-MetaTrader.md | TODO |
| 6 | Chart | DCD-06-Chart.md | TODO |
| 7 | Trade | DCD-07-Trade.md | TODO |
| 8 | Account History | DCD-08-Account-History.md | TODO |
| 9 | Create Account | DCD-09-Create-Account.md | TODO |
| 10 | Broker Installer | DCD-10-Broker-Installer.md | TODO |
| 11 | Notification | DCD-11-Notification.md | TODO |
| 12 | Feedback | DCD-12-Feedback.md | TODO |
| 13 | Setting | DCD-13-Setting.md | TODO |
| 14 | Data Collector | DCD-14-Data-Collector.md | TODO |
| 15 | Help | DCD-15-Help.md | TODO |

---

## Contoh DCD Minimal (DCD-02-Login)

### Penjelasan
Berikut adalah contoh minimal struktur DCD untuk komponen Login.

```markdown
# DCD-02-Login

## Overview
Komponen Login memfasilitasi user untuk autentikasi ke dalam sistem menggunakan form email/password atau SSO Google.
Fitur ini adalah gateway utama sebelum user bisa mengakses fitur-fitur lain di aplikasi.

## Object Identification

### Boundary
- LoginWindow.xaml (WPF window)
- AuthenticationViewModel
- POST /api/auth/login (API endpoint)
- Google OAuth Callback URI
- Error/Success Modal
- Remember Me Checkbox

### Control
- AuthenticationService
- ValidateEmailFormat()
- ValidatePassword()
- HandleAuthenticationResponse()
- GenerateSessionToken()
- EncryptCredentials()
- StoreSessionToken()

### Entity
- User (database entity)
- UserDTO (API DTO)
- AppSettings (encrypted credentials storage)
- Session (in-memory token)
- Backend Auth API Response

## Use Case List

| No | Use Case Name | Actor | Status | Detail |
|----|---------------|-------|--------|--------|
| 1 | Login Via Form | User | Active | See section 4.1 |
| 2 | Login Via SSO Google | User | Active | See section 4.2 |
| 3 | Forgot Password | User | Draft | See section 4.3 |
| 4 | Remember Me | User | Draft | See section 4.4 |

## Use Case Details

### 4.1 Login Via Form

| Kategori | Keterangan |
|----------|-----------|
| Use Case Name | Login Via Form |
| Actor | User |
| Deskripsi | User login menggunakan form email dan password di halaman login. |
| Pre-Condition | User sudah terdaftar dan berada di halaman login. |
| Post-Condition | User berhasil masuk ke dalam sistem dan mendapatkan akses ke dashboard. |
| Normal Flow | [step-by-step seperti di template] |
| Alternative Flow | [flow alternatif] |
| Exception Flow | [exception handling] |
| Related Use Case | • Login Via SSO Google • Forgot Password • Sign Up Via Form |
| Business Rule | [aturan bisnis sesuai template] |

### 4.2 Login Via SSO Google

| Kategori | Keterangan |
|----------|-----------|
| Use Case Name | Login Via SSO Google |
| Actor | User |
| Deskripsi | User login menggunakan akun Google melalui OAuth SSO. |
| Pre-Condition | User sudah memiliki akun Google. User berada di halaman login. |
| Post-Condition | User berhasil login dan redirect ke dashboard. |
| Normal Flow | [step-by-step] |
| Alternative Flow | [flow alternatif] |
| Exception Flow | [exception handling] |
| Related Use Case | • Login Via Form • Link Google Account (if account not yet linked) |
| Business Rule | [aturan bisnis] |

### 4.3 Forgot Password

| Kategori | Keterangan |
|----------|-----------|
| Use Case Name | Forgot Password |
| Actor | User |
| Deskripsi | User bisa mereset password jika lupa melalui email verification. |
| Pre-Condition | User sudah terdaftar tapi lupa password. User berada di halaman login. |
| Post-Condition | Password user berhasil direset. User bisa login dengan password baru. |
| Normal Flow | [step-by-step] |
| Alternative Flow | [flow alternatif] |
| Exception Flow | [exception handling] |
| Related Use Case | • Login Via Form |
| Business Rule | [aturan bisnis] |

### 4.4 Remember Me

| Kategori | Keterangan |
|----------|-----------|
| Use Case Name | Remember Me |
| Actor | User |
| Deskripsi | User bisa centang "Remember Me" agar credential tersimpan untuk login otomatis. |
| Pre-Condition | User berada di form login. User sebelumnya pernah centang "Remember Me". |
| Post-Condition | Credential tersimpan securely. Next login akan pre-filled email dan password. |
| Normal Flow | [step-by-step] |
| Alternative Flow | [flow alternatif] |
| Exception Flow | [exception handling] |
| Related Use Case | • Login Via Form |
| Business Rule | [aturan bisnis] |
```

---

## Panduan Penggunaan

### Kapan Membuat DCD
- Ketika akan membuat desain detail dari satu komponen produk
- Sebelum development dimulai untuk clarify requirement
- Saat onboarding developer baru untuk memahami fitur

### Workflow Membuat DCD
1. **Identifikasi Boundary, Control, Entity** dari komponen menggunakan template object identification
2. **Daftar semua use case** dalam Use Case List table
3. **Tulis detail setiap use case** menggunakan template use case detail
4. **Review** dengan product/BA team untuk memastikan lengkap dan akurat
5. **Simpan** di folder `.ai-doc/desain-component-document/`

### Validasi Checklist
- [ ] Object Identification Boundary items dapat dilihat di UI/codebase
- [ ] Object Identification Control items dapat ditelusuri ke services/handlers
- [ ] Object Identification Entity items dapat ditelusuri ke database/API
- [ ] Setiap use case memiliki clear normal flow
- [ ] Setiap use case memiliki exception/error handling
- [ ] Business rules sudah disepakati dengan stakeholder
- [ ] Related use cases sudah terdaftar di setiap use case detail

---

## Catatan Penting

- **Konsisten**: Gunakan format dan naming yang sama untuk semua DCD
- **Detail tapi ringkas**: Setiap section harus informatif tanpa berlebihan
- **Berbasis bukti**: Setiap item harus dapat ditelusuri ke codebase atau requirement
- **Iteratif**: DCD dapat diupdate saat ada perubahan requirement atau discovery baru
- **Non-duplication**: Jangan duplikasi informasi, gunakan cross-reference ke file lain jika perlu

---

## Referensi Template

- **Object Identification**: `.ai-doc/template/object-identification-section-template.md`
- **Use Case Detail**: `.ai-doc/template/usecase-detail-template.md`
- **Use Case Diagram**: `.ai-doc/template/usecase-diagram-template.md`

