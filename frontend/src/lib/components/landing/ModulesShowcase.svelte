<script lang="ts">
	import { ERP_MODULES } from '$lib/data/modules';
	import type { ERPModule } from '$lib/types/landing';

	interface Props {
		onSelectModule: (module: ERPModule) => void;
	}

	let { onSelectModule }: Props = $props();

	type CategoryType =
		'Semua' | 'Keuangan' | 'Rantai Pasok' | 'Penjualan & CRM' | 'Operasional' | 'SDM' | 'Sistem';

	let activeCategory = $state<CategoryType>('Semua');
	let searchQuery = $state('');

	const categories: { label: CategoryType; icon: string; count: number }[] = [
		{ label: 'Semua', icon: '🌐', count: ERP_MODULES.length },
		{
			label: 'Keuangan',
			icon: '💰',
			count: ERP_MODULES.filter((m) => m.category === 'Keuangan').length
		},
		{
			label: 'Rantai Pasok',
			icon: '📦',
			count: ERP_MODULES.filter((m) => m.category === 'Rantai Pasok').length
		},
		{
			label: 'Penjualan & CRM',
			icon: '📈',
			count: ERP_MODULES.filter((m) => m.category === 'Penjualan & CRM').length
		},
		{
			label: 'Operasional',
			icon: '⚙️',
			count: ERP_MODULES.filter((m) => m.category === 'Operasional').length
		},
		{ label: 'SDM', icon: '👥', count: ERP_MODULES.filter((m) => m.category === 'SDM').length },
		{
			label: 'Sistem',
			icon: '🛡️',
			count: ERP_MODULES.filter((m) => m.category === 'Sistem').length
		}
	];

	const filteredModules = $derived(
		ERP_MODULES.filter((m) => {
			const matchesCategory = activeCategory === 'Semua' || m.category === activeCategory;
			const query = searchQuery.trim().toLowerCase();
			const matchesSearch =
				query === '' ||
				m.name.toLowerCase().includes(query) ||
				m.code.toLowerCase().includes(query) ||
				m.description.toLowerCase().includes(query);
			return matchesCategory && matchesSearch;
		})
	);
</script>

<section id="modules" class="py-20 lg:py-28">
	<div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
		<!-- Section Header -->
		<div class="mx-auto max-w-3xl text-center">
			<div
				class="inline-flex items-center gap-1.5 rounded-full bg-indigo-50 px-3.5 py-1 text-xs font-bold text-indigo-700 ring-1 ring-indigo-500/20 dark:bg-indigo-950/70 dark:text-indigo-300"
			>
				<span>📦</span>
				<span>25 Modul Terintegrasi</span>
			</div>
			<h2
				class="mt-4 text-3xl font-black tracking-tight text-slate-900 sm:text-4xl lg:text-5xl dark:text-white"
			>
				Etalase Lengkap 25 Modul Enterprise
			</h2>
			<p class="mt-4 text-base text-slate-600 sm:text-lg dark:text-slate-300">
				Seluruh proses bisnis dirancang dalam arsitektur Hexagonal modular tanpa circular
				dependencies. Pilih modul sesuai skala bisnis Anda dan kembangkan secara bertahap.
			</p>
		</div>

		<!-- Category Filter Tabs & Search Bar -->
		<div class="mt-12 flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between">
			<!-- Category Tabs -->
			<div class="flex flex-wrap items-center gap-2">
				{#each categories as cat (cat.label)}
					<button
						type="button"
						onclick={() => (activeCategory = cat.label)}
						class="inline-flex cursor-pointer items-center gap-1.5 rounded-xl px-3.5 py-2 text-xs font-semibold transition sm:text-sm {activeCategory ===
						cat.label
							? 'bg-indigo-600 text-white shadow-sm shadow-indigo-600/30 dark:bg-indigo-500'
							: 'border border-slate-200 bg-white text-slate-700 hover:bg-slate-100 hover:text-slate-900 dark:border-slate-800 dark:bg-slate-900 dark:text-slate-300 dark:hover:bg-slate-800 dark:hover:text-white'}"
					>
						<span>{cat.icon}</span>
						<span>{cat.label}</span>
						<span
							class="py-0.2 ml-1 rounded-full px-1.5 text-[10px] {activeCategory === cat.label
								? 'bg-indigo-700/60 text-white'
								: 'bg-slate-100 text-slate-600 dark:bg-slate-800 dark:text-slate-400'}"
						>
							{cat.count}
						</span>
					</button>
				{/each}
			</div>

			<!-- Search Filter -->
			<div class="relative w-full lg:w-72">
				<input
					type="text"
					bind:value={searchQuery}
					placeholder="Cari modul (contoh: POS, GL)..."
					class="w-full rounded-xl border border-slate-200 bg-white py-2 pr-4 pl-9 text-sm text-slate-900 placeholder:text-slate-400 focus:border-indigo-500 focus:ring-2 focus:ring-indigo-500/20 focus:outline-hidden dark:border-slate-800 dark:bg-slate-900 dark:text-white dark:placeholder:text-slate-500"
				/>
				<svg
					class="pointer-events-none absolute top-2.5 left-3 h-4 w-4 text-slate-400"
					fill="none"
					viewBox="0 0 24 24"
					stroke="currentColor"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 1114 0z"
					/>
				</svg>
				{#if searchQuery}
					<button
						type="button"
						onclick={() => (searchQuery = '')}
						class="absolute top-2 right-2 text-xs text-slate-400 hover:text-slate-600 dark:hover:text-slate-200"
					>
						✕
					</button>
				{/if}
			</div>
		</div>

		<!-- Modules Grid -->
		<div class="mt-8 grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
			{#each filteredModules as mod (mod.id)}
				<div
					class="group flex flex-col justify-between rounded-2xl border border-slate-200/80 bg-white p-6 shadow-xs transition-all duration-200 hover:-translate-y-1 hover:border-indigo-300 hover:shadow-lg hover:shadow-indigo-500/10 dark:border-slate-800 dark:bg-slate-900 dark:hover:border-indigo-700/80"
				>
					<div>
						<!-- Card Top: Code badge, Name, Icon -->
						<div class="flex items-start justify-between gap-3">
							<div class="flex items-center gap-3">
								<div
									class="flex h-12 w-12 items-center justify-center rounded-xl bg-slate-100 text-2xl transition group-hover:scale-110 dark:bg-slate-800"
								>
									{mod.icon}
								</div>
								<div>
									<div class="flex items-center gap-2">
										<span
											class="font-mono text-xs font-black tracking-wider text-indigo-600 uppercase dark:text-indigo-400"
										>
											{mod.code}
										</span>
										{#if mod.badge}
											<span
												class="rounded-full bg-indigo-50 px-2 py-0.5 text-[10px] font-semibold text-indigo-700 ring-1 ring-indigo-500/20 dark:bg-indigo-950 dark:text-indigo-300"
											>
												{mod.badge}
											</span>
										{/if}
									</div>
									<h3
										class="mt-1 text-base font-bold text-slate-900 transition-colors group-hover:text-indigo-600 dark:text-white dark:group-hover:text-indigo-400"
									>
										{mod.name}
									</h3>
								</div>
							</div>
						</div>

						<!-- Description -->
						<p class="mt-3 text-xs leading-relaxed text-slate-600 dark:text-slate-400">
							{mod.description}
						</p>

						<!-- Key Features Bullets -->
						<div class="mt-4 space-y-1.5 border-t border-slate-100 pt-4 dark:border-slate-800/80">
							<span class="text-[11px] font-bold tracking-wider text-slate-400 uppercase">
								Fitur Unggulan:
							</span>
							{#each mod.features.slice(0, 3) as feat (feat)}
								<div class="flex items-start gap-2 text-xs text-slate-700 dark:text-slate-300">
									<span class="shrink-0 text-indigo-500">✓</span>
									<span class="line-clamp-1">{feat}</span>
								</div>
							{/each}
						</div>
					</div>

					<!-- Card Action: Detail Trigger -->
					<div
						class="mt-6 flex items-center justify-between border-t border-slate-100 pt-4 dark:border-slate-800/80"
					>
						<span class="text-[11px] font-medium text-slate-400 dark:text-slate-500">
							Kategori: <strong class="text-slate-600 dark:text-slate-300">{mod.category}</strong>
						</span>
						<button
							type="button"
							onclick={() => onSelectModule(mod)}
							class="inline-flex cursor-pointer items-center gap-1 text-xs font-bold text-indigo-600 transition hover:text-indigo-700 dark:text-indigo-400 dark:hover:text-indigo-300"
						>
							<span>Spesifikasi</span>
							<svg class="h-3.5 w-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
								<path
									stroke-linecap="round"
									stroke-linejoin="round"
									stroke-width="2"
									d="M9 5l7 7-7 7"
								/>
							</svg>
						</button>
					</div>
				</div>
			{/each}
		</div>

		{#if filteredModules.length === 0}
			<div
				class="mt-12 rounded-2xl border border-dashed border-slate-300 p-12 text-center dark:border-slate-700"
			>
				<p class="text-base font-semibold text-slate-600 dark:text-slate-400">
					Tidak ada modul yang sesuai dengan kata kunci "{searchQuery}".
				</p>
				<button
					type="button"
					onclick={() => {
						searchQuery = '';
						activeCategory = 'Semua';
					}}
					class="mt-3 text-sm font-bold text-indigo-600 hover:underline dark:text-indigo-400"
				>
					Reset Pencarian
				</button>
			</div>
		{/if}
	</div>
</section>
