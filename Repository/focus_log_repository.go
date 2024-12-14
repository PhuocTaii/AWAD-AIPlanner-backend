package repository

import (
	config "project/Config"
	models "project/Models"
	utils "project/Utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func InsertFocusLog(ctx *gin.Context, focusLog *models.FocusLog) (*models.FocusLog, error) {
	newFocusLog := &models.FocusLog{
		User:      focusLog.User,
		FocusTime: focusLog.FocusTime,
		Date:      utils.GetCurrent(),
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
