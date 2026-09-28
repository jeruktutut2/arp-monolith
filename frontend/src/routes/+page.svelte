<script lang="ts">
	import { onMount } from 'svelte';
	import type { ERPModule } from '$lib/types/landing';

	import Navbar from '$lib/components/landing/Navbar.svelte';
	import Hero from '$lib/components/landing/Hero.svelte';
	import TrustMetrics from '$lib/components/landing/TrustMetrics.svelte';
	import ModulesShowcase from '$lib/components/landing/ModulesShowcase.svelte';
	import DashboardPreview from '$lib/components/landing/DashboardPreview.svelte';
	import Solutions from '$lib/components/landing/Solutions.svelte';
	import ArchitectureSecurity from '$lib/components/landing/ArchitectureSecurity.svelte';
	import CtaBanner from '$lib/components/landing/CtaBanner.svelte';
	import Footer from '$lib/components/landing/Footer.svelte';
	import DemoModal from '$lib/components/landing/DemoModal.svelte';
	import ModuleModal from '$lib/components/landing/ModuleModal.svelte';
	import SignInModal from '$lib/components/landing/SignInModal.svelte';

	// Theme state
	let isDarkMode = $state(false);

	// Modals state
	let demoModalOpen = $state(false);
	let demoModalTitle = $state('Jadwalkan Konsultasi & Demo Produk');
	let signInModalOpen = $state(false);
	let selectedModule = $state<ERPModule | null>(null);

	onMount(() => {
		// Initialize theme from system preference or local storage
		const savedTheme = localStorage.getItem('erp_theme');
		if (savedTheme === 'dark') {
			isDarkMode = true;
			document.documentElement.classList.add('dark');
		} else if (savedTheme === 'light') {
			isDarkMode = false;
			document.documentElement.classList.remove('dark');
		} else if (window.matchMedia('(prefers-color-scheme: dark)').matches) {
			isDarkMode = true;
			document.documentElement.classList.add('dark');
		}
	});

	function toggleTheme() {
		isDarkMode = !isDarkMode;
		if (isDarkMode) {
			document.documentElement.classList.add('dark');
			localStorage.setItem('erp_theme', 'dark');
		} else {
			document.documentElement.classList.remove('dark');
			localStorage.setItem('erp_theme', 'light');
		}
	}

	function openDemoModal(title = 'Jadwalkan Konsultasi & Demo Produk') {
		demoModalTitle = title;
		demoModalOpen = true;
	}

	function closeDemoModal() {
		demoModalOpen = false;
	}

	function openSignInModal() {
		signInModalOpen = true;
	}

	function closeSignInModal() {
		signInModalOpen = false;
	}

	function openModuleDetails(mod: ERPModule) {
		selectedModule = mod;
	}

	function closeModuleDetails() {
		selectedModule = null;
	}
</script>

<svelte:head>
	<title>Enterprise ERP Monolith — Modern Business Management Platform</title>
	<meta
		name="description"
		content="Enterprise Resource Planning generasi baru berbasis Golang Echo v5, PostgreSQL, PgBouncer, dan SvelteKit 2. Satukan 25 modul bisnis dalam satu kendali."
	/>
</svelte:head>

<div
	class="flex min-h-screen flex-col justify-between bg-white text-slate-900 antialiased transition-colors duration-200 selection:bg-indigo-500 selection:text-white dark:bg-slate-950 dark:text-slate-100"
>
	<div>
		<!-- Navigation Header -->
		<Navbar
			{isDarkMode}
			onToggleTheme={toggleTheme}
			onOpenDemoModal={() => openDemoModal('Jadwalkan Konsultasi & Demo Produk')}
			onOpenSignInModal={openSignInModal}
		/>

		<!-- Hero Section -->
		<Hero
			onOpenDemoModal={() => openDemoModal('Jadwalkan Demo Langsung ERP')}
			onOpenTrialModal={() => openDemoModal('Mulai Uji Coba Gratis 30 Hari')}
		/>

		<!-- Live Metrics & Trust Banner -->
		<TrustMetrics />

		<!-- Interactive 25 Modules Showcase -->
		<ModulesShowcase onSelectModule={openModuleDetails} />

		<!-- Modern Dashboard Preview -->
		<DashboardPreview />

		<!-- Industry Solutions -->
		<Solutions />

		<!-- Architecture & Enterprise Security -->
		<ArchitectureSecurity />

		<!-- Closing Conversion CTA Banner -->
		<CtaBanner
			onOpenDemoModal={() => openDemoModal('Konsultasi Solusi Enterprise')}
			onOpenTrialModal={() => openDemoModal('Mulai Uji Coba Gratis 30 Hari')}
		/>
	</div>

	<!-- Global Enterprise Multi-Column Footer -->
	<Footer />

	<!-- Modals -->
	<DemoModal isOpen={demoModalOpen} title={demoModalTitle} onClose={closeDemoModal} />

	<ModuleModal module={selectedModule} onClose={closeModuleDetails} />

	<SignInModal isOpen={signInModalOpen} onClose={closeSignInModal} />
</div>
