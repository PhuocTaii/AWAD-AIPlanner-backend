package controllers

import (
	"net/http"
	config "project/Config"
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
