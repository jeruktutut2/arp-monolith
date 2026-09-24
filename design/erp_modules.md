# 📦 ERP System — Modul & Deskripsi

> Dokumen ini menjelaskan seluruh modul yang terdapat dalam sistem ERP beserta deskripsi, fitur utama, dan keterhubungan antar modul.

---

## Daftar Modul

| No | Kode Modul | Nama Modul | Kategori |
|----|-----------|------------|----------|
| 1 | `ACC` | Akuntansi & Keuangan | Keuangan |
| 2 | `AP` | Hutang Usaha (Accounts Payable) | Keuangan |
| 3 | `AR` | Piutang Usaha (Accounts Receivable) | Keuangan |
| 4 | `GL` | Buku Besar (General Ledger) | Keuangan |
| 5 | `FA` | Aset Tetap (Fixed Assets) | Keuangan |
| 6 | `BUD` | Anggaran (Budgeting) | Keuangan |
| 7 | `TAX` | Perpajakan | Keuangan |
| 8 | `PUR` | Pembelian (Purchasing) | Rantai Pasok |
| 9 | `INV` | Persediaan & Gudang (Inventory) | Rantai Pasok |
| 10 | `SAL` | Penjualan (Sales) | Penjualan & CRM |
| 11 | `CRM` | Customer Relationship Management | Penjualan & CRM |
| 12 | `POS` | Point of Sale | Penjualan & CRM |
| 13 | `MFG` | Manufaktur / Produksi | Operasional |
| 14 | `PRJ` | Manajemen Proyek | Operasional |
| 15 | `HRM` | Sumber Daya Manusia (HRM) | SDM |
| 16 | `PAY` | Penggajian (Payroll) | SDM |
| 17 | `ATT` | Absensi & Kehadiran | SDM |
| 18 | `REC` | Rekrutmen | SDM |
| 19 | `USR` | User Management | Sistem |
| 20 | `RPT` | Laporan & Analitik | Sistem |
| 21 | `ADM` | Administrasi & Pengaturan | Sistem |
| 22 | `WFL` | Workflow & Approval | Sistem |
| 23 | `DOC` | Manajemen Dokumen | Sistem |
| 24 | `MSG` | Pesan & Notifikasi | Sistem |
| 25 | `AUD` | Audit Trail | Sistem |

---

## Detail Modul

---

### 1. `ACC` — Akuntansi & Keuangan

**Kategori:** Keuangan
**Deskripsi:**
Modul inti yang mengelola seluruh proses akuntansi dan keuangan perusahaan. Mencakup pencatatan transaksi keuangan, pembuatan laporan keuangan, dan pengelolaan chart of accounts.

**Fitur Utama:**
- Chart of Accounts (CoA) — daftar akun dengan hierarki multi-level
- Jurnal Umum (Journal Entry) — pencatatan transaksi debit/kredit
- Laporan Keuangan — Neraca, Laba Rugi, Arus Kas, Perubahan Ekuitas
- Multi-mata uang (multi-currency) dengan kurs harian
- Multi-cabang / multi-company consolidation
- Periode akuntansi (fiscal year) — buka/tutup periode
- Rekonsiliasi bank
- Laporan aging (hutang & piutang)

**Relasi:**
- Menerima posting dari: `AP`, `AR`, `FA`, `PAY`, `SAL`, `PUR`
- Menyediakan data ke: `RPT`, `BUD`, `TAX`

---

### 2. `AP` — Hutang Usaha (Accounts Payable)

**Kategori:** Keuangan
**Deskripsi:**
Mengelola seluruh kewajiban perusahaan kepada vendor/supplier, mulai dari penerimaan invoice, persetujuan pembayaran, hingga pelunasan hutang.

**Fitur Utama:**
- Registrasi & verifikasi invoice supplier
- Three-way matching (PO — Goods Receipt — Invoice)
- Penjadwalan pembayaran (payment schedule)
- Pembayaran batch (bulk payment)
- Debit note & retur pembelian
- Laporan hutang & aging payable
- Manajemen data vendor/supplier

**Relasi:**
- Input dari: `PUR` (Purchase Order, Goods Receipt)
- Output ke: `ACC` (posting jurnal), `TAX` (PPN masukan)

---

### 3. `AR` — Piutang Usaha (Accounts Receivable)

**Kategori:** Keuangan
**Deskripsi:**
Mengelola seluruh piutang perusahaan dari pelanggan, mulai dari penerbitan invoice penjualan, pelacakan pembayaran, hingga penagihan.

**Fitur Utama:**
- Pembuatan & pengiriman invoice penjualan
- Pencatatan penerimaan pembayaran (receipt)
- Credit note & retur penjualan
- Laporan piutang & aging receivable
- Manajemen kredit pelanggan (credit limit)
- Reminder & surat penagihan otomatis
- Diskon & potongan penjualan

**Relasi:**
- Input dari: `SAL` (Sales Order, Delivery)
- Output ke: `ACC` (posting jurnal), `TAX` (PPN keluaran)

---

### 4. `GL` — Buku Besar (General Ledger)

**Kategori:** Keuangan
**Deskripsi:**
Sub-modul akuntansi yang mengelola buku besar sebagai pusat pencatatan seluruh transaksi keuangan. Semua modul keuangan lainnya melakukan posting ke General Ledger.

**Fitur Utama:**
- Buku besar per akun (ledger account)
- Trial balance
- Jurnal penyesuaian (adjusting entry)
- Jurnal penutup (closing entry)
- Intercompany transaction & eliminasi
- Dimensi akuntansi (department, project, cost center)
- Laporan buku besar detail & ringkasan

**Relasi:**
- Menerima posting dari seluruh modul keuangan
- Menyediakan data ke: `RPT`, `BUD`

---

### 5. `FA` — Aset Tetap (Fixed Assets)

**Kategori:** Keuangan
**Deskripsi:**
Mengelola siklus hidup aset tetap perusahaan mulai dari perolehan, penyusutan, pemeliharaan, hingga penghapusan/disposal aset.

**Fitur Utama:**
- Registrasi & kategori aset tetap
- Metode penyusutan (straight-line, declining balance, double declining, dll.)
- Perhitungan penyusutan otomatis (bulanan/tahunan)
- Revaluasi aset
- Transfer & mutasi aset antar lokasi/departemen
- Penghapusan & penjualan aset (disposal)
- Label & barcode aset untuk stock opname
- Laporan daftar aset, penyusutan, & nilai buku

**Relasi:**
- Input dari: `PUR` (pembelian aset)
- Output ke: `ACC` (posting penyusutan & disposal)

---

### 6. `BUD` — Anggaran (Budgeting)

**Kategori:** Keuangan
**Deskripsi:**
Mengelola perencanaan, penetapan, dan pengendalian anggaran perusahaan. Memungkinkan perbandingan antara anggaran (budget) dengan realisasi (actual).

**Fitur Utama:**
- Penyusunan anggaran per departemen/divisi/proyek
- Template anggaran (annual, quarterly, monthly)
- Revisi & versioning anggaran
- Budget vs Actual report
- Peringatan jika pengeluaran mendekati/melebihi anggaran
- Approval workflow anggaran
- Forecasting & proyeksi keuangan

**Relasi:**
- Input dari: `ACC`, `GL` (data realisasi)
- Output ke: `RPT` (laporan budget vs actual)
- Terintegrasi dengan: `WFL` (approval)

---

### 7. `TAX` — Perpajakan

**Kategori:** Keuangan
**Deskripsi:**
Mengelola perhitungan, pelaporan, dan kepatuhan pajak perusahaan sesuai regulasi yang berlaku (PPN, PPh, dll.).

**Fitur Utama:**
- Konfigurasi tarif pajak (PPN, PPh 21, PPh 23, PPh 4(2), dll.)
- Perhitungan PPN masukan & keluaran
- Faktur pajak (e-Faktur integration ready)
- Bukti potong PPh
- SPT Masa & SPT Tahunan
- Rekonsiliasi pajak
- Laporan pajak per periode

**Relasi:**
- Input dari: `AP` (PPN masukan), `AR` (PPN keluaran), `PAY` (PPh 21)
- Output ke: `RPT` (laporan perpajakan)

---

### 8. `PUR` — Pembelian (Purchasing)

**Kategori:** Rantai Pasok
**Deskripsi:**
Mengelola seluruh proses pengadaan barang dan jasa, mulai dari permintaan pembelian, pemilihan vendor, pembuatan Purchase Order, hingga penerimaan barang.

**Fitur Utama:**
- Purchase Requisition (PR) — permintaan pembelian internal
- Request for Quotation (RFQ) — permintaan penawaran ke vendor
- Vendor comparison & selection
- Purchase Order (PO) — pembuatan & pengiriman PO
- Goods Receipt (GR) — penerimaan barang di gudang
- Retur pembelian (purchase return)
- Kontrak & perjanjian pembelian (blanket order)
- Evaluasi kinerja vendor (vendor rating)
- Approval workflow bertingkat

**Relasi:**
- Output ke: `INV` (stok masuk), `AP` (invoice), `FA` (aset tetap baru)
- Terintegrasi dengan: `WFL` (approval), `DOC` (dokumen PO/kontrak)

---

### 9. `INV` — Persediaan & Gudang (Inventory)

**Kategori:** Rantai Pasok
**Deskripsi:**
Mengelola seluruh persediaan barang perusahaan di satu atau lebih gudang, termasuk penerimaan, penyimpanan, pengeluaran, dan stock opname.

**Fitur Utama:**
- Master barang/produk (item master) dengan kategori & satuan
- Multi-gudang (multi-warehouse) & multi-lokasi (bin/rack)
- Penerimaan barang (goods receipt)
- Pengeluaran barang (goods issue)
- Transfer antar gudang (stock transfer)
- Metode valuasi: FIFO, LIFO, Average, Specific Identification
- Stock opname / physical inventory
- Minimum stock & reorder point (auto purchase requisition)
- Serial number & batch/lot tracking
- Barcode / QR code integration
- Laporan stok: kartu stok, aging inventory, slow-moving items

**Relasi:**
- Input dari: `PUR` (barang masuk), `MFG` (hasil produksi)
- Output ke: `SAL` (barang keluar/delivery), `MFG` (bahan baku)
- Terintegrasi dengan: `ACC` (valuasi persediaan)

---

### 10. `SAL` — Penjualan (Sales)

**Kategori:** Penjualan & CRM
**Deskripsi:**
Mengelola seluruh proses penjualan mulai dari penawaran harga, pemesanan, pengiriman barang, hingga penagihan ke pelanggan.

**Fitur Utama:**
- Quotation / Penawaran Harga
- Sales Order (SO) — pemesanan penjualan
- Delivery Order (DO) — pengiriman barang
- Sales Invoice — faktur penjualan
- Retur penjualan (sales return)
- Diskon bertingkat (tiered pricing)
- Daftar harga (price list) per pelanggan/segmen
- Komisi sales
- Laporan penjualan: per produk, per pelanggan, per sales person, per periode
- Target penjualan & pencapaian

**Relasi:**
- Input dari: `CRM` (leads/opportunities converted)
- Output ke: `INV` (pengeluaran stok), `AR` (invoice), `ACC` (revenue)
- Terintegrasi dengan: `WFL` (approval diskon/harga khusus)

---

### 11. `CRM` — Customer Relationship Management

**Kategori:** Penjualan & CRM
**Deskripsi:**
Mengelola hubungan dengan pelanggan dan calon pelanggan, mulai dari akuisisi lead, nurturing, hingga konversi menjadi penjualan.

**Fitur Utama:**
- Manajemen lead & prospek
- Pipeline penjualan (sales funnel)
- Aktivitas & follow-up (call, meeting, email)
- Konversi lead → opportunity → quotation → sales order
- Manajemen kontak & akun pelanggan
- Segmentasi pelanggan
- Campaign management
- Customer support & ticketing
- Dashboard CRM & analytics
- Riwayat interaksi pelanggan

**Relasi:**
- Output ke: `SAL` (quotation & sales order)
- Terintegrasi dengan: `MSG` (notifikasi follow-up), `DOC` (dokumen proposal)

---

### 12. `POS` — Point of Sale

**Kategori:** Penjualan & CRM
**Deskripsi:**
Modul kasir/penjualan langsung untuk transaksi retail. Mendukung penjualan cepat dengan antarmuka yang dioptimalkan untuk transaksi di toko/outlet.

**Fitur Utama:**
- Antarmuka kasir touchscreen-friendly
- Scan barcode / QR code
- Multi metode pembayaran (tunai, kartu, e-wallet, transfer)
- Split payment
- Diskon & promosi otomatis
- Cetak struk / receipt (thermal printer)
- Buka/tutup kasir (cash register open/close)
- Laporan penjualan harian per kasir/outlet
- Mode offline (sync saat online)
- Manajemen shift kasir

**Relasi:**
- Output ke: `INV` (pengurangan stok real-time), `AR` (piutang jika kredit), `ACC` (revenue)
- Terintegrasi dengan: `SAL` (data penjualan tergabung)

---

### 13. `MFG` — Manufaktur / Produksi

**Kategori:** Operasional
**Deskripsi:**
Mengelola seluruh proses produksi/manufaktur, mulai dari perencanaan produksi, Bill of Materials, work order, hingga pencatatan hasil produksi.

**Fitur Utama:**
- Bill of Materials (BoM) — daftar bahan baku & komponen
- Multi-level BoM (sub-assembly)
- Production Planning — perencanaan produksi
- Work Order / Manufacturing Order
- Routing & work center
- Material Requirement Planning (MRP)
- Pencatatan hasil produksi (production receipt)
- Quality control / inspection
- Scrap & waste management
- Kapasitas produksi & scheduling
- Biaya produksi (cost of production)

**Relasi:**
- Input dari: `INV` (bahan baku), `SAL` (demand/forecast)
- Output ke: `INV` (barang jadi), `ACC` (biaya produksi)

---

### 14. `PRJ` — Manajemen Proyek

**Kategori:** Operasional
**Deskripsi:**
Mengelola proyek-proyek perusahaan termasuk perencanaan, alokasi sumber daya, pelacakan progress, dan pengendalian biaya proyek.

**Fitur Utama:**
- Pembuatan & konfigurasi proyek
- Work Breakdown Structure (WBS)
- Task management & assignment
- Gantt chart & timeline
- Milestone tracking
- Alokasi sumber daya (manpower, material)
- Timesheet & pencatatan jam kerja
- Budget proyek & cost tracking
- Progress reporting (% completion)
- Issue & risk management
- Laporan profitabilitas proyek

**Relasi:**
- Terintegrasi dengan: `HRM` (alokasi karyawan), `PUR` (pembelian material proyek), `ACC` (biaya proyek), `BUD` (anggaran proyek)

---

### 15. `HRM` — Sumber Daya Manusia (HRM)

**Kategori:** SDM
**Deskripsi:**
Mengelola seluruh data dan proses terkait sumber daya manusia, mulai dari data karyawan, struktur organisasi, hingga pengembangan karir.

**Fitur Utama:**
- Data karyawan (personal, keluarga, pendidikan, pengalaman)
- Struktur organisasi & departemen
- Jabatan & grade/level karyawan
- Kontrak kerja & masa kerja
- Mutasi, promosi, & demosi
- Cuti & izin (leave management)
- Pelatihan & pengembangan (training)
- Penilaian kinerja (performance appraisal)
- Surat peringatan & sanksi
- Resign & offboarding
- Employee self-service portal
- Laporan kepegawaian & dashboard HR

**Relasi:**
- Output ke: `PAY` (data untuk penggajian), `ATT` (data kehadiran)
- Input dari: `REC` (karyawan baru)
- Terintegrasi dengan: `PRJ` (alokasi proyek), `WFL` (approval cuti/mutasi)

---

### 16. `PAY` — Penggajian (Payroll)

**Kategori:** SDM
**Deskripsi:**
Mengelola perhitungan dan pembayaran gaji karyawan termasuk komponen gaji, potongan, tunjangan, dan kepatuhan pajak penghasilan.

**Fitur Utama:**
- Komponen gaji: gaji pokok, tunjangan, lembur, bonus, insentif
- Komponen potongan: PPh 21, BPJS Kesehatan, BPJS Ketenagakerjaan, pinjaman karyawan
- Perhitungan PPh 21 otomatis (TER / tarif efektif)
- Slip gaji (payslip) digital
- Proses payroll bulanan (batch processing)
- Bank transfer file generation
- THR & bonus tahunan
- Laporan payroll & rekapitulasi gaji
- Bukti potong 1721-A1

**Relasi:**
- Input dari: `HRM` (data karyawan), `ATT` (data kehadiran & lembur)
- Output ke: `ACC` (posting biaya gaji), `TAX` (PPh 21)

---

### 17. `ATT` — Absensi & Kehadiran

**Kategori:** SDM
**Deskripsi:**
Mengelola pencatatan kehadiran karyawan, jam kerja, lembur, dan integrasi dengan mesin absensi.

**Fitur Utama:**
- Integrasi mesin absensi (fingerprint, face recognition, kartu)
- Absensi mobile (GPS-based)
- Jadwal kerja & shift management
- Perhitungan jam kerja & lembur otomatis
- Rekap kehadiran harian/bulanan
- Keterlambatan & pulang cepat
- Work From Home (WFH) tracking
- Kalender hari libur & cuti bersama
- Laporan kehadiran per karyawan/departemen

**Relasi:**
- Output ke: `PAY` (data kehadiran untuk penggajian)
- Terintegrasi dengan: `HRM` (data karyawan & jadwal cuti)

---

### 18. `REC` — Rekrutmen

**Kategori:** SDM
**Deskripsi:**
Mengelola proses rekrutmen karyawan dari perencanaan kebutuhan tenaga kerja hingga onboarding karyawan baru.

**Fitur Utama:**
- Manpower planning — perencanaan kebutuhan karyawan
- Job posting — publikasi lowongan
- Applicant tracking — pelacakan pelamar
- Screening & shortlisting
- Penjadwalan interview
- Penilaian kandidat (scoring & evaluation)
- Offering letter generation
- Onboarding checklist
- Talent pool / database kandidat
- Laporan rekrutmen & analytics

**Relasi:**
- Output ke: `HRM` (data karyawan baru setelah onboarding)
- Terintegrasi dengan: `WFL` (approval rekrutmen), `DOC` (dokumen kontrak)

---

### 19. `USR` — User Management

**Kategori:** Sistem
**Deskripsi:**
Modul khusus untuk mengelola pengguna sistem ERP, termasuk autentikasi, otorisasi, profil pengguna, dan keamanan akses.

**Fitur Utama:**
- Registrasi & manajemen akun pengguna
- Autentikasi (login/logout) — email/password, SSO, OAuth2, LDAP/Active Directory
- Multi-Factor Authentication (MFA/2FA) — OTP, authenticator app
- Role-Based Access Control (RBAC) — role, permission, privilege
- Hak akses per modul, menu, dan aksi (create/read/update/delete)
- Profil pengguna (avatar, kontak, preferensi)
- Password policy (kompleksitas, expiry, riwayat password)
- Session management (active sessions, force logout, session timeout)
- Account lockout setelah percobaan login gagal
- Reset & forgot password (via email/SMS)
- User group & organisasi
- API key & token management untuk integrasi
- Login history & security log
- IP whitelist / blacklist
- Impersonation (admin login sebagai user lain untuk troubleshooting)

**Relasi:**
- Menyediakan autentikasi & otorisasi ke: Seluruh modul ERP
- Terintegrasi dengan: `ADM` (konfigurasi sistem), `HRM` (link ke data karyawan), `AUD` (log aktivitas login), `MSG` (notifikasi keamanan)

---

### 20. `RPT` — Laporan & Analitik

**Kategori:** Sistem
**Deskripsi:**
Modul pelaporan terpusat yang menyediakan laporan, dashboard, dan analitik dari seluruh modul ERP.

**Fitur Utama:**
- Dashboard eksekutif (KPI overview)
- Report builder — pembuatan laporan kustom
- Laporan standar per modul
- Grafik & visualisasi data interaktif
- Export laporan (PDF, Excel, CSV)
- Scheduled report — laporan otomatis terjadwal
- Drill-down dari ringkasan ke detail
- Comparative analysis (period-over-period)
- Filter & parameter laporan dinamis

**Relasi:**
- Input dari: Seluruh modul ERP
- Standalone: Tidak menghasilkan output ke modul lain

---

### 21. `ADM` — Administrasi & Pengaturan

**Kategori:** Sistem
**Deskripsi:**
Modul administrasi sistem untuk mengelola pengaturan umum dan konfigurasi ERP. Manajemen pengguna dan hak akses ditangani oleh modul `USR`.

**Fitur Utama:**
- Pengaturan perusahaan (company profile, logo, alamat)
- Konfigurasi multi-cabang / multi-company
- Pengaturan regional (mata uang, format tanggal, bahasa)
- Master data umum (negara, provinsi, kota, kode pos)
- Konfigurasi email (SMTP) & integrasi
- Konfigurasi numbering / penomoran dokumen
- Manajemen lisensi
- Backup & restore
- System health monitoring
- Konfigurasi integrasi pihak ketiga (API keys, webhook)

**Relasi:**
- Menyediakan konfigurasi ke: Seluruh modul ERP
- Terintegrasi dengan: `USR` (konfigurasi autentikasi & keamanan)

---

### 22. `WFL` — Workflow & Approval

**Kategori:** Sistem
**Deskripsi:**
Mengelola alur persetujuan (approval workflow) untuk berbagai proses bisnis yang memerlukan otorisasi bertingkat.

**Fitur Utama:**
- Workflow designer — desain alur approval visual
- Multi-level approval (sequential & parallel)
- Delegation & substitution (penunjukan pengganti)
- Auto-approval berdasarkan threshold/nilai
- Eskalasi otomatis jika melewati batas waktu
- Approval via email / mobile notification
- Riwayat approval (audit trail)
- Template workflow per jenis dokumen

**Relasi:**
- Terintegrasi dengan: `PUR` (approval PO/PR), `SAL` (approval diskon), `HRM` (approval cuti/mutasi), `BUD` (approval anggaran), `REC` (approval rekrutmen)
- Output ke: `MSG` (notifikasi approval)

---

### 23. `DOC` — Manajemen Dokumen

**Kategori:** Sistem
**Deskripsi:**
Mengelola penyimpanan, pengorganisasian, dan pencarian dokumen digital yang terkait dengan seluruh proses bisnis di ERP.

**Fitur Utama:**
- Upload & penyimpanan dokumen (PDF, image, dll.)
- Kategori & tagging dokumen
- Versioning dokumen
- Pencarian dokumen (full-text search)
- Lampiran dokumen ke transaksi (PO, invoice, kontrak, dll.)
- Hak akses dokumen per user/role
- Template dokumen (surat, kontrak, dll.)
- Digital signature integration ready
- Expiry & reminder dokumen (kontrak, lisensi, dll.)

**Relasi:**
- Terintegrasi dengan: Seluruh modul yang menghasilkan/membutuhkan dokumen

---

### 24. `MSG` — Pesan & Notifikasi

**Kategori:** Sistem
**Deskripsi:**
Mengelola pengiriman pesan dan notifikasi kepada pengguna untuk berbagai event di sistem ERP.

**Fitur Utama:**
- Notifikasi in-app (bell notification)
- Email notification
- Push notification (mobile)
- WhatsApp / SMS integration ready
- Konfigurasi template notifikasi
- Preferensi notifikasi per user
- Broadcast message ke grup/departemen
- Notification center — riwayat notifikasi
- Event-driven notification (trigger dari modul lain)

**Relasi:**
- Input dari: `WFL` (approval notification), `CRM` (follow-up reminder), dan modul lainnya
- Standalone: Tidak menghasilkan output ke modul lain

---

### 25. `AUD` — Audit Trail

**Kategori:** Sistem
**Deskripsi:**
Mencatat seluruh aktivitas pengguna dan perubahan data di sistem ERP untuk keperluan audit, keamanan, dan compliance.

**Fitur Utama:**
- Log aktivitas pengguna (login, logout, akses menu)
- Log perubahan data (create, update, delete) — field-level tracking
- Siapa, kapan, dari nilai apa, ke nilai apa (who, when, old value, new value)
- Filter & pencarian log berdasarkan user, modul, tanggal, tipe aktivitas
- Laporan audit per periode
- Data retention policy
- Tamper-proof log storage
- Export audit log

**Relasi:**
- Input dari: Seluruh modul ERP (setiap aksi pengguna tercatat)
- Standalone: Tidak menghasilkan output ke modul lain

---

## Diagram Relasi Antar Modul

```
                     ┌──────────┐          ┌──────────┐
                     │   ADM    │◄────────►│   USR    │
                     │ (Config) │          │  (Auth)  │
                     └────┬─────┘          └────┬─────┘
                          │ config               │ autentikasi & otorisasi
                          ▼                      ▼
┌────────────────────────────────────────────────────────────────────────────┐
│                            SELURUH MODUL                                   │
│                                                                            │
│  ┌──────────────── KEUANGAN ──────────────────┐  ┌── RANTAI PASOK ─────┐  │
│  │                                             │  │   & OPERASIONAL     │  │
│  │  ┌─────┐   ┌─────┐   ┌─────┐              │  │                      │  │
│  │  │ GL  │◄──┤ AP  │◄──┤ PUR │◄─────────────┼──┼───┌─────┐           │  │
│  │  │     │   │     │   │     │─────────────►┼──┼──►│ INV │           │  │
│  │  │     │   └──┬──┘   └──┬──┘──►┌─────┐   │  │   │     │◄──┐       │  │
│  │  │     │      │         │      │ FA  │   │  │   │     │───┤       │  │
│  │  │     │◄──┐  │         │      └──┬──┘   │  │   └──┬──┘ ┌─┴────┐  │  │
│  │  └──┬──┘   │  │         │         │      │  │      │    │ MFG  │  │  │
│  │     │      │  │         │         ▼      │  │      │    └──┬───┘  │  │
│  │  ┌──┴──┐   │  ▼         │      ┌─────┐  │  │      │       │      │  │
│  │  │ ACC │◄──┼──┤ TAX │   │      │ ACC │  │  │      │       │      │  │
│  │  │     │   │  └─────┘   │      └─────┘  │  │      │  ┌────┴──┐   │  │
│  │  └──▲──┘   │            │               │  │      │  │  PRJ  │   │  │
│  │     │      │  ┌─────┐   │    ┌─────┐    │  │      │  └───────┘   │  │
│  │     │      └──┤ AR  │◄──┼────┤ SAL │◄───┼──┼──────┘              │  │
│  │     │         │     │   │    │     │────┼──┼──►(INV, stok keluar)│  │
│  │     │         └─────┘   │    └──┬──┘    │  └──────────────────────┘  │
│  │     │                   │       │       │                             │
│  │  ┌──┴──┐   ┌─────┐   ┌─┴────┐  │       │                             │
│  │  │ BUD │   │ POS │   │ CRM  │──┘       │                             │
│  │  └─────┘   └──┬──┘   └──────┘          │                             │
│  │               │                         │                             │
│  └───────────────┼─────────────────────────┘                             │
│                  │ ──►INV, AR, ACC                                        │
│                                                                            │
│  ┌────────────── SDM ──────────────┐  ┌────────── SISTEM ──────────────┐  │
│  │                                  │  │                                │  │
│  │  ┌─────┐       ┌─────┐         │  │  ┌─────┐  ┌─────┐  ┌─────┐   │  │
│  │  │ REC │──────►│ HRM │         │  │  │ WFL │─►│ MSG │  │ RPT │   │  │
│  │  └─────┘       │     │         │  │  └─────┘  └─────┘  └─────┘   │  │
│  │                │     │──►PAY   │  │                                │  │
│  │  ┌─────┐       │     │         │  │  ┌─────┐  ┌─────┐            │  │
│  │  │ ATT │──────►└──┬──┘         │  │  │ DOC │  │ AUD │            │  │
│  │  └──┬──┘          │            │  │  └─────┘  └─────┘            │  │
│  │     │          ┌──┴──┐         │  │                                │  │
│  │     └─────────►│ PAY │──►ACC   │  └────────────────────────────────┘  │
│  │                │     │──►TAX   │                                       │
│  │                └─────┘         │                                       │
│  └──────────────────────────────── │                                       │
│                                                                            │
└────────────────────────────────────────────────────────────────────────────┘
```

### Matriks Relasi Detail Antar Modul

> Tabel berikut mendokumentasikan **seluruh relasi** yang dideklarasikan pada setiap spesifikasi modul.
> Panah `→` = Output, `←` = Input, `↔` = Terintegrasi dua arah.

| Modul | Input Dari (←) | Output Ke (→) | Terintegrasi Dengan (↔) |
|-------|----------------|---------------|--------------------------|
| **ACC** | AP, AR, FA, PAY, SAL, PUR | RPT, BUD, TAX | — |
| **AP** | PUR (PO, GR) | ACC (jurnal), TAX (PPN masukan) | — |
| **AR** | SAL (SO, Delivery) | ACC (jurnal), TAX (PPN keluaran) | — |
| **GL** | Seluruh modul keuangan | RPT, BUD | — |
| **FA** | PUR (pembelian aset) | ACC (penyusutan & disposal) | — |
| **BUD** | ACC, GL (data realisasi) | RPT (budget vs actual) | WFL (approval) |
| **TAX** | AP (PPN masukan), AR (PPN keluaran), PAY (PPh 21) | RPT (laporan pajak) | — |
| **PUR** | — | INV (stok masuk), AP (invoice), FA (aset baru) | WFL (approval), DOC (PO/kontrak) |
| **INV** | PUR (barang masuk), MFG (hasil produksi) | SAL (barang keluar), MFG (bahan baku) | ACC (valuasi) |
| **SAL** | CRM (leads converted) | INV (pengeluaran stok), AR (invoice), ACC (revenue) | WFL (approval diskon) |
| **CRM** | — | SAL (quotation & SO) | MSG (follow-up), DOC (proposal) |
| **POS** | — | INV (pengurangan stok), AR (piutang kredit), ACC (revenue) | SAL (data tergabung) |
| **MFG** | INV (bahan baku), SAL (demand) | INV (barang jadi), ACC (biaya produksi) | — |
| **PRJ** | — | — | HRM (alokasi SDM), PUR (material), ACC (biaya), BUD (anggaran) |
| **HRM** | REC (karyawan baru) | PAY (data gaji), ATT (data kehadiran) | PRJ (alokasi), WFL (approval cuti) |
| **PAY** | HRM (data karyawan), ATT (kehadiran & lembur) | ACC (biaya gaji), TAX (PPh 21) | — |
| **ATT** | — | PAY (data kehadiran) | HRM (jadwal cuti) |
| **REC** | — | HRM (data karyawan baru) | WFL (approval), DOC (kontrak) |
| **USR** | — | Seluruh modul (autentikasi & otorisasi) | ADM (config), HRM (link karyawan), AUD (log login), MSG (notifikasi keamanan) |
| **RPT** | Seluruh modul ERP | — (standalone) | — |
| **ADM** | — | Seluruh modul (konfigurasi) | USR (autentikasi) |
| **WFL** | — | MSG (notifikasi approval) | PUR, SAL, HRM, BUD, REC (approval) |
| **DOC** | — | — | Seluruh modul (lampiran dokumen) |
| **MSG** | WFL (approval), CRM (follow-up), modul lainnya | — (standalone) | — |
| **AUD** | Seluruh modul ERP (setiap aksi tercatat) | — (standalone) | — |

---

## Alur Proses Utama

### Procure-to-Pay (Pembelian → Pembayaran)

```
Purchase Requisition → RFQ → Purchase Order → Goods Receipt → Invoice Verification → Payment
      [PUR]           [PUR]     [PUR]            [INV]            [AP]              [AP/ACC]
```

### Order-to-Cash (Penjualan → Penerimaan)

```
Lead/Opportunity → Quotation → Sales Order → Delivery → Invoice → Payment Receipt
     [CRM]          [SAL]       [SAL]        [INV]      [AR]        [AR/ACC]
```

### Hire-to-Retire (Rekrutmen → Pensiun)

```
Manpower Plan → Job Posting → Selection → Onboarding → Active Employee → Payroll → Exit
    [REC]         [REC]        [REC]       [REC/HRM]      [HRM]         [PAY]    [HRM]
```

### Plan-to-Produce (Perencanaan → Produksi)

```
Demand Forecast → Production Plan → Work Order → Material Issue → Production → QC → Receipt
    [SAL]            [MFG]           [MFG]         [INV]          [MFG]      [MFG]  [INV]
```

---

## Catatan Teknis

> **Arsitektur:** Modul-modul dirancang secara modular dan dapat diaktifkan/dinonaktifkan sesuai kebutuhan bisnis.
>
> **Database:** Setiap modul berbagi database tunggal (single database) untuk memastikan integritas data dan real-time reporting.
>
> **API:** Setiap modul menyediakan RESTful API untuk integrasi dengan sistem eksternal.
>
> **Multi-tenant:** Mendukung arsitektur multi-tenant untuk SaaS deployment.

---

## Urutan Implementasi Modul (Roadmap)

> Urutan berikut disusun berdasarkan **analisis ketergantungan teknis** (topological dependency) — modul yang menjadi fondasi dan tidak bergantung pada modul lain dibangun lebih dahulu, modul yang membutuhkan data dari modul lain dibangun setelah dependensinya tersedia.

### Fase 1 — Fondasi Sistem & Keamanan `[Minggu 1–4]`

| Urutan | Modul | Alasan Prioritas |
|--------|-------|-------------------|
| 1 | **ADM** — Administrasi & Pengaturan | Menyediakan konfigurasi dasar (profil perusahaan, mata uang, regional, penomoran dokumen) yang dibutuhkan **seluruh modul**. Harus ada pertama. |
| 2 | **USR** — User Management & RBAC | Menyediakan autentikasi, otorisasi, & role-based access control untuk **seluruh modul**. Tanpa ini tidak ada pengguna yang bisa login. |
| 3 | **AUD** — Audit Trail & Log Keamanan | Harus aktif sejak hari pertama agar setiap aksi konfigurasi dan data entry sudah tercatat dalam audit log. Bersifat *append-only* dan tidak bergantung pada modul lain. |

---

### Fase 2 — Fondasi Keuangan `[Minggu 3–6]`

| Urutan | Modul | Alasan Prioritas |
|--------|-------|-------------------|
| 4 | **GL** — Buku Besar (General Ledger) | Pusat pencatatan seluruh transaksi keuangan. Chart of Accounts (CoA), periode fiskal, dan dimensi akuntansi harus didefinisikan sebelum modul keuangan lain bisa melakukan posting jurnal. |
| 5 | **ACC** — Akuntansi & Keuangan | Mesin posting jurnal utama. Menerima entri dari AP, AR, FA, PAY, SAL, PUR — sehingga harus siap sebagai *penerima* sebelum modul-modul tersebut dibangun. |

---

### Fase 3 — Infrastruktur Pendukung Lintas Modul `[Minggu 5–8]`

| Urutan | Modul | Alasan Prioritas |
|--------|-------|-------------------|
| 6 | **DOC** — Manajemen Dokumen | Dibutuhkan oleh PUR (lampiran PO/kontrak), REC (dokumen kontrak kerja), CRM (proposal), dan hampir seluruh modul lain untuk penyimpanan berkas. Harus siap sebelum modul transaksional. |
| 7 | **WFL** — Workflow & Approval | Dibutuhkan oleh PUR (approval PO), SAL (approval diskon), HRM (approval cuti), BUD (approval anggaran), REC (approval rekrutmen). Mesin persetujuan bertingkat harus ada sebelum transaksi yang memerlukan otorisasi. |
| 8 | **MSG** — Pesan & Notifikasi | Menerima output dari WFL (notifikasi approval) dan CRM (follow-up reminder). Gateway WhatsApp/Email/Push harus dikonfigurasi agar workflow approval bisa mengirim notifikasi. |

---

### Fase 4 — Siklus Pembelian & Persediaan (Procure-to-Pay) `[Minggu 7–12]`

| Urutan | Modul | Alasan Prioritas |
|--------|-------|-------------------|
| 9 | **PUR** — Pembelian (Purchasing) | Modul transaksional pertama. Menghasilkan output ke INV (stok masuk), AP (invoice), dan FA (aset baru). Terintegrasi dengan WFL (approval PO) dan DOC (lampiran kontrak) yang sudah tersedia dari Fase 3. |
| 10 | **INV** — Persediaan & Gudang | Menerima barang masuk dari PUR (Goods Receipt). Master barang, multi-gudang, valuasi stok (FIFO/Average), dan barcode tracking harus siap sebelum SAL bisa mengeluarkan stok. |
| 11 | **AP** — Hutang Usaha | Menerima invoice dari PUR (three-way matching: PO ↔ GR ↔ Invoice). Melakukan posting jurnal hutang ke ACC/GL. Siklus Procure-to-Pay lengkap setelah AP aktif. |
| 12 | **FA** — Aset Tetap | Menerima data pembelian aset dari PUR. Menghitung penyusutan dan memposting ke ACC. Bisa dibangun paralel dengan AP karena sama-sama menerima dari PUR. |

---

### Fase 5 — Siklus Penjualan & Piutang (Order-to-Cash) `[Minggu 11–16]`

| Urutan | Modul | Alasan Prioritas |
|--------|-------|-------------------|
| 13 | **CRM** — Customer Relationship Management | Tidak memiliki dependensi input — hanya menghasilkan output ke SAL (konversi leads → quotation → sales order). Dibangun lebih dulu agar pipeline penjualan sudah terisi saat SAL aktif. |
| 14 | **SAL** — Penjualan (Sales) | Menerima leads converted dari CRM. Mengeluarkan stok dari INV (Delivery Order), membuat invoice ke AR, dan memposting revenue ke ACC. INV harus sudah aktif (Fase 4). |
| 15 | **AR** — Piutang Usaha | Menerima invoice penjualan dari SAL. Memposting jurnal piutang ke ACC/GL. Siklus Order-to-Cash lengkap setelah AR aktif. |
| 16 | **POS** — Point of Sale | Menghasilkan output ke INV (pengurangan stok real-time), AR (piutang jika kredit), dan ACC (revenue). Membutuhkan INV dan AR yang sudah aktif. Bisa dibangun paralel dengan AR. |

---

### Fase 6 — Perpajakan & Anggaran `[Minggu 15–18]`

| Urutan | Modul | Alasan Prioritas |
|--------|-------|-------------------|
| 17 | **TAX** — Perpajakan | Menerima data dari AP (PPN masukan), AR (PPN keluaran), dan PAY (PPh 21). AP dan AR harus sudah aktif. PPh 21 dari PAY bisa ditambahkan setelah Fase 7. |
| 18 | **BUD** — Anggaran (Budgeting) | Menerima data realisasi dari ACC dan GL untuk perbandingan budget vs actual. Membutuhkan data transaksi keuangan yang sudah berjalan (Fase 4–5). Terintegrasi dengan WFL untuk approval. |

---

### Fase 7 — Siklus SDM & Penggajian (Hire-to-Retire) `[Minggu 17–24]`

| Urutan | Modul | Alasan Prioritas |
|--------|-------|-------------------|
| 19 | **HRM** — Sumber Daya Manusia | Fondasi data kepegawaian — master karyawan, struktur organisasi, jabatan, kontrak, cuti. Harus ada sebelum REC bisa mengirim data karyawan baru, dan sebelum PAY bisa menghitung gaji. |
| 20 | **REC** — Rekrutmen | Menghasilkan output ke HRM (data karyawan baru setelah onboarding). Terintegrasi dengan WFL (approval) dan DOC (kontrak kerja) yang sudah tersedia. |
| 21 | **ATT** — Absensi & Kehadiran | Menghasilkan output ke PAY (data kehadiran & jam lembur). Terintegrasi dengan HRM (jadwal cuti & shift). HRM harus sudah aktif. |
| 22 | **PAY** — Penggajian (Payroll) | Menerima data dari HRM (komponen gaji) dan ATT (kehadiran & lembur). Memposting biaya gaji ke ACC dan PPh 21 ke TAX. Modul SDM terakhir karena bergantung pada HRM + ATT. |

---

### Fase 8 — Operasional & Produksi `[Minggu 23–28]`

| Urutan | Modul | Alasan Prioritas |
|--------|-------|-------------------|
| 23 | **MFG** — Manufaktur / Produksi | Menerima bahan baku dari INV dan demand/forecast dari SAL. Mengembalikan barang jadi ke INV dan memposting biaya produksi ke ACC. Membutuhkan INV dan SAL yang sudah matang. |
| 24 | **PRJ** — Manajemen Proyek | Terintegrasi lintas-kategori dengan HRM (alokasi SDM), PUR (pembelian material), ACC (biaya proyek), dan BUD (anggaran proyek). Membutuhkan hampir seluruh modul pendukung sudah aktif. |

---

### Fase 9 — Analitik & Pelaporan `[Minggu 27–30]`

| Urutan | Modul | Alasan Prioritas |
|--------|-------|-------------------|
| 25 | **RPT** — Laporan & Analitik | Menerima data dari **seluruh modul ERP**. Dashboard eksekutif, report builder, dan scheduled reporting paling efektif jika seluruh modul sudah menghasilkan data transaksional. Dibangun terakhir. |

---

### Visualisasi Timeline Implementasi

```
Minggu:  1    4    8    12   16   20   24   28   30
         │    │    │    │    │    │    │    │    │
Fase 1:  ████████                                      ADM, USR, AUD
Fase 2:    ████████                                    GL, ACC
Fase 3:        ████████                                DOC, WFL, MSG
Fase 4:            ████████████                        PUR, INV, AP, FA
Fase 5:                ████████████                    CRM, SAL, AR, POS
Fase 6:                    ████████                    TAX, BUD
Fase 7:                        ████████████            HRM, REC, ATT, PAY
Fase 8:                            ████████████        MFG, PRJ
Fase 9:                                    ████████    RPT
```

### Catatan Implementasi

> **Paralelisasi:** Fase-fase yang saling overlap pada timeline di atas menandakan bahwa tim berbeda bisa mengerjakan modul secara paralel. Misalnya, tim Keuangan bisa mengerjakan Fase 5 (SAL/AR) sementara tim SDM memulai Fase 7 (HRM).
>
> **MVP vs Full:** Setiap modul bisa di-deploy secara bertahap — fitur inti (MVP) di-release lebih dulu, fitur lanjutan (integrasi penuh, laporan kompleks) di-release pada iterasi berikutnya.
>
> **Data Migration:** Migrasi data master (CoA, vendor, customer, karyawan, item) sebaiknya dilakukan pada fase yang relevan — jangan menunggu hingga seluruh modul selesai.

---

*Dokumen ini akan diperbarui seiring dengan perkembangan kebutuhan dan pengembangan sistem ERP.*
