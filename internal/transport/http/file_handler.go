package http

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type FileInitRequest struct {
	Filename  string    `json:"filename" binding:"required"`
	Size      int64     `json:"size" binding:"required"`
	MimeType  string    `json:"mime_type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type FileInitResponse struct {
	FileID    int64     `json:"file_id"`
	ObjectKey string    `json:"object_key"`
	UploadURL string    `json:"upload_url,omitempty"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// generateRandomToken generates a random 16-byte hex string
func generateRandomToken() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano()) // fallback
	}
	return hex.EncodeToString(b)
}

// POST /v1/files
func InitFileUploadHandler(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		telegramID := c.GetHeader("X-Telegram-User-Id")
		username := c.GetHeader("X-Telegram-Username")

		if telegramID == "" || username == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing Telegram headers"})
			return
		}

		var req FileInitRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Step 1: find or create user
		var userID int64
		err := db.QueryRow(`SELECT id FROM users WHERE telegram_user_id = $1`, telegramID).Scan(&userID)
		if err == sql.ErrNoRows {
			err = db.QueryRow(`
				INSERT INTO users (telegram_user_id, username)
				VALUES ($1, $2)
				RETURNING id
			`, telegramID, username).Scan(&userID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
				return
			}
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
			return
		}

		// Step 2: create file record
		objectKey := fmt.Sprintf("uploads/%d/%d_%s", userID, time.Now().Unix(), req.Filename)
		uploadToken := generateRandomToken() // replaced uuid.New().String()

		var fileID int64
		err = db.QueryRow(`
			INSERT INTO files (owner_user_id, object_key, filename, size, mime, status, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, 'pending', $6, $7)
			RETURNING id
		`, userID, objectKey, req.Filename, req.Size, req.MimeType, req.CreatedAt, req.UpdatedAt).Scan(&fileID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create file record"})
			return
		}

		// Step 3: (optional) generate presigned URL from S3 / MinIO
		uploadURL := fmt.Sprintf("https://minio.local/presigned/%s", uploadToken)

		resp := FileInitResponse{
			FileID:    fileID,
			ObjectKey: objectKey,
			UploadURL: uploadURL,
			Status:    "pending",
			CreatedAt: req.CreatedAt,
			UpdatedAt: req.UpdatedAt,
		}
		c.JSON(http.StatusOK, resp)
	}
}
