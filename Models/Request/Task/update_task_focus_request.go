package task

type UpdateTaskFocusRequest struct {
	FocusTime int `json:"focus_time" binding:"required"`
}
