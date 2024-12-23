package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FocusLog struct {
	ID        primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	User      *primitive.ObjectID `bson:"user" json:"user"`
	FocusTime int                 `bson:"focus_time" json:"focus_time"`
	Date      string              `bson:"date" json:"date"`
	CreatedAt *time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt *time.Time          `bson:"updated_at" json:"updated_at"`
}
