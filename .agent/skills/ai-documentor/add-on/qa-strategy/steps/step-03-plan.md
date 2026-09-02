# Step 3: Test Planning

> Turunkan strategy menjadi test plan konkret untuk satu sprint/release.
> Step ini opsional — jalankan hanya jika user memintanya.
> **HALT** — konfirmasi plan sebelum file ditulis.

---

## Input dari User

- **Target plan** — nomor sprint atau nama release (wajib)
- **Daftar fitur/scope** — fitur apa saja yang masuk target tsb (wajib)
- **Kapasitas** — opsional; batasan waktu/orang jika user mau estimasi realistis

---

## Yang Dilakukan Agent

### 1. Validasi Scope Terhadap Strategy

Untuk setiap fitur yang user sebutkan, cek:

- Fitur ada di artefak bukti (`DCD-*`, `rest-api-doc`) → langsung bisa dipetakan.
- Fitur TIDAK ada di artefak → tandai `Perlu Dikonfirmasi` dan tanyakan ke user;
  jangan mengarang use case-nya.

### 2. Mapping Feature → Test → Prioritas

Prioritas mengikuti risk matrix Step 2:

| Feature / Use Case | Jenis Test | Level | Prioritas | Risiko Terkait | Bukti |
|---|---|---|---|---|---|
| `<fitur>` | `<apa yang dites>` | Unit/Integration/E2E | P0/P1/P2 | `<R-00x>` | `<artefak>` |

Aturan prioritas:

- **P0** — menutup risiko skor tinggi (Impact × Likelihood besar) atau journey kritikal.
- **P1** — risiko menengah, endpoint publik, integritas data.
- **P2** — sisanya; kandidat `Deferred` jika kapasitas kurang.

### 3. Estimasi Effort Relatif

Gunakan ukuran S/M/L (tanpa jam, kecuali user minta):

| Ukuran | Arti praktis |
|---|---|
| S | Satu test case sederhana, data mudah |
| M | Perlu fixture/factory data, beberapa assertion |
| L | Perlu environment khusus, mock service, atau setup kompleks |

### 4. Tentukan Item Deferred

Item yang tidak muat di target:

- Tandai `Deferred` **dengan alasan** (kapasitas / menunggu fitur stabil).
- Jangan hapus diam-diam — deferred tetap tercatat di plan.

### 5. Konfirmasi Plan Sebelum Menulis

Tampilkan ringkasan:

```
=== TEST PLAN <target> ===
Total item: N
- P0: x  |  P1: y  |  P2: z
- Deferred: d
Estimasi: S=a, M=b, L=c
```

> Siap saya tulis ke `.ai-doc/qa/strategy/plans/test-plan-{slug}-{tanggal}.md`?
> - ✅ **Tulis** — generate dari template
> - ✏️ **Revisi** — ubah mapping/prioritas/deferred
> - ❌ **Batal**

**HALT** — tunggu konfirmasi. Setelah disetujui, generate file dari
`template/test-plan-template.md`.

---

## Output Step 3

```yaml
session_state:
  plan_target: "<sprint/release>"
  items:
    - feature: "..."
      test_type: "..."
      level: "unit|integration|e2e"
      priority: "P0|P1|P2"
      related_risk: "R-00x"
      effort: "S|M|L"
      status: "planned" | "deferred"
  plan_file: ".ai-doc/qa/strategy/plans/test-plan-<slug>-<date>.md"
```

---

## Aturan

- **JANGAN menambah scope baru** yang tidak ada di strategy maupun input user.
- **JANGAN memetakan use case yang tidak ada buktinya** — tandai `Perlu Dikonfirmasi`.
- **Deferred ≠ dihapus** — selalu sertakan alasan.
- Plan adalah turunan strategy; jika strategy berubah, plan harus disinkronkan.
