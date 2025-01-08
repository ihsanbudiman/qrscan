package main

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Generate JWT Token for testing
func GenerateJWT(user User) (string, error) {
	expirationTime := time.Now().Add(24 * 30 * time.Hour) // Token valid for 30 days
	claims := &Claims{
		User: user,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// Create the token with the claims and sign it with the secret key
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
