package repository

import (
	"context"
	"errors" // Use standard errors package

	"github.com/khaled2049/server/internal/domain"
)

// ErrUserNotFound is returned when a user is not found.
var ErrUserNotFound = errors.New("user not found")

// UserRepository defines the interface for interacting with user storage.
type UserRepository interface {
	// FindByID retrieves a user by their internal ID.
	FindByID(ctx context.Context, id string) (*domain.User, error)

	// FindByFirebaseUID retrieves a user by their Firebase UID.
	FindByFirebaseUID(ctx context.Context, firebaseUID string) (*domain.User, error)

	// Create saves a new user to the storage.
	// It should populate the User.ID and Timestamps.
	Create(ctx context.Context, user *domain.User) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	// Update modifies an existing user's details.
	// Update(ctx context.Context, user *domain.User) error
}