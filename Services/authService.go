package services

import (
	"context"
	"errors"
	middleware "project/Middleware"
	models "project/Models"
	"project/db"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func Register(ctx *gin.Context, user *models.User) error {
	user.ID = primitive.NewObjectID()
	user.Password, _ = middleware.HashPassword(user.Password)
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	_, err := db.UserCollection.InsertOne(ctx, user)
	return err
}

func Login(ctx context.Context, email, password string) (string, error) {
	// Xác định timeout cho context
	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Khai báo biến result để lưu kết quả tìm kiếm user
	var result models.User

	// Tìm kiếm người dùng trong MongoDB
	err := db.UserCollection.FindOne(timeoutCtx, bson.M{"email": email}).Decode(&result)
	if err != nil {
		// Kiểm tra nếu không tìm thấy user
		if err == mongo.ErrNoDocuments {
			return "", errors.New("User not found")
		}
		return "", err
	}

	// Kiểm tra mật khẩu có hợp lệ không
	isPasswordValid := middleware.CheckPasswordHash(password, result.Password)
	if !isPasswordValid {
		return "", errors.New("invalid password")
	}

	// Trả về thông báo đăng nhập thành công
	return "Logged in", nil
}
