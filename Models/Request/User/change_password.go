package user

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	Password    string `json:"password" binding:"required, min=8"`
}
