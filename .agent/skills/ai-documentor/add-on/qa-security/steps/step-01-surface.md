# Step 1: Surface Mapping

> Petakan attack surface dari bukti — artefak `.ai-doc/` jika ada,
> rekonstruksi codebase jika tidak. Dirancang tetap berguna untuk brownfield.
> **HALT** — konfirmasi scope & threat model ringkas sebelum Step 2.

---

## Input dari User

- **Scope audit** — seluruh sistem / satu komponen / satu endpoint group (wajib)
- **Auth mechanism** — session/JWT/OAuth/API key (boleh dijawab, boleh dibiarkan dideteksi)
- **Data sensitif** — PII, pembayaran, kredensial, atau lainnya
- **Deployment model** — cloud/container/serverless
- **Compliance context** — opsional; hanya konteks, bukan cakupan add-on ini

---

## Yang Dilakukan Agent

### 1. Baca Control Plane (jika ada)

- `.ai-doc/3p.md` — status project
- `.ai-doc/constitution.md` — aturan lokal

Jika tidak ada: lanjut saja, catat `bootstrap: none` di session state.

### 2. Kumpulkan Bukti — Artefak Dulu, Codebase Kemudian

| Prioritas | Sumber | Yang diambil | Jika tidak ada |
|---|---|---|---|
| 1 | `rest-api-doc` | Endpoint publik + auth per endpoint | Rekonstruksi dari route |
| 1 | `Dokumentasi-Codebase.md` / DCD | Komponen, integrasi pihak ketiga | Rekonstruksi ringan |
| 2 | Route/handler di kode | Daftar endpoint nyata, middleware auth | — |
| 2 | Config server/app | Security headers, CORS, debug flag, default credential | — |
| 3 | Dependency manifest + lockfile | Daftar dependency versi (bahan A03) | Tandai `Perlu Dikonfirmasi` |

Aturan brownfield:

- Rekonstruksi dari kode adalah **inferensi**, bukan dokumentasi final.
- Setiap klaim hasil inferensi ditandai `Asumsi` di session state.
- Jangan menulis/memperbarui artefak dokumentasi lain di step ini — itu tugas core flow.

### 3. Identifikasi Karakter Keamanan

Deteksi dan konfirmasikan ke user:

```yaml
security_profile:
  auth_mechanism: "jwt" | "session" | "oauth" | "api-key" | "<Perlu Dikonfirmasi>"
  sensitive_data: ["<PII>", "<payment>", "..."]
  deployment: "aws" | "gcp" | "container" | "serverless" | "..."
  api_style: "rest" | "graphql" | "grpc" | "mixed"
  llm_features: true/false          # menentukan apakah OWASP LLM Top 10 dinilai
  compliance_context: []            # hanya dicatat, tidak dipetakan di sini
```

### 4. Susun Daftar Attack Surface

| ID | Surface | Jenis | Auth? | Data Sensitif | Sumber Bukti |
|---|---|---|---|---|---|
| S-001 | `<endpoint/komponen>` | API publik / admin panel / webhook / file upload | Ya/Tidak | `<jenis>` | `<artefak/route file>` |

Surface yang wajib dicari (sering terlewat):

- Webhook receiver & fitur yang menerima URL user (vektor SSRF)
- File upload/import
- Admin panel & endpoint internal
- Job/background worker dengan input dari queue

### 5. Konfirmasi Scope & Threat Model Ringkas

> Dari [artefak/rekonstruksi] saya menemukan **[N] attack surface**:
> [ringkasan]. Profil keamanan terdeteksi: [auth], data sensitif: [jenis].
>
> Usulan scope audit: **[usulan]** karena [alasan].
> - ✅ **Lanjut dengan scope ini**
> - ✏️ **Ubah scope** — fokus ke surface tertentu
> - ➕ **Lengkapi bukti dulu** — jalankan core flow untuk artefak yang hilang

**HALT** — tunggu keputusan user.

---

## Output Step 1

```yaml
session_state:
  bootstrap: "ai-doc" | "codebase-reconstruction"
  scope: "system" | "component" | "surface-group"
  security_profile:
    auth_mechanism: "..."
    sensitive_data: [...]
    deployment: "..."
    api_style: "..."
    llm_features: false
  attack_surfaces:
    - id: "S-001"
      name: "..."
      kind: "..."
      auth: true/false
      sensitive_data: [...]
      source: "..."        # artefak atau file kode
      inference: true/false  # true = hasil rekonstruksi → klaim = Asumsi
```

---

## Aturan

- **JANGAN memodifikasi artefak dokumentasi lain** — add-on ini hanya membaca;
  penulisan hanya pada output sendiri di Step 3.
- **Klaim hasil rekonstruksi = `Asumsi`** sampai dikonfirmasi user.
- **JANGAN menebak auth mechanism** — kalau tidak ketemu di config/kode,
  tandai `Perlu Dikonfirmasi` dan tanya.
- **REKOMENDASIKAN scope** — saran dengan alasan, keputusan milik user.
