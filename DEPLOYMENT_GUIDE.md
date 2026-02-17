# 🚀 Muslimly Backend - Full Deployment Guide

Panduan lengkap untuk men-deploy, mengamankan, dan memonitoring Muslimly Backend di VPS Production.

---

## 🛠️ 1. Persiapan (Prerequisites)

Sebelum mulai, pastikan Anda memiliki:

1.  **VPS (Ubuntu 22.04/24.04)**: Minimal 1GB RAM.
2.  **Domain**: Terhubung ke IP VPS (A Record). Contoh: `muslimly.my.id`.
3.  **User VPS**: `lighthouse` (Tencent Cloud default) dengan akses sudo.
4.  **Local Tools**:
    - Terminal (PowerShell/Bash)
    - Database Client (DBeaver / PgAdmin / TablePlus)
    - SSH Client (OpenSSH)
5.  **File Rahasia (Local)**:
    - `cmd/tools/vps_key.pem` (Private Key SSH - **Jangan Commit ke Git!**)
    - `.secrets` (Database Credentials - **Jangan Commit ke Git!**)
    - `docker-compose.prod.yml` (Config Production + Security Fixes)
    - `config.yaml` (App Config)
    - `firebase-service-account.json` (Firebase Admin SDK)

---

## 📦 2. Deployment Step-by-Step

### Step 1: Install Docker & Docker Compose (di VPS)

Login ke VPS (`ssh root@IP_VPS`) dan jalankan:

```bash
# Update System
sudo apt update && sudo apt upgrade -y

# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# Install Docker Compose Plugin
sudo apt install -y docker-compose-plugin
```

### Step 2: Clone Repository

```bash
git clone https://github.com/imamrisnandar/muslimly-be.git
cd muslimly-be
```

### Step 3: Upload Konfigurasi Aman (Dari Local Laptop)

File konfigurasi sensitif **TIDAK BOLEH** ada di GitHub. Upload manual dari laptop Anda menggunakan `scp`:

```powershell
# Ganti IP_VPS dengan IP Server Anda
$VPS_IP = "43.134.180.17"

# 1. Upload App Config & Firebase Key
scp config.yaml root@$VPS_IP:/root/muslimly-be/
scp firebase-service-account.json root@$VPS_IP:/root/muslimly-be/

# 2. Upload Production Docker Compose (Rename jadi docker-compose.yml di server)
scp docker-compose.prod.yml root@$VPS_IP:/root/muslimly-be/docker-compose.yml
```

> **Note**: `docker-compose.prod.yml` berisi password asli database. File ini di-ignore oleh git agar aman.

### Step 4: Setup SSL / HTTPS (Nginx & Certbot)

Di dalam repo sudah ada folder `nginx/` dan `certbot/`.

1.  **Jalankan Stack**:

    ```bash
    cd ~/muslimly-be
    docker compose up -d
    ```

2.  **Request Certificate (Pertama Kali)**:
    Jika sertifikat belum ada, jalankan perintah ini untuk meminta ke Let's Encrypt:

    ```bash
    docker compose run --rm certbot certonly --webroot --webroot-path /var/www/certbot -d muslimly.my.id
    ```

    _(Ikuti instruksi di layar, pilih Agree)._

3.  **Restart Nginx** untuk memuat sertifikat:
    ```bash
    docker compose restart nginx
    ```

---

## 🗄️ 3. Cara Akses Database (Secure Access)

Database PostgreSQL di VPS **TIDAK** diekspos ke internet publik (Port 5432 ditutup firewall/docker bind). Anda **WAJIB** menggunakan **SSH Tunnel**.

### Konfigurasi DBeaver / PgAdmin / TablePlus:

**Tab SSH / SSH Tunnel:**

- **Host/IP**: `IP_VPS` (ex: 43.134.180.17)
- **Port**: `22`
- **Username**: `root`
- **Password**: _(Password VPS Anda)_

**Tab General / Connection:**

- **Host**: `localhost` (Karena kita masuk lewat tunnel, bagi tunnel ini adalah localhost server)
- **Port**: `5432`
- **Database**: `muslimly`
- **Username**: `postgres`
- **Password**: `Mu5l1mlyJ4nn4h` (Atau sesuai isi `docker-compose.prod.yml`)

---

## 📊 4. Monitoring Server

### Cek Status Container

Lihat apakah semua service (backend, db, nginx) berjalan:

```bash
docker compose ps
# Status harus "Up" atau "Up (healthy)"
```

### Cek Resource Usage (CPU/RAM)

Lihat beban server realtime:

```bash
docker stats
```

### Cek Logs (Debugging)

Melihat log aplikasi backend:

```bash
docker compose logs -f app --tail 100
```

Melihat log Nginx (Access/Error):

```bash
docker compose logs -f nginx --tail 100
```

Melihat log Database:

```bash
docker compose logs -f db --tail 100
```

---

## 🔄 5. Maintenance & Update

### Cara Update Aplikasi (Backend)

Jika ada update code di GitHub:

1.  **Pull Code**:
    ```bash
    cd ~/muslimly-be
    git pull origin main
    ```
2.  **Rebuild & Restart**:
    ```bash
    docker compose up -d --build app
    ```
    _(Hanya service `app` yang direstart, database aman)_.

### Cara Renew SSL (Setiap 2-3 Bulan)

Isi crontab (`crontab -e`) untuk otomatis renew:

```bash
0 0 1 * * docker compose run --rm certbot renew && docker compose restart nginx
```

---

## ⚡ Bonus: Utility Scripts (cmd/tools)

Untuk mempermudah manajemen server, telah disediakan 5 script Golang di folder `cmd/tools`. Anda bisa menjalankannya dari laptop lokal:

### 1. `deploy_vps.go` (Setup Awal)

Script untuk inisialisasi server baru. Melakukan instalasi Docker, Nginx, dan setup environment dasar.

```bash
go run cmd/tools/deploy_vps.go
```

### 2. `deploy_update.go` (Update Harian)

Script standar untuk deploy update code. Melakukan `git pull` dan `docker compose build`. Gunakan ini jika Anda push perubahan code ke GitHub.

```bash
go run cmd/tools/deploy_update.go
```

### 3. `deploy_reset.go` (Emergency Reset)

Script darurat jika terjadi konflik file di server. Melakukan `git reset --hard` untuk memaksa file di server sama persis dengan GitHub, lalu rebuild.

```bash
go run cmd/tools/deploy_reset.go
```

### 4. `check_health.go` (Validasi & Health Check)

Script untuk mengecek status deployment. Menampilkan status container, log error terakhir, schema database, dan mengetes endpoint API `/health`.

```bash
go run cmd/tools/check_health.go
```

### 5. `setup_ssl_robust.go` (Setup SSL)

Script khusus untuk install atau memperbaiki sertifikat SSL (HTTPS) menggunakan Certbot.

```bash
go run cmd/tools/setup_ssl_robust.go
```

---

## 🚨 Troubleshooting

**Q: Database Connection Refused?**
A: Pastikan di `docker list`, container `muslimly_db` statusnya Up. Cek logs: `docker logs muslimly_db`.

**Q: Website HTTPS Error?**
A: Cek masa berlaku sertifikat. Cek log nginx: `docker logs muslimly_nginx`. Pastikan port 443 terbuka di Firewall VPS.

**Q: Aplikasi Error "Connection Refused" ke DB?**
A: Pastikan di `config.yaml` dan `docker-compose.yml`, `DB_HOST` diset ke nama service container (`muslimly_db`), bukan localhost.
