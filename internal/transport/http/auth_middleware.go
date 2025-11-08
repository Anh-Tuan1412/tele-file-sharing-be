package http

import (
	"net/http"
	"strconv"
	"file-sharing/internal/storage"
	"file-sharing/internal/model"
	"github.com/gin-gonic/gin"
)

// ContextUserKey là key để lưu trữ thông tin người dùng trong context của Gin.
const ContextUserKey = "user"

// AuthMiddleware tạo ra một middleware cho Gin để xác thực người dùng.
func AuthMiddleware(userRepo storage.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		telegramIDStr := c.GetHeader("X-Telegram-User-Id")
		telegramUsername := c.GetHeader("X-Telegram-Username")

		if telegramIDStr == "" || telegramUsername == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Missing or invalid headers"})
			return
		}

		telegramID, err := strconv.ParseInt(telegramIDStr, 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid Telegram ID"})
			return
		}

		user, err := userRepo.UpsertByTelegram(c.Request.Context(), telegramID, telegramUsername)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to authenticate user"})
			return
		}

		c.Set(ContextUserKey, user)
		c.Next()
	}
}

// GetUserFromContext là một hàm helper để lấy thông tin user từ Gin context.
func GetUserFromContext(c *gin.Context) (*model.User, bool) {
	// Get trả về interface{}, nên cần ép kiểu an toàn
	user, exists := c.Get(ContextUserKey)
	if !exists {
		return nil, false
	}
	
	// Ép kiểu an toàn
	userTyped, ok := user.(*model.User)
	if !ok {
		return nil, false
	}

	return userTyped, true
}