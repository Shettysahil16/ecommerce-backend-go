package repositories

import (
	"backend/database"
	"backend/models"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func ProductCollection() *mongo.Collection {
	return database.DB.Collection("products")
}

func GetProductByID(ctx context.Context, productID bson.ObjectID) (*models.Product, error) {
	var product models.Product

	err := ProductCollection().FindOne(ctx, bson.M{"_id": productID}).Decode(&product)

	return &product, err
}

func GetProductsByIDs(ctx context.Context, ids []bson.ObjectID) ([]models.Product, error) {

	cursor, err := ProductCollection().Find(ctx, bson.M{
		"_id": bson.M{"$in": ids},
	},
		options.Find().SetProjection(bson.M{
			"description": 0,
			"createdAt":   0,
			"updatedAt":   0,
			"productImage": bson.M{
				"$slice": 1,
			},
		}),
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []models.Product
	if err := cursor.All(ctx, &products); err != nil {
		return nil, err
	}

	return products, nil
}

func GetSingleCategoryProduct(ctx context.Context) ([]bson.M, error) {

	pipeline := bson.A{
		bson.M{
			"$group": bson.M{
				"_id": "$category",
				"image": bson.M{
					"$first": bson.M{
						"$arrayElemAt": bson.A{"$productImage", 0},
					},
				},
			},
		},
		bson.M{
			"$project": bson.M{
				"_id":      0,
				"category": "$_id",
				"image":    1,
			},
		},
	}

	cursor, err := ProductCollection().Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var results []bson.M
	//var raw []bson.M

	err = cursor.All(ctx, &results)
	if err != nil {
		return nil, err
	}
	//fmt.Printf("%+v\n", raw)
	return results, nil
}

func GetCategoryProducts(ctx context.Context, category string) ([]bson.M, error) {

	pipeline := mongo.Pipeline{
		bson.D{
			{
				Key:   "$match",
				Value: bson.M{"category": category},
			},
		},
		bson.D{
			{
				Key: "$project",
				Value: bson.M{
					"_id":          0,
					"productName":  1,
					"brandName":    1,
					"category":     1,
					"price":        1,
					"sellingPrice": 1,
					"productImage": bson.M{
						"$arrayElemAt": bson.A{"$productImage", 0},
					},
				}},
		},
	}

	cursor, err := ProductCollection().Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var products []bson.M

	err = cursor.All(ctx, &products)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func GetProducts(ctx context.Context) ([]models.Product, error) {

	var products []models.Product

	cursor, err := ProductCollection().Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	err = cursor.All(ctx, &products)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func UpdateProduct(ctx context.Context, productId bson.ObjectID, updateData map[string]interface{}) (*models.Product, error) {

	var updatedProduct models.Product

	filter := bson.M{"_id": productId}

	updateData["updatedAt"] = time.Now()

	update := bson.M{
		"$set": updateData,
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	err := ProductCollection().FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedProduct)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("product not found")
		}

		return nil, err
	}

	return &updatedProduct, nil

}

func CreateProduct(ctx context.Context, productData *models.Product) (*models.Product, error) {

	productData.ID = bson.NewObjectID()

	now := time.Now()

	productData.CreatedAt = now
	productData.UpdatedAt = now

	_, err := ProductCollection().InsertOne(ctx, productData)

	if err != nil {
		return nil, err
	}

	return productData, nil
}
