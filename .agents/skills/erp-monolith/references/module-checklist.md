# 📋 Checklist Implementasi Modul ERP

Gunakan panduan ini saat membuat modul baru atau menambahkan fitur pada modul yang sudah ada.

---

## 🛠️ Langkah-Langkah Pengerjaan Modul Baru (Hexagonal Pattern)

### Tahap 1: Desain Basis Data & Migrasi PostgreSQL (golang-migrate)
- [ ] Buat berkas migrasi baru berpasangan di folder `migrations/` dengan penomoran berurutan:
  - Contoh: `migrations/000010_sal_sales_orders.up.sql` dan `migrations/000010_sal_sales_orders.down.sql`
- [ ] Tentukan nama tabel dengan prefiks modul (contoh: `sal_orders`, `sal_order_items`).
- [ ] Tambahkan kolom audit wajib:
  - `company_id UUID NOT NULL`
  - `branch_id UUID NOT NULL`
  - `created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP`
  - `created_by UUID NOT NULL`
  - `updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP`
  - `updated_by UUID NOT NULL`
- [ ] Pastikan kolom uang/angka menggunakan `NUMERIC(18, 4)` atau `NUMERIC(15, 2)`.
- [ ] Tambahkan index pada `company_id`, foreign keys, dan status transaksi.

### Tahap 2: Definisi Core Domain & Ports
- [ ] Buat folder `internal/modules/<nama_modul>/domain/`.
- [ ] Definisikan entitas struct utama (contoh: `SalesOrder`, `SalesOrderItem`).
- [ ] Definisikan Value Objects dan konstanta status (contoh: `OrderStatusDraft`, `OrderStatusApproved`).
- [ ] Definisikan sentinel errors (contoh: `ErrOrderNotFound`, `ErrInsufficientStock`).
- [ ] Definisikan **Inbound Ports** pada `ports.go` (interface yang diekspos usecase ke HTTP/CLI).
- [ ] Definisikan **Outbound Ports** pada `contracts.go` (interface repository dan consumer contract ke modul lain).
- [ ] Definisikan Domain Events (contoh: `OrderCreatedEvent`, `OrderApprovedEvent`).

### Tahap 3: Driven Adapter: Repository PostgreSQL (sqlc + PgBouncer)
- [ ] Tulis query SQL di file `sqlc/queries/<nama_modul>.sql`.
- [ ] Jalankan generator `sqlc generate`.
- [ ] Buat implementasi repository di `internal/modules/<nama_modul>/repository/` yang membungkus output sqlc dan mengembalikan domain entities.
- [ ] Pastikan query kompatibel dengan mode `pool_mode = transaction` PgBouncer (hindari prepared statement session-level).

### Tahap 4: Inbound Port Implementation: Usecase (Business Logic)
- [ ] Buat file di `internal/modules/<nama_modul>/usecase/`.
- [ ] Tulis logika validasi status, kalkulasi matematis (menggunakan `decimal.Decimal`), dan orkestrasi transaksi.
- [ ] Terbitkan event ke `eventBus` jika ada side-effects (contoh: posting jurnal atau notifikasi).

### Tahap 5: Driving Adapter: Echo v5 HTTP Delivery
- [ ] Buat DTO Request dan Response di `internal/modules/<nama_modul>/delivery/http/dto/`.
- [ ] Implementasikan handler **Echo v5** (`github.com/labstack/echo/v5`) di `internal/modules/<nama_modul>/delivery/http/handler.go`.
- [ ] Registrasikan sub-router di `internal/modules/<nama_modul>/delivery/http/routes.go` menggunakan `echo.Group`.
- [ ] Pasang middleware otorisasi RBAC (contoh: `RequirePermission("sales:create")`).

### Tahap 6: Wiring di Application Entry Point
- [ ] Buka `cmd/server/main.go`.
- [ ] Inisialisasi pool PostgreSQL/PgBouncer via `pgxpool`.
- [ ] Inisialisasi repository modul (Driven Adapter).
- [ ] Hubungkan dependensi usecase (masukkan service modul lain yang memenuhi outbound consumer contract).
- [ ] Daftarkan event subscriber modul pada `eventBus`.
- [ ] Pasang routes modul ke root router Echo v5 (`/api/v1/<modul>`).

### Tahap 7: Frontend UI (SvelteKit 2 + Svelte 5 Runes)
- [ ] Buka atau buat rute di `src/routes/(app)/<modul>/+page.svelte`.
- [ ] Terapkan layout standar App Shell (`0_HDR`, `0_SDB`, `0_FTR`).
- [ ] Terapkan reaktivitas berbasis Runes (`$state`, `$derived`, `$props`).
- [ ] Hubungkan library pihak ketiga jika diperlukan:
  - Tabular data masif: **AG Grid Community / TanStack Table**
  - Timeline WBS: **Frappe Gantt** (`14_PRJ`)
  - Workflow builder: **XYFlow** (`22_WFL`)
  - Kasir cepat: **hotkeys-js** & **on-scan.js** (`12_POS`)
  - Dashboard analitik: **Apache ECharts** (`20_RPT`)
- [ ] Pastikan mendukung responsivitas dan toggle Light / Dark mode.
