package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

const (
	SECRET_KEY = "randomsecretkey"
)

// defines the structure of the JWT payload
type JWTClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateJWT generates a JWT token
func GenerateJWT(userID string, secretKey []byte) (string, error) {
	claims := JWTClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)), // token expires in 24 hours
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// ValidateJWT validates the JWT token and returns the claims if valid
func ValidateJWT(tokenString string, secretKey []byte) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, err
	}

	return claims, nil
}
