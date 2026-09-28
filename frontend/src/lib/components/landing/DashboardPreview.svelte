<script lang="ts">
	type TabType = 'overview' | 'financials' | 'orders';
	let activeTab = $state<TabType>('overview');
</script>

<section
	id="preview"
	class="border-t border-slate-200/80 bg-slate-50/50 py-20 lg:py-28 dark:border-slate-800/80 dark:bg-slate-900/40"
>
	<div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
		<!-- Section Header -->
		<div class="mx-auto max-w-3xl text-center">
			<div
				class="inline-flex items-center gap-1.5 rounded-full bg-blue-50 px-3.5 py-1 text-xs font-bold text-blue-700 ring-1 ring-blue-500/20 dark:bg-blue-950/70 dark:text-blue-300"
			>
				<span>🖥️</span>
				<span>Executive Cockpit</span>
			</div>
			<h2
				class="mt-4 text-3xl font-black tracking-tight text-slate-900 sm:text-4xl lg:text-5xl dark:text-white"
			>
				Antarmuka Modern, Responsif & Cepat
			</h2>
			<p class="mt-4 text-base text-slate-600 sm:text-lg dark:text-slate-300">
				Pantau kesehatan finansial perusahaan, pipeline penjualan, dan efisiensi logistik secara
				langsung dari satu dashboard komprehensif.
			</p>
		</div>

		<!-- Interactive Tab Selector -->
		<div class="mt-10 flex justify-center">
			<div
				class="inline-flex rounded-xl border border-slate-200 bg-white p-1.5 shadow-xs dark:border-slate-800 dark:bg-slate-950"
			>
				<button
					type="button"
					onclick={() => (activeTab = 'overview')}
					class="cursor-pointer rounded-lg px-4 py-2 text-xs font-semibold transition sm:text-sm {activeTab ===
					'overview'
						? 'bg-indigo-600 text-white shadow-xs dark:bg-indigo-500'
						: 'text-slate-600 hover:text-slate-950 dark:text-slate-400 dark:hover:text-white'}"
				>
					Ringkasan Eksekutif
				</button>
				<button
					type="button"
					onclick={() => (activeTab = 'financials')}
					class="cursor-pointer rounded-lg px-4 py-2 text-xs font-semibold transition sm:text-sm {activeTab ===
					'financials'
						? 'bg-indigo-600 text-white shadow-xs dark:bg-indigo-500'
						: 'text-slate-600 hover:text-slate-950 dark:text-slate-400 dark:hover:text-white'}"
				>
					Arus Kas & GL
				</button>
				<button
					type="button"
					onclick={() => (activeTab = 'orders')}
					class="cursor-pointer rounded-lg px-4 py-2 text-xs font-semibold transition sm:text-sm {activeTab ===
					'orders'
						? 'bg-indigo-600 text-white shadow-xs dark:bg-indigo-500'
						: 'text-slate-600 hover:text-slate-950 dark:text-slate-400 dark:hover:text-white'}"
				>
					Operasional Rantai Pasok
				</button>
			</div>
		</div>

		<!-- Mockup App Frame Window -->
		<div
			class="mt-8 overflow-hidden rounded-2xl border border-slate-200/90 bg-white shadow-2xl shadow-indigo-500/10 dark:border-slate-800 dark:bg-slate-950"
		>
			<!-- Window Top Bar -->
			<div
				class="flex items-center justify-between border-b border-slate-200/80 bg-slate-100/80 px-4 py-3 dark:border-slate-800 dark:bg-slate-900"
			>
				<div class="flex items-center gap-2">
					<span class="inline-block h-3 w-3 rounded-full bg-rose-500/80"></span>
					<span class="inline-block h-3 w-3 rounded-full bg-amber-500/80"></span>
					<span class="inline-block h-3 w-3 rounded-full bg-emerald-500/80"></span>
					<span class="ml-2 font-mono text-xs font-semibold text-slate-500 dark:text-slate-400">
						https://erp.perusahaan.com/portal/dashboard
					</span>
				</div>
				<div class="flex items-center gap-2 text-xs font-medium text-slate-500 dark:text-slate-400">
					<span
						class="inline-flex items-center gap-1 rounded-full bg-emerald-50 px-2 py-0.5 text-[11px] font-bold text-emerald-700 dark:bg-emerald-950/80 dark:text-emerald-400"
					>
						<span class="h-1.5 w-1.5 animate-ping rounded-full bg-emerald-500"></span>
						Live Connected (PgBouncer: 4ms)
					</span>
					<span class="hidden sm:inline">FY 2026 | IDR</span>
				</div>
			</div>

			<!-- Window Content Area -->
			<div class="space-y-6 p-6 sm:p-8">
				<!-- KPI Cards Row -->
				<div class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
					<div
						class="rounded-xl border border-slate-200/80 bg-slate-50/50 p-4 dark:border-slate-800/80 dark:bg-slate-900/40"
					>
						<div
							class="flex items-center justify-between text-xs text-slate-500 dark:text-slate-400"
						>
							<span>Total Pendapatan (YTD)</span>
							<span class="font-bold text-emerald-600 dark:text-emerald-400">+18.4%</span>
						</div>
						<div class="mt-2 text-xl font-black text-slate-900 sm:text-2xl dark:text-white">
							Rp 14.825.000.000
						</div>
						<div class="mt-1 text-[11px] text-slate-400">Target Q3 tercapai 104.2%</div>
					</div>

					<div
						class="rounded-xl border border-slate-200/80 bg-slate-50/50 p-4 dark:border-slate-800/80 dark:bg-slate-900/40"
					>
						<div
							class="flex items-center justify-between text-xs text-slate-500 dark:text-slate-400"
						>
							<span>Margin Laba Bersih</span>
							<span class="font-bold text-emerald-600 dark:text-emerald-400">+3.1%</span>
						</div>
						<div class="mt-2 text-xl font-black text-slate-900 sm:text-2xl dark:text-white">
							26.4 %
						</div>
						<div class="mt-1 text-[11px] text-slate-400">EBITDA: Rp 3.910.000.000</div>
					</div>

					<div
						class="rounded-xl border border-slate-200/80 bg-slate-50/50 p-4 dark:border-slate-800/80 dark:bg-slate-900/40"
					>
						<div
							class="flex items-center justify-between text-xs text-slate-500 dark:text-slate-400"
						>
							<span>Pesanan Aktif (PO/SO)</span>
							<span class="font-bold text-indigo-600 dark:text-indigo-400">98.2% on-time</span>
						</div>
						<div class="mt-2 text-xl font-black text-slate-900 sm:text-2xl dark:text-white">
							1.842 Transaksi
						</div>
						<div class="mt-1 text-[11px] text-slate-400">12 menunggu persetujuan WFL</div>
					</div>

					<div
						class="rounded-xl border border-slate-200/80 bg-slate-50/50 p-4 dark:border-slate-800/80 dark:bg-slate-900/40"
					>
						<div
							class="flex items-center justify-between text-xs text-slate-500 dark:text-slate-400"
						>
							<span>Rasio Likuiditas Kas</span>
							<span class="font-bold text-emerald-600 dark:text-emerald-400">Sehat</span>
						</div>
						<div class="mt-2 text-xl font-black text-slate-900 sm:text-2xl dark:text-white">
							3.8x Rasio Lancar
						</div>
						<div class="mt-1 text-[11px] text-slate-400">Kas & Setara Kas: Rp 4.2M</div>
					</div>
				</div>

				<!-- Chart & Activity Grid -->
				<div class="grid grid-cols-1 gap-6 lg:grid-cols-3">
					<!-- Visual Chart Simulation -->
					<div
						class="rounded-xl border border-slate-200/80 bg-white p-5 lg:col-span-2 dark:border-slate-800/80 dark:bg-slate-900/60"
					>
						<div
							class="flex items-center justify-between border-b border-slate-100 pb-4 dark:border-slate-800"
						>
							<div>
								<h3 class="text-sm font-bold text-slate-900 dark:text-white">
									Tren Pendapatan & Pengeluaran Finansial (2026)
								</h3>
								<p class="text-xs text-slate-500">
									Visualisasi multi-axis terintegrasi Apache ECharts
								</p>
							</div>
							<span
								class="rounded-md bg-slate-100 px-2 py-1 text-[11px] font-semibold text-slate-700 dark:bg-slate-800 dark:text-slate-300"
							>
								Bulanan (Jan - Sep)
							</span>
						</div>

						<!-- Simulated Bar Chart Bars -->
						<div class="mt-6 flex h-48 items-end justify-between gap-2 pt-6 sm:gap-4">
							{#each [{ month: 'Jan', rev: 65, exp: 40 }, { month: 'Feb', rev: 72, exp: 42 }, { month: 'Mar', rev: 80, exp: 48 }, { month: 'Apr', rev: 78, exp: 46 }, { month: 'Mei', rev: 88, exp: 52 }, { month: 'Jun', rev: 95, exp: 55 }, { month: 'Jul', rev: 92, exp: 50 }, { month: 'Agu', rev: 104, exp: 60 }, { month: 'Sep', rev: 115, exp: 62 }] as bar (bar.month)}
								<div class="flex h-full flex-1 flex-col items-center justify-end gap-1.5">
									<div class="flex h-36 w-full items-end justify-center gap-1">
										<div
											style="height: {bar.rev}%"
											class="w-full max-w-[14px] rounded-t-sm bg-indigo-600 transition-all hover:bg-indigo-500"
											title="Pendapatan {bar.month}: {bar.rev}%"
										></div>
										<div
											style="height: {bar.exp}%"
											class="w-full max-w-[14px] rounded-t-sm bg-slate-300 transition-all hover:bg-slate-400 dark:bg-slate-700"
											title="Beban {bar.month}: {bar.exp}%"
										></div>
									</div>
									<span class="text-[10px] font-semibold text-slate-500">{bar.month}</span>
								</div>
							{/each}
						</div>

						<div class="mt-4 flex items-center justify-center gap-6 text-xs">
							<div class="flex items-center gap-2">
								<span class="h-3 w-3 rounded-xs bg-indigo-600"></span>
								<span class="font-medium text-slate-600 dark:text-slate-300"
									>Pendapatan Operasional</span
								>
							</div>
							<div class="flex items-center gap-2">
								<span class="h-3 w-3 rounded-xs bg-slate-300 dark:bg-slate-700"></span>
								<span class="font-medium text-slate-600 dark:text-slate-300"
									>Beban Pokok & Operasional</span
								>
							</div>
						</div>
					</div>

					<!-- Real-time Activity Feed / Approvals -->
					<div
						class="rounded-xl border border-slate-200/80 bg-white p-5 dark:border-slate-800/80 dark:bg-slate-900/60"
					>
						<div
							class="flex items-center justify-between border-b border-slate-100 pb-4 dark:border-slate-800"
						>
							<h3 class="text-sm font-bold text-slate-900 dark:text-white">
								Jurnal & Persetujuan Live
							</h3>
							<span
								class="rounded-full bg-emerald-50 px-2 py-0.5 text-[10px] font-bold text-emerald-700 dark:bg-emerald-950 dark:text-emerald-400"
							>
								WebSocket Aktif
							</span>
						</div>

						<div class="mt-4 space-y-3.5">
							<div class="flex items-start gap-3 text-xs">
								<div
									class="flex h-7 w-7 shrink-0 items-center justify-center rounded-lg bg-emerald-100 font-bold text-emerald-700 dark:bg-emerald-950 dark:text-emerald-400"
								>
									GL
								</div>
								<div class="flex-1">
									<div class="flex items-center justify-between">
										<span class="font-bold text-slate-800 dark:text-slate-200"
											>Posting Jurnal Otomatis</span
										>
										<span class="text-[10px] text-slate-400">12:44</span>
									</div>
									<p class="mt-0.5 text-slate-500">
										Faktur #INV-2026-092 diposting ke Piutang Usaha
									</p>
								</div>
							</div>

							<div class="flex items-start gap-3 text-xs">
								<div
									class="flex h-7 w-7 shrink-0 items-center justify-center rounded-lg bg-indigo-100 font-bold text-indigo-700 dark:bg-indigo-950 dark:text-indigo-400"
								>
									PO
								</div>
								<div class="flex-1">
									<div class="flex items-center justify-between">
										<span class="font-bold text-slate-800 dark:text-slate-200"
											>Three-way Matching Berhasil</span
										>
										<span class="text-[10px] text-slate-400">12:38</span>
									</div>
									<p class="mt-0.5 text-slate-500">
										PO #PO-8841 vs GRN-102 vs Inv Vendor klop 100%
									</p>
								</div>
							</div>

							<div class="flex items-start gap-3 text-xs">
								<div
									class="flex h-7 w-7 shrink-0 items-center justify-center rounded-lg bg-purple-100 font-bold text-purple-700 dark:bg-purple-950 dark:text-purple-400"
								>
									WFL
								</div>
								<div class="flex-1">
									<div class="flex items-center justify-between">
										<span class="font-bold text-slate-800 dark:text-slate-200"
											>Persetujuan Direksi Diperlukan</span
										>
										<span class="text-[10px] font-bold text-amber-500">Pending</span>
									</div>
									<p class="mt-0.5 text-slate-500">
										Pengadaan aset tetap > Rp 150.000.000 (Cabang BSD)
									</p>
								</div>
							</div>

							<div class="flex items-start gap-3 text-xs">
								<div
									class="flex h-7 w-7 shrink-0 items-center justify-center rounded-lg bg-sky-100 font-bold text-sky-700 dark:bg-sky-950 dark:text-sky-400"
								>
									TAX
								</div>
								<div class="flex-1">
									<div class="flex items-center justify-between">
										<span class="font-bold text-slate-800 dark:text-slate-200"
											>Sinkronisasi e-Faktur</span
										>
										<span class="text-[10px] text-slate-400">12:15</span>
									</div>
									<p class="mt-0.5 text-slate-500">
										64 faktur pajak keluaran berhasil terverifikasi DJP
									</p>
								</div>
							</div>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
</section>
