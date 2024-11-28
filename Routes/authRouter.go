package routes

import (
	controllers "project/Controllers"

	"github.com/gin-gonic/gin"
)

func SetupAuthRouter(apiGroup *gin.RouterGroup) {
	auth := apiGroup.Group("/auth")
	{
		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.Login)
	}
}
