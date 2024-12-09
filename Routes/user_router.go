package routes

import (
	controllers "project/Controllers"
	middleware "project/Middleware"

	"github.com/gin-gonic/gin"
)

func SetupUserRouter(apiGroup *gin.RouterGroup) {
	user := apiGroup.Group("/user")
	{
		user.GET("/profile", middleware.RequireAuth, controllers.UserProfile)
		user.PUT("/profile/password", middleware.RequireAuth, controllers.ChangeUserPassword)
		user.POST("/profile", middleware.RequireAuth, middleware.FileUploadMiddleware(), controllers.ModifyAvatar)
		user.PUT("/profile", middleware.RequireAuth, controllers.UpdateUserProfile)
	}
}
