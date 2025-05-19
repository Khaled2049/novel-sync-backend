package service

import (
	"context"
	"fmt"

	"github.com/khaled2049/server/internal/domain"
	"github.com/khaled2049/server/internal/repository"
)



type NovelService struct {
	novelRepo repository.NovelRepository
	chapterRepo repository.ChapterRepository
}

func NewNovelService(
		novelRepo repository.NovelRepository,
	 	chapterRepo repository.ChapterRepository) *NovelService {
	return &NovelService{
		novelRepo: novelRepo,
		chapterRepo: chapterRepo,
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
		NovelID:         createdNovel.ID,
		Title:           chapterTitle,
		Content:         initialContent,
		Status:          domain.ChapterStatusDraft,
		OrderIndex:      0, 
		WordCount:       0, 
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
	novelID string,
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
	chapter.NovelID = novelID
	chapter.OrderIndex = highestIndex + 1
	chapter.WordCount = len(chapter.Content) 
	chapter.Status = domain.ChapterStatusDraft 
	chapter.LastEditedByUserID = chapter.LastEditedByUserID 

	fmt.Printf("DEBUG: Adding chapter with OrderIndex = %d\n", chapter.LastEditedByUserID)

	return s.chapterRepo.Create(ctx, chapter)
}

func (s *NovelService) AutosaveChapter(
	ctx context.Context,
	chapterID string,
	content string,
	userID string,
) error {
	return s.chapterRepo.AutosaveContent(ctx, chapterID, content, userID)
}

// SaveChapterWithRevision saves chapter content and creates a new revision entry
func (s *NovelService) SaveChapterWithRevision(
	ctx context.Context,
	chapterID string,
	newContent string,
	userID string,
	notes string,
) error {
	// First update the chapter content
	chapter, err := s.chapterRepo.GetByID(ctx, chapterID)
	if err != nil {
		return fmt.Errorf("failed to get chapter for revision: %w", err)
	}
	
	// Update the chapter
	chapter.Content = newContent
	chapter.LastEditedByUserID = userID
	
	err = s.chapterRepo.Update(ctx, chapter)
	if err != nil {
		return fmt.Errorf("failed to update chapter: %w", err)
	}
	
	// Then save a revision
	err = s.chapterRepo.SaveRevision(ctx, chapterID, newContent, userID, notes)
	if err != nil {
		return fmt.Errorf("chapter updated but failed to save revision: %w", err)
	}
	
	return nil
}
