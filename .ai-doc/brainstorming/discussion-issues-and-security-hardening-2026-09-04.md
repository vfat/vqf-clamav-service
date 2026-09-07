# Discussion Summary — Penanganan Issues & Security Hardening (Vault AES-256-GCM, ConcurrentDatabaseReload & Anti-Memory Bomb)

| Field | Value |
|---|---|
| **Sesi** | Pembahasan Issues & Technical Debt `clamav-service` |
| **Area** | Troubleshooting & Security Hardening |
| **Sub-Topik** | Issue #1 (Vault AES-256-GCM & Backward Compatibility), Issue #2 (ConcurrentDatabaseReload no), Issue #3 (Quarantine Anti-Memory Bomb) |
| **Tanggal** | 2026-09-04 |
| **Kategori** | Security Remediation, Cryptography, Memory Protection, Daemon Configuration |
| **Teknik** | Five Whys, Risk Assessment, Comparative Trade-off, Solution Matrix |
| **Peserta** | User, Dinda 🎯 (Moderator), Melon 🏗️, Sultan ⚙️, Nindi 🔬, Lugi 📊, Bernadya ✨ |
| **Jumlah Ide** | 3 Ide (3 Approved, 0 Rejected, 0 Modified) |

---

## Konteks

Sesi brainstorming interaktif dipandu oleh **Dinda 🎯** bersama Full Team untuk membahas isu krusial pada [`.ai-doc/3p.md`](file:///home/ubuntu/workspace/plan/clamav-service/.ai-doc/3p.md) Bagian 3 (*Pending Issues*), khususnya masking statis XOR `0xA5` pada `internal/quarantine/vault.go` yang belum memenuhi standar kriptografi, serta isu tambahan yang diidentifikasi terkait risiko OOM daemon ClamAV (`ConcurrentDatabaseReload`) dan Memory Bomb pada buffer stream karantina.

---

## Ide yang Teridentifikasi

### 1. Self-Describing Storage via Magic Header (`VQF_AESGCM_V1\n`) & Dual Fallback Reader
File karantina baru disimpan dengan header 14-byte `VQF_AESGCM_V1\n` + Nonce 12-byte CSPRNG + Ciphertext AES-256-GCM + Auth Tag 16-byte dengan hak akses disk `0600`. Reader pada fungsi `RestoreFile` dan `DownloadFile` secara otomatis mendeteksi keberadaan magic header: jika cocok, lakukan dekripsi AES-GCM; jika absen, otomatis fallback ke algoritma XOR `0xA5` legacy untuk menjamin zero data loss tanpa skrip batch migration.
> **Sumber:** Sultan ⚙️ & Nindi 🔬 | **Status:** `approved` | **Kategori:** Storage Cryptography

### 2. Inject `ConcurrentDatabaseReload no` pada `configs/clamd.conf`
Menambahkan parameter `ConcurrentDatabaseReload no` pada konfigurasi daemon ClamAV. Mengatasi default bawaan `yes` yang memicu pemuatan database signature baru dan lama secara bersamaan ke RAM (~1.8GB - 2.4GB) yang memicu OOM Killer pada VPS 1GB/2GB.
> **Sumber:** User & Sultan ⚙️ | **Status:** `approved` | **Kategori:** Resource Limit / Daemon Config

### 3. Quarantine Anti-Memory Bomb (Strict LimitReader + Disk Spooling Staging)
Mengganti `io.ReadAll(r)` bebas dengan `io.LimitReader` di level Vault dan memanfaatkan temporary disk spooling (buffer 32 KB) saat menerima stream file infeksi. Mencegah alokasi heap RAM berlipat ganda saat file besar atau stream tak terbatas masuk ke karantina.
> **Sumber:** User, Nindi 🔬, & Melon 🏗️ | **Status:** `approved` | **Kategori:** Memory Protection / Reliability

---

## Pembahasan

### Pembahasan 1: Backward Compatibility & Determinisme Magic Header
Tim membedah alur tulis dan baca file di disk level byte. Magic header `VQF_AESGCM_V1\n` memiliki determinisme tinggi karena karakter pertama ASCII `'V'` (`0x56`) mustahil terbentuk secara tidak sengaja dari file lama yang di-XOR dengan `0xA5` (`0x56 ^ 0xA5 = 0xF3`). File lama tetap bisa di-restore 100% tanpa risiko crash.
**Keputusan:** Disetujui (Approved). File baru berformat AES-256-GCM dengan magic header; file lama di-handle via fallback reader.

### Pembahasan 2: Efek Samping `ConcurrentDatabaseReload no`
Sultan ⚙️ dan Melon 🏗️ mengklarifikasi bahwa saat diset `no`, `clamd` akan meng-unload database lama sebelum memuat database baru, menimbulkan jeda proses scanning sementara sekitar 5–15 detik saat freshclam selesai update harian. Request yang masuk akan tertahan di socket backlog queue (context timeout 30 detik) dan dilayani segera setelah reload selesai.
**Keputusan:** Disetujui (Approved). Trade-off jeda 10 detik sekali sehari jauh lebih baik daripada daemon mati permanen akibat OOM Killer.

### Pembahasan 3: Proteksi Memory Bomb di Quarantine
Nindi 🔬 dan Melon 🏗️ mengidentifikasi bahwa pembacaan langsung seluruh file ke memori sebelum enkripsi menyebabkan RAM Go menahan byte slice ganda. Dengan `io.LimitReader` dan buffer terkontrol, heap footprint Go tetap rendah (< 30 MB).
**Keputusan:** Disetujui (Approved). Terapkan batas ukuran ketat dan streaming spooling.

---

## Keputusan Final

- ✅ **Vault AES-256-GCM:** Mengganti XOR scrambling dengan AES-256-GCM authenticated cipher menggunakan `crypto.EncryptAESGCM` dan `crypto.DecryptAESGCM`.
- ✅ **Magic Header Identification:** Prefix 14-byte `VQF_AESGCM_V1\n` untuk membedakan file baru dan file lama tanpa migrasi batch.
- ✅ **Master Key Binding:** Constructor `quarantine.NewVault` menerima `masterKey []byte` dari `crypto.EnsureMasterKey`.
- ✅ **Anti-OOM clamd:** Tambahkan `ConcurrentDatabaseReload no` ke `configs/clamd.conf`.
- ✅ **Anti-Memory Bomb:** Ganti `io.ReadAll` dengan `io.LimitReader` dan streaming buffer di `internal/quarantine/vault.go`.

---

## Next Steps

- [ ] **TDD Unit Test Suite (TDD-006B)** — Owner: Sultan ⚙️ & Nindi 🔬 | Target: `internal/quarantine/vault_test.go`
- [ ] **Implementasi Magic Header & AES-GCM di Vault** — Owner: Sultan ⚙️ | Target: `internal/quarantine/vault.go`
- [ ] **Wiring Master Key di main.go** — Owner: Melon 🏗️ | Target: `cmd/server/main.go`
- [ ] **Konfigurasi clamd.conf ConcurrentDatabaseReload** — Owner: Sultan ⚙️ | Target: `configs/clamd.conf`
- [ ] **Update Dokumentasi 3P & Lampiran L-002** — Owner: Dinda 🎯 | Target: `.ai-doc/3p.md`

---

## Referensi & Context
- [`.ai-doc/3p.md`](file:///home/ubuntu/workspace/plan/clamav-service/.ai-doc/3p.md)
- [`internal/quarantine/vault.go`](file:///home/ubuntu/workspace/plan/clamav-service/internal/quarantine/vault.go)
- [`internal/crypto/aes.go`](file:///home/ubuntu/workspace/plan/clamav-service/internal/crypto/aes.go)
- [`configs/clamd.conf`](file:///home/ubuntu/workspace/plan/clamav-service/configs/clamd.conf)

---
*Generated by AI Documentor — Brainstorming Add-On*
