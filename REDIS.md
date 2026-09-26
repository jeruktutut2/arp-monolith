# 🚀 Panduan Instalasi & Penggunaan Redis — Enterprise ERP Monolith

Dokumen ini adalah panduan resmi untuk menginstal, mengonfigurasi, dan mengoperasikan **Redis 7+** pada sistem operasi **Ubuntu 22.04 / 24.04 LTS** sebagai basis data *in-memory* performa tinggi untuk keperluan **Distributed Caching, Session/Token Blacklist, Distributed Locks, dan Background Job Queue** pada **Enterprise ERP Monolith**.

---

## 📑 Daftar Isi
1. [Peran & Kasus Penggunaan Redis di ERP Monolith](#1-peran--kasus-penggunaan-redis-di-erp-monolith)
2. [Instalasi Redis di Ubuntu](#2-instalasi-redis-di-ubuntu)
3. [Konfigurasi Produksi Redis (`/etc/redis/redis.conf`)](#3-konfigurasi-produksi-redis)
4. [Menjalankan & Menguji Koneksi Redis](#4-menjalankan--menguji-koneksi-redis)
5. [Penggunaan di Proyek ERP Monolith (Backend Golang)](#5-penggunaan-di-proyek-erp-monolith-backend-golang)
6. [Pola Implementasi Praktis di Modul ERP](#6-pola-implementasi-praktis-di-modul-erp)
   - [A. Cache-Aside untuk Data Master (COA, Item, Currency)](#a-cache-aside-untuk-data-master)
   - [B. Distributed Lock untuk Mencegah Race Condition Stok & Jurnal](#b-distributed-lock-stok--jurnal)
   - [C. JWT Blacklist / Token Revocation](#c-jwt-blacklist--token-revocation)
   - [D. Penyimpanan State Rate Limiting Kong Gateway](#d-penyimpanan-state-rate-limiting-kong-gateway)
7. [CLI Cheat-sheet & Pemantauan (Monitoring)](#7-cli-cheat-sheet--pemantauan)

---

## 1. Peran & Kasus Penggunaan Redis di ERP Monolith

Di dalam arsitektur ERP Monolith berkinerja tinggi, Redis berperan sebagai akselerator *in-memory* yang mengurangi beban kueri ke PostgreSQL:

1. **Master Data Caching**: Menyimpan Chart of Accounts (COA), daftar mata uang dan nilai tukar harian, master barang, serta preferensi cabang tenant yang sering dibaca berulang kali.
2. **Casbin RBAC & User Permissions**: Menyimpan matriks hak akses pengguna dan role secara in-memory untuk otorisasi sub-milidetik di setiap HTTP middleware.
3. **Session & JWT Token Blacklist**: Saat user melakukan logout atau ganti password, token JWT lama dimasukkan ke dalam Redis Blacklist hingga masa berlakunya berakhir (*TTL auto-expiration*).
4. **Distributed Locks (`Redlock` / `SET NX PX`)**: Mencegah *race condition* kritis pada:
   - Pemotongan kuantitas stok gudang saat ada pemesanan serentak di kasir POS / e-commerce.
   - Penutupan periode buku akuntansi (*Period Closing*).
   - Penomoran otomatis dokumen faktur faktual (*Sequencer Generator*).
5. **Asynchronous Job Queue (`asynq`)**: Penjadwalan antrean latar belakang seperti pencetakan ribuan faktur PDF, pengiriman slip gaji bulanan via email, dan kalkulasi laporan laba rugi berkala.
6. **Backend Rate Limiting untuk Kong Gateway**: Kong dapat menggunakan Redis terpusat untuk menghitung kuota request API per user/IP.

---

## 2. Instalasi Redis di Ubuntu

Gunakan repositori resmi Redis agar mendapatkan versi terbaru (Redis 7.x+).

### Langkah 1: Update & Pasang Dependensi
```bash
sudo apt update && sudo apt install -y curl ca-certificates lsb-release gpg
```

### Langkah 2: Tambahkan GPG Key Resmi Redis
```bash
curl -fsSL https://packages.redis.io/gpg | sudo gpg --dearmor -o /usr/share/keyrings/redis-archive-keyring.gpg
```

### Langkah 3: Tambahkan Repositori Resmi Redis ke APT
```bash
echo "deb [signed-by=/usr/share/keyrings/redis-archive-keyring.gpg] https://packages.redis.io/deb $(lsb_release -cs) main" | sudo tee /etc/apt/sources.list.d/redis.list
```

### Langkah 4: Pasang Paket Redis Server & CLI
```bash
sudo apt update
sudo apt install -y redis-server redis-tools
```

---

## 3. Konfigurasi Produksi Redis

Berkas konfigurasi utama terletak di `/etc/redis/redis.conf`.

Buka berkas konfigurasi:
```bash
sudo nano /etc/redis/redis.conf
```

Sesuaikan parameter berikut:

### 1. Integrasi Systemd & Jaringan
```ini
# Ubah nilai 'supervised no' menjadi 'systemd'
supervised systemd

# Amankan bind hanya ke loopback lokal (kecuali Redis di host terpisah)
bind 127.0.0.1 ::1
port 6379
protected-mode yes
```

### 2. Autentikasi Password Kuat
Cari baris `# requirepass foobared` dan ubah menjadi:
```ini
requirepass erp_redis_secret
```
*(Gunakan password yang kuat dan aman untuk lingkungan produksi)*

### 3. Batas Memori & Strategi Eviction (Cache Policy)
Tambahkan konfigurasi alokasi RAM dan kebijakan pembersihan data kedaluwarsa:
```ini
# Batas memori maksimum untuk Redis
maxmemory 512mb

# Hapus key yang memiliki waktu kedaluwarsa (TTL) berdasarkan Least Recently Used (LRU)
maxmemory-policy volatile-lru
```

### 4. Persistensi Data (RDB Snapshot + AOF Durability)
```ini
# Aktifkan Append Only File (AOF) untuk durabilitas data finansial/queue
appendonly yes
appendfsync everysec
```

---

## 4. Menjalankan & Menguji Koneksi Redis

### Langkah 1: Restart Layanan & Aktifkan Autostart
```bash
sudo systemctl restart redis-server
sudo systemctl enable redis-server
```

### Langkah 2: Periksa Status Layanan
```bash
sudo systemctl status redis-server
```

### Langkah 3: Uji Menggunakan `redis-cli`
```bash
# Uji ping dengan autentikasi password
redis-cli -a erp_redis_secret ping
```
*Output sukses: `PONG`*

### Langkah 4: Uji Simpan & Ambil Key
```bash
redis-cli -a erp_redis_secret set test_key "ERP Redis Berhasil Terkoneksi"
redis-cli -a erp_redis_secret get test_key
```

---

## 5. Penggunaan di Proyek ERP Monolith (Backend Golang)

### A. Pengaturan Variabel Lingkungan (`backend/.env`)

Tambahkan konfigurasi Redis berikut pada berkas `.env` proyek:

```env
# ========================================================
# Redis Cache & Queue Configuration
# ========================================================
REDIS_HOST=127.0.0.1
REDIS_PORT=6379
REDIS_PASSWORD=erp_redis_secret
REDIS_DB=0
REDIS_POOL_SIZE=20
```

### B. Library Resmi Go: `go-redis/v9`

Di dalam folder `backend/`, pastikan dependency terpasang:
```bash
cd /opt/dev/erp_monolith/backend
go get github.com/redis/go-redis/v9
```

### C. Implementasi Wrapper Platform (`internal/platform/cache/redis.go`)

Buat modul platform koneksi Redis yang reusable di `backend/internal/platform/cache/redis.go`:

```go
package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	Host     string
	Port     string
	Password string
	DB       int
	PoolSize int
}

type Client struct {
	RDB *redis.Client
}

// NewRedisClient menginisialisasi koneksi pool Redis
func NewRedisClient(cfg Config) (*Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: 5,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("gagal terhubung ke Redis: %w", err)
	}

	return &Client{RDB: rdb}, nil
}

func (c *Client) Close() error {
	return c.RDB.Close()
}
```

---

## 6. Pola Implementasi Praktis di Modul ERP

### A. Cache-Aside untuk Data Master

Untuk data master yang jarang berubah namun sering diakses (seperti kode akun COA atau nilai tukar valuta):

```go
package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"erp_monolith/backend/internal/modules/acc/domain"
	"github.com/redis/go-redis/v9"
)

type AccountRepository struct {
	redisClient *redis.Client
	// dbPool *pgxpool.Pool
}

func (r *AccountRepository) GetAccountByID(ctx context.Context, companyID, accountID string) (*domain.Account, error) {
	cacheKey := fmt.Sprintf("tenant:%s:acc:%s", companyID, accountID)

	// 1. Cek dari Cache Redis terlebih dahulu
	val, err := r.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var acc domain.Account
		if err := json.Unmarshal([]byte(val), &acc); err == nil {
			return &acc, nil // Cache Hit 🚀
		}
	}

	// 2. Cache Miss: Ambil dari PostgreSQL via PgBouncer
	acc, err := r.fetchFromPostgres(ctx, companyID, accountID)
	if err != nil {
		return nil, err
	}

	// 3. Simpan kembali ke Redis dengan TTL 1 Jam
	data, _ := json.Marshal(acc)
	r.redisClient.Set(ctx, cacheKey, data, 1*time.Hour)

	return acc, nil
}
```

---

### B. Distributed Lock (Stok & Jurnal)

Mencegah dua kasir POS menjual 1 sisa barang terakhir secara bersamaan:

```go
package inventory

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// AcquireStockLock mengambil lock sementara untuk mutasi stok barang
func AcquireStockLock(ctx context.Context, rdb *redis.Client, warehouseID, itemID string, ttl time.Duration) (bool, string, error) {
	lockKey := fmt.Sprintf("lock:stock:%s:%s", warehouseID, itemID)
	lockValue := fmt.Sprintf("%d", time.Now().UnixNano())

	// SET key value NX (Not Exists) PX (TTL Milidetik)
	success, err := rdb.SetNX(ctx, lockKey, lockValue, ttl).Result()
	if err != nil {
		return false, "", err
	}

	return success, lockValue, nil
}

// ReleaseStockLock melepaskan lock dengan verifikasi token value
func ReleaseStockLock(ctx context.Context, rdb *redis.Client, warehouseID, itemID, lockValue string) error {
	lockKey := fmt.Sprintf("lock:stock:%s:%s", warehouseID, itemID)

	// Script Lua atomik untuk memastikan hanya pemilik lock yang dapat menghapus
	luaScript := `
		if redis.call("get", KEYS[1]) == ARGV[1] then
			return redis.call("del", KEYS[1])
		else
			return 0
		end
	`
	return rdb.Eval(ctx, luaScript, []string{lockKey}, lockValue).Err()
}
```

---

### C. JWT Blacklist / Token Revocation

Ketika pengguna logout sebelum masa berlaku JWT habis:

```go
package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// BlacklistToken memasukkan JWT token signature ke blacklist dengan TTL sisa masa berlaku token
func BlacklistToken(ctx context.Context, rdb *redis.Client, tokenSignature string, remainingTTL time.Duration) error {
	key := fmt.Sprintf("jwt:blacklist:%s", tokenSignature)
	return rdb.Set(ctx, key, "revoked", remainingTTL).Err()
}

// IsTokenBlacklisted memeriksa apakah token berada di daftar hitam
func IsTokenBlacklisted(ctx context.Context, rdb *redis.Client, tokenSignature string) bool {
	key := fmt.Sprintf("jwt:blacklist:%s", tokenSignature)
	exists, err := rdb.Exists(ctx, key).Result()
	return err == nil && exists > 0
}
```

---

### D. Penyimpanan State Rate Limiting Kong Gateway

Redis ini juga dapat langsung digunakan oleh **Kong API Gateway** untuk menyimpan counter request lintas kluster secara terpusat:
- **Host**: `127.0.0.1`
- **Port**: `6379`
- **Password**: `erp_redis_secret`
*(Panduan konfigurasi di sisi Kong terdapat pada file [KONG.md](file:///opt/dev/erp_monolith/KONG.md))*.

---

## 7. CLI Cheat-sheet & Pemantauan (Monitoring)

### Masuk ke CLI Interaktif:
```bash
redis-cli -a erp_redis_secret
```

### Perintah Penting:
```text
# Melihat informasi memori & CPU
INFO memory
INFO stats

# Mencari key tertentu (Gunakan SCAN di produksi, jangan gunakan KEYS *)
SCAN 0 MATCH "tenant:*" COUNT 50

# Melihat sisa masa hidup key (TTL dalam detik)
TTL "tenant:comp-123:acc:coa-001"

# Menghapus key tertentu
DEL "tenant:comp-123:acc:coa-001"

# Menghapus seluruh cache di database aktif (Hati-hati!)
FLUSHDB

# Memantau lalu lintas perintah yang masuk secara realtime (Debugging)
MONITOR
```
