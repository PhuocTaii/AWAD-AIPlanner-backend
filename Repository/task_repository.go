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
	err := config.TaskCollection.FindOne(ctx, bson.M{"_id": utils.ConvertStringToObjectID(id)}).Decode(&task)
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
		CreatedAt:          utils.GetCurrent(),
		UpdatedAt:          utils.GetCurrent(),
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

func FindTaskByIdAndUserId(ctx *gin.Context, id string, userId string) (*models.Task, error) {
	var task *models.Task
	err := config.TaskCollection.FindOne(ctx, bson.M{"_id": utils.ConvertStringToObjectID(id), "user": utils.ConvertStringToObjectID(userId)}).Decode(&task)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func UpdateTask(ctx *gin.Context, task *models.Task) (*models.Task, error) {
	filter := bson.M{"_id": task.ID, "is_deleted": false}
	update := bson.M{"$set": bson.M{
		"name":                 task.Name,
		"description":          task.Description,
		"subject":              task.Subject,
		"status":               task.Status,
		"priority":             task.Priority,
		"estimated_start_time": task.EstimatedStartTime,
		"estimated_end_time":   task.EstimatedEndTime,
		"actual_start_time":    task.ActualStartTime,
		"actual_end_time":      task.ActualEndTime,
		"updated_at":           primitive.DateTime(utils.GetCurrentTime()),
	}}

	_, err := config.TaskCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func DeleteTask(ctx *gin.Context, task *models.Task) (*models.Task, error) {
	filter := bson.M{"_id": task.ID, "is_deleted": false}
	update := bson.M{"$set": bson.M{
		"is_deleted": true,
		"updated_at": primitive.DateTime(utils.GetCurrentTime()),
	}}
	_, err := config.TaskCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func GetTasks(ctx *gin.Context, filter bson.M) ([]*models.Task, error) {
	cursor, err := config.TaskCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tasks []*models.Task
	if err = cursor.All(ctx, &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}
