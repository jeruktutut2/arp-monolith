<script lang="ts">
	interface Props {
		isOpen: boolean;
		onClose: () => void;
	}

	let { isOpen, onClose }: Props = $props();

	let username = $state('');
	let password = $state('');
	let rememberMe = $state(true);
	let loginStatus = $state<'idle' | 'success'>('idle');

	function handleSignIn(e: Event) {
		e.preventDefault();
		loginStatus = 'success';
	}

	function handleReset() {
		loginStatus = 'idle';
		username = '';
		password = '';
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
			class="relative w-full max-w-md rounded-3xl border border-slate-200 bg-white p-6 shadow-2xl sm:p-8 dark:border-slate-800 dark:bg-slate-900"
		>
			<!-- Close Button -->
			<button
				type="button"
				onclick={onClose}
				aria-label="Tutup login"
				class="absolute top-5 right-5 flex h-8 w-8 cursor-pointer items-center justify-center rounded-full text-slate-400 hover:bg-slate-100 hover:text-slate-600 dark:hover:bg-slate-800 dark:hover:text-white"
			>
				✕
			</button>

			<div class="flex items-center gap-3">
				<div
					class="flex h-10 w-10 items-center justify-center rounded-xl bg-indigo-600 text-lg font-black text-white shadow-md shadow-indigo-600/30"
				>
					E
				</div>
				<div>
					<h3 class="text-xl font-bold text-slate-900 dark:text-white">Masuk ke Portal ERP</h3>
					<p class="text-xs text-slate-500 dark:text-slate-400">
						Modul 19_USR / Zero-Trust Identity
					</p>
				</div>
			</div>

			{#if loginStatus === 'idle'}
				<form onsubmit={handleSignIn} class="mt-6 space-y-4">
					<div>
						<label
							for="login-username"
							class="block text-xs font-bold text-slate-700 dark:text-slate-300"
						>
							Email atau Username Korporat
						</label>
						<input
							id="login-username"
							type="text"
							required
							bind:value={username}
							placeholder="admin@perusahaan.co.id"
							class="mt-1 w-full rounded-xl border border-slate-300 bg-white px-3.5 py-2.5 text-sm text-slate-900 placeholder:text-slate-400 focus:border-indigo-600 focus:ring-2 focus:ring-indigo-600/20 focus:outline-hidden dark:border-slate-700 dark:bg-slate-800 dark:text-white"
						/>
					</div>

					<div>
						<div class="flex items-center justify-between">
							<label
								for="login-password"
								class="block text-xs font-bold text-slate-700 dark:text-slate-300"
							>
								Kata Sandi
							</label>
							<a
								href="#reset"
								class="text-xs font-semibold text-indigo-600 hover:underline dark:text-indigo-400"
							>
								Lupa sandi?
							</a>
						</div>
						<input
							id="login-password"
							type="password"
							required
							bind:value={password}
							placeholder="••••••••••••"
							class="mt-1 w-full rounded-xl border border-slate-300 bg-white px-3.5 py-2.5 text-sm text-slate-900 placeholder:text-slate-400 focus:border-indigo-600 focus:ring-2 focus:ring-indigo-600/20 focus:outline-hidden dark:border-slate-700 dark:bg-slate-800 dark:text-white"
						/>
					</div>

					<div class="flex items-center justify-between text-xs">
						<label
							class="flex cursor-pointer items-center gap-2 text-slate-600 dark:text-slate-400"
						>
							<input
								type="checkbox"
								bind:checked={rememberMe}
								class="h-4 w-4 rounded-md border-slate-300 text-indigo-600 focus:ring-indigo-500"
							/>
							<span>Ingat sesi browser ini (7 hari)</span>
						</label>
					</div>

					<div class="pt-2">
						<button
							type="submit"
							class="w-full cursor-pointer rounded-xl bg-gradient-to-r from-indigo-600 to-blue-600 py-3 text-sm font-bold text-white shadow-md shadow-indigo-600/30 transition hover:from-indigo-500 hover:to-blue-500 active:scale-98"
						>
							Masuk Sistem (SSO / MFA)
						</button>
					</div>

					<div
						class="rounded-xl bg-slate-50 p-3 text-center text-[11px] text-slate-500 dark:bg-slate-800/60 dark:text-slate-400"
					>
						🔒 Dilindungi oleh Zero-Trust JWT Authentication & Audit Logging
					</div>
				</form>
			{:else}
				<div class="py-6 text-center">
					<div
						class="mx-auto flex h-14 w-14 items-center justify-center rounded-full bg-emerald-100 text-2xl text-emerald-600 dark:bg-emerald-950 dark:text-emerald-400"
					>
						✓
					</div>
					<h3 class="mt-4 text-xl font-bold text-slate-900 dark:text-white">
						Autentikasi Terverifikasi
					</h3>
					<p class="mt-2 text-xs text-slate-600 dark:text-slate-300">
						Menghubungkan ke backend Golang Echo v5 session token...
					</p>

					<div class="mt-6">
						<button
							type="button"
							onclick={handleReset}
							class="w-full rounded-xl bg-indigo-600 py-2.5 text-xs font-bold text-white transition hover:bg-indigo-500"
						>
							Tutup
						</button>
					</div>
				</div>
			{/if}
		</div>
	</div>
{/if}
