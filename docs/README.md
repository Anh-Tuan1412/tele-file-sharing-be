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
1. Mục tiêu của flow

Cho phép người nhận (recipient) tải xuống file được chia sẻ qua link.
Bot thực hiện:

Lấy metadata của share link

Kiểm tra yêu cầu truy cập (hết hạn, thu hồi, password, TOTP, whitelist user)

Nếu đạt điều kiện → gọi API BE để lấy presigned URL

Gửi link tải xuống cho người nhận

2. Điều kiện kích hoạt (Trigger)

User: Nhấn vào link được chia sẻ → Bot nhận deep-link dạng:

/start dl_SHARETOKEN


Bot FE trích ra SHARE_TOKEN và gọi API BE tương ứng:

GET /v1/shares/:id

POST /v1/shares/:id/verify-password (nếu cần)

POST /v1/shares/:id/verify-totp (nếu cần)

GET /v1/shares/:id/download

3. Các actor liên quan

User (Recipient – người nhận file)

Bot Telegram (FE – xử lý UI + gọi API)

Backend API

Object Storage (MinIO/S3)

4. Conversation Flow (Chi tiết)
4.1. Bước 1 – User nhấn vào link chia sẻ

User:
(Nhấn vào link: https://domain.com/dl/abcd1234
 → Telegram mở bot với payload)

/start dl_abcd1234


Bot:
(Trích xuất token = "abcd1234")

4.2. Bước 2 – Bot gọi API kiểm tra share

API gọi (FE → BE):

Endpoint:

GET /v1/shares/abcd1234


Headers:

X-Telegram-ID: 987654321
X-Telegram-Username: user_b


Backend xử lý (Logic):

Tìm share theo share_token.

Kiểm tra:

file có tồn tại không

link có hết hạn không

link có bị revoke không

user nhận có nằm trong whitelist không (nếu có)

có cần password / TOTP hay không

Backend trả về (ví dụ link hợp lệ – không password – không TOTP):

{
  "data": {
    "share_id": "abcd1234",
    "file_name": "document.pdf",
    "file_size": 1340023,
    "expires_in": 7200,
    "password_required": false,
    "totp_required": false
  }
}


Bot gửi User:

 document.pdf (1.34 MB)
 Link còn hiệu lực trong 2 giờ.

Bạn muốn tải xuống chứ?
 Chọn Tải xuống.

4.3. Bước 3 – Người dùng nhấn “Tải xuống”

User:
(Bấm button “Tải xuống”)

API gọi (FE → BE):

GET /v1/shares/abcd1234/download


Backend xử lý:

Ghi audit log “USER_DOWNLOADED_FILE”

Tạo presigned URL từ S3/MinIO

Backend trả về:

{
  "download_url": "https://s3-presigned-url..."
}


Bot phản hồi User:

 File sẵn sàng!
 Đây là link tải của bạn:
https://s3-presigned-url
...

Nhánh bổ sung theo rule (nếu có)
4.4. Nếu share yêu cầu PASSWORD

Backend (Step 2) trả:

{
  "password_required": true
}


Bot:

 File này được bảo vệ bằng mật khẩu.
Vui lòng nhập mật khẩu.

User:

mypassword


Bot → API:

POST /v1/shares/abcd1234/verify-password
{
  "password": "mypassword"
}


Nếu đúng:

{ "status": "ok" }


Bot:

 Mật khẩu chính xác.
Bấm Tải xuống để nhận file.

4.5. Nếu share yêu cầu TOTP

Backend trả:

{
  "totp_required": true
}


Bot:

 File yêu cầu mã xác thực 2 lớp (TOTP).
Vui lòng nhập mã 6 số.

User:

123456


Bot → API:

POST /v1/shares/abcd1234/verify-totp
{
  "code": "123456"
}


Nếu đúng:

{ "status": "ok" }


Bot:

 Xác thực thành công.
Bấm Tải xuống để nhận file.

4.6. Nếu share hết hạn

Backend trả:

{ "error": "LINK_EXPIRED" }


Bot:

 Link đã hết hạn.
Vui lòng yêu cầu người gửi tạo link mới.

4.7. Nếu share bị thu hồi

Backend trả:

{ "error": "LINK_REVOKED" }


Bot:

 Link này đã bị thu hồi bởi người gửi.

4.8. Nếu user không nằm trong danh sách allowed recipients

Backend trả:

{ "error": "USER_NOT_ALLOWED" }


Bot:

 Bạn không có quyền tải file này.

4.9. Nếu file đã bị xóa

Backend trả:

{ "error": "FILE_NOT_FOUND" }


Bot:

 File không còn tồn tại trong hệ thống.

5. Error Handling
Tình huống|	Bot phản hồi|	Backend trả về
| :--- | :--- | :--- |
Sai password|	“ Mật khẩu không đúng.”|	INVALID_PASSWORD
Sai TOTP|	“ Mã xác thực không hợp lệ.”|	INVALID_TOTP
Link hết hạn|	“ Link đã hết hạn.”|	LINK_EXPIRED
Không được phép|	“ Bạn không có quyền.”|	USER_NOT_ALLOWED
File bị xóa|	“File không tồn tại.”|	FILE_NOT_FOUND
Lỗi hệ thống|	“Hệ thống đang bận, thử lại sau.”|	INTERNAL_ERROR

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
