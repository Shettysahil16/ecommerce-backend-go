package database

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var DB *mongo.Database

func ConnectDB() {

	uri := os.Getenv("MONGO_URI")

	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Ping(ctx, nil)

	if err != nil {
		log.Fatal("MongoDB connection failed")
	}

	DB = client.Database(os.Getenv("DB_NAME"))

	log.Println("MongoDB Connected Successfully")
}
