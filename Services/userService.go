package services

import (
	"context"
	models "project/Models"
	"project/db"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateUser(user *models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	_, err := db.UserCollection.InsertOne(ctx, user)
	return err
}
