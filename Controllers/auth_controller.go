package controllers

import (
	"net/http"
	config "project/Config"
	models "project/Models"
	auth "project/Models/Request/Auth"
	services "project/Services"

	"github.com/gin-gonic/gin"
)

// type GoogleUser struct {
// 	ID            string `json:"id"`
// 	Email         string `json:"email"`
// 	VerifiedEmail bool   `json:"verified_email"`
// 	Name          string `json:"name"`
// 	GivenName     string `json:"given_name"`
// 	FamilyName    string `json:"family_name"`
// 	Picture       string `json:"picture"`
// }

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

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  user,
	})
}
