package configs

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func EnvMongoURI() string {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    return os.Getenv("MONGOURI")
}

func ConnectDB() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(EnvMongoURI()))
	if err != nil {
        log.Fatal(err)
    }
	err = client.Ping(ctx, nil)
    if err != nil {
        log.Fatal(err)
    }
	return client
  }

var DB *mongo.Client = ConnectDB()


func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
    collection := client.Database("chat").Collection(collectionName)
    return collection
}