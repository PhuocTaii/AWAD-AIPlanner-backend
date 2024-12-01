package repository

import (
	"fmt"
	config "project/Config"
	models "project/Models"
	utils "project/Utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func FindUserById(ctx *gin.Context, id string, user *models.User) error {
	fmt.Println(id)
	err := config.UserCollection.FindOne(ctx, bson.M{"_id": utils.ConvertObjectIDToString(id)}).Decode(&user)
	return err
}

func FindUserByEmail(ctx *gin.Context, email string, user *models.User) error {
	err := config.UserCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	return err
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
