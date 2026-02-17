# Deployment Guide - Amaliah Ramadhan untuk Armbian/ARM

## 📦 Cara Build ARM Binary

### Prasyarat
- Go 1.21 atau lebih tinggi
- Git

### Build Single Binary untuk ARM

```bash
# Clone repository
git clone <repo-url>
cd apps_monitoring_ramadhan

# Build untuk ARM64 (64-bit)
./build-arm.sh

# Build untuk Linux Server (AMD64)
./build-linux.sh

# Untuk ARM 32-bit, edit build-arm.sh:
# Ubah GOARCH="arm64" menjadi GOARCH="arm"
```

### Output Build

Setelah build selesai, Anda akan mendapatkan:

```
dist/
├── amaliah-ramadhan-installer-linux-arm64  # Single binary installer (ARM)
├── amaliah-ramadhan-installer-linux-amd64  # Single binary installer (Linux/AMD64)
├── amaliah-ramadhan-1.0.0-linux-arm64.tar.gz  # Package lengkap (ARM)
├── amaliah-ramadhan-1.0.0-linux-amd64.tar.gz  # Package lengkap (Linux/AMD64)
└── README.md  # Dokumentasi deployment
```

---

## 🚀 Instalasi di Server Armbian

### Metode 1: Installer Wizard (Recommended)

1. **Transfer package ke server:**
   ```bash
   scp dist/amaliah-ramadhan-*.tar.gz user@192.168.1.100:/tmp/
   ```

2. **SSH ke server:**
   ```bash
   ssh user@192.168.1.100
   ```

3. **Ekstrak dan Jalankan Installer:**
   ```bash
   cd /tmp
   # Ekstrak package (sangat penting agar folder 'web' terbawa)
   tar -xzf amaliah-ramadhan-*.tar.gz
   
   # Masuk ke folder hasil ekstrak
   cd amaliah-ramadhan
   
   # Beri izin eksekusi
   chmod +x amaliah-ramadhan-installer-linux-*
   
   # Jalankan installer
   sudo ./amaliah-ramadhan-installer-linux-* -install
   ```

4. **Wizard akan menampilkan menu:**
   ```
   ╔══════════════════════════════════════════════════════════╗
   ║        🌙 AMALIAH RAMADHAN - SMK NIBA INSTALLER 🌙      ║
   ╚══════════════════════════════════════════════════════════╝
   
   1. 🆕 Install Baru
   2. 🔄 Update/Upgrade
   3. 🗑️  Uninstall
   4. ℹ️  Informasi Status
   5. ❌ Keluar
   
   Pilih (1-5):
   ```

5. **Pilih opsi 1 untuk Install Baru:**
   - Wizard akan menanyakan port (default: 8080)
   - Konfirmasi instalasi
   - Proses instalasi otomatis berjalan:
     ```
     ⏳ Membuat direktori instalasi ✅
     ⏳ Menyalin aplikasi ke sistem ✅
     ⏳ Menyalin file statis (templates, css, js) ✅
     ⏳ Membuat file konfigurasi ✅
     ⏳ Inisialisasi database ✅
     ⏳ Membuat systemd service ✅
     ⏳ Mengatur hak akses file ✅
     ⏳ Mengaktifkan dan memulai service ✅
     ```

6. **Selesai!** Aplikasi sudah berjalan di `http://server-ip:8080`

### Metode 2: Manual Installation

Lihat file `dist/README.md` untuk instruksi manual.

---

## 🔄 Update Aplikasi

### Menggunakan Installer:

1. **Transfer installer versi baru ke server:**
   ```bash
   scp dist/amaliah-ramadhan-installer-linux-arm64 user@server:/tmp/
   ```

2. **Jalankan installer:**
   ```bash
   sudo /tmp/amaliah-ramadhan-installer-linux-arm64 -install
   ```

3. **Pilih opsi 2 (Update/Upgrade)**

4. **Proses update otomatis:**
   ```
   ⏳ Menghentikan service ✅
   ⏳ Backup database ✅
   ⏳ Update aplikasi ✅
   ⏳ Update file statis ✅
   ⏳ Memulai service kembali ✅
   ```

**Note:** Database akan di-backup otomatis sebelum update!

---

## 📊 Cek Status Aplikasi

### Menggunakan Installer:

```bash
sudo /tmp/amaliah-ramadhan-installer-linux-arm64 -install
# Pilih opsi 4 (Informasi Status)
```

Output:
```
╔══════════════════════════════════════════════════════════╗
║                  STATUS APLIKASI                         ║
╚══════════════════════════════════════════════════════════╝

📦 Instalasi      : ✅ Terinstall
📁 Direktori      : /opt/amaliah-ramadhan
🟢 Service Status : ✅ Berjalan
🌐 Port           : 8080
🔗 URL            : http://localhost:8080
💾 Database       : ✅ (1024 KB)
```

### Menggunakan Systemctl:

```bash
# Status service
sudo systemctl status amaliah-ramadhan

# Lihat log real-time
sudo journalctl -u amaliah-ramadhan -f

# Restart service
sudo systemctl restart amaliah-ramadhan
```

---

## 🗑️ Uninstall

### Menggunakan Installer:

```bash
sudo /opt/amaliah-ramadhan/amaliah-ramadhan -install
# Pilih opsi 3 (Uninstall)
```

Proses uninstall akan:
- Menghentikan service
- Menghapus semua file
- Menghapus systemd service
- **⚠️ Menghapus database!**

---

## 🔧 Konfigurasi

### File Konfigurasi

Lokasi: `/opt/amaliah-ramadhan/.env`

```env
# Amaliah Ramadhan Configuration
APP_NAME=Amaliah Ramadhan
APP_PORT=8080
APP_ENV=production

# Database
DB_PATH=/opt/amaliah-ramadhan/amaliah.db

# Security
JWT_SECRET=<random-generated-key>
```

### Ubah Port

1. Edit file konfigurasi:
   ```bash
   sudo nano /opt/amaliah-ramadhan/.env
   ```

2. Ubah `APP_PORT=8080` ke port yang diinginkan

3. Restart service:
   ```bash
   sudo systemctl restart amaliah-ramadhan
   ```

---

## 🔐 Keamanan

### Firewall Setup

```bash
# Allow aplikasi port (contoh: 8080)
sudo ufw allow 8080/tcp

# Atau batasi akses dari subnet tertentu
sudo ufw allow from 192.168.1.0/24 to any port 8080
```

### Reverse Proxy dengan Nginx

```nginx
server {
    listen 80;
    server_name amaliah.yourdomain.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### SSL dengan Let's Encrypt

```bash
sudo apt install certbot python3-certbot-nginx
sudo certbot --nginx -d amaliah.yourdomain.com
```

---

## 📝 Lokasi File Penting

```
/opt/amaliah-ramadhan/
├── amaliah-ramadhan          # Binary executable
├── .env                       # Konfigurasi
├── amaliah.db                 # Database SQLite
├── amaliah.db.backup.*        # Backup database
└── web/                       # Static files
    ├── static/
    │   ├── css/
    │   ├── js/
    │   ├── images/
    │   └── fonts/
    └── templates/

/etc/systemd/system/
└── amaliah-ramadhan.service   # Systemd service

/var/log/
├── amaliah-ramadhan.log       # Application log
└── amaliah-ramadhan.error.log # Error log
```

---

## 🐛 Troubleshooting

### Service tidak mau start

```bash
# Cek log error
sudo journalctl -u amaliah-ramadhan -n 50

# Cek file permissions
ls -la /opt/amaliah-ramadhan/

# Test run manual
cd /opt/amaliah-ramadhan
sudo ./amaliah-ramadhan
```

### Port sudah digunakan

```bash
# Cek port yang digunakan
sudo netstat -tlnp | grep 8080

# Atau gunakan lsof
sudo lsof -i :8080

# Ganti port di .env
sudo nano /opt/amaliah-ramadhan/.env
```

### Database error

```bash
# Cek database file
ls -la /opt/amaliah-ramadhan/amaliah.db

# Restore dari backup
sudo cp /opt/amaliah-ramadhan/amaliah.db.backup.* /opt/amaliah-ramadhan/amaliah.db

# Reset database (HATI-HATI: Data hilang!)
sudo rm /opt/amaliah-ramadhan/amaliah.db
sudo systemctl restart amaliah-ramadhan
```

### Memory/CPU tinggi

```bash
# Monitor resource usage
htop

# Limit memory untuk service (edit service file)
sudo systemctl edit amaliah-ramadhan

# Tambahkan:
[Service]
MemoryLimit=512M
CPUQuota=50%
```

---

## 📊 Monitoring

### Resource Usage

```bash
# CPU & Memory
top -p $(pgrep amaliah-ramadhan)

# Disk usage
du -sh /opt/amaliah-ramadhan/*

# Database size
ls -lh /opt/amaliah-ramadhan/amaliah.db
```

### Log Monitoring

```bash
# Real-time log
sudo journalctl -u amaliah-ramadhan -f

# Log hari ini
sudo journalctl -u amaliah-ramadhan --since today

# Log dengan filter error
sudo journalctl -u amaliah-ramadhan -p err
```

---

## 🔄 Backup & Restore

### Backup Manual

```bash
# Backup database
sudo cp /opt/amaliah-ramadhan/amaliah.db \
       /backup/amaliah-$(date +%Y%m%d).db

# Backup full aplikasi
sudo tar -czf /backup/amaliah-full-$(date +%Y%m%d).tar.gz \
       /opt/amaliah-ramadhan/
```

### Restore

```bash
# Restore database
sudo systemctl stop amaliah-ramadhan
sudo cp /backup/amaliah-20260217.db /opt/amaliah-ramadhan/amaliah.db
sudo systemctl start amaliah-ramadhan
```

### Automated Backup (Cron)

```bash
# Edit crontab
sudo crontab -e

# Tambahkan (backup setiap hari jam 2 pagi)
0 2 * * * cp /opt/amaliah-ramadhan/amaliah.db /backup/amaliah-$(date +\%Y\%m\%d).db
```

---

## 🎯 Performance Tips

### 1. Gunakan Reverse Proxy
Nginx untuk serve static files lebih cepat

### 2. Enable Compression
Gzip untuk mengurangi bandwidth

### 3. Database Optimization
Vacuum SQLite secara berkala:
```bash
sqlite3 /opt/amaliah-ramadhan/amaliah.db "VACUUM;"
```

### 4. Log Rotation
Atur logrotate untuk mencegah log terlalu besar

---

## 📞 Support

Untuk bantuan lebih lanjut, hubungi administrator sistem atau tim development.

**SMK NIBA Isep Misbah**
Monitoring Ibadah Harian Ramadhan
