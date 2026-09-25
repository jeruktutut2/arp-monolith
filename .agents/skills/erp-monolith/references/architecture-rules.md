# 🏛️ Aturan Arsitektur & Prinsip Rekayasa ERP Monolith

Dokumen ini mendefinisikan aturan keras (*hard rules*) yang harus ditaati oleh setiap agent dan pengembang saat menulis kode untuk proyek ERP Monolith.

---

## 1. Hexagonal Architecture & Batas Modul

1. **Package by Domain / Bounded Context (Heksagon per Modul)**:
   - Setiap modul bisnis berada di bawah `internal/modules/<nama_modul>`.
   - Modul memiliki isolasi internal berbasis port dan adapter:
     - `domain/`: Core Domain (Entitas, Value Object, Invarian, Domain Errors), Inbound Ports (Usecase Interfaces), dan Outbound Ports (Repository & Inter-module Contracts).
     - `usecase/`: Implementasi Inbound Ports (Logika bisnis, kalkulasi finansial, orkestrasi transaksi).
     - `repository/`: Driven Adapter untuk akses PostgreSQL via `sqlc` dan `pgx/v5` melalui PgBouncer.
     - `delivery/http/`: Driving Adapter untuk HTTP REST API menggunakan **Echo v5** (`github.com/labstack/echo/v5`), DTO, binding, dan validasi input.

2. **Dilarang Import Konkrit Lintas Modul**:
   - 🔴 **DILARANG**:
     ```go
     // Di dalam package internal/modules/sales/usecase
     import "erp_monolith/internal/modules/inventory/usecase" // SALAH!
     ```
   - 🟢 **BENAR (Consumer-Defined Interface / Outbound Port)**:
     ```go
     // Di dalam package internal/modules/sales/domain/contracts.go
     type InventoryStockService interface {
         ReserveStock(ctx context.Context, itemID string, qty decimal.Decimal) error
     }
     ```
     Implementasi interface diberikan oleh adapter modul `inventory`, dan dihubungkan (*injected*) di `cmd/server/main.go`.

3. **Event Bus untuk Aksi Lanjutan (*Side-Effects*)**:
   - Jika suatu aksi tidak memerlukan konsistensi ACID instan (misal: posting jurnal buku besar setelah pesanan disetujui, pengiriman email notifikasi, atau pembuatan log audit), **wajib** menggunakan `EventBus`.
   - Modul pengirim hanya menerbitkan event, tanpa mengetahui siapa penerimanya.

---

## 2. PostgreSQL & PgBouncer Connection Pooling

1. **Topologi PgBouncer**:
   - Aplikasi Go terhubung ke PgBouncer (:6432) menggunakan connection pool `pgxpool`.
   - PgBouncer dikonfigurasi dengan `pool_mode = transaction`.
2. **Kepatuhan Protokol Driver (`pgx/v5`)**:
   - Karena transaksi dapat berganti koneksi backend pada mode transaksi PgBouncer, driver `pgxpool` harus menghindari prepared statement global yang mengikat session. Gunakan mode simple protocol (`QueryExecModeSimpleProtocol`) atau nameless prepared statements.
3. **Migrasi Skema dengan `golang-migrate`**:
   - Eksekusi migrasi DDL database wajib menggunakan **`golang-migrate`** (`github.com/golang-migrate/migrate/v4`) dengan pasangan berkas `.up.sql` dan `.down.sql`.
   - Perintah migrasi dijalankan langsung ke port server PostgreSQL (:5432), bukan melalui pooler transaksi PgBouncer, untuk mendukung statement DDL transaksional secara penuh.

---

## 3. Standar Presisi Finansial & Angka

> [!CAUTION]
> **Dilarang keras menggunakan tipe data `float32` atau `float64` untuk seluruh representasi nominal uang, diskon, kuantitas stok, dan tarif pajak.**

- **Library Wajib**: `github.com/shopspring/decimal`
- **Operasi Matematika**:
  - Penjumlahan: `a.Add(b)`
  - Pengurangan: `a.Sub(b)`
  - Perkalian: `a.Mul(b)`
  - Pembagian: `a.Div(b)` atau `a.DivRound(b, 4)`
  - Pembulatan Akuntansi: `val.RoundBanker(2)` (Banker's Rounding / Half-to-Even)
- **Tipe Kolom Database**:
  - Mata uang / saldo: `NUMERIC(18, 4)` atau `NUMERIC(15, 2)`
  - Kuantitas stok: `NUMERIC(12, 4)`

---

## 4. Isolasi Data & Transaksi Database

1. **Kepemilikan Tabel Tunggal (*Single Table Ownership*)**:
   - Hanya modul pemilik yang boleh mengeksekusi query `INSERT`, `UPDATE`, `DELETE` pada tabel miliknya.
   - Modul lain tidak boleh melakukan direct SQL `JOIN` ke tabel modul lain.
2. **Denormalisasi Transaksional**:
   - Pada saat pembuatan transaksi (contoh: *Sales Order*), simpan snapshot data penting (misal: `customer_name`, `item_name`, `unit_price`) di tabel transaksi penjualan. Hal ini mencegah ketergantungan join dan menjaga keaslian riwayat transaksi masa lalu saat data master berubah.
3. **Database Transaction Boundaries**:
   - Transaksi database (`pgx.Tx`) hanya boleh melingkupi satu modul.
   - Transaksi lintas modul diselesaikan melalui saga / event-driven kompensasi jika diperlukan, bukan dengan memanjangkan satu database transaction melintasi batas domain.

---

## 5. Multi-Tenancy & Konteks Keamanan

1. Setiap entitas master dan transaksional wajib memuat:
   - `company_id UUID`: Identitas perusahaan/anak usaha.
   - `branch_id UUID`: Identitas cabang operasional.
2. Setiap query baca dan tulis wajib menyertakan filter `WHERE company_id = $1`.
3. Informasi aktor (siapa yang mengeksekusi) diekstrak dari `echo.Context` melalui middleware platform:
   - `audit.GetActor(ctx) -> { UserID, CompanyID, BranchID, IP, UserAgent }`

---

## 6. Standar Frontend SvelteKit 2 & Svelte 5

1. **Svelte 5 Runes**:
   - Gunakan `$state()` untuk variabel reaktif lokal.
   - Gunakan `$derived()` untuk kalkulasi turunan (computed state).
   - Gunakan `$props()` untuk deklarasi props komponen.
   - Gunakan `$effect()` untuk efek samping berbasis siklus hidup DOM.
2. **Layout & Integrasi Libs**:
   - Gunakan layout bawaan SvelteKit `+layout.svelte` untuk menyematkan App Shell (`0_HDR`, `0_SDB`, `0_FTR`).
   - Hubungkan pustaka pihak ketiga secara modular melalui action Svelte `use:action` atau Svelte wrappers untuk meminimalisasi re-render tidak perlu.
