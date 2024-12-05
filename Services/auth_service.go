package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	config "project/Config"
	models "project/Models"
	auth "project/Models/Response/Auth"
	repository "project/Repository"
	utils "project/Utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func Register(ctx *gin.Context, user *models.User) (*models.User, *config.APIError) {

	fmt.Println(user)

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		err := &config.APIError{
			Code:    http.StatusInternalServerError,
			Message: "Error hashing password",
		}
		return nil, err
	}

	if config.UserCollection.FindOne(ctx, bson.M{"email": user.Email}).Err() == nil {
		err := &config.APIError{
			Code:    http.StatusBadRequest,
			Message: "Email already exists",
		}
		return nil, err
	}

	var tmp = &models.User{
		Name:     user.Name,
		Email:    user.Email,
		Password: string(hash),
	}

	newUser, err := repository.InsertUser(ctx, tmp)

	if err != nil {
		err := &config.APIError{
			Code:    http.StatusInternalServerError,
			Message: "Error inserting user",
		}
		return nil, err
	}

	return newUser, nil
}

func Login(ctx *gin.Context, email, password string) (string, *models.User, *config.APIError) {
	var user *models.User

	fmt.Println(email, password)
	user, err := repository.FindUserByEmail(ctx, email)
	if err != nil {
		err := &config.APIError{
			Code:    http.StatusBadRequest,
			Message: "Invalid email or password",
		}
		return "", nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		err := &config.APIError{
			Code:    http.StatusBadRequest,
			Message: "Invalid email or password",
		}
		return "", nil, err
	}

	stringToken, error := utils.GenerateJWT(ctx, user)
	if error != nil {
		err := &config.APIError{
			Code:    http.StatusInternalServerError,
			Message: "Error generating token",
		}
		return "", nil, err
	}

	return stringToken, user, nil
}

func GoogleLogin(c *gin.Context) (string, *models.User, *config.APIError) {
	state := c.Query("state")
	if state != "randomstate" {
		err := &config.APIError{
			Code:    http.StatusBadRequest,
			Message: "States don't Match!!",
		}
		return "", nil, err
	}

	code := c.Query("code")
	googlecon := config.GoogleConfig()
	token, err := googlecon.Exchange(context.Background(), code)
	if err != nil {
		err := &config.APIError{
			Code:    http.StatusInternalServerError,
			Message: "Code-Token Exchange Failed",
		}
		return "", nil, err
	}
	//Get reponse body from google
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		err := &config.APIError{
			Code:    http.StatusInternalServerError,
			Message: "User Data Fetch Failed",
		}
		return "", nil, err
	}
	userData, err := io.ReadAll(resp.Body)
	if err != nil {
		err := &config.APIError{
			Code:    http.StatusInternalServerError,
			Message: "JSON Parsing Failed",
		}
		return "", nil, err
	}
	//Unmarshal the data into a struct
	var user auth.GoogleUser

	err = json.Unmarshal(userData, &user)
	if err != nil {
		err := &config.APIError{
			Code:    http.StatusInternalServerError,
			Message: "JSON Decoding Failed",
		}
		return "", nil, err
	}
	var newUser *models.User
	var stringToken string

	newUser, err = repository.FindUserByEmail(c, user.Email) // Find user by email from google response

	if err != nil {
		tmp := &models.User{Name: user.Name, Email: user.Email, GoogleID: user.ID} // Create a new user
		newUser, err = repository.InsertUser(c, tmp)                               // Insert the user
		if err != nil {
			err := &config.APIError{
				Code:    http.StatusInternalServerError,
				Message: "Error inserting user",
			}
			return "", nil, err
		}
	} else {
		if newUser.GoogleID == "" {
			err := &config.APIError{
				Code:    http.StatusBadRequest,
				Message: "Email already exists",
			}
			return "", nil, err
		}
	}

	stringToken, _ = utils.GenerateJWT(c, newUser) // Generate JWT
	if stringToken == "" {
		err := &config.APIError{
			Code:    http.StatusInternalServerError,
			Message: "Error generating token",
		}
		return "", nil, err
	}

	return stringToken, newUser, nil
}
