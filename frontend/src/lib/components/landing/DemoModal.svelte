<script lang="ts">
	interface Props {
		isOpen: boolean;
		title?: string;
		onClose: () => void;
	}

	let { isOpen, title = 'Jadwalkan Konsultasi & Demo Produk', onClose }: Props = $props();

	let fullName = $state('');
	let email = $state('');
	let company = $state('');
	let employees = $state('50-250');
	let notes = $state('');
	let submitted = $state(false);

	function handleSubmit(e: Event) {
		e.preventDefault();
		submitted = true;
	}

	function handleReset() {
		submitted = false;
		fullName = '';
		email = '';
		company = '';
		notes = '';
		onClose();
	}
</script>

{#if isOpen}
	<div
		class="animate-in fade-in fixed inset-0 z-50 flex items-center justify-center bg-slate-950/70 p-4 backdrop-blur-xs duration-200"
		role="dialog"
		aria-modal="true"
	>
		<div
			class="relative w-full max-w-lg rounded-3xl border border-slate-200 bg-white p-6 shadow-2xl sm:p-8 dark:border-slate-800 dark:bg-slate-900"
		>
			<!-- Close Button -->
			<button
				type="button"
				onclick={onClose}
				aria-label="Tutup modal"
				class="absolute top-5 right-5 flex h-8 w-8 cursor-pointer items-center justify-center rounded-full text-slate-400 hover:bg-slate-100 hover:text-slate-600 dark:hover:bg-slate-800 dark:hover:text-white"
			>
				✕
			</button>

			{#if !submitted}
				<div class="flex items-center gap-3">
					<div
						class="flex h-10 w-10 items-center justify-center rounded-xl bg-indigo-100 text-xl text-indigo-600 dark:bg-indigo-950 dark:text-indigo-400"
					>
						📅
					</div>
					<div>
						<h3 class="text-xl font-bold text-slate-900 dark:text-white">
							{title}
						</h3>
						<p class="text-xs text-slate-500 dark:text-slate-400">
							Tim spesialis ERP kami akan menghubungi Anda dalam 1x24 jam kerja.
						</p>
					</div>
				</div>

				<form onsubmit={handleSubmit} class="mt-6 space-y-4">
					<div>
						<label
							for="demo-name"
							class="block text-xs font-bold text-slate-700 dark:text-slate-300"
						>
							Nama Lengkap *
						</label>
						<input
							id="demo-name"
							type="text"
							required
							bind:value={fullName}
							placeholder="contoh: Budi Santoso"
							class="mt-1 w-full rounded-xl border border-slate-300 bg-white px-3.5 py-2 text-sm text-slate-900 placeholder:text-slate-400 focus:border-indigo-600 focus:ring-2 focus:ring-indigo-600/20 focus:outline-hidden dark:border-slate-700 dark:bg-slate-800 dark:text-white"
						/>
					</div>

					<div>
						<label
							for="demo-email"
							class="block text-xs font-bold text-slate-700 dark:text-slate-300"
						>
							Email Perusahaan *
						</label>
						<input
							id="demo-email"
							type="email"
							required
							bind:value={email}
							placeholder="nama@perusahaan.co.id"
							class="mt-1 w-full rounded-xl border border-slate-300 bg-white px-3.5 py-2 text-sm text-slate-900 placeholder:text-slate-400 focus:border-indigo-600 focus:ring-2 focus:ring-indigo-600/20 focus:outline-hidden dark:border-slate-700 dark:bg-slate-800 dark:text-white"
						/>
					</div>

					<div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
						<div>
							<label
								for="demo-company"
								class="block text-xs font-bold text-slate-700 dark:text-slate-300"
							>
								Nama Perusahaan *
							</label>
							<input
								id="demo-company"
								type="text"
								required
								bind:value={company}
								placeholder="PT Rekayasa Gemilang"
								class="mt-1 w-full rounded-xl border border-slate-300 bg-white px-3.5 py-2 text-sm text-slate-900 placeholder:text-slate-400 focus:border-indigo-600 focus:ring-2 focus:ring-indigo-600/20 focus:outline-hidden dark:border-slate-700 dark:bg-slate-800 dark:text-white"
							/>
						</div>

						<div>
							<label
								for="demo-employees"
								class="block text-xs font-bold text-slate-700 dark:text-slate-300"
							>
								Jumlah Karyawan
							</label>
							<select
								id="demo-employees"
								bind:value={employees}
								class="mt-1 w-full rounded-xl border border-slate-300 bg-white px-3.5 py-2 text-sm text-slate-900 focus:border-indigo-600 focus:ring-2 focus:ring-indigo-600/20 focus:outline-hidden dark:border-slate-700 dark:bg-slate-800 dark:text-white"
							>
								<option value="1-50">1 - 50 Karyawan</option>
								<option value="50-250">50 - 250 Karyawan</option>
								<option value="250-1000">250 - 1.000 Karyawan</option>
								<option value="1000+">&gt; 1.000 Karyawan</option>
							</select>
						</div>
					</div>

					<div>
						<label
							for="demo-notes"
							class="block text-xs font-bold text-slate-700 dark:text-slate-300"
						>
							Modul yang Paling Dibutuhkan / Catatan Khusus
						</label>
						<textarea
							id="demo-notes"
							rows="2"
							bind:value={notes}
							placeholder="Contoh: Akuntansi multi-currency, Integrasi POS kasir & stok gudang"
							class="mt-1 w-full rounded-xl border border-slate-300 bg-white px-3.5 py-2 text-sm text-slate-900 placeholder:text-slate-400 focus:border-indigo-600 focus:ring-2 focus:ring-indigo-600/20 focus:outline-hidden dark:border-slate-700 dark:bg-slate-800 dark:text-white"
						></textarea>
					</div>

					<div class="pt-2">
						<button
							type="submit"
							class="w-full cursor-pointer rounded-xl bg-gradient-to-r from-indigo-600 to-blue-600 py-3 text-sm font-bold text-white shadow-md shadow-indigo-600/30 transition hover:from-indigo-500 hover:to-blue-500 active:scale-98"
						>
							Kirim Permintaan Demo
						</button>
					</div>
				</form>
			{:else}
				<div class="py-6 text-center">
					<div
						class="mx-auto flex h-16 w-16 items-center justify-center rounded-full bg-emerald-100 text-3xl text-emerald-600 dark:bg-emerald-950 dark:text-emerald-400"
					>
						✓
					</div>
					<h3 class="mt-4 text-2xl font-bold text-slate-900 dark:text-white">
						Terima Kasih, {fullName}!
					</h3>
					<p class="mt-2 text-sm text-slate-600 dark:text-slate-300">
						Permintaan demo Anda untuk <strong>{company}</strong> telah kami terima. Spesialis ERP
						kami akan mengirimkan tautan kalender konsultasi ke <strong>{email}</strong>.
					</p>

					<div class="mt-8">
						<button
							type="button"
							onclick={handleReset}
							class="rounded-xl bg-slate-900 px-6 py-2.5 text-sm font-bold text-white transition hover:bg-slate-800 dark:bg-slate-700 dark:hover:bg-slate-600"
						>
							Tutup
						</button>
					</div>
				</div>
			{/if}
		</div>
	</div>
{/if}
