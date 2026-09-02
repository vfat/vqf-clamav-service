# C4 Component Diagrams — `clamav-service`

## 1. All-in-One Container Component Diagram

### Deskripsi
Diagram komponen tingkat container untuk `clamav-service` yang dieksekusi di dalam satu Docker Container mandiri berbasis Alpine Linux. Diagram ini menggambarkan hubungan antar komponen internal Go backend, antarmuka Web UI, penyimpanan SQLite & Vault karantina, komunikasi Unix Domain Socket ke daemon ClamAV (`clamd`), serta integrasi alert eksternal.

### Diagram
```mermaid
C4Component
    title Component Diagram for clamav-service All-in-One Container

    Container(consumerApp, "Consumer Applications", "Microservices / Web Apps / Upload Gateways")
    Container(secOpsAdmin, "Security Admin Browser", "Web Browser")

    ContainerDb(sqliteDb, "SQLite Database (WAL)", "/data/clamav-service.db")
    ContainerDb(quarantineVault, "Quarantine Storage", "/data/quarantine/")

    System_Ext(telegramApi, "Telegram API", "Bot API HTTPS")
    System_Ext(discordWebhook, "Discord API", "Webhook HTTPS")
    System_Ext(smtpServer, "SMTP Mail Server", "TLS SMTP")
    System_Ext(clamavCvdNet, "Cisco Talos / ClamAV Network", "database.clamav.net")

    Container_Boundary(clamavContainer, "clamav-service Container") {
        Component(processSupervisor, "Process Supervisor (PID 1)", "Go Runtime / OS Signal Handler")
        Component(clamdDaemon, "clamd Daemon", "C ClamAV Engine (In-Memory Signatures)")
        Component(freshclamDaemon, "freshclam Daemon", "C Background Signature Updater")

        Component(adminWebUI, "Embedded Admin UI", "SPA (HTML5 / JS / Tailwind)")
        Component(apiGateway, "HTTP REST API Gateway", "Go net/http Router & Middleware")
        Component(authRateLimiter, "Auth & Token Bucket Limiter", "Go In-Memory Sliding Window")
        Component(scanController, "Scan & Stream Controller", "Go Stream Handler & Multipart Parser")
        Component(clamdBridge, "ClamAV Socket Bridge", "Go Unix Domain Socket Client")
        Component(quarantineManager, "Quarantine Vault Manager", "Go Storage Driver & De-scrambler")
        Component(sqliteManager, "SQLite Storage Manager", "Go modernc.org/sqlite / WAL")
        Component(alertDispatcher, "Alert & Flood Throttler", "Go Multi-Channel Notifier")

        Rel(processSupervisor, clamdDaemon, "Spawns and monitors", "os/exec")
        Rel(processSupervisor, freshclamDaemon, "Spawns and monitors", "os/exec")
        Rel(freshclamDaemon, clamavCvdNet, "Downloads CVD updates", "HTTPS:443")
        Rel(freshclamDaemon, clamdDaemon, "Sends RELOAD command", "Unix Socket")

        Rel(apiGateway, authRateLimiter, "Validates key & checks quota")
        Rel(apiGateway, adminWebUI, "Serves static assets & UI API")
        Rel(apiGateway, scanController, "Routes scan payloads")
        Rel(apiGateway, quarantineManager, "Routes quarantine & restore requests")

        Rel(scanController, clamdBridge, "Pipes binary stream", "Chunked INSTREAM")
        Rel(clamdBridge, clamdDaemon, "Streams payload for scan", "Unix Socket (/tmp/clamd.sock)")

        Rel(scanController, quarantineManager, "Dispatches infected payload")
        Rel(quarantineManager, quarantineVault, "Writes scrambled .quarantine file", "FS 0600")

        Rel(scanController, sqliteManager, "Writes scan audit log", "SQL Insert")
        Rel(quarantineManager, sqliteManager, "Manages quarantine metadata", "SQL CRUD")
        Rel(authRateLimiter, sqliteManager, "Loads API keys & policies", "SQL Read")

        Rel(scanController, alertDispatcher, "Triggers threat alert on infected")
        Rel(alertDispatcher, telegramApi, "Sends Markdown alert", "HTTPS POST")
        Rel(alertDispatcher, discordWebhook, "Sends Rich Embed alert", "HTTPS POST")
        Rel(alertDispatcher, smtpServer, "Sends HTML email alert", "SMTP TLS")
    }

    Rel(consumerApp, apiGateway, "Uploads file / binary stream / async job", "HTTP REST:8080")
    Rel(secOpsAdmin, apiGateway, "Accesses Admin Dashboard & Quarantine Vault", "HTTP:8080")
    Rel(sqliteManager, sqliteDb, "Reads and writes transactional data", "Local File I/O")
```

---

### Komponen Internal

| Komponen | Deskripsi | Teknologi |
|---|---|---|
| **Process Supervisor (PID 1)** | Master entrypoint container yang men-spawn, memonitor, dan me-restart daemon `clamd` & `freshclam`, serta menangani *graceful shutdown* sinyal OS (`SIGTERM`/`SIGINT`). | Go `os/exec`, `os/signal` |
| **HTTP REST API Gateway** | Router HTTP utama yang melayani endpoint pemindaian, manajemen karantina, health probe, dan menyajikan asset Embedded Web UI. | Go `net/http` / Chi Router |
| **Auth & Token Bucket Limiter** | Verifikasi API Key (salted hash) dan pembatasan laju request *in-memory* (Global RPS + Granular per-key). | Go `crypto/subtle`, Token Bucket |
| **Scan & Stream Controller** | Orkestrator alur pemindaian: menangani multipart upload, raw binary stream, parsing arsip terenkripsi password, dan async job queue. | Go Goroutines, `io.Reader/Writer` |
| **ClamAV Socket Bridge** | Client Unix Domain Socket yang menghubungkan API ke daemon `clamd` menggunakan protokol `INSTREAM` non-blocking. | Go `net.DialUnix` |
| **Quarantine Vault Manager** | Mengelola penyimpanan aman file terinfeksi (scrambling, penamaan `.quarantine`, izin `0600`), retensi TTL 7 hari, dan mekanisme pemulihan (*restore*). | Go File Driver, SHA-256 |
| **SQLite Storage Manager** | Layer persistensi database WAL mode untuk mencatat log audit pemindaian, metadata karantina, API keys, dan konfigurasi dinamis. | Go SQLite (WAL Mode), AES-256-GCM |
| **Alert & Flood Throttler** | Dispatcher notifikasi multi-channel (Telegram, Discord, Email) yang dilengkapi proteksi *Anti-Spam Throttling* (Batch Digest jika >5 malware/menit). | Go HTTP Client, SMTP Driver |
| **Embedded Admin UI** | Antarmuka web mandiri (SPA) untuk monitoring status daemon, drag-and-drop test scan, inspeksi karantina, dan export audit log. | HTML5, Vanilla JS, Tailwind CSS |

---

### External Systems & Subprocesses

| System / Subprocess | Tipe | Fungsi |
|---|---|---|
| **`clamd` Daemon** | Local Daemon Process | Engine antivirus inti ClamAV yang memuat ~8.5 juta signature ke RAM dan melayani pemindaian via socket `/tmp/clamd.sock`. |
| **`freshclam` Daemon** | Background Daemon Process | Proses terjadwal yang otomatis mengunduh database virus baru dari `database.clamav.net` setiap 2–4 jam. |
| **Telegram API** | External HTTPS Service | Endpoint Bot API untuk mengirimkan notifikasi instan ancaman malware ke grup/channel Telegram SecOps. |
| **Discord Webhook** | External HTTPS Service | Endpoint Webhook Discord untuk mengirimkan notifikasi Rich Embed ke channel monitoring. |
| **SMTP Mail Server** | External Mail Relay | Server surat untuk mengirimkan laporan email keamanan resmi ke daftar penerima administrator. |
| **Cisco Talos Network** | External Update Server | Server repositori resmi database signature virus ClamAV (`main.cvd`, `daily.cvd`, `bytecode.cvd`). |

---

### Relasi Utama

**Internal Flows:**
- `apiGateway` $\rightarrow$ `scanController`: Meneruskan request pemindaian file atau stream biner.
- `scanController` $\rightarrow$ `clamdBridge`: Meneruskan chunk data biner via protokol `INSTREAM`.
- `clamdBridge` $\rightarrow$ `clamdDaemon`: Mengirim payload ke Unix Socket `/tmp/clamd.sock` dan menerima verdict (`OK` / `FOUND`).
- `scanController` $\rightarrow$ `quarantineManager`: Memindahkan file terinfeksi ke Vault jika verdict `INFECTED`.
- `scanController` $\rightarrow$ `alertDispatcher`: Memicu pengiriman alert multi-channel saat malware terdeteksi.
- `quarantineManager` $\rightarrow$ `sqliteManager`: Mencatat status `QUARANTINED` / `RESTORED` dan hash whitelist.

**External Integrations:**
- `consumerApp` $\rightarrow$ `apiGateway`: Mengirimkan request pemindaian via REST API (Port 8080).
- `freshclamDaemon` $\rightarrow$ `clamavCvdNet`: Mengunduh delta update signature virus secara berkala.
- `alertDispatcher` $\rightarrow$ `telegramApi` / `discordWebhook` / `smtpServer`: Mengirimkan notifikasi alert keamanan.
