package controllers

import (
	"net/http"
	services "project/Services"

	"github.com/gin-gonic/gin"
)

func UserProfile(c *gin.Context) {
	userId := c.Param("id")
	user, _ := services.UserProfile(c, userId)
	if user == nil {
		return
	}
	c.JSON(http.StatusCreated, user)
}
