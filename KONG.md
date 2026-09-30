# 🦍 Panduan Instalasi & Penggunaan Kong API Gateway — Enterprise ERP Monolith

Dokumen ini adalah panduan resmi untuk menginstal, mengonfigurasi dengan basis data **PostgreSQL**, dan mendaftarkan endpoint/route pada **Kong API Gateway (OSS)** di depan backend **Echo v5** untuk **Enterprise ERP Monolith**.

---

## 📑 Daftar Isi
1. [Arsitektur Kong API Gateway di ERP Monolith](#1-arsitektur-kong-api-gateway-di-erp-monolith)
2. [Instalasi Kong di Ubuntu](#2-instalasi-kong-di-ubuntu)
3. [Konfigurasi Kong dengan Basis Data PostgreSQL](#3-konfigurasi-kong-dengan-basis-data-postgresql)
   - [A. Pembuatan User & Database Kong di PostgreSQL](#a-pembuatan-user--database-kong-di-postgresql)
   - [B. Konfigurasi Berkas `/etc/kong/kong.conf`](#b-konfigurasi-berkas-etckongkongconf)
   - [C. Menjalankan Bootstrap Migrasi Database](#c-menjalankan-bootstrap-migrasi-database)
   - [D. Menjalankan Layanan Kong](#d-menjalankan-layanan-kong)
4. [Cara Menambahkan Service, Route & Endpoint di Proyek Ini](#4-cara-menambahkan-service-route--endpoint-di-proyek-ini)
   - [A. Konsep: Upstream, Service, Route & Plugin](#a-konsep-upstream-service-route--plugin)
   - [B. Mendaftarkan Service Backend & Frontend](#b-mendaftarkan-service-backend--frontend)
   - [C. Mendaftarkan Route: Routing Lalu Lintas API vs Frontend](#c-mendaftarkan-route-routing-lalu-lintas-api-vs-frontend)
   - [D. Memasang Plugin: Rate Limiting (Redis), CORS, dan JWT](#d-memasang-plugin-rate-limiting-redis-cors-dan-jwt)
5. [Konfigurasi Deklaratif Alternatif (decK / `kong.yaml`)](#5-konfigurasi-deklaratif-alternatif-deck--kongyaml)
6. [Pengujian Endpoint & Integrasi Frontend (SvelteKit 2)](#6-pengujian-endpoint--integrasi-frontend-sveltekit-2)
7. [Troubleshooting & Pemeliharaan](#7-troubleshooting--pemeliharaan)

---

## 1. Arsitektur Kong API Gateway di ERP Monolith

Kong bertindak sebagai gerbang terpusat (*Reverse Proxy & API Gateway*) yang melindungi backend Echo v5 dan melayani frontend SvelteKit 2:

```text
┌────────────────────────────────────────────────────────┐
│ Client (Browser / SvelteKit 2 / Mobile / Integrator)  │
└───────────────────────────┬────────────────────────────┘
                            │ HTTP Request
                            ▼
┌────────────────────────────────────────────────────────┐
│ Kong API Gateway                                       │
│ • Proxy Port :8000 (HTTP) / :8443 (HTTPS)              │
│ • Admin API :8001 / Admin GUI (Manager) :8002          │
│ • Fitur: Rate Limiting (Redis), CORS, JWT, Logging     │
└──────┬──────────────────────┬──────────────────────┬───┘
       │                      │                      │
       │ (Direct :5432)       │ Path: /api/v1/*      │ Path: /* (Selain /api/v1)
       ▼                      ▼                      ▼
┌──────────────┐   ┌──────────────────────┐   ┌──────────────────────┐
│ PostgreSQL   │   │ Echo v5 Backend      │   │ SvelteKit 2 Frontend │
│ (kong_db)    │   │ Monolith Port :8080  │   │ Node/Bun Port :3000  │
│ Port 5432    │   └──────────┬───────────┘   └──────────────────────┘
└──────────────┘              │
                              ▼ (pgxpool via PgBouncer :6432)
                   ┌──────────────────────┐
                   │ PostgreSQL (erp_db)  │
                   └──────────────────────┘
```

> [!CAUTION]
> **ATURAN PENTING ARSITEKTUR BASIS DATA KONG:**
> Kong **WAJIB** terhubung langsung ke **PostgreSQL port 5432**, **BUKAN ke port PgBouncer 6432**!
> Kong mengelola skema internalnya sendiri dengan migrasi DDL, advisory locks, dan connection pooling internal yang tidak kompatibel dengan PgBouncer transaction pooling.

---

## 2. Instalasi Kong di Ubuntu

Gunakan repositori resmi Kong (Cloudsmith) untuk Ubuntu 22.04 (Jammy) atau 24.04 (Noble).

### Langkah 1: Pasang Dependensi
```bash
sudo apt update && sudo apt install -y curl apt-transport-https lsb-release
```

### Langkah 2: Tambahkan Repositori Resmi Kong Gateway OSS
Gunakan repositori versi 3.9 (`gateway-39`) dari Cloudsmith Kong resmi:
```bash
curl -1sLf "https://packages.konghq.com/public/gateway-39/setup.deb.sh" | sudo -E bash
```

> [!TIP]
> **Catatan Penamaan Repositori Cloudsmith:**
> Cloudsmith Kong menggunakan penamaan minor version seperti `gateway-39` (versi 3.9 OSS). Penggunaan URL lama seperti `gateway-3x` akan mengembalikan kode HTTP 404 sehingga script tidak mengeksekusi apapun dan menyebabkan error `E: Unable to locate package kong` pada saat `apt install`.

### Langkah 3: Pasang Paket Kong
```bash
sudo apt update
sudo apt install -y kong
```

> [!NOTE]
> **Metode Alternatif (Direct `.deb` Package):**
> Jika ingin memasang langsung tanpa script registrasi repositori:
> ```bash
> # Untuk Ubuntu 24.04 (Noble) x86_64:
> curl -Lo /tmp/kong.deb "https://packages.konghq.com/public/gateway-39/deb/ubuntu/pool/noble/main/k/ko/kong_3.9.3/kong_3.9.3_amd64.deb"
> sudo apt install -y /tmp/kong.deb
> ```

### Langkah 4: Verifikasi Instalasi
```bash
kong version
```

---

## 3. Konfigurasi Kong dengan Basis Data PostgreSQL

### A. Pembuatan User & Database Kong di PostgreSQL

Masuk ke PostgreSQL melalui user `postgres` dan buat database terpisah untuk Kong:

```bash
sudo -u postgres psql << 'EOF'
-- 1. Buat User kong dengan password aman
CREATE USER kong WITH PASSWORD 'kong_secret';

-- 2. Buat database kong dengan kepemilikan user kong
CREATE DATABASE kong OWNER kong;

-- 3. Berikan hak akses penuh ke schema public
\c kong
GRANT ALL PRIVILEGES ON SCHEMA public TO kong;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO kong;
EOF
```

---

### B. Konfigurasi Berkas `/etc/kong/kong.conf`

Salin template konfigurasi bawaan Kong:
```bash
sudo cp /etc/kong/kong.conf.default /etc/kong/kong.conf
sudo nano /etc/kong/kong.conf
```

Sesuaikan parameter konfigurasi utama berikut:

```ini
# ========================================================
# Kong Configuration for Enterprise ERP Monolith
# ========================================================

# ----------------------------------------
# 1. Konfigurasi Database PostgreSQL
# ----------------------------------------
database = postgres
pg_host = 127.0.0.1
pg_port = 5432                  # WAJIB DIRECT PORT 5432 (BUKAN 6432)
pg_user = kong
pg_password = kong_secret
pg_database = kong
pg_ssl = off

# ----------------------------------------
# 2. Port Jaringan & Listener
# ----------------------------------------
# Proxy Port: Jalur utama request dari frontend / user
proxy_listen = 0.0.0.0:8000, 0.0.0.0:8443 ssl

# Admin API: Hanya boleh diakses oleh localhost / admin internal
admin_listen = 127.0.0.1:8001, 127.0.0.1:8444 ssl

# Kong Manager (Web UI):
admin_gui_listen = 0.0.0.0:8002
admin_gui_url = http://localhost:8002

# ----------------------------------------
# 3. Tuning Performa & DNS
# ----------------------------------------
nginx_worker_processes = auto
dns_resolver = 127.0.0.53, 8.8.8.8
```

---

### C. Menjalankan Bootstrap Migrasi Database

Inisialisasi tabel dan skema internal Kong di basis data `kong`:

```bash
sudo kong migrations bootstrap -c /etc/kong/kong.conf
```

*Tunggu hingga seluruh migrasi schema Kong selesai (muncul pesan `Database is up-to-date`).*

---

### D. Menjalankan Layanan Kong

```bash
# Start Kong
sudo kong start -c /etc/kong/kong.conf

# Untuk memastikan kong berjalan otomatis via systemd:
sudo systemctl enable kong
sudo systemctl start kong

# Periksa status kesehatan Kong
curl -i http://localhost:8001/status
```

Output sukses akan mengembalikan kode HTTP `200 OK` dengan status database `reachable: true`.

---

## 4. Cara Menambahkan Service, Route & Endpoint di Proyek Ini

### A. Konsep Dasar Kong
- **Service**: Merepresentasikan aplikasi upstream yang dituju (backend Echo v5 di `http://127.0.0.1:8080` dan frontend SvelteKit 2 di `http://127.0.0.1:3000`).
- **Route**: Aturan pencocokan (*path matching*) yang menentukan ke Service mana request dari client akan diarahkan.
- **Plugin**: Middleware fungsional tambahan (Rate Limiting, CORS, JWT Auth, Logging) yang dapat dipasang ke Service atau Route tertentu.

---

### B. Mendaftarkan Service Backend & Frontend

#### 1. Mendaftarkan Service Backend Echo v5
Daftarkan aplikasi backend monolith (Echo v5) sebagai upstream service:

```bash
curl -i -X POST http://localhost:8001/services \
  --data "name=erp-backend-service" \
  --data "url=http://127.0.0.1:8080" \
  --data "connect_timeout=60000" \
  --data "write_timeout=60000" \
  --data "read_timeout=60000"
```

#### 2. Mendaftarkan Service Frontend SvelteKit 2
Daftarkan server frontend SvelteKit 2 sebagai upstream service:

```bash
curl -i -X POST http://localhost:8001/services \
  --data "name=erp-frontend-service" \
  --data "url=http://127.0.0.1:3000" \
  --data "connect_timeout=60000" \
  --data "write_timeout=60000" \
  --data "read_timeout=60000"
```

> [!NOTE]
> - **Mode Standalone / Production**: Port `3000` adalah port default server Node/Bun SvelteKit (`adapter-node`).
> - **Mode Pengembangan (Vite Dev Server)**: Jika sedang menjalankan `bun run dev`, ganti url menjadi `http://127.0.0.1:5173`.

---

### C. Mendaftarkan Route: Routing Lalu Lintas API vs Frontend

Kong Gateway menggunakan algoritma pencocokan prefix terpanjang (*longest prefix match precedence*):
1. **Endpoint `/api/v1`** (prefix 7 karakter) diprioritaskan langsung ke **`erp-backend-service`**.
2. **Endpoint selain `/api/v1`** (seperti root path `/`, prefix 1 karakter) diarahkan ke **`erp-frontend-service`**.

> [!NOTE]
> Parameter `strip_path=false` memastikan bahwa path URL (misal `/api/v1/acc` atau `/dashboard`) **tetap dikirim secara utuh** ke backend Echo v5 maupun frontend SvelteKit router.

#### 1. Mendaftarkan Endpoint `/api/v1` Langsung ke Service Backend
Perintah curl untuk mendaftarkan rute utama `/api/v1` langsung menuju backend Echo v5:

```bash
curl -i -X POST http://localhost:8001/services/erp-backend-service/routes \
  --data "name=route-backend-api-v1" \
  --data "paths[]=/api/v1" \
  --data "strip_path=false"
```

#### 2. Mendaftarkan Endpoint Selain `/api/v1` ke Service Frontend (Catch-All)
Perintah curl untuk mendaftarkan seluruh endpoint selain `/api/v1` (halaman root `/`, halaman UI aplikasi, assets statis SvelteKit `/_app/*`) langsung menuju frontend SvelteKit:

```bash
curl -i -X POST http://localhost:8001/services/erp-frontend-service/routes \
  --data "name=route-frontend-all" \
  --data "paths[]=/" \
  --data "strip_path=false"
```

> [!TIP]
> **Mengapa Rute Catch-All `/` Tidak Menimpa `/api/v1`?**
> Di Kong API Gateway, aturan pencocokan jalur (*path matching precedence*) memprioritaskan rute dengan karakter jalur terpanjang:
> - Request ke `http://localhost:8000/api/v1/...` cocok dengan rute `route-backend-api-v1` (panjang 7 karakter) ➔ **Diteruskan ke Backend (:8080)**.
> - Request selain `/api/v1` (seperti `http://localhost:8000/`, `/login`, `/dashboard`, `/_app/...`) cocok dengan rute `route-frontend-all` (panjang 1 karakter) ➔ **Diteruskan ke Frontend (:3000)**.

#### 3. Endpoint Khusus / Granular Backend (Opsional)
Jika Anda membutuhkan aturan plugin spesifik (misal rate-limiting khusus, method restriction, atau monitoring granular), Anda dapat mendaftarkan rute spesifik tambahan di bawah `erp-backend-service`:

##### a. Endpoint Health Check
```bash
curl -i -X POST http://localhost:8001/services/erp-backend-service/routes \
  --data "name=route-health" \
  --data "paths[]=/health" \
  --data "strip_path=false" \
  --data "methods[]=GET"
```

##### b. Sub-Rute Modul ERP Granular (Jika Perlu Dikelola Terpisah)
```bash
# Autentikasi & Pengguna (19_USR, 21_ADM)
curl -i -X POST http://localhost:8001/services/erp-backend-service/routes \
  --data "name=route-auth" \
  --data "paths[]=/api/v1/auth" \
  --data "strip_path=false"

# Modul Akuntansi & Keuangan (1_ACC, 4_GL, 2_AP, 3_AR)
curl -i -X POST http://localhost:8001/services/erp-backend-service/routes \
  --data "name=route-accounting" \
  --data "paths[]=/api/v1/acc" \
  --data "strip_path=false"

# Modul Persediaan & Pergudangan (9_INV)
curl -i -X POST http://localhost:8001/services/erp-backend-service/routes \
  --data "name=route-inventory" \
  --data "paths[]=/api/v1/inv" \
  --data "strip_path=false"

# Modul Penjualan & Kasir POS (10_SAL, 12_POS)
curl -i -X POST http://localhost:8001/services/erp-backend-service/routes \
  --data "name=route-sales" \
  --data "paths[]=/api/v1/sal" \
  --data "strip_path=false"

# Modul Pembelian & Pengadaan (8_PUR)
curl -i -X POST http://localhost:8001/services/erp-backend-service/routes \
  --data "name=route-purchasing" \
  --data "paths[]=/api/v1/pur" \
  --data "strip_path=false"

# Modul SDM & Penggajian (15_HRM, 16_PAY, 17_ATT)
curl -i -X POST http://localhost:8001/services/erp-backend-service/routes \
  --data "name=route-hrm" \
  --data "paths[]=/api/v1/hrm" \
  --data "strip_path=false"
```

---

### D. Memasang Plugin

#### 1. Rate Limiting Terintegrasi dengan Redis
Membatasi request maksimal 120 per menit per IP/Consumer, dengan counter tersimpan di Redis (lihat [REDIS.md](file:///opt/dev/erp_monolith/REDIS.md)):

```bash
curl -i -X POST http://localhost:8001/services/erp-backend-service/plugins \
  --data "name=rate-limiting" \
  --data "config.minute=120" \
  --data "config.policy=redis" \
  --data "config.redis_host=127.0.0.1" \
  --data "config.redis_port=6379" \
  --data "config.redis_password=erp_redis_secret" \
  --data "config.redis_database=0"
```

#### 2. CORS Plugin (Untuk Frontend SvelteKit 2)
Mengizinkan frontend SvelteKit di `http://localhost:5173` mengakses API dengan credential:

```bash
curl -i -X POST http://localhost:8001/services/erp-backend-service/plugins \
  --data "name=cors" \
  --data "config.origins[]=http://localhost:5173" \
  --data "config.origins[]=http://localhost:3000" \
  --data "config.methods[]=GET" \
  --data "config.methods[]=POST" \
  --data "config.methods[]=PUT" \
  --data "config.methods[]=PATCH" \
  --data "config.methods[]=DELETE" \
  --data "config.methods[]=OPTIONS" \
  --data "config.headers[]=Accept" \
  --data "config.headers[]=Accept-Version" \
  --data "config.headers[]=Content-Length" \
  --data "config.headers[]=Content-Type" \
  --data "config.headers[]=Authorization" \
  --data "config.headers[]=X-Company-ID" \
  --data "config.headers[]=X-Branch-ID" \
  --data "config.credentials=true" \
  --data "config.max_age=3600"
```

#### 3. Correlation ID Plugin
Menyematkan request tracing header unik `X-Request-ID` secara otomatis ke setiap request:

```bash
curl -i -X POST http://localhost:8001/services/erp-backend-service/plugins \
  --data "name=correlation-id" \
  --data "config.header_name=X-Request-ID" \
  --data "config.generator=uuid" \
  --data "config.echo_downstream=true"
```

---

## 5. Konfigurasi Deklaratif Alternatif (decK / `kong.yaml`)

Selain via Admin API curl, Anda dapat menggunakan format deklaratif GitOps dengan file `kong.yaml` dan tool `deck`.

Simpan file deklaratif berikut sebagai `kong.yaml` di root proyek:

```yaml
_format_version: "3.0"

services:
  # ----------------------------------------------------
  # 1. Service Backend Echo v5 Monolith
  # ----------------------------------------------------
  - name: erp-backend-service
    url: http://127.0.0.1:8080
    connect_timeout: 60000
    write_timeout: 60000
    read_timeout: 60000
    routes:
      - name: route-backend-api-v1
        paths:
          - /api/v1
        strip_path: false
      - name: route-health
        paths:
          - /health
        strip_path: false
        methods:
          - GET
    plugins:
      - name: correlation-id
        config:
          header_name: X-Request-ID
          generator: uuid
          echo_downstream: true
      - name: rate-limiting
        config:
          minute: 120
          policy: redis
          redis_host: 127.0.0.1
          redis_port: 6379
          redis_password: erp_redis_secret

  # ----------------------------------------------------
  # 2. Service Frontend SvelteKit 2
  # ----------------------------------------------------
  - name: erp-frontend-service
    url: http://127.0.0.1:3000
    connect_timeout: 60000
    write_timeout: 60000
    read_timeout: 60000
    routes:
      - name: route-frontend-all
        paths:
          - /
        strip_path: false
```

### Sinkronisasi Deklaratif dengan `deck`:
```bash
# Validasi sintaks berkas konfigurasi
deck file validate kong.yaml

# Sinkronisasikan konfigurasi ke Kong Gateway
deck gateway sync kong.yaml --kong-addr http://localhost:8001
```

---

## 6. Pengujian Endpoint & Integrasi Frontend (SvelteKit 2)

### A. Uji Coba Endpoint via Proxy Kong (:8000)

1. Jalankan backend Echo v5:
   ```bash
   cd /opt/dev/erp_monolith/backend
   go run cmd/server/main.go
   ```

2. Jalankan server frontend SvelteKit 2:
   ```bash
   cd /opt/dev/erp_monolith/frontend
   bun run build && bun run preview   # port 3000
   # atau untuk mode dev: bun run dev -- --port 3000
   ```

3. Uji Coba Akses Melalui Kong Proxy Port `:8000`:
   - **Akses Frontend (Selain `/api/v1`)**:
     Buka di browser atau via curl:
     ```bash
     curl -i http://localhost:8000/
     ```
     Kong akan meneruskan request ke **`erp-frontend-service`** (`http://127.0.0.1:3000`).

   - **Akses Backend API (`/api/v1`)**:
     ```bash
     curl -i http://localhost:8000/api/v1/auth/login
     # atau health check
     curl -i http://localhost:8000/health
     ```
     Kong akan memprioritaskan dan meneruskan request ke **`erp-backend-service`** (`http://127.0.0.1:8080`).
   *Respon:*
   ```http
   HTTP/1.1 200 OK
   Content-Type: application/json
   X-Kong-Proxy-Latency: 1
   X-Kong-Upstream-Latency: 2
   X-RateLimit-Limit-Minute: 120
   X-RateLimit-Remaining-Minute: 119
   X-Request-ID: 7b311749-9831-41ae-ba5c-9c9dc76ff309

   {"status":"healthy","database":"connected"}
   ```

### B. Konfigurasi Frontend SvelteKit 2

Pada aplikasi frontend (`frontend/.env`), arahkan seluruh panggilan API langsung ke **Kong Proxy (Port 8000)**:

```env
PUBLIC_API_BASE_URL=http://localhost:8000
```

Frontend tidak pernah memanggil backend port `8080` secara langsung, melainkan selalu melewati Kong Gateway untuk proteksi keamanan dan logging.

---

## 7. Troubleshooting & Pemeliharaan

### 1. `Kong error: connection refused connecting to PostgreSQL`
- Pastikan konfigurasi `pg_port = 5432` di `/etc/kong/kong.conf`. Jangan menggunakan port PgBouncer `6432`.
- Pastikan PostgreSQL sedang berjalan: `sudo systemctl status postgresql`.

### 2. Memeriksa Daftar Route & Service yang Terdaftar:
```bash
# Melihat seluruh service
curl -s http://localhost:8001/services | jq

# Melihat seluruh route
curl -s http://localhost:8001/routes | jq

# Melihat plugin aktif
curl -s http://localhost:8001/plugins | jq
```

### 3. Menghapus / Memperbarui Route:
```bash
# Menghapus route berdasarkan nama
curl -i -X DELETE http://localhost:8001/routes/route-inventory
```

### 4. Restart & Reload Kong:
```bash
# Reload tanpa memutus koneksi aktif
sudo kong reload -c /etc/kong/kong.conf

# Restart total
sudo kong restart -c /etc/kong/kong.conf
```
