package model

import "time"

type User struct {
	ID           int64     `json:"id"`
	TelegramID   int64     `json:"telegram_id" db:"telegram_user_id"`
	Username     string    `json:"username"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}