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
		user.PUT("/profile/password/:id", middleware.RequireAuth, controllers.ChangeUserPassword)
		user.POST("/profile/:id", middleware.RequireAuth, middleware.FileUploadMiddleware(), controllers.ModifyAvatar)
		user.PUT("/profile/:id", middleware.RequireAuth, controllers.UpdateUserProfile)
	}
}
