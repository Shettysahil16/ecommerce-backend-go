package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Product struct {
	ID           bson.ObjectID `bson:"_id,omitempty" json:"id"`
	ProductName  string        `bson:"productName" json:"productName"`
	BrandName    string        `bson:"brandName" json:"brandName"`
	Category     string        `bson:"category" json:"category"`
	ProductImage []string      `bson:"productImage" json:"productImage"`
	Description  string        `bson:"description" json:"description"`
	Price        float64       `bson:"price" json:"price"`
	SellingPrice float64       `bson:"sellingPrice" json:"sellingPrice"`
	CreatedAt    time.Time     `bson:"createdAt" json:"createdAt"`
	UpdatedAt    time.Time     `bson:"updatedAt" json:"updatedAt"`
}
