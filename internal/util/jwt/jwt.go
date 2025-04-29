package jwt

import (
	"fmt"
	"log" // Use structured logging
	"time"

	"github.com/golang-jwt/jwt/v5" // Using v5

	"github.com/khaled2049/server/internal/config" // Assuming config holds the secret
)

// Claims defines the structure of the JWT claims.
type Claims struct {
	UserID string `json:"userId"`
	// Add other claims like roles, email if needed
	jwt.RegisteredClaims // Includes issuer, expiry, etc.
}

// Generator provides methods to generate JWT tokens.
type Generator struct {
	secretKey []byte
	ttl       time.Duration // Token Time-To-Live
}

// NewGenerator creates a new JWT Generator.
func NewGenerator(cfg *config.JWTConfig) (*Generator, error) {
	if cfg.SecretKey == "" {
		return nil, fmt.Errorf("JWT secret key cannot be empty")
	}
	return &Generator{
		secretKey: []byte(cfg.SecretKey),
		ttl:       cfg.TTL,
	}, nil
}

// GenerateToken creates a new JWT token for a given user ID.
func (g *Generator) GenerateToken(userID string) (string, error) {
	// Set custom claims
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(g.ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "novel-platform-backend", // Optional: identify the issuer
			Subject:   userID,                   // Optional: subject identifies the user
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and return it as a string.
	tokenString, err := token.SignedString(g.secretKey)
	if err != nil {
		log.Printf("Error signing JWT token for user %s: %v", userID, err)
		return "", fmt.Errorf("failed to generate token") // Generic error
	}

	return tokenString, nil
}

// --- Add a corresponding ValidateToken function here later ---
// func (g *Generator) ValidateToken(tokenString string) (*Claims, error) { ... }