package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Task struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	//Reference to subject
	Subject Subject `bson:"subject" json:"subject"`
	User    User    `bson:"user" json:"user"`
	//Priority of the task using enum
	Priority           int                `bson:"priority" json:"priority"`
	Status             int                `bson:"status" json:"status"`
	EstimatedStartTime primitive.DateTime `bson:"estimated_start_time" json:"estimated_start_time"`
	EstimatedEndTime   primitive.DateTime `bson:"estimated_end_time" json:"estimated_end_time"`
	ActualStartTime    primitive.DateTime `bson:"actual_start_time" json:"actual_start_time"`
	ActualEndTime      primitive.DateTime `bson:"actual_end_time" json:"actual_end_time"`
	FocusTime          int                `bson:"focus_time" json:"focus_time"`
	IsDeleted          bool               `bson:"is_deleted" json:"is_deleted"`
	CreatedAt          primitive.DateTime `bson:"created_at" json:"created_at"`
	UpdatedAt          primitive.DateTime `bson:"updated_at" json:"updated_at"`
}
