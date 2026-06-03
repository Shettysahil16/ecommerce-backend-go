package models

type CategoryPreview struct {
	Category string  `bson:"_id" json:"category"`
	Product  Product `bson:"product" json:"product"`
}
