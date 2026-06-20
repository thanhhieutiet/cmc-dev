# 📦 EASM Digital Asset Management & Security Scanner Portal

Dự án này là cổng quản lý tài sản số (Asset Management) và quét an ninh bề mặt tấn công (EASM Scanner) được xây dựng theo mô hình **Clean Architecture** sử dụng **Go**, **PostgreSQL** và **React Portal (Vite)**.

Đây là phiên bản tích hợp đầy đủ của cả **Homework Day 1** (phiên bản In-memory ban đầu) và nâng cấp hoàn chỉnh cho **Homework Day 3 (Final)** (kết nối PostgreSQL, tự động chạy Migrations, mở rộng 7 Scanner Engines, viết bộ Unit Tests mới, và tích hợp giao diện Web Portal hỗ trợ Light/Dark mode).

---

## 📋 Mục Lục Tổng Quan

1. [Yêu Cầu Hệ Thống & Khởi Chạy Nhanh](#-yêu-cầu-hệ-thống--khởi-chạy-nhanh)
2. [Cấu Trúc Thư Mục](#-cấu-trúc-thư-mục)
3. [Danh Sách API Endpoints](#-danh-sách-api-endpoints)
4. [PHẦN I: HƯỚNG DẪN KIỂM THỬ DAY 3 (FINAL)](#phần-i-hướng-dẫn-kiểm-thử-day-3-final)
   - [Bài 1: Database Migration (PostgreSQL)](#bài-1-database-migration-postgresql)
   - [Bài 2: Mở Rộng EASM Scanner API](#bài-2-mở-rộng-easm-scanner-api)
   - [Bài 3: Chạy Unit Tests mới](#bài-3-chạy-unit-tests-mới)
   - [Bài 4: Tích Hợp Frontend React Portal](#bài-4-tích-hợp-frontend-react-portal)
5. [PHẦN II: HƯỚNG DẪN KIỂM THỬ DAY 1 (IN-MEMORY)](#phần-ii-hướng-dẫn-kiểm-thử-day-1-in-memory)
   - [Bài 5: Health Check (Day 1)](#bài-5-health-check-day-1)
   - [Bài 2: Batch Create (Day 1)](#bài-2-batch-create-day-1)
   - [Bài 1.1: Statistics (Day 1)](#bài-11-statistics-day-1)
   - [Bài 1.2: Count (Day 1)](#bài-12-count-day-1)
   - [Bài 3: Batch Delete (Day 1)](#bài-3-batch-delete-day-1)
   - [Bài 4: Concurrent-safe Create (Day 1)](#bài-4-concurrent-safe-create-day-1)
   - [Bài 6: Pagination (Day 1)](#bài-6-pagination-day-1)
   - [Bài 7: Search (Day 1)](#bài-7-search-day-1)

---

## 🔧 Yêu Cầu Hệ Thống & Khởi Chạy Nhanh

- **Go** phiên bản 1.21 trở lên (khuyến nghị 1.26+)
- **Docker Desktop** (để chạy database)
- **Node.js** phiên bản 18+ và **npm** (để chạy frontend)
- **curl**, **PowerShell** hoặc **Postman** để test API

### 1. Khởi chạy Database
Khởi động container PostgreSQL từ thư mục root của dự án bằng lệnh:
```bash
docker compose up -d
```

### 2. Khởi chạy Backend Server
```bash
go mod tidy
go run ./cmd/api/main.go
```
API server chạy tại cổng `http://localhost:8080`.

### 3. Khởi chạy Frontend React Portal
```bash
cd frontend
npm install
npm run dev
```
Giao diện portal chạy tại: **http://localhost:5173**

---

## 📁 Cấu Trúc Thư Mục

Dự án được tổ chức theo mô hình **Clean Architecture**:
```
homework-day1/
├── cmd/api/main.go              # Điểm khởi chạy ứng dụng (chạy DB migrations & init server)
├── internal/
│   ├── config/config.go         # Đọc cấu hình từ .env và biến môi trường
│   ├── domain/                  # Các thực thể nghiệp vụ & Định nghĩa Interfaces
│   │   ├── asset.go
│   │   └── scan.go
│   ├── usecase/                 # Quy trình nghiệp vụ chính (AssetUsecase & ScanUsecase)
│   │   ├── asset_usecase.go
│   │   └── scan_usecase.go
│   ├── repository/              # Lưu trữ dữ liệu (PostgreSQL & In-Memory fallback)
│   │   ├── postgres/
│   │   │   ├── asset_postgres.go
│   │   │   └── scan_postgres.go
│   │   └── memory/
│   └── delivery/http/           # Lớp phân phối HTTP (Router, Handlers & CORS)
│   │   ├── router.go
│   │   ├── asset_handler.go
│   │   └── scan_handler.go
│   └── scanner/                 # Các scanner engine quét an ninh EASM
│       ├── dns_scanner.go
│       ├── whois_scanner.go
│       ├── subdomain_scanner.go
│       ├── ip_scanner.go
│       ├── port_scanner.go
│       ├── ssl_scanner.go
│       └── tech_scanner.go
├── migrations/                  # Các file SQL Migration chạy tự động
├── docker-compose.yml
├── .env
├── go.mod
└── frontend/                    # Giao diện Web Portal (React + Vite)
```

---

## 🚀 Danh Sách API Endpoints

### 1. Quản lý Assets (Day 1 & Day 3)
- `POST   /assets`                - Tạo một asset mới.
- `GET    /assets/{id}`           - Lấy thông tin asset theo ID.
- `GET    /assets`                - Liệt kê assets có phân trang và bộ lọc.
- `GET    /assets/search`         - Tìm kiếm asset theo tên.
- `POST   /assets/batch`          - Tạo hàng loạt assets.
- `DELETE /assets/batch`          - Xóa hàng loạt assets (nhận IDs qua URL Query).
- `GET    /assets/stats`          - Thống kê số lượng asset theo Type & Status.
- `GET    /assets/count`          - Đếm số lượng asset theo bộ lọc.
- `GET    /health`                - Kiểm tra sức khỏe hệ thống (Database status & Uptime).

### 2. Quét An ninh EASM (Day 3 nâng cấp)
- `POST   /assets/{id}/scan`      - Trigger scan job cho asset (chạy ngầm).
- `GET    /scan-jobs/{id}`        - Lấy trạng thái của scan job.
- `GET    /scan-jobs/{id}/results`- Lấy kết quả chi tiết của scan job.
- `GET    /assets/{id}/scans`     - Lịch sử các scan jobs của asset.
- `GET    /assets/{id}/results`   - Tổng hợp toàn bộ kết quả scan của asset.
- `GET    /assets/{id}/dns`       - Lấy các bản ghi DNS phát hiện được của asset.
- `GET    /assets/{id}/whois`     - Lấy dữ liệu WHOIS phát hiện được của asset.
- `GET    /assets/{id}/subdomains`- Lấy danh sách subdomains phát hiện được của asset.

---

## PHẦN I: HƯỚNG DẪN KIỂM THỬ DAY 3 (FINAL)

### Bài 1: Database Migration (PostgreSQL)
Lưu trữ toàn bộ thông tin assets và kết quả scan vào PostgreSQL. Dữ liệu sẽ tồn tại bền vững kể cả khi khởi động lại server.

#### Hướng dẫn kiểm thử:
**Linux (bash):**
```bash
# 1. Tạo một asset
curl -X POST -H "Content-Type: application/json" -d '{"name": "google.com", "type": "domain", "status": "active"}' http://localhost:8080/assets

# 2. Lấy danh sách để kiểm tra
curl http://localhost:8080/assets
```

**Windows CMD:**
```cmd
:: 1. Tạo một asset
curl -X POST -H "Content-Type: application/json" -d "{\"name\":\"google.com\",\"type\":\"domain\",\"status\":\"active\"}" http://localhost:8080/assets

:: 2. Lấy danh sách
curl http://localhost:8080/assets
```

**Windows PowerShell:**
```powershell
# 1. Tạo một asset
Invoke-RestMethod -Uri "http://localhost:8080/assets" -Method POST -ContentType "application/json" -Body '{"name": "google.com", "type": "domain", "status": "active"}'

# 2. Lấy danh sách
Invoke-RestMethod -Uri "http://localhost:8080/assets" -Method GET
```

---

### Bài 2: Mở Rộng EASM Scanner API
Kích hoạt các tác vụ quét an ninh trên tài sản số. Đối với domain, hệ thống hỗ trợ quét: `dns`, `whois`, `subdomain`, `ssl`, `tech` hoặc `all`. Đối với IP, hỗ trợ: `ip` (ASN & Geolocation), `port` (Active TCP Port Scan) hoặc `all`.
*Chú ý: Port Scanner có cơ chế kiểm tra an toàn, chỉ cho phép quét localhost & dải IP Private (`127.0.0.1`, `10.x.x.x`, `172.16-31.x.x`, `192.168.x.x`).*

Giả sử ID của asset nhận được là `ae42152f-4e97-4556-9c39-9c8b300d4ba0`.

#### A. Trigger Scan Job (Chạy không đồng bộ, trả về 202 Accepted)
**Linux (bash):**
```bash
curl -X POST -H "Content-Type: application/json" -d '{"scan_type": "all"}' http://localhost:8080/assets/ae42152f-4e97-4556-9c39-9c8b300d4ba0/scan
```
**Windows CMD:**
```cmd
curl -X POST -H "Content-Type: application/json" -d "{\"scan_type\":\"all\"}" http://localhost:8080/assets/ae42152f-4e97-4556-9c39-9c8b300d4ba0/scan
```
**Windows PowerShell:**
```powershell
Invoke-RestMethod -Uri "http://localhost:8080/assets/ae42152f-4e97-4556-9c39-9c8b300d4ba0/scan" -Method POST -ContentType "application/json" -Body '{"scan_type": "all"}'
```
*Phản hồi mẫu sẽ trả về thông tin job bao gồm `id` của job.*

#### B. Kiểm tra trạng thái Job (Poll Status)
Giả sử ID của Scan Job nhận được là `924bb54e-0dbf-4359-be8d-f00607387f6c`.
**Linux (bash) & Windows CMD:**
```bash
curl http://localhost:8080/scan-jobs/924bb54e-0dbf-4359-be8d-f00607387f6c
```
**Windows PowerShell:**
```powershell
Invoke-RestMethod -Uri "http://localhost:8080/scan-jobs/924bb54e-0dbf-4359-be8d-f00607387f6c" -Method GET
```

#### C. Lấy kết quả Quét (khi Status đã chuyển sang `completed`)
**Linux (bash) & Windows CMD:**
```bash
curl http://localhost:8080/scan-jobs/924bb54e-0dbf-4359-be8d-f00607387f6c/results
```
**Windows PowerShell:**
```powershell
Invoke-RestMethod -Uri "http://localhost:8080/scan-jobs/924bb54e-0dbf-4359-be8d-f00607387f6c/results" -Method GET
```

---

### Bài 3: Chạy Unit Tests mới
Bộ kiểm thử chạy tự động xác thực các chức năng nghiệp vụ của Model Validation và các công cụ Scanner.

Để thực thi và xem tỷ lệ bao phủ mã nguồn (code coverage):
```bash
go test -v -cover ./...
```

---

### Bài 4: Tích Hợp Frontend React Portal
Giao diện Portal được phát triển giúp thao tác và trực quan hóa toàn bộ dữ liệu an ninh bề mặt tấn công.

#### Các tính năng chính:
1. **Light / Dark Mode Switcher:** Thay đổi hệ màu qua CSS variables bằng nút bấm trên thanh Header.
2. **Dashboard Overview:** Tổng hợp số lượng asset và trạng thái scan job gần đây bằng các thẻ panel. Tích hợp biểu đồ phân bố SVG sinh động.
3. **Quản lý Assets:** Hỗ trợ tạo mới, tìm kiếm thời gian thực, lọc theo loại/trạng thái và phân trang. Cung cấp tính năng Batch Import qua CSV/JSON. Hiển thị ID trực tiếp kèm nút Copy nhanh tại chỗ.
4. **Trigger Scan & Live Progress:** Khi bấm nút "Scan" ở bất kỳ asset nào, giao diện sẽ kích hoạt modal tùy chọn Plugin và thực hiện Polling trạng thái mỗi 2 giây. Khi hoàn tất, một báo cáo kết quả chi tiết (chấm điểm Grade cho SSL, sơ đồ Geolocation, open ports list, tech stack badges) sẽ tự động hiển thị trực quan.

---

## PHẦN II: HƯỚNG DẪN KIỂM THỬ DAY 1 (IN-MEMORY)

*Dưới đây là tài liệu kiểm thử cho các tính năng in-memory ban đầu của Homework Day 1 để đối chiếu.*

### Bài 5: Health Check (Day 1)
**Linux/macOS (Bash) & Windows CMD:**
```bash
curl http://localhost:8080/health
```
**Windows PowerShell:**
```powershell
Invoke-RestMethod -Uri "http://localhost:8080/health"
```

---

### Bài 2: Batch Create (Day 1)
#### 1. Tạo thành công (tạo 3 assets)
**Linux/macOS (Bash):**
```bash
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
**Windows CMD:**
```cmd
curl.exe -X POST http://localhost:8080/assets/batch -H "Content-Type: application/json" -d "{\"assets\": [{\"name\": \"example.com\", \"type\": \"domain\"}, {\"name\": \"192.168.1.1\", \"type\": \"ip\"}, {\"name\": \"web-server\", \"type\": \"service\"}]}"
```
**Windows PowerShell:**
```powershell
Invoke-RestMethod -Uri "http://localhost:8080/assets/batch" -Method POST -ContentType "application/json" -Body '{"assets": [{"name": "example.com", "type": "domain"}, {"name": "192.168.1.1", "type": "ip"}, {"name": "web-server", "type": "service"}]}'
```

---

### Bài 1.1: Statistics (Day 1)
**Linux/macOS (Bash) & Windows CMD:**
```bash
curl http://localhost:8080/assets/stats
```
**Windows PowerShell:**
```powershell
Invoke-RestMethod -Uri "http://localhost:8080/assets/stats"
```

---

### Bài 1.2: Count (Day 1)
#### 1. Đếm tất cả assets
**Linux/macOS (Bash) & Windows CMD:**
```bash
curl http://localhost:8080/assets/count
```
**Windows PowerShell:**
```powershell
Invoke-RestMethod -Uri "http://localhost:8080/assets/count"
```

#### 2. Đếm theo type VÀ status
**Linux/macOS (Bash) & Windows CMD:**
```bash
curl "http://localhost:8080/assets/count?type=domain&status=active"
```
**Windows PowerShell:**
```powershell
Invoke-RestMethod -Uri "http://localhost:8080/assets/count?type=domain&status=active"
```

---

### Bài 3: Batch Delete (Day 1)
**Windows PowerShell:**
```powershell
# Bước 1: Tạo assets và tự động lưu ID vào biến
$res1 = Invoke-RestMethod -Uri "http://localhost:8080/assets" -Method POST -ContentType "application/json" -Body '{"name":"delete-me-1.com","type":"domain"}'
$res2 = Invoke-RestMethod -Uri "http://localhost:8080/assets" -Method POST -ContentType "application/json" -Body '{"name":"delete-me-2.com","type":"domain"}'
$ID1 = $res1.id
$ID2 = $res2.id

# Bước 2: Thực hiện Batch Delete
Invoke-RestMethod -Uri "http://localhost:8080/assets/batch?ids=$ID1,$ID2,fake-uuid-123" -Method DELETE
```

---

### Bài 4: Concurrent-safe Create (Day 1)
**Windows PowerShell:**
```powershell
# Tạo 20 request song song
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

---

### Bài 6: Pagination (Day 1)
**Windows PowerShell:**
```powershell
Invoke-RestMethod -Uri "http://localhost:8080/assets?page=1&limit=5&type=domain&status=active"
```

---

### Bài 7: Search (Day 1)
**Windows PowerShell:**
```powershell
Invoke-RestMethod -Uri "http://localhost:8080/assets/search?q=example"
```

---

## 📮 Hướng Dẫn Test Với Postman (Day 1 & Day 3)

Bộ sưu tập Postman hỗ trợ toàn bộ các request từ Bài 1 đến Bài 4. Cấu hình kiểm thử:
1. Nhập file Postman Collection hoặc tạo các Request tương ứng với URL `http://localhost:8080`.
2. Tạo Asset mới: `POST http://localhost:8080/assets` với Header `Content-Type: application/json` và Body dạng raw JSON.
3. Batch Create: `POST http://localhost:8080/assets/batch` truyền danh sách assets.
4. Batch Delete: `DELETE http://localhost:8080/assets/batch?ids=id1,id2`.
5. Thực hiện Quét: `POST http://localhost:8080/assets/{id}/scan` kèm body chọn plugin.
6. Lấy kết quả: `GET http://localhost:8080/scan-jobs/{id}/results`.
