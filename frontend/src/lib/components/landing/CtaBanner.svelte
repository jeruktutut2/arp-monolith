<script lang="ts">
	interface Props {
		onOpenDemoModal: () => void;
		onOpenTrialModal: () => void;
	}

	let { onOpenDemoModal, onOpenTrialModal }: Props = $props();
</script>

<section
	class="relative overflow-hidden bg-gradient-to-br from-indigo-700 via-indigo-600 to-blue-700 py-20 text-white lg:py-28"
>
	<!-- Background Patterns -->
	<div
		class="pointer-events-none absolute inset-0 -z-10 bg-[radial-gradient(#ffffff_1px,transparent_1px)] [background-size:24px_24px] opacity-20"
	></div>
	<div
		class="pointer-events-none absolute -bottom-32 -left-32 -z-10 h-96 w-96 rounded-full bg-cyan-400/20 blur-3xl"
	></div>
	<div
		class="pointer-events-none absolute -top-32 -right-32 -z-10 h-96 w-96 rounded-full bg-blue-400/20 blur-3xl"
	></div>

	<div class="mx-auto max-w-7xl px-4 text-center sm:px-6 lg:px-8">
		<span
			class="inline-flex items-center gap-1.5 rounded-full bg-white/10 px-4 py-1.5 text-xs font-bold tracking-wider text-indigo-100 uppercase ring-1 ring-white/20 backdrop-blur-xs"
		>
			🚀 Siap Memulai Transformasi?
		</span>

		<h2 class="mt-6 text-3xl font-black tracking-tight text-white sm:text-5xl lg:text-6xl">
			Transformasikan Operasional Bisnis Anda Hari Ini
		</h2>

		<p class="mx-auto mt-6 max-w-2xl text-base leading-relaxed text-indigo-100/90 sm:text-xl">
			Tinggalkan sistem yang terpecah-pecah. Bergabunglah dengan ratusan korporasi modern yang
			menggunakan sistem ERP Monolith untuk efisiensi tak tertandingi.
		</p>

		<div class="mt-10 flex flex-col items-center justify-center gap-4 sm:flex-row">
			<button
				type="button"
				onclick={onOpenTrialModal}
				class="inline-flex w-full cursor-pointer items-center justify-center gap-2 rounded-xl bg-white px-8 py-4 text-base font-bold text-indigo-700 shadow-xl shadow-black/10 transition hover:bg-slate-50 active:scale-95 sm:w-auto"
			>
				<span>Mulai Uji Coba 30 Hari</span>
				<svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
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
				class="inline-flex w-full cursor-pointer items-center justify-center gap-2 rounded-xl border border-white/30 bg-white/10 px-8 py-4 text-base font-bold text-white backdrop-blur-xs transition hover:bg-white/20 active:scale-95 sm:w-auto"
			>
				<svg class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
					<path
						fill-rule="evenodd"
						d="M6 2a1 1 0 00-1 1v1H4a2 2 0 00-2 2v10a2 2 0 002 2h12a2 2 0 002-2V6a2 2 0 00-2-2h-1V3a1 1 0 10-2 0v1H7V3a1 1 0 00-1-1zm0 5a1 1 0 000 2h8a1 1 0 100-2H6z"
						clip-rule="evenodd"
					/>
				</svg>
				<span>Konsultasi Solusi Enterprise</span>
			</button>
		</div>

		<div class="mt-8 flex flex-wrap items-center justify-center gap-6 text-xs text-indigo-200">
			<span class="flex items-center gap-1.5">
				<span>✓</span> Tanpa kartu kredit
			</span>
			<span class="flex items-center gap-1.5">
				<span>✓</span> Onboarding & migrasi data didampingi
			</span>
			<span class="flex items-center gap-1.5">
				<span>✓</span> Dukungan teknis lokal 24/7
			</span>
		</div>
	</div>
</section>
