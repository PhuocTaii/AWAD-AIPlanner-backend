package services

import (
	"net/http"
	"os"
	initializers "project/Initializers"
	models "project/Models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func Register(ctx *gin.Context, user *models.User) (*models.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		defer initializers.HandleError(ctx, http.StatusInternalServerError, "Error hashing password", err)
		return nil, err
	}

	if initializers.UserCollection.FindOne(ctx, bson.M{"email": user.Email}).Err() == nil {
		defer initializers.HandleError(ctx, http.StatusBadRequest, "Email already exists", nil)
		return nil, err
	}

	res, err := initializers.UserCollection.InsertOne(ctx, user)
	if err != nil {
		initializers.HandleError(ctx, http.StatusInternalServerError, "Error inserting user", err)
		return nil, err
	}

	newUser := &models.User{
		ID:       res.InsertedID.(primitive.ObjectID),
		Name:     user.Name,
		Email:    user.Email,
		Password: string(hash),
	}

	return newUser, nil
}

func Login(ctx *gin.Context, email, password string) (string, models.User, error) {
	var user models.User

	err := initializers.UserCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		defer initializers.HandleError(ctx, http.StatusBadRequest, "Invalid email or password", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		defer initializers.HandleError(ctx, http.StatusBadRequest, "Invalid email or password", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(), //token expires in 30 days
	})

	stringToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		defer initializers.HandleError(ctx, http.StatusInternalServerError, "Error generating token", err)
	}

	return stringToken, user, nil
}
