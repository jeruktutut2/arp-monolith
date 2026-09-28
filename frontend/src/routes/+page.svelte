<script lang="ts">
	let count = $state(0);
	const doubleCount = $derived(count * 2);

	function increment() {
		count++;
	}

	function reset() {
		count = 0;
	}
</script>

<div
	class="min-h-screen bg-slate-50 text-slate-800 antialiased dark:bg-slate-900 dark:text-slate-100"
>
	<header
		class="border-b border-slate-200 bg-white shadow-xs dark:border-slate-800 dark:bg-slate-800/80"
	>
		<div class="mx-auto flex max-w-7xl items-center justify-between px-6 py-4">
			<div class="flex items-center gap-3">
				<div
					class="flex h-10 w-10 items-center justify-center rounded-lg bg-indigo-600 font-bold text-white shadow-sm"
				>
					ERP
				</div>
				<div>
					<h1 class="text-lg leading-tight font-bold tracking-tight">Enterprise ERP Monolith</h1>
					<p class="text-xs text-slate-500 dark:text-slate-400">Modular Monolith Frontend Client</p>
				</div>
			</div>
			<div class="flex items-center gap-2">
				<span
					class="inline-flex items-center rounded-full bg-emerald-50 px-2.5 py-1 text-xs font-semibold text-emerald-700 ring-1 ring-emerald-600/20 ring-inset dark:bg-emerald-950/60 dark:text-emerald-400"
				>
					● SvelteKit 2 + Svelte 5
				</span>
				<span
					class="inline-flex items-center rounded-full bg-amber-50 px-2.5 py-1 text-xs font-semibold text-amber-700 ring-1 ring-amber-600/20 ring-inset dark:bg-amber-950/60 dark:text-amber-400"
				>
					Bun Runtime
				</span>
			</div>
		</div>
	</header>

	<main class="mx-auto max-w-7xl px-6 py-12">
		<div class="mb-10 text-center">
			<h2 class="text-3xl font-extrabold tracking-tight text-slate-900 sm:text-4xl dark:text-white">
				Portal Antarmuka Enterprise ERP
			</h2>
			<p class="mx-auto mt-3 max-w-2xl text-base text-slate-600 dark:text-slate-400">
				Frontend modern dibangun dengan SvelteKit 2, Svelte 5 Runes, Tailwind CSS, dan dikelola
				menggunakan Bun untuk performa maksimal.
			</p>
		</div>

		<div class="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
			<div
				class="rounded-xl border border-slate-200 bg-white p-6 shadow-xs dark:border-slate-800 dark:bg-slate-800"
			>
				<div class="flex items-center gap-3">
					<div
						class="flex h-9 w-9 items-center justify-center rounded-lg bg-indigo-100 text-indigo-700 dark:bg-indigo-900/60 dark:text-indigo-300"
					>
						⚡
					</div>
					<h3 class="text-base font-semibold">Svelte 5 Runes</h3>
				</div>
				<p class="mt-2 text-sm text-slate-600 dark:text-slate-400">
					Reaktivitas mutakhir menggunakan rune resmi (<code
						class="rounded bg-slate-100 px-1.5 py-0.5 text-xs text-indigo-600 dark:bg-slate-700 dark:text-indigo-400"
						>$state</code
					>,
					<code
						class="rounded bg-slate-100 px-1.5 py-0.5 text-xs text-indigo-600 dark:bg-slate-700 dark:text-indigo-400"
						>$derived</code
					>,
					<code
						class="rounded bg-slate-100 px-1.5 py-0.5 text-xs text-indigo-600 dark:bg-slate-700 dark:text-indigo-400"
						>$props</code
					>).
				</p>
				<div class="mt-4 rounded-lg bg-slate-50 p-4 dark:bg-slate-900/50">
					<div class="flex items-center justify-between text-sm">
						<span class="text-slate-500">Nilai Count:</span>
						<span class="font-mono font-bold text-indigo-600 dark:text-indigo-400">{count}</span>
					</div>
					<div class="mt-1 flex items-center justify-between text-sm">
						<span class="text-slate-500">Nilai Derived (2x):</span>
						<span class="font-mono font-bold text-slate-700 dark:text-slate-300">{doubleCount}</span
						>
					</div>
					<div class="mt-4 flex gap-2">
						<button
							onclick={increment}
							class="flex-1 cursor-pointer rounded-md bg-indigo-600 px-3 py-1.5 text-xs font-semibold text-white transition hover:bg-indigo-500"
						>
							+ Tambah
						</button>
						<button
							onclick={reset}
							class="cursor-pointer rounded-md border border-slate-300 px-3 py-1.5 text-xs font-semibold text-slate-700 transition hover:bg-slate-100 dark:border-slate-700 dark:text-slate-300 dark:hover:bg-slate-800"
						>
							Reset
						</button>
					</div>
				</div>
			</div>

			<div
				class="rounded-xl border border-slate-200 bg-white p-6 shadow-xs dark:border-slate-800 dark:bg-slate-800"
			>
				<div class="flex items-center gap-3">
					<div
						class="flex h-9 w-9 items-center justify-center rounded-lg bg-emerald-100 text-emerald-700 dark:bg-emerald-900/60 dark:text-emerald-300"
					>
						🥟
					</div>
					<h3 class="text-base font-semibold">Bun Toolchain</h3>
				</div>
				<p class="mt-2 text-sm text-slate-600 dark:text-slate-400">
					Package manager dan runtime berkecepatan tinggi dengan instalasi instan dan eksekusi skrip
					kilat.
				</p>
				<div class="mt-4 space-y-2 font-mono text-xs text-slate-600 dark:text-slate-400">
					<div class="rounded bg-slate-50 p-2 dark:bg-slate-900/50">bun run dev</div>
					<div class="rounded bg-slate-50 p-2 dark:bg-slate-900/50">bun run build</div>
					<div class="rounded bg-slate-50 p-2 dark:bg-slate-900/50">bun run check</div>
				</div>
			</div>

			<div
				class="rounded-xl border border-slate-200 bg-white p-6 shadow-xs dark:border-slate-800 dark:bg-slate-800"
			>
				<div class="flex items-center gap-3">
					<div
						class="flex h-9 w-9 items-center justify-center rounded-lg bg-sky-100 text-sky-700 dark:bg-sky-900/60 dark:text-sky-300"
					>
						🖥️
					</div>
					<h3 class="text-base font-semibold">Backend Integration</h3>
				</div>
				<p class="mt-2 text-sm text-slate-600 dark:text-slate-400">
					Koneksi ke backend Golang Echo v5 modular monolith dan PgBouncer transaction pooling.
				</p>
				<ul class="mt-4 space-y-1.5 text-xs text-slate-600 dark:text-slate-400">
					<li class="flex items-center gap-1.5">
						<span class="text-emerald-500">✔</span> Echo v5 High-Performance Router
					</li>
					<li class="flex items-center gap-1.5">
						<span class="text-emerald-500">✔</span> Adapter Node dengan Precompression
					</li>
					<li class="flex items-center gap-1.5">
						<span class="text-emerald-500">✔</span> Systemd Service Ready (:3000)
					</li>
				</ul>
			</div>
		</div>
	</main>
</div>
