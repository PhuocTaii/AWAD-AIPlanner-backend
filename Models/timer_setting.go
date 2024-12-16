package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TimerSetting struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	User       primitive.ObjectID `bson:"user" json:"user"`
	FocusTime  int                `bson:"focus_time" json:"focus_time"`
	ShortBreak int                `bson:"short_break_time" json:"short_break_time"`
	LongBreak  int                `bson:"long_break_time" json:"long_break_time"`
	Interval   int                `bson:"interval" json:"interval"`
	CreatedAt  *time.Time         `bson:"created_at" json:"created_at"`
	UpdatedAt  *time.Time         `bson:"updated_at" json:"updated_at"`
}
