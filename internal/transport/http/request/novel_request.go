package request

import "github.com/khaled2049/server/internal/domain"

// CreateNovelWithFirstChapterRequest defines the payload for creating a novel with its first chapter.
type CreateNovelWithFirstChapterRequest struct {
	NovelData      domain.Novel `json:"novel_data"`
	ChapterTitle   string       `json:"chapter_title"`
	InitialContent string       `json:"initial_content"`
	// UserID is implicitly taken from NovelData.OwnerUserID by the service
}

type AddChapterToNovelRequest struct {
	Title           string `json:"title" binding:"required"`
	Content         string `json:"content"` // Initial content can be empty
	// Status, OrderIndex, WordCount will be set by the service or repository
	LastEditedByUserID string `json:"last_edited_by_user_id" binding:"required"` // Or get from auth context
}

// AutosaveChapterRequest defines the payload for autosaving chapter content.
type AutosaveChapterRequest struct {
	Content string `json:"content"`
	UserID  string `json:"user_id" binding:"required"` // Or get from auth context
}
