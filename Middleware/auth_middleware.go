package middleware

import (
	"fmt"
	"net/http"
	"os"
	initializers "project/Initializers"
	models "project/Models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func RequireAuth(c *gin.Context) {
	header := c.GetHeader("Authorization")

	if header == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	fmt.Println(header)

	tokenString := header[7:] // remove "Bearer " from token string

	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			initializers.HandleError(c, http.StatusUnauthorized, "Token expired", nil)
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		objectId, _ := primitive.ObjectIDFromHex(claims["sub"].(string))

		var user models.User
		err := initializers.UserCollection.FindOne(c, bson.M{"_id": objectId}).Decode(&user)

		if err != nil {
			initializers.HandleError(c, http.StatusUnauthorized, "Unauthorized", err)
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		c.Set("user", user)

		c.Next()

	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}
