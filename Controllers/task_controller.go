package controllers

import (
	"net/http"
	config "project/Config"
	constant "project/Models/Constant"
	task "project/Models/Request/Task"
	response "project/Models/Response"
	services "project/Services"
	utils "project/Utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func UpdateTaskStatus(c *gin.Context) {
	var request task.ModifyTaskStatusRequest
	var taskId = c.Param("id")
	if err := c.ShouldBindJSON(&request); err != nil {
		error := &config.APIError{
			Code:    http.StatusBadRequest,
			Message: "Invalid task data",
		}
		config.HandleError(c, error)
		return
	}
	task, err := services.ModifyTaskStatus(c, taskId, request)
	if err != nil {
		defer config.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, task)
}

func GetTasks(c *gin.Context) {

	limit, _ := strconv.Atoi(c.Query("limit"))

	page, _ := strconv.Atoi(c.Query("page"))

	name := c.Query("name")
	subject := utils.ConvertStringToObjectID(c.Query("subject"))
	priority, _ := constant.StringToPriority(c.Query("priority"))
	status, _ := constant.StringToStatus(c.Query("status"))

	//if any of the query params are not provided, they will be ignored in the filter
	filter := bson.M{}
	if name != "" {
		filter["name"] = bson.M{"$regex": primitive.Regex{Pattern: name, Options: "i"}}
	}
	if !subject.IsZero() {
		filter["subject"] = subject
	}
	if priority != -1 {
		filter["priority"] = priority
	}
	if status != -1 {
		filter["status"] = status
	}
	filter["is_deleted"] = false

	sortBy := c.Query("sort_by")
	direction := c.Query("direction")

	sort := bson.M{}
	if sortBy != "" && direction != "" {
		if direction == "desc" {
			sort[sortBy] = -1
		}
		if direction == "asc" {
			sort[sortBy] = 1
		}
	}

	tasks, totalPages, totalItems, error := services.GetPagingTask(c, limit, page, filter, sort)

	if error != nil {
		config.HandleError(c, error)
		return
	}

	response := &response.PagingResponse{
		Data:        tasks,
		CurrentPage: page,
		PerPage:     limit,
		TotalPage:   totalPages,
		TotalItems:  totalItems,
	}

	c.JSON(http.StatusOK, response)
}

func DeleteTask(c *gin.Context) {
	var taskId = c.Param("id")
	res, err := services.DeleteTask(c, taskId)
	if err != nil {
		defer config.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, res)
}

func UpdateTaskFocus(c *gin.Context) {
	var request task.UpdateTaskFocusRequest
	var taskId = c.Param("id")
	if err := c.ShouldBindJSON(&request); err != nil {
		error := &config.APIError{
			Code:    http.StatusBadRequest,
			Message: "Invalid task data",
		}
		config.HandleError(c, error)
		return
	}
	task, err := services.ModifyTaskFocus(c, taskId, request)
	if err != nil {
		defer config.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, task)
}
