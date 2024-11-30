package routes

import (
	controllers "project/Controllers"
	middleware "project/Middleware"

	"github.com/gin-gonic/gin"
)

func SetupUserRouter(apiGroup *gin.RouterGroup) {
	user := apiGroup.Group("/user")
	{
		user.GET("/profile/:id", middleware.RequireAuth, controllers.UserProfile)
	}
}
