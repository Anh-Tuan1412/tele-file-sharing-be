# Mock Conversation Flow cho Bot Telegram

### 1. Giới thiệu

### 2. Danh sách Command của Bot
## 2. Danh sách Command của Bot

Các command chính mà người dùng có thể sử dụng để tương tác với Bot, cùng với mô tả chức năng của chúng.

| Command | Mô tả | API Backend liên quan |
| :--- | :--- | :--- |
| **`/start`** | Bắt đầu phiên làm việc, xác thực người dùng với hệ thống. | `POST v1/api/me` |
| **`/upload`** | Kích hoạt luồng tải file mới lên. | `POST /v1/files` |
| **`/myfiles`** | Liệt kê tất cả các file mà người dùng đã tải lên. | `GET /v1/files` |
| **`/share`** | Bắt đầu luồng tạo link chia sẻ cho một file đã có. | `POST /v1/shares` |
| **`/myshares`** | Liệt kê các link chia sẻ mà người dùng đã tạo. | `GET /v1/shares` |
| **`/revoke`** | Bắt đầu luồng thu hồi (vô hiệu hóa) một link chia sẻ. | `POST /v1/shares/:id/revoke` |

### 3. Mock Conversation Flow

##### 3.1. Flow /start


###### 1\. Mục tiêu của flow

Khởi tạo phiên làm việc cho người dùng.

  * **Logic:** Backend kiểm tra thông tin từ Telegram. Nếu user chưa tồn tại trong bảng `users`, hệ thống tự động tạo mới (Auto-register). Nếu đã tồn tại, trả về thông tin hiện có.

###### 2\. Điều kiện kích hoạt (Trigger)

  * **User:** Gõ lệnh `/start`
  * **Backend Endpoint:** `GET /v1/me`

###### 3\. Các actor liên quan

  * User (Sender)
  * Bot Telegram (FE Logic)
  * Backend API

###### 4\. Conversation Flow (Chi tiết)

**4.1. Bước 1 – User bắt đầu**
**User:**
`/start`

**Bot:**
*(Bot nhận event, trích xuất `telegram_id` và `username` từ message của người dùng)*

**4.2. Bước 2 – Bot gọi API (Get or Create User)**
**API gọi (FE -\> BE):**

  * **Endpoint:** `GET /v1/me`
  * **Headers** (Dùng để định danh user):
    ```http
    X-Telegram-ID: 123456789
    X-Telegram-Username: nguyen_van_a
    ```

**Backend xử lý (Logic):**

1.  Query bảng `users` theo `telegram_id`.
2.  **Nếu không tồn tại:** `INSERT INTO users (telegram_id, username, created_at) ...`
3.  **Nếu tồn tại:** `SELECT * FROM users ...` (có thể update `username` nếu thay đổi).

**Backend trả về (Response 200 OK):**

```json
{
  "data": {
    "id": 1,
    "telegram_id": 123456789,
    "username": "nguyen_van_a",
    "created_at": "2025-11-18T10:00:00Z"
  }
}
```

**4.3. Bước 3 – Bot phản hồi User**

**Bot:**

> Xin chào **nguyen\_van\_a**\! 👋
> Tài khoản của bạn đã sẵn sàng.
>
> Bạn có thể:
> 📤 Gửi file trực tiếp để upload.
> ❓ Gõ /help để xem hướng dẫn.

###### 5\. Error Handling

| Tình huống | Bot phản hồi | Backend trả về |
| :--- | :--- | :--- |
| **Thiếu thông tin định danh** (Header rỗng) | "Lỗi hệ thống: Không xác định được người dùng." | **400 Bad Request**<br>`{ "error": "MISSING_TELEGRAM_INFO" }` |
| **Lỗi Database** (Connect/Insert fail) | "Hệ thống đang bận. Vui lòng thử lại sau." | **500 Internal Server Error**<br>`{ "error": "INTERNAL_ERROR" }` |



##### 3.2. Flow /upload

(Hội thoại mẫu + bước gửi file)

##### 3.3. Flow POST /v1/files + /complete

(Luồng hoàn tất upload)

##### 3.4. Flow /myfiles

###### 1\. Mục tiêu của flow

Hiển thị danh sách tất cả các file mà người dùng đã upload lên hệ thống, giúp người dùng dễ dàng quản lý và tạo link chia sẻ cho các file của mình.

###### 2\. Điều kiện kích hoạt (Trigger)

  * **Bot lệnh:** `/myfiles`
  * **Backend liên quan:** `GET /api/v1/files`

###### 3\. Các actor liên quan

  * User (Sender)
  * Bot Telegram (FE)
  * Backend API

###### 4\. Conversation Flow (Hội thoại chi tiết)

**4.1. Bước 1 – User gửi command**

**User:**
`/myfiles`

**Bot:**
```
⏳ Đang tải danh sách file của bạn...
```

**4.2. Bước 2 – Bot gọi API**

**API gọi:**
```http
GET /api/v1/files
Headers:
  X-Telegram-User-Id: 123456789
  X-Telegram-Username: nguyen_van_a
```

**Backend trả về (Response 200 OK - Có files):**

```json
[
  {
    "id": 101,
    "object_key": "uploads/5/1731312000_document.pdf",
    "filename": "document.pdf",
    "size": 2621440,
    "mime": "application/pdf",
    "status": "completed",
    "created_at": "2025-11-20T10:30:00Z",
    "updated_at": "2025-11-20T10:30:00Z"
  },
  {
    "id": 102,
    "object_key": "uploads/5/1731311100_image.jpg",
    "filename": "image.jpg",
    "size": 1258291,
    "mime": "image/jpeg",
    "status": "completed",
    "created_at": "2025-11-20T09:15:00Z",
    "updated_at": "2025-11-20T09:15:00Z"
  },
  {
    "id": 103,
    "object_key": "uploads/5/1731225900_presentation.pptx",
    "filename": "presentation.pptx",
    "size": 6082560,
    "mime": "application/vnd.openxmlformats-officedocument.presentationml.presentation",
    "status": "completed",
    "created_at": "2025-11-19T16:45:00Z",
    "updated_at": "2025-11-19T16:45:00Z"
  }
]
```

**4.3. Bước 3 – Bot hiển thị danh sách**

**Bot:**

```
📁 Danh sách file của bạn:

1️⃣ document.pdf
   📊 Size: 2.5 MB
   📅 Upload: 20/11/2025 10:30
   🔗 /share_101

2️⃣ image.jpg
   📊 Size: 1.2 MB
   📅 Upload: 20/11/2025 09:15
   🔗 /share_102

3️⃣ presentation.pptx
   📊 Size: 5.8 MB
   📅 Upload: 19/11/2025 16:45
   🔗 /share_103

💡 Gõ /share_<số> để tạo link chia sẻ
📝 Tổng: 3 files (9.5 MB)
```

###### 5\. Error Handling

| Tình huống | Bot phản hồi | Backend trả về |
| :--- | :--- | :--- |
| **Thiếu thông tin định danh** | "❌ Lỗi: Không xác định được người dùng. Vui lòng gõ /start để đăng nhập." | **400 Bad Request**<br>`{ "error": "missing Telegram headers" }` |
| **User chưa tồn tại** | "❌ Tài khoản chưa được khởi tạo. Vui lòng gõ /start trước." | **401 Unauthorized**<br>`{ "error": "user not found" }` |
| **Lỗi Database** | "❌ Hệ thống đang gặp sự cố. Vui lòng thử lại sau." | **500 Internal Server Error**<br>`{ "error": "database error" }` |
| **Timeout kết nối** | "❌ Không thể kết nối đến server. Vui lòng kiểm tra mạng và thử lại." | **503 Service Unavailable**<br>`{ "error": "connection timeout" }` |

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
