package utils

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client         *mongo.Client
	collection     *mongo.Collection
	userCollection *mongo.Collection
)

func ConnectDB() {
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGODB_URI"))
	var err error
	client, err = mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB")

	collection = client.Database("golang_db").Collection("todos")
	userCollection = client.Database("golang_db").Collection("users")
}

func GetTodoCollection() *mongo.Collection {
	return collection
}

func GetUserCollection() *mongo.Collection {
	return userCollection
}

func DisconnectDB() {
	if err := client.Disconnect(context.Background()); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Disconnected from MongoDB")
}
