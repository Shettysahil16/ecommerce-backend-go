package services

import (
	"backend/cache"
	"backend/models"
	"backend/repositories"
	"backend/utils"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
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

func LoginUser(ctx context.Context, user models.LoginRequest) (string, string, error) {

	existingUser, err := repositories.FindUserByEmail(ctx, user.Email)
	if err != nil {
		return "", "", errors.New(
			"invalid email or password",
		)
	}

	passwordValid := utils.CheckPasswordHash(user.Password, existingUser.Password)
	if !passwordValid {
		return "", "", errors.New(
			"invalid email or password",
		)
	}

	sessionID := uuid.NewString()

	accessToken, err := utils.GenerateAcccessToken(existingUser.ID.Hex(), sessionID)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := utils.GenerateRefreshToken(existingUser.ID.Hex(), sessionID)
	if err != nil {
		return "", "", err
	}

	hash := utils.HashToken(refreshToken)

	session := models.Session{
		SessionID:   sessionID,
		UserID:      existingUser.ID.Hex(),
		RefreshHash: hash,
	}

	err = cache.SetSessionCache(ctx, sessionID, session)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
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
