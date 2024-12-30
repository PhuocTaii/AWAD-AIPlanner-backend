package controllers

import (
	"net/http"
	config "project/Config"
	focuslog "project/Models/Request/FocusLog"
	services "project/Services"

	"github.com/gin-gonic/gin"
)

func CreateFocusLog(c *gin.Context) {
	var request focuslog.CreateFocusLogRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	focusLog, err := services.CreateFocusLog(c, request)
	if err != nil {
		config.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": focusLog})
}
