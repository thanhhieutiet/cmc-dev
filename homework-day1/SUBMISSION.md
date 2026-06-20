# 📦 EASM Digital Asset Management & Security Scanner - Final Submission

Bản báo cáo kết quả hoàn thành bài tập nâng cấp **Day 3 (Final) - Phase 1**. Dự án được tiếp tục phát triển trực tiếp trên cấu trúc Clean Architecture của Day 1 cũ, nâng cấp tầng lưu trữ sang **PostgreSQL** và tích hợp thêm hệ thống **EASM Scanning (DNS, WHOIS, Subdomain, GeoIP, Port scan, SSL TLS Audit, và Tech stack finder)** kèm giao diện **React Portal** hỗ trợ Light/Dark mode.

---

## 📋 Mục Lục
1. [Yêu cầu hệ thống](#1-yêu-cầu-hệ-thống)
2. [Hướng dẫn khởi chạy chi tiết](#2-hướng-dẫn-khởi-chạy-chi-tiết)
   - [Bước 1: Khởi chạy PostgreSQL bằng Docker](#bước-1-khởi-chạy-postgresql-bằng-docker)
   - [Bước 2: Khởi chạy Backend Go API](#bước-2-khởi-chạy-backend-go-api)
   - [Bước 3: Khởi chạy Frontend React Portal](#bước-3-khởi-chạy-frontend-react-portal)
3. [Cơ chế tự động chạy Database Migration](#3-cơ-chế-tự-động-chạy-database-migration)
4. [Danh sách API Endpoints mở rộng](#4-danh-sách-api-endpoints-mở-rộng)
5. [Hướng dẫn kiểm thử bằng Powershell & CMD (Windows & Linux)](#5-hướng-dẫn-kiểm-thử-bằng-powershell--cmd-windows--linux)
6. [Kết quả chạy Unit Tests](#6-kết-quả-chạy-unit-tests)

---

## 1. Yêu cầu hệ thống
- **Docker Desktop** (để chạy PostgreSQL container)
- **Go** (phiên bản 1.21 trở lên, khuyến nghị 1.26+)
- **Node.js** (phiên bản 18+ và npm)
- **PowerShell** hoặc **bash** để kiểm thử

---

## 2. Hướng dẫn khởi chạy chi tiết

Di chuyển vào thư mục dự án:
```bash
cd c:\STUDY\ThucTap\CMC_Training\DEV\homework-day1
```

### Bước 1: Khởi chạy PostgreSQL bằng Docker
Trong thư mục root `homework-day1/` có chứa file `docker-compose.yml`. Khởi động database bằng lệnh:
```bash
docker compose up -d
```
Database sẽ được khởi chạy ở cổng mặc định `5432` với tên DB là `mini_asm`.

### Bước 2: Khởi chạy Backend Go API
Tải dependencies và khởi chạy API server:
```bash
go mod tidy
go run ./cmd/api/main.go
```
Khi khởi chạy, server sẽ tự động thực thi các file migrations nằm trong thư mục `migrations/` để thiết lập cơ sở dữ liệu và in ra danh sách các API endpoints trên cổng `http://localhost:8080`.

### Bước 3: Khởi chạy Frontend React Portal
Mở một terminal mới, chuyển vào thư mục `frontend/`, tải dependencies và chạy dev server:
```bash
cd frontend
npm install
npm run dev
```
Trang giao diện Web Portal sẽ chạy tại: **http://localhost:5173**

---

## 3. Cơ chế tự động chạy Database Migration
Tầng lưu trữ được xây dựng với cơ chế **Auto-Migration** nằm trong `cmd/api/main.go`. Quy trình hoạt động:
1. Tạo bảng `schema_migrations` nếu chưa tồn tại để quản lý lịch sử các file SQL migration đã áp dụng.
2. Quét toàn bộ các file `.up.sql` trong thư mục `migrations/` theo thứ tự chữ cái (ví dụ: `001_create_assets.up.sql`, `002_create_scan_tables.up.sql`).
3. Với mỗi file chưa được áp dụng, server mở một transaction (`db.Begin()`), thực thi nội dung SQL, ghi log phiên bản vào bảng `schema_migrations` và commit. Điều này giúp khởi động server an toàn và nhất quán dữ liệu ở mọi môi trường dev/staging.

---

## 4. Danh sách API Endpoints mở rộng

### Nhóm Quản lý Asset (Day 1)
- `POST   /assets`                - Đăng ký asset đơn lẻ.
- `GET    /assets/{id}`           - Lấy thông tin chi tiết asset.
- `GET    /assets`                - Liệt kê assets có phân trang (`page`, `limit`) và bộ lọc (`type`, `status`).
- `GET    /assets/search`         - Tìm kiếm asset theo tên (`q=query`).
- `POST   /assets/batch`          - Tạo hàng loạt assets.
- `DELETE /assets/batch`          - Xóa hàng loạt assets.
- `GET    /assets/stats`          - Thống kê tổng quan số lượng asset.
- `GET    /assets/count`          - Đếm số asset theo bộ lọc.

### Nhóm Quét An ninh EASM (Day 3 nâng cấp)
- `POST   /assets/{id}/scan`      - Kích hoạt lượt quét mới (quét ngầm Goroutine, trả về `202 Accepted` kèm ID scan job).
- `GET    /scan-jobs/{id}`        - Kiểm tra trạng thái tiến trình quét (`pending`, `running`, `completed`, `failed`, `partial`).
- `GET    /scan-jobs/{id}/results`- Nhận kết quả chi tiết của lượt quét tương ứng.
- `GET    /assets/{id}/scans`     - Lịch sử các lượt quét của asset.
- `GET    /assets/{id}/results`   - Tổng hợp kết quả của tất cả các lần quét gần nhất trên asset đó.
- `GET    /assets/{id}/dns`       - Bản ghi DNS đã phát hiện.
- `GET    /assets/{id}/whois`     - Dữ liệu WHOIS đã phân tích.
- `GET    /assets/{id}/subdomains`- Các subdomain đã phát hiện.

---

## 5. Hướng dẫn kiểm thử bằng Powershell & CMD (Windows & Linux)

Dưới đây là so sánh các lệnh gọi cURL trên Windows (PowerShell và CMD) và Linux để tránh các lỗi phân tách ký tự JSON.

### A. Đăng ký Asset mới
**PowerShell (Windows):**
```powershell
Invoke-RestMethod -Uri "http://localhost:8080/assets" -Method POST -ContentType "application/json" -Body '{"name": "google.com", "type": "domain", "status": "active"}'
```
**CMD (Windows):**
```cmd
curl -X POST -H "Content-Type: application/json" -d "{\"name\":\"google.com\",\"type\":\"domain\",\"status\":\"active\"}" http://localhost:8080/assets
```
**Bash (Linux):**
```bash
curl -X POST -H "Content-Type: application/json" -d '{"name": "google.com", "type": "domain", "status": "active"}' http://localhost:8080/assets
```

### B. Kích hoạt lượt quét DNS & WHOIS
Giả sử ID của asset trả về từ lệnh trên là `ae42152f-4e97-4556-9c39-9c8b300d4ba0`.

**PowerShell (Windows):**
```powershell
Invoke-RestMethod -Uri "http://localhost:8080/assets/ae42152f-4e97-4556-9c39-9c8b300d4ba0/scan" -Method POST -ContentType "application/json" -Body '{"scan_type": "all"}'
```
Lệnh trên sẽ trả về mã `202 Accepted` kèm theo thông tin của Scan Job.
Ví dụ phản hồi:
```json
{
  "id": "924bb54e-0dbf-4359-be8d-f00607387f6c",
  "asset_id": "ae42152f-4e97-4556-9c39-9c8b300d4ba0",
  "scan_type": "all",
  "status": "pending",
  "started_at": "2026-06-18T16:52:00+07:00"
}
```

**CMD (Windows):**
```cmd
curl -X POST -H "Content-Type: application/json" -d "{\"scan_type\":\"all\"}" http://localhost:8080/assets/ae42152f-4e97-4556-9c39-9c8b300d4ba0/scan
```

### C. Xem trạng thái tiến độ Quét
Sử dụng ID của Scan Job để kiểm tra:

**PowerShell (Windows):**
```powershell
Invoke-RestMethod -Uri "http://localhost:8080/scan-jobs/924bb54e-0dbf-4359-be8d-f00607387f6c" -Method GET
```
**CMD (Windows) & Linux:**
```bash
curl http://localhost:8080/scan-jobs/924bb54e-0dbf-4359-be8d-f00607387f6c
```

### D. Xem kết quả chi tiết sau khi Job chuyển sang `completed`

**PowerShell (Windows):**
```powershell
Invoke-RestMethod -Uri "http://localhost:8080/scan-jobs/924bb54e-0dbf-4359-be8d-f00607387f6c/results" -Method GET
```
**CMD (Windows) & Linux:**
```bash
curl http://localhost:8080/scan-jobs/924bb54e-0dbf-4359-be8d-f00607387f6c/results
```

---

## 6. Kết quả chạy Unit Tests
Chúng ta có thể chạy kiểm thử tự động toàn bộ codebase Go bằng lệnh:
```bash
go test -v ./...
```
Tất cả các ca kiểm thử liên quan đến Model Validation và Scanner Engine (tra cứu, phân tích, trích xuất version) đều vượt qua thành công:
```
ok      homework-day1/internal/handler                 0.626s
ok      homework-day1/internal/model                   0.395s
ok      homework-day1/internal/repository/memory       0.379s
ok      homework-day1/internal/scanner                 8.328s
ok      homework-day1/internal/service                 0.521s
```
