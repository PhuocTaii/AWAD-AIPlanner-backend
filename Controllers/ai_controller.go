package controllers

import (
	"net/http"
	config "project/Config"
	services "project/Services"

	"github.com/gin-gonic/gin"
)

func AiFeedback(c *gin.Context) {

	aiType := c.Query("type")

	if aiType == "" {
		config.HandleError(c, &config.APIError{
			Code:    http.StatusBadRequest,
			Message: "Invalid AI type",
		})
		return
	}

	resp, err := services.AIGen(c, aiType)

	if err != nil {
		config.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}
