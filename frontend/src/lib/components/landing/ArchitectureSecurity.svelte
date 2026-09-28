<section
	id="architecture"
	class="border-t border-slate-200/80 bg-slate-900 py-20 text-white lg:py-28 dark:bg-black"
>
	<div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
		<!-- Section Header -->
		<div class="mx-auto max-w-3xl text-center">
			<div
				class="inline-flex items-center gap-1.5 rounded-full bg-indigo-500/20 px-3.5 py-1 text-xs font-bold text-indigo-300 ring-1 ring-indigo-400/30"
			>
				<span>🏛️</span>
				<span>Hexagonal Modular Monolith</span>
			</div>
			<h2 class="mt-4 text-3xl font-black tracking-tight sm:text-4xl lg:text-5xl">
				Arsitektur Performa Tinggi & Keamanan Tanpa Kompromi
			</h2>
			<p class="mt-4 text-base text-slate-300 sm:text-lg">
				Dirancang untuk skala korporasi dengan beban transaksi tinggi. Kami menggabungkan
				kesederhanaan deployment monolith dengan ketegasan arsitektur Hexagonal (Ports & Adapters).
			</p>
		</div>

		<!-- 3 Architecture & Security Pillars -->
		<div class="mt-16 grid grid-cols-1 gap-8 md:grid-cols-3">
			<!-- Pillar 1 -->
			<div
				class="rounded-3xl border border-slate-800 bg-slate-800/50 p-8 backdrop-blur-xs transition hover:border-slate-700"
			>
				<div
					class="mb-6 flex h-12 w-12 items-center justify-center rounded-2xl bg-indigo-500/20 text-2xl text-indigo-400"
				>
					⚡
				</div>
				<h3 class="text-xl font-bold text-white">Hexagonal & Anti-Circular</h3>
				<p class="mt-3 text-sm leading-relaxed text-slate-400">
					Setiap modul bisnis memiliki isolasi ketat antara Domain, Ports, dan Adapters. Tidak ada
					impor konkrit antar-modul; seluruh komunikasi sinkron memakai consumer contracts dan
					asinkron memakai EventBus decoupled.
				</p>
				<ul class="mt-6 space-y-2 border-t border-slate-800 pt-6 text-xs text-slate-300">
					<li class="flex items-center gap-2">
						<span class="text-indigo-400">●</span> Backend Golang Echo v5 Ultra-Fast
					</li>
					<li class="flex items-center gap-2">
						<span class="text-indigo-400">●</span> Native pgx/v5 simple query protocol
					</li>
					<li class="flex items-center gap-2">
						<span class="text-indigo-400">●</span> Event-driven background side-effects
					</li>
				</ul>
			</div>

			<!-- Pillar 2 -->
			<div
				class="rounded-3xl border border-slate-800 bg-slate-800/50 p-8 backdrop-blur-xs transition hover:border-slate-700"
			>
				<div
					class="mb-6 flex h-12 w-12 items-center justify-center rounded-2xl bg-emerald-500/20 text-2xl text-emerald-400"
				>
					🗄️
				</div>
				<h3 class="text-xl font-bold text-white">PgBouncer Transaction Pooling</h3>
				<p class="mt-3 text-sm leading-relaxed text-slate-400">
					Proxy PgBouncer pada mode transaksi menangani ribuan request koneksi aplikasi tanpa
					membebani alokasi RAM PostgreSQL 16. Skema migrasi otomatis dikelola oleh golang-migrate
					secara bertahap.
				</p>
				<ul class="mt-6 space-y-2 border-t border-slate-800 pt-6 text-xs text-slate-300">
					<li class="flex items-center gap-2">
						<span class="text-emerald-400">●</span> Pool Mode = Transaction (:6432)
					</li>
					<li class="flex items-center gap-2">
						<span class="text-emerald-400">●</span> Presisi finansial decimal.Decimal (No float)
					</li>
					<li class="flex items-center gap-2">
						<span class="text-emerald-400">●</span> Audit column multi-tenant otomatis
					</li>
				</ul>
			</div>

			<!-- Pillar 3 -->
			<div
				class="rounded-3xl border border-slate-800 bg-slate-800/50 p-8 backdrop-blur-xs transition hover:border-slate-700"
			>
				<div
					class="mb-6 flex h-12 w-12 items-center justify-center rounded-2xl bg-cyan-500/20 text-2xl text-cyan-400"
				>
					🛡️
				</div>
				<h3 class="text-xl font-bold text-white">Zero Trust & Compliance</h3>
				<p class="mt-3 text-sm leading-relaxed text-slate-400">
					Sertifikasi SOC-2 Type II dan kepatuhan penuh PSAK/IFRS. Dilengkapi kontrol akses granular
					RBAC hingga tingkat field data, enkripsi end-to-end, multi-factor authentication (MFA),
					dan jejak audit immutable.
				</p>
				<ul class="mt-6 space-y-2 border-t border-slate-800 pt-6 text-xs text-slate-300">
					<li class="flex items-center gap-2">
						<span class="text-cyan-400">●</span> Immutable changelog & before/after diff
					</li>
					<li class="flex items-center gap-2">
						<span class="text-cyan-400">●</span> Deteksi brute-force & pembatasan IP
					</li>
					<li class="flex items-center gap-2">
						<span class="text-cyan-400">●</span> JWT signed RBAC & branch authorization
					</li>
				</ul>
			</div>
		</div>
	</div>
</section>
