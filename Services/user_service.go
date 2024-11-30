package services

import (
	"net/http"
	initializers "project/Initializers"
	models "project/Models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func UserProfile(ctx *gin.Context, userId string) (*models.User, error) {
	var user models.User

	// Chuyển đổi userId (string) thành ObjectID
	objectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		defer initializers.HandleError(ctx, http.StatusBadRequest, "Invalid user ID format", err)
		return nil, err
	}

	// Tìm người dùng theo ObjectID
	err = initializers.UserCollection.FindOne(ctx, bson.M{"_id": objectId}).Decode(&user)
	if err != nil {
		defer initializers.HandleError(ctx, http.StatusNotFound, "User not found", err)
		return nil, err
	}

	return &user, nil
}
