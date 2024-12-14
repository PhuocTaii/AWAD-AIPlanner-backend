package task

import (
	"time"
)

type ModifyTaskRequest struct {
	Name               string     `json:"name"`
	Description        string     `json:"description"`
	SubjectId          string     `json:"subject_id"`
	Priority           string     `json:"priority"`
	Status             string     `json:"status"`
	EstimatedStartTime *time.Time `json:"estimated_start_time"`
	EstimatedEndTime   *time.Time `json:"estimated_end_time"`
}
