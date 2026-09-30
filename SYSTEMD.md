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
   - [E. Integrasi Alur Penerimaan Hasil Build dari GitHub Actions](#e-integrasi-alur-penerimaan-hasil-build-dari-github-actions)
   - [F. Langkah Persiapan Server Sebelum Deployment Pertama](#f-langkah-persiapan-server-sebelum-deployment-pertama)
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

> 💡 **User Service vs User Deployer:**
> User `erp` di atas adalah akun sistem non-login khusus untuk daemon Systemd. Untuk akun operasional interaktif yang bertugas melakukan git push/pull, rsync file, restart service, dan menerima koneksi SSH dari GitHub Actions, gunakan user `deployer`. Panduan lengkap konfigurasi user `deployer` dan hak akses SSH tersedia di [SSH_GUIDE.md](file:///opt/dev/erp_monolith/SSH_GUIDE.md).

### Langkah 2: Buat Direktori Konfigurasi & Log
```bash
# Direktori untuk berkas variabel lingkungan sensitif
sudo mkdir -p /etc/erp
sudo chown -R erp:erp /etc/erp
sudo chmod 750 /etc/erp

# Direktori target rilis jika menggunakan path terpisah (atau gunakan direktori workspace aktif)
sudo mkdir -p /opt/dev/erp_monolith/backend/bin
sudo mkdir -p /opt/apps/erp_monolith/frontend
sudo chown -R erp:erp /opt/dev/erp_monolith
sudo chown -R erp:erp /opt/apps
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

> 💡 **Build di VPS vs Menerima Build dari GitHub Actions (CI/CD):**
> Perintah di atas digunakan jika Anda memilih melakukan build lokal langsung di VPS. Namun, pada alur CI/CD produksi yang direkomendasikan ([GITHUB_ACTIONS.md](file:///opt/dev/erp_monolith/GITHUB_ACTIONS.md)), proses `bun run build` dijalankan di GitHub Actions runner. VPS hanya menerima hasil kompilasi folder `frontend/build/` yang ditaruh ke direktori target `/opt/apps/erp_monolith/frontend/` via SSH/Rsync sehingga server terbebas dari lonjakan RAM/CPU saat proses build. Lihat detail penerimaan dan konfigurasinya di [Sub-bab E](#e-integrasi-alur-penerimaan-hasil-build-dari-github-actions).

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

### E. Integrasi Alur Penerimaan Hasil Build dari GitHub Actions (CI/CD Artifact)

Jika Anda menerapkan alur pengiriman otomatis menggunakan **GitHub Actions** (seperti yang dikonfigurasi pada [.github/workflows/deploy.yml](file:///opt/dev/erp_monolith/.github/workflows/deploy.yml) dan [GITHUB_ACTIONS.md](file:///opt/dev/erp_monolith/GITHUB_ACTIONS.md)), proses deployment frontend ke Systemd memiliki alur dan karakteristik khusus:

#### 1. Arsitektur Penerimaan Artifact
- **Build di GitHub Actions Runner**: Runner GitHub mengompilasi SvelteKit 2 menggunakan Bun (`bun install --frozen-lockfile` & `bun run build`) menghasilkan artefak rilis di folder `frontend/build/`.
- **Pengiriman Tanpa Source Code Mentah**: Hanya direktori hasil build terkompilasi yang dikirimkan ke server VPS via SSH/Rsync ke path target produksi:
  ```text
  /opt/apps/erp_monolith/frontend/
  ```
- **Struktur Berkas di VPS**:
  ```text
  /opt/apps/erp_monolith/frontend/
  ├── index.js          # Entrypoint server (Polka / Node HTTP)
  ├── handler.js        # Request handler SvelteKit
  ├── env.js            # Injeksi runtime environment variable
  ├── shims.js          # Node polyfill compatibility
  ├── client/           # Aset statis client browser (_app/immutable/, CSS, JS)
  └── server/           # Server chunks & logika SSR (seluruh pustaka runtime telah dibundle)
  ```
- **Kelebihan**: Server VPS **tidak membutuhkan** instalasi compiler, devDependencies, maupun direktori `node_modules/` (ratusan MB). Hanya runtime minimal Node.js atau Bun yang dibutuhkan untuk mengeksekusi `index.js`.

#### 2. Penanganan Hak Akses & Kepemilikan (User `deployer` vs `erp`)
Ketika GitHub Actions mentransfer berkas menggunakan akun SSH `deployer` (yang tergabung dalam grup `erp`), berkas build baru akan tercatat sebagai milik `deployer:erp`. Agar service Systemd yang dijalankan oleh user `erp` dapat mengakses dan membaca hasil build:

- **Langkah 1 — Terapkan Flag SetGID pada Direktori `/opt/apps/erp_monolith/frontend` di VPS:**
  Jalankan perintah ini satu kali di server VPS:
  ```bash
  # Pastikan direktori tujuan frontend ada
  sudo mkdir -p /opt/apps/erp_monolith/frontend

  # Tetapkan kepemilikan deployer dan grup deployer
  sudo chown -R deployer:deployer /opt/apps/erp_monolith/frontend

  # Berikan izin baca-tulis-eksekusi untuk pemilik dan grup
  sudo chmod -R 775 /opt/apps/erp_monolith/frontend

  # Aktifkan SetGID agar file baru dari rsync otomatis mewarisi grup 'erp'
  sudo chmod g+s /opt/apps/erp_monolith/frontend
  ```

- **Langkah 2 — Penyesuaian Izin Pasca Sinkronisasi:**
  Jika tidak menggunakan SetGID, tambahkan perintah `chown` di blok script deployment GitHub Actions atau skrip transfer:
  ```bash
  sudo chown -R erp:erp /opt/apps/erp_monolith/frontend
  ```

#### 3. Konfigurasi Unit Systemd Frontend Khusus Hasil Build GitHub (`erp-monolith-frontend.service`)
Berikut adalah berkas unit `/etc/systemd/system/erp-monolith-frontend.service` (dengan alias/symlink `erp-frontend.service`) yang dirancang khusus untuk mengeksekusi hasil build yang diterima dari GitHub Actions di `/opt/apps/erp_monolith/frontend`:

```ini
[Unit]
Description=Enterprise ERP Monolith - SvelteKit 2 Frontend Service (GitHub Actions Build)
Documentation=https://github.com/jeruktutut2/arp-monolith
After=network.target erp-backend.service
Wants=erp-backend.service

[Service]
Type=exec
User=erp
Group=erp

# Direktori kerja frontend (lokasi hasil build dari GitHub Actions)
WorkingDirectory=/opt/apps/erp_monolith/frontend

# File environment produksi (memuat PORT=3000, HOST=127.0.0.1, ORIGIN, dll.)
EnvironmentFile=/etc/erp_monolith/erp-monolith-frontend.env

# Eksekusi entrypoint hasil build yang ditaruh di /opt/apps/erp_monolith/frontend
ExecStart=/usr/local/bin/bun /opt/apps/erp_monolith/frontend/index.js

# Siklus Hidup & Restart Otomatis
Restart=always
RestartSec=5s
KillMode=process
KillSignal=SIGTERM
TimeoutStopSec=20s

# Hardening Sandbox Linux
NoNewPrivileges=true
ProtectSystem=full
ProtectHome=true
PrivateTmp=true

# Batas Resource
LimitNOFILE=65536

# Integrasi Logging Journald
StandardOutput=journal
StandardError=journal
SyslogIdentifier=erp-frontend

[Install]
WantedBy=multi-user.target
```

> 📌 **Catatan Penting Penggunaan Bun di Systemd:**  
> Saat ini **BELUM BENAR**, karena file fisiknya masih berada di `/home/deployer/.bun/bin/bun`, sedangkan di `/usr/local/bin/bun` filenya belum ada. Jika dijalankan sekarang, systemd akan error: `No such file or directory`.
> 
> ⚠️ **Jangan Langsung Ganti ke `/home/deployer/.bun/bin/bun`!**  
> Anda mungkin berpikir untuk mengubahnya menjadi:
> ```ini
> ExecStart=/home/deployer/.bun/bin/bun ...   <-- JANGAN LAKUKAN INI
> ```
> Mengapa? Karena di file systemd Anda terdapat pengaturan keamanan ini:
> ```ini
> ProtectHome=true
> ```
> Fitur `ProtectHome=true` akan memblokir total akses ke folder `/home/`. Jika file bun ditaruh di `/home/deployer/`, systemd tidak akan bisa menjalankannya.
> 
> **Solusi yang Benar & Paling Aman:**  
> Salin (*copy*) file binary bun langsung ke direktori sistem `/usr/local/bin/`:
> ```bash
> # 1. Salin binary bun ke folder sistem
> sudo cp /home/deployer/.bun/bin/bun /usr/local/bin/bun
> # 2. Berikan izin eksekusi
> sudo chmod 755 /usr/local/bin/bun
> # 3. Verifikasi apakah sudah bisa dipanggil dari direktori sistem
> /usr/local/bin/bun -v
> ```
> 
> **Hasil Akhir:**  
> Setelah menjalankan 2 perintah di atas:
> - Perintah `which bun` nantinya akan mengenali `/usr/local/bin/bun`.
> - Baris konfigurasi Anda di systemd:
>   ```ini
>   ExecStart=/usr/local/bin/bun /opt/apps/erp_monolith/frontend/index.js
>   ```
>   sudah **100% BENAR**, aman dari blokiran `ProtectHome`, dan siap digunakan!
> 
> *(Catatan: Jika Anda memilih runtime Node.js alih-alih Bun, gunakan `ExecStart=/usr/bin/node /opt/apps/erp_monolith/frontend/index.js`)*.

#### 4. Alur Restart & Health Check Pasca Penerimaan Build
Setelah step transfer Rsync selesai memindahkan file build ke `/opt/apps/erp_monolith/frontend/`, skrip GitHub Actions akan mengeksekusi perintah reload & restart berikut melalui SSH:

```bash
# 1. Pastikan berkas entrypoint index.js dapat dibaca
test -f /opt/apps/erp_monolith/frontend/index.js || { echo "❌ /opt/apps/erp_monolith/frontend/index.js tidak ditemukan!"; exit 1; }

# 2. Restart unit service frontend secara graceful
sudo systemctl restart erp-frontend.service

# 3. Verifikasi ketersediaan service (Health Check internal port 3000)
sleep 3
curl -fs http://127.0.0.1:3000 > /dev/null || {
  echo "❌ Frontend gagal merespons setelah restart!"
  sudo journalctl -u erp-frontend -n 50 --no-pager
  exit 1
}

echo "✅ Frontend berhasil diperbarui dari build GitHub Actions di /opt/apps/erp_monolith/frontend dan aktif."
```

### F. Langkah Persiapan Server Sebelum Deployment Pertama

> ⚠️ **PERINGATAN PENTING:**  
> Karena hasil build frontend-nya belum ada, **JANGAN jalankan `systemctl start` sekarang**, karena pasti akan langsung gagal (*error: file not found / 203 EXEC*).

Berikut langkah-langkah yang harus dilakukan untuk mempersiapkan server sebelum deployment pertama:

#### 1. Pastikan User `erp` Sudah Ada di VPS
Pada service Anda tertulis `User=erp` dan `Group=erp`. Pastikan user sistem `erp` sudah dibuat:

```bash
# Cek apakah user erp sudah ada
id erp

# Jika belum ada, buat user sistem khusus (tanpa akses login shell):
sudo useradd -r -s /bin/false erp
```

#### 2. Buat Folder Tujuan & Atur Izin Akses
Folder `/opt/apps/erp_monolith/frontend` harus ada dan bisa ditulisi oleh user `deployer` (yang melakukan rsync dari GitHub Actions), serta bisa dibaca oleh user `erp`:

```bash
# Buat direktori tujuan
sudo mkdir -p /opt/apps/erp_monolith/frontend

# Berikan kepemilikan ke user deployer, dengan grup deployer
sudo chown -R deployer:deployer /opt/apps/erp_monolith/frontend

# Beri permission: deployer bisa tulis, user erp bisa membaca dan mengeksekusi
sudo chmod -R 755 /opt/apps/erp_monolith/frontend
```

#### 3. Pastikan File Environment Sudah Dibuat
Service membutuhkan `EnvironmentFile=/etc/erp_monolith/erp-monolith-frontend.env`. Jika file ini belum ada, systemd akan menolak jalan:

```bash
# Buat folder dan filenya jika belum ada
sudo mkdir -p /etc/erp_monolith
sudo nano /etc/erp_monolith/erp-monolith-frontend.env
```

Isi variabel minimal (sesuai kebutuhan SvelteKit), misalnya:

```env
NODE_ENV=production
PORT=3000
HOST=127.0.0.1
ORIGIN=https://erp.domainanda.com
```

Kunci hak aksesnya agar user `deployer` / `erp` bisa membaca file ini:

```bash
sudo chown root:deployer /etc/erp_monolith/erp-monolith-frontend.env
sudo chmod 640 /etc/erp_monolith/erp-monolith-frontend.env
```

#### 4. Cek Path Bun atau Node.js
Pastikan binary runtime tersedia di path sistem:

- **Jika Menggunakan Bun (Sesuai Konfigurasi Utama):**
  Pastikan binary bun sudah disalin ke `/usr/local/bin/bun` agar aman dari proteksi `ProtectHome=true`:
  ```bash
  # 1. Salin binary bun ke folder sistem
  sudo cp /home/deployer/.bun/bin/bun /usr/local/bin/bun
  # 2. Berikan izin eksekusi
  sudo chmod 755 /usr/local/bin/bun
  # 3. Verifikasi ketersediaan bun di sistem
  /usr/local/bin/bun -v
  ```
- **Jika Menggunakan Node.js:**
  ```bash
  which node
  ```
  Jika outputnya `/usr/bin/node`, pastikan `ExecStart=/usr/bin/node ...` pada berkas service.

#### 5. Daftarkan Service ke Systemd (daemon-reload & enable)
Jangan di-start sekarang, cukup reload dan enable agar service terdaftar di sistem:

```bash
# Muat ulang konfigurasi systemd
sudo systemctl daemon-reload

# Aktifkan agar service otomatis jalan tiap kali VPS restart
sudo systemctl enable erp-monolith-frontend.service
```

> 💡 **Catatan Unit Systemd:**  
> Pastikan meng-enable file unit `.service` (`erp-monolith-frontend.service`), bukan file environment (`.env`). File konfigurasi environment (`/etc/erp_monolith/erp-monolith-frontend.env`) dimuat secara otomatis oleh systemd melalui direktif `EnvironmentFile` di dalam berkas unit.

#### 6. Cek File Workflow GitHub Actions Anda
Perhatikan di file workflow GitHub Actions Anda (`Untitled-1` / `deploy.yml`), perintah restart saat ini masih dikomentari (`#`):

```yaml
# sudo systemctl restart erp-monolith-frontend.service
```

Buka komentar (*uncomment*) baris tersebut menjadi:

```yaml
sudo systemctl restart erp-monolith-frontend.service
sudo systemctl is-active erp-monolith-frontend.service
```

#### 7. Jalankan Deployment Pertama! 🚀
Sekarang server sudah siap menerima file:

1. Lakukan `git push` ke repository Anda untuk men-trigger GitHub Actions.
2. GitHub Actions akan:
   - Menjalankan `build-frontend`.
   - Mengirim folder hasil build ke `/opt/apps/erp_monolith/frontend/` melalui rsync.
   - Menjalankan `sudo systemctl restart erp-monolith-frontend.service`.
3. Setelah file masuk, barulah systemd otomatis menyalakan aplikasi frontend SvelteKit Anda!

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
sudo systemctl enable erp-monolith-frontend.service
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

# 4. Kompilasi / Penyiapan Frontend SvelteKit
# Catatan: Jika menggunakan GitHub Actions, hasil build langsung disinkronkan ke /opt/apps/erp_monolith/frontend
if [ -f "/opt/apps/erp_monolith/frontend/index.js" ]; then
    echo "🎨 Build frontend di /opt/apps/erp_monolith/frontend terdeteksi, memastikan hak akses..."
    chown -R erp:erp /opt/apps/erp_monolith/frontend
elif [ -f "${FRONTEND_DIR}/build/index.js" ]; then
    echo "🎨 Build frontend lokal terdeteksi, menyalin ke /opt/apps/erp_monolith/frontend..."
    mkdir -p /opt/apps/erp_monolith/frontend
    cp -r "${FRONTEND_DIR}/build/"* /opt/apps/erp_monolith/frontend/
    chown -R erp:erp /opt/apps/erp_monolith/frontend
else
    echo "🎨 Melakukan kompilasi frontend SvelteKit secara lokal..."
    cd "${FRONTEND_DIR}"
    bun install --frozen-lockfile
    bun run build
    mkdir -p /opt/apps/erp_monolith/frontend
    cp -r build/* /opt/apps/erp_monolith/frontend/
    chown -R erp:erp /opt/apps/erp_monolith/frontend
fi

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
  1. Pastikan pada `/etc/erp_monolith/erp-monolith-frontend.env`:
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
  1. Pada SvelteKit (`/etc/erp_monolith/erp-monolith-frontend.env`):
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
- **Penyebab:** Path binary Node.js atau Bun di direktori `/usr/local/bin` belum ada/tidak ditemukan oleh systemd, atau binary Bun ditaruh di `/home/deployer/.bun/bin/bun` yang diblokir oleh `ProtectHome=true`.
- **Solusi:** 
  1. Untuk Bun, salin binary ke folder sistem:
     ```bash
     sudo cp /home/deployer/.bun/bin/bun /usr/local/bin/bun
     sudo chmod 755 /usr/local/bin/bun
     ```
  2. Periksa lokasi absolut binary dengan:
     ```bash
     which node
     which bun
     ```
     Gunakan path absolut `/usr/local/bin/bun` pada direktif `ExecStart` di berkas `.service`.

---

### 5. Port Permission Denied: `bind: permission denied`
- **Penyebab:** Service dijalankan oleh user non-root (`erp`) dan mencoba mengikat (*bind*) ke port istimewa (*privileged port* < 1024) seperti port 80 atau 443 secara langsung.
- **Solusi:** Jangan bind Go atau SvelteKit langsung ke port 80/443. Biarkan backend mendengarkan port 8080 dan frontend port 3000, lalu gunakan Nginx atau izinkan capability:
  ```bash
  sudo setcap 'cap_net_bind_service=+ep' /opt/dev/erp_monolith/backend/bin/erp-backend
  ```
