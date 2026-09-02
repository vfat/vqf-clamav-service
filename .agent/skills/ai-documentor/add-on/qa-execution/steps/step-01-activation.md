# Step 1: Activation

> Validasi gate: cari test plan valid, deteksi runner yang tersedia,
> usulkan urutan eksekusi. **HALT** — konfirmasi batch item sebelum loop mulai.

---

## Input dari User

- **Test plan target** — nama file `test-plan-{slug}-{date}.md` (opsional;
  jika ada lebih dari satu, agent tampilkan daftar dan minta pilih)
- **Batch** — semua item / hanya P0 / item tertentu (default usulan: P0 dulu)

---

## Yang Dilakukan Agent

### 1. Validasi Gate

- Glob `.ai-doc/qa/strategy/plans/test-plan-*.md`.
- Jika tidak ditemukan, cek plan legacy `.ai-doc/test-plan-*.md`, gunakan sebagai
  referensi, dan tawarkan pemindahan ke `qa/strategy/plans/` sebelum eksekusi.
- Untuk tiap kandidat, cek: ada section Test Items dengan ID + kolom Status?
- Jika tidak ada file sama sekali ATAU file tidak valid:

> Saya tidak menemukan test plan yang valid di `.ai-doc/qa/strategy/plans/`. qa-execution butuh
> test plan sebagai sumber item — tanpa itu, status tidak punya jangkar.
> - ✅ **Jalankan qa-strategy dulu** — saya arahkan ke add-on tersebut
> - 📄 **Pakai plan lain** — sebut lokasi filenya
> - ❌ **Batal**

**HALT** — jangan pernah melanjutkan tanpa plan valid.

### 2. Baca Item Test Plan

Ekstrak dari section Test Items (§3) dan Hasil Eksekusi (§7):

| Data | Dipakai untuk |
|---|---|
| ID + Feature | Identitas item |
| Jenis / Level | Menentukan mode eksekusi (manual/automated/hybrid) |
| Prioritas (P0/P1/P2) | Urutan eksekusi |
| Effort (S/M/L) | Estimasi sesi |
| Status saat ini | Skip yang sudah PASSED kecuali user minta re-run |
| Bukti lama | Konteks re-run |

### 3. Deteksi Runner

Cari command test yang nyata tersedia di project:

| Prioritas | Sumber | Contoh |
|---|---|---|
| 1 | `package.json` scripts | `"test"`, `"test:e2e"` |
| 2 | Makefile / justfile | `make test`, `just e2e` |
| 3 | Config framework | `pytest.ini`, `playwright.config.ts`, `go.mod` |
| 4 | Catatan tooling | Section Tooling di test strategy / constitution |

Jika tidak ketemu runner untuk jenis item automated:

> Item [ID] seharusnya automated, tapi saya tidak menemukan suite test-nya.
> - 🔧 **Tulis testnya dulu** — dengan persetujuan Anda (unit wajib lewat TDD bila Enabled)
> - 🖐️ **Turunkan ke manual** — Anda jalankan langkahnya, saya catat
> - ⏸️ **BLOCKED** — catat dengan alasan konkret

### 4. Usulkan Urutan & Konfirmasi Batch

> Dari `test-plan-auth-2026-08-25.md` saya menemukan **12 item**:
> P0 = 5, P1 = 4, P2 = 3. Yang sudah PASSED = 2 (akan di-skip).
> Runner terdeteksi: `npm test` (unit), `npx playwright test` (E2E).
>
> Usulan batch sesi ini: **5 item P0** (estimasi ±30 menit).
> - ✅ **Setuju batch P0**
> - ✏️ **Ubah batch** — sebutkan itemnya
> - 🚀 **Semua item sekaligus**

**HALT** — tunggu konfirmasi.

---

## Output Step 1

```yaml
session_state:
  plan_file: ".ai-doc/qa/strategy/plans/test-plan-{slug}-{date}.md"
  items_total: 12
  items_batch: ["QA-001", "QA-002", "..."]
  items_skipped: ["QA-00x"]        # sudah PASSED
  runners:
    unit: "npm test"
    e2e: "npx playwright test"
    missing: []                     # jenis tanpa runner terdeteksi
  execution_order: "p0-first"       # atau custom dari user
```

---

## Aturan

- **JANGAN melanjutkan tanpa test plan valid** — gate ini keras.
- **JANGAN mengarang command runner** — hanya pakai yang terdeteksi di project;
  kalau ragu, tanya user.
- **REKOMENDASIKAN urutan P0 dulu** — tapi keputusan batch milik user.
- Skip item berstatus PASSED kecuali user minta re-run eksplisit.
