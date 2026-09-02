# Step 2: Coverage Assessment

> Petakan tiap attack surface ke kategori OWASP Top 10 (2025) dan tetapkan status
> coverage-nya. Tidak ada kategori yang boleh hilang diam-diam.
> **HALT** — review matrix bersama user sebelum lanjut ke findings.

---

## Input dari User

- **Konfirmasi status per kategori** — terutama yang berpotensi jadi `Accepted Risk`
- **Output scan awal (opsional)** — bila user sudah punya hasil scan untuk dijadikan bukti `Covered`

---

## Yang Dilakukan Agent

### 1. Bangun Matrix Surface × Kategori

Untuk setiap attack surface (S-00x), tentukan kategori OWASP yang relevan:

| Kategori | Kapan relevan |
|---|---|
| A01 Broken Access Control | Hampir selalu — IDOR, function-level access, SSRF |
| A02 Security Misconfiguration | Selalu ada server/config |
| A03 Supply Chain Failures | Selalu ada dependency/build pipeline |
| A04 Cryptographic Failures | Ada data sensitif / TLS / hashing password |
| A05 Injection | Ada input user → interpreter (SQL/XSS/CSRF/command) |
| A06 Insecure Design | Rate limit, business logic abuse, URL-fetch tanpa allow-list |
| A07 Authentication Failures | Ada mekanisme auth apa pun |
| A08 Data Integrity Failures | CDN script, deserialization, CSP |
| A09 Logging & Alerting Failures | Selalu — log kejadian keamanan + alertnya |
| A10 Mishandling Exceptional Conditions | Selalu — error path, fail-open vs fail-closed |

Kategori LLM Top 10 (LLM01–LLM10) hanya dibuka jika `llm_features: true`.

### 2. Tetapkan Status Tiap Kategori

| Status | Syarat |
|---|---|
| `Covered` | Ada bukti test/kontrol nyata (test file, config terverifikasi, output scan) |
| `Planned` | Relevan, belum ada bukti — kandidat item test plan |
| `N/A` | Attack surface-nya benar-benar tidak ada — wajib alasan berbasis bukti |
| `Accepted Risk` | Sadar risikonya dan user menyetujui penerimaannya secara eksplisit |

Contoh baris matrix:

| Kategori | Surface Terkait | Status | Bukti / Alasan |
|---|---|---|---|
| A01 | S-001, S-004 | `Planned` | Endpoint publik tanpa test akses — bukti: rest-api-doc |
| A04 | S-002 | `Covered` | Password hashing argon2 — bukti: `auth/hash.ts` |
| A06 | S-005 (webhook) | `Planned` | URL fetch tanpa allow-list terlihat di kode (`Asumsi`) |
| A08 | — | `N/A` | Tidak ada CDN script/deserialization — bukti: struktur frontend |

### 3. Kumpulkan Temuan Statis Sambil Menilai

Catat samping temuan yang terlihat dari inspeksi statis (diproses penuh di Step 3):

- Dependency dengan versi rentan yang diketahui (tanpa menjalankan scanner)
- Debug mode / verbose error aktif di config produksi
- Default credential / secret di file repo
- Header keamanan tidak diset (HSTS, CSP, nosniff)

### 4. Review Matrix Bersama User

Tampilkan matrix lengkap, lalu soroti:

- Kategori yang usulannya `Accepted Risk` → minta keputusan eksplisit sekarang:

> Kategori **A0x** saya usulkan `Accepted Risk` karena [alasan].
> Ini keputusan risiko Anda, bukan saya. Setuju diterima risikonya?
> (Setuju / Tetap Planned / Batal)

- Klaim bertanda `Asumsi` → minta konfirmasi atau turunkan jadi `Perlu Dikonfirmasi`.
- Kategori `N/A` dengan bukti lemah → tanyakan ulang.

**HALT** — setelah matrix disetujui user.

---

## Output Step 2

```yaml
session_state:
  coverage_matrix:
    - category: "A01"
      surfaces: ["S-001", "S-004"]
      status: "planned" | "covered" | "na" | "accepted-risk"
      evidence_or_reason: "..."
      decided_by: "agent-proposal" | "user"
  static_findings_prelim:
    - description: "..."
      location: "..."
      draft_severity: "high|medium|low"
  llm_matrix_opened: false
```

---

## Aturan

- **`Accepted Risk` butuh persetujuan eksplisit user** — agent tidak berhak memutuskan.
- **`N/A` tanpa bukti = `Perlu Dikonfirmasi`** — jangan pakai N/A sebagai jalan pintas.
- **Semua 10 kategori harus muncul di matrix** — meski statusnya N/A.
- **Temuan statis hanya dicatat**, belum diputuskan — itu Step 3.
