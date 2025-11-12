# Quick Reference - Upload Report Backend

Tham chiếu nhanh cho backend "Báo hoàn tất upload".

## 📦 Files to Add/Modify

### New Files (6)
1. `migrations/003_create_upload_reports_table.sql`
2. `internal/model/upload_report.go`
3. `internal/storage/file.go`
4. `internal/files/service.go`
5. `internal/files/error.go`
6. `internal/files/service_test.go`

### Modified Files (5)
1. `internal/model/user.go` - Add `FileWithOwner` struct
2. `internal/transport/http/file_handler.go` - Add 3 handlers
3. `internal/transport/http/validation.go` - Add validators
4. `internal/transport/http/error_handler.go` - Enhance error handling
5. `api/openapi.yaml` - Add 3 endpoints + 3 schemas
6. `cmd/example-client/main.go` - Add test client

### Documentation Files (3)
1. `docs/UPLOAD_REPORT_BACKEND.md`
2. `docs/INTEGRATION_GUIDE.md`
3. `docs/IMPLEMENTATION_GUIDE.md`

## 🔗 API Endpoints

| Method | Endpoint | Status | Purpose |
|--------|----------|--------|---------|
| POST | `/v1/files/{file_id}/report-complete` | 201 | Report upload completion |
| GET | `/v1/files/{file_id}/report` | 200 | Get latest report for file |
| GET | `/v1/upload-reports` | 200 | List user's all reports |

## 📋 Request Examples

### Report Completion (POST)
```bash
curl -X POST http://localhost:8080/v1/files/1/report-complete \
  -H "X-Telegram-User-Id: 123456789" \
  -H "X-Telegram-Username: testuser" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "completed",
    "report_type": "success",
    "message": "Upload successful",
    "file_checksum": "abc123",
    "file_size_actual": 1024,
    "upload_duration_ms": 5000,
    "bandwidth_kbps": 512.5
  }'
```

### Get Report (GET)
```bash
curl -X GET http://localhost:8080/v1/files/1/report \
  -H "X-Telegram-User-Id: 123456789" \
  -H "X-Telegram-Username: testuser"
```

### List Reports (GET)
```bash
curl -X GET "http://localhost:8080/v1/upload-reports?limit=20&offset=0" \
  -H "X-Telegram-User-Id: 123456789" \
  -H "X-Telegram-Username: testuser"
```

## 🗄️ Database Schema Quick View

```sql
-- Main table
CREATE TABLE upload_reports (
    id SERIAL PRIMARY KEY,
    file_id INT REFERENCES files(id) ON DELETE CASCADE,
    owner_user_id INT REFERENCES users(id) ON DELETE CASCADE,
    status VARCHAR(20),
    report_type VARCHAR(50),
    message TEXT,
    error_code VARCHAR(50),
    error_message TEXT,
    file_checksum TEXT,
    file_size_actual BIGINT,
    upload_duration_ms INT,
    bandwidth_kbps DECIMAL(10, 2),
    reported_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_upload_reports_file_id ON upload_reports(file_id);
CREATE INDEX idx_upload_reports_owner_user_id ON upload_reports(owner_user_id);
CREATE INDEX idx_upload_reports_status ON upload_reports(status);
CREATE INDEX idx_upload_reports_reported_at ON upload_reports(reported_at DESC);
```

## 🛠️ Setup Commands

### 1. Run Migration
```bash
psql -U user -d database -f migrations/003_create_upload_reports_table.sql
```

### 2. Test Connection
```bash
psql -U user -d database -c "SELECT * FROM upload_reports LIMIT 1;"
```

### 3. Run Unit Tests
```bash
go test -v ./internal/files/...
```

### 4. Run Example Client
```bash
# Report
go run cmd/example-client/main.go -cmd=report -file-id=1

# Get report
go run cmd/example-client/main.go -cmd=get -file-id=1

# List reports
go run cmd/example-client/main.go -cmd=list
```

## ✅ Validation Rules

### status
- Required: Yes
- Type: String (enum)
- Values: `completed`, `failed`

### report_type
- Required: Yes
- Type: String (enum)
- Values: `success`, `error`, `warning`

### message
- Required: No
- Type: String
- Max Length: 1000

### error_code
- Required: No
- Type: String
- Max Length: 50

### file_checksum
- Required: No
- Type: String
- Length: 8-128

### file_size_actual
- Required: No
- Type: Integer
- Min: 0

### upload_duration_ms
- Required: No
- Type: Integer
- Min: 0

### bandwidth_kbps
- Required: No
- Type: Float
- Min: 0

## 🔐 Security Checklist

- [x] Telegram authentication required
- [x] File ownership verification
- [x] Input validation for all fields
- [x] SQL injection prevention (parameterized queries)
- [x] Error messages don't leak sensitive data
- [x] Foreign key constraints
- [x] Unique constraints on critical data
- [x] Proper HTTP status codes
- [x] Rate limiting ready (can add middleware)

## 📊 Response Status Codes

| Code | Scenario |
|------|----------|
| 200 | GET successful |
| 201 | POST successful (resource created) |
| 400 | Bad request (validation error, missing headers) |
| 404 | Resource not found (file, report) |
| 500 | Server error (database, internal error) |

## 🎯 Response JSON Formats

### Success Response (POST)
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

### Error Response
```json
{
  "code": "VALIDATION_ERROR",
  "message": "Request validation failed",
  "errors": [
    {"field": "status", "message": "status is required"}
  ]
}
```

### Report Response (GET)
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

## 🚀 Integration Steps

1. **Database**: Run migration
2. **Models**: Add upload_report.go, update user.go
3. **Storage**: Add file.go repository
4. **Service**: Add service.go and error.go
5. **HTTP**: Update handlers, validation, error_handler
6. **API Docs**: Update openapi.yaml
7. **Routes**: Register endpoints in main.go
8. **Test**: Run tests and verify with curl
9. **Deploy**: Push to dev branch

## 📞 Common Issues & Solutions

| Issue | Solution |
|-------|----------|
| "database error" | Check DB connection and migration |
| "file not found" | Create file first, verify file_id |
| "missing headers" | Add X-Telegram-User-Id and X-Telegram-Username |
| "validation error" | Check required fields and enum values |
| "access denied" | Verify user owns the file |

## 📈 Performance Considerations

- Indexes on `file_id`, `owner_user_id`, `reported_at` for fast queries
- Pagination (limit: 20 default, 100 max) to avoid large result sets
- Database connection pooling recommended
- Consider caching for frequently accessed reports

## 🔄 Business Logic Flow

```
1. User uploads file
   → POST /v1/files
   ← file_id returned

2. File upload completes
   → POST /v1/files/{file_id}/report-complete
   → Service validates request
   → Service creates upload_reports record
   → Service updates files.status
   ← report_id returned

3. Check report status
   → GET /v1/files/{file_id}/report
   ← Latest report details returned

4. List all reports
   → GET /v1/upload-reports?limit=20
   ← Paginated list returned
```

## 💡 Tips & Best Practices

1. Always include Telegram headers in requests
2. Use pagination when listing reports
3. Include file_checksum for data integrity verification
4. Monitor upload_duration_ms and bandwidth_kbps for performance analysis
5. Use appropriate error reporting (success/error/warning types)
6. Keep error messages meaningful for debugging
7. Test with both successful and failed upload scenarios
8. Use the example client for quick testing

---

**Last Updated**: November 12, 2025
**Version**: 1.0.0-rc.1
