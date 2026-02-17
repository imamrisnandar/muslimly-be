# 🚀 Muslimly Backend - Full Deployment Guide

Panduan lengkap untuk men-deploy, mengamankan, dan memonitoring Muslimly Backend di VPS Production.

---

## 🛠️ 1. Persiapan (Prerequisites)

Sebelum mulai, pastikan Anda memiliki:

1.  **VPS (Ubuntu 22.04/24.04)**: Minimal 1GB RAM.
2.  **Domain**: Terhubung ke IP VPS (A Record). Contoh: `muslimly.my.id`.
3.  **Local Tools**:
    - Terminal (PowerShell/Bash)
    - Database Client (DBeaver / PgAdmin / TablePlus)
    - SSH Client (OpenSSH)
4.  **File Rahasia (Local)**:
    - `docker-compose.prod.yml` (Config Production)
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

## ⚡ Bonus: Script Deploy Otomatis (Golang)

Jika malas menjalankan command manual berulang kali, simpan script ini sebagai `deploy_tools.go` di laptop Anda.
Script ini akan:

1.  Mengupload file config terbaru.
2.  Upload `docker-compose.prod.yml` (sebagai `docker-compose.yml`).
3.  Merestart container di VPS.

### Script (`deploy_tools.go`)

```go
package main

import (
    "fmt"
    "io"
    "log"
    "os"
    "time"

    "golang.org/x/crypto/ssh"
)

const (
    // GANTI BAGIAN INI DENGAN DATA ANDA
    VPS_IP       = "43.134.180.17"
    VPS_USER     = "root"
    VPS_PASSWORD = "PASSWORD_VPS_ANDA" // Atau gunakan ssh.PublicKeys jika ada key file
    APP_DIR      = "/root/muslimly-be"
)

func main() {
    config := &ssh.ClientConfig{
        User: VPS_USER,
        Auth: []ssh.AuthMethod{
            ssh.Password(VPS_PASSWORD),
        },
        HostKeyCallback: ssh.InsecureIgnoreHostKey(),
        Timeout:         10 * time.Second,
    }

    client, err := ssh.Dial("tcp", VPS_IP+":22", config)
    if err != nil {
        log.Fatalf("Failed to dial: %v", err)
    }
    defer client.Close()

    fmt.Println("--- Connected to VPS ---")

    // 1. Upload Configs
    fmt.Println("\n[1] Uploading Configs...")
    uploadFile(client, "config.yaml", APP_DIR+"/config.yaml")

    // Upload Prod Compose as default compose
    uploadFile(client, "docker-compose.prod.yml", APP_DIR+"/docker-compose.yml")


    // 2. Restart Stack
    fmt.Println("\n[2] Restarting Stack...")
    runCommand(client, "cd "+APP_DIR+" && docker compose down --remove-orphans")
    runCommand(client, "cd "+APP_DIR+" && docker compose up -d --build")

    // 3. Status
    fmt.Println("\n[3] Checking Status...")
    time.Sleep(5 * time.Second)
    runCommand(client, "docker ps")
}

func runCommand(client *ssh.Client, cmd string) {
    session, err := client.NewSession()
    if err != nil {
        log.Fatalf("Failed: %v", err)
    }
    defer session.Close()
    session.Stdout = os.Stdout
    session.Stderr = os.Stderr
    if err := session.Run(cmd); err != nil {
        log.Printf("Cmd error: %v", err)
    }
}

func uploadFile(client *ssh.Client, localPath, remotePath string) {
    f, err := os.Open(localPath)
    if err != nil {
        log.Fatalf("Open failed: %v", err)
    }
    defer f.Close()
    session, err := client.NewSession()
    defer session.Close()

    go func() {
        w, _ := session.StdinPipe()
        defer w.Close()
        io.Copy(w, f)
    }()

    // Simple cat method
    if err := session.Run("cat > " + remotePath); err != nil {
        log.Fatalf("Upload failed: %v", err)
    }
}
```

**Cara Pakai:**

```bash
go run deploy_tools.go
```

---

## 🚨 Troubleshooting

**Q: Database Connection Refused?**
A: Pastikan di `docker list`, container `muslimly_db` statusnya Up. Cek logs: `docker logs muslimly_db`.

**Q: Website HTTPS Error?**
A: Cek masa berlaku sertifikat. Cek log nginx: `docker logs muslimly_nginx`. Pastikan port 443 terbuka di Firewall VPS.

**Q: Aplikasi Error "Connection Refused" ke DB?**
A: Pastikan di `config.yaml` dan `docker-compose.yml`, `DB_HOST` diset ke nama service container (`muslimly_db`), bukan localhost.
