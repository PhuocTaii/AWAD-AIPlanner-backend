package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name             string             `bson:"name" json:"name"`
	Email            string             `bson:"email" json:"email"`
	Password         string             `bson:"password" json:"password"`
	GoogleID         string             `bson:"google_id" json:"google_id"`
	Avatar           string             `bson:"avatar" json:"avatar"`
	IsVerified       bool               `bson:"is_verified" json:"is_verified"`
	VerificationCode string             `bson:"verification_code" json:"verification_code"`
	CreatedAt        *time.Time         `bson:"created_at" json:"created_at"`
	UpdatedAt        *time.Time         `bson:"updated_at" json:"updated_at"`
}
