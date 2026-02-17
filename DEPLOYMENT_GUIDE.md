# Cara Deploy ke VPS (Secure & Proper) 🚀

Guide ini akan membantu Anda men-deploy **Muslimly Backend** ke VPS (Ubuntu/Debian) menggunakan Docker.

## Persiapan

1.  **VPS (Virtual Private Server)**: OS Ubuntu 22.04 / 24.04.
2.  **Domain**: (Optional tapi Recommended) misal `api.muslimly.my.id`.
3.  **File Rahasia**: Pastikan Anda punya backup `config.yaml` dan `firebase-service-account.json` di Laptop.

## Langkah 1: Install Docker di VPS

Login ke VPS via SSH, lalu jalankan:

```bash
# Update Repo
sudo apt update && sudo apt upgrade -y

# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# Install Docker Compose
sudo apt install -y docker-compose-plugin
```

## Langkah 2: Clone Repository

```bash
git clone https://github.com/imamrisnandar/muslimly-be.git
cd muslimly-be
```

## Langkah 3: Upload File Rahasia (PENTING!)

Karena file ini tidak ada di GitHub (gitignore), Anda harus upload manual dari Laptop ke VPS.

**Cara Upload (Dari Laptop):**
Buka Terminal/CMD di folder project laptop, lalu jalankan `scp`:

```powershell
# Upload Config
scp config.yaml root@IP_VPS:/root/muslimly-be/

# Upload Firebase Key
scp firebase-service-account.json root@IP_VPS:/root/muslimly-be/
```

_(Ganti `IP_VPS` dengan IP Server Anda)_.

## Langkah 4: Setup Environment Variable

Di VPS, buat file `.env` untuk password database:

```bash
nano .env
```

Isi dengan:

```env
DB_USER=postgres
DB_PASSWORD=PasswordRahasiaAnda123
DB_NAME=muslimly
```

_(Simpan dengan CTRL+O, Enter, CTRL+X)_

## Langkah 5: Jalankan Server 🚀

```bash
docker compose up -d --build
```

Cek status:

```bash
docker compose ps
# Pastikan status "Up"
```

Cek logs:

```bash
docker compose logs -f app
```

## Langkah 6: (Optional) Setup SSL dengan Nginx & Certbot

Jika sudah punya domain, install Nginx Proxy Manager atau Caddy agar HTTPS aktif.

Selesai! API Anda sekarang aktif di `http://IP_VPS:8080`.
Jangan lupa update IP ini di source code Flutter (`network_module.dart`).
