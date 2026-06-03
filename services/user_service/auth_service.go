package services

import (
	"backend/models"
	"backend/repositories"
	"backend/utils"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func RegisterUser(ctx context.Context, user models.User) error {

	existingUser, _ := repositories.FindUserByEmail(ctx, user.Email)

	if existingUser.Email != "" {
		return errors.New("email already exists")
	}

	hashedPassword, err := utils.HashPassword(user.Password)

	if err != nil {
		return err
	}

	user.Password = hashedPassword
	user.Role = "general"
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	return repositories.CreateUser(ctx, user)

}

func LoginUser(ctx context.Context, user models.LoginRequest) (string, error) {

	existingUser, err := repositories.FindUserByEmail(ctx, user.Email)
	if err != nil {
		return "", errors.New(
			"invalid email or password",
		)
	}

	passwordValid := utils.CheckPasswordHash(user.Password, existingUser.Password)
	if !passwordValid {
		return "", errors.New(
			"invalid email or password",
		)
	}

	token, err := utils.GenerateToken(existingUser.ID.Hex())
	if err != nil {
		return "", err
	}

	return token, nil
}

func GetUserByID(ctx context.Context, userID string) (*models.User, error) {

	objectID, err := bson.ObjectIDFromHex(userID)

	if err != nil {
		return &models.User{}, err
	}
	return repositories.GetUserByID(ctx, objectID)
}

func GetAllUsers(ctx context.Context) ([]models.User, error) {

	return repositories.GetAllUsers(ctx)
}

func ChangeRole(ctx context.Context, userID string, role string) (*models.User, error) {

	objID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	return repositories.ChangeRole(ctx, objID, role)
}
