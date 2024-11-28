package controllers

import (
	"net/http"
	models "project/Models"
	services "project/Services"
	utils "project/Utils"

	"github.com/gin-gonic/gin"
)

func CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.HandleError(c, http.StatusBadRequest, "Invalid user data", err)
		return
	}

	if err := services.CreateUser(&user); err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to create user", err)
		return
	}
	c.JSON(http.StatusCreated, user)
}
