package domain

import (
	"database/sql"
	"time"
)

type ChapterStatus string

const (
	ChapterStatusDraft     ChapterStatus = "draft"
	ChapterStatusPublished ChapterStatus = "published"
	ChapterStatusArchived  ChapterStatus = "archived"
)

type Chapter struct {
	ID               string   `json:"id"`
	NovelID          string    `json:"novelId"`
	Title            string       `json:"title"`
	Content          string       `json:"content"`
	Status           ChapterStatus `json:"status"`
	OrderIndex       int          `json:"orderIndex"`
	WordCount        int          `json:"wordCount"`
	CreatedAt        time.Time    `json:"createdAt"`
	UpdatedAt        time.Time    `json:"updatedAt"`
	LastEditedByUserID string  `json:"lastEditedByUserId,omitempty"`
	PublishedAt      sql.NullTime   `json:"publishedAt,omitempty"`
}
