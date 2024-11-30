package controllers

import (
	"net/http"
	initializers "project/Initializers"
	models "project/Models"
	auth "project/Models/Request/Auth"
	services "project/Services"

	"github.com/gin-gonic/gin"
)

// type LoginRequest struct {
// 	Email    string `json:"email"`
// 	Password string `json:"password"`
// }

// type RegisterRequest struct {
// 	Name     string `json:"Name"`
// 	Email    string `json:"email"`
// 	Password string `json:"password"`
// }

func Register(c *gin.Context) {
	var request auth.RegisterRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		initializers.HandleError(c, http.StatusBadRequest, "Invalid user data", err)
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
		initializers.HandleError(c, http.StatusBadRequest, "Invalid user data", err)
		return
	}

	token, user, err := services.Login(c, request.Email, request.Password)

	if err != nil {
		initializers.HandleError(c, http.StatusInternalServerError, "Failed to login", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  user,
	})
}
