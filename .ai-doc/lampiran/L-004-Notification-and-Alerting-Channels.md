# Lampiran L-004: Sistem Notifikasi & Integrasi Multi-Channel Alert

| Metadata | Nilai |
|---|---|
| **Kode Lampiran** | `L-004` |
| **Nama Lampiran** | Sistem Notifikasi & Integrasi Multi-Channel Alert |
| **Target Service** | `clamav-service` |
| **Status** | Approved |
| **Versi** | 1.0 |
| **Tanggal Pembuatan** | 2026-09-02 |
| **Dokumen Utama** | [`.ai-doc/project-overview.md`](file:///home/ubuntu/workspace/plan/clamav-service/.ai-doc/project-overview.md) |

---

## 1. Ringkasan
Dokumen ini merinci arsitektur subsistem notifikasi, format payload template pesan untuk Telegram, Discord, dan Email (SMTP), mekanisme aktivasi otomatis/toggle UI, serta algoritma penanganan lonjakan serangan (*Alert Flood Throttling*).

---

## 2. Saluran Notifikasi yang Didukung

| Saluran | Parameter Kunci | Format Payload |
|---|---|---|
| **Telegram Bot** | `bot_token`, `chat_id` | MarkdownV2 / HTML Text Message |
| **Discord Webhook** | `webhook_url` | JSON Rich Embed dengan field terstruktur |
| **Email (SMTP)** | `smtp_host`, `port`, `user`, `pass`, `recipients` | Responsif HTML Template & Plaintext Fallback |

---

## 3. Template & Format Pesan

### 3.1. Template Pesan Telegram
```text
🚨 [CLAMAV ALERT] Malware Detected!
──────────────────────────────────
🦠 Threat : Win.Trojan.Agent-9482
📁 File   : invoice_payment.pdf.exe (2.4 MB)
🔑 Hash   : e3b0c44298fc1c149afbf4c8996fb92427ae...
🛡️ Vault  : Q-20260902-9812 (Quarantined)
🏷️ Source : Customer-Portal-API (10.0.4.12)
⏰ Time   : 2026-09-02 16:30:15 WIB

👉 [Buka File di Karantina](https://clamav.internal/admin/quarantine/Q-20260902-9812)
```

### 3.2. Template Pesan Discord (Rich Embed)
```json
{
  "username": "ClamAV Security Guard",
  "avatar_url": "https://clamav.internal/assets/shield-red.png",
  "embeds": [
    {
      "title": "🚨 Malware Threat Detected!",
      "color": 15158332,
      "fields": [
        { "name": "Threat Name", "value": "`Win.Trojan.Agent-9482`", "inline": true },
        { "name": "File Name", "value": "`invoice_payment.pdf.exe`", "inline": true },
        { "name": "Status", "value": "🛡️ **QUARANTINED**", "inline": true },
        { "name": "File SHA256", "value": "`e3b0c44298fc...`", "inline": false },
        { "name": "Source Consumer", "value": "Customer-Portal-API", "inline": true },
        { "name": "Vault ID", "value": "`Q-20260902-9812`", "inline": true }
      ],
      "footer": { "text": "ClamAV Security Telemetry • 2026-09-02 16:30:15 UTC" }
    }
  ]
}
```

---

## 4. Mekanisme Aktivasi & Kontrol Pengguna

1. **Auto-Detection via Environment Variable:**
   Jika variabel `ALERT_TELEGRAM_BOT_TOKEN` dan `ALERT_TELEGRAM_CHAT_ID` ada di `.env`, channel Telegram otomatis berstatus **ENABLED** saat startup.
2. **Interactive UI Toggle:**
   Pada Web Admin UI (Menu *Settings > Alerts*), admin dapat mematikan (*mute*) sementara notifikasi melalui tombol toggle switch `[ON / OFF]` tanpa menghapus kredensial token yang tersimpan.
3. **Penyaringan Event (Event Filtering):**
   * `[x] File Terinfeksi (INFECTED)` — *Default: ON*
   * `[x] File Mencurigakan / Terenkripsi (SUSPICIOUS)` — *Default: ON*
   * `[ ] Status Reload Database Signature` — *Default: OFF*

---

## 5. Proteksi Anti-Spam (*Alert Flood Throttling*)

Untuk mencegah *alert fatigue* atau spamming ratusan pesan saat terjadi serangan upload massal:

* **Threshold Trigger:** Jika terjadi lebih dari 5 deteksi malware dalam jendela waktu 60 detik.
* **Mekanisme:**
  1. Notifikasi ke-1 sampai ke-5 dikirimkan secara instan.
  2. Notifikasi ke-6 dan seterusnya ditahan (*buffered*) di memori.
  3. Setelah 60 detik, sistem mengirimkan **1 Pesan Rekapitulasi (*Batch Digest*)**:
     ```text
     ⚠️ [ALERT DIGEST] Lonjakan Deteksi Malware Massal!
     ──────────────────────────────────────────────────
     Sebanyak 48 ancaman malware tambahan terdeteksi dan 
     telah otomatis dikarantina dalam 1 menit terakhir.
     
     👉 [Buka Audit Log Dashboard](https://clamav.internal/admin/audit)
     ```
