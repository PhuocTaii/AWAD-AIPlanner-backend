package main

import (
	"log"
	"os"
	routes "project/Routes"
	"project/db"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")
	{
		// routes.SetupUserRouter(api)
		routes.SetupAuthRouter(api)
	}

	return r
}

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

	mongoURI := os.Getenv("MONGO_URI")
	dbName := os.Getenv("DB_NAME")
	port := os.Getenv("PORT")

	db.ConnectDB(mongoURI, dbName)

	r := SetupRouter()

	r.Run(":" + port)
}
