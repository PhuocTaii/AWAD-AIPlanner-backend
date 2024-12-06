package subject

type CreateSubjectRequest struct {
	Name string `json:"name" binding:"required"`
}
