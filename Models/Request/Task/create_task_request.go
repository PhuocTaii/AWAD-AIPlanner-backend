package task

import (
	"time"
)

type CreateTaskRequest struct {
	Name               string     `json:"name" binding:"required"`
	Description        string     `json:"description"`
	SubjectId          string     `json:"subject_id"`
	Priority           string     `json:"priority" binding:"required"`
	Status             string     `json:"status" binding:"required"`
	EstimatedStartTime *time.Time `json:"estimated_start_time"`
	EstimatedEndTime   *time.Time `json:"estimated_end_time"`
}
