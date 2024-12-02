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
	token, user, err := services.GoogleLogin(c)

	if err != nil {
		config.HandleError(c, http.StatusInternalServerError, "Failed to login", err)
		return
	}

	// Convert user model to JSON string
	userJSON, err := json.Marshal(user)
	if err != nil {
		config.HandleError(c, http.StatusInternalServerError, "Failed to parse user to JSON", err)
		return
	}

	url := os.Getenv("CLIENT_URL") + "?token=" + token + "&user=" + string(userJSON)
	c.Redirect(http.StatusSeeOther, url)
}

func Logout(c *gin.Context) {
	utils.ExpireToken(c)
	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}
