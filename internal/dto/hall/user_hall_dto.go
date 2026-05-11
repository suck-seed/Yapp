package dto

import (
	"time"

	"github.com/google/uuid"
)

type UserHallRes struct {
	ID               uuid.UUID `json:"id"`
	Name             string    `json:"name"`
	IsPrivate        bool      `json:"is_private"`
	IconURL          *string   `json:"icon_url"`
	IconThumbnailURL *string   `json:"icon_thumbnail_url"`
	BannerColor      *string   `json:"banner_color"`
	Description      *string   `json:"description"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	OwnerID          uuid.UUID `json:"owner_id"`

	IsPinned bool     `json:"is_pinned"`
	Position *float64 `json:"position"` // null when not pinned
}

// after_id = nil means place at very top of pinned halls.
type MovePinnedHallReq struct {
	AfterID *uuid.UUID `json:"after_id" binding:"omitempty"`
}
