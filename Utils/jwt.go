package utils

import (
	"net/http"
	"os"
	config "project/Config"
	models "project/Models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(ctx *gin.Context, user *models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(), //token expires in 30 days
	})

	stringToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		defer config.HandleError(ctx, http.StatusInternalServerError, "Error generating token", err)
	}

	return stringToken, err
}
