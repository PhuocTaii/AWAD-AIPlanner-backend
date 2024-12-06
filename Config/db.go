package config

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client

var UserCollection *mongo.Collection
var TaskCollection *mongo.Collection
var SubjectCollection *mongo.Collection
var FocusLogCollection *mongo.Collection
var TimerSettingsCollection *mongo.Collection

func ConnectDB(dbUri string, dbName string) {
	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(dbUri).SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}
	// Send a ping to confirm a successful connection
	var result bson.M
	if err := client.Database(dbName).RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Decode(&result); err != nil {
		panic(err)
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	UserCollection = client.Database(os.Getenv("DB_NAME")).Collection("users")
	TaskCollection = client.Database(os.Getenv("DB_NAME")).Collection("tasks")
	SubjectCollection = client.Database(os.Getenv("DB_NAME")).Collection("subjects")
	FocusLogCollection = client.Database(os.Getenv("DB_NAME")).Collection("focus_logs")
	TimerSettingsCollection = client.Database(os.Getenv("DB_NAME")).Collection("timer_settings")
}
