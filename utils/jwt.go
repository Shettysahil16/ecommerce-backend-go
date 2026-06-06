package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AccessClaims struct {
	UserID    string `json:"userId"`
	SessionID string `json:"sessionId"`

	jwt.RegisteredClaims
}

type RefreshClaims struct {
	UserID    string `json:"userId"`
	SessionID string `json:"sessionId"`

	jwt.RegisteredClaims
}

// ACCESS TOKEN GENERATION AND VALIDATION
func GenerateAcccessToken(userID string, sessionID string) (string, error) {

	access_secret := os.Getenv("ACCESS_SECRET")

	claims := AccessClaims{
		UserID:    userID,
		SessionID: sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:  userID,
			IssuedAt: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(
				time.Now().Add(15 * time.Minute),
			),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(access_secret))
}

func ValidateAccessToken(tokenString string) (*AccessClaims, error) {
	access_secret := os.Getenv("ACCESS_SECRET")

	token, err := jwt.ParseWithClaims(tokenString, &AccessClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(access_secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*AccessClaims)

	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// REFRESH TOKEN GENERATION AND VALIDATION
func GenerateRefreshToken(userID string, sessionID string) (string, error) {

	refresh_secret := os.Getenv("REFRESH_SECRET")

	claims := RefreshClaims{
		UserID:    userID,
		SessionID: sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:  userID,
			IssuedAt: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(
				time.Now().Add(7 * 24 * time.Hour),
			),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(refresh_secret))
}

func ValidateRefreshToken(tokenString string) (*RefreshClaims, error) {
	refresh_secret := os.Getenv("REFRESH_SECRET")

	token, err := jwt.ParseWithClaims(tokenString, &RefreshClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(refresh_secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*RefreshClaims)

	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func HashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
