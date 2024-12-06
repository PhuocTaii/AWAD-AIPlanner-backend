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

func Register(ctx *gin.Context, user *models.User) (*models.User, *models.TimerSetting, *config.APIError) {

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, nil, &config.APIError{
			Code:    http.StatusInternalServerError,
			Message: "Error hashing password",
		}
	}

	if config.UserCollection.FindOne(ctx, bson.M{"email": user.Email}).Err() == nil {
		return nil, nil, &config.APIError{
			Code:    http.StatusBadRequest,
			Message: "Email already exists",
		}
	}

	var tmp = &models.User{
		Name:     user.Name,
		Email:    user.Email,
		Password: string(hash),
	}

	newUser, err := repository.InsertUser(ctx, tmp)

	if err != nil {
		return nil, nil, &config.APIError{
			Code:    http.StatusInternalServerError,
			Message: "Error inserting user",
		}
	}

	timerSetting, e := CreateTimerSetting(ctx, newUser)
	if e != nil {
		return nil, nil, e
	}

	return newUser, timerSetting, nil
}

func Login(ctx *gin.Context, email, password string) (string, *models.User, *config.APIError) {
	var user *models.User

	fmt.Println(email, password)
	user, err := repository.FindUserByEmail(ctx, email)
	if err != nil {
		return "", nil, &config.APIError{
			Code:    http.StatusBadRequest,
			Message: "Invalid email or password",
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", nil, &config.APIError{
			Code:    http.StatusBadRequest,
			Message: "Invalid email or password",
		}
	}

	stringToken, error := utils.GenerateJWT(ctx, user)
	if error != nil {
		return "", nil, &config.APIError{
			Code:    http.StatusInternalServerError,
			Message: "Error generating token",
		}
	}

	return stringToken, user, nil
}

func GoogleLogin(c *gin.Context) (string, *models.User, *models.TimerSetting, *config.APIError) {
	state := c.Query("state")
	if state != "randomstate" {
		return "", nil, nil, &config.APIError{
			Code:    http.StatusBadRequest,
			Message: "States don't Match!!",
		}
	}

	code := c.Query("code")
	googlecon := config.GoogleConfig()
	token, err := googlecon.Exchange(context.Background(), code)
	if err != nil {
		return "", nil, nil, &config.APIError{
			Code:    http.StatusInternalServerError,
			Message: "Code-Token Exchange Failed",
		}
	}
	//Get reponse body from google
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return "", nil, nil, &config.APIError{
			Code:    http.StatusInternalServerError,
			Message: "User Data Fetch Failed",
		}
	}
	userData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil, nil, &config.APIError{
			Code:    http.StatusInternalServerError,
			Message: "JSON Parsing Failed",
		}
	}
	//Unmarshal the data into a struct
	var user auth.GoogleUser

	err = json.Unmarshal(userData, &user)
	if err != nil {
		return "", nil, nil, &config.APIError{
			Code:    http.StatusInternalServerError,
			Message: "JSON Decoding Failed",
		}
	}
	var newUser *models.User
	var stringToken string
	var timerSetting *models.TimerSetting

	newUser, err = repository.FindUserByEmail(c, user.Email) // Find user by email from google response

	if err != nil {
		tmp := &models.User{Name: user.Name, Email: user.Email, GoogleID: user.ID} // Create a new user
		newUser, err = repository.InsertUser(c, tmp)                               // Insert the user
		if err != nil {
			return "", nil, nil, &config.APIError{
				Code:    http.StatusInternalServerError,
				Message: "Error inserting user",
			}
		}
		timerSetting, _ = CreateTimerSetting(c, newUser)
		if timerSetting == nil {
			return "", nil, nil, &config.APIError{
				Code:    http.StatusInternalServerError,
				Message: "Error inserting user",
			}
		}

	} else {
		if newUser.GoogleID == "" {
			return "", nil, nil, &config.APIError{
				Code:    http.StatusBadRequest,
				Message: "Email already exists",
			}
		}
	}

	stringToken, _ = utils.GenerateJWT(c, newUser) // Generate JWT
	if stringToken == "" {
		return "", nil, nil, &config.APIError{
			Code:    http.StatusInternalServerError,
			Message: "Error generating token",
		}
	}

	return stringToken, newUser, timerSetting, nil
}
