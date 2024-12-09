package auth

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"` // email validation
	Password string `json:"password" binding:"required,min=8"` // min length 8
}
