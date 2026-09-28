<script lang="ts">
	import { resolve } from '$app/paths';

	interface Props {
		isDarkMode: boolean;
		onToggleTheme: () => void;
		onOpenDemoModal: () => void;
		onOpenSignInModal: () => void;
	}

	let { isDarkMode, onToggleTheme, onOpenDemoModal, onOpenSignInModal }: Props = $props();

	let mobileMenuOpen = $state(false);

	function toggleMobileMenu() {
		mobileMenuOpen = !mobileMenuOpen;
	}

	function closeMobileMenu() {
		mobileMenuOpen = false;
	}
</script>

<header
	class="sticky top-0 z-50 border-b border-slate-200/80 bg-white/90 backdrop-blur-md transition-colors duration-200 dark:border-slate-800/80 dark:bg-slate-950/90"
>
	<div class="mx-auto flex h-20 max-w-7xl items-center justify-between px-4 sm:px-6 lg:px-8">
		<!-- Brand Logo -->
		<a href={resolve('/')} class="group flex items-center gap-3">
			<div
				class="flex h-11 w-11 items-center justify-center rounded-xl bg-gradient-to-br from-indigo-600 via-blue-600 to-indigo-700 text-lg font-black text-white shadow-md ring-1 shadow-indigo-500/20 ring-white/20 transition-transform group-hover:scale-105"
			>
				E
			</div>
			<div>
				<div class="flex items-center gap-1.5">
					<span
						class="text-xl font-black tracking-tight text-slate-900 transition-colors dark:text-white"
					>
						ERP<span class="text-indigo-600 dark:text-indigo-400">System</span>
					</span>
					<span
						class="rounded-full bg-indigo-50 px-2 py-0.5 text-[10px] font-bold tracking-wider text-indigo-700 uppercase ring-1 ring-indigo-500/20 dark:bg-indigo-950/80 dark:text-indigo-300"
					>
						v2.4
					</span>
				</div>
				<p class="text-[11px] font-medium text-slate-500 dark:text-slate-400">
					Enterprise Monolith Platform
				</p>
			</div>
		</a>

		<!-- Desktop Navigation Links -->
		<nav class="hidden items-center gap-1 md:flex">
			<a
				href="#features"
				class="rounded-lg px-3.5 py-2 text-sm font-semibold text-slate-600 transition hover:bg-slate-100 hover:text-slate-950 dark:text-slate-300 dark:hover:bg-slate-800/60 dark:hover:text-white"
			>
				Keunggulan
			</a>
			<a
				href="#modules"
				class="rounded-lg px-3.5 py-2 text-sm font-semibold text-slate-600 transition hover:bg-slate-100 hover:text-slate-950 dark:text-slate-300 dark:hover:bg-slate-800/60 dark:hover:text-white"
			>
				25 Modul
			</a>
			<a
				href="#preview"
				class="rounded-lg px-3.5 py-2 text-sm font-semibold text-slate-600 transition hover:bg-slate-100 hover:text-slate-950 dark:text-slate-300 dark:hover:bg-slate-800/60 dark:hover:text-white"
			>
				Dashboard
			</a>
			<a
				href="#solutions"
				class="rounded-lg px-3.5 py-2 text-sm font-semibold text-slate-600 transition hover:bg-slate-100 hover:text-slate-950 dark:text-slate-300 dark:hover:bg-slate-800/60 dark:hover:text-white"
			>
				Solusi Industri
			</a>
			<a
				href="#architecture"
				class="rounded-lg px-3.5 py-2 text-sm font-semibold text-slate-600 transition hover:bg-slate-100 hover:text-slate-950 dark:text-slate-300 dark:hover:bg-slate-800/60 dark:hover:text-white"
			>
				Arsitektur & Security
			</a>
		</nav>

		<!-- Right Side Controls & Actions -->
		<div class="hidden items-center gap-3 sm:flex">
			<!-- Theme Toggle Button -->
			<button
				type="button"
				onclick={onToggleTheme}
				aria-label="Toggle theme"
				class="flex h-10 w-10 items-center justify-center rounded-xl border border-slate-200 bg-slate-50 text-slate-600 transition hover:border-slate-300 hover:bg-slate-100 dark:border-slate-800 dark:bg-slate-900 dark:text-slate-300 dark:hover:border-slate-700 dark:hover:bg-slate-800"
			>
				{#if isDarkMode}
					<svg class="h-5 w-5 text-amber-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z"
						/>
					</svg>
				{:else}
					<svg
						class="h-5 w-5 text-indigo-600"
						fill="none"
						viewBox="0 0 24 24"
						stroke="currentColor"
					>
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z"
						/>
					</svg>
				{/if}
			</button>

			<!-- Sign In Button -->
			<button
				type="button"
				onclick={onOpenSignInModal}
				class="rounded-xl px-4 py-2.5 text-sm font-semibold text-slate-700 transition hover:bg-slate-100 dark:text-slate-200 dark:hover:bg-slate-800"
			>
				Sign In
			</button>

			<!-- Schedule Demo CTA -->
			<button
				type="button"
				onclick={onOpenDemoModal}
				class="inline-flex items-center gap-2 rounded-xl bg-gradient-to-r from-indigo-600 to-blue-600 px-4 py-2.5 text-sm font-semibold text-white shadow-sm shadow-indigo-600/30 transition hover:from-indigo-500 hover:to-blue-500 active:scale-95"
			>
				<span>Schedule Demo</span>
				<svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"
					/>
				</svg>
			</button>
		</div>

		<!-- Mobile Menu Button -->
		<div class="flex items-center gap-2 sm:hidden">
			<button
				type="button"
				onclick={onToggleTheme}
				aria-label="Toggle theme"
				class="flex h-10 w-10 items-center justify-center rounded-lg border border-slate-200 bg-slate-50 text-slate-600 dark:border-slate-800 dark:bg-slate-900 dark:text-slate-300"
			>
				{#if isDarkMode}
					<svg class="h-5 w-5 text-amber-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z"
						/>
					</svg>
				{:else}
					<svg
						class="h-5 w-5 text-indigo-600"
						fill="none"
						viewBox="0 0 24 24"
						stroke="currentColor"
					>
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z"
						/>
					</svg>
				{/if}
			</button>

			<button
				type="button"
				onclick={toggleMobileMenu}
				aria-label="Open navigation menu"
				class="flex h-10 w-10 items-center justify-center rounded-lg border border-slate-200 bg-slate-50 text-slate-700 dark:border-slate-800 dark:bg-slate-900 dark:text-slate-200"
			>
				<svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
					{#if mobileMenuOpen}
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M6 18L18 6M6 6l12 12"
						/>
					{:else}
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M4 6h16M4 12h16M4 18h16"
						/>
					{/if}
				</svg>
			</button>
		</div>
	</div>

	<!-- Mobile Dropdown Navigation -->
	{#if mobileMenuOpen}
		<div
			class="border-b border-slate-200 bg-white px-4 pt-2 pb-6 sm:hidden dark:border-slate-800 dark:bg-slate-950"
		>
			<div class="flex flex-col gap-2">
				<a
					href="#features"
					onclick={closeMobileMenu}
					class="rounded-lg px-3 py-2 text-base font-semibold text-slate-700 hover:bg-slate-100 dark:text-slate-200 dark:hover:bg-slate-800"
				>
					Keunggulan
				</a>
				<a
					href="#modules"
					onclick={closeMobileMenu}
					class="rounded-lg px-3 py-2 text-base font-semibold text-slate-700 hover:bg-slate-100 dark:text-slate-200 dark:hover:bg-slate-800"
				>
					25 Modul
				</a>
				<a
					href="#preview"
					onclick={closeMobileMenu}
					class="rounded-lg px-3 py-2 text-base font-semibold text-slate-700 hover:bg-slate-100 dark:text-slate-200 dark:hover:bg-slate-800"
				>
					Dashboard
				</a>
				<a
					href="#solutions"
					onclick={closeMobileMenu}
					class="rounded-lg px-3 py-2 text-base font-semibold text-slate-700 hover:bg-slate-100 dark:text-slate-200 dark:hover:bg-slate-800"
				>
					Solusi Industri
				</a>
				<a
					href="#architecture"
					onclick={closeMobileMenu}
					class="rounded-lg px-3 py-2 text-base font-semibold text-slate-700 hover:bg-slate-100 dark:text-slate-200 dark:hover:bg-slate-800"
				>
					Arsitektur & Security
				</a>
				<div class="mt-4 flex flex-col gap-2 border-t border-slate-200 pt-4 dark:border-slate-800">
					<button
						type="button"
						onclick={() => {
							closeMobileMenu();
							onOpenSignInModal();
						}}
						class="w-full rounded-xl border border-slate-300 py-2.5 text-center text-sm font-semibold text-slate-800 dark:border-slate-700 dark:text-slate-200"
					>
						Sign In
					</button>
					<button
						type="button"
						onclick={() => {
							closeMobileMenu();
							onOpenDemoModal();
						}}
						class="w-full rounded-xl bg-indigo-600 py-2.5 text-center text-sm font-semibold text-white shadow-sm"
					>
						Schedule Demo
					</button>
				</div>
			</div>
		</div>
	{/if}
</header>
