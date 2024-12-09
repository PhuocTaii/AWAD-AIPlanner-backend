package utils

import (
	"net/http"
	config "project/Config"
	models "project/Models"
	"time"

	"github.com/gin-gonic/gin"
)

func GetCurrentTime() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func GetCurrent() *time.Time {
	now := time.Now()
	return &now
}

func GetCurrentUser(c *gin.Context) (*models.User, *config.APIError) {
	// var user *models.User

	// Get current user from context
	userInterface, _ := c.Get("user")
	if userInterface == nil {
		err := &config.APIError{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
		}
		return nil, err
	}
	user := userInterface.(*models.User)
	return user, nil
}
