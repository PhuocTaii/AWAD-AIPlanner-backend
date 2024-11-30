package main

import (
	"os"
	initializers "project/Initializers"
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

	initializers.LoadEnvVariables()

	mongoURI := os.Getenv("MONGO_URI")
	dbName := os.Getenv("DB_NAME")
	port := os.Getenv("PORT")

	initializers.ConnectDB(mongoURI, dbName)

	r := SetupRouter()

	r.Run(":" + port)
}
