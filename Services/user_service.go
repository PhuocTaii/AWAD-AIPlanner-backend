package services

import (
	"mime/multipart"
	"net/http"
	config "project/Config"
	models "project/Models"
	user "project/Models/Response/User"
	repository "project/Repository"
	utils "project/Utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func UserProfile(ctx *gin.Context) (*user.UserResponse, *config.APIError) {
	curUser, _ := utils.GetCurrentUser(ctx)
	if curUser == nil {
		return nil, &config.APIError{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
		}
	}

	foundUser, err := repository.FindUserById(ctx, curUser.ID.Hex())

	if err != nil {
		return nil, &config.APIError{
			Code:    http.StatusNotFound,
			Message: "User not found",
		}
	}

	res := &user.UserResponse{
		Name:     foundUser.Name,
		Email:    foundUser.Email,
		GoogleId: foundUser.GoogleID,
		Avatar:   foundUser.Avatar,
	}

	return res, nil
}

func ChangeUserPassword(ctx *gin.Context, oldPassword string, newPassword string) (*user.UserResponse, *config.APIError) {
	curUser, _ := utils.GetCurrentUser(ctx)
	if curUser == nil {
		return nil, &config.APIError{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
		}
	}

	foundUser, err := repository.FindUserById(ctx, curUser.ID.Hex())

	if err != nil {
		return nil, &config.APIError{
			Code:    http.StatusNotFound,
			Message: "User not found",
		}
	}

	if bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(oldPassword)) != nil {
		return nil, &config.APIError{
			Code:    http.StatusBadRequest,
			Message: "Old password is incorrect",
		}
	}

	if bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(newPassword)) == nil {
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

	foundUser.Password = string(hash)

	_, err = repository.UpdateUser(ctx, foundUser)

	if err != nil {
		return nil, &config.APIError{
			Code:    http.StatusBadRequest,
			Message: "Error updating user",
		}
	}

	res := &user.UserResponse{
		Name:     foundUser.Name,
		Email:    foundUser.Email,
		GoogleId: foundUser.GoogleID,
		Avatar:   foundUser.Avatar,
	}

	return res, nil

}

func ModifyAvatar(ctx *gin.Context, file multipart.File, filePath string) (*user.UserResponse, *config.APIError) {
	curUser, _ := utils.GetCurrentUser(ctx)
	if curUser == nil {
		return nil, &config.APIError{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
		}
	}

	foundUser, err := repository.FindUserById(ctx, curUser.ID.Hex())

	// user, err := repository.FindUserById(ctx, userId)
	if foundUser == nil {
		return nil, &config.APIError{
			Code:    http.StatusNotFound,
			Message: "User not found",
		}
	}
	if err != nil {
		return nil, &config.APIError{
			Code:    http.StatusNotFound,
			Message: "User not found",
		}
	}

	imageUrl, err := UploadToCloudinary(ctx, file, filePath)
	if err != nil {
		return nil, &config.APIError{
			Code:    http.StatusInternalServerError,
			Message: "Error uploading image",
		}
	}

	foundUser.Avatar = imageUrl
	_, err = repository.UpdateUser(ctx, foundUser)

	if err != nil {
		return nil, &config.APIError{
			Code:    http.StatusBadRequest,
			Message: "Error updating user",
		}
	}

	res := &user.UserResponse{
		Name:     foundUser.Name,
		Email:    foundUser.Email,
		GoogleId: foundUser.GoogleID,
		Avatar:   foundUser.Avatar,
	}

	return res, nil
}

func UpdateUserProfile(ctx *gin.Context, newName string) (*models.User, *config.APIError) {
	curUser, _ := utils.GetCurrentUser(ctx)
	if curUser == nil {
		return nil, &config.APIError{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
		}
	}
	user, err := repository.FindUserById(ctx, curUser.ID.Hex())
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
