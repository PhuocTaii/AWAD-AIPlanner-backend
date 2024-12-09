package task

import "go.mongodb.org/mongo-driver/bson/primitive"

type ModifyTaskRequest struct {
	Name               string             `json:"name"`
	Description        string             `json:"description"`
	SubjectId          string             `json:"subject_id"`
	Priority           string             `json:"priority"`
	Status             string             `json:"status"`
	EstimatedStartTime primitive.DateTime `json:"estimated_start_time"`
	EstimatedEndTime   primitive.DateTime `json:"estimated_end_time"`
}
