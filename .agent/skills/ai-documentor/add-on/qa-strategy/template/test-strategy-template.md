# Test Strategy — <Nama Project>

> Strategi testing berbasis bukti dari artefak `.ai-doc/` yang ada.
> Dibuat oleh add-on qa-strategy. Setiap klaim menunjuk sumber buktinya.

## 1. Metadata

- **Project:** `<nama project>`
- **Cakupan:** `System` | `Component: <nama>` | `Release: <nama>`
- **Created at:** `<YYYY-MM-DD>`
- **Sumber bukti:** `<daftar artefak .ai-doc yang dipakai>`
- **Last updated:** `<YYYY-MM-DD>`

## 2. Konteks & Tujuan

- **Apa yang diuji:** `<ringkasan sistem/komponen dari project-overview>`
- **Critical journeys:** `<daftar journey yang jika rusak berdampak besar>`
  - Sumber: `<artefak + bagian>`
- **Tujuan strategi:** `<1–3 kalimat, misal: menahan regresi pada journey kritikal dengan biaya maintenance test yang rendah>`

## 3. Risk Matrix

> Skor = Impact × Likelihood (masing-masing 1–5). Sumber bukti wajib;
> baris tanpa bukti ditandai `Asumsi`.

| ID | Risiko | Komponen/Area | Impact | Likelihood | Skor | Sumber Bukti |
|---|---|---|---|---|---|---|
| R-001 | `<deskripsi risiko>` | `<komponen>` | `<1–5>` | `<1–5>` | `<I×L>` | `<artefak + bagian>` |
| R-002 | `<deskripsi risiko>` | `<komponen>` | `<1–5>` | `<1–5>` | `<I×L>` | `Asumsi — Perlu Dikonfirmasi` |

**Prioritas pengujian:** urut skor tertinggi; risiko `Asumsi` divalidasi dulu
sebelum jadi dasar keputusan besar.

## 4. Test Pyramid

| Level | Proporsi | Cakupan Utama | Alasan (dari bukti) |
|---|---|---|---|
| Unit | `<%>` | `<logika apa saja>` | `<alasan dari karakter sistem>` |
| Integration | `<%>` | `<batas antar modul, DB, API pihak ketiga>` | `<alasan>` |
| E2E | `<%>` | `<journey kritikal lewat UI/API publik>` | `<alasan>` |

Catatan karakter sistem: `<misal API-heavy / UI-heavy / background-job-heavy>`.

## 5. Quality Gates

| Gate | Lokasi | Ambang / Syarat | Aksi Jika Gagal |
|---|---|---|---|
| `<unit + integration>` | PR CI | `<semua lulus; coverage ≥ X%>` | Block merge |
| `<smoke E2E>` | Pre-deploy staging | `<journey kritikal lulus>` | Block deploy |

Semua ambang objektif dan dapat dicek mesin.

## 6. Entry / Exit Criteria per Level

### Unit

- **Entry:** `<misal: modul dapat di-import terpisah>`
- **Exit:** `<misal: semua case P0/P1 lulus; tidak ada skip tanpa alasan tercatat>`

### Integration

- **Entry:** `<misal: environment test + data factory tersedia>`
- **Exit:** `<misal: kontrak endpoint publik terverifikasi>`

### E2E

- **Entry:** `<misal: staging stabil dan seed data siap>`
- **Exit:** `<misal: semua critical journey lulus di browser target>`

## 7. Tooling (Opsional)

> Section ini hanya diisi jika diminta saat penyusunan strategy.

| Kebutuhan | Pilihan Default | Alasan | Alternatif |
|---|---|---|---|
| `<E2E>` | `<tool>` | `<alasan terkait stack>` | `<alternatif — kondisi pemakaian>` |

## 8. Out of Scope

`<area yang sengaja tidak dicover strategi ini, beserta alasannya>`

## 9. Unknowns & Perlu Dikonfirmasi

| ID | Hal | Penanda | Tindak Lanjut |
|---|---|---|---|
| U-001 | `<area tanpa bukti>` | `Perlu Dikonfirmasi` | `<aksi: lengkapi artefak / konfirmasi user>` |

---

## Operating Rules

- Strategy adalah sumber prioritas untuk semua test plan turunannya.
- Perubahan risk matrix harus disinkronkan ke plan yang sudah ada.
- Update file ini dan `.ai-doc/3p.md` setelah revisi bermakna.
