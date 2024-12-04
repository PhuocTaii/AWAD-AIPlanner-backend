package repository

import (
	config "project/Config"
	models "project/Models"
	utils "project/Utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func FindUserById(ctx *gin.Context, id string) (*models.User, error) {
	var user *models.User
	err := config.UserCollection.FindOne(ctx, bson.M{"_id": utils.ConvertObjectIDToString(id)}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func FindUserByEmail(ctx *gin.Context, email string) (*models.User, error) {
	var user *models.User
	err := config.UserCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func InsertUser(ctx *gin.Context, user *models.User) (*models.User, error) {
	res, err := config.UserCollection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	newUser := &models.User{
		ID:       res.InsertedID.(primitive.ObjectID),
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}

	return newUser, nil
}
