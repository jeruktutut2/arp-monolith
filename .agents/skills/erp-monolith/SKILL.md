---
name: erp-monolith
description: >-
  Panduan dan prosedur standar pengembangan sistem ERP Monolith (Go Modular Monolith + SvelteKit 2).
  Gunakan skill ini setiap kali ada permintaan untuk membuat atau memodifikasi modul bisnis (ACC, GL,
  INV, SAL, PUR, HRM, POS, dll.), membuat migrasi SQL, menulis domain logic dan contracts, mengonfigurasi
  event bus, mengintegrasikan komponen UI dan library khusus (Gantt, Workflow, AG Grid, POS, ECharts),
  atau menjalankan fase-fase roadmap implementasi.
---

# 🏛️ ERP Monolith Development Skill

Skill ini memandu AI agent dalam merancang, menulis, menguji, dan memelihara kode pada proyek **ERP Monolith**. Proyek ini menerapkan arsitektur **Modular Monolith** dengan batas domain (*Bounded Context*) yang tegas di backend Golang dan antarmuka berbasis SvelteKit 2 / Svelte 5.

---

## 🧭 Prinsip Utama (Non-Negotiable Rules)

1. **Modular Monolith & Anti-Circular Dependency**:
   - Seluruh modul bisnis berada di `internal/modules/<nama_modul>/`.
   - Modul dilarang keras meng-import package usecase atau internal modul lain secara langsung.
   - Gunakan **Consumer-Defined Interface** pada `domain/contracts.go` untuk interaksi sinkron.
   - Gunakan **Event Bus** untuk interaksi asinkron / efek samping (contoh: posting jurnal dari invoice).
   - Baca panduan lengkap di: [Aturan Arsitektur](./references/architecture-rules.md).

2. **Presisi Finansial Wajib Desimal**:
   - Dilarang keras menggunakan `float32` atau `float64` untuk nominal uang, diskon, persentase pajak, dan stok.
   - Wajib gunakan `github.com/shopspring/decimal`.

3. **Kepemilikan Tabel & Skema PostgreSQL**:
   - Setiap tabel wajib menggunakan prefiks modul (`acc_`, `inv_`, `sal_`, `pur_`, `sys_`, dll.).
   - Tidak boleh ada query SQL JOIN lintas domain modul privat. Gunakan snapshot denormalisasi saat transaksi dibuat.

4. **Tech Stack Resmi**:
   - Backend: Go 1.22+, Chi Router, PostgreSQL 16+, `pgx/v5`, `sqlc`, `shopspring/decimal`.
   - Frontend: SvelteKit 2, Svelte 5 (Runes), Tailwind CSS, Vite.
   - Daftar pustaka UI lengkap: [Referensi Tech Stack](./references/tech-stack.md).

---

## 🛠️ Alur Kerja Implementasi Fitur/Modul

Saat diminta membuat fitur baru atau mengimplementasikan modul dari desain (`design/ui/`), ikuti langkah terstruktur berikut:

```mermaid
flowchart TD
    A["1. Analisis Spesifikasi UI & Domain<br>(design/ui/<modul> & SPEC.md)"] --> B["2. Buat DDL Migrasi SQL<br>(migrations/xxx_<modul>.sql)"]
    B --> C["3. Implementasikan Domain Layer<br>(Entities, Enums, Contracts, Events)"]
    C --> D["4. Query SQL & Repository<br>(sqlc / pgxpool)"]
    D --> E["5. Business Logic Usecase<br>(Decimal math, validation, events)"]
    E --> F["6. HTTP Delivery Handlers<br>(Chi router, DTO, RBAC check)"]
    F --> G["7. Wiring di cmd/server/main.go"]
    G --> H["8. SvelteKit UI Component & Testing"]
```

Gunakan checklist lengkap di [Checklist Implementasi Modul](./references/module-checklist.md) untuk verifikasi setiap langkah.

---

## 📂 Struktur File Standar per Modul

Setiap modul di bawah `internal/modules/<nama_modul>/` wajib memiliki struktur:

```text
internal/modules/<nama_modul>/
├── domain/
│   ├── entity.go             # Struct domain utama dan validation rules
│   ├── errors.go             # Sentinel errors domain (ErrNotFound, ErrInvalidState)
│   ├── events.go             # Struct payload event yang dipublikasikan
│   └── contracts.go          # Interface dependency ke modul lain (Consumer-Defined)
├── usecase/
│   ├── usecase.go            # Orchestrator usecase struct & interface
│   └── usecase_test.go       # Unit test dengan mock contracts
├── repository/
│   ├── repository.go         # SQL execution via pgx/sqlc
│   └── queries.sql           # Query SQL mentah untuk sqlc
└── delivery/
    └── http/
        ├── handler.go        # HTTP request parsing, DTO binding
        ├── routes.go         # Registrasi Chi sub-router
        └── dto/              # Request / Response JSON structs
```

---

## 🔍 Perintah Verifikasi & Testing

Setelah melakukan perubahan kode, selalu jalankan langkah verifikasi:

```bash
# 1. Format kode Go
go fmt ./...

# 2. Cek linter & circular dependencies
go vet ./...

# 3. Jalankan pengujian unit
go test -v -race ./internal/...

# 4. Generate query sqlc (jika ada perubahan file SQL)
sqlc generate

# 5. Type-check & lint Frontend (jika di folder frontend)
npm run check
npm run test
```

---

## 📚 Dokumen Terkait

- [Master System Specification (SPEC.md)](file:///opt/dev/erp_monolith/SPEC.md)
- [Modul & Fitur ERP (erp_modules.md)](file:///opt/dev/erp_monolith/design/erp_modules.md)
- [Arsitektur Backend Go (erp_backend_architecture.md)](file:///opt/dev/erp_monolith/design/erp_backend_architecture.md)
- [Panduan Library UI (erp_ui_thirdparty_libraries.md)](file:///opt/dev/erp_monolith/design/erp_ui_thirdparty_libraries.md)
- [Roadmap & Urutan Implementasi (roadmap_implementasi.md)](file:///opt/dev/erp_monolith/design/roadmap_implementasi.md)
