package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AccessClaims struct {
	UserID string `json:"userId"`

	jwt.RegisteredClaims
}

type RefreshClaims struct {
	UserID string `json:"userId"`

	jwt.RegisteredClaims
}

func GenerateAcccessToken(userID string) (string, error) {

	access_secret := os.Getenv("ACCESS_SECRET")

	claims := AccessClaims{
		UserID: userID,
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

func GenerateRefreshToken(userID string) (string, error) {

	refresh_secret := os.Getenv("REFRESH_SECRET")

	claims := RefreshClaims{
		UserID: userID,
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

// func ValidateAccessToken(tokenString string) (jwt.MapClaims, error) {
// 	secret := os.Getenv("JWT_SECRET")

// 	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, jwt.ErrTokenSignatureInvalid
// 		}
// 		return []byte(secret), nil
// 	})
// 	//fmt.Print("token inside validate token", token)

// 	if err != nil {
// 		return nil, err
// 	}

// 	claims, ok := token.Claims.(jwt.MapClaims)

// 	if !ok || !token.Valid {
// 		return nil, jwt.ErrTokenInvalidClaims
// 	}

// 	return claims, nil
// }

// token := &jwt.Token{
// 	Method: jwt.SigningMethodHS256,

// 	Header: map[string]interface{}{
// 		"alg": "HS256",
// 		"typ": "JWT",
// 	},

// 	Claims: jwt.MapClaims{
// 		"user_id": 5,
// 	},
// }
