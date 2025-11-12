package model

import "time"

type User struct {
	ID           int64     `json:"id"`
	TelegramID   int64     `json:"telegram_id" db:"telegram_user_id"`
	Username     string    `json:"username"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// FileWithOwner represents a file with its owner information
type FileWithOwner struct {
	ID            int64     `json:"id" db:"id"`
	OwnerUserID   int64     `json:"owner_user_id" db:"owner_user_id"`
	ObjectKey     string    `json:"object_key" db:"object_key"`
	Filename      string    `json:"filename" db:"filename"`
	Size          int64     `json:"size" db:"size"`
	Mime          string    `json:"mime" db:"mime"`
	Status        string    `json:"status" db:"status"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
	TelegramID    int64     `json:"telegram_id" db:"telegram_user_id"`
	Username      string    `json:"username" db:"username"`
}