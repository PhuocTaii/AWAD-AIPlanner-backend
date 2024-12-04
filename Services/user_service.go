package services

import (
	"net/http"
	config "project/Config"
	models "project/Models"
	repository "project/Repository"

	"github.com/gin-gonic/gin"
)

func UserProfile(ctx *gin.Context, userId string) (*models.User, error) {
	var user *models.User

	user, err := repository.FindUserById(ctx, userId)

	if err != nil {
		defer config.HandleError(ctx, http.StatusNotFound, "User not found", err)
		return nil, err
	}

	return user, nil
}
