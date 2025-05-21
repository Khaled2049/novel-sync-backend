package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/khaled2049/server/internal/domain"
)

type ChapterRepository interface {
	Create(ctx context.Context, chapter *domain.Chapter) (*domain.Chapter, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Chapter, error)
	Update(ctx context.Context, chapter *domain.Chapter) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListByNovelID(ctx context.Context, novelID uuid.UUID) ([]*domain.Chapter, error)
}
