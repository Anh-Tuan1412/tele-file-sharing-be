# Integration Guide - Upload Report Backend

Hướng dẫn nhanh để tích hợp backend "Báo hoàn tất upload" vào dự án.

## 📋 Checklist Triển Khai

### Phase 1: Database Setup
- [ ] Chạy migration: `migrations/003_create_upload_reports_table.sql`
  ```bash
  psql -U <user> -d <database> -f migrations/003_create_upload_reports_table.sql
  ```
- [ ] Verify tables được tạo
  ```sql
  SELECT * FROM upload_reports LIMIT 1;
  ```

### Phase 2: Code Integration
- [ ] Copy file `internal/model/upload_report.go`
- [ ] Copy file `internal/storage/file.go`
- [ ] Update `internal/model/user.go` (thêm `FileWithOwner` struct)
- [ ] Update `internal/transport/http/file_handler.go` (thêm 3 handlers mới)
- [ ] Update `api/openapi.yaml` (thêm 3 endpoints mới)

### Phase 3: Router Configuration
Thêm vào file khởi tạo routes (thường là `cmd/api/main.go`):

```go
// Import handler
import (
    "file-sharing/internal/transport/http"
)

// Trong hàm setup routes
func setupRoutes(router *gin.Engine, db *sql.DB) {
    // ... existing routes ...
    
    // Upload Report endpoints
    router.POST("/v1/files/:file_id/report-complete", http.ReportUploadCompleteHandler(db))
    router.GET("/v1/files/:file_id/report", http.GetUploadReportHandler(db))
    router.GET("/v1/upload-reports", http.ListUploadReportsHandler(db))
}
```

### Phase 4: Testing
- [ ] Test create report
  ```bash
  curl -X POST http://localhost:8080/v1/files/1/report-complete \
    -H "X-Telegram-User-Id: 123456" \
    -H "X-Telegram-Username: testuser" \
    -H "Content-Type: application/json" \
    -d '{
      "status": "completed",
      "report_type": "success",
      "message": "Upload completed successfully",
      "file_checksum": "abc123def456",
      "file_size_actual": 1024000,
      "upload_duration_ms": 5000,
      "bandwidth_kbps": 512.5
    }'
  ```

- [ ] Test get report for file
  ```bash
  curl -X GET http://localhost:8080/v1/files/1/report \
    -H "X-Telegram-User-Id: 123456" \
    -H "X-Telegram-Username: testuser"
  ```

- [ ] Test list all reports
  ```bash
  curl -X GET "http://localhost:8080/v1/upload-reports?limit=10&offset=0" \
    -H "X-Telegram-User-Id: 123456" \
    -H "X-Telegram-Username: testuser"
  ```

### Phase 5: Documentation
- [ ] Update project README
- [ ] Document API endpoints in CHANGELOG
- [ ] Provide client implementation guide

## 🔧 Troubleshooting

### Problem: "database error"
**Solution:**
- Kiểm tra migration: `SELECT * FROM information_schema.tables WHERE table_name = 'upload_reports';`
- Kiểm tra connection string
- Kiểm tra database user permissions

### Problem: "file not found or access denied"
**Solution:**
- Tạo file trước: POST `/v1/files` để get file_id
- Kiểm tra Telegram headers
- Kiểm tra file_id tồn tại: `SELECT * FROM files WHERE id = <file_id>;`

### Problem: "missing Telegram headers"
**Solution:**
- Thêm headers: `X-Telegram-User-Id` và `X-Telegram-Username`
- Headers không được để trống

### Problem: Validation error
**Solution:**
- status phải là: `completed` hoặc `failed`
- report_type phải là: `success`, `error`, hoặc `warning`
- file_size_actual phải là integer
- Tất cả required fields phải có

## 📝 API Endpoints Quick Reference

| Method | Endpoint | Purpose |
|--------|----------|---------|
| POST | `/v1/files/{file_id}/report-complete` | Báo cáo hoàn tất upload |
| GET | `/v1/files/{file_id}/report` | Lấy báo cáo của file |
| GET | `/v1/upload-reports` | Liệt kê tất cả báo cáo |

### Required Headers (All Endpoints)
```
X-Telegram-User-Id: <int64>
X-Telegram-Username: <string>
```

### Example Request - Report Completion
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

### Example Response - Report Completion
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

## 📊 Database Schema

### Table: upload_reports
```sql
id                   | SERIAL PRIMARY KEY
file_id              | INT REFERENCES files(id) ON DELETE CASCADE
owner_user_id        | INT REFERENCES users(id) ON DELETE CASCADE
status               | VARCHAR(20) DEFAULT 'pending'
report_type          | VARCHAR(50) NOT NULL
message              | TEXT
error_code           | VARCHAR(50)
error_message        | TEXT
file_checksum        | TEXT
file_size_actual     | BIGINT
upload_duration_ms   | INT
bandwidth_kbps       | DECIMAL(10, 2)
reported_at          | TIMESTAMPTZ DEFAULT NOW()
updated_at           | TIMESTAMPTZ DEFAULT NOW()

-- Indexes
CREATE INDEX idx_upload_reports_file_id ON upload_reports(file_id);
CREATE INDEX idx_upload_reports_owner_user_id ON upload_reports(owner_user_id);
CREATE INDEX idx_upload_reports_status ON upload_reports(status);
CREATE INDEX idx_upload_reports_reported_at ON upload_reports(reported_at DESC);
CREATE UNIQUE INDEX idx_upload_reports_file_id_unique ON upload_reports(file_id) WHERE status = 'completed';
```

## 🔒 Security Checklist

- [x] Authentication via Telegram headers
- [x] Authorization check (user owns file)
- [x] Input validation
- [x] SQL injection prevention (parameterized queries)
- [x] Proper error messages (no sensitive data leakage)
- [x] Database constraints (FK, unique indexes)

## 🚀 Performance Tips

1. **Indexing**: Indexes trên `file_id`, `owner_user_id`, `reported_at` đã được tạo
2. **Pagination**: Luôn sử dụng `limit` và `offset` khi liệt kê
3. **Caching**: Nếu cần, cache báo cáo gần đây
4. **Batch Operations**: Nếu cần batch reports, xem xét async processing

## 📞 Support

Nếu gặp lỗi:
1. Kiểm tra logs của API server
2. Verify database migration chạy thành công
3. Test cURL commands một cách từng bước
4. Kiểm tra Telegram headers format

---

Last Updated: November 12, 2025
