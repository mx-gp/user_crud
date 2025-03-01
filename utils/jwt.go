package utils

import (
	"errors"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Secret key (Use environment variables in production)
var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// Claims structure
type Claims struct {
	UserID uint `json:"user_id"`
	jwt.StandardClaims
}

// ValidateJWT verifies and parses the token
func ValidateJWT(tokenString string) bool {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return false
	}

	// Check if token is valid and not expired
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims.ExpiresAt > time.Now().Unix()
	}

	return false
}
