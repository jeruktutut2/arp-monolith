# 🧰 Tech Stack & Library Reference

Referensi teknologi yang digunakan dalam proyek ERP Monolith.

---

## 🖥️ Backend (Golang & Hexagonal Architecture)

| Kategori | Paket / Library | Keterangan & Penggunaan |
|:---|:---|:---|
| **Language** | **Go 1.22+ (Golang)** | Strongly typed, compiled, concurrency model tinggi, hemat resource. |
| **Architecture Pattern** | **Hexagonal Architecture** | Ports & Adapters per modul domain dalam satu kesatuan Modular Monolith. |
| **HTTP Framework / Router** | **Echo v5** (`github.com/labstack/echo/v5`) | Framework HTTP berkecepatan tinggi, middleware extensible, sub-router grouping bersih untuk Driving Adapters. |
| **Database Engine** | **PostgreSQL 16+** | Relational database ACID compliance penuh, dukungan native JSONB untuk workflow grafis, partitioning, dan indexing canggih. |
| **Connection Pooler Proxy** | **PgBouncer** | Connection pooler di depan PostgreSQL dengan mode `pool_mode = transaction` untuk menangani ribuan koneksi konkuren aplikasi tanpa membebani RAM database. |
| **Database Driver** | `github.com/jackc/pgx/v5` (`pgxpool`) | Driver PostgreSQL berperforma tinggi di Go. Dikonfigurasi dengan simple protocol mode agar sinkron dengan PgBouncer transaction pooling. |
| **Query Generator** | `github.com/sqlc-dev/sqlc` | Mengompilasi query SQL mentah menjadi kode Go type-safe tanpa overhead runtime reflection ORM. |
| **Database Migration**| `pressly/goose` / `golang-migrate` | Pengelolaan berkas DDL SQL bertahap (dijalankan langsung ke direct port Postgres). |
| **Precision Math** | `github.com/shopspring/decimal` | Penanganan seluruh nilai mata uang, kuantitas stok, dan tarif pajak tanpa floating-point rounding error. |
| **Event Bus** | `ThreeDotsLabs/watermill` / In-Memory Channel | Pengiriman pesan asinkron antar-modul (*decoupled domain events*). |
| **Job Queue** | `hibiken/asynq` / `riverqueue/river` | Background workers (email blast, sinkronisasi bulk, rekonsiliasi akhir bulan). |
| **Structured Log** | `log/slog` (Standard) atau `go.uber.org/zap` | Logging JSON dengan contextual trace_id dan tenant_id. |
| **RBAC / Auth** | Custom RBAC + `golang-jwt/jwt/v5` | Token JWT bertanda tangan dengan verifikasi peran dan branch. |

---

## 🎨 Frontend (SvelteKit 2 + Svelte 5)

| Fitur / Kebutuhan | Library Rekomendasi | Lisensi | Alasan & Modul Terkait |
|:---|:---|:---:|:---|
| **Framework** | **SvelteKit 2 + Svelte 5** | MIT | Reactivity berbasis Runes (`$state`, `$derived`, `$props`), performa tinggi, bundle size kecil. |
| **Styling** | **Tailwind CSS** | MIT | Utility-first CSS dengan dukungan penuh dark/light mode class. |
| **Interactive Gantt** | **Frappe Gantt** | MIT | Zero-dependency, ringan (~15KB), bawaan ERPNext untuk visualisasi timeline proyek (`14_PRJ`). |
| **Workflow Builder** | **XYFlow (`@xyflow/svelte`)** | MIT | Kanvas interaktif Node & Edge untuk visual approval flowchart (`22_WFL`). |
| **DataGrid Raksasa** | **AG Grid Community / TanStack**| MIT | DOM Virtualization untuk render 100.000+ baris buku besar dan kartu stok (`1_ACC`, `4_GL`, `9_INV`). |
| **POS Keyboard Bind**| **hotkeys-js** | MIT | Shortcut tombol keyboard (`F1`-`F12`, `Esc`) untuk operasional kasir cepat (`12_POS`). |
| **Barcode Ingestion**| **on-scan.js** | MIT | Intersepsi scanner barcode fisik agar tidak bentrok dengan active input focus (`12_POS`). |
| **Thermal Printing** | **esc-pos-encoder** + WebUSB | MIT | Cetak langsung struk belanja ke printer thermal POS tanpa dialog browser (`12_POS`). |
| **Executive Charts** | **Apache ECharts** | Apache 2.0 | Grafik multi-axis finansial, zoom slider, dan kalender heatmap kehadiran (`20_RPT`, `17_ATT`). |
