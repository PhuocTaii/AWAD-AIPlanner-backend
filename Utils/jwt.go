package utils

import (
	"fmt"
	"net/http"
	"os"
	config "project/Config"
	models "project/Models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var tokenBlacklist = make(map[string]bool)

func GenerateJWT(ctx *gin.Context, user *models.User) (string, *config.APIError) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(), //token expires in 30 days
	})

	stringToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		err := &config.APIError{
			Code:    http.StatusInternalServerError,
			Message: "Error generating token",
		}
		defer config.HandleError(ctx, err)
		return "", err
	}

	return stringToken, nil
}

func GetToken(ctx *gin.Context) string {
	header := ctx.GetHeader("Authorization")

	if header == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
		ctx.AbortWithStatus(http.StatusUnauthorized)
	}

	tokenString := header[7:] // remove "Bearer " from token string
	return tokenString
}

func ExpireToken(ctx *gin.Context) {
	tokenString := GetToken(ctx)
	claims, _ := GetClaims(ctx)

	// Add the token to the blacklist
	tokenBlacklist[tokenString] = true

	// Optionally, you can also set the expiration claim
	claims["exp"] = time.Now().Add(-time.Hour).Unix()

	ctx.JSON(http.StatusOK, gin.H{"expired": claims["exp"]})
}

func isTokenBlacklisted(tokenString string) bool {
	return tokenBlacklist[tokenString]
}

func GetClaims(ctx *gin.Context) (jwt.MapClaims, error) {
	tokenString := GetToken(ctx)

	// Check if the token is blacklisted
	if isTokenBlacklisted(tokenString) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Token is expired"})
		defer ctx.AbortWithStatus(http.StatusUnauthorized)
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
