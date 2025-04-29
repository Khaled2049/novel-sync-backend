// File: internal/util/password/password.go
package password

import (
	"fmt"
	"log" // Use structured logging in production

	"golang.org/x/crypto/bcrypt"
)

// HashPassword generates a bcrypt hash of the password.
// cost should typically be bcrypt.DefaultCost or slightly higher.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error generating password hash: %v", err) // Log the actual error internally
		return "", fmt.Errorf("failed to hash password")     // Return generic error
	}
	return string(bytes), nil
}

// CheckPasswordHash compares a plain text password with a stored bcrypt hash.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	// Returns nil on success, bcrypt.ErrMismatchedHashAndPassword on failure
	return err == nil
}