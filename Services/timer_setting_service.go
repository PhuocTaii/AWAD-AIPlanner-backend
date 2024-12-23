package services

import (
	"net/http"
	config "project/Config"
	models "project/Models"
	repository "project/Repository"
	utils "project/Utils"

	"github.com/gin-gonic/gin"
)

func CreateTimerSetting(c *gin.Context, curUser *models.User) (*models.TimerSetting, *config.APIError) {
	// Insert timer setting
	timerSetting := &models.TimerSetting{
		User:       &curUser.ID,
		FocusTime:  25,
		ShortBreak: 5,
		LongBreak:  15,
		Interval:   4,
	}
	res, _ := repository.InsertTimerSetting(c, timerSetting)
	if res == nil {
		return nil, &config.APIError{
			Code:    400,
			Message: "Failed to create timer setting",
		}
	}
	return res, nil
}

func GetTimeSetting(ctx *gin.Context) (*models.TimerSetting, *config.APIError) {
	curUser, _ := utils.GetCurrentUser(ctx)
	if curUser == nil {
		return nil, &config.APIError{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
		}
	}

	timerSetting, _ := repository.GetTimerSettingByUserId(ctx, curUser)
	if timerSetting == nil {
		return nil, &config.APIError{
			Code:    http.StatusNotFound,
			Message: "Timer setting not found",
		}
	}

	return timerSetting, nil
}
