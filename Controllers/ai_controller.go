package controllers

import (
	"net/http"
	config "project/Config"
	services "project/Services"

	"github.com/gin-gonic/gin"
)

func Feedback(c *gin.Context) {

	resp, err := services.Feedback(c)

	if err != nil {
		config.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}
