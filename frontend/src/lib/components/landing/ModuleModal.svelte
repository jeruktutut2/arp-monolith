<script lang="ts">
	import type { ERPModule } from '$lib/types/landing';

	interface Props {
		module: ERPModule | null;
		onClose: () => void;
	}

	let { module, onClose }: Props = $props();
</script>

{#if module}
	<div
		class="animate-in fade-in fixed inset-0 z-50 flex items-center justify-center bg-slate-950/70 p-4 backdrop-blur-xs duration-200"
		role="dialog"
		aria-modal="true"
	>
		<div
			class="relative max-h-[90vh] w-full max-w-2xl overflow-y-auto rounded-3xl border border-slate-200 bg-white p-6 shadow-2xl sm:p-8 dark:border-slate-800 dark:bg-slate-900"
		>
			<!-- Close Button -->
			<button
				type="button"
				onclick={onClose}
				aria-label="Tutup spesifikasi modul"
				class="absolute top-5 right-5 flex h-8 w-8 cursor-pointer items-center justify-center rounded-full text-slate-400 hover:bg-slate-100 hover:text-slate-600 dark:hover:bg-slate-800 dark:hover:text-white"
			>
				✕
			</button>

			<!-- Header -->
			<div class="flex items-start gap-4">
				<div
					class="flex h-16 w-16 shrink-0 items-center justify-center rounded-2xl bg-indigo-50 text-3xl dark:bg-slate-800"
				>
					{module.icon}
				</div>
				<div>
					<div class="flex flex-wrap items-center gap-2">
						<span
							class="font-mono text-sm font-black tracking-wider text-indigo-600 uppercase dark:text-indigo-400"
						>
							Modul {module.code}
						</span>
						<span
							class="rounded-full bg-slate-100 px-2.5 py-0.5 text-xs font-semibold text-slate-700 dark:bg-slate-800 dark:text-slate-300"
						>
							Pilar {module.category}
						</span>
						{#if module.badge}
							<span
								class="rounded-full bg-emerald-50 px-2.5 py-0.5 text-xs font-semibold text-emerald-700 ring-1 ring-emerald-500/20 dark:bg-emerald-950/80 dark:text-emerald-300"
							>
								{module.badge}
							</span>
						{/if}
					</div>
					<h3 class="mt-1 text-2xl font-black text-slate-900 dark:text-white">
						{module.name}
					</h3>
				</div>
			</div>

			<!-- Description -->
			<p class="mt-5 text-sm leading-relaxed text-slate-600 dark:text-slate-300">
				{module.description}
			</p>

			<!-- All Features -->
			<div class="mt-6 border-t border-slate-100 pt-5 dark:border-slate-800">
				<h4 class="text-xs font-bold tracking-wider text-slate-400 uppercase">
					Daftar Kapabilitas & Fitur Lengkap:
				</h4>
				<ul class="mt-3 grid grid-cols-1 gap-2.5 sm:grid-cols-2">
					{#each module.features as feat (feat)}
						<li
							class="flex items-start gap-2 text-xs font-medium text-slate-700 dark:text-slate-300"
						>
							<span class="shrink-0 font-bold text-emerald-500">✓</span>
							<span>{feat}</span>
						</li>
					{/each}
				</ul>
			</div>

			<!-- Relational Flow -->
			<div
				class="mt-6 grid grid-cols-1 gap-4 border-t border-slate-100 pt-5 sm:grid-cols-2 dark:border-slate-800"
			>
				<div class="rounded-xl bg-slate-50 p-3.5 dark:bg-slate-800/60">
					<span
						class="text-[11px] font-bold tracking-wider text-slate-500 uppercase dark:text-slate-400"
					>
						Menerima Input Dari Modul:
					</span>
					<div class="mt-2 flex flex-wrap gap-1.5">
						{#if module.inputsFrom && module.inputsFrom.length > 0}
							{#each module.inputsFrom as inputMod (inputMod)}
								<span
									class="rounded-md border border-slate-200 bg-white px-2 py-1 font-mono text-xs font-bold text-slate-700 shadow-xs dark:border-slate-600 dark:bg-slate-700 dark:text-slate-200"
								>
									{inputMod}
								</span>
							{/each}
						{:else}
							<span class="text-xs text-slate-400 italic">Data Master / Independen</span>
						{/if}
					</div>
				</div>

				<div class="rounded-xl bg-slate-50 p-3.5 dark:bg-slate-800/60">
					<span
						class="text-[11px] font-bold tracking-wider text-slate-500 uppercase dark:text-slate-400"
					>
						Meneruskan Output / Posting Ke:
					</span>
					<div class="mt-2 flex flex-wrap gap-1.5">
						{#if module.outputsTo && module.outputsTo.length > 0}
							{#each module.outputsTo as outputMod (outputMod)}
								<span
									class="rounded-md border border-slate-200 bg-white px-2 py-1 font-mono text-xs font-bold text-indigo-700 shadow-xs dark:border-slate-600 dark:bg-slate-700 dark:text-indigo-300"
								>
									{outputMod}
								</span>
							{/each}
						{:else}
							<span class="text-xs text-slate-400 italic">Modul Konsumen Akhir / Laporan</span>
						{/if}
					</div>
				</div>
			</div>

			<!-- Footer Action -->
			<div class="mt-8 flex justify-end">
				<button
					type="button"
					onclick={onClose}
					class="cursor-pointer rounded-xl bg-slate-900 px-5 py-2.5 text-xs font-bold text-white transition hover:bg-slate-800 dark:bg-slate-700 dark:hover:bg-slate-600"
				>
					Tutup Spesifikasi
				</button>
			</div>
		</div>
	</div>
{/if}
