package routes

import (
	controllers "project/Controllers"
	middleware "project/Middleware"

	"github.com/gin-gonic/gin"
)

func SetupTimeSettingRouter(apiGroup *gin.RouterGroup) {
	timeSetting := apiGroup.Group("/time-setting")
	{
		timeSetting.GET("/", middleware.RequireAuth, controllers.GetTimeSetting)
	}
}
