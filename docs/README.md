# Mock Conversation Flow cho Bot Telegram

### 1. Giới thiệu

### 2. Danh sách Command của Bot

### 3. Mock Conversation Flow

##### 3.1. Flow /start

(Hội thoại mẫu + hành vi dự kiến)



**1. Mục tiêu của flow**
Khởi tạo phiên làm việc. Backend kiểm tra định danh người dùng từ Telegram:

  * Nếu chưa tồn tại: Tự động tạo mới (Insert DB).
  * Nếu đã tồn tại: Trả về thông tin profile hiện tại.

**2. Điều kiện kích hoạt (Trigger)**

  * **User:** Gõ lệnh `/start`
  * **Backend Endpoint:** `GET /v1/me`

**3. Các actor liên quan**

  * User (Sender)
  * Bot Telegram (FE)
  * Backend API

**4. Conversation Flow (Chi tiết)**

  * **Bước 1 – User gửi lệnh**

      * **User:** `/start`
      * **Bot:** (Nhận event, trích xuất `telegram_id` và `username`).

  * **Bước 2 – Bot gọi API (Get or Create)**

      * **API Request (FE -\> BE):**
          * Endpoint: `GET /v1/me`
          * Headers (Định danh):
            ```http
            X-Telegram-ID: <telegram_id>
            X-Telegram-Username: <username>
            ```
      * **Backend Processing:**
          * Query bảng `users` theo `telegram_id`.
          * Nếu không có dữ liệu: `INSERT` bản ghi mới.
          * Nếu có dữ liệu: `SELECT` bản ghi cũ (kèm update username nếu thay đổi).
      * **Backend Response (JSON):**
        ```json
        {
          "data": {
            "id": 1,
            "telegram_id": 123456789,
            "username": "nguyen_van_a",
            "created_at": "2025-11-18T10:00:00Z",
            "status": "active"
          }
        }
        ```

  * **Bước 3 – Bot phản hồi User**

      * **Bot:**
        > Xin chào **nguyen\_van\_a**\! 👋
        > Tài khoản của bạn đã được kích hoạt.
        > Bạn có thể gửi file vào đây để upload hoặc gõ /help để xem hướng dẫn.

**5. Error Handling**

| Tình huống | Bot phản hồi | Backend trả về |
| :--- | :--- | :--- |
| **Thiếu Header** (Lỗi code FE) | "Lỗi hệ thống: Không xác định được danh tính." | **400 Bad Request**<br>`{ "error": "MISSING_TELEGRAM_INFO" }` |
| **Lỗi Database** (Connect fail) | "Hệ thống đang bận, vui lòng thử lại sau." | **500 Internal Server Error**<br>`{ "error": "INTERNAL_ERROR" }` |





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
