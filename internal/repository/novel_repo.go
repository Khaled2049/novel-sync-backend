package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/khaled2049/server/internal/domain"
)

var ErrNovelNotFound = errors.New("novel not found")

type NovelRepository interface {
	Create(ctx context.Context, novel *domain.Novel) (*domain.Novel, error)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Novel, error)
	Update(ctx context.Context, novel *domain.Novel) (*domain.Novel, error)
	Delete(ctx context.Context, id string) error
	GetAll(ctx context.Context) ([]*domain.Novel, error)

	// GetByOwner(ctx context.Context, ownerID string) ([]*domain.Novel, error)
	// Search(query string, limit, offset int) ([]*domain.Novel, error)
	// GetCollaborativeNovels(userID string) ([]*domain.Novel, error)
}
