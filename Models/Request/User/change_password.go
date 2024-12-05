package user

type ChangePasswordRequest struct {
	Password string `json:"password" binding:"required"`
}
