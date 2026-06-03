package repositories

import (
	"backend/database"
	"backend/models"
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func UserCollection() *mongo.Collection {
	return database.DB.Collection("users")
}

func FindUserByEmail(ctx context.Context, email string) (models.User, error) {

	var user models.User
	err := UserCollection().FindOne(ctx, bson.M{"email": email}).Decode(&user)

	return user, err
}

func CreateUser(ctx context.Context, user models.User) error {

	_, err := UserCollection().InsertOne(ctx, user)

	return err
}

func GetUserByID(ctx context.Context, userID bson.ObjectID) (*models.User, error) {
	var user models.User

	err := UserCollection().FindOne(ctx, bson.M{"_id": userID}).Decode(&user)

	return &user, err
}

func GetAllUsers(ctx context.Context) ([]models.User, error) {

	cursor, err := UserCollection().Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var users []models.User

	err = cursor.All(ctx, &users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func ChangeRole(ctx context.Context, userID bson.ObjectID, role string) (*models.User, error) {

	filter := bson.M{"_id": userID}

	update := bson.M{
		"$set": bson.M{
			"role": role,
		},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updatedUser models.User

	err := UserCollection().FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedUser)

	if err != nil {
		return nil, err
	}

	return &updatedUser, nil

}
