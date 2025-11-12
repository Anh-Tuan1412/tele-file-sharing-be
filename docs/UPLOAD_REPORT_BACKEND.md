# Backend cho Báo hoàn tất Upload

Hướng dẫn triển khai và hệ thống cho tính năng "Báo hoàn tất upload" trong dịch vụ chia sẻ file.

## Tổng Quan

Tính năng "Báo hoàn tất upload" cho phép người dùng gửi báo cáo về trạng thái hoàn tất upload file, bao gồm thông tin về:
- Trạng thái upload (thành công/thất bại)
- Loại báo cáo (success/error/warning)
- Chi tiết lỗi (nếu có)
- Thông tin hiệu năng (thời gian upload, tốc độ, kích thước)
- Checksum file

## Thay Đổi Trong Codebase

### 1. Migrations

#### Tệp mới: `migrations/003_create_upload_reports_table.sql`

Tạo bảng `upload_reports` với các trường:
- `id`: Khóa chính
- `file_id`: Khóa ngoại tới bảng `files`
- `owner_user_id`: Khóa ngoại tới bảng `users`
- `status`: Trạng thái báo cáo (pending, completed, failed, processing)
- `report_type`: Loại báo cáo (success, error, warning)
- `message`: Thông điệp mô tả
- `error_code`: Mã lỗi (nếu có)
- `error_message`: Chi tiết lỗi
- `file_checksum`: MD5/SHA256 checksum
- `file_size_actual`: Kích thước thực tế của file
- `upload_duration_ms`: Thời gian upload (milliseconds)
- `bandwidth_kbps`: Tốc độ upload (Kbps)
- `reported_at`: Timestamp báo cáo
- `updated_at`: Timestamp cập nhật lần cuối

Cấu trúc chỉ mục:
- PK trên `id`
- FK `file_id` với cascade delete
- FK `owner_user_id` với cascade delete
- Index trên `file_id`, `owner_user_id`, `status`
- Index trên `reported_at DESC` để tìm báo cáo gần đây
- Unique index trên `file_id` cho báo cáo hoàn tất (status = 'completed')

### 2. Models

#### Tệp mới: `internal/model/upload_report.go`

Định nghĩa struct `UploadReport` với các trường và hằng số:

```go
type UploadReport struct {
    ID                   int64     // ID báo cáo
    FileID               int64     // ID file
    OwnerUserID          int64     // ID chủ sở hữu
    Status               string    // pending, completed, failed, processing
    ReportType           string    // success, error, warning
    Message              string    // Thông điệp
    ErrorCode            string    // Mã lỗi
    ErrorMessage         string    // Chi tiết lỗi
    FileChecksum         string    // Checksum file
    FileSizeActual       int64     // Kích thước file
    UploadDurationMs     int       // Thời gian upload
    BandwidthKbps        float64   // Tốc độ upload
    ReportedAt           time.Time // Thời gian báo cáo
    UpdatedAt            time.Time // Thời gian cập nhật
}

// Hằng số Status
const (
    ReportStatusPending   = "pending"
    ReportStatusCompleted = "completed"
    ReportStatusFailed    = "failed"
    ReportStatusProcessing = "processing"
)

// Hằng số ReportType
const (
    ReportTypeSuccess = "success"
    ReportTypeError   = "error"
    ReportTypeWarning = "warning"
)
```

#### Cập nhật: `internal/model/user.go`

Thêm struct `FileWithOwner` để chứa thông tin file kèm thông tin chủ sở hữu:

```go
type FileWithOwner struct {
    ID            int64     // ID file
    OwnerUserID   int64     // ID chủ sở hữu
    ObjectKey     string    // Đường dẫn object
    Filename      string    // Tên file
    Size          int64     // Kích thước
    Mime          string    // MIME type
    Status        string    // Trạng thái
    CreatedAt     time.Time
    UpdatedAt     time.Time
    TelegramID    int64     // Telegram ID của chủ sở hữu
    Username      string    // Username của chủ sở hữu
}
```

### 3. Storage Layer

#### Tệp mới: `internal/storage/file.go`

Tạo repository interface `FileRepository` và implementation PostgreSQL `postgresFileRepository` với các phương thức:

**Các phương thức chính:**
- `GetFileByID(ctx, fileID)`: Lấy file theo ID
- `GetFileByIDAndUserID(ctx, fileID, userID)`: Lấy file theo ID và ID chủ sở hữu
- `UpdateFileStatus(ctx, fileID, status)`: Cập nhật trạng thái file
- `CreateUploadReport(ctx, report)`: Tạo báo cáo upload mới
- `GetUploadReport(ctx, fileID)`: Lấy báo cáo gần đây nhất cho file
- `GetUploadReportsByUserID(ctx, userID, limit, offset)`: Lấy danh sách báo cáo của user
- `UpdateUploadReportStatus(ctx, reportID, status)`: Cập nhật trạng thái báo cáo

**Tính năng bảo mật:**
- Kiểm tra chủ sở hữu file trước khi trả về thông tin
- Xác minh quyền truy cập của user
- Sử dụng parameterized queries để phòng chống SQL injection

### 4. HTTP Handlers

#### Cập nhật: `internal/transport/http/file_handler.go`

Thêm ba handler mới:

**1. POST /v1/files/{file_id}/report-complete - ReportUploadCompleteHandler**

Xử lý báo cáo hoàn tất upload:
- Xác thực user từ Telegram headers
- Kiểm tra file tồn tại và thuộc sở hữu của user
- Tạo record trong bảng `upload_reports`
- Cập nhật trạng thái file (completed/failed)
- Trả về `ReportUploadCompleteResponse` với HTTP 201

Request body:
```json
{
    "status": "completed",
    "report_type": "success",
    "message": "Upload completed successfully",
    "error_code": "",
    "error_message": "",
    "file_checksum": "a1b2c3d4e5f6",
    "file_size_actual": 204800,
    "upload_duration_ms": 5000,
    "bandwidth_kbps": 512.5
}
```

Response (201):
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

**2. GET /v1/files/{file_id}/report - GetUploadReportHandler**

Lấy báo cáo gần đây nhất cho file:
- Xác thực user từ Telegram headers
- Kiểm tra quyền truy cập
- Trả về báo cáo mới nhất (sắp xếp theo `reported_at DESC`)

Response (200):
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
    "file_size_actual": 204800,
    "upload_duration_ms": 5000,
    "bandwidth_kbps": 512.5,
    "reported_at": "2025-11-12T10:30:00Z",
    "updated_at": "2025-11-12T10:30:00Z"
}
```

**3. GET /v1/upload-reports - ListUploadReportsHandler**

Liệt kê tất cả báo cáo của user:
- Xác thực user từ Telegram headers
- Hỗ trợ phân trang (limit, offset)
- Sắp xếp theo thời gian gần nhất trước
- Limit mặc định: 20, tối đa: 100

Query parameters:
- `limit`: Số lượng tối đa (default: 20, max: 100)
- `offset`: Vị trí bắt đầu (default: 0)

Response (200):
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
    ...
]
```

### 5. API Documentation

#### Cập nhật: `api/openapi.yaml`

Thêm ba endpoint mới:

**1. POST /v1/files/{file_id}/report-complete**
- Tag: Files
- Request schema: `ReportUploadCompleteRequest`
- Response: `ReportUploadCompleteResponse` (201)

**2. GET /v1/files/{file_id}/report**
- Tag: Files
- Response: `UploadReport` (200)

**3. GET /v1/upload-reports**
- Tag: Files
- Query parameters: `limit`, `offset`
- Response: Array of `UploadReport` (200)

**Schemas mới:**
- `ReportUploadCompleteRequest`: Request body cho báo cáo
- `ReportUploadCompleteResponse`: Response cho tạo báo cáo
- `UploadReport`: Model báo cáo upload hoàn chỉnh

## Quy Trình Hoạt Động

### Upload Completion Flow

```
1. Client khởi tạo upload
   → POST /v1/files
   ← file_id, status = "pending"

2. Client upload file
   → Upload file via presigned URL hoặc direct

3. Client báo cáo hoàn tất
   → POST /v1/files/{file_id}/report-complete
   ← report_id, status = "completed/failed"

4. Client có thể:
   a) Xem báo cáo của file cụ thể
      → GET /v1/files/{file_id}/report
      ← UploadReport details
   
   b) Xem tất cả báo cáo của user
      → GET /v1/upload-reports?limit=20&offset=0
      ← Array of UploadReport
```

## Các Tính Năng Bảo Mật

### Authentication & Authorization
- Tất cả endpoint yêu cầu Telegram headers (`X-Telegram-User-Id`, `X-Telegram-Username`)
- User chỉ có thể:
  - Tạo báo cáo cho file thuộc sở hữu của mình
  - Xem báo cáo của file mình sở hữu
  - Xem danh sách báo cáo của chính mình

### Validation
- Kiểm tra quyền sở hữu file trước mỗi thao tác
- Validate enum values (status: completed/failed, report_type: success/error/warning)
- Validate file_id format
- Kiểm tra file tồn tại trong database

### Data Integrity
- Sử dụng foreign keys với cascade delete
- Unique index trên `upload_reports.file_id` cho báo cáo hoàn tất
- Parameterized queries để chống SQL injection

## Integration Notes

### Database Connection
- Sử dụng `*sql.DB` hoặc `*sqlx.DB` (tùy vào architecture)
- Transactions nên được sử dụng nếu cần atomicity
- Connection pooling được khuyến nghị

### Error Handling
- 404: File không tồn tại hoặc không có quyền truy cập
- 400: Bad request (missing headers, invalid data)
- 500: Database error, failed queries

### Logging
- Tất cả database errors được log bằng `log.Printf`
- Lưu file ID và user ID trong log message

## Các Bước Triển Khai

### 1. Database Migration
```bash
# Chạy migration để tạo bảng upload_reports
# Ví dụ: flyway, migrate, hoặc manual SQL
psql -U user -d database -f migrations/003_create_upload_reports_table.sql
```

### 2. Cập nhật Code
- Copy file `internal/model/upload_report.go`
- Copy file `internal/storage/file.go`
- Cập nhật `internal/model/user.go` (thêm FileWithOwner)
- Cập nhật `internal/transport/http/file_handler.go` (thêm handlers)

### 3. Cập nhật Routes
Trong file khởi tạo routes (ví dụ `cmd/api/main.go`), đăng ký các endpoint mới:

```go
// POST /v1/files/{file_id}/report-complete
router.POST("/v1/files/:file_id/report-complete", http.ReportUploadCompleteHandler(db))

// GET /v1/files/{file_id}/report
router.GET("/v1/files/:file_id/report", http.GetUploadReportHandler(db))

// GET /v1/upload-reports
router.GET("/v1/upload-reports", http.ListUploadReportsHandler(db))
```

### 4. Testing
```bash
# Test create report
curl -X POST http://localhost:8080/v1/files/1/report-complete \
  -H "X-Telegram-User-Id: 123456" \
  -H "X-Telegram-Username: testuser" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "completed",
    "report_type": "success",
    "file_checksum": "abc123",
    "file_size_actual": 1024,
    "upload_duration_ms": 5000,
    "bandwidth_kbps": 512.5
  }'

# Test get report
curl -X GET http://localhost:8080/v1/files/1/report \
  -H "X-Telegram-User-Id: 123456" \
  -H "X-Telegram-Username: testuser"

# Test list reports
curl -X GET http://localhost:8080/v1/upload-reports?limit=10 \
  -H "X-Telegram-User-Id: 123456" \
  -H "X-Telegram-Username: testuser"
```

## Future Enhancements

1. **Async Processing**: Sử dụng message queue để xử lý báo cáo không đồng bộ
2. **Retry Logic**: Tự động thử lại khi upload thất bại
3. **Batch Reporting**: Cho phép báo cáo nhiều file cùng một lúc
4. **Webhook Notifications**: Thông báo client khi upload hoàn tất
5. **Analytics**: Thống kê về tốc độ upload, tỉ lệ thành công
6. **Caching**: Cache báo cáo gần đây cho hiệu năng tốt hơn
7. **Rate Limiting**: Giới hạn số báo cáo per user per minute

## Hỗ Trợ & Troubleshooting

### Common Issues

**Issue: "database error"**
- Kiểm tra kết nối database
- Kiểm tra migration đã được chạy
- Kiểm tra permissions của database user

**Issue: "file not found or access denied"**
- Kiểm tra file_id có tồn tại
- Kiểm tra user có sở hữu file
- Kiểm tra Telegram headers đúng

**Issue: "missing Telegram headers"**
- Kiểm tra request có headers `X-Telegram-User-Id` và `X-Telegram-Username`
- Headers không được thiếu hoặc trống

---

Được tạo: November 12, 2025
Backend cho tính năng "Báo hoàn tất upload" - File Sharing API
