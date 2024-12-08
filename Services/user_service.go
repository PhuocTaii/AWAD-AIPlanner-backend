package services

import (
	"mime/multipart"
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
		return nil, &config.APIError{
			Code:    http.StatusNotFound,
			Message: "User not found",
		}
	}

	return user, nil
}

func ChangeUserPassword(ctx *gin.Context, userId string, newPassword string) (*models.User, *config.APIError) {
	var user *models.User

	user, err := repository.FindUserByIdAndGoogleID(ctx, userId, "")

	if err != nil {
		return nil, &config.APIError{
			Code:    http.StatusNotFound,
			Message: "User not found",
		}
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(newPassword)) == nil {
		return nil, &config.APIError{
			Code:    http.StatusBadRequest,
			Message: "Password is the same as the old one",
		}
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, &config.APIError{
			Code:    http.StatusInternalServerError,
			Message: "Error hashing password",
		}
	}

	user.Password = string(hash)

	_, err = repository.UpdateUser(ctx, user)

	if err != nil {
		return nil, &config.APIError{
			Code:    http.StatusBadRequest,
			Message: "Error updating user",
		}
	}

	return user, nil
}

func ModifyAvatar(ctx *gin.Context, userId string, file multipart.File, filePath string) (string, *config.APIError) {
	var user *models.User

	user, err := repository.FindUserById(ctx, userId)
	if user == nil {
		return "", &config.APIError{
			Code:    http.StatusNotFound,
			Message: "User not found",
		}
	}
	if err != nil {
		return "", &config.APIError{
			Code:    http.StatusNotFound,
			Message: "User not found",
		}
	}

	imageUrl, err := UploadToCloudinary(ctx, file, filePath)
	if err != nil {
		return "", &config.APIError{
			Code:    http.StatusInternalServerError,
			Message: "Error uploading image",
		}
	}

	user.Avatar = imageUrl
	_, err = repository.UpdateUser(ctx, user)

	if err != nil {
		return "", &config.APIError{
			Code:    http.StatusBadRequest,
			Message: "Error updating user",
		}
	}

	return imageUrl, nil
}

func UpdateUserProfile(ctx *gin.Context, userId string, newName string) (*models.User, *config.APIError) {
	var user *models.User

	user, err := repository.FindUserById(ctx, userId)
	if err != nil {
		return nil, &config.APIError{
			Code:    http.StatusNotFound,
			Message: "User not found",
		}
	}
	user.Name = newName

	_, err = repository.UpdateUser(ctx, user)

	if err != nil {
		return nil, &config.APIError{
			Code:    http.StatusBadRequest,
			Message: "Error updating user",
		}
	}

	return user, nil
}
