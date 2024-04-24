package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var TasksCollection *mongo.Collection
var UsersCollection *mongo.Collection

func init() {
	var connectionString string = os.Getenv("MONGO_URI")
	var dbName string = os.Getenv("DB_NAME")
	var tasksCollectionName string = os.Getenv("TASKS_COLLECTION_NAME")
	var usersCollectionName string = os.Getenv("USERS_COLLECTION_NAME")

	clientOptions := options.Client().ApplyURI(connectionString)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	TasksCollection = client.Database(dbName).Collection(tasksCollectionName)
	UsersCollection = client.Database(dbName).Collection(usersCollectionName)
}
