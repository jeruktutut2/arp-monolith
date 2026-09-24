# 🎨 Design UI — Profil Pengguna (User Profile)

> Design mockup halaman Profil Pengguna untuk modul `USR` (User Management) dalam ERP System.

---

## Light Mode

![Profil Pengguna - Light Mode](./profile-light-mode.jpg)

**Spesifikasi Warna Light Mode:**

| Elemen | Warna | Kode |
|--------|-------|------|
| Background halaman | Abu-abu terang | `#f3f4f6` |
| Card / Panel | Putih | `#ffffff` |
| Sidebar | Putih | `#ffffff` |
| Teks utama | Gelap | `#1f2937` |
| Teks sekunder | Abu-abu | `#6b7280` |
| Label field | Abu-abu muda | `#9ca3af` |
| Border | Abu-abu tipis | `#e5e7eb` |
| Primary / Aksen | Biru | `#2563eb` |
| Active nav item bg | Biru muda | `#eff6ff` |

---

## Dark Mode

![Profil Pengguna - Dark Mode](./profile-dark-mode.jpg)

**Spesifikasi Warna Dark Mode:**

| Elemen | Warna | Kode |
|--------|-------|------|
| Background halaman | Navy gelap | `#0f172a` |
| Card / Panel | Slate gelap | `#1e293b` |
| Sidebar | Navy gelap | `#0f172a` |
| Teks utama | Putih | `#f8fafc` |
| Teks sekunder | Abu-abu terang | `#94a3b8` |
| Label field | Abu-abu | `#64748b` |
| Border | Slate | `#334155` |
| Primary / Aksen | Biru | `#3b82f6` |
| Active nav item bg | Biru tinted | `#1e3a5f` |

---

## Default Button Colors

| Tipe | Light Mode | Dark Mode | Kode Warna |
|------|-----------|-----------|------------|
| **Default** | Putih + border abu-abu | Gelap + border slate | bg `#fff` / `#1e293b`, border `#d1d5db` / `#475569` |
| **Primary** | Biru solid | Biru solid | `#2563eb` (light) / `#3b82f6` (dark) |
| **Success** | Hijau solid | Hijau solid | `#16a34a` (light) / `#22c55e` (dark) |
| **Warning** | Amber solid | Amber solid | `#f59e0b` (light) / `#f59e0b` (dark) |
| **Danger** | Merah solid | Merah solid | `#ef4444` (light) / `#ef4444` (dark) |

### Button States

| State | Perubahan |
|-------|-----------|
| **Hover** | Warna sedikit lebih gelap + shadow lebih kuat |
| **Active** | Scale 0.98 (sedikit mengecil) |
| **Disabled** | Opacity 50%, cursor not-allowed |

---

## Komponen Halaman

1. **Sidebar** — Navigasi modul ERP dengan ikon
2. **Top Header** — Breadcrumb + ikon search & notifikasi
3. **Profile Header Card** — Banner gradient, avatar, nama, role badge, meta info, tombol edit
4. **Tab Navigation** — Informasi Pribadi, Kontak, Preferensi
5. **Informasi Pribadi** — Grid 2 kolom data karyawan
6. **Keamanan Akun** — Status login, 2FA, password
7. **Aktivitas Terbaru** — Feed aktivitas dengan timestamp
8. **Button Showcase** — Default, Primary, Success, Warning, Danger

---

*File design disimpan di folder `ui/`*
