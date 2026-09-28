<script lang="ts">
	import { SOLUTIONS_DATA } from '$lib/data/modules';

	let activeSolutionIndex = $state(0);
	const activeSolution = $derived(SOLUTIONS_DATA[activeSolutionIndex]);
</script>

<section id="solutions" class="py-20 lg:py-28">
	<div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
		<!-- Section Header -->
		<div class="mx-auto max-w-3xl text-center">
			<div
				class="inline-flex items-center gap-1.5 rounded-full bg-cyan-50 px-3.5 py-1 text-xs font-bold text-cyan-700 ring-1 ring-cyan-500/20 dark:bg-cyan-950/70 dark:text-cyan-300"
			>
				<span>🏢</span>
				<span>Use Cases Spesifik Industri</span>
			</div>
			<h2
				class="mt-4 text-3xl font-black tracking-tight text-slate-900 sm:text-4xl lg:text-5xl dark:text-white"
			>
				Solusi Khusus Sesuai Karakter Industri Anda
			</h2>
			<p class="mt-4 text-base text-slate-600 sm:text-lg dark:text-slate-300">
				Setiap industri memiliki proses bisnis unik. Sistem ERP Monolith dirancang fleksibel untuk
				menjawab kompleksitas operasional holding, manufaktur, distribusi, ritel, hingga proyek
				kontraktor.
			</p>
		</div>

		<!-- Interactive Industry Showcase -->
		<div class="mt-14 grid grid-cols-1 gap-8 lg:grid-cols-12 lg:items-center">
			<!-- Industry Navigation Buttons (Left) -->
			<div class="space-y-3 lg:col-span-5">
				{#each SOLUTIONS_DATA as sol, idx (sol.title)}
					<button
						type="button"
						onclick={() => (activeSolutionIndex = idx)}
						class="w-full cursor-pointer rounded-2xl border p-5 text-left transition {activeSolutionIndex ===
						idx
							? 'border-indigo-600 bg-indigo-50/60 shadow-md shadow-indigo-600/10 dark:border-indigo-500 dark:bg-indigo-950/40'
							: 'border-slate-200/80 bg-white hover:border-slate-300 hover:bg-slate-50 dark:border-slate-800 dark:bg-slate-900 dark:hover:border-slate-700 dark:hover:bg-slate-800/60'}"
					>
						<div class="flex items-center gap-3.5">
							<span class="text-2xl">{sol.icon}</span>
							<div>
								<h3
									class="text-base font-bold text-slate-900 transition-colors {activeSolutionIndex ===
									idx
										? 'text-indigo-600 dark:text-indigo-400'
										: 'dark:text-white'}"
								>
									{sol.title}
								</h3>
								<p class="mt-0.5 line-clamp-1 text-xs text-slate-500 dark:text-slate-400">
									{sol.tagline}
								</p>
							</div>
						</div>
					</button>
				{/each}
			</div>

			<!-- Active Industry Card Details (Right) -->
			<div class="lg:col-span-7">
				<div
					class="rounded-3xl border border-slate-200/90 bg-gradient-to-br from-white via-slate-50/50 to-indigo-50/30 p-8 shadow-xl dark:border-slate-800 dark:from-slate-900 dark:via-slate-900/80 dark:to-indigo-950/30"
				>
					<div class="flex items-center gap-3">
						<span class="text-4xl">{activeSolution.icon}</span>
						<div>
							<span
								class="text-xs font-bold tracking-wider text-indigo-600 uppercase dark:text-indigo-400"
							>
								{activeSolution.tagline}
							</span>
							<h3 class="text-2xl font-black text-slate-900 dark:text-white">
								{activeSolution.title}
							</h3>
						</div>
					</div>

					<p class="mt-6 text-sm leading-relaxed text-slate-600 dark:text-slate-300">
						{activeSolution.description}
					</p>

					<!-- Key Highlights -->
					<div class="mt-6 space-y-3">
						<h4 class="text-xs font-bold tracking-wider text-slate-400 uppercase">
							Nilai Tambah Operasional:
						</h4>
						<div class="grid grid-cols-1 gap-2.5 sm:grid-cols-2">
							{#each activeSolution.highlights as hl (hl)}
								<div
									class="flex items-start gap-2 text-xs font-medium text-slate-700 dark:text-slate-300"
								>
									<div
										class="flex h-4 w-4 shrink-0 items-center justify-center rounded-full bg-emerald-100 font-bold text-emerald-700 dark:bg-emerald-950 dark:text-emerald-400"
									>
										✓
									</div>
									<span>{hl}</span>
								</div>
							{/each}
						</div>
					</div>

					<!-- Relevant Modules Pill List -->
					<div class="mt-8 border-t border-slate-200/80 pt-6 dark:border-slate-800">
						<span class="text-xs font-semibold text-slate-500 dark:text-slate-400">
							Modul Kunci Terintegrasi:
						</span>
						<div class="mt-3 flex flex-wrap gap-2">
							{#each activeSolution.modules as modCode (modCode)}
								<span
									class="inline-flex items-center gap-1 rounded-xl border border-slate-200 bg-white px-3 py-1.5 font-mono text-xs font-bold text-indigo-700 shadow-xs dark:border-slate-700 dark:bg-slate-800 dark:text-indigo-300"
								>
									<span>⚡</span>
									<span>Modul {modCode}</span>
								</span>
							{/each}
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
</section>
