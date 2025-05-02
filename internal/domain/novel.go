package domain

import (
	"time"
)

// Novel represents a writing project in the platform
type Novel struct {
	ID            string          `json:"id" db:"id"`
	Title         string          `json:"title" db:"title"`
	Logline       string          `json:"logline,omitempty" db:"logline"`
	Description   string          `json:"description,omitempty" db:"description"`
	Genre         string          `json:"genre,omitempty" db:"genre"`
	Visibility    NovelVisibility `json:"visibility" db:"visibility"`
	OwnerUserID   string          `json:"owner_user_id" db:"owner_user_id"`
	CoverImageURL string          `json:"cover_image_url,omitempty" db:"cover_image_url"`
	CreatedAt     time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at" db:"updated_at"`
}


type NovelVisibility string

const (
	NovelVisibilityPrivate NovelVisibility = "private"
	NovelVisibilityInviteOnly NovelVisibility = "invite_only"
	NovelVisibilityPublic NovelVisibility = "public"
)
