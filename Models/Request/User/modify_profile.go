package user

type UpdateProfileRequest struct {
	Name string `json:"name" binding:"required"`
}
