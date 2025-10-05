package models

import (
	"time"

	"github.com/google/uuid"
)

type Hall struct {
	ID               uuid.UUID `json:"id" db:"id"`
	Name             string    `json:"name" db:"name"`
	IconURL          *string   `json:"icon_url,omitempty" db:"icon_url"`
	IconThumbnailURL *string   `json:"icon_thumbnail_url,omitempty" db:"icon_thumbnail_url"`
	BannerColor      *string   `json:"banner_color,omitempty" db:"banner_color"`
	Description      *string   `json:"description,omitempty" db:"description"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`

	// createdby
	OwnerID uuid.UUID `json:"owner_id" db:"owner_id"`
}
