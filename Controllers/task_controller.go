package controllers

import (
	"net/http"
	config "project/Config"
	task "project/Models/Request/Task"
	services "project/Services"

	"github.com/gin-gonic/gin"
)

func CreateTask(c *gin.Context) {
	var request task.CreateTaskRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		error := &config.APIError{
			Code:    http.StatusBadRequest,
			Message: "Invalid task data",
		}
		config.HandleError(c, error)
		return
	}
	task, err := services.CreateTask(c, request)
	if err != nil {
		defer config.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, task)
}
