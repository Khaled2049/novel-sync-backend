package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/khaled2049/server/internal/domain"
	"github.com/khaled2049/server/internal/repository"
)

type NovelService struct {
	novelRepo     repository.NovelRepository
	chapterRepo   repository.ChapterRepository
	characterRepo repository.CharacterRepository
}

func NewNovelService(
	novelRepo repository.NovelRepository,
	chapterRepo repository.ChapterRepository,
	characterRepo repository.CharacterRepository) *NovelService {
	return &NovelService{
		novelRepo:     novelRepo,
		chapterRepo:   chapterRepo,
		characterRepo: characterRepo,
	}
}

func (s *NovelService) GetNovelByID(ctx context.Context, id uuid.UUID) (*domain.Novel, error) {
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

func (s *NovelService) CreateNovelWithFirstChapter(
	ctx context.Context,
	novel *domain.Novel,
	chapterTitle string,
	initialContent string,
	userID string,
) (*domain.Novel, *domain.Chapter, error) {
	createdNovel, err := s.novelRepo.Create(ctx, novel)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create novel: %w", err)
	}

	fmt.Printf("DEBUG: createdNovel.ID = %s\n", createdNovel.ID)

	chapter := &domain.Chapter{
		NovelID:            createdNovel.ID,
		Title:              chapterTitle,
		Content:            initialContent,
		Status:             domain.ChapterStatusDraft,
		OrderIndex:         0,
		WordCount:          0,
		LastEditedByUserID: createdNovel.OwnerUserID,
	}

	createdChapter, err := s.chapterRepo.Create(ctx, chapter)
	if err != nil {
		return createdNovel, nil, fmt.Errorf("novel created but failed to create first chapter: %w", err)
	}

	return createdNovel, createdChapter, nil
}

// AddChapterToNovel adds a new chapter to an existing novel
func (s *NovelService) AddChapterToNovel(
	ctx context.Context,
	novelID uuid.UUID,
	chapter *domain.Chapter,
) (*domain.Chapter, error) {
	// Verify the novel exists
	_, err := s.novelRepo.GetByID(ctx, novelID)
	if err != nil {
		return nil, fmt.Errorf("novel not found: %w", err)
	}

	// Get the highest order index to append at the end
	chapters, err := s.chapterRepo.ListByNovelID(ctx, novelID)
	if err != nil {
		return nil, fmt.Errorf("failed to get chapters: %w", err)
	}

	highestIndex := 0
	for _, ch := range chapters {
		if ch.OrderIndex > highestIndex {
			highestIndex = ch.OrderIndex
		}
	}

	// Set the new chapter's order index
	chapter.NovelID = novelID.String()
	chapter.OrderIndex = highestIndex + 1
	chapter.WordCount = len(chapter.Content)
	chapter.Status = domain.ChapterStatusDraft

	fmt.Printf("DEBUG: Adding chapter with LastEditedByUserID = %s\n", chapter.LastEditedByUserID)

	return s.chapterRepo.Create(ctx, chapter)
}

// CreateCharacter creates a new character and associates it with a novel
func (s *NovelService) CreateCharacter(
	ctx context.Context,
	novelID uuid.UUID,
	character *domain.Character,
) (*domain.Character, error) {
	// Verify the novel exists
	_, err := s.novelRepo.GetByID(ctx, novelID)
	if err != nil {
		return nil, fmt.Errorf("novel not found: %w", err)
	}

	character.NovelID = novelID
	character.Source = "user"

	return s.characterRepo.Create(ctx, character)
}
