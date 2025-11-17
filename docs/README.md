# Mock Conversation Flow cho Bot Telegram

### 1. Danh sách Command của Bot

### 2. Danh sách conversation flow

| STT | Tên Flow                 | API Chính                                                       | Giải thích Flow                                                                                                                                               |
| --- | ------------------------ | --------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 1   | Flow `/start`            | GET `/v1/api/me`                                                | User bắt đầu tương tác với bot. Bot gửi lời chào và FE gọi API để đăng ký user mới nếu chưa tồn tại.                                                          |
| 2   | Flow `/upload`           | POST `/v1/files`                                                | Người dùng yêu cầu upload file. Bot yêu cầu user gửi file và FE tạo bản ghi file ở trạng thái “pending”.                                                      |
| 3   | Flow `/upload complete`  | POST `/v1/files/:id/complete`                                   | Sau khi file đã được bot nhận và lưu trữ, FE gọi API này để đánh dấu file đã upload xong (complete).                                                          |
| 4   | Flow `/myfiles`          | GET `/v1/files`                                                 | Người dùng muốn xem danh sách file họ đã upload. Bot gọi API lấy danh sách và hiển thị cho user.                                                              |
| 5   | Flow `/share`            | POST `/v1/shares`                                               | Người dùng chọn file → nhập password (nếu có) → đặt hạn dùng. FE gửi dữ liệu này lên BE để tạo link chia sẻ.                                                  |
| 6   | Flow `/myshares`         | GET `/v1/shares`                                                | Người dùng muốn xem danh sách các link chia sẻ đã tạo. Bot gọi API để liệt kê và hiển thị cho user.                                                           |
| 7   | Flow `/revoke`           | POST `/v1/shares/:id/revoke`                                    | Người dùng chọn một link chia sẻ muốn thu hồi. FE gọi API revoke để BE đánh dấu share bị vô hiệu hóa.                                                         |
| 8   | Flow tải file (metadata) | GET `/v1/shares/:id`                                            | Người nhận click link chia sẻ → Bot cần kiểm tra metadata link (có tồn tại không, hết hạn hay chưa, có yêu cầu password không).                               |
| 9   | Flow password + download | POST `/v1/shares/:id/authorize` + GET `/v1/shares/:id/download` | Nếu link có password: bot yêu cầu nhập pass → FE gọi authorize. Nếu hợp lệ → FE gọi API download để lấy presigned URL hoặc gửi file trực tiếp cho người nhận. |

(Liệt kê các command bot hỗ trợ và mô tả ngắn)

### 3. Mock Conversation Flow

##### 3.1. Flow /start

(Hội thoại mẫu + hành vi dự kiến)

##### 3.2. Flow /upload

(Hội thoại mẫu + bước gửi file)

##### 3.3. Flow POST /v1/files + /complete

(Luồng hoàn tất upload)

##### 3.4. Flow /myfiles

(Luồng FE gọi GET /v1/files)

##### 3.5. Flow /share

(Luồng chọn file → đặt password → thời hạn → recipients)

##### 3.6. Flow /myshares

(Luồng FE gọi GET /v1/shares)

##### 3.7. Flow /revoke

(Luồng gọi POST /v1/shares/:id/revoke)

##### 3.8. Flow người nhận tải file

(Luồng GET /v1/shares/:id → authorize → download)

### 4. State Diagram cho Bot

(Sơ đồ FSM liệt kê tất cả các trạng thái bot từ upload, share, nhập pass…)

### 5. API mapping

## 5. API Mapping (Bot → Backend)

| Bot Command / Hành động người dùng | API Backend được gọi            |
| ---------------------------------- | ------------------------------- |
| Người dùng start                   | `POST v1/api/me`                |
| `/upload`                          | `POST /v1/files`                |
| (User gửi file cho bot)            | `POST /v1/files/:id/complete`   |
| `/myfiles`                         | `GET /v1/files`                 |
| `/share`                           | `POST /v1/shares`               |
| `/myshares`                        | `GET /v1/shares`                |
| `/revoke`                          | `POST /v1/shares/:id/revoke`    |
| Người nhận mở link share           | `GET /v1/shares/:id`            |
| Người nhận nhập password (nếu có)  | `POST /v1/shares/:id/authorize` |
| Người nhận tải file                | `GET /v1/shares/:id/download`   |
