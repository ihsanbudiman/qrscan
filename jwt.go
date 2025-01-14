package main

import (
	"log"
	"os"
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

func GenerateJWTForWebsocket(channel string) string {
	// Get the secret key from the environment variable
	secret := os.Getenv("CENTRIFUGO_TOKEN_HMAC_SECRET_KEY")
	if secret == "" {
		log.Fatal("CENTRIFUGO_TOKEN_HMAC_SECRET_KEY is not set")
	}

	// Define claims
	claims := jwt.MapClaims{
		"sub": channel,                               // Channel name
		"exp": time.Now().Add(24 * time.Hour).Unix(), // Token expiry
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Fatalf("Failed to sign token: %v", err)
	}

	return signedToken
}
