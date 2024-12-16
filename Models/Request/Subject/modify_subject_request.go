package subject

type ModifySubjectRequest struct {
	Name string `json:"name" binding:"required"`
}
