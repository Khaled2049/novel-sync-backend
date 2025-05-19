package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/khaled2049/server/internal/domain"
)


type CharacterRepository interface {
	Create(ctx context.Context, character *domain.Character) (*domain.Character, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Character, error)
	Update(ctx context.Context, character *domain.Character) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListByNovelID(ctx context.Context, novelID uuid.UUID) ([]*domain.Character, error)
}

