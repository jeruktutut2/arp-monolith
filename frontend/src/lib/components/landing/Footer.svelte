<script lang="ts">
	import { resolve } from '$app/paths';

	let selectedLang = $state<'ID' | 'EN'>('ID');
	let selectedCurrency = $state<'IDR' | 'USD'>('IDR');
</script>

<footer
	class="border-t border-slate-200 bg-white text-slate-600 transition-colors duration-200 dark:border-slate-800 dark:bg-slate-950 dark:text-slate-400"
>
	<!-- Main Multi-Column Section -->
	<div class="mx-auto max-w-7xl px-4 py-16 sm:px-6 lg:px-8">
		<div class="grid grid-cols-1 gap-10 md:grid-cols-2 lg:grid-cols-6">
			<!-- Col 1: Brand & Compliance (Span 2 cols on lg) -->
			<div class="space-y-4 lg:col-span-2">
				<a href={resolve('/')} class="flex items-center gap-3">
					<div
						class="flex h-10 w-10 items-center justify-center rounded-xl bg-gradient-to-br from-indigo-600 to-blue-600 text-lg font-black text-white shadow-md shadow-indigo-500/20"
					>
						E
					</div>
					<span class="text-xl font-black tracking-tight text-slate-900 dark:text-white">
						ERP<span class="text-indigo-600 dark:text-indigo-400">System</span>
					</span>
				</a>

				<p class="max-w-sm text-xs leading-relaxed text-slate-500 dark:text-slate-400">
					Platform Enterprise Resource Planning terpadu berbasis arsitektur Hexagonal Modular
					Monolith dan PgBouncer transaction pooling. Mengintegrasikan seluruh siklus bisnis dari
					akuntansi hingga logistik.
				</p>

				<!-- Compliance Badges -->
				<div class="pt-2">
					<span
						class="text-[11px] font-bold tracking-wider text-slate-400 uppercase dark:text-slate-500"
					>
						Standar & Kepatuhan:
					</span>
					<div class="mt-2 flex flex-wrap gap-2">
						<span
							class="inline-flex items-center gap-1 rounded-md bg-slate-100 px-2.5 py-1 text-[11px] font-semibold text-slate-700 dark:bg-slate-800 dark:text-slate-300"
						>
							<span class="text-emerald-500">✓</span> SOC-2 Type II
						</span>
						<span
							class="inline-flex items-center gap-1 rounded-md bg-slate-100 px-2.5 py-1 text-[11px] font-semibold text-slate-700 dark:bg-slate-800 dark:text-slate-300"
						>
							<span class="text-emerald-500">✓</span> ISO 27001
						</span>
						<span
							class="inline-flex items-center gap-1 rounded-md bg-slate-100 px-2.5 py-1 text-[11px] font-semibold text-slate-700 dark:bg-slate-800 dark:text-slate-300"
						>
							<span class="text-emerald-500">✓</span> PSAK / IFRS
						</span>
					</div>
				</div>

				<!-- Live Uptime Status Pill -->
				<div class="pt-2">
					<div
						class="inline-flex items-center gap-2 rounded-full border border-emerald-200 bg-emerald-50 px-3 py-1 text-xs font-semibold text-emerald-800 dark:border-emerald-900 dark:bg-emerald-950/60 dark:text-emerald-300"
					>
						<span class="h-2 w-2 animate-pulse rounded-full bg-emerald-500"></span>
						<span>All Systems Operational (99.99%)</span>
					</div>
				</div>
			</div>

			<!-- Col 2: Keuangan & Akuntansi -->
			<div>
				<h3 class="text-xs font-bold tracking-wider text-slate-900 uppercase dark:text-white">
					Keuangan
				</h3>
				<ul class="mt-4 space-y-2 text-xs">
					<li>
						<a href="#modules" class="transition hover:text-indigo-600 dark:hover:text-indigo-400">
							Akuntansi (ACC)
						</a>
					</li>
					<li>
						<a href="#modules" class="transition hover:text-indigo-600 dark:hover:text-indigo-400">
							Hutang Usaha (AP)
						</a>
					</li>
					<li>
						<a href="#modules" class="transition hover:text-indigo-600 dark:hover:text-indigo-400">
							Piutang Usaha (AR)
						</a>
					</li>
					<li>
						<a href="#modules" class="transition hover:text-indigo-600 dark:hover:text-indigo-400">
							Buku Besar (GL)
						</a>
					</li>
					<li>
						<a href="#modules" class="transition hover:text-indigo-600 dark:hover:text-indigo-400">
							Aset Tetap (FA)
						</a>
					</li>
					<li>
						<a href="#modules" class="transition hover:text-indigo-600 dark:hover:text-indigo-400">
							Anggaran (BUD)
						</a>
					</li>
					<li>
						<a href="#modules" class="transition hover:text-indigo-600 dark:hover:text-indigo-400">
							Perpajakan (TAX)
						</a>
					</li>
				</ul>
			</div>

			<!-- Col 3: Pasok & Penjualan -->
			<div>
				<h3 class="text-xs font-bold tracking-wider text-slate-900 uppercase dark:text-white">
					Supply Chain & CRM
				</h3>
				<ul class="mt-4 space-y-2 text-xs">
					<li>
						<a href="#modules" class="transition hover:text-indigo-600 dark:hover:text-indigo-400">
							Pembelian (PUR)
						</a>
					</li>
					<li>
						<a href="#modules" class="transition hover:text-indigo-600 dark:hover:text-indigo-400">
							Persediaan Gudang (INV)
						</a>
					</li>
					<li>
						<a href="#modules" class="transition hover:text-indigo-600 dark:hover:text-indigo-400">
							Penjualan (SAL)
						</a>
					</li>
					<li>
						<a href="#modules" class="transition hover:text-indigo-600 dark:hover:text-indigo-400">
							Customer CRM
						</a>
					</li>
					<li>
						<a href="#modules" class="transition hover:text-indigo-600 dark:hover:text-indigo-400">
							Kasir Point of Sale (POS)
						</a>
					</li>
					<li>
						<a href="#modules" class="transition hover:text-indigo-600 dark:hover:text-indigo-400">
							Manufaktur (MFG)
						</a>
					</li>
					<li>
						<a href="#modules" class="transition hover:text-indigo-600 dark:hover:text-indigo-400">
							Proyek & Gantt (PRJ)
						</a>
					</li>
				</ul>
			</div>

			<!-- Col 4: SDM & Sistem -->
			<div>
				<h3 class="text-xs font-bold tracking-wider text-slate-900 uppercase dark:text-white">
					SDM & Platform
				</h3>
				<ul class="mt-4 space-y-2 text-xs">
					<li>
						<a href="#modules" class="transition hover:text-indigo-600 dark:hover:text-indigo-400">
							SDM / HRM
						</a>
					</li>
					<li>
						<a href="#modules" class="transition hover:text-indigo-600 dark:hover:text-indigo-400">
							Penggajian (PAY)
						</a>
					</li>
					<li>
						<a href="#modules" class="transition hover:text-indigo-600 dark:hover:text-indigo-400">
							Absensi Biometrik (ATT)
						</a>
					</li>
					<li>
						<a href="#modules" class="transition hover:text-indigo-600 dark:hover:text-indigo-400">
							User & RBAC (USR)
						</a>
					</li>
					<li>
						<a href="#modules" class="transition hover:text-indigo-600 dark:hover:text-indigo-400">
							Workflow Desainer (WFL)
						</a>
					</li>
					<li>
						<a href="#modules" class="transition hover:text-indigo-600 dark:hover:text-indigo-400">
							Laporan Eksekutif (RPT)
						</a>
					</li>
					<li>
						<a href="#modules" class="transition hover:text-indigo-600 dark:hover:text-indigo-400">
							Audit Trail Diff (AUD)
						</a>
					</li>
				</ul>
			</div>

			<!-- Col 5: Bantuan & Legal -->
			<div>
				<h3 class="text-xs font-bold tracking-wider text-slate-900 uppercase dark:text-white">
					Bantuan & Legal
				</h3>
				<ul class="mt-4 space-y-2 text-xs">
					<li>
						<a href="#privacy" class="transition hover:text-indigo-600 dark:hover:text-indigo-400">
							Kebijakan Privasi
						</a>
					</li>
					<li>
						<a href="#terms" class="transition hover:text-indigo-600 dark:hover:text-indigo-400">
							Syarat & Ketentuan
						</a>
					</li>
					<li>
						<a href="#sla" class="transition hover:text-indigo-600 dark:hover:text-indigo-400">
							Jaminan SLA Layanan
						</a>
					</li>
					<li>
						<a href="#security" class="transition hover:text-indigo-600 dark:hover:text-indigo-400">
							Praktik Keamanan Siber
						</a>
					</li>
					<li>
						<a href="#contact" class="transition hover:text-indigo-600 dark:hover:text-indigo-400">
							Kontak Support 24/7
						</a>
					</li>
				</ul>
			</div>
		</div>
	</div>

	<!-- Bottom Legal Bar with Localization Switcher -->
	<div
		class="border-t border-slate-200 bg-slate-50/80 px-4 py-4 text-xs dark:border-slate-800 dark:bg-slate-900/60"
	>
		<div
			class="mx-auto flex max-w-7xl flex-col items-center justify-between gap-4 sm:flex-row sm:px-6 lg:px-8"
		>
			<!-- Copyright & Version -->
			<div class="flex items-center gap-2">
				<span>© 2026 ERPSystem Inc. All rights reserved.</span>
				<span
					class="rounded-md bg-slate-200/80 px-1.5 py-0.5 font-mono text-[10px] font-bold dark:bg-slate-800"
				>
					v2.4.0-enterprise
				</span>
			</div>

			<!-- Localization Selector Controls -->
			<div class="flex items-center gap-4">
				<!-- Language Selector -->
				<div class="flex items-center gap-1">
					<span class="text-slate-400">Bahasa:</span>
					<div
						class="inline-flex rounded-lg border border-slate-200 bg-white p-0.5 dark:border-slate-700 dark:bg-slate-800"
					>
						<button
							type="button"
							onclick={() => (selectedLang = 'ID')}
							class="cursor-pointer rounded-md px-2 py-0.5 text-[11px] font-bold transition {selectedLang ===
							'ID'
								? 'bg-indigo-600 text-white'
								: 'text-slate-600 hover:text-slate-900 dark:text-slate-300 dark:hover:text-white'}"
						>
							ID
						</button>
						<button
							type="button"
							onclick={() => (selectedLang = 'EN')}
							class="cursor-pointer rounded-md px-2 py-0.5 text-[11px] font-bold transition {selectedLang ===
							'EN'
								? 'bg-indigo-600 text-white'
								: 'text-slate-600 hover:text-slate-900 dark:text-slate-300 dark:hover:text-white'}"
						>
							EN
						</button>
					</div>
				</div>

				<!-- Currency Selector -->
				<div class="flex items-center gap-1">
					<span class="text-slate-400">Mata Uang:</span>
					<div
						class="inline-flex rounded-lg border border-slate-200 bg-white p-0.5 dark:border-slate-700 dark:bg-slate-800"
					>
						<button
							type="button"
							onclick={() => (selectedCurrency = 'IDR')}
							class="cursor-pointer rounded-md px-2 py-0.5 text-[11px] font-bold transition {selectedCurrency ===
							'IDR'
								? 'bg-indigo-600 text-white'
								: 'text-slate-600 hover:text-slate-900 dark:text-slate-300 dark:hover:text-white'}"
						>
							IDR (Rp)
						</button>
						<button
							type="button"
							onclick={() => (selectedCurrency = 'USD')}
							class="cursor-pointer rounded-md px-2 py-0.5 text-[11px] font-bold transition {selectedCurrency ===
							'USD'
								? 'bg-indigo-600 text-white'
								: 'text-slate-600 hover:text-slate-900 dark:text-slate-300 dark:hover:text-white'}"
						>
							USD ($)
						</button>
					</div>
				</div>
			</div>
		</div>
	</div>

	<!-- Live Telemetry Status Bar (from 0_FTR/footer.md) -->
	<div
		class="border-t border-slate-200 bg-slate-100/90 px-4 py-2 font-mono text-[11px] text-slate-500 dark:border-slate-800/90 dark:bg-slate-950 dark:text-slate-400"
	>
		<div
			class="mx-auto flex max-w-7xl flex-wrap items-center justify-between gap-3 sm:px-6 lg:px-8"
		>
			<div class="flex items-center gap-4">
				<span class="inline-flex items-center gap-1.5">
					<span class="h-2 w-2 rounded-full bg-emerald-500"></span>
					<span>Gateway: <strong class="text-slate-700 dark:text-slate-200">8ms</strong></span>
				</span>
				<span class="inline-flex items-center gap-1.5">
					<span class="h-2 w-2 rounded-full bg-emerald-500"></span>
					<span>WS: <strong class="text-slate-700 dark:text-slate-200">Connected</strong></span>
				</span>
			</div>

			<div class="flex items-center gap-4">
				<span
					>Cabang: <strong class="text-slate-700 dark:text-slate-200">Jakarta HQ (Main)</strong
					></span
				>
				<span
					>Tahun Fiskal: <strong class="text-indigo-600 dark:text-indigo-400">FY 2026</strong></span
				>
			</div>
		</div>
	</div>
</footer>
