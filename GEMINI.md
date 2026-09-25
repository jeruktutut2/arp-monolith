# 📜 Enterprise ERP Monolith — Panduan & Aturan Agent

Dokumen ini adalah aturan wajib (*project rules*) untuk seluruh AI Agent saat beroperasi di repositori ini.

---

## 🏛️ 1. Arsitektur Inti: Hexagonal Modular Monolith
- **Bahasa & Backend Framework**: **Golang (Go 1.22+)** dengan **Echo v5** (`github.com/labstack/echo/v5`) sebagai routing & HTTP delivery adapter.
- **Pola Arsitektur**: **Hexagonal Architecture (Ports & Adapters)** di setiap modul bisnis dalam satu kesatuan Modular Monolith.
  - **Inbound Ports**: Interface Usecase/Application service.
  - **Driving Adapters**: HTTP Handlers dengan Echo v5 (`delivery/http`).
  - **Core Domain**: Model entitas bisnis, kalkulasi finansial (`decimal.Decimal`), dan aturan validasi invarian (`domain/`).
  - **Outbound Ports**: Interface Repository, EventBus, dan Consumer-Defined Contracts (`domain/contracts.go`).
  - **Driven Adapters**: Akses PostgreSQL via `pgx/v5` & `sqlc` melalui **PgBouncer** (`repository/`).
- **Dilarang Keras Circular Dependency**:
  - Dilarang meng-import kode konkrit modul lain secara langsung.
  - Untuk dependensi sinkron, gunakan **Consumer-Defined Interface (Outbound Port)** pada `domain/contracts.go`.
  - Untuk integrasi asinkron dan efek samping (seperti posting jurnal otomatis atau notifikasi), gunakan `EventBus`.

---

## 🗄️ 2. Database, PgBouncer & Migrasi: PostgreSQL + golang-migrate
- **Database Engine**: **PostgreSQL 16+**.
- **Connection Pooler Proxy**: **PgBouncer** dengan mode `pool_mode = transaction`.
  - Seluruh koneksi dari aplikasi Go (`pgxpool`) diarahkan ke port PgBouncer (:6432).
  - Gunakan mode eksekusi query sederhana (`simple protocol` / tanpa prepared statement bentrok) agar kompatibel dengan pooling transaksi PgBouncer.
- **Migration Tool Resmi**: **golang-migrate** (`github.com/golang-migrate/migrate/v4`).
  - Berkas migrasi disimpan berpasangan di folder `migrations/` dengan format: `<seq>_<nama_modul>_<deskripsi>.up.sql` dan `<seq>_<nama_modul>_<deskripsi>.down.sql`.
  - DDL migrasi skema dijalankan langsung ke port direct PostgreSQL (:5432).
- **Kepemilikan Tabel**: Setiap tabel memiliki prefiks domain modul (contoh: `acc_`, `inv_`, `sal_`, `pur_`, `sys_`).
- **Dilarang Direct SQL JOIN Lintas Modul**: Gunakan data snapshot / denormalisasi saat transaksi dibuat atau panggil via interface port modul pemilik data.
- **Kolom Audit Multi-Tenant Wajib**: `company_id`, `branch_id`, `created_at`, `created_by`, `updated_at`, `updated_by`.

---

## 💰 3. Presisi Finansial: Wajib `decimal.Decimal`
- **Haram menggunakan `float32` atau `float64`** untuk segala nilai moneter, harga pokok, diskon, tarif pajak, dan kuantitas persediaan.
- Selalu gunakan library `github.com/shopspring/decimal`.
- Pada PostgreSQL, gunakan tipe kolom `NUMERIC(18, 4)` atau `NUMERIC(15, 2)`.

---

## 🎨 4. Frontend: SvelteKit 2 + Svelte 5 Runes
- **Framework Utama**: **SvelteKit 2** dengan **Svelte 5 Runes** (`$state`, `$derived`, `$props`, `$effect`) dan Tailwind CSS.
- **Pustaka Khusus Terintegrasi**:
  - `14_PRJ`: Frappe Gantt
  - `22_WFL`: XYFlow (`@xyflow/svelte`)
  - `12_POS`: `hotkeys-js` + `on-scan.js` + WebUSB/ESC-POS
  - `1_ACC`, `4_GL`, `9_INV`: AG Grid Community / TanStack Table (DOM Virtualization)
  - `20_RPT`, `17_ATT`: Apache ECharts

---

## 📚 5. Referensi Spesifikasi & Skill
- Master Specification: [SPEC.md](file:///opt/dev/erp_monolith/SPEC.md)
- Development Skill: [.agents/skills/erp-monolith/SKILL.md](file:///opt/dev/erp_monolith/.agents/skills/erp-monolith/SKILL.md)
- Roadmap Implementasi: [design/roadmap_implementasi.md](file:///opt/dev/erp_monolith/design/roadmap_implementasi.md)
- Desain Arsitektur: [design/erp_backend_architecture.md](file:///opt/dev/erp_monolith/design/erp_backend_architecture.md)
