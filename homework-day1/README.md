# 📦 Asset Management API - Homework Day 1

REST API quản lý tài sản số (Asset Management) được xây dựng theo mô hình **Clean Architecture** sử dụng **Go** và **In-memory Storage**.

---

## 📋 Mục Lục

- [Yêu cầu hệ thống](#-yêu-cầu-hệ-thống)
- [Cài đặt](#-cài-đặt)
- [Cách chạy](#-cách-chạy)
- [Cấu trúc thư mục](#-cấu-trúc-thư-mục)
- [Danh sách API Endpoints](#-danh-sách-api-endpoints)
- [Hướng dẫn test bằng curl](#-hướng-dẫn-test-bằng-curl)
  - [Bài 5: Health Check](#bài-5-health-check)
  - [Bài 2: Batch Create Assets](#bài-2-batch-create-assets)
  - [Bài 1.1: Get Assets Statistics](#bài-11-get-assets-statistics)
  - [Bài 1.2: Count Assets by Filter](#bài-12-count-assets-by-filter)
  - [Bài 3: Batch Delete Assets](#bài-3-batch-delete-assets)
  - [Bài 4: Concurrent-safe Create](#bài-4-concurrent-safe-create)
  - [Bài 6: Pagination & Filtering (Bonus)](#bài-6-pagination--filtering-bonus)
  - [Bài 7: Search by Name (Bonus)](#bài-7-search-by-name-bonus)
- [Hướng dẫn test bằng Postman](#-hướng-dẫn-test-bằng-postman)
- [Chạy Unit Tests](#-chạy-unit-tests)

---

## 🔧 Yêu cầu hệ thống

- **Go** phiên bản 1.21 trở lên (khuyến nghị 1.26+)
- **curl** hoặc **Postman** để test API

Kiểm tra Go đã cài đặt:

```bash
go version
# Output mong đợi: go version go1.26.x ...
```

---

## 📥 Cài đặt

### 1. Clone/Di chuyển vào thư mục project

```bash
cd homework-day1
```

### 2. Tải dependencies

```bash
go mod tidy
```

Lệnh này sẽ tự động tải về các thư viện cần thiết:
- `github.com/go-chi/chi/v5` — HTTP Router
- `github.com/google/uuid` — Sinh UUID ngẫu nhiên

---

## 🚀 Cách chạy

### Khởi chạy server

```bash
go run ./cmd/api/main.go
```

Server sẽ chạy tại: **http://localhost:8080**

Khi khởi chạy thành công, bạn sẽ thấy output như sau:

```
======================================
  Asset Management API
  Clean Architecture + In-Memory
  Server started on http://localhost:8080
======================================

Available endpoints:
  GET    /health          - Health check (Bài 5)
  GET    /assets/stats    - Statistics (Bài 1.1)
  GET    /assets/count    - Count with filters (Bài 1.2)
  POST   /assets/batch    - Batch create (Bài 2)
  DELETE /assets/batch    - Batch delete (Bài 3)
  POST   /assets          - Create single asset
  GET    /assets/{id}     - Get asset by ID
  GET    /assets          - List with pagination (Bài 6)
  GET    /assets/search   - Search by name (Bài 7)
```

### Dừng server

Nhấn `Ctrl + C` trong terminal đang chạy server.

---

## 📁 Cấu trúc thư mục

```
homework-day1/
├── cmd/
│   └── api/
│       └── main.go                          # Điểm bắt đầu: khởi tạo và nối các layer
├── internal/
│   ├── domain/
│   │   └── asset.go                         # Lớp 1 (Entities): Struct, Interface, DTO, Errors
│   ├── usecase/
│   │   └── asset_usecase.go                 # Lớp 2 (Use Cases): Logic nghiệp vụ
│   ├── repository/
│   │   └── memory/
│   │       ├── asset_memory.go              # Lớp 3 (Repository): In-memory storage + sync.RWMutex
│   │       └── asset_memory_test.go         # Unit tests cho concurrent safety
│   └── delivery/
│       └── http/
│           ├── asset_handler.go             # Lớp 3 (Delivery): HTTP handlers
│           ├── health_handler.go            # Health check handler
│           └── router.go                    # Đăng ký routes với chi router
├── go.mod
├── go.sum
└── README.md                                # File hướng dẫn này
```

### Giải thích Clean Architecture

```
Request → [Handler/Delivery] → [Usecase] → [Repository] → In-Memory Storage
                ↑                   ↑            ↑
             Lớp 3              Lớp 2         Lớp 3
         (HTTP adapter)    (Business Logic)  (Data adapter)
                                  ↓
                            [Domain/Entity]
                               Lớp 1
                        (Struct, Interface)
```

- **Domain** (`internal/domain/`): Định nghĩa thực thể `Asset`, các interface `AssetRepository` và `AssetUsecase`, các DTO request/response, và các lỗi nghiệp vụ.
- **Usecase** (`internal/usecase/`): Chứa logic nghiệp vụ — validate all-or-nothing, phân trang, đếm, thống kê.
- **Repository** (`internal/repository/memory/`): Lưu trữ dữ liệu bằng `map[string]*Asset` trong bộ nhớ, bảo vệ bằng `sync.RWMutex`.
- **Delivery** (`internal/delivery/http/`): Nhận HTTP request, gọi Usecase, trả HTTP response.

---

## 📡 Danh sách API Endpoints

| Method   | Endpoint          | Mô tả                          | Bài tập |
|----------|-------------------|---------------------------------|---------|
| `GET`    | `/health`         | Health check                    | Bài 5   |
| `GET`    | `/assets/stats`   | Thống kê tổng quan              | Bài 1.1 |
| `GET`    | `/assets/count`   | Đếm theo bộ lọc                 | Bài 1.2 |
| `POST`   | `/assets/batch`   | Tạo nhiều assets (all-or-nothing)| Bài 2   |
| `DELETE` | `/assets/batch`   | Xóa nhiều assets                | Bài 3   |
| `POST`   | `/assets`         | Tạo 1 asset                     | CRUD    |
| `GET`    | `/assets/{id}`    | Lấy 1 asset theo ID             | CRUD    |
| `GET`    | `/assets`         | Danh sách có phân trang + lọc   | Bài 6   |
| `GET`    | `/assets/search`  | Tìm kiếm theo tên               | Bài 7   |

### Asset Model

```json
{
  "id": "uuid-string",
  "name": "example.com",
  "type": "domain",         // domain | ip | service
  "status": "active",       // active | inactive (mặc định: active)
  "created_at": "2026-06-12T10:00:00Z"
}
```

---

## 🧪 Hướng dẫn test bằng curl

> **Lưu ý:** Hãy đảm bảo server đang chạy trước khi test (`go run ./cmd/api/main.go`).

---

### Bài 5: Health Check

**Linux/macOS (Bash) & Windows CMD:**
```bash
curl http://localhost:8080/health
```

**Windows PowerShell:**
```powershell
Invoke-RestMethod -Uri "http://localhost:8080/health"
```

**Response mong đợi (200 OK):**
```json
{
  "status": "ok",
  "storage": {
    "type": "in-memory",
    "asset_count": 0
  },
  "uptime_seconds": 5.123,
  "timestamp": "2026-06-12T10:00:00Z"
}
```

---

### Bài 2: Batch Create Assets

#### 1. Trường hợp tạo thành công (tạo 3 assets)

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

**Response mong đợi (201 Created):**
```json
{
  "created": 3,
  "ids": ["uuid-1", "uuid-2", "uuid-3"]
}
```

#### 2. Trường hợp lỗi (chứa type không hợp lệ -> Không tạo asset nào - All or Nothing)

**Linux/macOS (Bash):**
```bash
curl -X POST http://localhost:8080/assets/batch \
  -H "Content-Type: application/json" \
  -d '{
    "assets": [
      {"name": "good.com", "type": "domain"},
      {"name": "bad.com", "type": "invalid_type"}
    ]
  }'
```

**Windows CMD:**
```cmd
curl.exe -X POST http://localhost:8080/assets/batch -H "Content-Type: application/json" -d "{\"assets\": [{\"name\": \"good.com\", \"type\": \"domain\"}, {\"name\": \"bad.com\", \"type\": \"invalid_type\"}]}"
```

**Windows PowerShell:**
```powershell
Invoke-RestMethod -Uri "http://localhost:8080/assets/batch" -Method POST -ContentType "application/json" -Body '{"assets": [{"name": "good.com", "type": "domain"}, {"name": "bad.com", "type": "invalid_type"}]}'
```

**Response mong đợi (400 Bad Request):**
```json
{
  "error": "invalid asset type: must be domain, ip, or service"
}
```

---

### Bài 1.1: Get Assets Statistics

**Linux/macOS (Bash) & Windows CMD:**
```bash
curl http://localhost:8080/assets/stats
```

**Windows PowerShell:**
```powershell
Invoke-RestMethod -Uri "http://localhost:8080/assets/stats"
```

**Response mong đợi (200 OK):**
```json
{
  "total": 3,
  "by_type": {
    "domain": 1,
    "ip": 1,
    "service": 1
  },
  "by_status": {
    "active": 3
  }
}
```

---

### Bài 1.2: Count Assets by Filter

#### 1. Đếm tất cả assets
**Linux/macOS (Bash) & Windows CMD:**
```bash
curl http://localhost:8080/assets/count
```
**Windows PowerShell:**
```powershell
Invoke-RestMethod -Uri "http://localhost:8080/assets/count"
```

#### 2. Đếm theo type (ví dụ: domain)
**Linux/macOS (Bash) & Windows CMD:**
```bash
curl "http://localhost:8080/assets/count?type=domain"
```
**Windows PowerShell:**
```powershell
Invoke-RestMethod -Uri "http://localhost:8080/assets/count?type=domain"
```

#### 3. Đếm theo type VÀ status
**Linux/macOS (Bash) & Windows CMD:**
```bash
curl "http://localhost:8080/assets/count?type=domain&status=active"
```
**Windows PowerShell:**
```powershell
Invoke-RestMethod -Uri "http://localhost:8080/assets/count?type=domain&status=active"
```

**Response mong đợi (200 OK):**
```json
{
  "count": 1,
  "filters": {
    "type": "domain",
    "status": "active"
  }
}
```

---

### Bài 3: Batch Delete Assets

#### Linux/macOS (Bash):
```bash
# Bước 1: Tạo assets và lấy ID
ID1=$(curl -s -X POST http://localhost:8080/assets \
  -H "Content-Type: application/json" \
  -d '{"name":"delete-me-1.com","type":"domain"}' | jq -r '.id')

ID2=$(curl -s -X POST http://localhost:8080/assets \
  -H "Content-Type: application/json" \
  -d '{"name":"delete-me-2.com","type":"domain"}' | jq -r '.id')

echo "Created: $ID1, $ID2"

# Bước 2: Batch delete (2 ID thật + 1 ID giả)
curl -X DELETE "http://localhost:8080/assets/batch?ids=$ID1,$ID2,fake-uuid-123"
```

#### Windows CMD:
```cmd
:: Bước 1: Tạo các assets và sao chép ID từ response trả về
curl.exe -X POST http://localhost:8080/assets -H "Content-Type: application/json" -d "{\"name\":\"delete-me-1.com\",\"type\":\"domain\"}"
curl.exe -X POST http://localhost:8080/assets -H "Content-Type: application/json" -d "{\"name\":\"delete-me-2.com\",\"type\":\"domain\"}"

:: Bước 2: Thay thế <ID1> và <ID2> bằng ID thực tế của bạn nhận từ response trên
curl.exe -X DELETE "http://localhost:8080/assets/batch?ids=<ID1>,<ID2>,fake-uuid-123"
```

#### Windows PowerShell:
```powershell
# Bước 1: Tạo assets và tự động lưu ID vào biến
$res1 = Invoke-RestMethod -Uri "http://localhost:8080/assets" -Method POST -ContentType "application/json" -Body '{"name":"delete-me-1.com","type":"domain"}'
$res2 = Invoke-RestMethod -Uri "http://localhost:8080/assets" -Method POST -ContentType "application/json" -Body '{"name":"delete-me-2.com","type":"domain"}'
$ID1 = $res1.id
$ID2 = $res2.id
Write-Host "Created: $ID1, $ID2"

# Bước 2: Thực hiện Batch Delete
Invoke-RestMethod -Uri "http://localhost:8080/assets/batch?ids=$ID1,$ID2,fake-uuid-123" -Method DELETE
```

**Response mong đợi (200 OK):**
```json
{
  "deleted": 2,
  "not_found": 1
}
```

---

### Bài 4: Concurrent-safe Create

#### Linux/macOS (Bash):
```bash
# Bắn 20 request tạo asset đồng thời
for i in $(seq 1 20); do
  curl -s -X POST http://localhost:8080/assets \
    -H "Content-Type: application/json" \
    -d "{\"name\":\"concurrent-$i.com\",\"type\":\"domain\"}" &
done
wait

# Kiểm tra tổng số lượng - phải tăng đúng 20
curl http://localhost:8080/assets/count
```

#### Windows PowerShell:
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
Write-Host "All 20 requests completed"

# Kiểm tra count
Invoke-RestMethod -Uri "http://localhost:8080/assets/count"
```

---

### Bài 6: Pagination & Filtering (Bonus)

#### 1. Trang 1, mỗi trang 10 items
**Linux/macOS (Bash) & Windows CMD:**
```bash
curl "http://localhost:8080/assets?page=1&limit=10"
```
**Windows PowerShell:**
```powershell
Invoke-RestMethod -Uri "http://localhost:8080/assets?page=1&limit=10"
```

#### 2. Lọc theo type (ví dụ: domain)
**Linux/macOS (Bash) & Windows CMD:**
```bash
curl "http://localhost:8080/assets?type=domain"
```
**Windows PowerShell:**
```powershell
Invoke-RestMethod -Uri "http://localhost:8080/assets?type=domain"
```

#### 3. Kết hợp phân trang + lọc theo type/status
**Linux/macOS (Bash) & Windows CMD:**
```bash
curl "http://localhost:8080/assets?page=1&limit=5&type=domain&status=active"
```
**Windows PowerShell:**
```powershell
Invoke-RestMethod -Uri "http://localhost:8080/assets?page=1&limit=5&type=domain&status=active"
```

**Response mong đợi (200 OK):**
```json
{
  "data": [
    {
      "id": "...",
      "name": "example.com",
      "type": "domain",
      "status": "active",
      "created_at": "..."
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 3,
    "total_pages": 1
  }
}
```

---

### Bài 7: Search by Name (Bonus)

#### 1. Tìm kiếm theo tên (partial match)
**Linux/macOS (Bash) & Windows CMD:**
```bash
curl "http://localhost:8080/assets/search?q=example"
```
**Windows PowerShell:**
```powershell
Invoke-RestMethod -Uri "http://localhost:8080/assets/search?q=example"
```

#### 2. Tìm theo đuôi tên miền
**Linux/macOS (Bash) & Windows CMD:**
```bash
curl "http://localhost:8080/assets/search?q=.com"
```
**Windows PowerShell:**
```powershell
Invoke-RestMethod -Uri "http://localhost:8080/assets/search?q=.com"
```

#### 3. Tìm không phân biệt hoa thường (Case-insensitive)
**Linux/macOS (Bash) & Windows CMD:**
```bash
curl "http://localhost:8080/assets/search?q=WEB"
```
**Windows PowerShell:**
```powershell
Invoke-RestMethod -Uri "http://localhost:8080/assets/search?q=WEB"
```

**Response mong đợi (200 OK):**
```json
[
  {
    "id": "...",
    "name": "web-server",
    "type": "service",
    "status": "active",
    "created_at": "..."
  }
]
```

---

## 📮 Hướng dẫn test bằng Postman

> **Lưu ý:** Hãy đảm bảo server đang chạy trước khi test (`go run ./cmd/api/main.go`).

### Cài đặt Postman

1. Tải Postman tại: [https://www.postman.com/downloads/](https://www.postman.com/downloads/)
2. Cài đặt và mở Postman
3. Tạo một **Collection** mới đặt tên `Homework Day 1` để quản lý các request

### Thiết lập chung

- **Base URL:** `http://localhost:8080`
- Với các request có Body (POST), luôn chọn:
  - Tab **Body** → chọn **raw** → chọn kiểu **JSON** từ dropdown

---

### Bài 5: Health Check

| Cấu hình | Giá trị |
|-----------|---------|
| **Method** | `GET` |
| **URL** | `http://localhost:8080/health` |
| **Body** | Không cần |

**Các bước thực hiện:**
1. Tạo request mới trong Collection → đặt tên `Bài 5 - Health Check`
2. Chọn method **GET**
3. Nhập URL: `http://localhost:8080/health`
4. Nhấn **Send**

**Response mong đợi (200 OK):**
```json
{
  "status": "ok",
  "storage": {
    "type": "in-memory",
    "asset_count": 0
  },
  "uptime_seconds": 5.123,
  "timestamp": "2026-06-12T10:00:00Z"
}
```

---

### Bài 2: Batch Create Assets

#### Test Case 1: Tạo thành công 3 assets

| Cấu hình | Giá trị |
|-----------|---------|
| **Method** | `POST` |
| **URL** | `http://localhost:8080/assets/batch` |
| **Body** | raw → JSON |

**Các bước thực hiện:**
1. Tạo request mới → đặt tên `Bài 2 - Batch Create (Success)`
2. Chọn method **POST**
3. Nhập URL: `http://localhost:8080/assets/batch`
4. Chọn tab **Body** → **raw** → chọn **JSON**
5. Dán nội dung JSON sau:

```json
{
  "assets": [
    {"name": "example.com", "type": "domain"},
    {"name": "192.168.1.1", "type": "ip"},
    {"name": "web-server", "type": "service"}
  ]
}
```

6. Nhấn **Send**

**Response mong đợi (201 Created):**
```json
{
  "created": 3,
  "ids": ["uuid-1", "uuid-2", "uuid-3"]
}
```

#### Test Case 2: All-or-Nothing (chứa type không hợp lệ)

| Cấu hình | Giá trị |
|-----------|---------|
| **Method** | `POST` |
| **URL** | `http://localhost:8080/assets/batch` |
| **Body** | raw → JSON |

**Các bước thực hiện:**
1. Tạo request mới → đặt tên `Bài 2 - Batch Create (Error)`
2. Cấu hình giống Test Case 1 nhưng dán JSON sau:

```json
{
  "assets": [
    {"name": "good.com", "type": "domain"},
    {"name": "bad.com", "type": "invalid_type"}
  ]
}
```

3. Nhấn **Send**

**Response mong đợi (400 Bad Request):**
```json
{
  "error": "invalid asset type: must be domain, ip, or service"
}
```

> 💡 **Kiểm tra:** Sau khi nhận lỗi, gọi `GET /assets/stats` để xác nhận không có asset nào mới được tạo (All-or-Nothing).

---

### Bài 1.1: Get Assets Statistics

| Cấu hình | Giá trị |
|-----------|---------|
| **Method** | `GET` |
| **URL** | `http://localhost:8080/assets/stats` |
| **Body** | Không cần |

**Các bước thực hiện:**
1. Tạo request mới → đặt tên `Bài 1.1 - Statistics`
2. Chọn method **GET**
3. Nhập URL: `http://localhost:8080/assets/stats`
4. Nhấn **Send**

**Response mong đợi (200 OK):**
```json
{
  "total": 3,
  "by_type": {
    "domain": 1,
    "ip": 1,
    "service": 1
  },
  "by_status": {
    "active": 3
  }
}
```

---

### Bài 1.2: Count Assets by Filter

#### 1. Đếm tất cả assets

| Cấu hình | Giá trị |
|-----------|---------|
| **Method** | `GET` |
| **URL** | `http://localhost:8080/assets/count` |

#### 2. Đếm theo type

| Cấu hình | Giá trị |
|-----------|---------|
| **Method** | `GET` |
| **URL** | `http://localhost:8080/assets/count?type=domain` |

Hoặc sử dụng tab **Params** trong Postman:

| Key | Value |
|-----|-------|
| `type` | `domain` |

#### 3. Đếm theo type VÀ status

| Cấu hình | Giá trị |
|-----------|---------|
| **Method** | `GET` |
| **URL** | `http://localhost:8080/assets/count?type=domain&status=active` |

Hoặc sử dụng tab **Params**:

| Key | Value |
|-----|-------|
| `type` | `domain` |
| `status` | `active` |

**Các bước thực hiện:**
1. Tạo request mới → đặt tên `Bài 1.2 - Count by Filter`
2. Chọn method **GET**
3. Nhập URL: `http://localhost:8080/assets/count`
4. Chọn tab **Params** → thêm query params theo bảng ở trên
5. Nhấn **Send**

**Response mong đợi (200 OK):**
```json
{
  "count": 1,
  "filters": {
    "type": "domain",
    "status": "active"
  }
}
```

---

### Bài 3: Batch Delete Assets

**Bước 1: Tạo 2 assets mới để test xóa**

| Cấu hình | Giá trị |
|-----------|---------|
| **Method** | `POST` |
| **URL** | `http://localhost:8080/assets` |
| **Body** | raw → JSON |

Request lần 1:
```json
{"name": "delete-me-1.com", "type": "domain"}
```

Request lần 2:
```json
{"name": "delete-me-2.com", "type": "domain"}
```

> 📋 **Lưu ý:** Sau mỗi lần Send, sao chép giá trị `"id"` từ response để dùng ở bước 2.

**Bước 2: Batch Delete**

| Cấu hình | Giá trị |
|-----------|---------|
| **Method** | `DELETE` |
| **URL** | `http://localhost:8080/assets/batch?ids=<ID1>,<ID2>,fake-uuid-123` |

Hoặc sử dụng tab **Params**:

| Key | Value |
|-----|-------|
| `ids` | `<ID1>,<ID2>,fake-uuid-123` |

> Thay `<ID1>` và `<ID2>` bằng ID thực tế bạn vừa sao chép ở bước 1.

**Các bước thực hiện:**
1. Tạo request mới → đặt tên `Bài 3 - Batch Delete`
2. Chọn method **DELETE**
3. Nhập URL và thay thế ID
4. Nhấn **Send**

**Response mong đợi (200 OK):**
```json
{
  "deleted": 2,
  "not_found": 1
}
```

> 💡 **Kiểm tra:** Gọi `GET /assets/<ID1>` để xác nhận trả về **404 Not Found**.

---

### Bài 4: Concurrent-safe Create

> ⚠️ **Lưu ý:** Postman không hỗ trợ gửi request đồng thời trực tiếp. Tuy nhiên bạn có thể sử dụng **Postman Collection Runner** để test tuần tự, hoặc dùng curl/PowerShell cho concurrent test (xem phần hướng dẫn curl ở trên).

**Cách test bằng Postman Collection Runner:**

1. Tạo request mới → đặt tên `Bài 4 - Create Asset`

| Cấu hình | Giá trị |
|-----------|---------|
| **Method** | `POST` |
| **URL** | `http://localhost:8080/assets` |
| **Body** | raw → JSON |

```json
{"name": "concurrent-{{iteration}}.com", "type": "domain"}
```

2. Mở **Collection Runner** (nút ▶ **Run** trên Collection)
3. Chọn request `Bài 4 - Create Asset`
4. Đặt **Iterations** = `20`
5. Nhấn **Run**
6. Sau khi hoàn thành, tạo request `GET http://localhost:8080/assets/count` để kiểm tra count đã tăng đúng 20

**Tiêu chí đạt:**
- ✅ Server không crash
- ✅ Không mất dữ liệu
- ✅ Count tăng đúng bằng số request

---

### Bài 6: Pagination & Filtering (Bonus)

#### 1. Phân trang

| Cấu hình | Giá trị |
|-----------|---------|
| **Method** | `GET` |
| **URL** | `http://localhost:8080/assets` |

Sử dụng tab **Params**:

| Key | Value |
|-----|-------|
| `page` | `1` |
| `limit` | `10` |

#### 2. Lọc theo type

| Key | Value |
|-----|-------|
| `type` | `domain` |

#### 3. Kết hợp phân trang + lọc

| Key | Value |
|-----|-------|
| `page` | `1` |
| `limit` | `5` |
| `type` | `domain` |
| `status` | `active` |

**Các bước thực hiện:**
1. Tạo request mới → đặt tên `Bài 6 - Pagination & Filtering`
2. Chọn method **GET**
3. Nhập URL: `http://localhost:8080/assets`
4. Chọn tab **Params** → thêm các query params theo bảng
5. Nhấn **Send**

**Response mong đợi (200 OK):**
```json
{
  "data": [
    {
      "id": "...",
      "name": "example.com",
      "type": "domain",
      "status": "active",
      "created_at": "..."
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 3,
    "total_pages": 1
  }
}
```

**Query params hỗ trợ:**

| Param    | Mặc định | Mô tả                              |
|----------|----------|-------------------------------------|
| `page`   | 1        | Số trang                            |
| `limit`  | 20       | Số items/trang (tối đa 100)        |
| `type`   | —        | Lọc theo type: domain, ip, service  |
| `status` | —        | Lọc theo status: active, inactive   |

---

### Bài 7: Search by Name (Bonus)

| Cấu hình | Giá trị |
|-----------|---------|
| **Method** | `GET` |
| **URL** | `http://localhost:8080/assets/search` |

Sử dụng tab **Params**:

| Key | Value | Mô tả |
|-----|-------|--------|
| `q` | `example` | Tìm theo từ khóa |
| `q` | `.com` | Tìm theo đuôi tên miền |
| `q` | `WEB` | Case-insensitive search |

**Các bước thực hiện:**
1. Tạo request mới → đặt tên `Bài 7 - Search by Name`
2. Chọn method **GET**
3. Nhập URL: `http://localhost:8080/assets/search`
4. Chọn tab **Params** → thêm key `q` với value muốn tìm
5. Nhấn **Send**

**Response mong đợi (200 OK):**
```json
[
  {
    "id": "...",
    "name": "web-server",
    "type": "service",
    "status": "active",
    "created_at": "..."
  }
]
```

---

### 📝 Mẹo sử dụng Postman hiệu quả

| Mẹo | Mô tả |
|------|--------|
| **Lưu response** | Click **Save Response** để lưu lại kết quả test |
| **Environment** | Tạo Environment với biến `base_url = http://localhost:8080`, rồi dùng `{{base_url}}` trong URL |
| **Variables** | Dùng **Tests** script để tự động lưu ID: `pm.environment.set("asset_id", pm.response.json().id)` |
| **Collection Runner** | Chạy toàn bộ Collection tự động để kiểm tra tất cả endpoints cùng lúc |
| **Export** | Export Collection ra file `.json` để chia sẻ với bạn cùng nhóm |

---

## 🧪 Chạy Unit Tests

```bash
# Chạy tất cả tests
go test ./... -v

# Output mong đợi:
# === RUN   TestConcurrentCreate
# --- PASS: TestConcurrentCreate (0.00s)
# === RUN   TestConcurrentBatchCreateAndStats
# --- PASS: TestConcurrentBatchCreateAndStats (0.00s)
# PASS
```

```bash
# Chạy với race detector (yêu cầu GCC trên Windows)
go test ./... -race
```

---

## 📝 Ghi chú

- **Storage**: Dữ liệu lưu trong bộ nhớ (RAM), sẽ mất khi tắt server.
- **Concurrent Safety**: Sử dụng `sync.RWMutex` — đọc đồng thời, ghi tuần tự.
- **Validation**: Type chỉ chấp nhận `domain`, `ip`, `service`. Status chỉ chấp nhận `active`, `inactive`.
- **Batch limit**: Tối đa 100 assets mỗi request batch create.
