# 🐧 Panduan Deployment dengan Systemd (Golang & SvelteKit 2) — Enterprise ERP Monolith

Dokumen ini adalah panduan resmi untuk men-deploy dan mengelola layanan backend **Golang (Echo v5)** dan frontend **SvelteKit 2 (Svelte 5 Runes)** pada server Linux (**Ubuntu 22.04 / 24.04 LTS**) menggunakan **Systemd Service Manager**.

---

## 📑 Daftar Isi
1. [Arsitektur Deployment & Topologi Layanan](#1-arsitektur-deployment--topologi-layanan)
2. [Persiapan Sistem & Akun Service Dedicated](#2-persiapan-sistem--akun-service-dedicated)
3. [Deployment Backend (Golang Echo v5)](#3-deployment-backend-golang-echo-v5)
   - [A. Build Binary Go Produksi](#a-build-binary-go-produksi)
   - [B. Berkas Environment Backend (`/etc/erp/erp-backend.env`)](#b-berkas-environment-backend)
   - [C. Berkas Unit Systemd Backend (`erp-backend.service`)](#c-berkas-unit-systemd-backend)
4. [Deployment Frontend (SvelteKit 2 + Bun / Node.js)](#4-deployment-frontend-sveltekit-2--bun--nodejs)
   - [A. Konfigurasi Adapter Node (`svelte.config.js`)](#a-konfigurasi-adapter-node)
   - [B. Build Produksi Frontend](#b-build-produksi-frontend)
   - [C. Berkas Environment Frontend (`/etc/erp/erp-frontend.env`)](#c-berkas-environment-frontend)
   - [D. Berkas Unit Systemd Frontend (`erp-frontend.service`)](#d-berkas-unit-systemd-frontend)
5. [Manajemen Terpadu dengan Systemd Target (`erp.target`)](#5-manajemen-terpadu-dengan-systemd-target)
6. [Integrasi Reverse Proxy (Nginx)](#6-integrasi-reverse-proxy-nginx)
7. [Skrip Otomasi Deployment (`deploy.sh`)](#7-skrip-otomasi-deployment)
8. [Pemantauan, Status & Manajemen Log (`journalctl`)](#8-pemantauan-status--manajemen-log)
9. [Troubleshooting & Solusi Masalah Umum](#9-troubleshooting--solusi-masalah-umum)

---

## 1. Arsitektur Deployment & Topologi Layanan

Pada lingkungan produksi bare-metal atau Cloud VPS, sistem ERP Monolith beroperasi dengan pemisahan peran yang terlindungi di balik Reverse Proxy / API Gateway:

```text
                                  INTERNET / LAN
                                        │
                                        ▼ HTTPS (:443)
                         ┌─────────────────────────────┐
                         │   Nginx / Kong Gateway      │
                         │   (SSL Termination & Cache) │
                         └──────────────┬──────────────┘
                                        │
             ┌──────────────────────────┴──────────────────────────┐
             │ Path: /api/*, /health                               │ Path: /* (SSR & Static)
             ▼                                                     ▼
┌──────────────────────────────┐                         ┌──────────────────────────────┐
│  erp-backend.service         │                         │  erp-frontend.service        │
│  • Golang Echo v5 Native     │                         │  • SvelteKit 2 (Node / Bun)  │
│  • Internal Port: :8080      │                         │  • Internal Port: :3000      │
│  • Managed by Systemd        │                         │  • Managed by Systemd        │
└──────────────┬───────────────┘                         └──────────────────────────────┘
               │
               ▼ TCP (:6432)
┌──────────────────────────────┐
│  pgbouncer.service           │
│  • Transaction Pooler Proxy  │
└──────────────┬───────────────┘
               │
               ▼ TCP (:5432)
┌──────────────────────────────┐
│  postgresql.service          │
│  • PostgreSQL 16 DB Server   │
└──────────────────────────────┘
```

### Keunggulan Menggunakan Systemd:
- **Auto-restart:** Memastikan aplikasi otomatis bangkit saat crash (*self-healing*) atau setelah server reboot.
- **Isolasi Keamanan (Hardening):** Menjalankan service dengan user non-privileged, membatasi akses filesystem (`ProtectSystem`, `ProtectHome`), dan membatasi kapabilitas OS (`NoNewPrivileges`).
- **Standardisasi Logging:** Semua output `stdout`/`stderr` otomatis dikelola oleh `journalctl` lengkap dengan timestamp presisi dan rotasi log terpadu.
- **Dependency Ordering:** Memastikan backend tidak menyala sebelum database dan pooler PgBouncer benar-benar siap menerima koneksi.

---

## 2. Persiapan Sistem & Akun Service Dedicated

Untuk alasan keamanan, **jangan pernah menjalankan aplikasi backend atau frontend dengan user `root`**. Buat akun service sistem khusus tanpa hak login shell interaktif.

### Langkah 1: Buat User & Grup Khusus `erp`
```bash
sudo useradd -r -s /usr/sbin/nologin -d /opt/dev/erp_monolith erp
```

### Langkah 2: Buat Direktori Konfigurasi & Log
```bash
# Direktori untuk berkas variabel lingkungan sensitif
sudo mkdir -p /etc/erp
sudo chown -R erp:erp /etc/erp
sudo chmod 750 /etc/erp

# Direktori target rilis jika menggunakan path terpisah (atau gunakan direktori workspace aktif)
sudo mkdir -p /opt/dev/erp_monolith/backend/bin
sudo chown -R erp:erp /opt/dev/erp_monolith
```

---

## 3. Deployment Backend (Golang Echo v5)

### A. Build Binary Go Produksi
Kompilasi source code Go menjadi binary mandiri (*statically linked*) yang telah dioptimasi dan dibersihkan dari simbol debug untuk meminimalkan ukuran file dan meningkatkan performa:

```bash
cd /opt/dev/erp_monolith/backend

# Kompilasi rilis produksi
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/erp-backend cmd/server/main.go

# Pastikan permission eksekusi
sudo chmod +x bin/erp-backend
sudo chown erp:erp bin/erp-backend
```

### B. Berkas Environment Backend
Simpan variabel lingkungan di `/etc/erp/erp-backend.env` dan kunci hak aksesnya (`600`):

```bash
sudo nano /etc/erp/erp-backend.env
```

Isi berkas konfigurasi produksi:
```ini
# ========================================================
# Enterprise ERP Monolith — Backend Production Environment
# ========================================================

SERVER_PORT=8080
SERVER_ENV=production

# Koneksi ke PgBouncer (Port 6432 - Transaction Pooling)
DATABASE_URL=postgres://erp_user:erp_secret@127.0.0.1:6432/erp_db?sslmode=disable

# Koneksi Direct PostgreSQL untuk DDL Migrasi (Port 5432)
DIRECT_DATABASE_URL=postgres://erp_user:erp_secret@127.0.0.1:5432/erp_db?sslmode=disable

# Pengaturan Connection Pool pgxpool
DB_MAX_CONNS=50
DB_MIN_CONNS=10
DB_MAX_CONN_LIFETIME=1h
DB_MAX_CONN_IDLE_TIME=30m

# Keamanan & JWT
JWT_SECRET=GantiDenganStringRahasiaSangatPanjangMinimal32KarakterAcak!
JWT_EXPIRATION_HOURS=24

# CORS (Izinkan domain frontend publik dan reverse proxy)
CORS_ALLOWED_ORIGINS=https://erp.perusahaan.com,http://localhost:3000
```

Kunci permission:
```bash
sudo chown erp:erp /etc/erp/erp-backend.env
sudo chmod 600 /etc/erp/erp-backend.env
```

### C. Berkas Unit Systemd Backend
Buat berkas unit systemd di `/etc/systemd/system/erp-backend.service`:

```bash
sudo nano /etc/systemd/system/erp-backend.service
```

Salin konfigurasi standar enterprise berikut:

```ini
[Unit]
Description=Enterprise ERP Monolith - Golang Echo v5 Backend Service
Documentation=https://github.com/jeruktutut2/arp-monolith
After=network.target network-online.target pgbouncer.service redis.service
Wants=network-online.target pgbouncer.service
# Jika pgbouncer dijalankan di mesin yang sama, pastikan pgbouncer start terlebih dahulu

[Service]
Type=exec
User=erp
Group=erp
WorkingDirectory=/opt/dev/erp_monolith/backend
EnvironmentFile=/etc/erp/erp-backend.env
ExecStart=/opt/dev/erp_monolith/backend/bin/erp-backend

# Graceful restart & reload
ExecReload=/bin/kill -HUP $MAINPID
KillMode=process
KillSignal=SIGTERM
TimeoutStopSec=30s
Restart=always
RestartSec=5s

# -------------------------------------------------------------
# Hardening Keamanan (Security Sandbox)
# -------------------------------------------------------------
NoNewPrivileges=true
ProtectSystem=full
ProtectHome=true
PrivateTmp=true
ProtectKernelTunables=true
ProtectKernelModules=true
ProtectControlGroups=true
RestrictAddressFamilies=AF_INET AF_INET6 AF_UNIX
RestrictRealtime=true
CapabilityBoundingSet=

# Resource Limits
LimitNOFILE=65536
LimitNPROC=32768

# Logging
StandardOutput=journal
StandardError=journal
SyslogIdentifier=erp-backend

[Install]
WantedBy=multi-user.target
```

---

## 4. Deployment Frontend (SvelteKit 2 + Bun / Node.js)

### A. Konfigurasi Adapter Node
Aplikasi SvelteKit yang di-deploy dengan systemd membutuhkan adapter Node.js (`@sveltejs/adapter-node`) agar dapat dijalankan sebagai web server mandiri.

Pastikan pada `frontend/svelte.config.js`:
```javascript
import adapter from '@sveltejs/adapter-node';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';

/** @type {import('@sveltejs/kit').Config} */
const config = {
  preprocess: vitePreprocess(),
  kit: {
    adapter: adapter({
      out: 'build',
      precompress: true, // Gzip & Brotli pre-compression
      polyfill: false
    })
  }
};

export default config;
```

### B. Build Produksi Frontend
Jalankan kompilasi menggunakan **Bun** (atau Node.js):

```bash
cd /opt/dev/erp_monolith/frontend

# Pasang dependencies yang presisi
bun install --frozen-lockfile

# Kompilasi aset produksi
bun run build

# Pastikan kepemilikan folder dimiliki oleh user erp
sudo chown -R erp:erp /opt/dev/erp_monolith/frontend
```

Hasil build akan tersimpan di direktori `/opt/dev/erp_monolith/frontend/build`.

### C. Berkas Environment Frontend
Buat berkas konfigurasi di `/etc/erp/erp-frontend.env`:

```bash
sudo nano /etc/erp/erp-frontend.env
```

Isi berkas:
```ini
# ========================================================
# Enterprise ERP Monolith — Frontend Production Environment
# ========================================================

NODE_ENV=production
HOST=127.0.0.1
PORT=3000

# ORIGIN WAJIB didefinisikan untuk mencegah error proteksi SvelteKit CSRF
ORIGIN=https://erp.perusahaan.com

# Ukuran maksimum payload request upload dokumen/berkas (default 512KB, naikkan ke 50MB)
BODY_SIZE_LIMIT=52428800

# URL Internal Gateway / Backend untuk SSR (Server-Side Rendering call)
INTERNAL_API_URL=http://127.0.0.1:8080

# URL Publik untuk pemanggilan dari browser client
PUBLIC_API_URL=https://erp.perusahaan.com/api
```

Kunci hak akses:
```bash
sudo chown erp:erp /etc/erp/erp-frontend.env
sudo chmod 600 /etc/erp/erp-frontend.env
```

### D. Berkas Unit Systemd Frontend
Buat berkas unit di `/etc/systemd/system/erp-frontend.service`:

```bash
sudo nano /etc/systemd/system/erp-frontend.service
```

Pilih runtime yang digunakan pada `ExecStart`:

#### Opsi 1: Menggunakan Node.js (Disarankan untuk stabilitas `@sveltejs/adapter-node`)
```ini
[Unit]
Description=Enterprise ERP Monolith - SvelteKit 2 Frontend Service (Node.js)
Documentation=https://github.com/jeruktutut2/arp-monolith
After=network.target erp-backend.service
Wants=erp-backend.service

[Service]
Type=exec
User=erp
Group=erp
WorkingDirectory=/opt/dev/erp_monolith/frontend
EnvironmentFile=/etc/erp/erp-frontend.env
ExecStart=/usr/bin/node ./build/index.js

KillMode=process
KillSignal=SIGTERM
TimeoutStopSec=20s
Restart=always
RestartSec=5s

# Hardening Keamanan
NoNewPrivileges=true
ProtectSystem=full
ProtectHome=true
PrivateTmp=true

# Resource Limits
LimitNOFILE=65536

# Logging
StandardOutput=journal
StandardError=journal
SyslogIdentifier=erp-frontend

[Install]
WantedBy=multi-user.target
```

#### Opsi 2: Menggunakan Bun Runtime
Jika memilih menjalankan langsung di atas Bun, ubah baris `ExecStart`:
```ini
# Cari path bun terlebih dahulu dengan: which bun (biasanya /usr/local/bin/bun atau /home/user/.bun/bin/bun)
ExecStart=/usr/local/bin/bun ./build/index.js
```

---

## 5. Manajemen Terpadu dengan Systemd Target (`erp.target`)

Agar administrator dapat menyalakan, mematikan, me-restart, atau memeriksa status seluruh stack ERP (backend + frontend) dengan satu perintah, buat unit target gabungan:

```bash
sudo nano /etc/systemd/system/erp.target
```

Isi berkas `erp.target`:
```ini
[Unit]
Description=Enterprise ERP Monolith Application Stack
Wants=erp-backend.service erp-frontend.service

[Install]
WantedBy=multi-user.target
```

### Mengaktifkan dan Menjalankan Service:
```bash
# 1. Reload systemd daemon untuk membaca berkas unit baru
sudo systemctl daemon-reload

# 2. Aktifkan auto-start saat boot sistem
sudo systemctl enable erp-backend.service
sudo systemctl enable erp-frontend.service
sudo systemctl enable erp.target

# 3. Nyalakan seluruh stack aplikasi
sudo systemctl start erp.target

# 4. Verifikasi status kedua layanan
sudo systemctl status erp-backend erp-frontend
```

---

## 6. Integrasi Reverse Proxy (Nginx)

Nginx bertindak sebagai front-facing reverse proxy yang menangani SSL, kompresi, dan mendistribusikan request ke backend Go (`:8080`) atau frontend SvelteKit (`:3000`).

Buat konfigurasi virtual host di `/etc/nginx/sites-available/erp.perusahaan.com`:

```nginx
# ========================================================
# Nginx Reverse Proxy — Enterprise ERP Monolith
# ========================================================

upstream erp_frontend_upstream {
    server 127.0.0.1:3000;
    keepalive 32;
}

upstream erp_backend_upstream {
    server 127.0.0.1:8080;
    keepalive 32;
}

server {
    listen 80;
    server_name erp.perusahaan.com;
    return 301 https://$host$request_uri;
}

server {
    listen 443 ssl http2;
    server_name erp.perusahaan.com;

    # Konfigurasi Sertifikat SSL
    ssl_certificate /etc/letsencrypt/live/erp.perusahaan.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/erp.perusahaan.com/privkey.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;

    # Ukuran maksimum body request (Upload lampiran invoice/faktur)
    client_max_body_size 50M;

    # Header Keamanan
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;

    # -------------------------------------------------------------
    # 1. Routing API ke Golang Echo v5 Backend (:8080)
    # -------------------------------------------------------------
    location /api/ {
        proxy_pass http://erp_backend_upstream;
        proxy_http_version 1.1;
        
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        # WebSocket / Server-Sent Events (SSE) Support
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";

        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }

    # Health check endpoint langsung ke Go backend
    location = /health {
        proxy_pass http://erp_backend_upstream/health;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
    }

    # -------------------------------------------------------------
    # 2. Routing Aplikasi Utama ke SvelteKit 2 Frontend (:3000)
    # -------------------------------------------------------------
    location / {
        proxy_pass http://erp_frontend_upstream;
        proxy_http_version 1.1;

        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        # Mendukung WebSocket (HMR atau fitur real-time SvelteKit)
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }

    # Cache untuk aset statis SvelteKit (_app/immutable)
    location /_app/immutable/ {
        proxy_pass http://erp_frontend_upstream;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        expires 1y;
        add_header Cache-Control "public, immutable";
    }
}
```

Aktifkan konfigurasi Nginx:
```bash
sudo ln -s /etc/nginx/sites-available/erp.perusahaan.com /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

---

## 7. Skrip Otomasi Deployment (`deploy.sh`)

Buat skrip deployment zero-downtime atomik di `/opt/dev/erp_monolith/deploy.sh`:

```bash
#!/usr/bin/env bash
set -e

APP_DIR="/opt/dev/erp_monolith"
BACKEND_DIR="${APP_DIR}/backend"
FRONTEND_DIR="${APP_DIR}/frontend"

echo "=========================================================="
echo "🚀 Memulai Deployment Enterprise ERP Monolith..."
echo "=========================================================="

cd "${APP_DIR}"

# 1. Tarik pembaruan kode Git (jika di server CI/CD)
# git fetch origin main && git reset --hard origin/main

# 2. Migrasi Database (golang-migrate ke port direct 5432)
echo "📦 Menjalankan migrasi database..."
if command -v migrate >/dev/null 2>&1; then
    migrate -path "${BACKEND_DIR}/migrations" -database "${DIRECT_DATABASE_URL:-postgres://erp_user:erp_secret@127.0.0.1:5432/erp_db?sslmode=disable}" up || {
        echo "⚠️ Catatan: Migrasi gagal atau tidak ada perubahan skema baru."
    }
fi

# 3. Kompilasi Backend Golang (Atomic Swap)
echo "🔨 Mengompilasi backend Golang..."
cd "${BACKEND_DIR}"
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/erp-backend.new cmd/server/main.go
chmod +x bin/erp-backend.new

# Atomic replace binary
mv -f bin/erp-backend.new bin/erp-backend
chown erp:erp bin/erp-backend

# 4. Kompilasi Frontend SvelteKit
echo "🎨 Mengompilasi frontend SvelteKit..."
cd "${FRONTEND_DIR}"
bun install --frozen-lockfile
bun run build
chown -R erp:erp "${FRONTEND_DIR}/build"

# 5. Restart Layanan via Systemd
echo "🔄 Merestart service systemd..."
sudo systemctl restart erp-backend.service
sudo systemctl restart erp-frontend.service

# 6. Verifikasi Health Check
echo "🔍 Melakukan verifikasi health check..."
sleep 3
if curl -fs "http://127.0.0.1:8080/health" > /dev/null; then
    echo "✅ Backend Health Check: OK"
else
    echo "❌ ERROR: Backend gagal merespons /health!"
    sudo systemctl status erp-backend.service --no-pager
    exit 1
fi

if curl -fs "http://127.0.0.1:3000" > /dev/null; then
    echo "✅ Frontend Health Check: OK"
else
    echo "⚠️ Peringatan: Frontend belum siap atau merespons kode non-200."
fi

echo "=========================================================="
echo "🎉 Deployment Berhasil Selesai!"
echo "=========================================================="
```

Beri izin eksekusi:
```bash
chmod +x /opt/dev/erp_monolith/deploy.sh
```

---

## 8. Pemantauan, Status & Manajemen Log (`journalctl`)

### Perintah Operasional Harian
| Operasi | Perintah |
|---|---|
| **Cek Status Backend** | `sudo systemctl status erp-backend` |
| **Cek Status Frontend** | `sudo systemctl status erp-frontend` |
| **Cek Semua Layanan ERP** | `sudo systemctl status erp.target erp-backend erp-frontend` |
| **Restart Backend Saja** | `sudo systemctl restart erp-backend` |
| **Restart Frontend Saja** | `sudo systemctl restart erp-frontend` |
| **Restart Keseluruhan** | `sudo systemctl restart erp.target` |
| **Stop Layanan** | `sudo systemctl stop erp.target` |

### Membaca Log Aplikasi Real-time (`journalctl`)
```bash
# Streaming log backend secara live
sudo journalctl -u erp-backend -f

# Streaming log frontend secara live
sudo journalctl -u erp-frontend -f

# Streaming kedua service secara bersamaan
sudo journalctl -u erp-backend -u erp-frontend -f

# Membaca 100 baris log terakhir
sudo journalctl -u erp-backend -n 100 --no-pager

# Membaca log error saja sejak 1 jam yang lalu
sudo journalctl -u erp-backend --since "1 hour ago" -p err
```

### Konfigurasi Batas Ukuran Log (Mencegah Disk Penuh)
Tambahkan konfigurasi pada `/etc/systemd/journald.conf`:
```ini
[Journal]
SystemMaxUse=500M
SystemMaxFileSize=50M
MaxRetentionSec=1month
```
Terapkan dengan:
```bash
sudo systemctl restart systemd-journald
```

---

## 9. Troubleshooting & Solusi Masalah Umum

### 1. SvelteKit CSRF Error: `Cross-site POST form submissions are forbidden`
- **Penyebab:** SvelteKit 2 memeriksa header `Origin` saat menerima form submission POST. Jika reverse proxy tidak mengirimkan header yang cocok dengan `ORIGIN`, request akan ditolak (status 403).
- **Solusi:** 
  1. Pastikan pada `/etc/erp/erp-frontend.env`:
     ```ini
     ORIGIN=https://erp.perusahaan.com
     ```
  2. Pastikan di konfigurasi Nginx terdapat:
     ```nginx
     proxy_set_header Host $host;
     proxy_set_header X-Forwarded-Proto $scheme;
     ```

---

### 2. Request Entity Too Large (Status 413) saat Upload File
- **Penyebab:** Batas default SvelteKit adalah 512KB dan Nginx adalah 1MB.
- **Solusi:**
  1. Pada SvelteKit (`/etc/erp/erp-frontend.env`):
     ```ini
     BODY_SIZE_LIMIT=52428800 # 50MB dalam bytes
     ```
  2. Pada Nginx block `server`:
     ```nginx
     client_max_body_size 50M;
     ```

---

### 3. Backend Gagal Konek ke Database saat Booting Ulang Server
- **Penyebab:** `erp-backend.service` menyala lebih cepat daripada `pgbouncer.service`.
- **Solusi:** Pastikan pada `/etc/systemd/system/erp-backend.service` terdapat dependensi:
  ```ini
  [Unit]
  After=network.target pgbouncer.service
  Wants=pgbouncer.service
  ```

---

### 4. Eksekusi Node atau Bun Gagal: `code=exited, status=203/EXEC`
- **Penyebab:** Path binary Node.js atau Bun di direktori `/usr/local/bin` atau nvm tidak ditemukan oleh systemd. Systemd tidak memuat file profile shell (`.bashrc` / `.zshrc`).
- **Solusi:** Periksa lokasi absolut binary dengan:
  ```bash
  which node
  which bun
  ```
  Gunakan path absolut lengkap tersebut pada direktif `ExecStart` di berkas `.service`.

---

### 5. Port Permission Denied: `bind: permission denied`
- **Penyebab:** Service dijalankan oleh user non-root (`erp`) dan mencoba mengikat (*bind*) ke port istimewa (*privileged port* < 1024) seperti port 80 atau 443 secara langsung.
- **Solusi:** Jangan bind Go atau SvelteKit langsung ke port 80/443. Biarkan backend mendengarkan port 8080 dan frontend port 3000, lalu gunakan Nginx atau izinkan capability:
  ```bash
  sudo setcap 'cap_net_bind_service=+ep' /opt/dev/erp_monolith/backend/bin/erp-backend
  ```
