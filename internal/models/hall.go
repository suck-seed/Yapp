package models

import (
	"time"

	"github.com/google/uuid"
)

type Hall struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	IconURL     *string   `json:"icon_url,omitempty" db:"icon_url"`
	BannerColor *string   `json:"banner_color,omitempty" db:"banner_color"`
	Description *string   `json:"description,omitempty" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`

	// createdby
	CreatedBy uuid.UUID `json:"created_by_id" db:"created_by_id"`
}
