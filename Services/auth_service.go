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

func Register(ctx *gin.Context, user *models.User) (*models.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		defer config.HandleError(ctx, http.StatusInternalServerError, "Error hashing password", err)
		return nil, err
	}

	if config.UserCollection.FindOne(ctx, bson.M{"email": user.Email}).Err() == nil {
		defer config.HandleError(ctx, http.StatusBadRequest, "Email already exists", nil)
		return nil, err
	}

	user = &models.User{
		Name:     user.Name,
		Email:    user.Email,
		Password: string(hash),
	}

	newUser, err := repository.InsertUser(ctx, user)

	if err != nil {
		defer config.HandleError(ctx, http.StatusInternalServerError, "Error inserting user", err)
		return nil, err
	}

	return newUser, nil
}

func Login(ctx *gin.Context, email, password string) (string, *models.User, error) {
	var user *models.User

	user, err := repository.FindUserByEmail(ctx, email)
	if err != nil {
		defer config.HandleError(ctx, http.StatusBadRequest, "Invalid email or password", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		defer config.HandleError(ctx, http.StatusBadRequest, "Invalid email or password", err)
	}

	stringToken, err := utils.GenerateJWT(ctx, user)
	if err != nil {
		defer config.HandleError(ctx, http.StatusInternalServerError, "Error generating token", err)
	}

	return stringToken, user, nil
}

func GoogleLogin(c *gin.Context) (string, *models.User, error) {
	state := c.Query("state")
	if state != "randomstate" {
		defer c.String(http.StatusBadRequest, "States don't Match!!")
	}

	code := c.Query("code")
	googlecon := config.GoogleConfig()
	token, err := googlecon.Exchange(context.Background(), code)
	if err != nil {
		defer c.String(http.StatusInternalServerError, "Code-Token Exchange Failed")
	}
	//Get reponse body from google
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		defer c.String(http.StatusInternalServerError, "User Data Fetch Failed")
	}
	userData, err := io.ReadAll(resp.Body)
	if err != nil {
		defer c.String(http.StatusInternalServerError, "JSON Parsing Failed")
	}
	//Unmarshal the data into a struct
	var user auth.GoogleUser

	err = json.Unmarshal(userData, &user)
	if err != nil {
		defer c.String(http.StatusInternalServerError, "JSON Decoding Failed")
	}
	var newUser *models.User
	var stringToken string

	newUser, err = repository.FindUserByEmail(c, user.Email) // Find user by email from google response

	if err != nil {
		tmp := &models.User{Name: user.Name, Email: user.Email}
		newUser, err = repository.InsertUser(c, tmp) // Insert the user
		if err != nil {
			defer config.HandleError(c, http.StatusInternalServerError, "Error inserting user", err)
			return "", nil, err
		}
	} else {
		newUser, err := repository.FindUserByEmail(c, user.Email)
		if newUser == nil {
			defer config.HandleError(c, http.StatusBadRequest, "Invalid email or password", err)
		}
	}

	stringToken, err = utils.GenerateJWT(c, newUser) // Generate JWT
	if err != nil {
		defer config.HandleError(c, http.StatusInternalServerError, "Error generating token", err)
		return "", nil, err
	}

	fmt.Println(newUser.ID)

	return stringToken, newUser, nil
}
