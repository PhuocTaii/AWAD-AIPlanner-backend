package routes

import (
	controllers "project/Controllers"
	middleware "project/Middleware"

	"github.com/gin-gonic/gin"
)

func SetupAuthRouter(apiGroup *gin.RouterGroup) {
	auth := apiGroup.Group("/auth")
	{
		auth.POST("/register", controllers.Register)
		auth.GET("/verify", controllers.Verify)
		auth.POST("/login", controllers.Login)
		auth.GET("/google_login", controllers.GoogleLogin)
		auth.GET("/google_callback", controllers.GoogleCallback)
		auth.POST("/logout", middleware.RequireAuth, controllers.Logout)
	}
}
