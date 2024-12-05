package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type FocusLog struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserId    User               `bson:"user_id" json:"user_id"`
	Date      primitive.DateTime `bson:"date" json:"date"`
	FocusTime int                `bson:"focus_time" json:"focus_time"`
}
