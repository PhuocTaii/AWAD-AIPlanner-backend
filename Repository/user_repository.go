package repository

import (
	"fmt"
	config "project/Config"
	models "project/Models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func convertObjectIDToString(id string) primitive.ObjectID {
	objectId, _ := primitive.ObjectIDFromHex(id)
	fmt.Println(objectId.Hex())
	return objectId
}

func FindUserById(ctx *gin.Context, id string, user *models.User) error {
	fmt.Println(id)
	err := config.UserCollection.FindOne(ctx, bson.M{"_id": convertObjectIDToString(id)}).Decode(&user)
	return err
}
