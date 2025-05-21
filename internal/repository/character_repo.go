package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/khaled2049/server/internal/domain"
)

// CharacterRepository defines the interface for character data operations
type CharacterRepository interface {
	// Core CRUD operations
	Create(ctx context.Context, character *domain.Character) (*domain.Character, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Character, error)
	Update(ctx context.Context, character *domain.Character) error
	Delete(ctx context.Context, id uuid.UUID) error

	// Novel-specific operations
	ListByNovelID(ctx context.Context, novelID uuid.UUID) ([]*domain.Character, error)

	// Search operations
	SearchByName(ctx context.Context, novelID uuid.UUID, nameQuery string) ([]*domain.Character, error)
}
