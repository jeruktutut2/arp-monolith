# 📜 Enterprise ERP Monolith — Panduan & Aturan Agent

Dokumen ini adalah aturan wajib (*project rules*) untuk seluruh AI Agent saat beroperasi di repositori ini.

---

## 🏛️ 1. Arsitektur Inti: Modular Monolith
- **Bahasa & Framework**: Go 1.22+ di backend (Clean Architecture per modul) dan SvelteKit 2 + Svelte 5 di frontend.
- **Batas Domain**: Setiap modul bisnis terisolasi di dalam `internal/modules/<nama_modul>/`.
- **Dilarang Keras Circular Dependency**:
  - Dilarang meng-import kode konkrit modul lain secara langsung.
  - Untuk dependensi sinkron, gunakan **Consumer-Defined Interface** pada `domain/contracts.go`.
  - Untuk integrasi asinkron dan efek samping (seperti posting jurnal otomatis atau notifikasi), gunakan `EventBus`.

---

## 💰 2. Presisi Finansial: Wajib `decimal.Decimal`
- **Haram menggunakan `float32` atau `float64`** untuk segala nilai moneter, harga pokok, diskon, tarif pajak, dan kuantitas persediaan.
- Selalu gunakan library `github.com/shopspring/decimal`.
- Pada PostgreSQL, gunakan tipe kolom `NUMERIC(18, 4)` atau `NUMERIC(15, 2)`.

---

## 🗄️ 3. Kepemilikan Database & Isolasi Tabel
- Setiap tabel harus memiliki prefiks domain modul (contoh: `acc_`, `inv_`, `sal_`, `pur_`, `sys_`).
- Dilarang membuat query SQL `JOIN` langsung ke tabel milik modul lain. Gunakan snapshot denormalisasi saat transaksi dibuat atau ambil lewat interface modul pemilik data.
- Setiap tabel wajib menyertakan kolom audit multi-tenant: `company_id`, `branch_id`, `created_at`, `created_by`, `updated_at`, `updated_by`.

---

## 🎨 4. Frontend & Library UI Pihak Ketiga
- Gunakan Svelte 5 Runes (`$state`, `$derived`, `$props`, `$effect`) dan Tailwind CSS.
- Ikuti standar library khusus:
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
