package repository

import (
	config "project/Config"
	models "project/Models"
	utils "project/Utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func InsertTimerSetting(ctx *gin.Context, timerSetting *models.TimerSetting) (*models.TimerSetting, error) {
	newTimerSetting := &models.TimerSetting{
		User:       timerSetting.User,
		FocusTime:  timerSetting.FocusTime,
		ShortBreak: timerSetting.ShortBreak,
		LongBreak:  timerSetting.LongBreak,
		Interval:   timerSetting.Interval,
		CreatedAt:  primitive.DateTime(utils.GetCurrentTime()),
		UpdatedAt:  primitive.DateTime(utils.GetCurrentTime()),
	}

	res, err := config.TimerSettingsCollection.InsertOne(ctx, newTimerSetting)
	if err != nil {
		return nil, err
	}

	response := &models.TimerSetting{
		ID:         res.InsertedID.(primitive.ObjectID),
		User:       newTimerSetting.User,
		FocusTime:  newTimerSetting.FocusTime,
		ShortBreak: newTimerSetting.ShortBreak,
		LongBreak:  newTimerSetting.LongBreak,
		Interval:   newTimerSetting.Interval,
		CreatedAt:  newTimerSetting.CreatedAt,
		UpdatedAt:  newTimerSetting.UpdatedAt,
	}

	return response, nil
}
