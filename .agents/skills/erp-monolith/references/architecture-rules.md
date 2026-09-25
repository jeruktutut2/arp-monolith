# 🏛️ Aturan Arsitektur & Prinsip Rekayasa ERP Monolith

Dokumen ini mendefinisikan aturan keras (*hard rules*) yang harus ditaati oleh setiap agent dan pengembang saat menulis kode untuk proyek ERP Monolith.

---

## 1. Batas Modul & Larangan Circular Import

1. **Package by Domain / Bounded Context**:
   - Setiap modul bisnis berada di bawah `internal/modules/<modul_name>`.
   - Modul memiliki isolasi internal sendiri:
     - `domain/`: Entitas, Value Object, Domain Errors, dan Consumer-Defined Interfaces.
     - `usecase/`: Aturan bisnis, kalkulasi finansial, orkestrasi transaksi.
     - `repository/`: Akses database Postgres (sqlc / pgx).
     - `delivery/http/`: Handler HTTP Chi, DTO, binding, validasi input.

2. **Dilarang Import Konkrit Lintas Modul**:
   - 🔴 **DILARANG**:
     ```go
     // Di dalam package internal/modules/sales/usecase
     import "erp_monolith/internal/modules/inventory/usecase" // SALAH!
     ```
   - 🟢 **BENAR (Consumer-Defined Interface)**:
     ```go
     // Di dalam package internal/modules/sales/domain/contracts.go
     type InventoryStockService interface {
         ReserveStock(ctx context.Context, itemID string, qty decimal.Decimal) error
     }
     ```
     Implementasi interface diberikan oleh modul `inventory`, dan dihubungkan (*injected*) di `cmd/server/main.go`.

3. **Event Bus untuk Aksi Lanjutan (*Side-Effects*)**:
   - Jika suatu aksi tidak memerlukan konsistensi ACID instan (misal: posting jurnal buku besar setelah pesanan disetujui, pengiriman email notifikasi, atau pembuatan log audit), **wajib** menggunakan `EventBus`.
   - Modul pengirim hanya menerbitkan event, tanpa mengetahui siapa penerimanya.

---

## 2. Standar Presisi Finansial & Angka

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

## 3. Isolasi Data & Transaksi Database

1. **Kepemilikan Tabel Tunggal (*Single Table Ownership*)**:
   - Hanya modul pemilik yang boleh mengeksekusi query `INSERT`, `UPDATE`, `DELETE` pada tabel miliknya.
   - Modul lain tidak boleh melakukan direct SQL `JOIN` ke tabel modul lain.
2. **Denormalisasi Transaksional**:
   - Pada saat pembuatan transaksi (contoh: *Sales Order*), simpan snapshot data penting (misal: `customer_name`, `item_name`, `unit_price`) di tabel transaksi penjualan. Hal ini mencegah ketergantungan join dan menjaga keaslian riwayat transaksi masa lalu saat data master berubah.
3. **Database Transaction Boundaries**:
   - Transaksi database (`pgx.Tx`) hanya boleh melingkupi satu modul.
   - Transaksi lintas modul diselesaikan melalui saga / event-driven kompensasi jika diperlukan, bukan dengan memanjangkan satu database transaction melintasi batas domain.

---

## 4. Multi-Tenancy & Konteks Keamanan

1. Setiap entitas master dan transaksional wajib memuat:
   - `company_id UUID`: Identitas perusahaan/anak usaha.
   - `branch_id UUID`: Identitas cabang operasional.
2. Setiap query baca dan tulis wajib menyertakan filter `WHERE company_id = $1`.
3. Informasi aktor (siapa yang mengeksekusi) diekstrak dari `context.Context` melalui middleware platform:
   - `audit.GetActor(ctx) -> { UserID, CompanyID, BranchID, IP, UserAgent }`
