package controllers

import (
	"encoding/json"
	"net/http"
	"os"
	config "project/Config"
	models "project/Models"
	auth "project/Models/Request/Auth"
	services "project/Services"
	utils "project/Utils"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	var request auth.LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		error := &config.APIError{
			Code:    http.StatusBadRequest,
			Message: "Invalid user data",
		}
		config.HandleError(c, error)
		return
	}
	token, user, err := services.Login(c, request.Email, request.Password)

	if token == "" || user == nil {
		defer config.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful!",
		"token":   token,
		"user":    user,
	})
}

func GoogleLogin(c *gin.Context) {
	url := config.AppConfig.GoogleLoginConfig.AuthCodeURL("randomstate")

	c.Redirect(http.StatusSeeOther, url)
}

func GoogleCallback(c *gin.Context) {
	token, user, _, err := services.GoogleLogin(c)

	if err != nil {
		config.HandleError(c, err)
		return
	}

	// Convert user model to JSON string
	userJSON, error := json.Marshal(user)
	if error != nil {
		config.HandleError(c, err)
		return
	}

	url := os.Getenv("CLIENT_URL") + "?token=" + token + "&user=" + string(userJSON)
	c.Redirect(http.StatusSeeOther, url)
}

func Logout(c *gin.Context) {
	utils.ExpireToken(c)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully!"})
}

func Register(c *gin.Context) {
	var request auth.RegisterRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		error := &config.APIError{
			Code:    http.StatusBadRequest,
			Message: "Invalid user data",
		}
		config.HandleError(c, error)
		return
	}
	user := models.User{Name: request.Name, Email: request.Email, Password: request.Password}
	err := services.Register(c, &user)

	if err != nil {
		config.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Please verify your gmail!",
	})
}

func Verify(c *gin.Context) {
	code := c.Query("code")
	err := services.Verify(c, code)
	if err != nil {
		config.HandleError(c, err)
		return
	}
	c.Redirect(http.StatusSeeOther, os.Getenv("CLIENT_URL")+"/login")
}
