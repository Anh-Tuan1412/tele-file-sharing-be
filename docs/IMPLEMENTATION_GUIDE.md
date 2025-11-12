# Implementation Guide - Báo Hoàn Tất Upload Backend

Hướng dẫn chi tiết để triển khai backend "Báo hoàn tất upload" vào project dath-tele.

## 📁 Cấu Trúc File

```
dath-tele/
├── migrations/
│   └── 003_create_upload_reports_table.sql    [NEW]
├── internal/
│   ├── model/
│   │   ├── upload_report.go                    [NEW]
│   │   └── user.go                             [UPDATED]
│   ├── files/
│   │   ├── service.go                          [NEW]
│   │   ├── service_test.go                     [NEW]
│   │   └── error.go                            [NEW]
│   ├── storage/
│   │   └── file.go                             [NEW]
│   └── transport/http/
│       ├── file_handler.go                     [UPDATED]
│       ├── error_handler.go                    [UPDATED]
│       └── validation.go                       [UPDATED]
├── cmd/
│   └── example-client/
│       └── main.go                             [UPDATED]
└── api/
    └── openapi.yaml                            [UPDATED]
```

## 🚀 Bước Triển Khai

### Step 1: Database Migration

Chạy migration để tạo bảng `upload_reports`:

```bash
# Với migrate tool
migrate -path migrations -database "postgresql://user:password@localhost/database" up

# Hoặc với psql
psql -U <user> -d <database> -f migrations/003_create_upload_reports_table.sql

# Hoặc với sqlc (nếu sử dụng)
sqlc migrate apply -d "postgresql://user:password@localhost/database"
```

Xác nhận migration thành công:

```sql
-- Connect to your database
psql -U <user> -d <database>

-- Check table
\dt upload_reports

-- Check structure
\d upload_reports

-- Check indexes
\di upload_reports*
```

### Step 2: Copy Model Files

**File: `internal/model/upload_report.go`**
- Định nghĩa struct `UploadReport`
- Constants cho status và report_type

**File: `internal/model/user.go`** (UPDATED)
- Thêm struct `FileWithOwner`

### Step 3: Copy Storage Layer

**File: `internal/storage/file.go`**
- Interface `FileRepository`
- Implementation `postgresFileRepository`
- Các method để làm việc với upload reports

### Step 4: Copy Service Layer

**File: `internal/files/service.go`**
- Interface `Service`
- Implementation `fileService`
- Business logic cho upload reports

**File: `internal/files/error.go`**
- Custom error types
- Error constructors

### Step 5: Update Transport Layer

**File: `internal/transport/http/validation.go`** (UPDATED)
- `ValidateReportUploadCompleteRequest()` - Validate request data
- `ValidateTelegramHeaders()` - Validate headers
- `ValidateFileID()` - Validate file ID
- `ValidatePaginationParams()` - Validate pagination

**File: `internal/transport/http/error_handler.go`** (UPDATED)
- `AppError` struct
- `HandleError()` function
- Error response middleware

**File: `internal/transport/http/file_handler.go`** (UPDATED)
- `ReportUploadCompleteHandler()` - POST /v1/files/{file_id}/report-complete
- `GetUploadReportHandler()` - GET /v1/files/{file_id}/report
- `ListUploadReportsHandler()` - GET /v1/upload-reports

### Step 6: Update API Documentation

**File: `api/openapi.yaml`** (UPDATED)
- Thêm 3 endpoint mới
- Thêm 3 schema mới

### Step 7: Update Router Configuration

Trong file khởi tạo route (thường là `cmd/api/main.go`):

```go
package main

import (
    "github.com/gin-gonic/gin"
    "file-sharing/internal/transport/http"
    "file-sharing/internal/storage"
    "file-sharing/internal/files"
)

func setupRoutes(router *gin.Engine, db *sql.DB, sqlxDB *sqlx.DB) {
    // Initialize repositories and services
    fileRepo := storage.NewFileRepository(sqlxDB)
    fileService := files.NewFileService(db, fileRepo)
    
    // File upload report endpoints
    router.POST("/v1/files/:file_id/report-complete", 
        http.ReportUploadCompleteHandler(db))
    router.GET("/v1/files/:file_id/report", 
        http.GetUploadReportHandler(db))
    router.GET("/v1/upload-reports", 
        http.ListUploadReportsHandler(db))
    
    // Apply error middleware
    router.Use(http.ErrorResponseMiddleware())
    router.Use(http.LoggerMiddleware())
}
```

### Step 8: Testing

Chạy unit tests:

```bash
go test -v ./internal/files/...
```

Test endpoints với curl:

```bash
# 1. Report upload completion
curl -X POST http://localhost:8080/v1/files/1/report-complete \
  -H "Content-Type: application/json" \
  -H "X-Telegram-User-Id: 123456789" \
  -H "X-Telegram-Username: testuser" \
  -d '{
    "status": "completed",
    "report_type": "success",
    "message": "Upload completed successfully",
    "file_checksum": "a1b2c3d4e5f6",
    "file_size_actual": 1024000,
    "upload_duration_ms": 5000,
    "bandwidth_kbps": 512.5
  }'

# 2. Get upload report
curl -X GET http://localhost:8080/v1/files/1/report \
  -H "X-Telegram-User-Id: 123456789" \
  -H "X-Telegram-Username: testuser"

# 3. List upload reports
curl -X GET "http://localhost:8080/v1/upload-reports?limit=10&offset=0" \
  -H "X-Telegram-User-Id: 123456789" \
  -H "X-Telegram-Username: testuser"
```

Hoặc sử dụng example client:

```bash
# Report completion
go run cmd/example-client/main.go \
  -url=http://localhost:8080 \
  -telegram-id=123456789 \
  -username=testuser \
  -cmd=report \
  -file-id=1 \
  -status=completed \
  -type=success

# Get report
go run cmd/example-client/main.go \
  -url=http://localhost:8080 \
  -telegram-id=123456789 \
  -username=testuser \
  -cmd=get \
  -file-id=1

# List reports
go run cmd/example-client/main.go \
  -url=http://localhost:8080 \
  -telegram-id=123456789 \
  -username=testuser \
  -cmd=list
```

## 📝 API Endpoints

### 1. POST /v1/files/{file_id}/report-complete

**Purpose**: Báo cáo hoàn tất upload file

**Request**:
```json
{
  "status": "completed",
  "report_type": "success",
  "message": "Upload completed successfully",
  "error_code": "",
  "error_message": "",
  "file_checksum": "a1b2c3d4e5f6",
  "file_size_actual": 1024000,
  "upload_duration_ms": 5000,
  "bandwidth_kbps": 512.5
}
```

**Response (201)**:
```json
{
  "report_id": 1,
  "file_id": 101,
  "status": "completed",
  "report_type": "success",
  "message": "Upload completed successfully",
  "reported_at": "2025-11-12T10:30:00Z"
}
```

### 2. GET /v1/files/{file_id}/report

**Purpose**: Lấy báo cáo gần đây nhất cho file

**Response (200)**:
```json
{
  "id": 1,
  "file_id": 101,
  "owner_user_id": 1,
  "status": "completed",
  "report_type": "success",
  "message": "Upload completed successfully",
  "error_code": "",
  "error_message": "",
  "file_checksum": "a1b2c3d4e5f6",
  "file_size_actual": 1024000,
  "upload_duration_ms": 5000,
  "bandwidth_kbps": 512.5,
  "reported_at": "2025-11-12T10:30:00Z",
  "updated_at": "2025-11-12T10:30:00Z"
}
```

### 3. GET /v1/upload-reports

**Purpose**: Liệt kê tất cả báo cáo của user

**Query Parameters**:
- `limit`: Số lượng tối đa (default: 20, max: 100)
- `offset`: Vị trí bắt đầu (default: 0)

**Response (200)**:
```json
[
  {
    "id": 1,
    "file_id": 101,
    "owner_user_id": 1,
    "status": "completed",
    "report_type": "success",
    ...
  },
  {
    "id": 2,
    "file_id": 102,
    "owner_user_id": 1,
    "status": "completed",
    "report_type": "error",
    ...
  }
]
```

## 🧪 Testing Scenarios

### Scenario 1: Successful Upload Report

```bash
# 1. Upload file
curl -X POST http://localhost:8080/v1/files \
  -H "Content-Type: application/json" \
  -H "X-Telegram-User-Id: 123456789" \
  -H "X-Telegram-Username: testuser" \
  -d '{"filename":"test.txt","size":1024}'

# Response: {"file_id":1,"object_key":"uploads/1/1234567890_test.txt",...}

# 2. Report completion
curl -X POST http://localhost:8080/v1/files/1/report-complete \
  -H "Content-Type: application/json" \
  -H "X-Telegram-User-Id: 123456789" \
  -H "X-Telegram-Username: testuser" \
  -d '{
    "status": "completed",
    "report_type": "success",
    "file_size_actual": 1024,
    "upload_duration_ms": 3000
  }'

# Response: {"report_id":1,"file_id":1,"status":"completed",...}

# 3. Verify report
curl -X GET http://localhost:8080/v1/files/1/report \
  -H "X-Telegram-User-Id: 123456789" \
  -H "X-Telegram-Username: testuser"
```

### Scenario 2: Failed Upload Report

```bash
curl -X POST http://localhost:8080/v1/files/1/report-complete \
  -H "Content-Type: application/json" \
  -H "X-Telegram-User-Id: 123456789" \
  -H "X-Telegram-Username: testuser" \
  -d '{
    "status": "failed",
    "report_type": "error",
    "error_code": "FILE_TOO_LARGE",
    "error_message": "File size exceeds the limit of 100MB"
  }'
```

### Scenario 3: Validation Error

```bash
# Missing required field
curl -X POST http://localhost:8080/v1/files/1/report-complete \
  -H "Content-Type: application/json" \
  -H "X-Telegram-User-Id: 123456789" \
  -H "X-Telegram-Username: testuser" \
  -d '{
    "report_type": "success"
    # status is missing
  }'

# Response (400):
# {
#   "code": "VALIDATION_ERROR",
#   "message": "Request validation failed",
#   "errors": [
#     {"field": "status", "message": "status is required"}
#   ]
# }
```

### Scenario 4: Unauthorized Access

```bash
# Try to access file not owned
curl -X GET http://localhost:8080/v1/files/999/report \
  -H "X-Telegram-User-Id: 123456789" \
  -H "X-Telegram-Username: testuser"

# Response (404):
# {
#   "code": "NOT_FOUND",
#   "message": "File not found or access denied"
# }
```

## 📊 Database Schema

### upload_reports table

| Column | Type | Description |
|--------|------|-------------|
| id | SERIAL PRIMARY KEY | ID báo cáo |
| file_id | INT FK | ID file |
| owner_user_id | INT FK | ID chủ sở hữu |
| status | VARCHAR(20) | pending, completed, failed, processing |
| report_type | VARCHAR(50) | success, error, warning |
| message | TEXT | Thông điệp |
| error_code | VARCHAR(50) | Mã lỗi |
| error_message | TEXT | Chi tiết lỗi |
| file_checksum | TEXT | Checksum |
| file_size_actual | BIGINT | Kích thước |
| upload_duration_ms | INT | Thời gian upload |
| bandwidth_kbps | DECIMAL | Tốc độ upload |
| reported_at | TIMESTAMPTZ | Thời gian báo cáo |
| updated_at | TIMESTAMPTZ | Lần cập nhật cuối |

### Indexes

```sql
CREATE INDEX idx_upload_reports_file_id ON upload_reports(file_id);
CREATE INDEX idx_upload_reports_owner_user_id ON upload_reports(owner_user_id);
CREATE INDEX idx_upload_reports_status ON upload_reports(status);
CREATE INDEX idx_upload_reports_reported_at ON upload_reports(reported_at DESC);
CREATE UNIQUE INDEX idx_upload_reports_file_id_unique ON upload_reports(file_id) WHERE status = 'completed';
```

## 🔒 Security Features

- ✅ Authentication via Telegram headers
- ✅ Authorization check (user owns file)
- ✅ Input validation (enum, length, format)
- ✅ SQL injection prevention
- ✅ Error handling (no sensitive data leakage)
- ✅ Database constraints (FK, unique indexes)

## 🐛 Troubleshooting

### Database Connection Error

**Problem**: "database error"

**Solution**:
```bash
# Check connection
psql -U <user> -d <database> -c "SELECT 1"

# Check table exists
psql -U <user> -d <database> -c "\dt upload_reports"
```

### File Not Found Error

**Problem**: "file not found or access denied"

**Solution**:
```bash
# 1. Create file first
curl -X POST http://localhost:8080/v1/files \
  -H "Content-Type: application/json" \
  -H "X-Telegram-User-Id: 123456789" \
  -H "X-Telegram-Username: testuser" \
  -d '{"filename":"test.txt","size":1024}'

# 2. Use returned file_id
```

### Validation Error

**Problem**: "validation errors"

**Solution**:
- Ensure all required fields are present
- Check enum values (status, report_type)
- Verify field formats and lengths

## 📚 Documentation Files

- `docs/UPLOAD_REPORT_BACKEND.md` - Tài liệu chi tiết backend
- `docs/INTEGRATION_GUIDE.md` - Hướng dẫn tích hợp nhanh
- `api/openapi.yaml` - OpenAPI specification

## ✅ Checklist Hoàn Thiện

- [ ] Database migration chạy thành công
- [ ] Copy tất cả file model
- [ ] Copy tất cả file storage
- [ ] Copy tất cả file service
- [ ] Update file handler
- [ ] Update file validation
- [ ] Update file error_handler
- [ ] Update OpenAPI spec
- [ ] Update router configuration
- [ ] Run unit tests
- [ ] Test endpoints với curl
- [ ] Test với example client
- [ ] Code review
- [ ] Deploy to staging
- [ ] Deploy to production

---

**Last Updated**: November 12, 2025
**Version**: 1.0.0
