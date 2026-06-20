# 📦 EASM Digital Asset Management & Security Scanner Portal

<p align="center">
  <b>External Attack Surface Management (EASM)</b><br>
  Hệ thống quản lý tài sản số và quét bề mặt tấn công an ninh mạng
</p>

---

Dự án này là cổng quản lý tài sản số (Asset Management) và quét an ninh bề mặt tấn công (EASM Scanner) được xây dựng theo mô hình **Clean Architecture** sử dụng **Go**, **PostgreSQL** và **React Portal (Vite)**.

Đây là phiên bản tích hợp đầy đủ của cả **Homework Day 1** (phiên bản In-memory ban đầu) và nâng cấp hoàn chỉnh cho **Homework Day 3 (Final)** — kết nối PostgreSQL, tự động chạy Migrations, mở rộng 7 Scanner Engines, bộ Unit Tests đa tầng, tích hợp giao diện Web Portal hỗ trợ Light/Dark mode, CI/CD pipeline, deploy Cloud VM, HTTPS và Auto Deploy.

---

## 📋 Mục Lục

- [Yêu Cầu Hệ Thống](#-yêu-cầu-hệ-thống)
- [Khởi Chạy Nhanh](#-khởi-chạy-nhanh)
  - [Cách 1: Chạy bằng Docker Compose (Khuyến nghị)](#cách-1-chạy-bằng-docker-compose-khuyến-nghị)
  - [Cách 2: Chạy thủ công (Development)](#cách-2-chạy-thủ-công-development)
- [Cấu Trúc Thư Mục](#-cấu-trúc-thư-mục)
- [Biến Môi Trường](#-biến-môi-trường)
- [Danh Sách API Endpoints](#-danh-sách-api-endpoints)
- [Hướng Dẫn Test API](#-hướng-dẫn-test-api)
  - [Bài 1: Database Migration (PostgreSQL)](#bài-1-database-migration-postgresql)
  - [Bài 2: Mở Rộng EASM Scanner API](#bài-2-mở-rộng-easm-scanner-api)
  - [Bài 3: Chạy Unit Tests](#bài-3-chạy-unit-tests)
  - [Bài 4: Tích Hợp Frontend React Portal](#bài-4-tích-hợp-frontend-react-portal)
- [Hướng Dẫn Test với Postman](#-hướng-dẫn-test-với-postman)
- [CI/CD Pipeline](#-cicd-pipeline)
- [Deploy & Production](#-deploy--production)
- [Kiểm Thử API Day 1 (In-Memory)](#-kiểm-thử-api-day-1-in-memory)

---

## 🔧 Yêu Cầu Hệ Thống

| Thành phần         | Phiên bản tối thiểu | Ghi chú                               |
| ------------------ | ------------------- | ------------------------------------- |
| **Go**             | 1.21+               | Khuyến nghị 1.23+                     |
| **Docker Desktop** | Latest              | Để chạy PostgreSQL & full stack       |
| **Node.js**        | 18+                 | Để chạy Frontend (dev mode)           |
| **npm**            | 9+                  | Đi kèm Node.js                        |
| **curl**           | Any                 | Hoặc PowerShell/Postman để test API   |
| **PostgreSQL**     | 15+                 | Chạy qua Docker (không cần cài riêng) |

---

## 🚀 Khởi Chạy Nhanh

### Cách 1: Chạy bằng Docker Compose (Khuyến nghị)

Cách này sẽ tự động khởi động đầy đủ **3 services**: PostgreSQL, Backend API và Frontend Portal.

#### Linux / macOS (Bash)

```bash
# Clone repository
git clone https://github.com/thanhhieutiet/cmc-dev.git
cd cmc-dev/homework-day1

# Khởi động toàn bộ stack
docker compose up -d

# Kiểm tra trạng thái
docker compose ps

# Xem logs
docker compose logs -f
```

#### Windows PowerShell

```powershell
# Clone repository
git clone https://github.com/thanhhieutiet/cmc-dev.git
cd cmc-dev\homework-day1

# Khởi động toàn bộ stack
docker compose up -d

# Kiểm tra trạng thái
docker compose ps

# Xem logs
docker compose logs -f
```

**Sau khi khởi động:**

- 🌐 Frontend Portal: [http://localhost:3000](http://localhost:3000)
- 🔌 Backend API: [http://localhost:8080](http://localhost:8080)
- 🗄️ PostgreSQL: `localhost:5432`

**Dừng services:**

```bash
docker compose down          # Dừng và xóa containers
docker compose down -v       # Dừng và xóa cả volumes (dữ liệu DB)
```

---

### Cách 2: Chạy thủ công (Development)

Phù hợp khi phát triển và debug. Cần chạy 3 terminal riêng biệt.

#### Bước 1: Khởi động Database

```bash
# Chạy chỉ PostgreSQL container
docker compose up -d db

# Đợi DB sẵn sàng
docker compose logs db
```

#### Bước 2: Khởi động Backend

##### Linux / macOS (Bash)

```bash
cd homework-day1
go mod tidy
go run ./cmd/api/main.go
```

##### Windows PowerShell

```powershell
cd homework-day1
go mod tidy
go run ./cmd/api/main.go
```

> **Lưu ý:** Backend sẽ tự động chạy database migrations khi khởi động. API server lắng nghe tại `http://localhost:8080`.

#### Bước 3: Khởi động Frontend

##### Linux / macOS (Bash)

```bash
cd homework-day1/frontend
npm install
npm run dev
```

##### Windows PowerShell

```powershell
cd homework-day1\frontend
npm install
npm run dev
```

> Frontend chạy tại [http://localhost:5173](http://localhost:5173) (Vite dev server).

---

## 📁 Cấu Trúc Thư Mục

Dự án được tổ chức theo mô hình **Clean Architecture**:

```
homework-day1/
├── cmd/
│   └── api/
│       └── main.go                  # Entry point — khởi tạo DB, migrations, server
├── internal/
│   ├── config/
│   │   └── config.go                # Đọc cấu hình từ .env và biến môi trường
│   ├── model/
│   │   ├── asset.go                 # Entity Asset + validation + interfaces
│   │   ├── asset_test.go            # Unit tests cho Asset model
│   │   ├── scan.go                  # Entity Scan + ScanJob + ScanResult
│   │   ├── scan_test.go             # Unit tests cho Scan model
│   │   └── alert.go                 # Entity Alert cho notification system
│   ├── handler/
│   │   ├── router.go                # HTTP Router (Chi) + CORS middleware
│   │   ├── asset_handler.go         # REST handlers cho Asset CRUD
│   │   ├── asset_handler_test.go    # Handler tests với mock service
│   │   ├── scan_handler.go          # REST handlers cho Scan operations
│   │   ├── alert_handler.go         # REST handlers cho Alerts
│   │   └── health_handler.go        # Health check endpoint
│   ├── service/
│   │   ├── asset_service.go         # Business logic cho Asset management
│   │   ├── asset_service_test.go    # Service tests với mock storage
│   │   └── scan_service.go          # Business logic cho 7 scan engines
│   ├── repository/
│   │   ├── postgres/                # PostgreSQL implementations
│   │   │   ├── asset_postgres.go    # Asset CRUD với PostgreSQL
│   │   │   └── scan_postgres.go     # Scan CRUD với PostgreSQL
│   │   └── memory/                  # In-memory fallback (Day 1)
│   ├── scanner/                     # 7 Scanner Engines
│   │   ├── dns_scanner.go           # DNS record lookup (A, MX, NS, TXT)
│   │   ├── dns_scanner_test.go
│   │   ├── whois_scanner.go         # WHOIS registration lookup
│   │   ├── subdomain_scanner.go     # Subdomain enumeration
│   │   ├── ip_scanner.go            # IP Geolocation & ASN lookup
│   │   ├── ip_scanner_test.go
│   │   ├── port_scanner.go          # TCP Port Scan (localhost only)
│   │   ├── port_scanner_test.go
│   │   ├── ssl_scanner.go           # SSL/TLS Certificate analysis
│   │   ├── ssl_scanner_test.go
│   │   ├── tech_scanner.go          # Technology stack detection
│   │   └── tech_scanner_test.go
│   └── scheduler/
│       └── cron.go                  # Scheduled auto-scan scheduler
├── migrations/                      # SQL Migration files (auto-run on startup)
│   ├── 001_create_assets.up.sql
│   ├── 001_create_assets.down.sql
│   ├── 002_create_scan_tables.up.sql
│   ├── 002_create_scan_tables.down.sql
│   ├── 003_add_tags_to_assets.up.sql
│   ├── 003_add_tags_to_assets.down.sql
│   ├── 004_create_alerts.up.sql
│   └── 004_create_alerts.down.sql
├── frontend/                        # React + Vite Portal
│   ├── src/
│   │   ├── App.jsx                  # Main application component
│   │   ├── index.css                # Global styles (Light/Dark mode)
│   │   ├── components/              # Reusable UI components
│   │   ├── hooks/                   # Custom React hooks
│   │   └── utils/                   # Utility functions
│   ├── Dockerfile                   # Multi-stage build (Node → Nginx)
│   ├── nginx.conf                   # Nginx config (port 3000)
│   └── package.json
├── .github/
│   └── workflows/
│       ├── ci.yml                   # CI pipeline (tests + security scans)
│       └── deploy.yml               # Auto deploy on merge to 'deploy' branch
├── docker-compose.yml               # Full stack orchestration
├── Dockerfile                       # Backend multi-stage build
├── .env                             # Environment variables
├── go.mod                           # Go module definition
└── go.sum
```

---

## ⚙️ Biến Môi Trường

File `.env` cấu hình kết nối database và server:

```env
# Database Configuration
DB_HOST=localhost          # Host của PostgreSQL (dùng 'db' trong Docker Compose)
DB_PORT=5432               # Port của PostgreSQL
DB_USER=postgres           # Username
DB_PASSWORD=postgres       # Password
DB_NAME=mini_asm           # Tên database

# Server Configuration
SERVER_PORT=8080           # Port của backend API server
```

> **Lưu ý:** Khi chạy bằng Docker Compose, `DB_HOST` sẽ được override thành `db` (tên service trong compose).

---

## 🚀 Danh Sách API Endpoints

### 1. Quản lý Assets

| Method   | Endpoint                 | Mô tả                             |
| -------- | ------------------------ | --------------------------------- |
| `POST`   | `/assets`                | Tạo một asset mới                 |
| `GET`    | `/assets`                | Liệt kê assets (phân trang, lọc)  |
| `GET`    | `/assets/{id}`           | Lấy thông tin asset theo ID       |
| `POST`   | `/assets/batch`          | Tạo hàng loạt assets              |
| `DELETE` | `/assets/batch`          | Xóa hàng loạt assets              |
| `GET`    | `/assets/stats`          | Thống kê asset theo type & status |
| `GET`    | `/assets/count`          | Đếm số lượng asset theo bộ lọc    |
| `GET`    | `/assets/search`         | Tìm kiếm asset theo tên           |
| `PUT`    | `/assets/{id}/auto-scan` | Bật/tắt auto-scan cho asset       |

### 2. Quét An ninh EASM

| Method | Endpoint                  | Mô tả                                    |
| ------ | ------------------------- | ---------------------------------------- |
| `POST` | `/assets/{id}/scan`       | Kích hoạt scan job (async, 202 Accepted) |
| `GET`  | `/scan-jobs/{id}`         | Lấy trạng thái scan job                  |
| `GET`  | `/scan-jobs/{id}/results` | Lấy kết quả chi tiết scan job            |
| `GET`  | `/scan-jobs`              | Liệt kê tất cả scan jobs                 |
| `GET`  | `/assets/{id}/scans`      | Lịch sử scan của asset                   |
| `GET`  | `/assets/{id}/results`    | Tổng hợp toàn bộ kết quả scan            |
| `GET`  | `/assets/{id}/dns`        | Lấy DNS records của asset                |
| `GET`  | `/assets/{id}/whois`      | Lấy WHOIS data của asset                 |
| `GET`  | `/assets/{id}/subdomains` | Lấy subdomains của asset                 |

### 3. Alerts & Notifications

| Method | Endpoint            | Mô tả                        |
| ------ | ------------------- | ---------------------------- |
| `GET`  | `/alerts`           | Liệt kê tất cả alerts        |
| `GET`  | `/alerts/unread`    | Lấy alerts chưa đọc          |
| `PUT`  | `/alerts/read`      | Đánh dấu tất cả đã đọc       |
| `PUT`  | `/alerts/{id}/read` | Đánh dấu alert cụ thể đã đọc |

### 4. Health Check

| Method | Endpoint  | Mô tả                                    |
| ------ | --------- | ---------------------------------------- |
| `GET`  | `/health` | Kiểm tra sức khỏe hệ thống (DB & Uptime) |

### Scan Types Hỗ Trợ

| Scan Type   | Asset Type | Mô tả                                    | Loại      |
| ----------- | ---------- | ---------------------------------------- | --------- |
| `dns`       | domain     | DNS record lookup (A, MX, NS, TXT)       | Passive   |
| `whois`     | domain     | WHOIS registration information           | Passive   |
| `subdomain` | domain     | Subdomain enumeration                    | Passive   |
| `ip`        | ip         | IP Geolocation & ASN lookup              | Passive   |
| `port`      | ip         | TCP Port Scan (chỉ localhost/private IP) | ⚠️ Active |
| `ssl`       | domain     | SSL/TLS Certificate analysis             | ⚠️ Active |
| `tech`      | domain     | Technology stack detection               | Passive   |
| `all`       | domain/ip  | Chạy tất cả scans phù hợp với asset type | Mixed     |

> ⚠️ **Port Scanner Safety:** Chỉ cho phép scan `127.0.0.1`, `10.x.x.x`, `172.16-31.x.x`, `192.168.x.x`. Từ chối scan public IP.

---

## 🧪 Hướng Dẫn Test API

### Bài 1: Database Migration (PostgreSQL)

Tạo asset và kiểm tra dữ liệu persist sau khi restart server.

#### Linux / macOS (Bash)

```bash
# 1. Khởi động DB
docker compose up -d db

# 2. Chạy server (migrations tự động)
go run ./cmd/api/main.go

# 3. Tạo một asset
curl -X POST http://localhost:8080/assets \
  -H "Content-Type: application/json" \
  -d '{"name": "google.com", "type": "domain", "status": "active"}'

# 4. Kiểm tra danh sách
curl http://localhost:8080/assets

# 5. Restart server (Ctrl+C → chạy lại) → data vẫn tồn tại
go run ./cmd/api/main.go
curl http://localhost:8080/assets
```

#### Windows PowerShell

```powershell
# 1. Khởi động DB
docker compose up -d db

# 2. Chạy server
go run ./cmd/api/main.go

# 3. Tạo một asset
Invoke-RestMethod -Uri "http://localhost:8080/assets" `
  -Method POST `
  -ContentType "application/json" `
  -Body '{"name": "google.com", "type": "domain", "status": "active"}'

# 4. Kiểm tra danh sách
Invoke-RestMethod -Uri "http://localhost:8080/assets" -Method GET

# 5. Restart server → data vẫn tồn tại
```

#### Windows CMD

```cmd
:: 1. Tạo một asset
curl -X POST -H "Content-Type: application/json" -d "{\"name\":\"google.com\",\"type\":\"domain\",\"status\":\"active\"}" http://localhost:8080/assets

:: 2. Lấy danh sách
curl http://localhost:8080/assets
```

---

### Bài 2: Mở Rộng EASM Scanner API

Kích hoạt scan trên asset và lấy kết quả. Quy trình: **Tạo Asset → Start Scan → Poll Status → Get Results**.

#### Linux / macOS (Bash)

```bash
# --- Bước 1: Tạo Domain Asset ---
DOMAIN_ID=$(curl -s -X POST http://localhost:8080/assets \
  -H "Content-Type: application/json" \
  -d '{"name":"google.com","type":"domain"}' | jq -r '.id')
echo "Domain Asset ID: $DOMAIN_ID"

# --- Bước 2: Tạo IP Asset ---
IP_ID=$(curl -s -X POST http://localhost:8080/assets \
  -H "Content-Type: application/json" \
  -d '{"name":"127.0.0.1","type":"ip"}' | jq -r '.id')
echo "IP Asset ID: $IP_ID"

# --- Bước 3: Trigger Scans ---

# DNS Scan
DNS_JOB=$(curl -s -X POST "http://localhost:8080/assets/$DOMAIN_ID/scan" \
  -H "Content-Type: application/json" \
  -d '{"scan_type":"dns"}' | jq -r '.id')
echo "DNS Job: $DNS_JOB"

# Port Scan (chỉ localhost!)
PORT_JOB=$(curl -s -X POST "http://localhost:8080/assets/$IP_ID/scan" \
  -H "Content-Type: application/json" \
  -d '{"scan_type":"port"}' | jq -r '.id')
echo "Port Job: $PORT_JOB"

# SSL Scan
SSL_JOB=$(curl -s -X POST "http://localhost:8080/assets/$DOMAIN_ID/scan" \
  -H "Content-Type: application/json" \
  -d '{"scan_type":"ssl"}' | jq -r '.id')
echo "SSL Job: $SSL_JOB"

# Tech Scan
TECH_JOB=$(curl -s -X POST "http://localhost:8080/assets/$DOMAIN_ID/scan" \
  -H "Content-Type: application/json" \
  -d '{"scan_type":"tech"}' | jq -r '.id')
echo "Tech Job: $TECH_JOB"

# All Scans (chạy tất cả plugin phù hợp)
ALL_JOB=$(curl -s -X POST "http://localhost:8080/assets/$DOMAIN_ID/scan" \
  -H "Content-Type: application/json" \
  -d '{"scan_type":"all"}' | jq -r '.id')
echo "All Scan Job: $ALL_JOB"

# --- Bước 4: Kiểm tra trạng thái (đợi vài giây) ---
sleep 5
curl -s "http://localhost:8080/scan-jobs/$DNS_JOB" | jq

# --- Bước 5: Lấy kết quả ---
curl -s "http://localhost:8080/scan-jobs/$DNS_JOB/results" | jq

# --- Bước 6: Xem toàn bộ kết quả của asset ---
curl -s "http://localhost:8080/assets/$DOMAIN_ID/results" | jq

# --- Bước 7: Lịch sử scan ---
curl -s "http://localhost:8080/assets/$DOMAIN_ID/scans" | jq
```

#### Windows PowerShell

```powershell
# --- Bước 1: Tạo Domain Asset ---
$domain = Invoke-RestMethod -Uri "http://localhost:8080/assets" `
  -Method POST -ContentType "application/json" `
  -Body '{"name":"google.com","type":"domain"}'
$DOMAIN_ID = $domain.id
Write-Host "Domain Asset ID: $DOMAIN_ID"

# --- Bước 2: Tạo IP Asset ---
$ip = Invoke-RestMethod -Uri "http://localhost:8080/assets" `
  -Method POST -ContentType "application/json" `
  -Body '{"name":"127.0.0.1","type":"ip"}'
$IP_ID = $ip.id
Write-Host "IP Asset ID: $IP_ID"

# --- Bước 3: Trigger Scans ---

# DNS Scan
$dnsJob = Invoke-RestMethod -Uri "http://localhost:8080/assets/$DOMAIN_ID/scan" `
  -Method POST -ContentType "application/json" `
  -Body '{"scan_type":"dns"}'
Write-Host "DNS Job ID: $($dnsJob.id)"

# Port Scan (chỉ localhost!)
$portJob = Invoke-RestMethod -Uri "http://localhost:8080/assets/$IP_ID/scan" `
  -Method POST -ContentType "application/json" `
  -Body '{"scan_type":"port"}'
Write-Host "Port Job ID: $($portJob.id)"

# SSL Scan
$sslJob = Invoke-RestMethod -Uri "http://localhost:8080/assets/$DOMAIN_ID/scan" `
  -Method POST -ContentType "application/json" `
  -Body '{"scan_type":"ssl"}'
Write-Host "SSL Job ID: $($sslJob.id)"

# Tech Scan
$techJob = Invoke-RestMethod -Uri "http://localhost:8080/assets/$DOMAIN_ID/scan" `
  -Method POST -ContentType "application/json" `
  -Body '{"scan_type":"tech"}'
Write-Host "Tech Job ID: $($techJob.id)"

# All Scans
$allJob = Invoke-RestMethod -Uri "http://localhost:8080/assets/$DOMAIN_ID/scan" `
  -Method POST -ContentType "application/json" `
  -Body '{"scan_type":"all"}'
Write-Host "All Scan Job ID: $($allJob.id)"

# --- Bước 4: Kiểm tra trạng thái ---
Start-Sleep -Seconds 5
Invoke-RestMethod -Uri "http://localhost:8080/scan-jobs/$($dnsJob.id)" | ConvertTo-Json -Depth 10

# --- Bước 5: Lấy kết quả ---
Invoke-RestMethod -Uri "http://localhost:8080/scan-jobs/$($dnsJob.id)/results" | ConvertTo-Json -Depth 10

# --- Bước 6: Xem toàn bộ kết quả của asset ---
Invoke-RestMethod -Uri "http://localhost:8080/assets/$DOMAIN_ID/results" | ConvertTo-Json -Depth 10
```

#### Windows CMD

```cmd
:: Tạo asset
curl -X POST -H "Content-Type: application/json" -d "{\"name\":\"google.com\",\"type\":\"domain\"}" http://localhost:8080/assets

:: Trigger scan (thay YOUR_ASSET_ID bằng ID thực)
curl -X POST -H "Content-Type: application/json" -d "{\"scan_type\":\"dns\"}" http://localhost:8080/assets/YOUR_ASSET_ID/scan

:: Kiểm tra trạng thái (thay YOUR_JOB_ID)
curl http://localhost:8080/scan-jobs/YOUR_JOB_ID

:: Lấy kết quả
curl http://localhost:8080/scan-jobs/YOUR_JOB_ID/results
```

---

### Bài 3: Chạy Unit Tests

Bộ kiểm thử bao gồm: Model validation, Scanner tests, Handler tests (mock) và Service tests (mock).

#### Chạy toàn bộ tests

```bash
# Chạy tất cả tests (Linux/macOS/Windows)
go test -v ./...

# Chạy tests với coverage
go test -v -cover ./...

# Generate HTML coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

#### Chạy tests cụ thể

```bash
# Chỉ chạy model tests
go test -v ./internal/model/...

# Chỉ chạy scanner tests
go test -v ./internal/scanner/...

# Chỉ chạy handler tests
go test -v ./internal/handler/...

# Chỉ chạy service tests
go test -v ./internal/service/...

# Chạy test cụ thể theo tên
go test -v -run TestAssetValidation ./internal/model/...
```

---

### Bài 4: Tích Hợp Frontend React Portal

Giao diện Portal cung cấp trải nghiệm quản lý tài sản và quét an ninh trực quan.

#### Khởi chạy (Development)

```bash
# Terminal 1: Backend
cd homework-day1
go run ./cmd/api/main.go

# Terminal 2: Frontend
cd homework-day1/frontend
npm install
npm run dev

# Mở browser: http://localhost:5173
```

#### Các tính năng chính

1. **🌓 Light / Dark Mode Switcher** — Chuyển đổi theme qua nút bấm trên Header
2. **📊 Dashboard Overview** — Thống kê tổng quan với biểu đồ phân bố
3. **📋 Quản lý Assets** — CRUD, tìm kiếm real-time, lọc, phân trang, Batch Import CSV/JSON
4. **🔍 Trigger Scan & Live Progress** — Chọn plugin, polling status 2s, hiển thị kết quả chi tiết
5. **🔔 Alerts & Notifications** — Cảnh báo bảo mật tự động
6. **📥 Export Reports** — Xuất báo cáo scan results

---

## 📮 Hướng Dẫn Test với Postman

### Setup Postman Collection

1. Mở Postman → **New Collection** → đặt tên `EASM API`
2. Tạo biến Collection:
   - `base_url` = `http://localhost:8080`
3. Tạo các request theo hướng dẫn bên dưới

### Các Request cần tạo

#### 1. Health Check

- **Method:** `GET`
- **URL:** `{{base_url}}/health`

#### 2. Tạo Asset

- **Method:** `POST`
- **URL:** `{{base_url}}/assets`
- **Headers:** `Content-Type: application/json`
- **Body (raw JSON):**

```json
{
  "name": "google.com",
  "type": "domain",
  "status": "active"
}
```

#### 3. Liệt kê Assets

- **Method:** `GET`
- **URL:** `{{base_url}}/assets?page=1&limit=10`

#### 4. Lấy Asset theo ID

- **Method:** `GET`
- **URL:** `{{base_url}}/assets/{{asset_id}}`

#### 5. Batch Create

- **Method:** `POST`
- **URL:** `{{base_url}}/assets/batch`
- **Body (raw JSON):**

```json
{
  "assets": [
    { "name": "example.com", "type": "domain" },
    { "name": "192.168.1.1", "type": "ip" },
    { "name": "web-server", "type": "service" }
  ]
}
```

#### 6. Batch Delete

- **Method:** `DELETE`
- **URL:** `{{base_url}}/assets/batch?ids=id1,id2,id3`

#### 7. Trigger Scan

- **Method:** `POST`
- **URL:** `{{base_url}}/assets/{{asset_id}}/scan`
- **Body (raw JSON):**

```json
{
  "scan_type": "dns"
}
```

> Các giá trị `scan_type` hợp lệ: `dns`, `whois`, `subdomain`, `ip`, `port`, `ssl`, `tech`, `all`

#### 8. Kiểm tra Scan Status

- **Method:** `GET`
- **URL:** `{{base_url}}/scan-jobs/{{job_id}}`

#### 9. Lấy Scan Results

- **Method:** `GET`
- **URL:** `{{base_url}}/scan-jobs/{{job_id}}/results`

#### 10. Thống kê

- **Method:** `GET`
- **URL:** `{{base_url}}/assets/stats`

#### 11. Tìm kiếm

- **Method:** `GET`
- **URL:** `{{base_url}}/assets/search?q=google`

### Postman Test Script (Optional)

Thêm vào tab **Tests** của request tạo asset để tự động lưu ID:

```javascript
if (pm.response.code === 201) {
  var jsonData = pm.response.json();
  pm.collectionVariables.set("asset_id", jsonData.id);
}
```

---

## 🔄 CI/CD Pipeline

### GitHub Actions Workflows

#### 1. CI Pipeline (`.github/workflows/ci.yml`)

Chạy tự động khi push hoặc tạo PR vào `main`, `master`, `dev`, `homework-final`.

| Job                  | Tool       | Mô tả                                 |
| -------------------- | ---------- | ------------------------------------- |
| `Unit Tests`         | Go test    | Chạy toàn bộ unit tests               |
| `Gosec Scanner`      | Gosec      | Quét lỗ hổng bảo mật trong Go code    |
| `Gitleaks Scanner`   | Gitleaks   | Phát hiện secrets/API keys trong code |
| `Trivy Scanner`      | Trivy      | Quét CVE trong dependencies           |
| `TruffleHog Scanner` | TruffleHog | Quét deep commit history tìm secrets  |

#### 2. Auto Deploy (`.github/workflows/deploy.yml`)

Chạy tự động khi push vào nhánh `deploy`. SSH vào AWS EC2, pull code mới, rebuild Docker containers.

---

## ☁️ Deploy & Production

### Deploy lên Cloud VM (AWS EC2)

```bash
# 1. SSH vào VM
ssh -i your-key.pem ubuntu@your-vm-ip

# 2. Update system
sudo apt update && sudo apt upgrade -y

# 3. Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER

# 4. Clone repository
git clone https://github.com/thanhhieutiet/cmc-dev.git
cd cmc-dev/homework-day1

# 5. Deploy
docker compose up -d

# 6. Cấu hình firewall
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw allow 22/tcp
sudo ufw enable
```

### Domain & HTTPS

- **Domain:** `hieu-easm.ddns.net` (DDNS)
- **HTTPS:** Configured với auto certificate
- **Health check:** `curl -f https://hieu-easm.ddns.net/health`

---

## 📦 Kiểm Thử API Day 1 (In-Memory)

> Các tính năng Day 1 vẫn hoạt động đầy đủ trong phiên bản Day 3 (đã migrate sang PostgreSQL).

### Health Check

```bash
# Linux/macOS/CMD
curl http://localhost:8080/health
```

```powershell
# PowerShell
Invoke-RestMethod -Uri "http://localhost:8080/health"
```

### Batch Create

```bash
# Linux/macOS
curl -X POST http://localhost:8080/assets/batch \
  -H "Content-Type: application/json" \
  -d '{
    "assets": [
      {"name": "example.com", "type": "domain"},
      {"name": "192.168.1.1", "type": "ip"},
      {"name": "web-server", "type": "service"}
    ]
  }'
```

```powershell
# PowerShell
Invoke-RestMethod -Uri "http://localhost:8080/assets/batch" `
  -Method POST -ContentType "application/json" `
  -Body '{"assets": [{"name": "example.com", "type": "domain"}, {"name": "192.168.1.1", "type": "ip"}, {"name": "web-server", "type": "service"}]}'
```

### Statistics & Count

```bash
# Thống kê
curl http://localhost:8080/assets/stats

# Đếm tất cả
curl http://localhost:8080/assets/count

# Đếm theo filter
curl "http://localhost:8080/assets/count?type=domain&status=active"
```

```powershell
# PowerShell
Invoke-RestMethod -Uri "http://localhost:8080/assets/stats"
Invoke-RestMethod -Uri "http://localhost:8080/assets/count"
Invoke-RestMethod -Uri "http://localhost:8080/assets/count?type=domain&status=active"
```

### Batch Delete

```powershell
# PowerShell - Tạo assets rồi xóa
$res1 = Invoke-RestMethod -Uri "http://localhost:8080/assets" -Method POST -ContentType "application/json" -Body '{"name":"delete-me-1.com","type":"domain"}'
$res2 = Invoke-RestMethod -Uri "http://localhost:8080/assets" -Method POST -ContentType "application/json" -Body '{"name":"delete-me-2.com","type":"domain"}'
$ID1 = $res1.id
$ID2 = $res2.id

Invoke-RestMethod -Uri "http://localhost:8080/assets/batch?ids=$ID1,$ID2" -Method DELETE
```

### Concurrent-safe Create

```powershell
# PowerShell - Tạo 20 request song song
$jobs = @()
for ($i = 1; $i -le 20; $i++) {
    $jobs += Start-Job -ScriptBlock {
        param($n)
        Invoke-RestMethod -Uri "http://localhost:8080/assets" `
          -Method POST `
          -ContentType "application/json" `
          -Body "{`"name`":`"concurrent-$n.com`",`"type`":`"domain`"}"
    } -ArgumentList $i
}
$jobs | Wait-Job | Out-Null

# Kiểm tra count
Invoke-RestMethod -Uri "http://localhost:8080/assets/count"
```

### Pagination & Search

```bash
# Phân trang
curl "http://localhost:8080/assets?page=1&limit=5&type=domain&status=active"

# Tìm kiếm
curl "http://localhost:8080/assets/search?q=example"
```

```powershell
# PowerShell
Invoke-RestMethod -Uri "http://localhost:8080/assets?page=1&limit=5&type=domain&status=active"
Invoke-RestMethod -Uri "http://localhost:8080/assets/search?q=example"
```

---

## 📝 License

This project is part of the CMC Internship Training Program.

**Author:** Tiết Thanh Minh Hiếu
