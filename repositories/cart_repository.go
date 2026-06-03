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

func CartCollection() *mongo.Collection {
	return database.DB.Collection("cartproducts")
}

func UpsertCartProduct(ctx context.Context, userID string, productID bson.ObjectID) (*models.Cart, error) {

	now := time.Now()

	filter := bson.M{
		"userId":    userID,
		"productId": productID,
	}

	update := bson.M{
		"$inc": bson.M{
			"quantity": 1,
		},
		"$set": bson.M{
			"updatedAt": now,
		},
		"$setOnInsert": bson.M{
			"_id":       bson.NewObjectID(),
			"userId":    userID,
			"productId": productID,
			"createdAt": now,
		},
	}

	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	var updatedCart models.Cart

	err := CartCollection().FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedCart)

	if err != nil {
		return nil, err
	}

	return &updatedCart, nil
}

func GetCartItems(ctx context.Context, userID string) ([]models.Cart, error) {

	var cartProducts []models.Cart

	cursor, err := CartCollection().Find(ctx, bson.M{"userId": userID})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	err = cursor.All(ctx, &cartProducts)
	if err != nil {
		return nil, err
	}

	return cartProducts, nil
}

func DeleteCartItem(ctx context.Context, cartID bson.ObjectID) error {

	result, err := CartCollection().DeleteOne(ctx, bson.M{"_id": cartID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("cart item not found")
	}

	return nil
}

func ReduceCartItem(ctx context.Context, userID string, productID bson.ObjectID) (*models.Cart, error) {

	filter := bson.M{
		"userId":    userID,
		"productId": productID,
	}

	update := bson.M{
		"$inc": bson.M{
			"quantity": -1,
		},

		"$set": bson.M{
			"updateAt": time.Now(),
		},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updatedCart models.Cart

	err := CartCollection().FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedCart)
	if err != nil {
		return nil, err
	}

	if updatedCart.Quantity <= 0 {
		_, err := CartCollection().DeleteOne(ctx, bson.M{"_id": updatedCart.ID})
		if err != nil {
			return nil, err
		}
	}

	return &updatedCart, nil
}
