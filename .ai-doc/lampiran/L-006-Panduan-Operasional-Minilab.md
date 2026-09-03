# 📖 Lampiran L-006: Panduan Operasional ClamAV Service (Minilab)

Panduan praktis untuk menjalankan (*start*), menghentikan (*stop*), mengelola firewall, dan memelihara container **`vqf-clamav-service`** di lingkungan server VPS.

---

## 🚀 1. Cara Menjalankan (START)

### Langkah 1: Buka Akses Port di Firewall (UFW)
Jika dashboard ingin diakses melalui browser dari jaringan publik/laptop Anda:
```bash
sudo ufw allow 8085/tcp
```

### Langkah 2: Jalankan Container
Masuk ke direktori minilab dan jalankan Docker Compose:
```bash
cd /home/ubuntu/workspace/minilab/clamav-service
docker compose up -d
```

### Langkah 3: Verifikasi Status
Pastikan container berstatus **Up (healthy)**:
```bash
docker compose ps
```
*Catatan:* Pada awal start, ClamAV membutuhkan ~15–20 detik untuk memuat 8+ juta database virus signature ke RAM sebelum statusnya berubah menjadi `healthy`.

### 🌐 Akses Web Dashboard
* **URL:** `http://<IP-VPS-ANDA>:8085/`
* **Default UI Password:** `123456`

---

## 🛑 2. Cara Menghentikan (STOP)

### Langkah 1: Matikan Container
Hentikan dan bersihkan container serta network-nya:
```bash
cd /home/ubuntu/workspace/minilab/clamav-service
docker compose down
```
*(Data karantina dan database di `./data` serta virus signature tetap aman dan tidak akan hilang).*

### Langkah 2: Tutup Kembali Port di Firewall
Tutup kembali port `8085` agar server tidak terbuka ke publik saat tidak digunakan:
```bash
sudo ufw delete allow 8085/tcp
```

### Langkah 3: Verifikasi Port Telah Tertutup
```bash
sudo ufw status | grep 8085
# Jika output kosong, port 8085 sudah tertutup rapat.
```

---

## 🛠️ 3. Operasi Harian & Maintenance (Cheat Sheet)

### Melihat Live Logs Daemon & Scanning
Untuk memantau aktivitas pemindaian file secara real-time:
```bash
docker compose logs -f
```

### Restart Service Saja (Tanpa Hapus Network)
```bash
docker compose restart
```

### Update ke Versi Image Terbaru dari Docker Hub
Jika ada rilis image baru di Docker Hub (`vickyfatrian/vqf-clamav-service:latest`):
```bash
docker compose pull
docker compose up -d
```

### Cek Endpoint Health via Terminal VPS
```bash
curl -s http://localhost:8085/api/v1/health | jq
```

---

## 📂 4. Lokasi Penyimpanan Data Persisten

Semua data penting disimpan di luar container agar aman saat container di-recreate:
1. **Quarantine Vault & SQLite DB:**  
   Tersimpan di host pada folder `./data/` (`/home/ubuntu/workspace/minilab/clamav-service/data/`).
2. **Virus Signatures (`main.cvd`, `daily.cvd`):**  
   Tersimpan di Docker Named Volume `clamav_signatures_data` agar saat restart server tidak perlu download ulang 300+ MB signature.

### 💾 Backup Data Karantina & Audit Log:
```bash
tar -czvf backup-clamav-data-$(date +%F).tar.gz ./data
```
