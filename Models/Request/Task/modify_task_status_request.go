package task

type ModifyTaskStatusRequest struct {
	Status string `json:"status" binding:"required"`
}
