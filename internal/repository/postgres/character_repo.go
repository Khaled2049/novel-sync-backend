package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/khaled2049/server/internal/domain"
	"github.com/khaled2049/server/internal/repository"
)

var ErrCharacterNotFound = errors.New("character not found")
var ErrNovelNotFound = errors.New("novel not found")

type postgresCharacterRepository struct {
	pool *pgxpool.Pool
}

func NewCharacterRepository(pool *pgxpool.Pool) repository.CharacterRepository {
	return &postgresCharacterRepository{pool: pool}
}

func (r *postgresCharacterRepository) Create(ctx context.Context, character *domain.Character) (*domain.Character, error) {
	// Make sure we have a novel ID

	fmt.Println("Creating character with ID:", character.ID)

	if character.NovelID == uuid.Nil {
		return nil, errors.New("novel ID is required")
	}

	if character.ID == uuid.Nil {
		character.ID = uuid.New()
	}

	query := `
		INSERT INTO characters (
			id, novel_id, name, description, backstory, motivations, 
			physical_description, image_url, source, created_by_user_id
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
		) RETURNING id, novel_id, name, description, backstory, motivations, 
			physical_description, image_url, source, created_at, updated_at, created_by_user_id
	`

	row := r.pool.QueryRow(
		ctx,
		query,
		character.ID,
		character.NovelID,
		character.Name,
		character.Description,
		character.Backstory,
		character.Motivations,
		character.PhysicalDescription,
		character.ImageURL,
		character.Source,
		character.CreatedByUserID,
	)

	err := row.Scan(
		&character.ID,
		&character.NovelID,
		&character.Name,
		&character.Description,
		&character.Backstory,
		&character.Motivations,
		&character.PhysicalDescription,
		&character.ImageURL,
		&character.Source,
		&character.CreatedAt,
		&character.UpdatedAt,
		&character.CreatedByUserID,
	)

	if err != nil {
		return nil, fmt.Errorf("error creating character: %w", err)
	}

	return character, nil
}

func (r *postgresCharacterRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Character, error) {
	query := `
		SELECT 
			id, novel_id, name, description, backstory, motivations, 
			physical_description, image_url, source, created_at, updated_at, created_by_user_id
		FROM characters
		WHERE id = $1
	`

	character := &domain.Character{}
	row := r.pool.QueryRow(ctx, query, id)

	err := row.Scan(
		&character.ID,
		&character.NovelID,
		&character.Name,
		&character.Description,
		&character.Backstory,
		&character.Motivations,
		&character.PhysicalDescription,
		&character.ImageURL,
		&character.Source,
		&character.CreatedAt,
		&character.UpdatedAt,
		&character.CreatedByUserID,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCharacterNotFound
		}
		return nil, fmt.Errorf("error getting character: %w", err)
	}

	return character, nil
}

func (r *postgresCharacterRepository) Update(ctx context.Context, character *domain.Character) error {
	query := `
		UPDATE characters
		SET 
			name = $1,
			description = $2,
			backstory = $3,
			motivations = $4,
			physical_description = $5,
			image_url = $6
		-- novel_id, source, created_by_user_id are not updated here
		-- updated_at is handled by the trigger
		WHERE id = $7
		RETURNING updated_at 
	`

	err := r.pool.QueryRow(
		ctx,
		query,
		character.Name,
		character.Description,
		character.Backstory,
		character.Motivations,
		character.PhysicalDescription,
		character.ImageURL,
		character.ID,
	).Scan(&character.UpdatedAt) // Scan the returned updated_at

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrCharacterNotFound // Character to update was not found
		}
		return fmt.Errorf("error updating character: %w", err)
	}

	return nil
}

func (r *postgresCharacterRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM characters WHERE id = $1`

	commandTag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting character: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return ErrCharacterNotFound
	}

	return nil
}

func (r *postgresCharacterRepository) ListByNovelID(ctx context.Context, novelID uuid.UUID) ([]*domain.Character, error) {
	query := `
		SELECT 
			id, novel_id, name, description, backstory, motivations, 
			physical_description, image_url, source, created_at, updated_at, created_by_user_id
		FROM characters
		WHERE novel_id = $1
		ORDER BY name
	`

	rows, err := r.pool.Query(ctx, query, novelID)
	if err != nil {
		return nil, fmt.Errorf("error listing characters: %w", err)
	}
	defer rows.Close()

	var characters []*domain.Character
	for rows.Next() {
		character := &domain.Character{}
		err := rows.Scan(
			&character.ID,
			&character.NovelID,
			&character.Name,
			&character.Description,
			&character.Backstory,
			&character.Motivations,
			&character.PhysicalDescription,
			&character.ImageURL,
			&character.Source,
			&character.CreatedAt,
			&character.UpdatedAt,
			&character.CreatedByUserID,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning character row: %w", err)
		}
		characters = append(characters, character)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating character rows: %w", err)
	}

	return characters, nil
}

func (r *postgresCharacterRepository) SearchByName(ctx context.Context, novelID uuid.UUID, nameQuery string) ([]*domain.Character, error) {
	query := `
		SELECT 
			id, novel_id, name, description, backstory, motivations, 
			physical_description, image_url, source, created_at, updated_at, created_by_user_id
		FROM characters
		WHERE novel_id = $1 AND name ILIKE $2
		ORDER BY name
		LIMIT 20
	`

	rows, err := r.pool.Query(ctx, query, novelID, "%"+nameQuery+"%")
	if err != nil {
		return nil, fmt.Errorf("error searching characters: %w", err)
	}
	defer rows.Close()

	var characters []*domain.Character
	for rows.Next() {
		character := &domain.Character{}
		err := rows.Scan(
			&character.ID,
			&character.NovelID,
			&character.Name,
			&character.Description,
			&character.Backstory,
			&character.Motivations,
			&character.PhysicalDescription,
			&character.ImageURL,
			&character.Source,
			&character.CreatedAt,
			&character.UpdatedAt,
			&character.CreatedByUserID,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning character search row: %w", err)
		}
		characters = append(characters, character)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating character search rows: %w", err)
	}

	return characters, nil
}
