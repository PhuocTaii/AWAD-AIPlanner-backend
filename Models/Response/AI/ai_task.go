package ai

import "time"

type AiTask struct {
	Name        string `bson:"name" json:"name"`
	Description string `bson:"description" json:"description"`
	//Reference to subject
	Subject string `bson:"subject" json:"subject"`
	User    string `bson:"user" json:"user"`
	//Priority and Status of the task using enum
	Priority           string     `bson:"priority" json:"priority"`
	Status             string     `bson:"status" json:"status"`
	FocusTime          int        `bson:"focus_time" json:"focus_time"`
	EstimatedStartTime *time.Time `bson:"estimated_start_time" json:"estimated_start_time"`
	EstimatedEndTime   *time.Time `bson:"estimated_end_time" json:"estimated_end_time"`
}
