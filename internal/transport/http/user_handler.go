package http

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

// UserHandler là một struct để quản lý các handler cho user.
type UserHandler struct {
	// Trong tương lai, có thể inject một "user service" vào đây
}

// NewUserHandler tạo một UserHandler mới.
func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

// GetCurrentUser lấy thông tin user từ context trả về
// Handler này chỉ được gọi sau khi AuthMiddleware được gọi.
func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	// Lấy user từ context
	user, exists := GetUserFromContext(c)
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	// Trả về thông tin user
	c.JSON(http.StatusOK, user)
}