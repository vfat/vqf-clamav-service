# Template Object Identification Section

## Tujuan
Template ini mendefinisikan struktur untuk membuat section identifikasi object dalam dokumentasi codebase.
Section ini membantu mengklasifikasikan dan mendokumentasikan komponen sistem berdasarkan kategori **Boundary**, **Control**, dan **Entity**.

## Struktur Section

```markdown
### Boundary
- Screen, Form, Modal, UI Component
- Gateway, IPC Channel
- File System Interface

### Control
- Service, Handler, Validator
- Manager, Repository
- Business Logic, Event Handler

### Entity
- Database Table, DTO
- Domain Model, Configuration
- Cache, External Data Source
```

**Penjelasan**:
- Section ini terdiri dari 3 bagian utama: **Boundary**, **Control**, **Entity**
- Setiap section berisi list item yang termasuk dalam kategori tersebut
- Gunakan struktur ini sebagai **reference definition** ketika mendokumentasikan fitur/komponen tertentu.

---

## Deskripsi Kategori

### Boundary
**Definisi**: Antarmuka atau titik masuk ke dalam sistem.

**Yang Termasuk**:
- Screen/Form (UI layer)
- Modal/Dialog
- UI Component
- External Gateway (UI-based interaction)
- File System Interface (user-triggered file operations)
- IPC Channel (direct user/system interaction)

**Fokus pada**: Apa yang user/sistem gunakan untuk berinteraksi dengan aplikasi? (UI-centric interfaces)

### Control
**Definisi**: Alur proses, aturan bisnis, dan koordinasi.

**Yang Termasuk**:
- Service/Handler
- Validator/Validator Logic
- Business Logic/Rules
- Event Handler
- Manager
- Repository/Gateway
- Operation/Command

**Fokus pada**: Apa yang terjadi di balik layar? Bagaimana data diproses, divalidasi, atau dikoordinasikan?

### Entity
**Definisi**: Daftar data yang diakses atau disimpan oleh sistem.

**Yang Termasuk**:
- Database Table
- API Response/DTO
- Domain Model
- Configuration/Settings
- External Data Source
- Cache/Memory
- Log Files
- **API Endpoint** (data/service endpoint yang di-akses oleh sistem)

**Fokus pada**: Apa data/entitas yang diakses atau dimanipulasi oleh sistem?

---

## Panduan Pembuatan

### Step 1: Identifikasi Scope
Tentukan **komponen produk**, **fitur**, atau **use case** yang akan didokumentasikan.

Contoh:
- "Login Feature"
- "Dashboard Component"
- "MetaTrader Integration"
- "Trade Execution Flow"

### Step 2: Daftar Boundary
Tanyakan: **Bagaimana user atau sistem eksternal berinteraksi dengan komponen ini?**

Contoh untuk Login Feature:
- LoginWindow.xaml (screen)
- AuthenticationViewModel (UI layer)
- Google OAuth Callback URI (gateway)
- Error/Success Modal (modal dialog)

### Step 3: Daftar Control
Tanyakan: **Apa proses, validasi, atau logika bisnis yang terjadi di balik layar?**

Contoh untuk Login Feature:
- AuthenticationService (service)
- ValidateEmailFormat() (validator)
- ValidatePassword() (validator)
- HandleAuthenticationResponse() (handler)
- GenerateSessionToken() (operation)

### Step 4: Daftar Entity
Tanyakan: **Data apa saja yang diakses, disimpan, atau dimanipulasi?**

Contoh untuk Login Feature:
- User (database entity)
- UserDTO (API DTO)
- AppSettings (configuration)
- Session (in-memory token)
- Backend Auth API Response

### Step 5: Validasi terhadap Codebase
- **Boundary**: Apakah ada UI component, form, dialog, gateway UI, atau channel interaksi di codebase? (cek `Views/`, `ViewModels/`, `Controllers/`, `Features/*/Presentation/`)
- **Control**: Apakah ada service, handler, atau validator yang sesuai? (cek `Services/`, `Features/*/Application/`, `Core/Services/`)
- **Entity**: Apakah ada model, DTO, atau database schema yang sesuai? (cek `Models/`, `Entities/`, `Core/DTO/`, `Migrations/`)

### Step 6: Sesuaikan dengan Codebase
Gunakan nama actual dari codebase, bukan nama generik.

---

## Contoh Lengkap

### Login Feature

#### Boundary
- LoginWindow.xaml
- AuthenticationViewModel
- Google OAuth Callback URI
- Error/Success Modal

#### Control
- AuthenticationService
- ValidateEmailFormat()
- ValidatePassword()
- HandleAuthenticationResponse()
- GenerateSessionToken()
- EncryptCredentials()

#### Entity
- User (database entity)
- UserDTO (API DTO)
- AppSettings (credentials storage)
- Session (in-memory token)
- Backend Auth API Response
- POST /api/auth/login (API endpoint)
- GET /api/auth/verify (API endpoint)

---

### Dashboard Feature

#### Boundary
- DashboardView.xaml
- DashboardViewModel
- RefreshButton Click
- NotificationPanel

#### Control
- DashboardService
- FetchUserSummary()
- ApplyFinancialCalculations()
- AggregateMetrics()
- GenerateSummaryReport()
- SubscribeToRealtimeUpdates()

#### Entity
- Dashboard (database entity)
- DashboardDTO (API DTO)
- Account (database entity)
- Trade (database entity)
- AccountBalance (calculated entity)
- MetaTrader Terminal Data
- Settings (user preferences)
- GET /api/dashboard/summary (API endpoint)
- WebSocket /dashboard/updates (API endpoint)

---

### MetaTrader Integration

#### Boundary
- MetaTraderMonitorView.xaml
- TerminalStatusIndicator
- Terminal Start/Stop Button
- AutoStart Toggle
- Terminal Task Grid
- IPC Channel (client)

#### Control
- MetaTraderService
- DetectTerminalInstances()
- SyncTerminalState()
- ValidateTerminalPath()
- HandleTerminalError()
- SubscribeToTerminalEvents()
- FolderSyncManager()

#### Entity
- MetaTrader Terminal Instance
- Terminal Configuration (stored)
- Terminal Tasks (database)
- MT Installation Log
- MQL5 Log File
- Windows Registry (MT settings)
- Local File System

---

## Alat dan Format

- Format: Markdown section dengan title dan list
- Tools: VS Code, atau editor markdown lainnya
- Validasi: Pastikan setiap item dapat ditelusuri ke codebase

## Kapan Menggunakan Section Ini

1. **Dokumentasi Komponen**: Setiap kali mendokumentasikan komponen produk baru.
2. **Desain Component Document**: Di section identifikasi object/interface.
3. **Analisis Use Case**: Untuk menunjukkan implementasi teknis dari use case.
4. **Refactoring Planning**: Untuk memahami scope dan dependency sebelum refactoring.

## Kapan Update Template

- Jika ada kategori baru yang diperlukan (misalnya "Shared Library"), update template ini.
- Jika ada perubahan konvensi naming atau klasifikasi, catat di constitution.
- Jika ada contoh baru yang relevan ditemukan, tambahkan ke section "Contoh Lengkap".

---

## Catatan Penting

- **Jangan campur**: Boundary ≠ Control ≠ Entity. Masing-masing punya peran yang berbeda.
- **Berbasis Bukti**: Setiap item HARUS dapat ditelusuri ke codebase, bukan asumsi.
- **Fokus Scope**: Jika list terlalu panjang (>15 item), pertimbangkan membagi menjadi sub-fitur atau feature area terpisah.
- **Iteratif**: Section ini dapat diupdate saat eksplorasi codebase berlanjut.
