package controllers

import (
	"net/http"
	config "project/Config"
	subject "project/Models/Request/Subject"
	services "project/Services"

	"github.com/gin-gonic/gin"
)

func CreateSubject(c *gin.Context) {
	// code here
	var request subject.CreateSubjectRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		error := &config.APIError{
			Code:    http.StatusBadRequest,
			Message: "Invalid task data",
		}
		config.HandleError(c, error)
		return
	}
	task, err := services.CreateSubject(c, request)
	if err != nil {
		defer config.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, task)
}

func GetSubjects(c *gin.Context) {
	tasks, err := services.GetSubjects(c)
	if err != nil {
		defer config.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, tasks)
}

func UpdateSubject(c *gin.Context) {
	var subjectId = c.Param("id")

	var request subject.ModifySubjectRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		error := &config.APIError{
			Code:    http.StatusBadRequest,
			Message: "Invalid task data",
		}
		config.HandleError(c, error)
		return
	}
	subject, err := services.ModifySubject(c, subjectId, request)
	if err != nil {
		defer config.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, subject)
}
