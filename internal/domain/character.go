package domain

import (
	"time"

	"github.com/google/uuid"
)

type Character struct {
	ID                  uuid.UUID  `json:"id"`
	NovelID             uuid.UUID  `json:"novelId"`
	Name                string     `json:"name"`
	Description         string     `json:"description"`
	Backstory           string     `json:"backstory"`
	Motivations         string     `json:"motivations"`
	PhysicalDescription string     `json:"physicalDescription"`
	ImageURL            string     `json:"imageUrl"`
	Source              string     `json:"source"`
	CreatedAt           time.Time  `json:"createdAt"`
	UpdatedAt           time.Time  `json:"updatedAt"`
	CreatedByUserID     *uuid.UUID `json:"createdByUserId,omitempty"`
}

