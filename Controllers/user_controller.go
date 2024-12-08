package controllers

import (
	"mime/multipart"
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

func ModifyAvatar(c *gin.Context) {
	userId := c.Param("id")
	file, _ := c.Get("file")
	filePath, _ := c.Get("filePath")
	imageUrl, err := services.ModifyAvatar(c, userId, file.(multipart.File), filePath.(string))
	if err != nil {
		defer config.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"imageUrl": imageUrl,
	})
}

func UpdateUserProfile(c *gin.Context) {
	userId := c.Param("id")
	var request user.UpdateProfileRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		error := &config.APIError{
			Code:    http.StatusBadRequest,
			Message: "Invalid user data",
		}
		config.HandleError(c, error)
		return
	}
	user, err := services.UpdateUserProfile(c, userId, request.Name)
	if err != nil {
		defer config.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, user)
}
