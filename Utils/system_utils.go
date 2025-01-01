package utils

import (
	"encoding/hex"
	"net/http"
	config "project/Config"
	models "project/Models"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/rand"
)

func GetCurrentTime() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func GetCurrent() *time.Time {
	now := time.Now()
	return &now
}

func GenerateVerifcationCode() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func GeneratePassword() string {
	// random a 10 character password
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 10
	seededRand := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(result)
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
