package middleware

import (
	"fmt"
	"net/http"
	"os"
	config "project/Config"
	models "project/Models"
	repository "project/Repository"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func RequireAuth(c *gin.Context) {
	header := c.GetHeader("Authorization")

	if header == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	tokenString := header[7:] // remove "Bearer " from token string

	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			config.HandleError(c, http.StatusUnauthorized, "Token expired", nil)
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		var user models.User
		err := repository.FindUserById(c, claims["sub"].(string), &user)

		if err != nil {
			config.HandleError(c, http.StatusUnauthorized, "Unauthorized", err)
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		c.Set("user", user)

		c.Next()

	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}
