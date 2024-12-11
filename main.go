package main

import (
	"os"
	config "project/Config"
	routes "project/Routes"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(config.CORSConfig())

	api := r.Group("/api")
	{
		routes.SetupAuthRouter(api)
		routes.SetupUserRouter(api)
		routes.SetupTaskRouter(api)
		routes.SetupSubjectRouter(api)
		routes.SetupTimeSettingRouter(api)
		routes.SetupFocusLogRouter(api)
	}

	return r
}

func main() {

	config.LoadEnvVariables()

	mongoURI := os.Getenv("MONGO_URI")
	dbName := os.Getenv("DB_NAME")

	config.ConnectDB(mongoURI, dbName)

	config.GoogleConfig()

	// r.Use(config.CORSConfig())\

	r := SetupRouter()

	r.Run(":8080")
}
