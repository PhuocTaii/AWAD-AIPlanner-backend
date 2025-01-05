package controllers

import (
	"net/http"
	config "project/Config"
	models "project/Models"
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
	subject := &models.Subject{
		Name: request.Name,
	}
	task, err := services.CreateSubject(c, subject)
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

func DeleteSubject(c *gin.Context) {
	var subjectId = c.Param("id")
	res, err := services.DeleteSubject(c, subjectId)
	if err != nil {
		defer config.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, res)
}

func GetTaskAmountsBySubject(c *gin.Context) {
	taskAmounts, err := services.GetTaskAmountBySubject(c)
	if err != nil {
		defer config.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, taskAmounts)
}
