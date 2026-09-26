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
   - [B. Mendaftarkan Service Backend Echo v5](#b-mendaftarkan-service-backend-echo-v5)
   - [C. Mendaftarkan Endpoint / Rute Modul ERP](#c-mendaftarkan-endpoint--rute-modul-erp)
   - [D. Memasang Plugin: Rate Limiting (Redis), CORS, dan JWT](#d-memasang-plugin-rate-limiting-redis-cors-dan-jwt)
5. [Konfigurasi Deklaratif Alternatif (decK / `kong.yaml`)](#5-konfigurasi-deklaratif-alternatif-deck--kongyaml)
6. [Pengujian Endpoint & Integrasi Frontend (SvelteKit 2)](#6-pengujian-endpoint--integrasi-frontend-sveltekit-2)
7. [Troubleshooting & Pemeliharaan](#7-troubleshooting--pemeliharaan)

---

## 1. Arsitektur Kong API Gateway di ERP Monolith

Kong bertindak sebagai gerbang terpusat (*Reverse Proxy & API Gateway*) yang melindungi backend Echo v5:

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
└─────────────┬───────────────────────────┬──────────────┘
              │                           │
   (Direct :5432 Internal State)          │ Reverse Proxy Forward
              │                           ▼
┌─────────────▼──────────┐   ┌───────────────────────────┐
│ PostgreSQL (kong_db)   │   │ Echo v5 Backend Monolith  │
│ Port 5432 (Bukan 6432) │   │ Port 8080 (:8080)         │
└────────────────────────┘   └─────────────┬─────────────┘
                                           │
                              (pgxpool via PgBouncer :6432)
                                           ▼
                             ┌───────────────────────────┐
                             │ PostgreSQL (erp_db)       │
                             └───────────────────────────┘
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
```bash
curl -1sLf "https://packages.konghq.com/public/gateway-3x/setup.deb.sh" | sudo -E bash
```

### Langkah 3: Pasang Paket Kong
```bash
sudo apt update
sudo apt install -y kong
```

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
- **Service**: Merepresentasikan aplikasi backend upstream yang dituju (yaitu backend Echo v5 di `http://127.0.0.1:8080`).
- **Route**: Aturan pencocokan (*path matching*) yang menentukan ke Service mana request dari client akan diarahkan.
- **Plugin**: Middleware fungsional tambahan (Rate Limiting, CORS, JWT Auth, Logging) yang dapat dipasang ke Service atau Route tertentu.

---

### B. Mendaftarkan Service Backend Echo v5

Daftarkan aplikasi backend monolith sebagai sebuah Service di Kong:

```bash
curl -i -X POST http://localhost:8001/services \
  --data "name=erp-backend-service" \
  --data "url=http://127.0.0.1:8080" \
  --data "connect_timeout=60000" \
  --data "write_timeout=60000" \
  --data "read_timeout=60000"
```

---

### C. Mendaftarkan Endpoint / Rute Modul ERP

Setiap endpoint modul ERP didaftarkan sebagai **Route** di bawah service `erp-backend-service`.

> [!NOTE]
> Parameter `strip_path=false` memastikan bahwa prefix path URL (misal `/api/v1/acc`) **tetap dikirim secara utuh** ke Echo v5 router.

#### 1. Endpoint Health Check
```bash
curl -i -X POST http://localhost:8001/services/erp-backend-service/routes \
  --data "name=route-health" \
  --data "paths[]=/health" \
  --data "strip_path=false" \
  --data "methods[]=GET"
```

#### 2. Endpoint Autentikasi & Pengguna (19_USR, 21_ADM)
```bash
curl -i -X POST http://localhost:8001/services/erp-backend-service/routes \
  --data "name=route-auth" \
  --data "paths[]=/api/v1/auth" \
  --data "strip_path=false"
```

#### 3. Endpoint Modul Akuntansi & Keuangan (1_ACC, 4_GL, 2_AP, 3_AR)
```bash
curl -i -X POST http://localhost:8001/services/erp-backend-service/routes \
  --data "name=route-accounting" \
  --data "paths[]=/api/v1/acc" \
  --data "strip_path=false"
```

#### 4. Endpoint Modul Persediaan & Pergudangan (9_INV)
```bash
curl -i -X POST http://localhost:8001/services/erp-backend-service/routes \
  --data "name=route-inventory" \
  --data "paths[]=/api/v1/inv" \
  --data "strip_path=false"
```

#### 5. Endpoint Modul Penjualan & Kasir POS (10_SAL, 12_POS)
```bash
curl -i -X POST http://localhost:8001/services/erp-backend-service/routes \
  --data "name=route-sales" \
  --data "paths[]=/api/v1/sal" \
  --data "strip_path=false"
```

#### 6. Endpoint Modul Pembelian & Pengadaan (8_PUR)
```bash
curl -i -X POST http://localhost:8001/services/erp-backend-service/routes \
  --data "name=route-purchasing" \
  --data "paths[]=/api/v1/pur" \
  --data "strip_path=false"
```

#### 7. Endpoint Modul SDM & Penggajian (15_HRM, 16_PAY, 17_ATT)
```bash
curl -i -X POST http://localhost:8001/services/erp-backend-service/routes \
  --data "name=route-hrm" \
  --data "paths[]=/api/v1/hrm" \
  --data "strip_path=false"
```

#### 8. Catch-All Route untuk Seluruh Endpoint API Lainnya
```bash
curl -i -X POST http://localhost:8001/services/erp-backend-service/routes \
  --data "name=route-api-fallback" \
  --data "paths[]=/api" \
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
  - name: erp-backend-service
    url: http://127.0.0.1:8080
    connect_timeout: 60000
    write_timeout: 60000
    read_timeout: 60000
    routes:
      - name: route-health
        paths:
          - /health
        strip_path: false
        methods:
          - GET
      - name: route-auth
        paths:
          - /api/v1/auth
        strip_path: false
      - name: route-accounting
        paths:
          - /api/v1/acc
        strip_path: false
      - name: route-inventory
        paths:
          - /api/v1/inv
        strip_path: false
      - name: route-sales
        paths:
          - /api/v1/sal
        strip_path: false
      - name: route-purchasing
        paths:
          - /api/v1/pur
        strip_path: false
      - name: route-hrm
        paths:
          - /api/v1/hrm
        strip_path: false
      - name: route-api-fallback
        paths:
          - /api
        strip_path: false
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

2. Panggil endpoint `/health` melalui Kong Proxy Port `:8000`:
   ```bash
   curl -i http://localhost:8000/health
   ```
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
