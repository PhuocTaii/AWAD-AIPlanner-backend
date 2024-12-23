package repository

import (
	config "project/Config"
	models "project/Models"
	utils "project/Utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func InsertTimerSetting(ctx *gin.Context, timerSetting *models.TimerSetting) (*models.TimerSetting, error) {
	newTimerSetting := &models.TimerSetting{
		User:       timerSetting.User,
		FocusTime:  timerSetting.FocusTime,
		ShortBreak: timerSetting.ShortBreak,
		LongBreak:  timerSetting.LongBreak,
		Interval:   timerSetting.Interval,
		CreatedAt:  utils.GetCurrent(),
		UpdatedAt:  utils.GetCurrent(),
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

func GetTimerSettingByUserId(ctx *gin.Context, user *models.User) (*models.TimerSetting, error) {
	var timerSetting models.TimerSetting
	filter := bson.M{
		"user": user.ID,
	}
	err := config.TimerSettingsCollection.FindOne(ctx, filter).Decode(&timerSetting)
	if err != nil {
		return nil, err
	}

	return &timerSetting, nil
}

func UpdateTimerSetting(ctx *gin.Context, timerSetting *models.TimerSetting) (*models.TimerSetting, error) {
	filter := bson.M{
		"user": timerSetting.User,
	}
	update := bson.M{
		"$set": bson.M{
			"focus_time":       timerSetting.FocusTime,
			"short_break_time": timerSetting.ShortBreak,
			"long_break_time":  timerSetting.LongBreak,
			"interval":         timerSetting.Interval,
			"updated_at":       utils.GetCurrent(),
		},
	}
	_, err := config.TimerSettingsCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return timerSetting, nil
}
