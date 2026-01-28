package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateHallReq struct {
	Name             string  `json:"name" binding:"required"`
	IconURL          *string `json:"icon_url" binding:"omitempty,url"`
	IconThumbnailURL *string `json:"icon_thumbnail_url" binding:"omitempty,url"`
	BannerColor      *string `json:"banner_color" binding:"omitempty"`
	Description      *string `json:"description" binding:"omitempty,max=500"`
}

type CreateHallRes struct {
	ID               uuid.UUID `json:"id"`
	Name             string    `json:"name"`
	IconURL          *string   `json:"icon_url"`
	IconThumbnailURL *string   `json:"icon_thumbnail_url"`
	BannerColor      *string   `json:"banner_color"`
	Description      *string   `json:"description"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	OwnerID          uuid.UUID `json:"owner_id"`
}
