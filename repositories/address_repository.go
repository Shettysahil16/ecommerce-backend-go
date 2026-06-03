package repositories

import (
	"backend/database"
	"backend/models"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func AddressCollection() *mongo.Collection {
	return database.DB.Collection("address")
}

func GetAddress(ctx context.Context, addressID bson.ObjectID, userID bson.ObjectID) (*models.DefaultAddressResponse, error) {

	var address models.DefaultAddressResponse

	err := AddressCollection().FindOne(ctx, bson.M{"_id": addressID, "userId": userID}).Decode(&address)
	if err != nil {
		return nil, err
	}

	return &address, nil
}

func GetAddresses(ctx context.Context, userID bson.ObjectID) (*[]models.Address, error) {

	var addresses []models.Address

	cursor, err := AddressCollection().Find(ctx, bson.M{"userId": userID})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	err = cursor.All(ctx, &addresses)
	if err != nil {
		return nil, err
	}

	return &addresses, nil
}

func CreateAddress(ctx context.Context, AddressData models.Address, userID bson.ObjectID) (*models.Address, error) {

	AddressData.ID = bson.NewObjectID()
	AddressData.UserID = userID

	_, err := AddressCollection().InsertOne(ctx, AddressData)
	if err != nil {
		return nil, err
	}

	return &AddressData, nil
}

func UpdateAddress(ctx context.Context, addressID bson.ObjectID, addressData models.AddressResponse) (*models.AddressResponse, error) {

	filter := bson.M{
		"_id": addressID,
	}

	var updateFields bson.M

	data, _ := bson.Marshal(addressData)
	bson.Unmarshal(data, &updateFields)

	update := bson.M{
		"$set": updateFields,
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updatedAddress models.AddressResponse

	err := AddressCollection().FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedAddress)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("address not found")
		}

		return nil, err
	}

	return &updatedAddress, nil
}

func UpdateManyAddress(ctx context.Context, addressID bson.ObjectID, userID bson.ObjectID) error {

	_, err := AddressCollection().UpdateMany(
		ctx,
		bson.M{
			"userId": userID,
			"_id": bson.M{
				"$ne": addressID,
			},
			"isDefault": true,
		},
		bson.M{
			"$set": bson.M{
				"isDefault": false,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func UpdateDefaultAddress(ctx context.Context, addressID bson.ObjectID, userID bson.ObjectID) (*models.DefaultAddressResponse, error) {

	filter := bson.M{
		"_id":    addressID,
		"userId": userID,
	}

	update := bson.M{
		"$set": bson.M{
			"isDefault": true,
		},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updatedAddress models.DefaultAddressResponse

	err := AddressCollection().FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedAddress)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("address not found")
		}

		return nil, err
	}

	return &updatedAddress, nil
}

func DeleteAddress(ctx context.Context, userID bson.ObjectID, addressID bson.ObjectID) error {

	filter := bson.M{
		"_id":    addressID,
		"userId": userID,
	}

	result, err := AddressCollection().DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("address not found")
	}

	return err
}

func CountAddresses(ctx context.Context, userID bson.ObjectID) (int64, error) {

	totalAddresses, err := AddressCollection().CountDocuments(ctx, bson.M{"userId": userID})
	if err != nil {
		return 0, err
	}

	return totalAddresses, nil
}

func GetFirstAddress(ctx context.Context, userID bson.ObjectID, addressID bson.ObjectID) (*models.Address, error) {

	var firstAddress models.Address

	err := AddressCollection().FindOne(
		ctx,
		bson.M{
			"userId": userID,
			"_id": bson.M{
				"$ne": addressID,
			},
			"isDefault": false,
		},
		options.FindOne().SetSort(bson.M{"_id": 1}),
	).Decode(&firstAddress)

	if err != nil {
		return nil, err
	}

	return &firstAddress, nil

}

func GetFirstAddressUpdate(ctx context.Context, userID bson.ObjectID, addressID bson.ObjectID) (*models.Address, error) {

	filter := bson.M{
		"userId": userID,
		"_id":    bson.M{"$ne": addressID},
	}

	update := bson.M{
		"$set": bson.M{"isDefault": true},
	}

	options := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updatedAddress models.Address

	err := AddressCollection().FindOneAndUpdate(ctx, filter, update, options).Decode(&updatedAddress)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &updatedAddress, nil
}
