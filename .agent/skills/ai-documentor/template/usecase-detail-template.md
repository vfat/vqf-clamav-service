# Template Use Case Detail

## Tujuan
Template ini mendefinisikan struktur untuk mendokumentasikan detail dari setiap use case.
Setiap use case harus memiliki dokumentasi lengkap yang mencakup deskripsi, flow, business rules, dan related use cases.

## Struktur Template

```markdown
| Kategori | Keterangan |
|----------|-----------|
| Use Case Name | [Nama use case yang jelas dan deskriptif] |
| Actor | [Siapa yang melakukan use case ini] |
| Deskripsi | [Penjelasan singkat tentang tujuan use case] |
| Pre-Condition | [Kondisi yang harus terpenuhi sebelum use case bisa dijalankan] |
| Post-Condition | [Kondisi setelah use case berhasil dijalankan] |
| Normal Flow | [Alur normal/happy path step-by-step] |
| Alternative Flow | [Alur alternatif jika ada pilihan lain] |
| Exception Flow | [Alur exception/error handling] |
| Related Use Case | [Daftar use case yang terkait atau berinteraksi] |
| Business Rule | [Aturan bisnis atau constraint yang berlaku] |
```

---

## Penjelasan Setiap Field

### Use Case Name
**Definisi**: Nama identifikasi dari use case ini.

**Contoh**:
- Login Via Form
- Login Via SSO Google
- Sign Up Via Form
- Forgot Password

**Fokus pada**: Nama harus deskriptif dan unik, menunjukkan action dan context yang jelas.

### Actor
**Definisi**: Siapa atau apa yang melakukan use case ini (user, system, external service).

**Contoh**:
- User
- Operator
- System
- Backend API
- MetaTrader Terminal

**Fokus pada**: Dapat berupa human actor (user, operator) atau system actor (API, service).

### Deskripsi
**Definisi**: Penjelasan singkat tentang tujuan dan konteks use case.

**Contoh**:
- "User login menggunakan form email dan password di halaman login."
- "User membuat akun baru dengan mengisi form registrasi."
- "Sistem secara otomatis mendeteksi terminal MetaTrader yang terinstall."

**Fokus pada**: Ringkas, jelas, dan menunjukkan value/tujuan use case.

### Pre-Condition
**Definisi**: Kondisi atau state yang harus ada sebelum use case bisa dijalankan.

**Contoh**:
- "User sudah terdaftar dan berada di halaman login."
- "User sudah login dan berada di halaman settings."
- "MetaTrader terminal sudah terinstall di sistem."

**Fokus pada**: Apa yang harus sudah terjadi agar use case bisa mulai.

### Post-Condition
**Definisi**: Kondisi atau state yang dihasilkan setelah use case berhasil dijalankan.

**Contoh**:
- "User berhasil masuk ke dalam sistem dan mendapatkan akses ke dashboard."
- "Data user tersimpan di database dan dashboard menampilkan data terbaru."
- "MetaTrader terminal status berubah menjadi 'running'."

**Fokus pada**: Apa yang seharusnya terjadi sebagai hasil dari use case yang berhasil.

### Normal Flow
**Definisi**: Alur step-by-step dari use case dalam kondisi normal/ideal (happy path).

**Format**: Numbered list dengan actor yang melakukan setiap step.

**Contoh**:
```
1. User mengisi email dan password.
2. Submit data dengan menekan tombol login.
3. Sistem memvalidasi inputan form.
4. Jika valid, sistem melakukan API call ke backend login.
5. Sistem menerima token dan data user dari API.
6. Sistem menyimpan credential securely.
7. Sistem redirect ke dashboard.
8. Dashboard menampilkan data user.
```

**Fokus pada**: Step harus jelas, urutan harus logis, actor harus disebutkan jika ada perubahan.

### Alternative Flow
**Definisi**: Alur alternatif atau branch dari normal flow jika ada pilihan/opsi lain yang valid.

**Format**: Deskripsi alur alternatif dengan kondisi kapan digunakan.

**Contoh**:
- "Jika User klik 'Login via Google', alihkan ke proses SSO Google."
- "Jika User memilih 'Remember me', simpan credential untuk login otomatis kali berikutnya."
- "Jika User klik 'Forgot Password', alihkan ke halaman reset password."

**Fokus pada**: Variasi valid dari normal flow tanpa error condition.

### Exception Flow
**Definisi**: Alur yang menangani error, validasi gagal, atau kondisi exception lainnya.

**Format**: Kondisi exception dan aksi yang diambil.

**Contoh**:
```
3a. Jika inputan form tidak valid (email/password kosong):
    - Sistem menampilkan error message di form.
    - User dapat memperbaiki inputan dan submit ulang.

5a. Jika validasi API login gagal (wrong password):
    - Sistem menampilkan error message: "Email atau password salah".
    - User dapat retry atau klik 'Forgot Password'.

5b. Jika API timeout:
    - Sistem menampilkan error message: "Koneksi timeout, silakan coba lagi".
    - User dapat retry atau kembali ke halaman awal.
```

**Fokus pada**: Exception harus jelas, action harus konsisten dengan UX yang baik.

### Related Use Case
**Definisi**: Daftar use case lain yang terkait atau berinteraksi dengan use case ini.

**Format**: Bullet list dengan tipe relasi jika perlu (e.g., "alternative to", "extension of", "triggers").

**Contoh**:
- Login Via SSO Google (alternative to)
- Forgot Password (related to)
- Sign Up Via Form (prerequisite: user harus sign up dulu)
- Dashboard (postcondition: user diarahkan ke use case ini)
- Two-Factor Authentication (extension of)

**Fokus pada**: Menunjukkan ketergantungan dan relasi dalam user journey yang lebih besar.

### Business Rule
**Definisi**: Aturan bisnis, constraint, atau requirement yang harus dipenuhi dalam use case ini.

**Format**: Bullet list dengan aturan-aturan spesifik.

**Contoh**:
```
- Email dan Password wajib diisi.
- Email harus memiliki format yang valid (RFC 5322).
- Password minimal 8 karakter.
- Password harus mengandung minimal satu angka [0-9].
- Password harus mengandung minimal satu huruf uppercase [A-Z].
- Karakter yang diizinkan: huruf A-Z, a-z, angka 0-9, dan simbol: ! @ # $ % ^ & * - +
- Login attempt dibatasi maksimal 5x dalam 15 menit.
- Token expiration: 1 jam.
- Refresh token expiration: 7 hari.
```

**Fokus pada**: Aturan yang harus diimplementasikan dalam kode, tidak ada kompromi.

---

## Contoh Lengkap: Login Via Form

| Kategori | Keterangan |
|----------|-----------|
| Use Case Name | Login Via Form |
| Actor | User |
| Deskripsi | User login menggunakan form email dan password di halaman login. |
| Pre-Condition | User sudah terdaftar dan berada di halaman login. |
| Post-Condition | User berhasil masuk ke dalam sistem dan mendapatkan akses ke dashboard. |
| Normal Flow | 1. User mengisi email dan password.<br>2. User menekan tombol "Login".<br>3. Sistem memvalidasi inputan form (format email, length password, required fields).<br>4. Jika valid, sistem melakukan API call ke endpoint POST /api/auth/login dengan credentials.<br>5. Backend memvalidasi credentials terhadap database.<br>6. Jika valid, backend mengembalikan access token, refresh token, dan user data.<br>7. Sistem menyimpan token ke secure storage (encrypted).<br>8. Sistem menyimpan user data ke local storage/database.\<br>9. Sistem redirect ke dashboard page.\<br>10. Dashboard menampilkan data user dan ringkasan finansial. |
| Alternative Flow | **Login via SSO Google**: Jika User klik tombol "Login with Google", alihkan ke SSO Google flow.<br>**Remember Me**: Jika User centang checkbox "Remember me", maka pada login berikutnya credential akan pre-filled. |
| Exception Flow | **3a. Form validation error**: Jika email kosong, password kosong, atau format email invalid:<br>&nbsp;&nbsp;- Sistem menampilkan error message di bawah field yang salah.<br>&nbsp;&nbsp;- Field berubah warna merah dan error message berwarna merah.<br>&nbsp;&nbsp;- User dapat memperbaiki dan submit ulang.<br><br>**5a. Wrong credentials**: Jika email tidak terdaftar atau password salah:<br>&nbsp;&nbsp;- Backend mengembalikan error 401 Unauthorized.<br>&nbsp;&nbsp;- Sistem menampilkan alert: "Email atau password salah. Silakan coba lagi.".<br>&nbsp;&nbsp;- User dapat retry atau klik "Forgot Password".<br><br>**5b. Account locked**: Jika user login gagal 5x dalam 15 menit:<br>&nbsp;&nbsp;- Akun akan di-lock sementara selama 30 menit.<br>&nbsp;&nbsp;- Sistem menampilkan pesan: "Akun Anda terkunci. Silakan coba lagi dalam 30 menit atau reset password.".<br><br>**6a. API timeout/network error**: Jika koneksi ke backend timeout:<br>&nbsp;&nbsp;- Sistem menampilkan error: "Koneksi timeout. Silakan periksa internet Anda dan coba lagi.".<br>&nbsp;&nbsp;- User dapat retry. |
| Related Use Case | • Login Via SSO Google (alternative to)<br>• Forgot Password (linked from exception flow)<br>• Sign Up Via Form (user harus sign up dulu)<br>• Two-Factor Authentication (if enabled)<br>• Logout<br>• Dashboard (post-condition)<br>• Session Management |
| Business Rule | • Email dan password wajib diisi.<br>• Email harus format valid (RFC 5322) dan max 255 karakter.<br>• Password min 8 karakter, max 128 karakter.<br>• Password harus terdiri dari kombinasi: huruf (A-Z, a-z), angka (0-9), dan simbol (!@#$%^&*-+).<br>• Tidak boleh ada whitespace di awal/akhir email atau password.<br>• Login attempt dibatasi maksimal 5x gagal dalam 15 menit → akun lock 30 menit.<br>• Access token expiration: 1 jam.<br>• Refresh token expiration: 7 hari.<br>• Session maksimal per user: 3 device aktif simultaneously.<br>• Token harus disimpan dalam secure storage (encrypted, tidak hardcoded).<br>• API call harus menggunakan HTTPS/TLS 1.3+.<br>• Password harus di-hash menggunakan bcrypt dengan minimum 12 rounds sebelum disimpan. |

---

## Contoh Lengkap: Detect MetaTrader Terminal

| Kategori | Keterangan |
|----------|-----------|
| Use Case Name | Detect MetaTrader Terminal |
| Actor | System (Data Collector Service) |
| Deskripsi | Sistem secara otomatis mendeteksi MetaTrader terminal (MT4/MT5) yang terinstall di komputer user. |
| Pre-Condition | • Data Collector Service berjalan di background.<br>• Windows OS sudah terinstall.<br>• Minimal satu MetaTrader terminal sudah terinstall atau ada registry entry. |
| Post-Condition | • Daftar terminal yang terdeteksi tersimpan di database lokal (SQLite).<br>• Status terminal (running/stopped) ditampilkan di UI MetaTrader monitoring page.<br>• Notification dikirim ke user jika ada terminal baru terdeteksi (optional). |
| Normal Flow | 1. Data Collector Service trigger di startup aplikasi atau setiap 5 menit (background job).<br>2. Service melakukan scan pada Windows Registry untuk mencari entry MetaTrader:<br>&nbsp;&nbsp;- HKLM\SOFTWARE\MetaQuotes Software Corp\Terminal 64-Bit\Installations<br>&nbsp;&nbsp;- HKLM\SOFTWARE\MetaQuotes Software Corp\Terminal\Installations<br>3. Untuk setiap entry, extract informasi: Path, Version, Type (MT4/MT5).<br>4. Service lakukan file system check pada path yang di-extract untuk konfirmasi.<br>5. Untuk setiap terminal yang ditemukan, cek process list untuk status (running/stopped).<br>6. Collect additional info: terminal name, build number, account count (jika bisa diakses).<br>7. Compare dengan data terakhir di database lokal untuk deteksi perubahan (new/removed/updated).<br>8. Simpan data terbaru ke SQLite database (tabel: MetaTraderTerminals).<br>9. Publish event: TerminalsDetected dengan daftar terminal dan status-nya.<br>10. UI MetaTrader monitoring page subscribe ke event ini dan update tampilan. |
| Alternative Flow | **Scan from custom path**: Jika user menginput custom installation path di settings:<br>&nbsp;&nbsp;- Service akan juga melakukan scan pada path custom tersebut.<br>&nbsp;&nbsp;- Jika ditemukan terminal, add ke daftar deteksi.<br><br>**Offline mode**: Jika user offline atau API tidak bisa diakses:<br>&nbsp;&nbsp;- Service tetap melakukan deteksi lokal dari registry dan filesystem.<br>&nbsp;&nbsp;- Data akan disync ke cloud saat internet kembali terhubung. |
| Exception Flow | **1a. Registry access denied**: Jika User Account Control (UAC) membatasi akses registry:<br>&nbsp;&nbsp;- Service coba lakukan dengan elevated privileges (Run as Admin).<br>&nbsp;&nbsp;- Jika tetap gagal, log error dan tampilkan notification: "Akses registry terbatas. Jalankan aplikasi sebagai Administrator."<br><br>**4a. Terminal path tidak valid atau terdeleted**: Jika file system path tidak ditemukan:<br>&nbsp;&nbsp;- Tandai terminal sebagai "offline" atau "uninstalled".<br>&nbsp;&nbsp;- Jangan hapus dari database, tapi mark status sebagai inactive.<br><br>**5a. Process detection error**: Jika API Windows untuk cek process gagal:<br>&nbsp;&nbsp;- Default status ke "unknown".<br>&nbsp;&nbsp;- Log error untuk debugging.<br><br>**Database write error**: Jika SQLite insert/update gagal:<br>&nbsp;&nbsp;- Retry maksimal 3x dengan backoff exponential.<br>&nbsp;&nbsp;- Jika tetap gagal, log error critical dan alert admin. |
| Related Use Case | • Monitor MetaTrader Terminal (triggered by detection result)<br>• Sync Terminal Tasks (uses detected terminals)<br>• MetaTrader Settings Management<br>• Auto-start MetaTrader<br>• Folder Sync (uses terminal data)<br>• Data Collector Main Process |
| Business Rule | • Detection harus dilakukan setiap startup aplikasi dan setiap 5 menit background.<br>• Registry scan harus cover semua standard registry path untuk MT4 dan MT5.<br>• File system verification mandatory untuk confirm terminal path valid.<br>• Process detection harus check: terminal.exe, terminalApp.exe (MT4), metatrader4.exe, metatrader5.exe (MT5).<br>• Maksimal timeout untuk setiap detection attempt: 10 detik.<br>• Data perubahan (new/removed) harus log untuk audit trail.<br>• Terminal data harus di-cache minimal 1 menit untuk avoid redundant scans.<br>• Jika user run aplikasi tanpa admin privilege, detection akan fallback ke limited mode (hanya filesystem scan).<br>• Custom path yang user input harus validate: path format valid, path exist, readable permission. |

---

## Panduan Pembuatan

### Step 1: Identifikasi Use Case Scope
Tentukan **nama use case**, **actor**, dan **tujuan** dengan jelas.

### Step 2: Tulis Pre/Post Condition
- Pre-condition: apa yang harus sudah ada
- Post-condition: apa hasil yang diharapkan

### Step 3: Rancang Normal Flow
- Step-by-step dari start hingga finish
- Sebutkan actor di setiap step jika ada perubahan
- Hindari detail teknis yang terlalu spesifik

### Step 4: Identifikasi Alternative dan Exception Flow
- Alternative: variasi valid tanpa error
- Exception: error handling dan recovery

### Step 5: Daftar Related Use Cases
- Gunakan untuk menunjukkan user journey yang lebih besar

### Step 6: Tentukan Business Rules
- Constraint teknis dan bisnis yang tidak bisa dikompromikan
- Dari requirement gathering dan discussion dengan product/BA

### Step 7: Validasi terhadap Codebase
- Setiap step harus bisa ditelusuri ke actual code/component
- Flow harus konsisten dengan implementasi yang ada atau direncanakan

---

## Kapan Menggunakan Template Ini

1. **Dokumentasi Use Case Detail**: Setiap use case dari use case diagram perlu dokumentasi detail ini.
2. **Requirement Gathering**: Saat discussi dengan product team atau BA untuk clarify requirement.
3. **Development Planning**: Saat planning sprint dan breakdown task dari use case.
4. **QA/Testing**: Untuk membuat test case dari normal flow, alternative, dan exception flow.
5. **Onboarding Developer Baru**: Untuk menjelaskan feature dan business logic.

---

## Catatan Penting

- **Detail tapi ringkas**: Setiap field harus informatif tapi tidak berlebihan.
- **Berbasis bukti**: Setiap flow harus bisa diimplementasikan dengan code yang ada atau direncanakan.
- **Business rule non-negotiable**: Pastikan business rule sudah disepakati dengan stakeholder.
- **Flow harus clear**: Gunakan language yang mudah dipahami, hindari ambiguitas.
- **Test case ready**: Normal flow, alternative flow, dan exception flow harus bisa langsung di-convert ke test case.

---

## Referensi

- **Use Case Writing Best Practices**: https://www.ib.edu.au/the-craft-of-use-cases/
- **UML Use Case Diagram**: https://www.uml-diagrams.org/use-cases.html
