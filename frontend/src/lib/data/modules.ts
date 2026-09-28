import type { ERPModule, SolutionItem } from '$lib/types/landing';

export const ERP_MODULES: ERPModule[] = [
	// 1. Keuangan (Finance & Accounting)
	{
		id: '1_ACC',
		code: 'ACC',
		name: 'Akuntansi & Keuangan',
		category: 'Keuangan',
		description:
			'Pencatatan transaksi keuangan terpadu, bagan akun (CoA) hierarkis, jurnal umum, rekonsiliasi bank, dan laporan keuangan komprehensif.',
		icon: '💰',
		color: 'from-blue-500 to-indigo-600',
		badge: 'Core Financials',
		features: [
			'Chart of Accounts (CoA) multi-level tak terbatas',
			'Jurnal umum debit/kredit dengan multi-currency',
			'Neraca, Laba Rugi, & Arus Kas real-time',
			'Rekonsiliasi mutasi bank otomatis',
			'Konsolidasi holding & multi-cabang'
		],
		inputsFrom: ['AP', 'AR', 'FA', 'PAY', 'SAL', 'PUR'],
		outputsTo: ['RPT', 'BUD', 'TAX']
	},
	{
		id: '2_AP',
		code: 'AP',
		name: 'Hutang Usaha (Accounts Payable)',
		category: 'Keuangan',
		description:
			'Manajemen kewajiban vendor dengan three-way matching otomatis antara PO, Goods Receipt, dan Vendor Invoice.',
		icon: '📑',
		color: 'from-blue-600 to-cyan-600',
		features: [
			'Verifikasi invoice supplier & Three-way matching',
			'Jadwal pembayaran bertahap (payment terms)',
			'Pembayaran massal (bulk/batch payment)',
			'Penerbitan Debit Note & penanganan retur beli',
			'Analisis umur hutang (Aging Payable analysis)'
		],
		inputsFrom: ['PUR'],
		outputsTo: ['ACC', 'TAX']
	},
	{
		id: '3_AR',
		code: 'AR',
		name: 'Piutang Usaha (Accounts Receivable)',
		category: 'Keuangan',
		description:
			'Pengelolaan penagihan dan piutang pelanggan, batas kredit (credit limit), credit note, dan pelunasan bertahap.',
		icon: '💳',
		color: 'from-cyan-500 to-teal-600',
		features: [
			'Penerbitan faktur penjualan terintegrasi',
			'Pencatatan penerimaan pelunasan piutang',
			'Pengawasan limit kredit & batas jatuh tempo',
			'Reminder penagihan otomatis via email & pesan',
			'Analisis umur piutang (Aging Receivable analysis)'
		],
		inputsFrom: ['SAL'],
		outputsTo: ['ACC', 'TAX']
	},
	{
		id: '4_GL',
		code: 'GL',
		name: 'Buku Besar (General Ledger)',
		category: 'Keuangan',
		description:
			'Pusat pencatatan seluruh buku besar akun, neraca saldo (trial balance), jurnal penyesuaian, dan eliminasi transaksi intercompany.',
		icon: '📖',
		color: 'from-indigo-600 to-violet-600',
		badge: 'PSAK / IFRS',
		features: [
			'Buku besar terinci per sub-akun & cost center',
			'Neraca Saldo (Trial Balance) instan & akurat',
			'Jurnal penyesuaian & jurnal penutup akhir periode',
			'Eliminasi transaksi antar-entitas (intercompany)',
			'Dimensi multi-cabang & alokasi proyek'
		],
		inputsFrom: ['ACC', 'AP', 'AR', 'FA', 'PAY'],
		outputsTo: ['RPT', 'BUD']
	},
	{
		id: '5_FA',
		code: 'FA',
		name: 'Aset Tetap (Fixed Assets)',
		category: 'Keuangan',
		description:
			'Siklus hidup aset tetap mulai dari kapitalisasi pengadaan, perhitungan depresiasi otomatis bulanan, hingga mutasi dan disposal.',
		icon: '🏢',
		color: 'from-sky-500 to-blue-600',
		features: [
			'Metode penyusutan garis lurus & saldo menurun',
			'Jadwal depresiasi otomatis terjadwal',
			'Pelacakan nomor seri, barcode & QR fisik aset',
			'Mutasi lokasi & transfer tanggung jawab aset',
			'Pelepasan aset (disposal) & kalkulasi untung/rugi'
		],
		inputsFrom: ['PUR'],
		outputsTo: ['ACC']
	},
	{
		id: '6_BUD',
		code: 'BUD',
		name: 'Anggaran (Budgeting)',
		category: 'Keuangan',
		description:
			'Perencanaan plafon anggaran departemen dan proyek, validasi budget ceiling sebelum transaksi, serta analisa varians realisasi.',
		icon: '📊',
		color: 'from-teal-500 to-emerald-600',
		features: [
			'Penyusunan anggaran tahunan per divisi/cost center',
			'Pencegahan overbudgeting transaksi real-time',
			'Revisi & persetujuan perubahan anggaran berjenjang',
			'Analisis komparasi Budget vs Actual dinamis',
			'Peramalan arus kas & kebutuhan belanja modal'
		],
		inputsFrom: ['ACC', 'GL'],
		outputsTo: ['RPT', 'PUR']
	},
	{
		id: '7_TAX',
		code: 'TAX',
		name: 'Perpajakan & Compliance',
		category: 'Keuangan',
		description:
			'Kepatuhan regulasi pajak Indonesia (PPN, PPh 21/23/4 ayat 2), integrasi e-Faktur, e-Bupot, dan persiapan SPT masa otomatis.',
		icon: '⚖️',
		color: 'from-emerald-600 to-green-600',
		badge: 'DJP Compliant',
		features: [
			'Perhitungan PPN 11% & 12% otomatis per transaksi',
			'Ekspor & impor format e-Faktur DJP Indonesia',
			'Pengelolaan bukti potong PPh e-Bupot unifikasi',
			'Rekonsiliasi pajak masukan vs keluaran terpadu',
			'Kompilasi lampiran SPT Masa & SPT Tahunan'
		],
		inputsFrom: ['AP', 'AR', 'PAY'],
		outputsTo: ['RPT', 'ACC']
	},

	// 2. Rantai Pasok & Persediaan (Supply Chain & Inventory)
	{
		id: '8_PUR',
		code: 'PUR',
		name: 'Pembelian (Purchasing)',
		category: 'Rantai Pasok',
		description:
			'Siklus pengadaan terstruktur dari Purchase Requisition, Request for Quotation (RFQ), evaluasi vendor komparatif, hingga Purchase Order.',
		icon: '🛒',
		color: 'from-amber-500 to-orange-600',
		badge: 'Smart Sourcing',
		features: [
			'Permintaan pembelian (Purchase Requisition) berjenjang',
			'Komparasi matriks penawaran harga vendor (RFQ)',
			'Penerbitan PO multi-mata uang dan termin pengiriman',
			'Penerimaan barang fisik (Goods Receipt Note / GRN)',
			'Evaluasi rating performa & SLA ketepatan vendor'
		],
		inputsFrom: ['INV', 'BUD'],
		outputsTo: ['AP', 'INV']
	},
	{
		id: '9_INV',
		code: 'INV',
		name: 'Persediaan & Gudang (Inventory)',
		category: 'Rantai Pasok',
		description:
			'Manajemen persediaan multi-gudang, penomoran batch & serial number, pergerakan stok real-time, dan stock opname fisik berakurasi tinggi.',
		icon: '📦',
		color: 'from-orange-500 to-amber-600',
		badge: 'Multi-Warehouse',
		features: [
			'Hierarki multi-gudang, lorong, rak, & bin lokasi',
			'Pelacakan nomor lot, batch produksi, & expiry date',
			'Kalkulasi nilai persediaan metode FIFO dan Rata-rata',
			'Surat jalan transfer antar-gudang dan cabang',
			'Stock opname dengan pemindaian barcode mobile'
		],
		inputsFrom: ['PUR', 'MFG', 'SAL'],
		outputsTo: ['ACC', 'SAL', 'MFG']
	},

	// 3. Penjualan & CRM (Sales, CRM & POS)
	{
		id: '10_SAL',
		code: 'SAL',
		name: 'Penjualan (Sales)',
		category: 'Penjualan & CRM',
		description:
			'Orkestrasi alur penjualan B2B/B2C mulai dari penawaran harga (Quotation), konfirmasi Sales Order, surat jalan pengiriman, hingga invoice.',
		icon: '📈',
		color: 'from-rose-500 to-red-600',
		features: [
			'Penyusunan penawaran harga (Quotation) interaktif',
			'Konfirmasi Sales Order & reservasi kuota stok otomatis',
			'Penerbitan Surat Jalan pengiriman bertahap (DO)',
			'Dukungan multi-pricelist, diskon bertingkat, & promo',
			'Analisis performa sales rep & target komisi'
		],
		inputsFrom: ['CRM', 'INV'],
		outputsTo: ['AR', 'INV']
	},
	{
		id: '11_CRM',
		code: 'CRM',
		name: 'Customer Relationship Management',
		category: 'Penjualan & CRM',
		description:
			'Pipeline penjualan visual Kanban untuk konversi lead, manajemen kontak pelanggan, log interaksi meeting/telepon, dan layanan tiket bantuan.',
		icon: '🤝',
		color: 'from-pink-500 to-rose-600',
		features: [
			'Visual Kanban sales opportunity & deal pipeline',
			'Skoring lead otomatis dan alokasi ke sales team',
			'Riwayat sentuhan pelanggan (call, email, visit log)',
			'Helpdesk ticket pelanggan terintegrasi dengan SLA',
			'Kampanye marketing email & broadcast segmentasi'
		],
		inputsFrom: ['MSG'],
		outputsTo: ['SAL']
	},
	{
		id: '12_POS',
		code: 'POS',
		name: 'Point of Sale (Retail & F&B)',
		category: 'Penjualan & CRM',
		description:
			'Antarmuka kasir ritel berkecepatan tinggi dengan navigasi keyboard shortcuts (F1-F12), scanner barcode global, dan cetak struk ESC/POS via WebUSB.',
		icon: '🏪',
		color: 'from-red-500 to-pink-600',
		badge: 'Ultra Fast POS',
		features: [
			'Shortcut keyboard cepat tanpa mouse untuk kasir',
			'Intersepsi scanner barcode fisik tanpa bentrok input',
			'Direct thermal receipt printing tanpa dialog OS (ESC/POS)',
			'Dukungan multi-metode pembayaran: QRIS, EDC, Tunai',
			'Buka & tutup kasir harian dengan rekonsiliasi kas register'
		],
		inputsFrom: ['INV', 'SAL'],
		outputsTo: ['ACC', 'INV']
	},

	// 4. Operasional (Manufacturing & Projects)
	{
		id: '13_MFG',
		code: 'MFG',
		name: 'Manufaktur & Produksi',
		category: 'Operasional',
		description:
			'Perencanaan kebutuhan bahan (MRP), Bill of Materials (BOM) multi-level, Perintah Kerja (Work Order), dan kalkulasi harga pokok produksi (HPP).',
		icon: '⚙️',
		color: 'from-slate-600 to-zinc-700',
		features: [
			'Bill of Materials (BOM) multi-level & rumus resep',
			'Work Order (WO) dan routing stasiun kerja mesin',
			'Kalkulasi Harga Pokok Produksi (HPP) presisi desimal',
			'Pencatatan barang dalam proses (Work In Progress / WIP)',
			'Pemeriksaan kualitas hasil produksi (Quality Control)'
		],
		inputsFrom: ['INV', 'PUR'],
		outputsTo: ['INV', 'ACC']
	},
	{
		id: '14_PRJ',
		code: 'PRJ',
		name: 'Manajemen Proyek',
		category: 'Operasional',
		description:
			'Manajemen proyek konstruksi & layanan dengan Gantt Chart interaktif (Frappe Gantt), pelacakan tonggak pencapaian (milestone), dan timesheet biaya tim.',
		icon: '📅',
		color: 'from-violet-600 to-purple-700',
		badge: 'Interactive Gantt',
		features: [
			'Visualisasi linimasa interaktif (Frappe Gantt)',
			'Ketergantungan tugas (task dependencies) drag-and-drop',
			'Pencatatan timesheet jam kerja tenaga ahli/kontraktor',
			'Budgeting proyek dan pemantauan biaya riil vs estimasi',
			'Penagihan kemajuan progres proyek (Progressive Billing)'
		],
		inputsFrom: ['HRM', 'BUD'],
		outputsTo: ['AR', 'ACC']
	},

	// 5. SDM & Penggajian (Human Resources & Payroll)
	{
		id: '15_HRM',
		code: 'HRM',
		name: 'Human Resource Management',
		category: 'SDM',
		description:
			'Pangkalan data karyawan sentral, struktur organisasi hierarkis, pelacakan kontrak kerja, riwayat karir, dan pengajuan cuti mandiri (Employee Self-Service).',
		icon: '👥',
		color: 'from-emerald-500 to-teal-600',
		features: [
			'Direktori data induk profil lengkap karyawan',
			'Bagan struktur organisasi & hierarki jabatan dinamis',
			'Manajemen masa berlaku kontrak & evaluasi kerja',
			'Pengajuan dan persetujuan cuti & izin karyawan (ESS)',
			'Penyimpanan berkas digital identitas & ijazah karyawan'
		],
		inputsFrom: ['REC'],
		outputsTo: ['PAY', 'ATT']
	},
	{
		id: '16_PAY',
		code: 'PAY',
		name: 'Penggajian (Payroll)',
		category: 'SDM',
		description:
			'Perhitungan gaji bulanan otomatis, tunjangan, lembur, BPJS Ketenagakerjaan/Kesehatan, pemotongan pajak PPh 21 tarif efektif (TER), dan slip gaji digital.',
		icon: '💵',
		color: 'from-green-500 to-emerald-600',
		badge: 'BPJS & TER PPh21',
		features: [
			'Kalkulasi formula gaji pokok, tunjangan & insentif',
			'Perhitungan otomatis tarif efektif (TER) PPh Pasal 21',
			'Pemotongan BPJS Kesehatan & Ketenagakerjaan resmi',
			'Distribusi slip gaji elektronik terenkripsi ke email',
			'File transfer gaji massal kompatibel bank (BCA, Mandiri, BRI)'
		],
		inputsFrom: ['HRM', 'ATT'],
		outputsTo: ['ACC', 'TAX']
	},
	{
		id: '17_ATT',
		code: 'ATT',
		name: 'Absensi & Kehadiran',
		category: 'SDM',
		description:
			'Integrasi mesin absensi biometrik sidik jari/wajah, mobile clock-in berbasis GPS geofencing, penjadwalan shift kerja dinamis, dan kalkulasi lembur.',
		icon: '⏱️',
		color: 'from-teal-600 to-cyan-700',
		badge: 'Biometric & GPS',
		features: [
			'Sinkronisasi real-time mesin biometrik (ZKTeco dll.)',
			'Clock-in aplikasi mobile dengan GPS & selfie anti-fake',
			'Pengaturan shift roster bergilir & rotasi kerja',
			'Rekapitulasi jam lembur otomatis terhubung ke payroll',
			'Kalender heatmap kehadiran Apache ECharts'
		],
		inputsFrom: ['HRM'],
		outputsTo: ['PAY']
	},
	{
		id: '18_REC',
		code: 'REC',
		name: 'Rekrutmen & Seleksi',
		category: 'SDM',
		description:
			'Portal pelamar kerja (Career Site), pelacakan pipeline kandidat (ATS), penjadwalan wawancara bertahap, dan konversi rekrutmen menjadi karyawan.',
		icon: '🎯',
		color: 'from-cyan-600 to-blue-700',
		features: [
			'Papan lowongan kerja publik (Career Portal)',
			'Applicant Tracking System (ATS) tahap seleksi kandidat',
			'Penjadwalan wawancara & integrasi kalender tim',
			'Penilaian hasil tes kompetensi & rekam jejak pewawancara',
			'Penerbitan surat penawaran kerja (Offering Letter)'
		],
		inputsFrom: [],
		outputsTo: ['HRM']
	},

	// 6. Sistem, Keamanan & Intelijen (System, Security & Intelligence)
	{
		id: '19_USR',
		code: 'USR',
		name: 'User & Access Security',
		category: 'Sistem',
		description:
			'Manajemen identitas terpusat, kontrol akses berbasis peran (RBAC) granular hingga tingkat field, Single Sign-On (SSO), dan otentikasi MFA.',
		icon: '🔐',
		color: 'from-violet-500 to-indigo-700',
		badge: 'Zero Trust RBAC',
		features: [
			'Role-Based Access Control (RBAC) per aksi & endpoint',
			'Otentikasi Multi-Factor (MFA) dengan TOTP / Authenticator',
			'Proteksi Brute-Force dan monitoring sesi aktif',
			'Kebijakan kata sandi kuat & kedaluwarsa berkala',
			'Dukungan SSO OAuth2, Google Workspace & LDAP'
		],
		inputsFrom: [],
		outputsTo: ['AUD']
	},
	{
		id: '20_RPT',
		code: 'RPT',
		name: 'Laporan & Analitik Eksekutif',
		category: 'Sistem',
		description:
			'Dashboard analitik eksekutif dengan Apache ECharts, peramalan tren bisnis, laporan lintas departemen interaktif, dan ekspor instan PDF/Excel.',
		icon: '📉',
		color: 'from-indigo-600 to-purple-600',
		badge: 'Apache ECharts',
		features: [
			'Dashboard KPI visual interaktif dengan Apache ECharts',
			'Analitik multi-dimensi per cabang, periode, & produk',
			'Penyusun laporan custom (Report Builder) drag-and-drop',
			'Jadwal pengiriman laporan otomatis ke email direksi',
			'Ekspor format Excel, PDF, CSV berkecepatan tinggi'
		],
		inputsFrom: ['ACC', 'SAL', 'PUR', 'INV', 'HRM'],
		outputsTo: []
	},
	{
		id: '21_ADM',
		code: 'ADM',
		name: 'Administrasi & Multi-Tenant',
		category: 'Sistem',
		description:
			'Konfigurasi multi-perusahaan (holding/subsidiaries), hierarki cabang fisik/virtual, penomoran urut dokumen dinamis, dan parameter global ERP.',
		icon: '🛠️',
		color: 'from-blue-700 to-slate-800',
		features: [
			'Pengelolaan multi-entitas perusahaan & struktur holding',
			'Hierarki cabang & pusat pertanggungjawaban laba',
			'Pola penomoran dokumen otomatis fleksibel per modul',
			'Konfigurasi zona waktu, kalender kerja, & kurs valuta',
			'Pengaturan tema perusahaan, logo, & kop surat dokumen'
		],
		inputsFrom: [],
		outputsTo: ['ACC', 'SAL', 'PUR', 'INV']
	},
	{
		id: '22_WFL',
		code: 'WFL',
		name: 'Workflow & Multi-Tier Approval',
		category: 'Sistem',
		description:
			'Perancang alur persetujuan visual (XYFlow @xyflow/svelte) berbasis node & edge, delegasi wewenang saat cuti, dan eskalasi SLA persetujuan otomatis.',
		icon: '🔀',
		color: 'from-purple-600 to-pink-600',
		badge: 'XYFlow Canvas',
		features: [
			'Visual Flowchart Builder node & edge (XYFlow)',
			'Kondisi bertingkat dinamis (misal: nominal > Rp 50jt)',
			'Persetujuan satu klik via email atau aplikasi mobile',
			'Delegasi wewenang penandatangan sementara saat cuti',
			'Notifikasi eskalasi SLA bila approval tertunda'
		],
		inputsFrom: ['PUR', 'SAL', 'PAY', 'PRJ'],
		outputsTo: ['MSG']
	},
	{
		id: '23_DOC',
		code: 'DOC',
		name: 'Manajemen Dokumen Digital',
		category: 'Sistem',
		description:
			'Penyimpanan terpusat dokumen perusahaan, lampiran file per transaksi, riwayat versi dokumen (versioning), dan tanda tangan digital terverifikasi.',
		icon: '📁',
		color: 'from-slate-700 to-sky-800',
		features: [
			'Penyimpanan berkas terenkripsi per transaksi ERP',
			'Riwayat versi dokumen (version control & diff file)',
			'Pencarian teks OCR di dalam berkas faktur & kontrak',
			'Akses hak baca & unduh dokumen berbasis peran',
			'Integrasi tanda tangan digital tersertifikasi'
		],
		inputsFrom: [],
		outputsTo: []
	},
	{
		id: '24_MSG',
		code: 'MSG',
		name: 'Pesan & Notifikasi Sistem',
		category: 'Sistem',
		description:
			'Saluran komunikasi real-time WebSocket dalam aplikasi, blast notifikasi email/WhatsApp gateway, serta peringatan otomatis kondisi darurat.',
		icon: '🔔',
		color: 'from-blue-600 to-indigo-600',
		badge: 'Realtime WS',
		features: [
			'Notifikasi WebSocket instan tanpa reload browser',
			'Gateway integrasi WhatsApp Business API & SMS blast',
			'Pemberitahuan email otomatis invoice jatuh tempo & PO',
			'Ruang diskusi internal per transaksi dokumen ERP',
			'Pusat preferensi langganan notifikasi per pengguna'
		],
		inputsFrom: ['WFL', 'SAL', 'PUR'],
		outputsTo: []
	},
	{
		id: '25_AUD',
		code: 'AUD',
		name: 'Audit Trail & Compliance Diff',
		category: 'Sistem',
		description:
			'Perekaman jejak audit permanen yang tidak dapat diubah (immutable), inspeksi perbedaan nilai data (before vs after diff), dan monitor keamanan.',
		icon: '🛡️',
		color: 'from-slate-800 to-zinc-900',
		badge: 'SOC-2 Compliance',
		features: [
			'Pencatatan data changelog immutable (siapa, kapan, apa)',
			'Visualisasi Before vs After Data Diff per kolom',
			'Pendeteksian aktivitas login anomali & brute-force',
			'Laporan kepatuhan audit PSAK, IFRS, dan SOC-2',
			'Penyimpanan log tamper-proof dengan cryptographic hash'
		],
		inputsFrom: ['USR', 'ACC', 'GL'],
		outputsTo: []
	}
];

export const SOLUTIONS_DATA: SolutionItem[] = [
	{
		title: 'Holding & Korporasi Multi-Entitas',
		tagline: 'Satu platform terpadu untuk puluhan anak perusahaan',
		description:
			'Konsolidasi laporan keuangan otomatis antar-perusahaan, eliminasi transaksi intercompany yang presisi, dan standarisasi operasional seluruh unit bisnis dalam satu database terpusat.',
		modules: ['ACC', 'GL', 'ADM', 'BUD', 'RPT'],
		icon: '🏢',
		highlights: [
			'Multi-company chart of accounts & eliminations',
			'Sentralisasi pengadaan & vendor management',
			'Laporan eksekutif konsolidasi real-time'
		]
	},
	{
		title: 'Pabrikasi & Industri Manufaktur',
		tagline: 'Kendali penuh dari bahan baku hingga barang jadi',
		description:
			'Integrasikan rumus Bill of Materials, perencanaan kapasitas stasiun kerja (routing), penghitungan HPP berbasis desimal presisi, dan pemantauan barang dalam proses (WIP) bebas selisih.',
		modules: ['MFG', 'INV', 'PUR', 'ACC', 'PRJ'],
		icon: '🏭',
		highlights: [
			'Perencanaan kebutuhan bahan otomatis (MRP)',
			'Penetapan harga pokok produksi desimal akurat',
			'Integrasi quality control di lantai produksi'
		]
	},
	{
		title: 'Distribusi & Rantai Pasok Terpadu',
		tagline: 'Perputaran stok cepat, akurat, dan transparan',
		description:
			'Kelola jaringan gudang bertingkat, penomoran lot/batch kadaluarsa, surat jalan otomatis, dan integrasi penagihan pelanggan secara instan tanpa tumpang tindih.',
		modules: ['INV', 'PUR', 'SAL', 'AR', 'AP'],
		icon: '🚚',
		highlights: [
			'Sinkronisasi stok multi-gudang real-time',
			'Three-way matching invoice vendor otomatis',
			'Otomasi sales order hingga surat jalan pengiriman'
		]
	},
	{
		title: 'Retail Chain & Jaringan Kasir POS',
		tagline: 'Transaksi kasir kilat dengan keandalan operasional',
		description:
			'Antarmuka kasir cepat berorientasi keyboard (F1-F12), scanning barcode fisik global, integrasi cetak struk direct thermal via WebUSB, serta rekonsiliasi kas register harian.',
		modules: ['POS', 'INV', 'SAL', 'ACC', 'TAX'],
		icon: '🛍️',
		highlights: [
			'POS keyboard-first tanpa hambatan mouse',
			'Direct thermal receipt printing tanpa dialog print',
			'Kompilasi pajak PPN dan e-Faktur otomatis'
		]
	},
	{
		title: 'Kontraktor & Manajemen Proyek',
		tagline: 'Presisi jadwal dan anggaran dari awal hingga serah terima',
		description:
			'Visualisasikan jadwal pekerjaan dengan Frappe Gantt interaktif, pantau ketergantungan antar-tugas, lacak timesheet konsultan, dan terbitkan termin tagihan progresif.',
		modules: ['PRJ', 'BUD', 'AR', 'PUR', 'DOC'],
		icon: '🏗️',
		highlights: [
			'Interactive Gantt timeline dengan drag-and-drop',
			'Pengawasan realisasi biaya vs plafon anggaran',
			'Penagihan kemajuan progres proyek bersertifikat'
		]
	}
];
