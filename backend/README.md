# 🏛️ Enterprise ERP Monolith — Backend Service

Backend service sistem ERP Monolith berbasis **Golang (Go 1.22+)** yang mengadopsi arsitektur **Hexagonal Modular Monolith (Ports & Adapters)** dengan HTTP router **Echo v5**, basis data **PostgreSQL 16+**, connection pooler proxy **PgBouncer**, dan tool migrasi **golang-migrate**.

---

## 🗂️ Struktur Direktori

```text
backend/
├── cmd/
│   └── server/
│       └── main.go                 # Entry Point: Dependency injection, wiring module, start Echo v5
│
├── internal/
│   ├── platform/                   # Infrastruktur Teknis (Non-Bisnis)
│   │   ├── database/               # pgxpool connection ke PgBouncer (:6432)
│   │   ├── eventbus/               # In-Memory EventBus untuk domain events asinkron
│   │   ├── middleware/             # TenantScope (Company/Branch), Auth, Logger, CORS
│   │   └── response/               # Standar response JSON & pagination
│   │
│   ├── shared/                     # Shared Kernel
│   │   ├── money/                  # Tipe decimal arbitrary precision (shopspring/decimal)
│   │   ├── types/                  # UUID CompanyID, BranchID, UserID
│   │   ├── audit/                  # Context Actor (UserID, IP, UserAgent)
│   │   └── apperrors/              # Sentinel app errors
│   │
│   └── modules/                    # BOUNDED CONTEXTS (Modul Bisnis Heksagonal)
│       ├── acc/                    # Modul Akuntansi & Keuangan (1_ACC, 4_GL, 2_AP, 3_AR)
│       ├── inv/                    # Modul Persediaan & Gudang (9_INV)
│       ├── sal/                    # Modul Penjualan & Kasir (10_SAL, 12_POS, 11_CRM)
│       ├── pur/                    # Modul Pembelian (8_PUR)
│       ├── hrm/                    # Modul SDM & Payroll (15_HRM, 16_PAY, 17_ATT, 18_REC)
│       └── system/                 # Modul Pendukung (19_USR, 21_ADM, 22_WFL, 25_AUD)
│
├── migrations/                     # Berkas DDL terurut golang-migrate (*.up.sql & *.down.sql)
├── sqlc.yaml                       # Konfigurasi sqlc untuk type-safe SQL query generator
├── .env.example                    # Template variabel lingkungan
└── go.mod                          # Go module dependencies
```

---

## 🚀 Menjalankan Server & Migrasi

### 1. Konfigurasi Lingkungan
Salin berkas `.env.example` ke `.env`:
```bash
cp .env.example .env
```

### 2. Menjalankan Migrasi Skema (golang-migrate)
Pastikan diarahkan langsung ke port PostgreSQL direct (`:5432`):
```bash
# Jalankan migrasi naik (Up)
migrate -path migrations -database "postgres://erp_user:erp_secret@localhost:5432/erp_db?sslmode=disable" up

# Rollback migrasi terakhir (Down)
migrate -path migrations -database "postgres://erp_user:erp_secret@localhost:5432/erp_db?sslmode=disable" down 1
```

### 3. Menjalankan Backend Server
Server akan terhubung ke port PgBouncer (`:6432`):
```bash
go run cmd/server/main.go
```
Endpoint health check dapat diuji di `GET http://localhost:8080/health`.
