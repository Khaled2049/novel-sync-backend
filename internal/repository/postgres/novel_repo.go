// File: internal/repository/postgres/novel_repo.go
package postgres

import (
	"context"
	"errors"
	"fmt"
	"log" // Consider using structured logging

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/khaled2049/server/internal/domain"
	"github.com/khaled2049/server/internal/repository"
)

// postgresNovelRepository implements the repository.NovelRepository interface.
type postgresNovelRepository struct {
	pool *pgxpool.Pool
}

// NewNovelRepository creates a new instance of postgresNovelRepository.
func NewNovelRepository(pool *pgxpool.Pool) repository.NovelRepository {
	return &postgresNovelRepository{pool: pool}
}

// func gets all novels
func (r *postgresNovelRepository) GetAll(ctx context.Context) ([]*domain.Novel, error) {
	query := `
		SELECT id, owner_user_id, title, logline, description, genre, visibility, cover_image_url, created_at, updated_at
		FROM novels
		ORDER BY updated_at DESC;`

	var novels []*domain.Novel
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		log.Printf("Error querying all novels: %v", err) // Replace with structured logging
		return nil, fmt.Errorf("failed to find all novels: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		novel := &domain.Novel{}
		err := rows.Scan(
			&novel.ID, &novel.OwnerUserID, &novel.Title, &novel.Logline, &novel.Description,
			&novel.Genre, &novel.Visibility, &novel.CoverImageURL, &novel.CreatedAt, &novel.UpdatedAt,
		)
		if err != nil {
			log.Printf("Error scanning novel row: %v", err) // Replace with structured logging
			return nil, fmt.Errorf("failed to scan novel: %w", err)
		}
		novels = append(novels, novel)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over novel rows: %v", err) // Replace with structured logging
		return nil, fmt.Errorf("failed to iterate novel rows: %w", err)
	}

	return novels, nil
}

// Create saves a new novel to the storage.
func (r *postgresNovelRepository) Create(ctx context.Context, novel *domain.Novel) (*domain.Novel, error) {
	query := `
		INSERT INTO novels (
			owner_user_id, title, logline, description, genre, visibility, cover_image_url
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7
		)
		RETURNING id, created_at, updated_at;`

	err := r.pool.QueryRow(ctx, query,
		novel.OwnerUserID, novel.Title, novel.Logline, novel.Description, novel.Genre, novel.Visibility, novel.CoverImageURL,
	).Scan(&novel.ID, &novel.CreatedAt, &novel.UpdatedAt)

	if err != nil {
		log.Printf("Error creating novel: %v", err) // Replace with structured logging
		return nil, fmt.Errorf("failed to create novel: %w", err)
	}

	return novel, nil
}

// FindByID retrieves a novel by its internal ID.
func (r *postgresNovelRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Novel, error) {
	query := `
		SELECT id, owner_user_id, title, logline, description, genre, visibility, cover_image_url, created_at, updated_at
		FROM novels
		WHERE id = $1;`

	novel := &domain.Novel{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&novel.ID, &novel.OwnerUserID, &novel.Title, &novel.Logline, &novel.Description,
		&novel.Genre, &novel.Visibility, &novel.CoverImageURL, &novel.CreatedAt, &novel.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNovelNotFound
		}
		log.Printf("Error scanning novel by ID %s: %v", id, err) // Replace with structured logging
		return nil, fmt.Errorf("failed to find novel by ID: %w", err)
	}

	return novel, nil
}

// FindByOwnerID retrieves all novels owned by a specific user ID.
func (r *postgresNovelRepository) FindByOwnerID(ctx context.Context, ownerID string) ([]*domain.Novel, error) {
	query := `
		SELECT id, owner_user_id, title, logline, description, genre, visibility, cover_image_url, created_at, updated_at
		FROM novels
		WHERE owner_user_id = $1
		ORDER BY updated_at DESC;`

	var novels []*domain.Novel
	rows, err := r.pool.Query(ctx, query, ownerID)
	if err != nil {
		log.Printf("Error querying novels by owner ID %s: %v", ownerID, err) // Replace with structured logging
		return nil, fmt.Errorf("failed to find novels by owner ID: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		novel := &domain.Novel{}
		err := rows.Scan(
			&novel.ID, &novel.OwnerUserID, &novel.Title, &novel.Logline, &novel.Description,
			&novel.Genre, &novel.Visibility, &novel.CoverImageURL, &novel.CreatedAt, &novel.UpdatedAt,
		)
		if err != nil {
			log.Printf("Error scanning novel row: %v", err) // Replace with structured logging
			return nil, fmt.Errorf("failed to scan novel: %w", err)
		}
		novels = append(novels, novel)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over novel rows: %v", err) // Replace with structured logging
		return nil, fmt.Errorf("failed to iterate novel rows: %w", err)
	}

	return novels, nil
}

// Update updates an existing novel in the storage.
func (r *postgresNovelRepository) Update(ctx context.Context, novel *domain.Novel) (*domain.Novel, error) {
	query := `
		UPDATE novels
		SET title = $2,
			logline = $3,
			description = $4,
			genre = $5,
			visibility = $6,
			cover_image_url = $7,
			updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at;`

	err := r.pool.QueryRow(ctx, query,
		novel.ID, novel.Title, novel.Logline, novel.Description, novel.Genre, novel.Visibility, novel.CoverImageURL,
	).Scan(&novel.UpdatedAt)

	if err != nil {
		log.Printf("Error updating novel with ID %s: %v", novel.ID, err) // Replace with structured logging
		return nil, fmt.Errorf("failed to update novel: %w", err)
	}

	return novel, nil
}

// Delete removes a novel from the storage by its ID.
func (r *postgresNovelRepository) Delete(ctx context.Context, id string) error {
	query := `
		DELETE FROM novels
		WHERE id = $1;`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		log.Printf("Error deleting novel with ID %s: %v", id, err) // Replace with structured logging
		return fmt.Errorf("failed to delete novel: %w", err)
	}

	rowsAffected := result.RowsAffected()

	if rowsAffected == 0 {
		return repository.ErrNovelNotFound
	}

	return nil
}

// Search retrieves novels matching a search query.
func (r *postgresNovelRepository) Search(ctx context.Context, query string, limit, offset int) ([]*domain.Novel, error) {
	sqlQuery := `
		SELECT id, owner_user_id, title, logline, description, genre, visibility, cover_image_url, created_at, updated_at
		FROM novels
		WHERE title ILIKE $1 OR logline ILIKE $1 OR description ILIKE $1
		ORDER BY updated_at DESC
		LIMIT $2 OFFSET $3;`

	searchPattern := "%" + query + "%"
	var novels []*domain.Novel
	rows, err := r.pool.Query(ctx, sqlQuery, searchPattern, limit, offset)
	if err != nil {
		log.Printf("Error searching novels with query '%s': %v", query, err) // Replace with structured logging
		return nil, fmt.Errorf("failed to search novels: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		novel := &domain.Novel{}
		err := rows.Scan(
			&novel.ID, &novel.OwnerUserID, &novel.Title, &novel.Logline, &novel.Description,
			&novel.Genre, &novel.Visibility, &novel.CoverImageURL, &novel.CreatedAt, &novel.UpdatedAt,
		)
		if err != nil {
			log.Printf("Error scanning search result row: %v", err) // Replace with structured logging
			return nil, fmt.Errorf("failed to scan novel: %w", err)
		}
		novels = append(novels, novel)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over search result rows: %v", err) // Replace with structured logging
		return nil, fmt.Errorf("failed to iterate search result rows: %w", err)
	}

	return novels, nil
}

// FindCollaborativeNovels retrieves all novels a user is collaborating on.
func (r *postgresNovelRepository) FindCollaborativeNovels(ctx context.Context, userID string) ([]*domain.Novel, error) {
	query := `
		SELECT n.id, n.owner_user_id, n.title, n.logline, n.description, n.genre, n.visibility, n.cover_image_url, n.created_at, n.updated_at
		FROM novels n
		JOIN novel_collaborators nc ON n.id = nc.novel_id
		WHERE nc.user_id = $1
		ORDER BY n.updated_at DESC;`

	var novels []*domain.Novel
	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		log.Printf("Error finding collaborative novels for user %s: %v", userID, err) // Replace with structured logging
		return nil, fmt.Errorf("failed to find collaborative novels: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		novel := &domain.Novel{}
		err := rows.Scan(
			&novel.ID, &novel.OwnerUserID, &novel.Title, &novel.Logline, &novel.Description,
			&novel.Genre, &novel.Visibility, &novel.CoverImageURL, &novel.CreatedAt, &novel.UpdatedAt,
		)
		if err != nil {
			log.Printf("Error scanning collaborative novel row: %v", err) // Replace with structured logging
			return nil, fmt.Errorf("failed to scan novel: %w", err)
		}
		novels = append(novels, novel)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over collaborative novel rows: %v", err) // Replace with structured logging
		return nil, fmt.Errorf("failed to iterate collaborative novel rows: %w", err)
	}

	return novels, nil
}

// AddCollaborator adds a user as a collaborator to a novel.
func (r *postgresNovelRepository) AddCollaborator(ctx context.Context, novelID, userID, role string) error {
	query := `
		INSERT INTO novel_collaborators (novel_id, user_id, role)
		VALUES ($1, $2, $3)
		ON CONFLICT (novel_id, user_id) DO UPDATE SET role = $3;`

	_, err := r.pool.Exec(ctx, query, novelID, userID, role)
	if err != nil {
		log.Printf("Error adding collaborator %s to novel %s with role %s: %v", userID, novelID, role, err) // Replace with structured logging
		return fmt.Errorf("failed to add collaborator: %w", err)
	}

	return nil
}

// RemoveCollaborator removes a user as a collaborator from a novel.
func (r *postgresNovelRepository) RemoveCollaborator(ctx context.Context, novelID, userID string) error {
	query := `
		DELETE FROM novel_collaborators
		WHERE novel_id = $1 AND user_id = $2;`

	result, err := r.pool.Exec(ctx, query, novelID, userID)
	if err != nil {
		log.Printf("Error removing collaborator %s from novel %s: %v", userID, novelID, err) // Replace with structured logging
		return fmt.Errorf("failed to remove collaborator: %w", err)
	}

	rowsAffected := result.RowsAffected()

	if rowsAffected == 0 {
		return errors.New("collaborator not found for this novel")
	}

	return nil
}
