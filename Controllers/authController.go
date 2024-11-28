package controllers

import (
	"net/http"
	models "project/Models"
	services "project/Services"
	utils "project/Utils"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.HandleError(c, http.StatusBadRequest, "Invalid user data", err)
		return
	}

	if err := services.Register(c, &user); err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to create user", err)
		return
	}
	c.JSON(http.StatusCreated, user)
}

func Login(c *gin.Context) {
	var request LoginRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		utils.HandleError(c, http.StatusBadRequest, "Invalid user data", err)
		return
	}

	message, err := services.Login(c, request.Email, request.Password)
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to login", err)
		return
	}

	// Trả về kết quả đăng nhập thành công
	c.JSON(http.StatusOK, gin.H{"message": message})
}
