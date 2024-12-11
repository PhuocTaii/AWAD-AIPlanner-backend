package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type FocusLog struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	User      User               `bson:"user_id" json:"user_id"`
	FocusTime int                `bson:"focus_time" json:"focus_time"`
	Date      primitive.DateTime `bson:"date" json:"date"`
	CreatedAt primitive.DateTime `bson:"created_at" json:"created_at"`
	UpdatedAt primitive.DateTime `bson:"updated_at" json:"updated_at"`
}
