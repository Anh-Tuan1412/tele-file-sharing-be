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

(Luồng FE gọi GET /v1/files)

##### 3.5. Flow /share

(Luồng chọn file → đặt password → thời hạn → recipients)

##### 3.6. Flow /myshares
###### 1\. Mục tiêu của flow

Liệt kê các link chia sẻ mà người dùng đã tạo.

* **Logic:** FE gọi GET /v1/shares → Backend lấy danh sách shares theo telegram_id → trả về danh sách theo thứ tự mới nhất trước.

###### 2\. Điều kiện kích hoạt (Trigger)

  * **User:** Gõ lệnh `/myshares`
  * **Backend Endpoint:**  `GET /v1/shares`

  ###### 3\. Các actor liên quan

  * User (Sender)
  * Bot Telegram (FE Logic)
  * Backend API

###### 4\. Conversation Flow (Chi tiết)

**4.1. Bước 1 – User yêu cầu xem danh sách share**
**User:**
`/myshares`

**Bot:**
*(Bot nhận event, trích xuất `telegram_id` và `username` từ message của người dùng)*

**4.2. Bot gọi API lấy danh sách shares**
**API gọi (FE -\> BE):**

  * **Endpoint:** `GET /v1/shares`
  * **Headers** (Dùng để định danh user):
    ```http
    X-Telegram-ID: 123456789
    X-Telegram-Username: nguyen_van_a
    ```

**Backend xử lý (Logic):**

1.  Validate headers → nếu thiếu → trả lỗi 401.
2.  Query bảng shares theo `owner_telegram_id`.
3.  Áp dụng limit + offset.
4.  Sort theo `created_at DESC`.
5.  Trả list về FE.

**Backend trả về (Response 200 OK):**

```json
[
  {
    "id": 10,
    "file_id": 101,
    "owner_user_id": 5,
    "hash": "a1b2c3d4",
    "require_password": true,
    "hash_password": "$2a$10$....",
    "revoked": false,
    "expires_at": null,
    "created_at": "2025-11-18T12:00:00Z",
    "updated_at": "2025-11-18T12:00:00Z"
  },
  {
    "id": 8,
    "file_id": 95,
    "owner_user_id": 5,
    "hash": "xyz123",
    "require_password": false,
    "hash_password": null,
    "revoked": true,
    "expires_at": "2025-11-20T00:00:00Z",
    "created_at": "2025-11-17T09:20:00Z",
    "updated_at": "2025-11-18T10:30:00Z"
  }
]
```

**4.3. Bước 3 – Bot phản hồi User**

**Bot:**

**Nếu có kết quả:**

> 📂 Danh sách link bạn đã chia sẻ:
> 1️⃣ baocao.pdf – https://abc.xyz/s/xyz123
> 2️⃣ anh.png – (đã thu hồi)

**Nếu danh sách trống:**

> Bạn chưa tạo link chia sẻ nào.
> Hãy gửi một file để bắt đầu chia sẻ nhé! 📤

###### 5\. Error Handling

| Tình huống | Bot phản hồi | Backend trả về |
| :--- | :--- | :--- |
| **Thiếu thông tin định danh** (Header rỗng) | "Lỗi hệ thống: Không xác định được người dùng." | **400 Bad Request**<br>`{ "error": "MISSING_TELEGRAM_INFO" }` |
| **Lỗi Database** (Connect/Insert fail) | "Hệ thống đang bận. Vui lòng thử lại sau." | **500 Internal Server Error**<br>`{ "error": "INTERNAL_ERROR" }` |
 



##### 3.7. Flow /revoke
###### 1\. Mục tiêu của flow

Thu hồi (vô hiệu hóa) một link chia sẻ.

###### 2\. Điều kiện kích hoạt (Trigger)

  * **User:** Gõ lệnh `/revoke <id>`
  * **Backend Endpoint:**  `POST /v1/shares/:id/revoke`

  ###### 3\. Các actor liên quan

  * User 
  * Bot Telegram (FE Logic)
  * Backend API

  ###### 4\. Conversation Flow (Chi tiết)

**4.1. Bước 1 – User yêu cầu thu hồi share**
**User:**
 `/revoke 10`

 **Bot:**
*(Parse id = 10, chuẩn bị gọi API)*

**4.2. Bước 2 – Bot gọi API revoke share**
**API gọi (FE -\> BE):**

  * **Endpoint:** `POST /v1/shares/:id/revoke`
  * **Headers** (Dùng để định danh user):
    ```http
    X-Telegram-ID: 123456789
    X-Telegram-Username: nguyen_van_a
    ```

**Backend xử lý (Logic):**

1.  Validate headers.
2.  Lấy share theo id.
3.  Kiểm tra chủ sở hữu.
4.  Nếu không phải chủ sở hữu → trả lỗi 404.
5.  Nếu share đã bị revoke → vẫn trả 200.
6.  Set `revoked = true`, update DB.
7.  Trả về object kết quả.

**Backend trả về (Response 200 OK):**

```json
{
  "message": "Share revoked successfully",
  "share_id": 10,
  "status": "revoked"
}
```

**4.3. Bước 3 – Bot phản hồi User**

**Bot:**

> 🔒 Link chia sẻ #10 đã được thu hồi thành công.
> Người khác sẽ không thể truy cập link này nữa.

**Nếu người dùng không sở hữu link:**

> ❌ Bạn không có quyền thu hồi link này.

###### 5\. Error Handling

| Tình huống | Bot phản hồi | Backend trả về |
| :--- | :--- | :--- |
| **Thiếu thông tin định danh** (Header rỗng) | "Lỗi hệ thống: Không xác định được người dùng." | **400 Bad Request**<br>`{ "error": "MISSING_TELEGRAM_INFO" }` |
| **Share không tồn tại hoặc không thuộc sở hữu** | "Không tìm thấy link hoặc bạn không có quyền thu hồi."| **404 Not Found**<br>`{ "error": "NOT_FOUND" }` |
| **Lỗi Database** (Connect/Insert fail) | "Hệ thống đang bận. Vui lòng thử lại sau." | **500 Internal Server Error**<br>`{ "error": "INTERNAL_ERROR" }` |


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
