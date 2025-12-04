package model

import "time"

// File đại diện cho thông tin metadata của file được lưu trữ
type File struct {
	ID          int64     `json:"id" db:"id"`
	OwnerUserID int64     `json:"owner_user_id" db:"owner_user_id"`
	ObjectKey   string    `json:"object_key" db:"object_key"` // Khóa duy nhất trên MinIO
	Filename    string    `json:"filename" db:"filename"`
	Size        int64     `json:"size" db:"size"`
	Mime        string    `json:"mime" db:"mime"`
	Status      string    `json:"status" db:"status"` // Trạng thái: pending, uploaded, failed
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Hằng số trạng thái cho File
const (
	FileStatusPending  = "pending"  // Khởi tạo upload thành công, chờ client upload
	FileStatusUploaded = "uploaded" // File đã được upload thành công lên MinIO
	FileStatusFailed   = "failed"   // Upload thất bại hoặc có lỗi
)
