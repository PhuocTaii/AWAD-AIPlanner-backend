package routes

import (
	controllers "project/Controllers"
	middleware "project/Middleware"

	"github.com/gin-gonic/gin"
)

func SetupFocusLogRouter(apiGroup *gin.RouterGroup) {
	focusLog := apiGroup.Group("/focus_log")
	{
		focusLog.POST("/", middleware.RequireAuth, controllers.CreateFocusLog)
		// focusLog.POST("/", middleware.RequireAuth, controllers.GetFocusLog)
	}
}
