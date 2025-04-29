// File: internal/domain/user.go
package domain

import "time"

// User represents a user in the system.
type User struct {
	ID           string    `json:"id"`
	FirebaseUID  string    `json:"-"` // Usually not exposed
	Email        string    `json:"email"`
	Name         string    `json:"name"`
	PasswordHash string    `json:"-"` // Store the hash, NEVER the plain password. Exclude from JSON.
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}