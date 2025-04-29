// File: internal/repository/postgres/user_repo.go
package postgres

import (
	"context"
	"errors"
	"fmt"
	"log" // Use structured logging in production

	"github.com/google/uuid" // For generating internal IDs
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn" // For checking specific error codes
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/khaled2049/server/internal/domain"
	"github.com/khaled2049/server/internal/repository"
)

// postgresUserRepository implements the repository.UserRepository interface.
type postgresUserRepository struct {
	pool *pgxpool.Pool
}

// NewUserRepository creates a new instance of postgresUserRepository.
func NewUserRepository(pool *pgxpool.Pool) repository.UserRepository {
	return &postgresUserRepository{pool: pool}
}

// FindByID retrieves a user by their internal ID.
func (r *postgresUserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	query := `
		SELECT id, firebase_uid, email, name, created_at, updated_at
		FROM users
		WHERE id = $1;`

	user := &domain.User{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.FirebaseUID, &user.Email, &user.Name, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrUserNotFound
		}
		log.Printf("Error scanning user by ID %s: %v", id, err) // Replace with structured logging
		return nil, fmt.Errorf("failed to find user by ID: %w", err)
	}

	return user, nil
}

// FindByFirebaseUID retrieves a user by their Firebase UID.
func (r *postgresUserRepository) FindByFirebaseUID(ctx context.Context, firebaseUID string) (*domain.User, error) {
	query := `
		SELECT id, firebase_uid, email, name, created_at, updated_at
		FROM users
		WHERE firebase_uid = $1;`

	user := &domain.User{}
	err := r.pool.QueryRow(ctx, query, firebaseUID).Scan(
		&user.ID, &user.FirebaseUID, &user.Email, &user.Name, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrUserNotFound
		}
		log.Printf("Error scanning user by Firebase UID %s: %v", firebaseUID, err) // Replace with structured logging
		return nil, fmt.Errorf("failed to find user by Firebase UID: %w", err)
	}

	return user, nil
}

// Create saves a new user to the storage.
func (r *postgresUserRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	// Generate a new UUID for the internal ID if not provided
	if user.ID == "" {
		user.ID = uuid.NewString()
	}

	query := `
		INSERT INTO users (id, firebase_uid, email, name)
		VALUES ($1, $2, $3, $4)
		RETURNING created_at, updated_at;` // Get generated timestamps

	err := r.pool.QueryRow(ctx, query,
		user.ID, user.FirebaseUID, user.Email, user.Name,
	).Scan(&user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		var pgErr *pgconn.PgError
		// Check for unique constraint violation (e.g., duplicate email or firebase_uid)
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // 23505 is unique_violation
			log.Printf("Unique constraint violation during user creation: %v", err)
			return nil, fmt.Errorf("user with provided details already exists: %w", err) // Consider a more specific error type
		}
		log.Printf("Error creating user: %v", err) // Replace with structured logging
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (r *postgresUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, firebase_uid, email, name, password_hash, created_at, updated_at
		FROM users
		WHERE email = $1;`

	user := &domain.User{}
	err := r.pool.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.FirebaseUID, // Make sure your DB has this or handle nullable appropriately
		&user.Email,
		&user.Name,
		&user.PasswordHash, // Scan the password hash
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrUserNotFound // Use the canonical error
		}
		log.Printf("Error scanning user by email %s: %v", email, err)
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}

	return user, nil
}

