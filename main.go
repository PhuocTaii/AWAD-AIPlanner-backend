package main

import (
	"os"
	config "project/Config"
	routes "project/Routes"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")
	{
		routes.SetupAuthRouter(api)
		routes.SetupUserRouter(api)
	}

	return r
}

func main() {

	config.LoadEnvVariables()

	mongoURI := os.Getenv("MONGO_URI")
	dbName := os.Getenv("DB_NAME")

	config.ConnectDB(mongoURI, dbName)

	config.GoogleConfig()

	r := SetupRouter()

	r.Run(":8080")
}
