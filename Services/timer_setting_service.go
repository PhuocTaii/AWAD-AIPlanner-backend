package services

import (
	config "project/Config"
	models "project/Models"
	repository "project/Repository"

	"github.com/gin-gonic/gin"
)

func CreateTimerSetting(c *gin.Context, curUser *models.User) (*models.TimerSetting, *config.APIError) {
	// Insert timer setting
	timerSetting := &models.TimerSetting{
		User:       *curUser,
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
