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
		task.GET("/:id", middleware.RequireAuth, controllers.GetTaskById)
		task.DELETE("/:id", middleware.RequireAuth, controllers.DeleteTask)
		task.PUT("/status/:id", middleware.RequireAuth, controllers.UpdateTaskStatus)
		task.PUT("/focus/:id", middleware.RequireAuth, controllers.UpdateTaskFocus)
	}
}
