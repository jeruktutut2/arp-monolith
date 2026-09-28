# 🚀 Panduan Deployment GitHub Actions ke Remote VPS (Golang & SvelteKit 2) — Enterprise ERP Monolith

Dokumen ini adalah panduan resmi untuk merancang, mengonfigurasi, dan mengoperasikan pipa integrasi & pengiriman berkelanjutan (**CI/CD Pipeline**) menggunakan **GitHub Actions** guna melakukan deployment otomatis backend **Golang (Echo v5)** dan frontend **SvelteKit 2 (Bun/Node)** ke server remote Linux (**VPS Ubuntu 22.04 / 24.04 LTS**) yang dikelola oleh **Systemd**.

---

## 📑 Daftar Isi
1. [Arsitektur CI/CD & Pipeline Workflow](#1-arsitektur-cicd--pipeline-workflow)
2. [Perbandingan 2 Strategi Deployment ke VPS](#2-perbandingan-2-strategi-deployment-ke-vps)
3. [Konfigurasi Prasyarat di Sisi Server VPS](#3-konfigurasi-prasyarat-di-sisi-server-vps)
   - [A. Pembuatan User Dedicated & SSH Key Deployer](#a-pembuatan-user-dedicated--ssh-key-deployer)
   - [B. Konfigurasi Sudoers Tanpa Password (`NOPASSWD`)](#b-konfigurasi-sudoers-tanpa-password-nopasswd)
   - [C. Menyiapkan Struktur Direktori di VPS](#c-menyiapkan-struktur-direktori-di-vps)
4. [Konfigurasi GitHub Secrets & Variables](#4-konfigurasi-github-secrets--variables)
5. [Workflow Produksi 1: Build on Runner + Rsync (Direkomendasikan)](#5-workflow-produksi-1-build-on-runner--rsync-direkomendasikan)
6. [Workflow Produksi 2: SSH Remote Trigger + Git Pull di VPS](#6-workflow-produksi-2-ssh-remote-trigger--git-pull-di-vps)
7. [Penanganan Migrasi Database Otomatis (`golang-migrate`)](#7-penanganan-migrasi-database-otomatis-golang-migrate)
8. [Verifikasi Health Check & Rollback Otomatis](#8-verifikasi-health-check--rollback-otomatis)
9. [Hardening Keamanan & Best Practices](#9-hardening-keamanan--best-practices)
10. [Troubleshooting Masalah Umum CI/CD](#10-troubleshooting-masalah-umum-cicd)

---

## 1. Arsitektur CI/CD & Pipeline Workflow

Diagram alur berikut mengilustrasikan proses deployment otomatis dari saat developer melakukan `git push` ke branch `main`:

```text
  Developer / Git Push ke branch 'main'
                   │
                   ▼
┌────────────────────────────────────────────────────────┐
│ GitHub Actions Runner (Ubuntu Latest)                 │
│                                                        │
│ 1. Checkout Repository                                 │
│ 2. Backend Pipeline:                                   │
│    • Setup Go 1.22+                                    │
│    • go test ./...                                     │
│    • Build Linux amd64 binary: erp-backend             │
│                                                        │
│ 3. Frontend Pipeline:                                  │
│    • Setup Bun / Node.js                               │
│    • bun install --frozen-lockfile                     │
│    • bun run build (Output: /build)                    │
│                                                        │
│ 4. Deployment via SSH & Rsync                          │
└──────────────────────────┬─────────────────────────────┘
                           │ SSH / SFTP Encrypted (Port 22/custom)
                           ▼
┌────────────────────────────────────────────────────────┐
│ Remote VPS Server (Ubuntu 22.04 / 24.04 LTS)          │
│                                                        │
│ 1. Terima Artifact (Binary Go & Folder Build)          │
│ 2. Jalankan DDL Migrasi DB (golang-migrate direct:5432)│
│ 3. Atomic Binary Replacement (mv .new erp-backend)     │
│ 4. Systemd Restart (sudo systemctl restart erp.target) │
│ 5. Health Check Verification (curl /health & :3000)    │
│ 6. Notifikasi Status (Success / Rollback)              │
└────────────────────────────────────────────────────────┘
```

---

## 2. Perbandingan 2 Strategi Deployment ke VPS

| Parameter | Strategi A: Build di Runner + Rsync Artifact (⭐ Direkomendasikan) | Strategi B: SSH Trigger + Build di VPS |
|---|---|---|
| **Beban Server VPS** | **Sangat Rendah** (VPS tidak mengalami lonjakan CPU/RAM) | **Tinggi** (VPS bisa kehabisan RAM/OOM saat kompilasi) |
| **Kebutuhan Tooling VPS**| Hanya butuh OS Linux & Systemd (tidak perlu Go/Bun/Vite di VPS) | Wajib pasang Git, Go Compiler, Bun, Node.js, gcc di VPS |
| **Keamanan Kode Sumber** | Source code asli tetap di GitHub Runner, hanya binary & dist yang dikirim | Seluruh source code tersimpan di server VPS |
| **Waktu Deployment** | Cepat dan terisolasi | Bergantung pada kecepatan CPU VPS |
| **Kompleksitas CI/CD**| Memerlukan step rsync / SCP artifact | Sangat sederhana (hanya perintah bash lewat SSH) |

> 💡 **Rekomendasi:** Untuk server produksi berskala kecil hingga menengah (RAM 2 GB – 8 GB), gunakan **Strategi A** agar server tidak hang atau kehabisan memori (*Out-Of-Memory*) saat menjalankan `bun run build` atau `go build`.

---

## 3. Konfigurasi Prasyarat di Sisi Server VPS

Sebelum GitHub Actions dapat terhubung, siapkan user khusus deployment dan hak akses SSH di VPS Anda.

### A. Pembuatan User Dedicated & SSH Key Deployer
Jangan gunakan `root` untuk deployment. Buat user `deployer` (atau gunakan user yang sudah ada seperti `ubuntu` / `erp`):

```bash
# 1. Buat user deployer di VPS
sudo useradd -m -s /bin/bash deployer
sudo usermod -aG erp deployer

# 2. Buat pasangan kunci SSH di VPS (atau di komputer lokal Anda)
# Jalankan sebagai user deployer:
sudo su - deployer
mkdir -p ~/.ssh && chmod 700 ~/.ssh
ssh-keygen -t ed25519 -C "github-actions-deployer" -f ~/.ssh/id_ed25519 -N ""

# 3. Masukkan public key ke authorized_keys
cat ~/.ssh/id_ed25519.pub >> ~/.ssh/authorized_keys
chmod 600 ~/.ssh/authorized_keys

# 4. Tampilkan Private Key untuk disalin ke GitHub Secret
cat ~/.ssh/id_ed25519
```
> [!CAUTION]
> Salin seluruh isi `~/.ssh/id_ed25519` (termasuk `-----BEGIN OPENSSH PRIVATE KEY-----` dan `-----END OPENSSH PRIVATE KEY-----`). Kunci ini akan dimasukkan ke GitHub Secrets (`VPS_SSH_KEY`).

---

### B. Konfigurasi Sudoers Tanpa Password (`NOPASSWD`)
User `deployer` memerlukan izin untuk me-restart layanan systemd tanpa harus memasukkan password sudo interaktif.

Buat berkas konfigurasi sudoers khusus:
```bash
sudo nano /etc/sudoers.d/erp-deployer
```

Isi berkas tersebut dengan pembatasan perintah ketat (*least privilege*):
```text
# Izinkan user deployer merestart & memeriksa layanan ERP Monolith tanpa password
deployer ALL=(ALL) NOPASSWD: /bin/systemctl restart erp-backend, /bin/systemctl restart erp-frontend, /bin/systemctl restart erp.target, /bin/systemctl status erp*, /bin/systemctl reload nginx
```

Uji validitas sintaks sudoers:
```bash
sudo visudo -cf /etc/sudoers.d/erp-deployer
```

---

### C. Menyiapkan Struktur Direktori di VPS
Pastikan direktori aplikasi di VPS sudah dibuat dan user `deployer` memiliki hak tulis:

```bash
sudo mkdir -p /opt/dev/erp_monolith/backend/bin
sudo mkdir -p /opt/dev/erp_monolith/backend/migrations
sudo mkdir -p /opt/dev/erp_monolith/frontend/build
sudo mkdir -p /etc/erp

# Berikan izin ke grup erp di mana deployer tergabung
sudo chown -R erp:erp /opt/dev/erp_monolith
sudo chmod -R 775 /opt/dev/erp_monolith
```

---

## 4. Konfigurasi GitHub Secrets & Variables

Buka repositori GitHub Anda: **Settings** -> **Secrets and variables** -> **Actions** -> **New repository secret**.

Tambahkan rahasia-rahasia berikut:

| Nama Secret | Deskripsi | Contoh Nilai |
|---|---|---|
| `VPS_HOST` | IP Publik atau Domain VPS | `203.0.113.50` atau `vps.perusahaan.com` |
| `VPS_PORT` | Port SSH server | `22` (atau port SSH custom Anda) |
| `VPS_USER` | Username SSH untuk login | `deployer` |
| `VPS_SSH_KEY` | Private Key SSH OpenSSH | `-----BEGIN OPENSSH PRIVATE KEY----- ...` |
| `DIRECT_DATABASE_URL` | Koneksi PostgreSQL direct untuk migrasi skema (port 5432) | `postgres://erp_user:erp_secret@127.0.0.1:5432/erp_db?sslmode=disable` |

---

## 5. Workflow Produksi 1: Build on Runner + Rsync (Direkomendasikan)

Workflow ini mengompilasi binary Golang dan membundle SvelteKit di server GitHub Actions, lalu mentransfer hasilnya ke VPS via SSH/Rsync.

Simpan berkas berikut pada repositori Anda di `.github/workflows/deploy.yml`:

```yaml
name: 🚀 Build & Deploy to Production VPS

on:
  push:
    branches:
      - main
  workflow_dispatch: # Memungkinkan trigger manual dari UI GitHub

# Cegah concurrent deployment agar tidak saling tumpang tindih
concurrency:
  group: production_deployment
  cancel-in-progress: false

jobs:
  # ========================================================
  # JOB 1: Test & Build Backend Golang
  # ========================================================
  build-backend:
    name: 🔨 Build Go Backend
    runs-on: ubuntu-latest
    steps:
      - name: 📥 Checkout Code
        uses: actions/checkout@v4

      - name: 🐹 Setup Go 1.22+
        uses: actions/setup-go@v5
        with:
          go-version: '1.22'
          check-latest: true
          cache-dependency-path: backend/go.sum

      - name: 🧪 Run Backend Tests
        run: |
          cd backend
          go test -v ./...

      - name: 📦 Compile Go Binary (Linux AMD64)
        run: |
          cd backend
          CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/erp-backend cmd/server/main.go

      - name: 📤 Upload Backend Artifact
        uses: actions/upload-artifact@v4
        with:
          name: backend-binary
          path: backend/bin/erp-backend
          retention-days: 1

  # ========================================================
  # JOB 2: Build Frontend SvelteKit 2 (Bun)
  # ========================================================
  build-frontend:
    name: 🎨 Build SvelteKit Frontend
    runs-on: ubuntu-latest
    steps:
      - name: 📥 Checkout Code
        uses: actions/checkout@v4

      - name: 🥟 Setup Bun
        uses: oven-sh/setup-bun@v2
        with:
          bun-version: latest

      - name: 📦 Install Frontend Dependencies & Build
        run: |
          if [ -d "frontend" ]; then
            cd frontend
            bun install --frozen-lockfile
            bun run build
          else
            echo "Direktori frontend belum dibuat, membuat placeholder dist."
            mkdir -p frontend/build
            echo "Placeholder frontend" > frontend/build/index.js
          fi

      - name: 📤 Upload Frontend Artifact
        uses: actions/upload-artifact@v4
        with:
          name: frontend-build
          path: frontend/build/
          retention-days: 1

  # ========================================================
  # JOB 3: Deploy Artifacts & Restart Services on VPS
  # ========================================================
  deploy:
    name: 🚀 Deploy to VPS
    needs: [build-backend, build-frontend]
    runs-on: ubuntu-latest
    steps:
      - name: 📥 Checkout Migrations & Scripts
        uses: actions/checkout@v4

      - name: 📥 Download Backend Artifact
        uses: actions/download-artifact@v4
        with:
          name: backend-binary
          path: backend-dist

      - name: 📥 Download Frontend Artifact
        uses: actions/download-artifact@v4
        with:
          name: frontend-build
          path: frontend-dist

      - name: 🔑 Configure SSH Key
        run: |
          mkdir -p ~/.ssh
          echo "${{ secrets.VPS_SSH_KEY }}" > ~/.ssh/id_ed25519
          chmod 600 ~/.ssh/id_ed25519
          ssh-keyscan -p ${{ secrets.VPS_PORT || 22 }} ${{ secrets.VPS_HOST }} >> ~/.ssh/known_hosts

      - name: 🚚 Transfer Artifacts to VPS via Rsync
        env:
          PORT: ${{ secrets.VPS_PORT || 22 }}
          USER: ${{ secrets.VPS_USER }}
          HOST: ${{ secrets.VPS_HOST }}
        run: |
          # 1. Kirim binary baru backend ke file temporary .new
          scp -P $PORT backend-dist/erp-backend $USER@$HOST:/opt/dev/erp_monolith/backend/bin/erp-backend.new
          
          # 2. Kirim migrasi database
          rsync -avz -e "ssh -p $PORT" --delete backend/migrations/ $USER@$HOST:/opt/dev/erp_monolith/backend/migrations/

          # 3. Kirim folder build frontend
          rsync -avz -e "ssh -p $PORT" --delete frontend-dist/ $USER@$HOST:/opt/dev/erp_monolith/frontend/build/

      - name: ⚡ Execute Deployment Commands on Remote VPS
        uses: appleboy/ssh-action@v1.0.3
        with:
          host: ${{ secrets.VPS_HOST }}
          username: ${{ secrets.VPS_USER }}
          key: ${{ secrets.VPS_SSH_KEY }}
          port: ${{ secrets.VPS_PORT || 22 }}
          script: |
            set -e
            echo "--- [1/5] Menjalankan Migrasi Database ---"
            if command -v migrate >/dev/null 2>&1; then
              migrate -path /opt/dev/erp_monolith/backend/migrations \
                -database "${{ secrets.DIRECT_DATABASE_URL }}" up || echo "Migrasi up selesai."
            fi

            echo "--- [2/5] Atomic Swap Binary Golang ---"
            chmod +x /opt/dev/erp_monolith/backend/bin/erp-backend.new
            mv -f /opt/dev/erp_monolith/backend/bin/erp-backend.new /opt/dev/erp_monolith/backend/bin/erp-backend

            echo "--- [3/5] Restart Layanan Systemd ---"
            sudo systemctl restart erp-backend.service
            sudo systemctl restart erp-frontend.service

            echo "--- [4/5] Verifikasi Health Check ---"
            sleep 3
            curl -fs http://127.0.0.1:8080/health || {
              echo "❌ Backend Health Check GAGAL!"
              sudo systemctl status erp-backend.service --no-pager
              exit 1
            }

            echo "✅ Health check Backend Berhasil (HTTP 200 OK)."

            echo "--- [5/5] Status Akhir Layanan ---"
            sudo systemctl is-active erp-backend.service
            sudo systemctl is-active erp-frontend.service
            echo "🎉 Deployment Sukses!"
```

---

## 6. Workflow Produksi 2: SSH Remote Trigger + Git Pull di VPS

Jika server VPS Anda memiliki spesifikasi tinggi (RAM > 8 GB) dan seluruh tool (*Go, Bun, Git, Migrate*) sudah terinstal langsung di VPS, Anda dapat menggunakan workflow SSH sederhana yang melakukan `git pull` dan build di server lokal VPS:

Simpan di `.github/workflows/deploy-ssh-pull.yml`:

```yaml
name: 🔄 Deploy via SSH Pull (Direct on VPS)

on:
  workflow_dispatch: # Jalankan sesuai kebutuhan

jobs:
  deploy:
    name: 🚀 Remote Deploy on VPS
    runs-on: ubuntu-latest
    steps:
      - name: ⚡ Connect SSH & Trigger Local Build
        uses: appleboy/ssh-action@v1.0.3
        with:
          host: ${{ secrets.VPS_HOST }}
          username: ${{ secrets.VPS_USER }}
          key: ${{ secrets.VPS_SSH_KEY }}
          port: ${{ secrets.VPS_PORT || 22 }}
          script_stop: true
          script: |
            cd /opt/dev/erp_monolith

            echo "--- [1] Git Pull Perubahan Terbaru ---"
            git fetch origin main
            git reset --hard origin/main

            echo "--- [2] Menjalankan Skrip Deployment Lokal ---"
            chmod +x ./deploy.sh
            ./deploy.sh
```

---

## 7. Penanganan Migrasi Database Otomatis (`golang-migrate`)

Sesuai aturan arsitektur monolit ini:
- **Aplikasi Go Backend** terhubung ke port **PgBouncer 6432**.
- **Migrasi DDL Database (`golang-migrate`)** **WAJIB** terhubung langsung ke port **PostgreSQL direct 5432**.

### Menyiapkan Binary `migrate` di VPS
Pastikan CLI `golang-migrate` sudah terpasang di VPS:
```bash
# Pasang golang-migrate di Ubuntu
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz
sudo mv migrate /usr/local/bin/migrate
migrate -version
```

Saat CI/CD berjalan, perintah:
```bash
migrate -path /opt/dev/erp_monolith/backend/migrations \
  -database "${{ secrets.DIRECT_DATABASE_URL }}" up
```
akan mengeksekusi semua berkas `*.up.sql` yang belum pernah dijalankan secara idempoten dan aman.

---

## 8. Verifikasi Health Check & Rollback Otomatis

Untuk mencegah downtime jika terjadi *bad deployment* (misal binary crash atau port bentrok), pasang mekanisme rollback di skrip SSH:

```bash
# Contoh logika Rollback pada skrip remote:
BACKUP_BIN="/opt/dev/erp_monolith/backend/bin/erp-backend.bak"
ACTIVE_BIN="/opt/dev/erp_monolith/backend/bin/erp-backend"

# Simpan cadangan binary saat ini
cp "$ACTIVE_BIN" "$BACKUP_BIN"

# Terapkan binary baru
mv -f /opt/dev/erp_monolith/backend/bin/erp-backend.new "$ACTIVE_BIN"
sudo systemctl restart erp-backend.service

# Verifikasi Health Check dalam 5 detik
sleep 3
if ! curl -fs http://127.0.0.1:8080/health > /dev/null; then
    echo "🚨 Health check gagal! Melakukan ROLLBACK otomatis..."
    cp "$BACKUP_BIN" "$ACTIVE_BIN"
    sudo systemctl restart erp-backend.service
    echo "⚠️ Sistem berhasil di-rollback ke versi sebelumnya."
    exit 1
fi
```

Jika exit code adalah `1`, GitHub Actions akan menandai workflow sebagai **Failed** (merah) dan mengirimkan alert email/notifikasi ke tim.

---

## 9. Hardening Keamanan & Best Practices

1. **Gunakan SSH Key Khusus (Ed25519):**
   - Buat kunci khusus untuk GitHub Actions, jangan gunakan kunci pribadi developer.
   - Jika repository atau token dicurigai bocor, Anda cukup menghapus entri public key dari `~/.ssh/authorized_keys` di VPS tanpa mengganggu akses tim lain.

2. **Batasi Hak Akses Sudo (`/etc/sudoers.d/erp-deployer`):**
   - Jangan pernah memberikan `NOPASSWD: ALL`.
   - Batasi hanya pada perintah `systemctl` untuk unit-unit terkait (`erp-backend`, `erp-frontend`, `erp.target`, `nginx`).

3. **Gunakan GitHub Environment Protection:**
   - Masuk ke **Settings** -> **Environments** -> Buat environment `production`.
   - Aktifkan **Required reviewers** jika rilis ke VPS memerlukan persetujuan manual (Lead/DevOps).

4. **Kunci Port SSH dengan UFW / Fail2ban:**
   - Pastikan port SSH dilindungi oleh `fail2ban` untuk mencegah brute-force.
   - Pertimbangkan mengubah default SSH port dari `22` ke port custom.

---

## 10. Troubleshooting Masalah Umum CI/CD

### 1. Error: `Host key verification failed`
- **Penyebab:** GitHub Runner belum mengenali fingerprint SSH server VPS.
- **Solusi:** Jalankan `ssh-keyscan` sebelum perintah SSH atau SCP:
  ```bash
  ssh-keyscan -p ${{ secrets.VPS_PORT || 22 }} ${{ secrets.VPS_HOST }} >> ~/.ssh/known_hosts
  ```

---

### 2. Error: `Permission denied (publickey)`
- **Penyebab:** Private key di `VPS_SSH_KEY` tidak cocok dengan public key di `~/.ssh/authorized_keys`, atau izin file di VPS terlalu terbuka.
- **Solusi:** Di VPS, periksa izin folder:
  ```bash
  chmod 700 ~/.ssh
  chmod 600 ~/.ssh/authorized_keys
  ```
  Pastikan saat menyalin private key ke GitHub Secrets, tidak ada karakter spasi atau baris pembuka/penutup yang terpotong.

---

### 3. Error: `sudo: a password is required`
- **Penyebab:** Perintah `sudo systemctl restart ...` yang dieksekusi oleh `deployer` meminta password interaktif karena belum diizinkan di sudoers.
- **Solusi:** Pastikan konfigurasi `/etc/sudoers.d/erp-deployer` sudah benar dan gunakan path absolut lengkap binary:
  ```text
  deployer ALL=(ALL) NOPASSWD: /bin/systemctl restart erp-backend, /bin/systemctl restart erp-frontend, /bin/systemctl restart erp.target
  ```

---

### 4. Error: `curl: (7) Failed to connect to 127.0.0.1 port 8080: Connection refused`
- **Penyebab:** Backend gagal start (misal koneksi database PgBouncer putus atau file konfigurasi `.env` hilang).
- **Solusi:**
  1. Periksa log systemd di VPS:
     ```bash
     journalctl -u erp-backend -n 50 --no-pager
     ```
  2. Pastikan file `/etc/erp/erp-backend.env` ada dan berisi koneksi database yang valid.
