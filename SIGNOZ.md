# 📊 Panduan Lengkap SigNoz & ClickHouse di Ubuntu

Dokumentasi komprehensif mengenai arsitektur, instalasi di sistem operasi **Ubuntu**, konfigurasi database **ClickHouse**, cara integrasi dengan aplikasi (khususnya backend Golang & frontend), serta pemeliharaan (*maintenance* & *troubleshooting*).

---

## 📑 Daftar Isi
1. [Pengenalan SigNoz & Arsitektur](#1-pengenalan-signoz--arsitektur)
2. [Spesifikasi & Persyaratan Sistem Ubuntu](#2-spesifikasi--persyaratan-sistem-ubuntu)
3. [Langkah-Langkah Instalasi di Ubuntu](#3-langkah-langkah-instalasi-di-ubuntu)
4. [Manajemen & Eksplorasi Database ClickHouse](#4-manajemen--eksplorasi-database-clickhouse)
5. [Cara Menggunakan SigNoz (Integrasi Aplikasi)](#5-cara-menggunakan-signoz-integrasi-aplikasi)
6. [Fitur-Fitur Utama Web UI SigNoz](#6-fitur-fitur-utama-web-ui-signoz)
7. [Operasional, Pemeliharaan & Troubleshooting](#7-operasional-pemeliharaan--troubleshooting)

---

## 1. Pengenalan SigNoz & Arsitektur

**SigNoz** adalah platform observabilitas *open-source* (APM - *Application Performance Monitoring*) berbasis standar **OpenTelemetry (OTel)**. SigNoz menggabungkan tiga pilar utama observabilitas dalam satu tampilan (*single pane of glass*):
- **Traces**: Pelacakan jejak request terdistribusi (*distributed tracing*) dari frontend hingga query database.
- **Metrics**: Pengukuran performa seperti CPU, Memory, RPS (*Requests per Second*), dan Latensi (p50, p95, p99).
- **Logs**: Pengumpulan dan pencarian log terstruktur dengan korelasi otomatis ke Trace ID dan Span ID.

### 🏛️ Komponen & Arsitektur Penyimpanan

```
+-------------------------------------------------------------+
|                     Aplikasi / Klien                        |
|   (Backend Go / Frontend SvelteKit / Node.js / Python)      |
+-------------------------------------------------------------+
                              │
               OTLP gRPC (4317) / HTTP (4318)
                              ▼
+-------------------------------------------------------------+
|                OpenTelemetry (OTel) Collector               |
|      (Menerima, memvalidasi, memproses, & mem-batch data)   |
+-------------------------------------------------------------+
                              │
                    Batch Insert (Port 9000)
                              ▼
+-------------------------------------------------------------+
|                 ClickHouse Database (Kolom)                 |
|   - signoz_traces  : Data distributed tracing & span        |
|   - signoz_metrics : Metrik performa, histogram & time-series|
|   - signoz_logs    : Data log terstruktur                   |
+-------------------------------------------------------------+
                              ▲
                      SQL Query Engine
                              │
+-------------------------------------------------------------+
|                    SigNoz Query Service                     |
|           (Backend query handler & alerting engine)         |
+-------------------------------------------------------------+
                              │
                              ▼
+-------------------------------------------------------------+
|                 SigNoz Frontend (Web UI :3301)              |
+-------------------------------------------------------------+
```

### 🗄️ Mengapa ClickHouse Dipilih Sebagai Database Utama?
SigNoz tidak menggunakan database relasional biasa (seperti PostgreSQL atau MySQL) untuk menyimpan telemetri karena volume data telemetri dapat mencapai jutaan hingga miliaran baris per hari.
- **Penyimpanan Berorientasi Kolom (*Columnar Storage*)**: Data dikompresi hingga 5x - 10x lebih hemat kapasitas dibanding basis data baris biasa.
- **Agregasi Sangat Cepat**: Query analitik (*aggregation queries* seperti `AVG`, `PERCENTILE`, `COUNT`) atas miliaran baris selesai dalam hitungan milidetik.
- **Skalabilitas**: Mampu menangani beban *high-throughput write* tanpa mengalami *lock bottleneck*.

---

## 2. Spesifikasi & Persyaratan Sistem Ubuntu

### 🖥️ Kebutuhan Perangkat Keras (Hardware)
| Kebutuhan | Dev / Testing Lingkungan Kecil | Production (Skala Menengah - Besar) |
|---|---|---|
| **OS** | Ubuntu 20.04 / 22.04 / 24.04 LTS (x86_64 atau ARM64) | Ubuntu 22.04 / 24.04 LTS (64-bit) |
| **CPU** | Minimal 2 Core | 4 - 8+ Core vCPU |
| **RAM** | Minimal 4 GB (8 GB direkomendasikan) | 16 GB - 32 GB+ RAM |
| **Disk** | 30 GB SSD | 100 GB - 500 GB+ NVMe SSD |

> [!WARNING]
> ClickHouse dan OTel Collector membutuhkan alokasi memori yang stabil. Jangan menjalankan SigNoz pada server Ubuntu dengan RAM kurang dari 4 GB karena proses ClickHouse dapat terbunuh oleh Linux OOM (*Out Of Memory*) Killer.

### 🔌 Alokasi Port Jaringan (Firewall / UFW)
Pastikan port-port berikut diizinkan pada firewall jika diakses dari luar server:
- **`3301`** : Dashboard Web UI SigNoz (HTTP)
- **`4317`** : OpenTelemetry gRPC Receiver (digunakan oleh SDK aplikasi)
- **`4318`** : OpenTelemetry HTTP Receiver (digunakan oleh SDK frontend/browser atau HTTP clients)
- **`8123`** : ClickHouse HTTP Interface (opsional untuk audit internal/query)
- **`9000`** : ClickHouse Native TCP Client (internal antar kontainer)

---

## 3. Langkah-Langkah Instalasi di Ubuntu

Metode resmi dan paling stabil untuk menjalankan SigNoz di Ubuntu adalah menggunakan **Docker Engine** dan **Docker Compose v2**.

### Langkah 1: Update Sistem & Instalasi Dependency Dasar
Buka terminal server Ubuntu Anda dan jalankan pembaruan sistem:
```bash
sudo apt update && sudo apt upgrade -y
sudo apt install -y curl wget git apt-transport-https ca-certificates gnupg lsb-release
```

---

### Langkah 2: Instalasi Docker Engine & Docker Compose
Jika server Ubuntu Anda belum terpasang Docker resmi, ikuti langkah berikut:

```bash
# Tambahkan GPG key resmi Docker
sudo install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
sudo chmod a+r /etc/apt/keyrings/docker.gpg

# Tambahkan repositori Docker ke APT
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

# Instal Docker CE dan plugin Docker Compose
sudo apt update
sudo apt install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

# Masukkan user saat ini ke grup docker (agar tidak perlu selalu sudo)
sudo usermod -aG docker $USER
```
> [!NOTE]
> Setelah menjalankan `usermod`, lakukan relog terminal atau jalankan `newgrp docker` agar hak akses grup aktif.

Verifikasi instalasi:
```bash
docker --version
docker compose version
```

---

### Langkah 3: Clone Repositori SigNoz & Instalasi

#### Opsi A — Skrip Instalasi Otomatis Resmi (Rekomendasi)
SigNoz menyediakan skrip instalasi yang secara otomatis mendeteksi konfigurasi sistem dan memasang seluruh stack:
```bash
# Clone repositori deploy
git clone -b main https://github.com/SigNoz/signoz.git
cd signoz/deploy/

# Jalankan skrip instalasi
./install.sh
```
Skrip akan memvalidasi ketersediaan port, memverifikasi alokasi memori, mengunduh image container, lalu menjalankan seluruh kontainer di latar belakang.

---

#### Opsi B — Menjalankan via Docker Compose Manual
Jika Anda ingin kontrol penuh atas file konfigurasi dan port:
```bash
git clone -b main https://github.com/SigNoz/signoz.git
cd signoz/deploy/docker/clickhouse-setup

# Jalankan semua service di latar belakang
docker compose up -d
```

---

### Langkah 4: Verifikasi Status Kontainer
Pastikan seluruh kontainer berjalan normal dengan status `Up` atau `healthy`:
```bash
docker compose ps
```
Kontainer yang harus aktif meliputi:
1. `signoz-clickhouse`: Mesin database ClickHouse.
2. `signoz-otel-collector`: Pengumpul dan pemroses data OTel.
3. `signoz-query-service`: Backend query dan agregasi log/trace.
4. `signoz-frontend`: Antarmuka Web Dashboard.
5. `signoz-alertmanager`: Komponen pemicu notifikasi alert.
6. `signoz-schema-migrator`: Mengelola migrasi tabel di ClickHouse.

---

### Langkah 5: Mengakses Web UI SigNoz
1. Buka browser dan arahkan ke alamat IP server Ubuntu Anda:
   ```
   http://<IP_SERVER_UBUNTU>:3301
   ```
2. Anda akan disambut halaman pendaftaran pertama kali (*First-time onboarding*):
   - Masukkan **Nama**, **Email**, dan **Password Admin**.
   - Selesaikan wizard setup awal untuk masuk ke Dashboard utama.

---

## 4. Manajemen & Eksplorasi Database ClickHouse

SigNoz mengelola database ClickHouse secara otomatis, namun sebagai administrator sistem, memahami struktur dan cara mengelola ClickHouse sangat penting.

### 🔍 Akses Command Line Interface (CLI) ClickHouse
Anda dapat langsung berinteraksi dengan database ClickHouse menggunakan perintah `clickhouse-client` bawaan dari dalam kontainer:

```bash
docker exec -it signoz-clickhouse clickhouse-client
```

---

### 📂 Struktur Database & Tabel SigNoz
Di dalam CLI ClickHouse, jalankan perintah SQL berikut:

```sql
-- Melihat daftar database
SHOW DATABASES;
```
Anda akan melihat database utama:
- `signoz_traces`
- `signoz_metrics`
- `signoz_logs`

#### Tabel-Tabel Kunci:
```sql
-- Pindah ke database traces
USE signoz_traces;
SHOW TABLES;
-- Berisi: signoz_index_v2, signoz_spans, signoz_error_index_v2, durationSort, dll.

-- Pindah ke database logs
USE signoz_logs;
SHOW TABLES;
-- Berisi: logs, logs_v2, schema_migrations, dll.

-- Pindah ke database metrics
USE signoz_metrics;
SHOW TABLES;
-- Berisi: samples_v4, time_series_v4, metadata, dll.
```

---

### 📊 Query Berguna untuk Monitoring ClickHouse

#### 1. Mengetahui Penggunaan Disk Penyimpanan Tiap Tabel:
```sql
SELECT 
    table, 
    formatReadableSize(sum(bytes)) AS size, 
    sum(rows) AS total_rows 
FROM system.parts 
WHERE active AND database IN ('signoz_traces', 'signoz_metrics', 'signoz_logs') 
GROUP BY table 
ORDER BY sum(bytes) DESC;
```

#### 2. Mengetahui Jumlah Request / Span Terakhir:
```sql
SELECT 
    serviceName, 
    count() AS total_spans, 
    formatReadableTimeDelta(avg(durationNano)/1000000000) AS avg_duration 
FROM signoz_traces.signoz_index_v2 
WHERE timestamp > now() - INTERVAL 1 HOUR 
GROUP BY serviceName 
ORDER BY total_spans DESC;
```

---

### ⏳ Pengaturan Data Retention (TTL - Time to Live)

Secara default, data telemetri dapat menghabiskan ruang disk jika tidak dibatasi durasi penyimpanannya.

#### Cara 1: Mengatur Melalui Web UI SigNoz (Direkomendasikan)
1. Buka Web UI SigNoz (`http://<IP_UBUNTU>:3301`).
2. Masuk ke menu **Settings** (ikon gerigi) -> **General Settings** -> **Data Retention Period**.
3. Atur durasi penyimpanan (misalnya: Traces = `15 hari`, Logs = `7 hari`, Metrics = `30 hari`).
4. Klik **Save**. SigNoz akan mengirimkan query `ALTER TABLE ... MODIFY TTL` secara otomatis ke ClickHouse.

#### Cara 2: Memeriksa atau Mengubah TTL Manual di ClickHouse
Untuk memeriksa TTL aktif pada tabel logs:
```sql
SHOW CREATE TABLE signoz_logs.logs_v2;
```
Untuk memperbarui retensi secara manual (contoh: 7 hari):
```sql
ALTER TABLE signoz_logs.logs_v2 MODIFY TTL toDateTime(timestamp / 1000000000) + INTERVAL 7 DAY;
```

---

### 💾 Backup dan Restore ClickHouse
Untuk membuat cadangan database ClickHouse pada lingkungan produksi:
1. **Docker Volume Snapshot**: Backup direktori data volume Docker (`/var/lib/docker/volumes/...` atau direktori mapping `clickhouse-data`).
2. **clickhouse-backup**: Gunakan utility [clickhouse-backup](https://github.com/AlexAkulov/clickhouse-backup) untuk snapshot tabel ke cloud storage (S3 / MinIO / Local Backup).

---

## 5. Cara Menggunakan SigNoz (Integrasi Aplikasi)

Setelah SigNoz berjalan di server Ubuntu, aplikasi Anda dapat mulai mengirimkan data trace, metrik, dan log ke port OTel Collector (`4317` untuk gRPC atau `4318` untuk HTTP).

### A. Integrasi Backend Golang (Echo / Standard Library)

Berikut contoh implementasi standar menggunakan library resmi OpenTelemetry Go SDK:

#### 1. Pasang Dependency OTel di Proyek Go:
```bash
go get go.opentelemetry.io/otel \
       go.opentelemetry.io/otel/trace \
       go.opentelemetry.io/otel/sdk \
       go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc \
       go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho
```

#### 2. Inisialisasi Tracer Provider (`telemetry.go`):
```go
package telemetry

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// InitTracer menginisialisasi OTel Tracer Provider dan mengekspor span ke SigNoz
func InitTracer(ctx context.Context, serviceName, signozEndpoint string) (*sdktrace.TracerProvider, error) {
	// Koneksi ke OTel Collector di SigNoz (port 4317)
	conn, err := grpc.DialContext(ctx, signozEndpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithTimeout(3*time.Second),
	)
	if err != nil {
		return nil, err
	}

	exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, err
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
			semconv.DeploymentEnvironmentKey.String("production"),
		),
	)
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return tp, nil
}
```

#### 3. Menerapkan Middleware pada Echo Server (`main.go`):
```go
package main

import (
	"context"
	"log"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"yourproject/telemetry"
)

func main() {
	ctx := context.Background()

	// Inisialisasi OTel: gRPC port 4317 SigNoz
	tp, err := telemetry.InitTracer(ctx, "erp-backend-monolith", "localhost:4317")
	if err != nil {
		log.Printf("Gagal inisialisasi OTel ke SigNoz: %v", err)
	} else {
		defer func() {
			_ = tp.Shutdown(ctx)
		}()
	}

	e := echo.New()

	// Tambahkan middleware tracing ke Echo
	e.Use(otelecho.Middleware("erp-backend-monolith"))

	e.GET("/api/v1/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "ok"})
	})

	e.Logger.Fatal(e.Start(":8080"))
}
```

Setiap request HTTP ke Echo server Anda sekarang otomatis tercatat sebagai Span di SigNoz lengkap dengan status kode HTTP, durasi proses, dan URI endpoint.

---

### B. Integrasi Frontend (SvelteKit / Web Browser)
Untuk melacak latensi halaman, resource load, dan user action dari browser:
1. Pasang paket npm OpenTelemetry Web:
   ```bash
   bun add @opentelemetry/sdk-trace-web \
           @opentelemetry/exporter-trace-otlp-http \
           @opentelemetry/instrumentation-document-load \
           @opentelemetry/instrumentation-user-interaction
   ```
2. Konfigurasikan exporter agar mengirim ke endpoint SigNoz HTTP (`http://<IP_UBUNTU>:4318/v1/traces`).
3. Seluruh navigasi halaman dan AJAX request akan terhubung langsung ke trace backend (*distributed trace context* melalui header `traceparent`).

---

## 6. Fitur-Fitur Utama Web UI SigNoz

1. **Services / APM Dashboard**:
   - Menampilkan daftar semua microservice / monolit yang aktif.
   - Grafik otomatis: **Application Latency** (p99, p95, p50), **Operations per Second (RPS)**, dan **Error Percentage**.
   - Menampilkan dependency map (komunikasi antar-service).

2. **Distributed Tracing**:
   - Pencarian span berdasarkan filter: HTTP Status (`status_code >= 400`), Service Name, Min Duration, atau Custom Tags.
   - **Flamegraph & Gantt Chart View**: Melihat bagian kode atau query database mana yang memakan waktu paling lama.

3. **Log Explorer**:
   - Pencarian log secara terpusat dengan kemampuan filtering live tail.
   - Klik satu baris log untuk langsung melompat ke **Trace Terkait** (*Log-to-Trace correlation*).

4. **Alerts & Notification Channels**:
   - Pembuatan aturan alert berdasarkan ambang batas Latensi, Error Rate, atau ekspresi PromQL.
   - Pengiriman alert langsung ke **Slack**, **Discord**, **PagerDuty**, **Webhook**, atau **Email**.

---

## 7. Operasional, Pemeliharaan & Troubleshooting

### 🛠️ Perintah Operasional Docker SigNoz
Masuk ke direktori deploy SigNoz terlebih dahulu (`cd signoz/deploy/docker/clickhouse-setup`):

```bash
# Menjalankan SigNoz
docker compose up -d

# Menghentikan SigNoz sementara
docker compose stop

# Menyalakan kembali setelah stop
docker compose start

# Mematikan seluruh stack kontainer
docker compose down

# Melihat log dari semua service
docker compose logs -f

# Melihat log khusus OTel Collector
docker compose logs -f signoz-otel-collector

# Melihat log khusus ClickHouse
docker compose logs -f signoz-clickhouse
```

---

### 🚨 Panduan Troubleshooting Masalah Umum

#### 1. Port Sudah Terpakai (*Port Conflict*)
- **Gejala**: Error `bind: address already in use` saat `docker compose up`.
- **Solusi**: Periksa service yang menggunakan port tersebut:
  ```bash
  sudo ss -tulpn | grep -E '3301|4317|4318|8123|9000'
  ```
  Jika ada service lain (misalnya clickhouse lokal atau web server lain), ubah mapping port di `docker-compose.yaml` (contoh: ubah `"3301:3301"` menjadi `"3302:3301"`).

#### 2. ClickHouse Sering Restart atau Status `Exited (137)` (OOM Killer)
- **Gejala**: Kontainer `signoz-clickhouse` tiba-tiba mati dan hidup kembali.
- **Penyebab**: Server kehabisan RAM (*Linux Out Of Memory Killer* mematikan proses ClickHouse).
- **Solusi**:
  1. Periksa log kernel: `dmesg -T | grep -i oom`.
  2. Tambahkan Swap file pada server Ubuntu jika RAM terbatas:
     ```bash
     sudo fallocate -l 4G /swapfile
     sudo chmod 600 /swapfile
     sudo mkswap /swapfile
     sudo swapon /swapfile
     echo '/swapfile none swap sw 0 0' | sudo tee -a /etc/fstab
     ```
  3. Konfigurasikan batas memori ClickHouse pada file konfigurasi `config.xml` (`<max_server_memory_usage>`).

#### 3. Data Telemetri Tidak Muncul di Dashboard SigNoz
- **Penyebab & Solusi**:
  1. Pastikan OTel Collector dapat dijangkau dari aplikasi:
     ```bash
     nc -zv <IP_SERVER_UBUNTU> 4317
     ```
  2. Periksa log OTel Collector:
     ```bash
     docker compose logs -f --tail=100 signoz-otel-collector
     ```
  3. Pastikan waktu/jam server Ubuntu dan jam mesin aplikasi sinkron (Gunakan `chrony` atau `systemd-timesyncd`). Jika selisih waktu terlalu besar, data trace bisa dianggap kadaluarsa dan tidak muncul di filter default rentang waktu.

---

### 🔄 Cara Upgrade SigNoz ke Versi Terbaru
Untuk memperbarui SigNoz ke versi rilis terbaru:
```bash
cd signoz/deploy/docker/clickhouse-setup

# Tarik perubahan kode terbaru
git pull origin main

# Unduh image Docker versi terbaru
docker compose pull

# Jalankan ulang kontainer dengan image baru
docker compose up -d --remove-orphans
```
Proses migrasi skema tabel ClickHouse akan otomatis dijalankan oleh kontainer `signoz-schema-migrator`.
