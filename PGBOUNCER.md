# ⚡ Panduan Instalasi & Penggunaan PgBouncer — Enterprise ERP Monolith

Dokumen ini adalah panduan resmi untuk menginstal, mengonfigurasi, dan mengoperasikan **PgBouncer** pada sistem operasi **Ubuntu 22.04 / 24.04 LTS** sebagai **Connection Pooler Proxy** transaksi performa tinggi di depan basis data PostgreSQL untuk **Enterprise ERP Monolith**.

---

## 📑 Daftar Isi
1. [Arsitektur & Mengapa PgBouncer Wajib](#1-arsitektur--mengapa-pgbouncer-wajib)
2. [Instalasi PgBouncer di Ubuntu](#2-instalasi-pgbouncer-di-ubuntu)
3. [Konfigurasi Lengkap PgBouncer (`pgbouncer.ini` & `userlist.txt`)](#3-konfigurasi-lengkap-pgbouncer)
4. [Menjalankan & Mengaktifkan Layanan](#4-menjalankan--mengaktifkan-layanan)
5. [Penggunaan di Proyek ERP Monolith](#5-penggunaan-di-proyek-erp-monolith)
6. [Kesesuaian Driver `pgx/v5` & Mode Simple Protocol](#6-kesesuaian-driver-pgxv5--mode-simple-protocol)
7. [Aturan Pemisahan Port: 6432 (App) vs 5432 (DDL Migrate)](#7-aturan-pemisahan-port)
8. [Konsol Administrasi & Pemantauan (Monitoring)](#8-konsol-administrasi--pemantauan)
9. [Troubleshooting & Solusi Masalah Umum](#9-troubleshooting--solusi-masalah-umum)

---

## 1. Arsitektur & Mengapa PgBouncer Wajib

PostgreSQL mengalokasikan 1 proses sistem operasi (OS process) per koneksi client, dengan kebutuhan RAM berkisar 5–10 MB per koneksi. Pada sistem ERP berskala enterprise dengan banyak modul (akuntansi, kasir POS, gudang, payroll) dan ratusan user aktif konkuren, ribuan koneksi langsung ke PostgreSQL akan menyebabkan lonjakan memori (*OOM - Out of Memory*) dan penurunan drastis performa (*connection exhaustion*).

**PgBouncer** bertindak sebagai perantara ringan (*lightweight connection pooler*):

```text
┌──────────────────────────────────────────────┐
│  Echo v5 Backend Apps (Ribuan Goroutine)    │
└──────────────────────┬───────────────────────┘
                       │ TCP (:6432)
                       ▼
┌──────────────────────────────────────────────┐
│  PgBouncer Connection Pooler Proxy           │
│  Mode: TRANSACTION POOLING                   │
└──────────────────────┬───────────────────────┘
                       │ TCP (:5432) - Multiplexing ke 20-50 Koneksi Fisik
                       ▼
┌──────────────────────────────────────────────┐
│  PostgreSQL 16 Database Server               │
└──────────────────────────────────────────────┘
```

### Karakteristik Transaction Pooling:
- Koneksi fisik PostgreSQL hanya dipinjamkan ke client **selama transaksi berjalan** (`BEGIN ... COMMIT / ROLLBACK`).
- Begitu transaksi selesai, koneksi fisik langsung dikembalikan ke pool untuk digunakan request lain.
- Menghemat memori server secara radikal dan memungkinkan ribuan koneksi client secara bersamaan.

---

## 2. Instalasi PgBouncer di Ubuntu

### Langkah 1: Update & Pasang Paket
```bash
sudo apt update
sudo apt install -y pgbouncer
```

### Langkah 2: Periksa Versi yang Terpasang
```bash
pgbouncer -V
```

---

## 3. Konfigurasi Lengkap PgBouncer

Konfigurasi PgBouncer berada di direktori `/etc/pgbouncer/`.

### A. Berkas Konfigurasi Utama: `/etc/pgbouncer/pgbouncer.ini`

Cadangkan berkas default dan buat konfigurasi baru:
```bash
sudo cp /etc/pgbouncer/pgbouncer.ini /etc/pgbouncer/pgbouncer.ini.bak
sudo nano /etc/pgbouncer/pgbouncer.ini
```

Salin dan simpan konfigurasi produksi berikut:

```ini
;; ========================================================
;; Enterprise ERP Monolith — PgBouncer Configuration
;; ========================================================

[databases]
;; Pemetaan database ERP Monolith ke PostgreSQL direct port 5432
erp_db = host=127.0.0.1 port=5432 dbname=erp_db auth_user=erp_user

;; Fallback database wildcard (opsional)
* = host=127.0.0.1 port=5432

[pgbouncer]
;; ----------------------------------------
;; Pengaturan Jaringan & Port
;; ----------------------------------------
logfile = /var/log/postgresql/pgbouncer.log
pidfile = /var/run/postgresql/pgbouncer.pid
listen_addr = 0.0.0.0
listen_port = 6432
unix_socket_dir = /var/run/postgresql

;; ----------------------------------------
;; Autentikasi
;; ----------------------------------------
auth_type = scram-sha-256
auth_file = /etc/pgbouncer/userlist.txt
;; Mengizinkan user pgbouncer membaca auth query langsung dari postgres (opsional):
;; auth_query = SELECT usename, passwd FROM pg_shadow WHERE usename=$1

;; ----------------------------------------
;; Aturan Connection Pooling (KRITIKAL)
;; ----------------------------------------
;; WAJIB 'transaction' untuk efisiensi maksimum di ERP Monolith
pool_mode = transaction

;; Batas maksimal koneksi client dari aplikasi Go / frontend
max_client_conn = 5000

;; Jumlah koneksi fisik PostgreSQL default per (user, database)
default_pool_size = 30

;; Jumlah minimum koneksi standby
min_pool_size = 5

;; Cadangan koneksi instan saat lonjakan beban transaksi tinggi
reserve_pool_size = 10
reserve_pool_timeout = 3.0

;; Batas absolut total koneksi fisik ke server PostgreSQL
max_db_connections = 60

;; ----------------------------------------
;; Timeouts & Housekeeping
;; ----------------------------------------
server_reset_query = DISCARD ALL
server_check_delay = 30
server_check_query = select 1
server_idle_timeout = 600
client_idle_timeout = 0
query_timeout = 120

;; Hindari konflik startup parameter pada driver pgx/PostgreSQL
ignore_startup_parameters = extra_float_digits, application_name, search_path

;; User yang memiliki akses ke konsol admin pgbouncer
admin_users = postgres, erp_user
stats_users = postgres, erp_user
```

---

### B. Berkas Autentikasi: `/etc/pgbouncer/userlist.txt`

PgBouncer membutuhkan daftar user dan password yang diizinkan untuk login.

Buka atau buat berkas `/etc/pgbouncer/userlist.txt`:
```bash
sudo nano /etc/pgbouncer/userlist.txt
```

Isi dengan kredensial `erp_user` dan `postgres`:
```text
"erp_user" "erp_secret"
"postgres" "postgres"
```

> [!TIP]
> Jika Anda menggunakan enkripsi `scram-sha-256` dari PostgreSQL, Anda dapat mengekstrak hash langsung dari tabel `pg_shadow`:
> ```bash
> sudo -u postgres psql -t -A -c "SELECT '\"' || usename || '\" \"' || passwd || '\"' FROM pg_shadow WHERE usename IN ('erp_user', 'postgres');" | sudo tee /etc/pgbouncer/userlist.txt
> ```

### C. Amankan Hak Akses File
```bash
sudo chown -R postgres:postgres /etc/pgbouncer
sudo chmod 640 /etc/pgbouncer/userlist.txt
sudo chmod 644 /etc/pgbouncer/pgbouncer.ini
```

---

## 4. Menjalankan & Mengaktifkan Layanan

Pastikan PgBouncer diatur berjalan otomatis saat server booting:

```bash
# Aktifkan dan jalankan PgBouncer
sudo systemctl enable pgbouncer
sudo systemctl restart pgbouncer

# Periksa status
sudo systemctl status pgbouncer
```

### Verifikasi Port 6432 Terbuka:
```bash
ss -tulpn | grep 6432
```

### Uji Koneksi ke PgBouncer:
```bash
PGPASSWORD=erp_secret psql -h 127.0.0.1 -p 6432 -U erp_user -d erp_db -c "SELECT 1 as pgbouncer_connected;"
```

Jika output menghasilkan `pgbouncer_connected: 1`, PgBouncer telah aktif dan terhubung sempurna ke PostgreSQL!

---

## 5. Penggunaan di Proyek ERP Monolith

### Konfigurasi Lingkungan (`backend/.env`)

Aplikasi Go ERP Monolith harus diarahkan ke port **6432** (PgBouncer), bukan 5432:

```env
# Koneksi Aplikasi Go via PgBouncer (Port 6432)
DATABASE_URL=postgres://erp_user:erp_secret@localhost:6432/erp_db?sslmode=disable

# Konfigurasi Pool pgxpool di sisi Golang
DB_MAX_CONNS=50
DB_MIN_CONNS=10
DB_MAX_CONN_LIFETIME=1h
DB_MAX_CONN_IDLE_TIME=30m
```

---

## 6. Kesesuaian Driver `pgx/v5` & Mode Simple Protocol

Dalam arsitektur PgBouncer dengan `pool_mode = transaction`:
- Setiap query dalam transaksi dapat berpindah ke backend connection PostgreSQL yang berbeda secara transparan.
- Fitur bawaan PostgreSQL *server-side prepared statements* terikat pada *session connection*.
- Jika driver client mencoba mengeksekusi prepared statement bernama (*named prepared statement*) yang tersimpan di koneksi A namun dikirim ke koneksi B, PostgreSQL akan mengembalikan error fatal:
  `ERROR: prepared statement "stmt_cache_xyz" does not exist`.

### Solusi Wajib di Proyek Ini:
Pada file [postgres.go](file:///opt/dev/erp_monolith/backend/internal/platform/database/postgres.go), koneksi pgxpool dikonfigurasi menggunakan:

```go
// internal/platform/database/postgres.go
poolConfig.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol
```

Mode `QueryExecModeSimpleProtocol` mengirimkan kueri langsung via protocol sederhana tanpa server-side caching prepared statement yang bentrok, sehingga 100% aman dan berkecepatan tinggi saat melalui PgBouncer transaction pooler.

---

## 7. Aturan Pemisahan Port

| Operasi | Port Target | URL Lingkungan | Alasan |
|---|---|---|---|
| **Aplikasi ERP Backend (`Echo v5`)** | `:6432` *(PgBouncer)* | `DATABASE_URL` | Menangani ribuan request transaksi user simultan dengan efisiensi memori tinggi. |
| **Migrasi DDL (`golang-migrate`)** | `:5432` *(Direct Postgres)* | `DIRECT_DATABASE_URL` | Migrasi DDL memerlukan *Advisory Lock* dan *Session Lock* pada tabel `schema_migrations`, yang dilarang pada transaction pooling. |
| **Kong API Gateway** | `:5432` *(Direct Postgres)* | N/A (Internal Kong) | Kong membutuhkan koneksi database stateful untuk bootstrapping dan sync konfigurasi internal. |

### Menjalankan Migrasi DDL yang Tepat:
```bash
cd /opt/dev/erp_monolith/backend

# BENAR: Menggunakan port 5432 direct
migrate -path migrations -database "postgres://erp_user:erp_secret@localhost:5432/erp_db?sslmode=disable" up

# SALAH: DILARANG menghubungkan migrate ke port 6432!
```

---

## 8. Konsol Administrasi & Pemantauan (Monitoring)

PgBouncer memiliki database virtual bernama `pgbouncer` untuk memantau performa pool secara realtime.

### Masuk ke Virtual Admin Console:
```bash
psql -h 127.0.0.1 -p 6432 -U postgres pgbouncer
```

### Perintah Penting di Admin Console:

#### 1. Melihat Status Pool (`SHOW POOLS;`)
```sql
SHOW POOLS;
```
Menampilkan metrik:
- `cl_active`: Client connections yang sedang aktif mengeksekusi query.
- `cl_waiting`: Client connections yang sedang antre menunggu koneksi database kosong.
- `sv_active`: Koneksi fisik ke PostgreSQL yang sedang memproses query.
- `sv_idle`: Koneksi fisik yang siap digunakan seketika.

#### 2. Melihat Statistik Query & Transaksi (`SHOW STATS;`)
```sql
SHOW STATS;
```
Menampilkan jumlah transaksi total, durasi rata-rata transaksi, total byte jaringan terkirim/diterima.

#### 3. Melihat Daftar Client yang Terhubung (`SHOW CLIENTS;`)
```sql
SHOW CLIENTS;
```

#### 4. Reload Konfigurasi Tanpa Downtime:
Jika Anda mengubah berkas `pgbouncer.ini` atau `userlist.txt`:
```sql
RELOAD;
```
Atau dari terminal Ubuntu:
```bash
sudo systemctl reload pgbouncer
```

#### 5. Menjeda & Melanjutkan Koneksi (Maintenance):
```sql
PAUSE;   -- Menunggu query selesai lalu menahan query baru
RESUME;  -- Melanjutkan koneksi normal
```

---

## 9. Troubleshooting & Solusi Masalah Umum

### 1. `ERROR: prepared statement does not exist`
- **Penyebab**: Aplikasi mengirim query menggunakan prepared statement bernama pada mode transaction pooling.
- **Solusi**: Pastikan driver `pgxpool` di [postgres.go](file:///opt/dev/erp_monolith/backend/internal/platform/database/postgres.go) menyertakan:
  ```go
  poolConfig.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol
  ```

### 2. `server login failed: password authentication failed`
- **Penyebab**: Hash kredensial di `/etc/pgbouncer/userlist.txt` tidak cocok dengan password `erp_user` di PostgreSQL.
- **Solusi**: Perbarui password di `userlist.txt` atau ekspor ulang dari `pg_shadow`, lalu jalankan `sudo systemctl reload pgbouncer`.

### 3. `no more connections allowed (max_client_conn)`
- **Penyebab**: Jumlah client koneksi melebihi kuota `max_client_conn`.
- **Solusi**: Naikkan `max_client_conn = 10000` di `/etc/pgbouncer/pgbouncer.ini` dan periksa apakah ada koneksi aplikasi yang mengalami leak / lupa `defer dbPool.Close()`.

### 4. `golang-migrate error: pq: transaction pooling not supported`
- **Penyebab**: Tool migrasi dijalankan mengarah ke port `6432`.
- **Solusi**: Ubah koneksi migrasi ke `DIRECT_DATABASE_URL` pada port `5432`.
