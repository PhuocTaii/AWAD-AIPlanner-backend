package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Task struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	//Reference to subject
	Subject Subject `bson:"subject" json:"subject"`
	User    User    `bson:"user" json:"user"`
	//Priority and Status of the task using enum
	Priority           int        `bson:"priority" json:"priority"`
	Status             int        `bson:"status" json:"status"`
	EstimatedStartTime *time.Time `bson:"estimated_start_time" json:"estimated_start_time"`
	EstimatedEndTime   *time.Time `bson:"estimated_end_time" json:"estimated_end_time"`
	ActualStartTime    *time.Time `bson:"actual_start_time" json:"actual_start_time"`
	ActualEndTime      *time.Time `bson:"actual_end_time" json:"actual_end_time"`
	FocusTime          int        `bson:"focus_time" json:"focus_time"`
	IsDeleted          bool       `bson:"is_deleted" json:"is_deleted"`
	CreatedAt          *time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt          *time.Time `bson:"updated_at" json:"updated_at"`
}
