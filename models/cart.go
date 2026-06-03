package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Cart struct {
	ID        bson.ObjectID `bson:"_id" json:"id"`
	ProductID bson.ObjectID `bson:"productId" json:"productId"`
	UserID    string        `bson:"userId" json:"userId"`
	Quantity  int           `bson:"quantity" json:"quantity" binding:"required,min=1"`
	CreatedAt time.Time     `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time     `bson:"updatedAt" json:"updatedAt"`
}
