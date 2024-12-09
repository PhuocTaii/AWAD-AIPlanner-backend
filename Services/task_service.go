package services

import (
	"fmt"
	"net/http"
	config "project/Config"
	models "project/Models"
	constant "project/Models/Constant"
	task "project/Models/Request/Task"
	repository "project/Repository"
	utils "project/Utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
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

func ModifyTask(c *gin.Context, id string, request task.ModifyTaskRequest) (*models.Task, *config.APIError) {
	//Get current user
	curUser, err := utils.GetCurrentUser(c)
	if err != nil {
		return nil, err
	}

	task, _ := repository.FindTaskByIdAndUserId(c, id, curUser.ID.Hex())
	if task == nil {
		return nil, &config.APIError{
			Code:    http.StatusNotFound,
			Message: "Task not found",
		}
	}

	//cannot modify task to expired
	if request.Status == "Expired" {
		return nil, &config.APIError{
			Code:    http.StatusBadRequest,
			Message: "Cannot modify task to expired",
		}
	}

	// Modify task
	if request.Name == "" {
		request.Name = task.Name
	}
	if request.Description == "" {
		request.Description = task.Description
	}
	if request.SubjectId == "" {
		request.SubjectId = task.Subject.ID.Hex()
	}
	if request.Priority == "" {
		request.Priority = constant.PriorityToString(task.Priority)
	}
	if request.Status == "" {
		request.Status = constant.StatusToString(task.Status)
	}
	if request.EstimatedStartTime == 0 {
		request.EstimatedStartTime = task.EstimatedStartTime
	}
	if request.EstimatedEndTime == 0 {
		request.EstimatedEndTime = task.EstimatedEndTime
	}

	if request.EstimatedStartTime > request.EstimatedEndTime {
		return nil, &config.APIError{
			Code:    http.StatusBadRequest,
			Message: "Invalid estimated start time and estimated end time",
		}
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

	status, _ := constant.StringToStatus(request.Status)
	if status == -1 {
		return nil, &config.APIError{
			Code:    http.StatusInternalServerError,
			Message: "Invalid status",
		}
	}

	task.Name = request.Name
	task.Description = request.Description
	task.Priority = priority
	task.Subject = *subject
	task.EstimatedStartTime = request.EstimatedStartTime
	task.EstimatedEndTime = request.EstimatedEndTime

	//cannot modify status of expired task
	if task.Status == constant.Expired {
		fmt.Println("hello, this is thanh vân")
		//check valid new estimated start and end time

		if task.ActualStartTime == nil || task.ActualStartTime.After(*utils.GetCurrent()) {
			task.Status = constant.ToDo
		} else {
			task.Status = constant.InProgress
		}
	} else {
		//if change status to completed, set actual end time
		if status == constant.Completed {
			task.ActualEndTime = utils.GetCurrent()
			//if change status from to do to completed, set actual start time to current time
			if task.Status == constant.ToDo {
				task.ActualStartTime = utils.GetCurrent()
			}
		}
		//if change status to to do, set actual start time and end time to nil
		if status == constant.ToDo {
			task.ActualStartTime = nil
			task.ActualEndTime = nil
		}
		//if change status to in progress, set actual start time to current time
		if status == constant.InProgress {
			task.ActualStartTime = utils.GetCurrent()
			if task.Status == constant.Completed {
				task.ActualEndTime = nil
			}
		}

		task.Status = status
	}

	// Update task
	res, _ := repository.UpdateTask(c, task)
	if res == nil {
		return nil, &config.APIError{
			Code:    http.StatusBadRequest,
			Message: "Failed to update task",
		}
	}
	return res, nil
}

func GetPagingTask(c *gin.Context, limit, page int, filter, sort bson.M) ([]*models.Task, int, int, *config.APIError) {
	paginateConfig := config.NewPagingConfig(c, limit, page)

	tasks, totalPages, totalItems, err := config.PaginatedFind[*models.Task](c, config.TaskCollection, paginateConfig, filter, sort)
	if err != nil {
		return nil, 0, 0, &config.APIError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get tasks",
		}
	}

	return tasks, totalPages, totalItems, nil

}

func DeleteTask(c *gin.Context, id string) (*models.Task, *config.APIError) {
	//Get current user
	curUser, err := utils.GetCurrentUser(c)
	if err != nil {
		return nil, err
	}

	task, _ := repository.FindTaskByIdAndUserId(c, id, curUser.ID.Hex())
	if task == nil {
		return nil, &config.APIError{
			Code:    http.StatusNotFound,
			Message: "Task not found",
		}
	}

	// Update task is_deleted to true
	task.IsDeleted = true

	res, _ := repository.DeleteTask(c, task)
	if res == nil {
		return nil, &config.APIError{
			Code:    http.StatusBadRequest,
			Message: "Failed to delete task",
		}
	}
	return res, nil
}
