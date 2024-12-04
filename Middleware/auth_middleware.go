package middleware

import (
	"net/http"
	config "project/Config"
	models "project/Models"
	repository "project/Repository"
	utils "project/Utils"

	"github.com/gin-gonic/gin"
)

func RequireAuth(c *gin.Context) {
	claims, error := utils.GetClaims(c)

	if error != nil {
		defer c.AbortWithStatus(http.StatusUnauthorized)
	}

	var user *models.User
	user, err := repository.FindUserById(c, claims["sub"].(string))

	if err != nil {
		config.HandleError(c, http.StatusUnauthorized, "Unauthorized", err)
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	c.Set("user", user)

	c.Next()
}
