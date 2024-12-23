package repository

import (
	config "project/Config"
	models "project/Models"
	utils "project/Utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetTodayFocusLog(ctx *gin.Context, user *models.User) (*models.FocusLog, error) {
	curDay := utils.GetCurrent().Format("2006-01-02")
	filter := bson.M{
		"user": user.ID,
		"date": curDay,
	}

	var focusLog *models.FocusLog
	err := config.FocusLogCollection.FindOne(ctx, filter).Decode(&focusLog)
	if err != nil {
		return nil, err
	}
	return focusLog, nil
}

func InsertFocusLog(ctx *gin.Context, focusLog *models.FocusLog) (*models.FocusLog, error) {
	newFocusLog := &models.FocusLog{
		User:      focusLog.User,
		FocusTime: focusLog.FocusTime,
		Date:      utils.GetCurrent().Format("2006-01-02"),
		CreatedAt: utils.GetCurrent(),
		UpdatedAt: utils.GetCurrent(),
	}

	res, err := config.FocusLogCollection.InsertOne(ctx, newFocusLog)
	if err != nil {
		return nil, err
	}

	response := &models.FocusLog{
		ID:        res.InsertedID.(primitive.ObjectID),
		User:      newFocusLog.User,
		FocusTime: newFocusLog.FocusTime,
		Date:      newFocusLog.Date,
		CreatedAt: newFocusLog.CreatedAt,
		UpdatedAt: newFocusLog.UpdatedAt,
	}

	return response, nil
}

func UpdateFocusLog(ctx *gin.Context, focusLog *models.FocusLog) (*models.FocusLog, error) {
	filter := bson.M{
		"_id": focusLog.ID,
	}
	update := bson.M{
		"$set": bson.M{
			"focus_time": focusLog.FocusTime,
			"updated_at": utils.GetCurrent(),
		},
	}

	_, err := config.FocusLogCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return focusLog, nil
}
