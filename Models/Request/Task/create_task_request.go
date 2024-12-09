package task

import "go.mongodb.org/mongo-driver/bson/primitive"

type CreateTaskRequest struct {
	Name               string             `json:"name" binding:"required"`
	Description        string             `json:"description" binding:"required"`
	SubjectId          string             `json:"subject_id" binding:"required"`
	Priority           string             `json:"priority" binding:"required"`
	EstimatedStartTime primitive.DateTime `json:"estimated_start_time"`
	EstimatedEndTime   primitive.DateTime `json:"estimated_end_time"`
}
