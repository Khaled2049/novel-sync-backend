package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/khaled2049/server/internal/domain"
)

type PlaceRepository interface {
	Create(ctx context.Context, place *domain.Place) (*domain.Place, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Place, error)
	Update(ctx context.Context, place *domain.Place) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListByNovelID(ctx context.Context, novelID uuid.UUID) ([]*domain.Place, error)
}