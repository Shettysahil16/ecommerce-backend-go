package database

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func CreateIndexes() {

	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				bson.E{Key: "userId", Value: 1},
			},
		},
	}

	_, err := CartCollection().Indexes().CreateMany(context.Background(), indexes)

	if err != nil {
		log.Fatal(err)
	}
}
