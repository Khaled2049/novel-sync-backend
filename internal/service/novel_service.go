package service

import (
	"context"

	"github.com/khaled2049/server/internal/domain"
	"github.com/khaled2049/server/internal/repository"
)



type NovelService struct {
	novelRepo repository.NovelRepository
}

func NewNovelService(novelRepo repository.NovelRepository) *NovelService {
	return &NovelService{
		novelRepo: novelRepo,
	}
}

func (s *NovelService) GetNovelByID(ctx context.Context, id string) (*domain.Novel, error) {
	novel, err := s.novelRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return novel, nil
}

// Create
func (s *NovelService) CreateNovel(ctx context.Context, novel *domain.Novel) (*domain.Novel, error) {
	createdNovel, err := s.novelRepo.Create(ctx, novel)
	if err != nil {
		return nil, err
	}
	return createdNovel, nil
}

// get all novels
func (s *NovelService) GetAllNovels(ctx context.Context) ([]*domain.Novel, error) {
	novels, err := s.novelRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return novels, nil
}