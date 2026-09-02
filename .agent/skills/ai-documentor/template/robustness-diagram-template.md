# Template Robustness Diagram

## Tujuan
Template ini mendefinisikan struktur dan konvensi untuk membuat diagram robustness (Robustness Diagram) dalam dokumentasi codebase. Diagram ini digunakan untuk mendokumentasikan interaksi antara elemen UI (boundary), logic (control), dan data (entity) dalam skenario utama dan alternatif.

## Prinsip
- Gunakan elemen `boundary` untuk UI, `control` untuk proses logika, dan `entity` untuk data.
- Sertakan title yang jelas dan deskripsi skenario.
- Label setiap alur dengan nomor dan keterangan [Sukses] / [Alt A] / [Alt B] dll.
- Hindari inferring dari nama file saja; semua klaim harus traceable ke code.
- Gunakan `skinparam packageStyle rect` untuk tampilan yang rapi.
- Update `.ai-doc/3p.md` setelah langkah yang bermakna.

## Struktur Diagram

```plantuml
@startuml
top to bottom direction
skinparam packageStyle rect

title Robustness Diagram: Skenario Login (Utama & Alternatif)

' --- BOUNDARY ELEMENTS (UI) ---
boundary "Halaman Login" as LoginUI
boundary "Halaman Beranda (Dashboard)" as HomeUI
boundary "Halaman Reset Password" as ResetUI
boundary "Notifikasi Akun Terkunci" as LockUI

' --- CONTROL ELEMENTS (LOGIC) ---
control "Proses Validasi Input" as ValidasiInputCtrl
control "Proses Cek Kredensial" as CekKredensialCtrl
control "Proses Penghitung Gagal Login" as CounterCtrl
control "Proses Kunci Akun" as LockCtrl

' --- ENTITY ELEMENTS (DATA) ---
entity "Data Pengguna (Database)" as UserDB

' --- HUBUNGAN SKENARIO UTAMA (SUKSES) ---
LoginUI -- ValidasiInputCtrl : 1. Input Username & Password
ValidasiInputCtrl -- CekKredensialCtrl : 2. Teruskan Data Valid
CekKredensialCtrl -- UserDB : 3. Cocokkan Kredensial
CekKredensialCtrl -- HomeUI : 4. [Sukses] Alihkan ke Beranda

' --- HUBUNGAN JALUR ALTERNATIF (GAGAL / ERROR) ---
ValidasiInputCtrl -- LoginUI : [Alt A] Format Input Salah (Tampilkan Error)

CekKredensialCtrl -- CounterCtrl : [Alt B] Password Salah (Tambah Counter Gagal)
CounterCtrl -- LoginUI : [Alt B1] Gagal < 3 Kali (Tampilkan Pesan Salah)
CounterCtrl -- LockCtrl : [Alt B2] Gagal >= 3 Kali

LockCtrl -- UserDB : 5. Ubah Status Akun Jadi 'Locked'
LockCtrl -- LockUI : 6. Alihkan ke Halaman Blokir

LoginUI -- ResetUI : [Alt C] Klik "Lupa Password"

@enduml
```
