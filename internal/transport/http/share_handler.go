package http

import (
	"net/http"
	"strconv"

	"file-sharing/internal/share"
	"github.com/gin-gonic/gin"
)

type ShareHandler struct {
	service share.Service
}

func NewShareHandler(service share.Service) *ShareHandler {
	return &ShareHandler{
		service: service,
	}
}

// POST /v1/shares/:id/revoke
func (h *ShareHandler) HandleRevoke(c *gin.Context) {
	// Lấy User hiện tại từ Context (do AuthMiddleware nạp vào)
	user, exists := GetUserFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Lấy Share ID từ URL
	shareIDStr := c.Param("id")
	shareID, err := strconv.ParseInt(shareIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid share ID"})
		return
	}

	// Gọi Service
	err = h.service.RevokeShare(c.Request.Context(), shareID, user.ID)
	if err != nil {
        // Nếu lỗi là do không tìm thấy hoặc không đúng chủ sở hữu
		if err.Error() == "share not found or access denied" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	// Trả về thành công
	c.JSON(http.StatusOK, gin.H{
		"message":  "Share revoked successfully",
		"share_id": shareID,
		"status":   "revoked",
	})
}