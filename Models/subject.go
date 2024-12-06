package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Subject struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name" json:"name"`
	IsDeleted bool               `bson:"is_deleted" json:"is_deleted"`
	User      User               `bson:"user" json:"user"`
	CreatedAt primitive.DateTime `bson:"created_at" json:"created_at"`
	UpdatedAt primitive.DateTime `bson:"updated_at" json:"updated_at"`
}
