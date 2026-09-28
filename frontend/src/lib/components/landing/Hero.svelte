<script lang="ts">
	interface Props {
		onOpenDemoModal: () => void;
		onOpenTrialModal: () => void;
	}

	let { onOpenDemoModal, onOpenTrialModal }: Props = $props();
</script>

<section class="relative overflow-hidden pt-12 pb-20 lg:pt-20 lg:pb-28">
	<!-- Background Glow Effects -->
	<div
		class="pointer-events-none absolute -top-40 left-1/2 -z-10 h-[550px] w-[900px] -translate-x-1/2 transform rounded-full bg-gradient-to-tr from-indigo-500/20 via-blue-500/20 to-sky-400/10 blur-[130px] dark:from-indigo-900/30 dark:via-blue-900/20 dark:to-cyan-900/15"
	></div>
	<div
		class="pointer-events-none absolute top-1/2 -right-60 -z-10 h-[400px] w-[500px] rounded-full bg-indigo-500/10 blur-[100px] dark:bg-indigo-900/20"
	></div>

	<div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
		<div class="mx-auto max-w-4xl text-center">
			<!-- Enterprise Badge Pill -->
			<div
				class="inline-flex items-center gap-2 rounded-full border border-indigo-200/80 bg-indigo-50/70 px-4 py-1.5 text-xs font-semibold text-indigo-800 backdrop-blur-xs transition hover:border-indigo-300 dark:border-indigo-900/60 dark:bg-indigo-950/60 dark:text-indigo-300"
			>
				<span class="flex h-2 w-2 animate-pulse rounded-full bg-indigo-600"></span>
				<span class="font-bold tracking-wide uppercase">Enterprise ERP Monolith</span>
				<span class="text-indigo-400">|</span>
				<span>Hexagonal Modular Architecture</span>
			</div>

			<!-- Main Headline -->
			<h1
				class="mt-8 text-4xl font-black tracking-tight text-slate-900 sm:text-6xl sm:leading-[1.12] dark:text-white"
			>
				Next-Gen Enterprise <br class="hidden sm:inline" />
				<span
					class="bg-gradient-to-r from-indigo-600 via-blue-600 to-cyan-500 bg-clip-text text-transparent dark:from-indigo-400 dark:via-blue-400 dark:to-cyan-300"
				>
					Resource Planning
				</span>
				for Modern Teams
			</h1>

			<!-- Subtitle -->
			<p class="mx-auto mt-6 max-w-2xl text-lg text-slate-600 sm:text-xl dark:text-slate-300">
				Platform ERP generasi baru yang menggabungkan kecepatan backend Golang Echo v5 dengan
				arsitektur <strong>Modular Monolith Hexagonal</strong>, keandalan transaksi PgBouncer, dan
				antarmuka reaktif SvelteKit 2. Satukan 25 modul bisnis dalam satu kendali terpusat.
			</p>

			<!-- Dual CTA Buttons -->
			<div class="mt-10 flex flex-col items-center justify-center gap-4 sm:flex-row">
				<button
					type="button"
					onclick={onOpenTrialModal}
					class="group relative inline-flex w-full items-center justify-center gap-2.5 rounded-xl bg-gradient-to-r from-indigo-600 via-blue-600 to-indigo-700 px-7 py-3.5 text-base font-bold text-white shadow-lg shadow-indigo-600/30 transition hover:from-indigo-500 hover:to-blue-500 active:scale-95 sm:w-auto"
				>
					<span>Mulai Uji Coba Gratis</span>
					<svg
						class="h-5 w-5 transition-transform group-hover:translate-x-1"
						fill="none"
						viewBox="0 0 24 24"
						stroke="currentColor"
					>
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M14 5l7 7m0 0l-7 7m7-7H3"
						/>
					</svg>
				</button>

				<button
					type="button"
					onclick={onOpenDemoModal}
					class="inline-flex w-full items-center justify-center gap-2.5 rounded-xl border border-slate-300 bg-white/80 px-7 py-3.5 text-base font-bold text-slate-700 backdrop-blur-xs transition hover:bg-slate-100 hover:text-slate-950 sm:w-auto dark:border-slate-700 dark:bg-slate-900/80 dark:text-slate-200 dark:hover:border-slate-600 dark:hover:bg-slate-800"
				>
					<svg
						class="h-5 w-5 text-indigo-600 dark:text-indigo-400"
						viewBox="0 0 20 20"
						fill="currentColor"
					>
						<path
							fill-rule="evenodd"
							d="M10 18a8 8 0 100-16 8 8 0 000 16zM9.555 7.168A1 1 0 008 8v4a1 1 0 001.555.832l3-2a1 1 0 000-1.664l-3-2z"
							clip-rule="evenodd"
						/>
					</svg>
					<span>Jadwalkan Demo Langsung</span>
				</button>
			</div>

			<!-- Tech Stack Badges -->
			<div
				class="mt-12 flex flex-wrap items-center justify-center gap-3 border-t border-slate-200/60 pt-6 dark:border-slate-800/60"
			>
				<span
					class="text-xs font-semibold tracking-wider text-slate-400 uppercase dark:text-slate-500"
				>
					Didukung Teknologi Unggulan:
				</span>
				<div
					class="flex flex-wrap items-center justify-center gap-2 text-xs font-medium text-slate-600 dark:text-slate-300"
				>
					<span class="rounded-lg bg-slate-100 px-2.5 py-1 dark:bg-slate-800/80"
						>Golang 1.22 + Echo v5</span
					>
					<span class="rounded-lg bg-slate-100 px-2.5 py-1 dark:bg-slate-800/80"
						>PostgreSQL 16 + PgBouncer</span
					>
					<span class="rounded-lg bg-slate-100 px-2.5 py-1 dark:bg-slate-800/80"
						>SvelteKit 2 + Svelte 5 Runes</span
					>
					<span class="rounded-lg bg-slate-100 px-2.5 py-1 dark:bg-slate-800/80">Bun Runtime</span>
					<span
						class="rounded-lg bg-emerald-50 px-2.5 py-1 text-emerald-700 dark:bg-emerald-950/60 dark:text-emerald-400"
					>
						PSAK / IFRS Compliant
					</span>
				</div>
			</div>
		</div>
	</div>
</section>
