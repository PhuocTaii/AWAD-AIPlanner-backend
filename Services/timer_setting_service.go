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

func UpdateTimeSetting(c *gin.Context, timerSetting *models.TimerSetting) (*models.TimerSetting, *config.APIError) {
	if timerSetting.Interval < 1 || timerSetting.FocusTime < 1 || timerSetting.ShortBreak < 0 || timerSetting.LongBreak < 0 {
		return nil, &config.APIError{
			Code:    http.StatusBadRequest,
			Message: "Invalid time setting",
		}
	}
	curUser, _ := utils.GetCurrentUser(c)
	if curUser == nil {
		return nil, &config.APIError{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
		}
	}

	userTimeSetting, _ := repository.GetTimerSettingByUserId(c, curUser)

	if userTimeSetting == nil {
		return nil, &config.APIError{
			Code:    http.StatusNotFound,
			Message: "Timer setting not found",
		}
	}

	userTimeSetting.FocusTime = timerSetting.FocusTime
	userTimeSetting.ShortBreak = timerSetting.ShortBreak
	userTimeSetting.LongBreak = timerSetting.LongBreak
	userTimeSetting.Interval = timerSetting.Interval

	res, _ := repository.UpdateTimerSetting(c, userTimeSetting)
	if res == nil {
		return nil, &config.APIError{
			Code:    http.StatusNotFound,
			Message: "Timer setting not found",
		}
	}
	return res, nil
}
