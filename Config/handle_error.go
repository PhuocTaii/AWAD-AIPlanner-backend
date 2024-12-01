package config

import (
	"github.com/gin-gonic/gin"
)

// Custon error struct
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func HandleError(c *gin.Context, statusCode int, message string, err error) {
	if err != nil {
		c.Error(err) // Thêm lỗi vào gin context để dễ dàng quản lý sau này
	}

	// Trả về mã lỗi và thông báo lỗi dưới dạng JSON
	c.JSON(statusCode, APIError{
		Code:    statusCode,
		Message: message,
	})
}
