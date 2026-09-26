# 🐘 Panduan Instalasi & Penggunaan PostgreSQL 16+ — Enterprise ERP Monolith

Dokumen ini adalah panduan operasional resmi untuk menginstal, mengonfigurasi, dan menggunakan **PostgreSQL 16+** pada sistem operasi **Ubuntu 22.04 / 24.04 LTS** sebagai basis data relasional utama untuk **Enterprise ERP Monolith**.

---

## 📑 Daftar Isi
1. [Arsitektur & Aturan Basis Data ERP](#1-arsitektur--aturan-basis-data-erp)
2. [Instalasi PostgreSQL 16 di Ubuntu](#2-instalasi-postgresql-16-di-ubuntu)
3. [Konfigurasi User & Database Proyek](#3-konfigurasi-user--database-proyek)
4. [Konfigurasi Jaringan & Autentikasi (`postgresql.conf` & `pg_hba.conf`)](#4-konfigurasi-jaringan--autentikasi)
5. [Penggunaan di Proyek ERP Monolith](#5-penggunaan-di-proyek-erp-monolith)
6. [Migrasi Skema (golang-migrate Direct Connection)](#6-migrasi-skema-golang-migrate-direct-connection)
7. [Aturan Skema & Konvensi Modul](#7-aturan-skema--konvensi-modul)
8. [Pemeliharaan, Backup & Troubleshooting](#8-pemeliharaan-backup--troubleshooting)

---

## 1. Arsitektur & Aturan Basis Data ERP

Pada arsitektur ERP Monolith ini, PostgreSQL memiliki dua jalur koneksi yang **harus dipisahkan**:

```text
[ Echo v5 App Backend ]  ──(Port 6432: Simple Protocol)──>  [ PgBouncer Proxy ]  ──> [ PostgreSQL :5432 ]
                                                                                         ▲
[ golang-migrate (DDL) ] ──(Port 5432: Direct TCP Connection)────────────────────────────┘
```

> [!IMPORTANT]
> - **Aplikasi Backend Go (`pgxpool`)** terhubung melalui **PgBouncer** pada port `:6432` dengan mode `transaction pooling`.
> - **Tool Migrasi (`golang-migrate`)** terhubung **LANGSUNG** ke PostgreSQL pada port `:5432` (`DIRECT_DATABASE_URL`). Dilarang menjalankan migrasi skema DDL melalui PgBouncer karena PgBouncer mode transaksi tidak mendukung advisory locks dan transactional DDL locking.
> - **Presisi Finansial**: Seluruh kolom nominal uang, kuantitas, harga pokok, tarif diskon, dan persentase pajak **WAJIB** menggunakan tipe data `NUMERIC(18, 4)` atau `NUMERIC(15, 2)`. Dilarang keras menggunakan `FLOAT` atau `DOUBLE PRECISION`.

---

## 2. Instalasi PostgreSQL 16 di Ubuntu

Gunakan repositori resmi PostgreSQL Global Development Group (PGDG) untuk memastikan mendapatkan PostgreSQL versi 16 terbaru.

### Langkah 1: Update Sistem & Pasang Dependensi
```bash
sudo apt update && sudo apt install -y curl ca-certificates gnupg lsb-release
```

### Langkah 2: Tambahkan GPG Key Resmi PostgreSQL
```bash
sudo install -d /etc/apt/keyrings
curl -fsSL https://www.postgresql.org/media/keys/ACCC4CF8.asc | sudo gpg --dearmor -o /etc/apt/keyrings/postgresql.gpg
```

### Langkah 3: Tambahkan Repositori PGDG ke APT Sources
```bash
echo "deb [signed-by=/etc/apt/keyrings/postgresql.gpg] http://apt.postgresql.org/pub/repos/apt $(lsb_release -cs)-pgdg main" | sudo tee /etc/apt/sources.list.d/pgdg.list
```

### Langkah 4: Pasang Paket PostgreSQL 16 & Contrib
```bash
sudo apt update
sudo apt install -y postgresql-16 postgresql-contrib-16
```

### Langkah 5: Verifikasi Status Service
```bash
sudo systemctl status postgresql
sudo systemctl enable postgresql
```

Pastikan output menampilkan status `active (running)`.

---

## 3. Konfigurasi User & Database Proyek

Sesuai template konfigurasi lingkungan proyek (`backend/.env.example`), kredensial standar yang digunakan untuk ERP adalah:
- **Database**: `erp_db`
- **Username**: `erp_user`
- **Password**: `erp_secret`

### Buat User & Database via `psql`

Jalankan perintah berikut menggunakan user sistem `postgres`:

```bash
sudo -u postgres psql << 'EOF'
-- 1. Buat User ERP dengan password terenkripsi SCRAM-SHA-256
CREATE USER erp_user WITH PASSWORD 'erp_secret';

-- 2. Berikan izin pembuatan database (opsional)
ALTER USER erp_user CREATEDB;

-- 3. Buat Database erp_db dengan kepemilikan erp_user
CREATE DATABASE erp_db OWNER erp_user ENCODING 'UTF8' LC_COLLATE 'en_US.UTF-8' LC_CTYPE 'en_US.UTF-8';

-- 4. Hubungkan ke database erp_db dan aktifkan ekstensi penting
\c erp_db

-- Ekstensi UUID generator
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Ekstensi performa & audit (opsional untuk reporting)
CREATE EXTENSION IF NOT EXISTS "pg_trgm";
CREATE EXTENSION IF NOT EXISTS "btree_gist";

-- Berikan semua hak akses pada skema public
GRANT ALL PRIVILEGES ON SCHEMA public TO erp_user;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO erp_user;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO erp_user;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO erp_user;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO erp_user;

EOF
```

---

## 4. Konfigurasi Jaringan & Autentikasi

### A. Konfigurasi `postgresql.conf`
Lokasi berkas: `/etc/postgresql/16/main/postgresql.conf`

Buka berkas konfigurasi:
```bash
sudo nano /etc/postgresql/16/main/postgresql.conf
```

Sesuaikan parameter berikut:
```ini
# Izinkan koneksi dari localhost (atau '*' jika PostgreSQL berada di mesin terpisah)
listen_addresses = 'localhost'
port = 5432

# Manajemen Memori (Penyesuaian rekomendasi untuk server ERP)
max_connections = 100
shared_buffers = 1GB                  # Sesuaikan 25% dari total RAM
effective_cache_size = 3GB            # Sesuaikan 75% dari total RAM
work_mem = 16MB
maintenance_work_mem = 256MB
min_wal_size = 1GB
max_wal_size = 4GB

# Autentikasi modern
password_encryption = scram-sha-256
```

### B. Konfigurasi `pg_hba.conf`
Lokasi berkas: `/etc/postgresql/16/main/pg_hba.conf`

Buka berkas akses kontrol:
```bash
sudo nano /etc/postgresql/16/main/pg_hba.conf
```

Tambahkan baris berikut agar `erp_user` dan PgBouncer dapat mengautentikasi melalui loopback TCP:
```ini
# TYPE  DATABASE        USER            ADDRESS                 METHOD
local   all             postgres                                peer
local   all             all                                     md5
host    all             all             127.0.0.1/32            scram-sha-256
host    all             all             ::1/128                 scram-sha-256
```

### C. Terapkan Perubahan
```bash
sudo systemctl restart postgresql
```

### D. Uji Koneksi Langsung
```bash
PGPASSWORD=erp_secret psql -h 127.0.0.1 -p 5432 -U erp_user -d erp_db -c "SELECT version();"
```

---

## 5. Penggunaan di Proyek ERP Monolith

### A. Pengaturan Variabel Lingkungan (`backend/.env`)

Pastikan berkas `backend/.env` memuat dua URL database berikut:

```env
# Koneksi Aplikasi Operasional (Via PgBouncer Port 6432)
DATABASE_URL=postgres://erp_user:erp_secret@localhost:6432/erp_db?sslmode=disable

# Koneksi Langsung Migrasi DDL (Direct PostgreSQL Port 5432)
DIRECT_DATABASE_URL=postgres://erp_user:erp_secret@localhost:5432/erp_db?sslmode=disable

# Parameter PgBouncer Pool di backend (pgxpool)
DB_MAX_CONNS=50
DB_MIN_CONNS=10
DB_MAX_CONN_LIFETIME=1h
DB_MAX_CONN_IDLE_TIME=30m
```

### B. Inisialisasi Koneksi pgxpool di Kode Go

Kode inisialisasi koneksi berada di [postgres.go](file:///opt/dev/erp_monolith/backend/internal/platform/database/postgres.go).

```go
package database

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(ctx context.Context, cfg Config) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(cfg.URL)
	if err != nil {
		return nil, err
	}

	// WAJIB: Mode query sederhana agar kompatibel dengan transaction pooling PgBouncer
	poolConfig.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	return pgxpool.NewWithConfig(ctx, poolConfig)
}
```

---

## 6. Migrasi Skema (golang-migrate Direct Connection)

Gunakan tool resmi **golang-migrate** untuk seluruh pembuatan tabel dan perubahan skema DDL.

### Instalasi golang-migrate CLI di Ubuntu:
```bash
curl -L -s https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz
sudo mv migrate /usr/local/bin/migrate
migrate -version
```

### Menjalankan Migrasi ke PostgreSQL Direct:

> [!CAUTION]
> Selalu jalankan migrasi ke port `:5432` (`DIRECT_DATABASE_URL`), bukan port `:6432`!

Masuk ke folder `backend/`:
```bash
cd /opt/dev/erp_monolith/backend

# Jalankan semua migrasi naik (Up)
migrate -path migrations -database "postgres://erp_user:erp_secret@localhost:5432/erp_db?sslmode=disable" up

# Rollback 1 migrasi terakhir (Down)
migrate -path migrations -database "postgres://erp_user:erp_secret@localhost:5432/erp_db?sslmode=disable" down 1

# Memeriksa versi skema saat ini
migrate -path migrations -database "postgres://erp_user:erp_secret@localhost:5432/erp_db?sslmode=disable" version

# Paksa reset state jika terjadi dirty migration (misal pada versi 1)
migrate -path migrations -database "postgres://erp_user:erp_secret@localhost:5432/erp_db?sslmode=disable" force 1
```

---

## 7. Aturan Skema & Konvensi Modul

Setiap tabel di dalam ERP Monolith ini **wajib** mengikuti pedoman berikut:

### 1. Prefiks Nama Tabel Berdasarkan Bounded Context
| Domain | Prefiks Tabel | Contoh Modul |
|---|---|---|
| **System & Admin** | `adm_`, `sys_` | `adm_companies`, `adm_branches`, `sys_users` |
| **Akuntansi & Keuangan** | `acc_` | `acc_chart_of_accounts`, `acc_journal_entries`, `acc_ledger` |
| **Persediaan & Pergudangan** | `inv_` | `inv_items`, `inv_stock_balances`, `inv_warehouses` |
| **Penjualan & Kasir** | `sal_`, `pos_` | `sal_orders`, `sal_invoices`, `pos_shifts` |
| **Pembelian & Vendor** | `pur_` | `pur_orders`, `pur_invoices`, `pur_vendors` |
| **SDM & Penggajian** | `hrm_`, `pay_` | `hrm_employees`, `pay_slips`, `hrm_attendance` |

### 2. Kolom Standar Multi-Tenant & Audit Trail Wajib
```sql
CREATE TABLE IF NOT EXISTS acc_invoices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    -- Kolom Multi-Tenant Wajib
    company_id UUID NOT NULL REFERENCES adm_companies(id) ON DELETE CASCADE,
    branch_id UUID NOT NULL REFERENCES adm_branches(id) ON DELETE RESTRICT,
    
    -- Field Bisnis
    invoice_number VARCHAR(100) NOT NULL,
    total_amount NUMERIC(18, 4) NOT NULL DEFAULT 0.0000, -- WAJIB NUMERIC, HARAM FLOAT
    
    -- Kolom Audit Trail Wajib
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by UUID NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by UUID NOT NULL,

    CONSTRAINT uq_inv_number_per_company UNIQUE (company_id, invoice_number)
);

-- Indeks Multi-Tenant Wajib
CREATE INDEX idx_acc_invoices_tenant ON acc_invoices(company_id, branch_id);
```

### 3. Generate Kode Go Type-Safe Menggunakan `sqlc`
Setelah menulis query SQL di folder internal modul, jalankan:
```bash
cd /opt/dev/erp_monolith/backend
sqlc generate
```

---

## 8. Pemeliharaan, Backup & Troubleshooting

### Cadangkan Basis Data (Backup):
```bash
# Backup penuh terkompresi
pg_dump -h 127.0.0.1 -p 5432 -U erp_user -F c -b -v -f erp_db_$(date +%Y%m%d_%H%M%S).dump erp_db

# Backup plain SQL
pg_dump -h 127.0.0.1 -p 5432 -U erp_user --schema-only -f erp_db_schema.sql erp_db
```

### Pemulihan Basis Data (Restore):
```bash
pg_restore -h 127.0.0.1 -p 5432 -U erp_user -d erp_db -v erp_db_20260926_160000.dump
```

### Monitoring Koneksi Aktif:
```sql
SELECT pid, usename, client_addr, application_name, state, query, age(clock_timestamp(), query_start) 
FROM pg_stat_activity 
WHERE datname = 'erp_db' 
ORDER BY query_start DESC;
```

### Mematikan Kueri yang Menggantung (Hung Query):
```sql
-- Membatalkan kueri tertentu
SELECT pg_cancel_backend(<pid>);

-- Memutus paksa koneksi tertentu
SELECT pg_terminate_backend(<pid>);
```
