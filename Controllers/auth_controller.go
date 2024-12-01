package controllers

import (
	"context"
	"io/ioutil"
	"net/http"
	config "project/Config"
	models "project/Models"
	auth "project/Models/Request/Auth"
	services "project/Services"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var request auth.RegisterRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		config.HandleError(c, http.StatusBadRequest, "Invalid user data", err)
		return
	}

	user := models.User{Name: request.Name, Email: request.Email, Password: request.Password}

	res, _ := services.Register(c, &user)

	if res == nil {
		return
	}

	c.JSON(http.StatusCreated, res)
}

func Login(c *gin.Context) {
	var request auth.LoginRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		config.HandleError(c, http.StatusBadRequest, "Invalid user data", err)
		return
	}

	token, user, err := services.Login(c, request.Email, request.Password)

	if err != nil {
		config.HandleError(c, http.StatusInternalServerError, "Failed to login", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  user,
	})
}

func GoogleLogin(c *gin.Context) {
	url := config.AppConfig.GoogleLoginConfig.AuthCodeURL("randomstate")

	c.Redirect(http.StatusSeeOther, url)
}

func GoogleCallback(c *gin.Context) {
	state := c.Query("state")
	if state != "randomstate" {
		c.String(http.StatusBadRequest, "States don't Match!!")
		return
	}

	code := c.Query("code")

	googlecon := config.GoogleConfig()

	token, err := googlecon.Exchange(context.Background(), code)
	if err != nil {
		c.String(http.StatusInternalServerError, "Code-Token Exchange Failed")
		return
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		c.String(http.StatusInternalServerError, "User Data Fetch Failed")
		return
	}

	userData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.String(http.StatusInternalServerError, "JSON Parsing Failed")
		return
	}

	c.String(http.StatusOK, string(userData))
}
