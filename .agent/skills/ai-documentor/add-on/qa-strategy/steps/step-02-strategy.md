# Step 2: Strategy Building

> Susun strategi testing berbasis bukti Step 1: risk matrix, test pyramid,
> quality gates, entry/exit criteria.
> **HALT** — review strategy bersama user sebelum ditulis ke file.

---

## Input dari User

- **Konfirmasi prioritas risiko** — user boleh menaikkan/menurunkan skor hasil ekstraksi
- **Preferensi quality gate** — seberapa ketat merge/deploy diblokir
- **Permintaan rekomendasi tool** — opsional; default tidak membahas tool

---

## Yang Dilakukan Agent

### 1. Susun Risk Matrix

Setiap baris wajib punya sumber bukti dari `evidence_inventory`:

| ID | Risiko | Komponen/Area | Impact (1–5) | Likelihood (1–5) | Skor | Sumber Bukti |
|---|---|---|---|---|---|---|
| R-001 | `<deskripsi>` | `<komponen>` | `<1–5>` | `<1–5>` | `I×L` | `<artefak + bagian>` |

Aturan penilaian:

- Impact = seberapa parah jika benar-benar terjadi di produksi.
- Likelihood = frekuensi perubahan + riwayat masalah + kompleksitas integrasi.
- Risiko dengan confidence `asumsi` tetap boleh masuk, tapi ditandai jelas.

Tampilkan matrix ke user:

> Berikut risk matrix awal. Ada yang ingin dinaikkan/menurunkan prioritasnya?
> Risiko R-00x masih `Asumsi` — mau konfirmasi atau pertahankan sebagai asumsi?

**HALT** — setelah user setuju pada matrix.

### 2. Tentukan Test Pyramid

Usulkan proporsi per level dengan alasan dari karakter sistem:

| Level | Proporsi usulan | Alasan (dari bukti) |
|---|---|---|
| Unit | `<%>` | `<misal: logika dominan di service layer>` |
| Integration | `<%>` | `<misal: banyak integrasi DB/pihak ketiga>` |
| E2E | `<%>` | `<misal: journey kritikal lewat UI publik>` |

Sistem API-heavy biasanya mengecilkan E2E; UI-heavy sebaliknya.
Jangan pakai angka generik tanpa alasan.

### 3. Tetapkan Quality Gates

| Gate | Lokasi | Ambang / Syarat | Aksi jika gagal |
|---|---|---|---|
| `<misal unit test>` | PR CI | `<misal lulus semua + coverage ≥ X%>` | Block merge |
| `<misal smoke E2E>` | Pre-deploy staging | `<journey kritikal lulus>` | Block deploy |

Ambang harus objektif dan dapat dicek mesin, bukan "cukup baik".

### 4. Definisikan Entry/Exit Criteria

Per level test:

- **Entry** — syarat mulai menulis/menjalankan level tsb (misal: fixture tersedia).
- **Exit** — syarat menyatakan level tsb selesai (misal: semua case prioritas P0/P1 lulus).

### 5. Rekomendasi Tool (Opsional)

Hanya jika user meminta:

> Rekomendasi default: **[tool]** karena [alasan terkait stack di project].
> Alternatif: [tool lain] — pilih jika [kondisi].

Satu default + satu alternatif cukup.

### 6. Review Strategy Sebelum Menulis

Tampilkan ringkasan strategy lengkap, lalu:

> Siap saya tulis ke `.ai-doc/qa/strategy/test-strategy.md`?
> - ✅ **Tulis** — generate dari template
> - ✏️ **Revisi** — sebutkan section mana
> - ❌ **Batal**

**HALT** — tunggu konfirmasi. Setelah disetujui, generate
`.ai-doc/qa/strategy/test-strategy.md` dari `template/test-strategy-template.md`.

---

## Output Step 2

```yaml
session_state:
  risk_matrix: [...]
  pyramid:
    unit: "..."
    integration: "..."
    e2e: "..."
  quality_gates: [...]
  entry_exit_criteria: [...]
  tools: []            # kosong kecuali diminta
  strategy_file: ".ai-doc/qa/strategy/test-strategy.md"
```

---

## Aturan

- **JANGAN menulis file sebelum user konfirmasi.**
- **JANGAN memberi skor risiko tanpa sumber bukti** — gunakan tanda `Asumsi`.
- **WAJIB ambang gate objektif** — dapat dicek mesin, bukan opini.
- **Tool hanya bila diminta** — jangan membebani strategy dengan diskusi tool.
