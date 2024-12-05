package controllers

import (
	"net/http"
	config "project/Config"
	user "project/Models/Request/User"
	services "project/Services"

	"github.com/gin-gonic/gin"
)

func UserProfile(c *gin.Context) {
	userId := c.Param("id")
	user, _ := services.UserProfile(c, userId)
	if user == nil {
		return
	}
	c.JSON(http.StatusOK, user)
}

func ChangeUserPassword(c *gin.Context) {
	userId := c.Param("id")
	var request user.ChangePasswordRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		error := &config.APIError{
			Code:    http.StatusBadRequest,
			Message: "Invalid user data",
		}
		config.HandleError(c, error)
		return
	}
	user, err := services.ChangeUserPassword(c, userId, request.Password)
	if err != nil {
		defer config.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, user)
}
