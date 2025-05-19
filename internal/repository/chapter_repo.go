package repository

import (
	"context"

	"github.com/khaled2049/server/internal/domain"
)



type ChapterRepository interface {
	Create(ctx context.Context, chapter *domain.Chapter) (*domain.Chapter, error)
	GetByID(ctx context.Context, id string) (*domain.Chapter, error)
	Update(ctx context.Context, chapter *domain.Chapter) error
	Delete(ctx context.Context, id string) error
	ListByNovelID(ctx context.Context, novelID string) ([]*domain.Chapter, error)
	SaveRevision(ctx context.Context, chapterID string, content string, userID string, notes string) error
	// GetRevisions(ctx context.Context, chapterID string) ([]*domain., error)
	AutosaveContent(ctx context.Context, chapterID string, content string, userID string) error
}