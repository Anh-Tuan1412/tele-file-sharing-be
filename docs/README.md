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

1. Mục tiêu của flow:

Cho phép người dùng tạo một chia sẻ cho một file mà mình đã upload lên hệ thống

2. Điều kiện kích hoạt:

Bot lệnh:

/share

Backend liên quan:

POST /v1/shares

3. Các actor liên quan:

User(Sender)

Bot Telegram

Backend API 

4. Conservation flow

4.1. Bước 1 - User gửi command

User:

/share

Bot:

Hãy chọn một file mà bạn muốn tạo chia sẻ

4.2 Bước 2 - User chọn file

User:

Chọn một file

Bot:

Hãy nhập ngày bắt đầu và ngày kết thúc hiệu lực của link chia sẻ của file

4.3 Bước 3 - User nhập ngày

User:

Nhập ngày bắt đầu và ngày hết hiệu lực

Bot:

Bạn có muốn đặt mật khẩu không? Nhập "Có"/"Không"

4.4 Bước 4 - User đồng ý 

User:

Nhập "Có"

Bot: 

Hãy nhập mật khẩu mà bạn muốn

4.5 Bước 5 - User nhập mật khẩu

User:

Nhập mật khẩu

Bot:

Hãy nhập username Telegram của những người được phép tải file 

4.6 Bước 6 - User nhập thông tin

User:

Nhập các username Telegram

Bot:

Gọi API:

POST /v1/shares

{

  "file_id": 1234,

  "from_ts": "2025-11-05T00:00:00Z",

  "to_ts": "2025-11-10T00::00:00Z",

  "require_password": true,

  "password": "1234",

  "recipients": ["@huytran", "@anle"]

}

Backend trả về:

{

  "share_id": "1324",

  "Link": "https://api.fileshare.com/s/1324",

  "recipients": ["@huytran", "@anle"]

}

Bot:

Đây là link chia sẻ của bạn: <u>Link</u>↗️. Chỉ có @huytran và @anle mới tải được.

5. Error Handling:

|Tình huống|Bot phản hồi|Backend trả|
|:---|:---|:---|
|Dữ liệu không hợp lệ|"Dữ liệu không hợp lệ. Vui lòng kiểm tra và nhập lại"|400 BAD REQUEST<br>{"error: INVALID_INPUT"}|
|Lỗi hệ thống(máy chủ không truy cập được database hệ thống, backend gặp lỗi,...)|"Máy chủ xảy ra vấn đề. Vui lòng quay lại sau"|500 INTERNAL SERVER ERROR<br>{"error: INTERNAL_ERROR"}|



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
