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

func UpdateTask(c *gin.Context) {
	var request task.ModifyTaskRequest
	var taskId = c.Param("id")
	if err := c.ShouldBindJSON(&request); err != nil {
		error := &config.APIError{
			Code:    http.StatusBadRequest,
			Message: "Invalid task data",
		}
		config.HandleError(c, error)
		return
	}
	task, err := services.ModifyTask(c, taskId, request)
	if err != nil {
		defer config.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, task)
}
