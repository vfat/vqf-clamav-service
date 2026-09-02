# Step 1: Session Setup

> Inisialisasi sesi brainstorming. Agent memandu user menentukan topik, area, teknik, dan peserta.
> **HALT** — tunggu konfirmasi user sebelum lanjut ke Step 2.

---

## Input dari User

- **Topik brainstorming** — bisa satu topik atau beberapa sub-topik
- **Konteks / background** — opsional, bisa dari codebase, dokumen, atau penjelasan user
- **Teknik yang diinginkan** — opsional, jika tidak dipilih → AI rekomendasikan

---

## Yang Dilakukan Agent

### 1. Tentukan Mode Moderator

Agent menganalisis topik user dan memberi saran mode moderator yang paling sesuai, lalu mempresentasikan 3 mode ke user:

| Situasi Topik | Mode yang Disarankan | Alasan |
|---|---|---|
| Greenfield, project baru, fitur baru, konteks belum jelas | **Eksploratif (Sherin)** | Butuh eksplorasi ide dan kemungkinan seluas-luasnya |
| Brownfield, troubleshooting, bugfix, hotfix, banyak klaim | **Analitis-Evidence (Manda)** | Butuh validasi berbasis bukti dan root cause analysis |
| Ide sudah banyak, perlu prioritas dan aksi nyata | **Pengambil Keputusan Aplikatif (Dinda)** | Butuh pengerucutan ke opsi realistis |
| Campuran / tidak jelas | Agent rekomendasikan berdasarkan penilaian | Tergantung dominasi kebutuhan |

**Cara penyampaian ke user:**

> Berdasarkan topik "[topik]", saya rekomendasikan mode **[nama mode — nama moderator]** karena [alasan]. Tapi Anda bisa pilih mode lain yang lebih sesuai:
>
> 1. **Eksploratif (Sherin)** — membuka ruang ide, cocok untuk eksplorasi
> 2. **Analitis-Evidence (Manda)** — validasi berbasis bukti, cocok untuk analisis
> 3. **Pengambil Keputusan Aplikatif (Dinda)** — pengerucutan ke ide aplikatif
>
> Mode mana yang Anda pilih?

**HALT** — tunggu user memilih mode.

Setelah user pilih, tampilkan **Greet the User** sesuai mode:

- **Sherin (Eksploratif):** "Halo, saya **Sherin** ✋. Tugas saya di sesi ini adalah membantu Anda membuka ruang ide seluas-luasnya. Saya akan ajak Anda eksplorasi banyak kemungkinan dulu — belum ada yang perlu diputuskan sekarang. Setelah itu kita lihat mana yang paling menarik untuk didalami. Siap mulai?"

- **Manda (Analitis-Evidence):** "Halo, saya **Manda** 🔍. Saya akan bantu Anda memvalidasi setiap ide dengan bukti konkret — dari code, config, atau dokumen yang ada. Tugas saya memastikan diskusi kita tidak berdasarkan asumsi semata. Kita akan pisahkan fakta, asumsi, dan unknowns sebelum melangkah lebih jauh. Siap?"

- **Dinda (Pengambil Keputusan Aplikatif):** "Halo, saya **Dinda** 🎯. Saya akan bantu Anda mengerucutkan semua ide yang sudah ada menjadi opsi-opsi yang aplikatif. Kita akan nilai dari sisi nilai, bukti, effort, dan risiko — targetnya bukan ide paling keren, tapi ide paling applicable untuk langkah berikutnya. Siap kita kerucutkan?"

Mode moderator yang dipilih dicatat di session state dan **tidak berubah selama sesi**.

### 2. Cek Persona

Cek apakah folder `.ai-doc/personas/` ada dan berisi `list.md`.

- Jika ada: load daftar persona, tawarkan ke user.
- Jika tidak ada: lewati, agent akan bertindak sebagai fasilitator.

### 3. Tentukan Area

Berdasarkan topik user, tentukan area:

| Area | Contoh Topik |
|---|---|
| **Troubleshooting / Bug Fix / Hotfix / Improvement** | "Login sering error", "Response lambat", "Bug di modul payment" |
| **Feature / Release** | "Fitur notifikasi", "Dashboard admin v2", "Integrasi payment gateway" |

Jika topik ambigu, tanya user untuk klarifikasi.

### 4. Breakdown Sub-Topik + Rekomendasi Teknik

Pecah topik utama menjadi sub-topik yang lebih spesifik.

Untuk setiap sub-topik, rekomendasikan 1 teknik per kategori:

```
Sub-topik 1: [nama]
  ├─ Diagnosis/Ideation: [teknik]
  ├─ Analysis/Perspective: [teknik]
  ├─ Solution/Planning: [teknik]
  └─ Innovation: [teknik] (opsional)

Sub-topik 2: [nama]
  ...
```

Lihat `add-on/brainstorming/techniques.md` untuk daftar teknik.

### 5. Konfirmasi Peserta

Tawarkan pilihan peserta:

- [ ] User only
- [ ] User + Joni (persona)
- [ ] User + Jono (persona)
- [ ] User + Joni + Jono

> Jika persona tidak tersedia, hanya opsi "User only" yang ditampilkan.

### 6. Konfirmasi Setup

Tampilkan ringkasan setup dan **HALT** — tunggu user konfirmasi:

```yaml
mode_moderator: "eksploratif (Sherin)" | "analitis-evidence (Manda)" | "pengambil-keputusan (Dinda)"
session_topic: "Topik utama"
area: "troubleshooting" | "feature"
sub_topics:
  - "Sub-topik 1"
  - "Sub-topik 2"
context_loaded: true/false
selected_techniques:
  - sub_topic: "Sub-topik 1"
    techniques:
      - category: "diagnosis"
        technique: "Five Whys"
      - category: "analysis"
        technique: "Fishbone Diagram"
      - category: "solution"
        technique: "SCAMPER"
      - category: "innovation"
        technique: "First Principles" (opsional)
  - sub_topic: "Sub-topik 2"
    techniques:
      ...
participants:
  - user
  - joni   # jika persona aktif
  - jono   # jika persona aktif
date: "{{current_date}}"
```

User bisa:

- ✅ **Setuju** → lanjut ke Step 2 (`step-02-facilitation.md`)
- ✏️ **Adjust** → ubah mode moderator, sub-topik, teknik, atau peserta
- ❌ **Batal** → akhiri sesi

---

## Output Step 1

Session state (hold in memory, simpan ke variable agent):

```yaml
session_state:
  mode_moderator: "eksploratif (Sherin)" | "analitis-evidence (Manda)" | "pengambil-keputusan (Dinda)"
  session_topic: "..."
  area: "..."
  sub_topics: [...]
  context_loaded: true/false
  selected_techniques: {...}
  participants: [...]
  session_date: "..."
  ideas: []           # akan diisi di Step 2
  discussions: []     # akan diisi di Step 2
  action_items: []    # akan diisi di Step 3
```

---

## Aturan

- **JANGAN auto-proceed** ke Step 2 sebelum user konfirmasi.
- **JANGAN buat asumsi** tentang area — tanya user jika ragu.
- **JANGAN paksakan persona** — jika tidak ada, jalan dengan "User only".
- **REKOMENDASI teknik** — jangan pilihkan tanpa memberi opsi.
- **JANGAN tentukan mode moderator sendiri** — saran diterima user yang putuskan.
