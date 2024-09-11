package database

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectWithMongodb() *mongo.Client {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		log.Fatal("mongo url is missing")
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Database connected successfully")
	}
	return client
}

var Client *mongo.Client = ConnectWithMongodb()

func OpenConnection(client *mongo.Client, collectionName string) *mongo.Collection {

	var collection *mongo.Collection = client.Database("tentree").Collection(collectionName)
	return collection
}
