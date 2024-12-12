package routes

import (
	controllers "project/Controllers"
	middleware "project/Middleware"

	"github.com/gin-gonic/gin"
)

func SetupAiRouter(apiGroup *gin.RouterGroup) {
	ai := apiGroup.Group("/ai")
	{
		ai.GET("/create", middleware.RequireAuth, controllers.CreatePlan)
	}
}
