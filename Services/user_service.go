package services

import (
	"net/http"
	config "project/Config"
	models "project/Models"
	repository "project/Repository"

	"github.com/gin-gonic/gin"
)

func UserProfile(ctx *gin.Context, userId string) (*models.User, error) {
	var user models.User

	// Chuyển đổi userId (string) thành ObjectID
	// objectId, err := primitive.ObjectIDFromHex(userId)
	// if err != nil {
	// 	defer config.HandleError(ctx, http.StatusBadRequest, "Invalid user ID format", err)
	// 	return nil, err
	// }

	// Tìm người dùng theo ObjectID
	// err = config.UserCollection.FindOne(ctx, bson.M{"_id": objectId}).Decode(&user)

	err := repository.FindUserById(ctx, userId, &user)

	if err != nil {
		defer config.HandleError(ctx, http.StatusNotFound, "User not found", err)
		return nil, err
	}

	return &user, nil
}
