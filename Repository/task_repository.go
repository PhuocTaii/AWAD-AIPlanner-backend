package repository

import (
	config "project/Config"
	models "project/Models"
	utils "project/Utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func FindTaskById(ctx *gin.Context, id string) (*models.Task, error) {
	var task *models.Task
	err := config.TaskCollection.FindOne(ctx, bson.M{"_id": utils.ConvertObjectIDToString(id)}).Decode(&task)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func InsertTask(ctx *gin.Context, task *models.Task) (*models.Task, error) {
	newTask := &models.Task{
		Name:               task.Name,
		Description:        task.Description,
		Subject:            task.Subject,
		User:               task.User,
		Priority:           task.Priority,
		Status:             task.Status,
		EstimatedStartTime: task.EstimatedStartTime,
		EstimatedEndTime:   task.EstimatedEndTime,
		IsDeleted:          false,
		CreatedAt:          primitive.DateTime(utils.GetCurrentTime()),
		UpdatedAt:          primitive.DateTime(utils.GetCurrentTime()),
	}

	res, err := config.TaskCollection.InsertOne(ctx, newTask)
	if err != nil {
		return nil, err
	}

	response := &models.Task{
		ID:                 res.InsertedID.(primitive.ObjectID),
		Name:               newTask.Name,
		Description:        newTask.Description,
		Subject:            newTask.Subject,
		User:               newTask.User,
		Priority:           newTask.Priority,
		Status:             newTask.Status,
		EstimatedStartTime: newTask.EstimatedStartTime,
		EstimatedEndTime:   newTask.EstimatedEndTime,
		ActualStartTime:    newTask.ActualStartTime,
		ActualEndTime:      newTask.ActualEndTime,
		IsDeleted:          false,
		CreatedAt:          newTask.CreatedAt,
		UpdatedAt:          newTask.UpdatedAt,
	}

	return response, nil
}
