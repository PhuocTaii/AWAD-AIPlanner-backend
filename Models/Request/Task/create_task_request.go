package task

import (
	"time"
)

type CreateTaskRequest struct {
	Name               string     `json:"name" binding:"required"`
	Description        string     `json:"description" binding:"required"`
	SubjectId          string     `json:"subject_id" binding:"required"`
	Priority           string     `json:"priority" binding:"required"`
	EstimatedStartTime *time.Time `json:"estimated_start_time"`
	EstimatedEndTime   *time.Time `json:"estimated_end_time"`
}
