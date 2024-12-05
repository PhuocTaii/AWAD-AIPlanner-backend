package config

import (
	"github.com/gin-gonic/gin"
)

// Custon error struct
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func HandleError(c *gin.Context, err *APIError) {

	// Trả về mã lỗi và thông báo lỗi dưới dạng JSON
	c.JSON(err.Code, APIError{
		Code:    err.Code,
		Message: err.Message,
	})
}
