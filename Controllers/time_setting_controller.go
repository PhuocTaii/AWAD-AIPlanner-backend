package controllers

import (
	"net/http"
	config "project/Config"
	models "project/Models"
	timesetting "project/Models/Request/TimeSetting"
	services "project/Services"

	"github.com/gin-gonic/gin"
)

func GetTimeSetting(c *gin.Context) {
	timerSetting, err := services.GetTimeSetting(c)
	if err != nil {
		defer config.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, timerSetting)
}

func UpdateTimeSetting(c *gin.Context) {
	var request timesetting.UpdateTimeSettingRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		error := &config.APIError{
			Code:    http.StatusBadRequest,
			Message: "Invalid task data",
		}
		config.HandleError(c, error)
		return
	}

	tmp := &models.TimerSetting{
		FocusTime:  request.FocusTime,
		ShortBreak: request.ShortBreak,
		LongBreak:  request.LongBreak,
		Interval:   request.Interval,
	}

	timerSetting, err := services.UpdateTimeSetting(c, tmp)
	if err != nil {
		defer config.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, timerSetting)
}
