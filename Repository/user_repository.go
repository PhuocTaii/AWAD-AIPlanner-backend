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

func FindUserByIdAndGoogleID(ctx *gin.Context, id string, googleID string) (*models.User, error) {
	var user *models.User
	err := config.UserCollection.FindOne(ctx, bson.M{"_id": utils.ConvertObjectIDToString(id), "google_id": googleID}).Decode(&user)
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
	newUser := &models.User{
		Name:      user.Name,
		Email:     user.Email,
		Password:  user.Password,
		GoogleID:  user.GoogleID,
		CreatedAt: primitive.DateTime(utils.GetCurrentTime()),
		UpdatedAt: primitive.DateTime(utils.GetCurrentTime()),
		Avatar:    "https://res.cloudinary.com/dl6v6a4nk/image/upload/v1733640750/tz0unlvxxyrnbgyp9x0y.jpg",
	}

	res, err := config.UserCollection.InsertOne(ctx, newUser)
	if err != nil {
		return nil, err
	}

	response := &models.User{
		ID:        res.InsertedID.(primitive.ObjectID),
		Name:      newUser.Name,
		Email:     newUser.Email,
		Password:  newUser.Password,
		GoogleID:  newUser.GoogleID,
		CreatedAt: newUser.CreatedAt,
		UpdatedAt: newUser.UpdatedAt,
		Avatar:    newUser.Avatar,
	}

	return response, nil
}

func UpdateUser(ctx *gin.Context, user *models.User) (*models.User, error) {
	filter := bson.M{"_id": user.ID}
	update := bson.M{"$set": bson.M{
		"name":       user.Name,
		"email":      user.Email,
		"password":   user.Password,
		"google_id":  user.GoogleID,
		"updated_at": primitive.DateTime(utils.GetCurrentTime()),
	}}

	_, err := config.UserCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return user, nil
}
