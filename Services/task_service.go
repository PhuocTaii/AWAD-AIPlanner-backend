package services

import (
	"net/http"
	config "project/Config"
	models "project/Models"
	constant "project/Models/Constant"
	task "project/Models/Request/Task"
	repository "project/Repository"
	utils "project/Utils"

	"github.com/gin-gonic/gin"
)

func CreateTask(c *gin.Context, request task.CreateTaskRequest) (*models.Task, *config.APIError) {
	//Get current user
	curUser, err := utils.GetCurrentUser(c)
	if err != nil {
		return nil, err
	}

	subject, err := FindSubjectById(c, request.SubjectId)
	if err != nil {
		return nil, err
	}

	priority, _ := constant.StringToPriority(request.Priority)
	if priority == -1 {
		return nil, &config.APIError{
			Code:    http.StatusInternalServerError,
			Message: "Invalid priority",
		}
	}

	if request.EstimatedStartTime > request.EstimatedEndTime {
		return nil, &config.APIError{
			Code:    http.StatusBadRequest,
			Message: "Invalid estimated start time and estimated end time",
		}
	}

	// Create task
	task := &models.Task{
		Name:               request.Name,
		Description:        request.Description,
		User:               *curUser,
		Subject:            *subject,
		Priority:           priority,
		Status:             constant.ToDo,
		EstimatedStartTime: request.EstimatedStartTime,
		EstimatedEndTime:   request.EstimatedEndTime,
	}

	// Insert task
	res, _ := repository.InsertTask(c, task)
	if res == nil {
		return nil, &config.APIError{
			Code:    http.StatusBadRequest,
			Message: "Failed to create task",
		}
	}
	return res, nil
}

// func ExpiredTask() *config.APIError{
// 	//
// }