package user

type UserResponse struct {
	Name     string `bson:"name" json:"name"`
	Email    string `bson:"email" json:"email"`
	GoogleId string `bson:"google_id" json:"google_id"`
	Avatar   string `bson:"avatar" json:"avatar"`
}
