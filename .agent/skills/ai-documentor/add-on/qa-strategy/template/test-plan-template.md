# Test Plan — <Sprint / Release>

> Turunan langsung dari `.ai-doc/qa/strategy/test-strategy.md`. Jangan menambah scope baru.
> Dibuat oleh add-on qa-strategy.

## 1. Metadata

- **Project:** `<nama project>`
- **Target:** `<sprint N / release X.Y>`
- **Strategy referensi:** `.ai-doc/qa/strategy/test-strategy.md`
- **Created at:** `<YYYY-MM-DD>`
- **Status:** `Draft` | `Active` | `Executed`

## 2. Scope

**Masuk plan:**

- `<fitur/use case 1>`
- `<fitur/use case 2>`

**Tidak masuk (Deferred):**

| Item | Alasan | Target Berikutnya? |
|---|---|---|
| `<item>` | `<kapasitas / menunggu stabil>` | Ya / Tidak |

## 3. Test Items

> Prioritas mengikuti risk matrix di strategy. Bukti = artefak asal item.

| ID | Feature / Use Case | Jenis Test | Level | Prioritas | Risiko Terkait | Effort (S/M/L) | Status | Bukti |
|---|---|---|---|---|---|---|---|---|
| T-001 | `<fitur>` | `<apa yang dites>` | Unit | P0 | R-001 | M | Planned | `<DCD-x / rest-api-doc>` |
| T-002 | `<fitur>` | `<apa yang dites>` | E2E | P0 | R-002 | L | Planned | `<project-overview journey>` |
| T-003 | `<fitur>` | `<apa yang dites>` | Integration | P1 | — | S | Deferred | `<rest-api-doc>` |

## 4. Ringkasan

```
Total item : N
P0         : x    ← wajib tuntas sebelum exit criteria terpenuhi
P1         : y
P2         : z
Deferred   : d
Effort     : S=a, M=b, L=c
```

## 5. Data & Environment

| Kebutuhan | Item Terkait | Keterangan |
|---|---|---|
| `<fixture/factory/mock service>` | `<T-00x>` | `<cara menyediakan>` |

## 6. Exit Criteria Plan Ini

- [ ] Semua item P0 berstatus `Passed`
- [ ] Semua item P1 berstatus `Passed` atau punya keputusan eksplisit user
- [ ] Semua `Deferred` tercatat dengan alasan
- [ ] Quality gates strategy yang relevan lulus

## 7. Hasil Eksekusi (diisi belakangan)

| ID | Status | Evidence | Catatan |
|---|---|---|---|
| T-001 | `Planned` | — | — |

## 8. Unknowns & Perlu Dikonfirmasi

| ID | Hal | Penanda | Tindak Lanjut |
|---|---|---|---|
| U-001 | `<misal: fitur belum ada use case di DCD>` | `Perlu Dikonfirmasi` | `<aksi>` |

---

## Operating Rules

- Jika strategy berubah, sinkronkan plan ini sebelum eksekusi lanjut.
- Status item: `Planned → In Progress → Passed / Failed / Blocked / Deferred`.
- `Failed` atau `Blocked` jangan disamarkan sebagai selesai.
- Update `.ai-doc/3p.md` setelah plan dibuat atau status signifikan berubah.
