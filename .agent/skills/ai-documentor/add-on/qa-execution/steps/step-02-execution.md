# Step 2: Execution Loop

> Jalankan item satu per satu sesuai mode-nya, tentukan status baru HANYA dari
> bukti nyata, dan update control plane secara incremental.
> **HALT** — setiap FAILED/BLOCKED dikonfirmasi sebelum lanjut item berikutnya.

---

## Input dari User

- **Mode per item** — setuju dengan mode usulan agent / minta ganti
- **Hasil langkah manual** — untuk item manual/hybrid
- **Keputusan lanjut/berhenti** — saat ada FAILED atau BLOCKED

---

## Yang Dilakukan Agent

### 1. Tentukan Mode per Item

| Jenis di Plan | Mode Default | Alasan |
|---|---|---|
| Unit / Integration | Automated | Suite biasanya sudah ada |
| E2E / API | Automated bila suite ada; manual bila belum | Tergantung kesiapan repo |
| UX / Visual / Eksploratif | Manual | Sulit diautomasi |

Sampaikan mode usulan sekali di awal batch; user bisa override per item.

### 2. Eksekusi Automated

Jalankan command nyata di terminal, tangkap exit code + output:

```text
$ npx playwright test login-gagal --project=chromium
Exit code: 0
Output: 1 passed (12.3s)
```

Aturan penangkapan bukti:

- Command ditulis persis seperti dijalankan (termasuk filter/prefix).
- Exit code dicatat apa adanya — termasuk non-zero.
- Output diringkas tapi bagian menentukan (pass/fail summary) dikutip utuh.
- **Dilarang mengarang output** — kalau command gagal jalan, itu data, bukan aib.

### 3. Eksekusi Manual

Pandu user langkah demi langkah:

> **QA-006** [Manual] Verifikasi pesan error tampil dalam Bahasa Indonesia
> ketika server mati.
>
> Langkah:
> 1. Matikan service backend
> 2. Buka halaman login, isi form, submit
> 3. Amati pesan yang muncul
>
> Kabari saya hasilnya: sesuai harapan / tidak sesuai / tidak bisa dites.

Catat konfirmasi user sebagai bukti (`user-confirmed`, tanggal).

### 4. Tentukan Status Baru dari Bukti

| Bukti yang Ada | Status Sah |
|---|---|
| Exit code 0 + output lolos ATAU konfirmasi manual user | `PASSED` |
| Output kegagalan + hipotesis penyebab | `FAILED` |
| Alasan environment/data konkret | `BLOCKED` |
| Alasan + target waktu | `DEFERRED` |

Untuk FAILED, susun hipotesis:

> QA-005 **FAILED** — exit code 1, `expect(response.status).toBe(401)`
> menerima 500. Hipotesis: error handler melempar exception sebelum
> mengembalikan 401. Mau saya telusuri sekarang, atau lanjut item berikut?

**HALT** pada setiap FAILED/BLOCKED — konfirmasi dulu sebelum lanjut.

### 5. Update Incremental

Setiap item selesai, langsung:

1. Perbarui baris item di §7 test plan (Status + Bukti).
2. Perbarui registry qa-overview.md + Progress Summary.
3. Tambah baris Change Log.

Jangan menumpuk semua update di akhir sesi — kalau sesi putus di tengah,
control plane harus tetap konsisten dengan kondisi terakhir.

### 6. Interaksi TDD

Jika constitution menyatakan `TDD: Enabled` dan item berlevel unit:

- Eksekusi lewat siklus RED→GREEN add-on tdd.
- qa-execution hanya mencatat referensi: ID target TDD + status akhirnya.
- JANGAN menduplikasi siklus RED/GREEN di qa-overview.

---

## Output Step 2

```yaml
session_state:
  results:
    - id: "QA-005"
      mode: "automated"
      command: "npx playwright test login-gagal"
      exit_code: 1
      new_status: "FAILED"
      hypothesis: "error handler lempar exception sebelum 401"
    - id: "QA-006"
      mode: "manual"
      evidence: "user-confirmed @ 2026-08-25"
      new_status: "PASSED"
  halt_pending_decision: null | "QA-004"
```

---

## Aturan

- **Tidak ada transisi status tanpa bukti** — sama kerasnya dengan aturan TDD.
- **JANGAN menjalankan command destruktif** tanpa persetujuan eksplisit
  (drop database, reset environment, load test ke produksi).
- **JANGAN batch-update di akhir sesi** — update incremental wajib.
- **FAILED dicatat jujur** — jangan haluskan atau sembunyikan output gagal.
- Item BLOCKED butuh alasan konkret, bukan "sementara belum bisa".
