package services

import (
	"net/http"
	config "project/Config"
	models "project/Models"
	repository "project/Repository"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func UserProfile(ctx *gin.Context, userId string) (*models.User, *config.APIError) {
	var user *models.User

	user, err := repository.FindUserById(ctx, userId)

	if err != nil {
		err := &config.APIError{
			Code:    http.StatusNotFound,
			Message: "User not found",
		}
		return nil, err
	}

	return user, nil
}

func ChangeUserPassword(ctx *gin.Context, userId string, newPassword string) (*models.User, *config.APIError) {
	var user *models.User

	user, err := repository.FindUserByIdAndGoogleID(ctx, userId, "")

	if err != nil {
		err := &config.APIError{
			Code:    http.StatusNotFound,
			Message: "User not found",
		}
		return nil, err
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(newPassword)) == nil {
		err := &config.APIError{
			Code:    http.StatusBadRequest,
			Message: "Password is the same as the old one",
		}
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		err := &config.APIError{
			Code:    http.StatusInternalServerError,
			Message: "Error hashing password",
		}
		return nil, err
	}

	user.Password = string(hash)

	_, err = repository.UpdateUser(ctx, user)

	if err != nil {
		err := &config.APIError{
			Code:    http.StatusBadRequest,
			Message: "Error updating user",
		}
		return nil, err
	}

	return user, nil
}
