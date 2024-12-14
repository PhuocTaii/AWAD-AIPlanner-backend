package routes

import (
	controllers "project/Controllers"
	middleware "project/Middleware"

	"github.com/gin-gonic/gin"
)

func SetupSubjectRouter(apiGroup *gin.RouterGroup) {
	subject := apiGroup.Group("/subject")
	{
		subject.POST("/", middleware.RequireAuth, controllers.CreateSubject)
		subject.GET("/", middleware.RequireAuth, controllers.GetSubjects)
	}
}
