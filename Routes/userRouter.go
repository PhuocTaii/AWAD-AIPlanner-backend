package routes

import (
	controllers "project/Controllers"

	"github.com/gin-gonic/gin"
)

func SetupUserRouter(apiGroup *gin.RouterGroup) {
	user := apiGroup.Group("/users")
	{
		user.POST("/", controllers.CreateUser)
	}
}
