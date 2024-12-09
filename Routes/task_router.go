package routes

import (
	controllers "project/Controllers"
	middleware "project/Middleware"

	"github.com/gin-gonic/gin"
)

func SetupTaskRouter(apiGroup *gin.RouterGroup) {
	task := apiGroup.Group("/task")
	{
		task.POST("/", middleware.RequireAuth, controllers.CreateTask)
		task.PUT("/:id", middleware.RequireAuth, controllers.UpdateTask)
		task.GET("/", middleware.RequireAuth, controllers.GetTasks)
		task.PUT("/delete/:id", middleware.RequireAuth, controllers.DeleteTask)
	}
}
