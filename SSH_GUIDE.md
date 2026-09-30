# 🔑 Panduan Lengkap SSH & User Deployer — Enterprise ERP Monolith

Dokumen ini adalah panduan resmi dan terstandarisasi untuk mengonfigurasi **User Deployer**, **Konektivitas SSH Dua Arah (Server Ubuntu ⇄ GitHub)**, **Manajemen Berkas `~/.ssh/config`**, serta **Hardening Keamanan Server** pada infrastruktur Enterprise ERP Monolith berbasis Ubuntu 22.04 / 24.04 LTS.

---

## 📑 Daftar Isi
1. [Arsitektur Konektivitas & Model Keamanan SSH](#1-arsitektur-konektivitas--model-keamanan-ssh)
2. [Bagian 1: Pembuatan & Konfigurasi User `deployer` di Ubuntu](#2-bagian-1-pembuatan--konfigurasi-user-deployer-di-ubuntu)
   - [A. Pembuatan Akun & Penugasan Grup](#a-pembuatan-akun--penugasan-grup)
   - [B. Konfigurasi Sudoers Tanpa Password Terbatas (`NOPASSWD`)](#b-konfigurasi-sudoers-tanpa-password-terbatas-nopasswd)
   - [C. Hak Akses Direktori Proyek & SetGID](#c-hak-akses-direktori-proyek--setgid)
   - [D. Inisialisasi Direktori `~/.ssh` dengan Izin Presisi](#d-inisialisasi-direktori-ssh-dengan-izin-presisi)
3. [Bagian 2: SSH dari Server Ubuntu ke GitHub (Push & Pull)](#3-bagian-2-ssh-dari-server-ubuntu-ke-github-push--pull)
   - [A. Skenario & Use Case VPS ke GitHub](#a-skenario--use-case-vps-ke-github)
   - [B. Pembuatan Pasangan Kunci SSH Ed25519 di VPS](#b-pembuatan-pasangan-kunci-ssh-ed25519-di-vps)
   - [C. Pendaftaran Kunci di GitHub: Deploy Key vs Personal/Machine Key](#c-pendaftaran-kunci-di-github-deploy-key-vs-personalmachine-key)
   - [D. Verifikasi Host Key & Tes Koneksi SSH](#d-verifikasi-host-key--tes-koneksi-ssh)
   - [E. Konfigurasi Git & Migrasi Remote URL ke SSH](#e-konfigurasi-git--migrasi-remote-url-ke-ssh)
   - [F. Pengujian Push & Pull Langsung dari Server](#f-pengujian-push--pull-langsung-dari-server)
4. [Bagian 3: Pembuatan & Manajemen Berkas `~/.ssh/config`](#4-bagian-3-pembuatan--manajemen-berkas-sshconfig)
   - [A. Manfaat & Parameter Kunci `~/.ssh/config`](#a-manfaat--parameter-kunci-sshconfig)
   - [B. Konfigurasi `~/.ssh/config` di Sisi Server VPS (`deployer`)](#b-konfigurasi-sshconfig-di-sisi-server-vps-deployer)
   - [C. Konfigurasi Multi-Akun / Multi-Repo GitHub](#c-konfigurasi-multi-akun--multi-repo-github)
   - [D. Konfigurasi `~/.ssh/config` di Komputer Lokal Developer](#d-konfigurasi-sshconfig-di-komputer-lokal-developer)
   - [E. Izin Akses Wajib File Config](#e-izin-akses-wajib-file-config)
5. [Bagian 4: SSH dari GitHub Actions Deploy ke Server Ubuntu (CI/CD)](#5-bagian-4-ssh-dari-github-actions-deploy-ke-server-ubuntu-cicd)
   - [A. Arsitektur Alur Deployment Runner ke VPS](#a-arsitektur-alur-deployment-runner-ke-vps)
   - [B. Pembuatan Pasangan Kunci SSH Khusus Deployment](#b-pembuatan-pasangan-kunci-ssh-khusus-deployment)
   - [C. Pemasangan Public Key di VPS (`authorized_keys`) & Opsi Restriksi](#c-pemasangan-public-key-di-vps-authorized_keys--opsi-restriksi)
   - [D. Konfigurasi GitHub Repository Secrets](#d-konfigurasi-github-repository-secrets)
   - [E. Integrasi di GitHub Actions Workflow (`rsync` & `ssh-action`)](#e-integrasi-di-github-actions-workflow-rsync--ssh-action)
6. [Bagian 5: Hardening Keamanan OpenSSH Server (`sshd_config`)](#6-bagian-5-hardening-keamanan-openssh-server-sshd_config)
7. [Bagian 6: Cheatsheet Perintah Cepat & Troubleshooting](#7-bagian-6-cheatsheet-perintah-cepat--troubleshooting)

---

## 1. Arsitektur Konektivitas & Model Keamanan SSH

Dalam siklus hidup operasional ERP Monolith, terdapat **3 entitas pengguna** dan **3 arah alur koneksi SSH** yang saling melengkapi:

### Pemisahan Peran User Sistem di Server Ubuntu:
1. **`root`**: Administrator sistem level tertinggi. Akses SSH langsung dinonaktifkan (`PermitRootLogin no`) demi keamanan.
2. **`erp`**: Akun *system service* (`/usr/sbin/nologin`). Digunakan secara eksklusif oleh **Systemd** untuk menjalankan binary Go (`erp-backend`) dan proses SSR frontend SvelteKit (`erp-frontend`). User ini tidak memiliki password dan tidak bisa login via SSH.
3. **`deployer`**: Akun operasional interaktif bertindak sebagai *Deployment Agent*. Memiliki izin sudo terbatas (*least privilege*) untuk restart service, mengelola file di `/opt/dev/erp_monolith`, dan memegang otentikasi SSH.

### Topologi Tiga Arah Alur SSH:

```text
 ┌────────────────────────┐                   ┌────────────────────────┐
 │   Developer Workstation│                   │     GitHub Actions     │
 │    (Komputer Lokal)    │                   │   Runner (Cloud CI/CD) │
 └───────────┬────────────┘                   └───────────┬────────────┘
             │                                            │
             │ [Arah 1: Maintenance]                      │ [Arah 2: CI/CD Deploy]
             │ ssh erp-prod                               │ scp & ssh via Secrets
             ▼                                            ▼
 ┌─────────────────────────────────────────────────────────────────────┐
 │                         Server VPS Ubuntu                           │
 │                                                                     │
 │   User: deployer (Home: /home/deployer)                             │
 │   • /home/deployer/.ssh/authorized_keys <── Public Key Dev & CI/CD  │
 │   • /home/deployer/.ssh/config          <── Alias Host & Parameter  │
 │   • /home/deployer/.ssh/id_ed25519_github <── Private Key ke GitHub │
 │                                                                     │
 │   User: erp (System Service, nologin)                               │
 │   • Menjalankan erp-backend.service & erp-frontend.service          │
 └──────────────────────────────────┬──────────────────────────────────┘
                                    │
                                    │ [Arah 3: Git Push / Pull / Fetch]
                                    │ Otentikasi SSH Key
                                    ▼
                         ┌─────────────────────┐
                         │ GitHub Repositories │
                         │   (github.com)      │
                         └─────────────────────┘
```

### Tabel Matriks Pasangan Kunci SSH:

| Alur / Arah | SSH Client (Asal) | SSH Server (Tujuan) | Lokasi Private Key | Lokasi Public Key | Fungsi Utama |
|---|---|---|---|---|---|
| **Arah 1: Dev Maintenance** | Laptop Developer | VPS Ubuntu | `~/.ssh/id_ed25519` di laptop | VPS: `/home/deployer/.ssh/authorized_keys` | Remote login, debugging terminal, rsync manual |
| **Arah 2: GitHub CI/CD Deploy** | GitHub Actions Runner | VPS Ubuntu | GitHub Repo Secrets (`VPS_SSH_KEY`) | VPS: `/home/deployer/.ssh/authorized_keys` | Mengirim artifact binary dan restart service Systemd |
| **Arah 3: VPS ke GitHub** | Server VPS Ubuntu | GitHub (github.com) | VPS: `/home/deployer/.ssh/id_ed25519_github` | GitHub: **Deploy Keys** (Write) atau **SSH Keys** akun | `git pull`, `git push`, update tag rilis dari server |

---

## 2. Bagian 1: Pembuatan & Konfigurasi User `deployer` di Ubuntu

### A. Pembuatan Akun & Penugasan Grup
Jalankan perintah berikut sebagai user dengan hak `sudo` atau `root`:

```bash
# 1. Pastikan grup service 'erp' sudah ada di sistem
sudo groupadd -f -r erp

# 2. Buat user 'deployer' dengan shell bash dan direktori home
sudo useradd -m -s /bin/bash -g erp -G sudo deployer

# 3. Tetapkan password sementara yang kuat untuk user deployer
sudo passwd deployer
```

> [!NOTE]
> Menambahkan `deployer` ke grup primer `erp` (`-g erp`) memastikan semua file yang dibuat oleh `deployer` otomatis dapat dibaca dan dieksekusi oleh service `erp` yang dikelola Systemd.

---

### B. Konfigurasi Sudoers Tanpa Password Terbatas (`NOPASSWD`)
User `deployer` memerlukan hak untuk me-restart layanan Systemd dan memuat ulang reverse proxy tanpa meminta interaksi password (karena GitHub Actions berjalan secara non-interaktif).

Gunakan prinsip hak akses paling minim (*Principle of Least Privilege*). **Jangan pernah memberikan `deployer ALL=(ALL) NOPASSWD: ALL`**.

Buat berkas aturan sudoers khusus:
```bash
sudo nano /etc/sudoers.d/erp-deployer
```

Isi dengan baris berikut:
```text
# Aturan hak akses terbatas untuk otomasi deployment ERP Monolith
deployer ALL=(ALL) NOPASSWD: /bin/systemctl restart erp-backend, \
                             /bin/systemctl restart erp-frontend, \
                             /bin/systemctl restart erp.target, \
                             /bin/systemctl reload erp.target, \
                             /bin/systemctl status erp*, \
                             /bin/systemctl is-active erp*, \
                             /bin/journalctl -u erp*, \
                             /bin/systemctl reload nginx, \
                             /bin/systemctl restart nginx
```

Uji validitas sintaks berkas sudoers sebelum disimpan (sangat penting untuk mencegah sistem sudo terkunci/rusak):
```bash
sudo visudo -cf /etc/sudoers.d/erp-deployer
```
*Output yang benar:* `/etc/sudoers.d/erp-deployer: parsed OK`

---

### C. Hak Akses Direktori Proyek & SetGID
Pastikan direktori instalasi aplikasi ERP dimiliki oleh user `deployer` dan grup `erp`, serta aktifkan bit **SetGID** agar file atau subdirektori baru yang tercipta mewarisi grup `erp`:

```bash
# 1. Buat direktori struktur jika belum ada
sudo mkdir -p /opt/dev/erp_monolith/backend/bin
sudo mkdir -p /opt/dev/erp_monolith/backend/migrations
sudo mkdir -p /opt/dev/erp_monolith/frontend/build
sudo mkdir -p /etc/erp

# 2. Tetapkan kepemilikan deployer:erp
sudo chown -R deployer:erp /opt/dev/erp_monolith
sudo chown -R deployer:erp /etc/erp

# 3. Izin direktori: Pemilik (rwx), Grup (rwx), Lainnya (r-x)
sudo chmod -R 775 /opt/dev/erp_monolith
sudo chmod 750 /etc/erp

# 4. Terapkan SetGID bit pada direktori (sub-file baru otomatis ber-grup 'erp')
sudo find /opt/dev/erp_monolith -type d -exec chmod g+s {} +
```

---

### D. Inisialisasi Direktori `~/.ssh` dengan Izin Presisi
OpenSSH menolak otentikasi jika izin direktori atau berkas kunci terlalu longgar (*permissive*). Siapkan direktori `.ssh` untuk user `deployer`:

```bash
# Masuk ke shell user deployer
sudo su - deployer

# Buat direktori .ssh dengan izin 700 (hanya pemilik yang boleh baca/tulis/buka)
mkdir -p ~/.ssh
chmod 700 ~/.ssh

# Buat berkas authorized_keys dengan izin 600 (hanya pemilik yang boleh baca/tulis)
touch ~/.ssh/authorized_keys
chmod 600 ~/.ssh/authorized_keys

# Pastikan kepemilikan absolut berada pada deployer:erp
chown -R deployer:erp ~/.ssh
```

---

## 3. Bagian 2: SSH dari Server Ubuntu ke GitHub (Push & Pull)

### A. Skenario & Use Case VPS ke GitHub
Skenario ini terjadi saat **Server Ubuntu bertindak sebagai SSH Client** dan **GitHub bertindak sebagai Server Tujuan**. 

Kondisi ini digunakan ketika:
1. Menjalankan skrip `deploy.sh` langsung di VPS yang mengeksekusi `git fetch` dan `git pull origin main`.
2. Melakukan update file konfigurasi, dokumentasi, atau migrasi di VPS lalu melakukan `git push` kembali ke repositori GitHub.
3. Server membuat commit otomatis (misal pembaruan versi rilis, tag git, atau backup metadata otomatis).

---

### B. Pembuatan Pasangan Kunci SSH Ed25519 di VPS
Masuk sebagai user `deployer` di VPS, lalu generate pasangan kunci berbasis algoritma modern **Ed25519** (lebih cepat dan lebih aman daripada RSA):

```bash
sudo su - deployer

# Generate SSH Key khusus untuk koneksi ke GitHub (tanpa passphrase agar otomasi lancar)
ssh-keygen -t ed25519 -C "deployer-vps-erp@github" -f ~/.ssh/id_ed25519_github -N ""
```

Perintah di atas menghasilkan 2 berkas di `~/.ssh/`:
- `id_ed25519_github`: **Private Key** (Rahasia! Jangan pernah disebarkan keluar dari VPS).
- `id_ed25519_github.pub`: **Public Key** (Kunci yang akan ditempel di GitHub).

Tampilkan isi public key untuk disalin:
```bash
cat ~/.ssh/id_ed25519_github.pub
```
*Contoh output:*
`ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIExampleKeyStringDeployerVPS deployer-vps-erp@github`

---

### C. Pendaftaran Kunci di GitHub: Deploy Key vs Personal/Machine Key

Pilih salah satu metode pendaftaran sesuai kebijakan tim Anda:

#### Pilihan 1: Deploy Keys (Spesifik per Repositori) — ⭐ Direkomendasikan
Deploy Key membatasi akses SSH hanya pada satu repositori tertentu. Sangat aman karena jika kunci VPS terekspos, repositori Anda yang lain tidak terpengaruh.

1. Buka repositori proyek di GitHub: `https://github.com/<organisasi>/<nama_repo>`.
2. Masuk ke menu **Settings** ➔ **Deploy keys** ➔ klik tombol **Add deploy key**.
3. Isi data:
   - **Title**: `VPS Production ERP - Deployer (Push/Pull)`
   - **Key**: Paste seluruh isi berkas `id_ed25519_github.pub`.
   - **Allow write access**: Centang opsi ini **Wajib** jika server Ubuntu perlu melakukan `git push` ke GitHub! (Jika tidak dicentang, server hanya memiliki izin read-only / pull).
4. Klik **Add key**.

> [!WARNING]
> GitHub melarang 1 Deploy Key yang sama digunakan di lebih dari 1 repositori. Jika Anda memiliki submodul atau repositori terpisah, gunakan *Machine User* (Pilihan 2) atau buat key terpisah per repo.

#### Pilihan 2: Personal SSH Key / Machine User Account
Jika server VPS perlu mengakses banyak repositori sekaligus (misalnya backend, frontend, konfigurasi k8s terpisah):
1. Buat akun GitHub khusus (misal: `bot-erp-deployer`) atau gunakan akun GitHub admin.
2. Masuk ke **Settings Akun** ➔ **SSH and GPG keys** ➔ **New SSH key**.
3. Tempelkan isi `id_ed25519_github.pub` dan simpan. Berikan akses `Write` pada tim atau kolaborator di organisasi.

---

### D. Verifikasi Host Key & Tes Koneksi SSH
Tambahkan fingerprint host GitHub ke daftar `known_hosts` VPS untuk mencegah prompt interaktif:

```bash
# Tambahkan fingerprint GitHub secara otomatis
ssh-keyscan -t ed25519 github.com >> ~/.ssh/known_hosts
chmod 644 ~/.ssh/known_hosts
```

Uji otentikasi SSH dari VPS ke GitHub menggunakan kunci yang baru dibuat:
```bash
ssh -i ~/.ssh/id_ed25519_github -T git@github.com
```

*Output yang benar jika sukses:*
```text
Hi <nama_repo_atau_user>! You've successfully authenticated, but GitHub does not provide shell access.
```

---

### E. Konfigurasi Git & Migrasi Remote URL ke SSH
Agar Git di VPS dapat mengenali identitas pembuat commit dan menggunakan SSH (bukan HTTPS yang membutuhkan Personal Access Token):

```bash
# Konfigurasi identitas Git global di user deployer
git config --global user.name "ERP Deployer Bot"
git config --global user.email "deployer@perusahaan.com"

# Masuk ke direktori repositori di VPS
cd /opt/dev/erp_monolith

# Periksa remote URL saat ini
git remote -v

# Jika remote URL masih menggunakan format HTTPS (https://github.com/...), ubah ke format SSH:
git remote set-url origin git@github.com:perusahaan/erp_monolith.git

# Pastikan URL sudah berformat git@github.com:...
git remote -v
```

---

### F. Pengujian Push & Pull Langsung dari Server
Lakukan pengujian dua arah secara langsung dari terminal VPS:

```bash
cd /opt/dev/erp_monolith

# 1. Uji Git Fetch & Pull
git fetch origin main
git status

# 2. Uji Git Push (Contoh: membuat tag rilis baru atau branch pengujian)
# Membuat tag release lokal
git tag -a v1.0.0-test -m "Uji coba SSH push dari server Ubuntu VPS"

# Push tag ke GitHub
git push origin v1.0.0-test

# Hapus tag pengujian setelah sukses
git push origin --delete v1.0.0-test
git tag -d v1.0.0-test
```

Jika perintah di atas berhasil tanpa meminta password atau token, berarti autentikasi SSH dari VPS ke GitHub telah aktif 100%.

---

## 4. Bagian 3: Pembuatan & Manajemen Berkas `~/.ssh/config`

### A. Manfaat & Parameter Kunci `~/.ssh/config`
Berkas `~/.ssh/config` adalah berkas konfigurasi klien SSH yang mengatur profil koneksi untuk setiap host.

**Keuntungan Menggunakan SSH Config:**
- **Pintasan Praktis:** Tidak perlu mengetik argumen panjang seperti `ssh -i ~/.ssh/id_ed25519_github -p 22 git@github.com`.
- **Mencegah Kegagalan Autentikasi:** Parameter `IdentitiesOnly yes` mencegah SSH mengirim semua kunci yang ada di SSH agent (yang bisa memicu error `Too many authentication failures`).
- **Koneksi Stabil (KeepAlive):** Mencegah koneksi SSH *freezing* atau terputus di tengah proses transfer data yang lama.

---

### B. Konfigurasi `~/.ssh/config` di Sisi Server VPS (`deployer`)
Buat berkas konfigurasi di VPS sebagai user `deployer`:

```bash
sudo su - deployer
nano ~/.ssh/config
```

Isi dengan template berikut:

```ssh-config
# ========================================================
# Profil GitHub Default untuk Git Clone / Fetch / Push
# ========================================================
Host github.com
    HostName github.com
    User git
    IdentityFile ~/.ssh/id_ed25519_github
    IdentitiesOnly yes
    ServerAliveInterval 30
    ServerAliveCountMax 4
    StrictHostKeyChecking accept-new
```

Setelah konfigurasi ini dipasang, setiap perintah `git clone`, `git fetch`, `git pull`, atau `git push` yang mengarah ke `github.com` akan otomatis menggunakan file kunci `~/.ssh/id_ed25519_github` tanpa perlu parameter tambahan!

---

### C. Konfigurasi Multi-Akun / Multi-Repo GitHub
Jika server VPS Anda perlu mengelola dua akun atau repositori yang berbeda (misalnya akun organisasi perusahaan dan akun repositori vendor/pribadi) di mana GitHub tidak mengizinkan satu Deploy Key dipakai berulang:

```bash
# Buat kunci kedua
ssh-keygen -t ed25519 -C "deployer-vendor@github" -f ~/.ssh/id_ed25519_vendor -N ""
```

Tambahkan ke `~/.ssh/config`:

```ssh-config
# Akun Utama: Enterprise ERP Monolith
Host github.com
    HostName github.com
    User git
    IdentityFile ~/.ssh/id_ed25519_github
    IdentitiesOnly yes

# Akun Kedua: Plugin / Vendor Modules
Host github-vendor
    HostName github.com
    User git
    IdentityFile ~/.ssh/id_ed25519_vendor
    IdentitiesOnly yes
```

Cara clone dari akun kedua:
```bash
# Ganti host github.com dengan alias host yang didefinisikan:
git clone git@github-vendor:vendor-org/erp-custom-plugin.git
```

---

### D. Konfigurasi `~/.ssh/config` di Komputer Lokal Developer
Di komputer kerja lokal developer (Mac / Linux / Windows WSL), konfigurasikan `~/.ssh/config` agar login ke VPS cukup mengetik satu kata singkat:

```ssh-config
# ========================================================
# Koneksi Cepat ke VPS ERP Production
# ========================================================
Host erp-prod
    HostName 203.0.113.50             # Ganti dengan IP Publik atau Domain VPS
    Port 22                           # Ganti jika menggunakan port custom (misal 2222)
    User deployer
    IdentityFile ~/.ssh/id_ed25519_erp_deployer
    IdentitiesOnly yes
    ServerAliveInterval 60
    ServerAliveCountMax 3
    ForwardAgent no

# ========================================================
# Koneksi Cepat ke VPS ERP Staging
# ========================================================
Host erp-stage
    HostName 203.0.113.51
    Port 22
    User deployer
    IdentityFile ~/.ssh/id_ed25519_erp_deployer
    IdentitiesOnly yes
```

Dengan konfigurasi di atas, Anda dapat mengakses VPS dan mentransfer file dengan sangat ringkas:
```bash
# Login SSH
ssh erp-prod

# Rsync file konfigurasi
rsync -avz ./deploy.sh erp-prod:/opt/dev/erp_monolith/

# Melihat log systemd secara live
ssh erp-prod "journalctl -u erp-backend -f"
```

---

### E. Izin Akses Wajib File Config
File `~/.ssh/config` berisi informasi arsitektur server Anda. Batasi izin akses agar tidak bisa dibaca oleh user lain di OS:

```bash
chmod 600 ~/.ssh/config
```

---

## 5. Bagian 4: SSH dari GitHub Actions Deploy ke Server Ubuntu (CI/CD)

### A. Arsitektur Alur Deployment Runner ke VPS
Pada arsitektur ini:
- **SSH Client**: GitHub Actions Runner (mesin Ubuntu sementara milik GitHub).
- **SSH Server**: VPS Ubuntu Anda (menerima koneksi pada port 22/custom).
- **Tujuan**: Mengirim file biner Golang, file build frontend SvelteKit, dan mengeksekusi perintah restart Systemd.

---

### B. Pembuatan Pasangan Kunci SSH Khusus Deployment
**Jangan pernah memakai kunci pribadi developer!** Buat pasangan kunci khusus yang didedikasikan hanya untuk GitHub Actions CI/CD:

Jalankan perintah ini di komputer lokal developer atau langsung di VPS:
```bash
ssh-keygen -t ed25519 -C "github-actions-ci-cd-erp" -f ~/github_actions_deploy -N ""
```

Hasil:
- `github_actions_deploy.pub` ➔ Kunci Publik untuk **Server VPS**.
- `github_actions_deploy` ➔ Kunci Privat untuk **GitHub Repository Secrets**.

---

### C. Pemasangan Public Key di VPS (`authorized_keys`) & Opsi Restriksi
Tambahkan isi `github_actions_deploy.pub` ke berkas `authorized_keys` milik user `deployer` di VPS:

```bash
sudo su - deployer

# Tambahkan public key ke daftar kunci yang diizinkan
cat ~/github_actions_deploy.pub >> ~/.ssh/authorized_keys

# Pastikan permission aman
chmod 600 ~/.ssh/authorized_keys
```

#### Hardening Ekstra: Opsi Restriksi OpenSSH di `authorized_keys`
Untuk keamanan level Enterprise, Anda dapat mengunci fungsi kunci ini agar **tidak dapat digunakan untuk port forwarding atau interactive pseudo-terminal**.

Buka `~/.ssh/authorized_keys`, lalu tambahkan parameter pembatasan di depan baris kunci tersebut:
```text
no-port-forwarding,no-X11-forwarding,no-agent-forwarding ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAI... github-actions-ci-cd-erp
```

---

### D. Konfigurasi GitHub Repository Secrets
Buka GitHub di browser: **Repository Settings** ➔ **Secrets and variables** ➔ **Actions** ➔ **New repository secret**.

Tambahkan 4 secret berikut:

| Nama Secret | Nilai (*Value*) | Keterangan |
|---|---|---|
| `VPS_HOST` | `203.0.113.50` (IP Publik VPS Anda) | Alamat IP atau hostname VPS tujuan |
| `VPS_PORT` | `22` (atau port SSH custom Anda) | Port SSH daemon VPS |
| `VPS_USER` | `deployer` | Username operasional deployment |
| `VPS_SSH_KEY` | *(Salin seluruh isi file `github_actions_deploy`)* | Private key (harus utuh dari BEGIN hingga END) |

> [!CAUTION]
> Saat menyalin `VPS_SSH_KEY`, pastikan baris `-----BEGIN OPENSSH PRIVATE KEY-----` dan `-----END OPENSSH PRIVATE KEY-----` ikut tersalin tanpa ada spasi atau karakter kosong di awal/akhir.

---

### E. Integrasi di GitHub Actions Workflow (`rsync` & `ssh-action`)
Berikut adalah implementasi standar produksi di `.github/workflows/deploy.yml`:

```yaml
name: 🚀 Deploy to Production VPS

on:
  push:
    branches:
      - main

jobs:
  deploy:
    name: 📦 Transfer Artifacts & Restart Services
    runs-on: ubuntu-latest
    steps:
      - name: 📥 Checkout Code
        uses: actions/checkout@v4

      # Setup SSH Agent & Known Hosts pada Runner GitHub
      - name: 🔑 Configure SSH Key on Runner
        run: |
          mkdir -p ~/.ssh
          echo "${{ secrets.VPS_SSH_KEY }}" > ~/.ssh/id_ed25519
          chmod 600 ~/.ssh/id_ed25519
          # Pindai host key VPS agar tidak muncul prompt 'Host key verification failed'
          ssh-keyscan -p ${{ secrets.VPS_PORT || 22 }} ${{ secrets.VPS_HOST }} >> ~/.ssh/known_hosts
          chmod 644 ~/.ssh/known_hosts

      # Transfer File menggunakan Rsync / SCP Native OpenSSH
      - name: 🚚 Transfer Deployment Artifacts
        env:
          PORT: ${{ secrets.VPS_PORT || 22 }}
          USER: ${{ secrets.VPS_USER }}
          HOST: ${{ secrets.VPS_HOST }}
        run: |
          # Mengirim binary backend ke ekstensi sementara .new
          scp -P $PORT backend/bin/erp-backend $USER@$HOST:/opt/dev/erp_monolith/backend/bin/erp-backend.new
          
          # Sinkronisasi direktori build frontend
          rsync -avz -e "ssh -p $PORT" --delete frontend/build/ $USER@$HOST:/opt/dev/erp_monolith/frontend/build/

      # Eksekusi Restart Layanan via Remote SSH Command
      - name: ⚡ Atomic Swap & Restart Systemd
        uses: appleboy/ssh-action@v1.0.3
        with:
          host: ${{ secrets.VPS_HOST }}
          username: ${{ secrets.VPS_USER }}
          key: ${{ secrets.VPS_SSH_KEY }}
          port: ${{ secrets.VPS_PORT || 22 }}
          script_stop: true
          script: |
            set -e
            echo "--- [1/3] Atomic Swap Binary Go ---"
            chmod +x /opt/dev/erp_monolith/backend/bin/erp-backend.new
            mv -f /opt/dev/erp_monolith/backend/bin/erp-backend.new /opt/dev/erp_monolith/backend/bin/erp-backend

            echo "--- [2/3] Restart Systemd Services ---"
            sudo systemctl restart erp-backend.service
            sudo systemctl restart erp-frontend.service

            echo "--- [3/3] Verifikasi Health Check ---"
            sleep 3
            curl -fs http://127.0.0.1:8080/health || {
              echo "❌ Health check gagal!"
              sudo systemctl status erp-backend.service --no-pager
              exit 1
            }
            echo "✅ Deployment Sukses & Service Aktif!"
```

---

## 6. Bagian 5: Hardening Keamanan OpenSSH Server (`sshd_config`)

Setelah semua kunci SSH terpasang dan teruji, perkuat konfigurasi daemon SSH server (`/etc/ssh/sshd_config`) untuk menutup celah serangan brute-force:

Edit berkas konfigurasi daemon SSH:
```bash
sudo nano /etc/ssh/sshd_config
# atau jika di Ubuntu modern: /etc/ssh/sshd_config.d/50-cloud-init.conf
```

Terapkan parameter rekomendasi berikut:
```text
# 1. Nonaktifkan otentikasi menggunakan password teks biasa (Wajib SSH Key)
PasswordAuthentication no
ChallengeResponseAuthentication no

# 2. Nonaktifkan login langsung akun root
PermitRootLogin no

# 3. Batasi hanya user tertentu yang boleh login via SSH
AllowUsers deployer developeradmin

# 4. Batasi percobaan login gagal
MaxAuthTries 3

# 5. Nonaktifkan otentikasi berbasis file .rhosts kosong
PermitEmptyPasswords no

# 6. Set timeout sesi tidak aktif
ClientAliveInterval 300
ClientAliveCountMax 2
```

Uji sintaks berkas konfigurasi sebelum me-restart daemon SSH:
```bash
sudo sshd -t
```
*Jika tidak ada pesan error yang muncul, berarti konfigurasi valid.*

Terapkan perubahan dengan me-restart layanan SSH:
```bash
sudo systemctl restart ssh || sudo systemctl restart sshd
```

> [!CAUTION]
> **PENTING:** Jangan pernah menutup sesi terminal SSH Anda saat ini sebelum membuka tab terminal baru dan memastikan Anda dapat login kembali menggunakan `ssh -i ~/.ssh/key deployer@<IP_VPS>`.

---

## 7. Bagian 6: Cheatsheet Perintah Cepat & Troubleshooting

### Tabel Referensi Izin Berkas SSH:
| Jalur Berkas / Direktori | Izin (*Mode*) | Pemilik (*Owner*) | Keterangan |
|---|---|---|---|
| `~/.ssh` | `700` (`drwx------`) | `deployer:erp` | Direktori SSH user |
| `~/.ssh/authorized_keys` | `600` (`-rw-------`) | `deployer:erp` | Public keys yang boleh login ke VPS |
| `~/.ssh/id_ed25519` / private key | `600` (`-rw-------`) | `deployer:erp` | Private key klien |
| `~/.ssh/id_ed25519.pub` | `644` (`-rw-r--r--`) | `deployer:erp` | Public key klien |
| `~/.ssh/config` | `600` (`-rw-------`) | `deployer:erp` | File alias & konfigurasi klien SSH |
| `~/.ssh/known_hosts` | `644` (`-rw-r--r--`) | `deployer:erp` | Fingerprint server yang dikenal |

Perintah 1 baris untuk merapikan semua izin berkas SSH:
```bash
sudo chown -R deployer:erp /home/deployer/.ssh && \
chmod 700 /home/deployer/.ssh && \
chmod 600 /home/deployer/.ssh/authorized_keys /home/deployer/.ssh/config /home/deployer/.ssh/id_* 2>/dev/null || true && \
chmod 644 /home/deployer/.ssh/*.pub /home/deployer/.ssh/known_hosts 2>/dev/null || true
```

---

### Diagnosa & Troubleshooting Masalah Umum

#### 1. `Permission denied (publickey)`
- **Penyebab:**
  1. Izin folder `~/.ssh` atau `~/.ssh/authorized_keys` terlalu terbuka (misal `777` atau `755`). OpenSSH akan menolak kunci secara otomatis (*StrictModes*).
  2. Public key yang didaftarkan di `authorized_keys` tidak cocok dengan private key yang dikirim oleh client.
  3. Pemilik file bukan user `deployer` (misal masih dimiliki oleh `root`).
- **Solusi:**
  Jalankan perintah normalisasi izin di atas. Jalankan koneksi dengan mode verbose untuk melihat proses handshake:
  ```bash
  ssh -vvv deployer@<IP_VPS>
  ```

#### 2. `Host key verification failed`
- **Penyebab:** Fingerprint server tujuan belum tercatat di file `~/.ssh/known_hosts` atau IP server telah berganti (*man-in-the-middle alert*).
- **Solusi:**
  Tambahkan fingerprint host menggunakan `ssh-keyscan`:
  ```bash
  ssh-keyscan -p <PORT> <HOST> >> ~/.ssh/known_hosts
  ```

#### 3. `fatal: Could not read from remote repository` (Saat Git Push dari VPS ke GitHub)
- **Penyebab:**
  1. Kunci SSH di VPS belum didaftarkan di GitHub Deploy Keys.
  2. Kunci sudah didaftarkan di Deploy Keys, tetapi checkbox **"Allow write access"** belum dicentang.
  3. URL remote repository masih HTTPS bukan SSH (`git@github.com:...`).
- **Solusi:**
  1. Buka GitHub Repo ➔ **Settings** ➔ **Deploy keys** ➔ pastikan opsi **Allow write access** aktif.
  2. Periksa `git remote -v`. Jika HTTPS, ubah ke:
     `git remote set-url origin git@github.com:<org>/<repo>.git`.
  3. Uji: `ssh -T git@github.com`.

#### 4. `Too many authentication failures for deployer`
- **Penyebab:** Klien SSH lokal atau agent memiliki lebih dari 5 private key yang dicoba satu per satu hingga mencapai batas limit server.
- **Solusi:**
  Tambahkan opsi `IdentitiesOnly yes` pada berkas `~/.ssh/config` di host yang bersangkutan:
  ```ssh-config
  Host erp-prod
      HostName 203.0.113.50
      User deployer
      IdentityFile ~/.ssh/id_ed25519_erp_deployer
      IdentitiesOnly yes
  ```

#### 5. `sudo: a password is required`
- **Penyebab:** Perintah yang dipanggil oleh skrip deployment di GitHub Actions tidak cocok dengan daftar perintah yang diizinkan dalam `/etc/sudoers.d/erp-deployer`.
- **Solusi:**
  Pastikan perintah menggunakan full path (contoh: `/bin/systemctl restart erp-backend`). Periksa status parsing sudoers:
  ```bash
  sudo visudo -cf /etc/sudoers.d/erp-deployer
  ```
