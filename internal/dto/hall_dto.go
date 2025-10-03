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
	ID               uuid.UUID `json:"id" binding:"required"`
	Name             string    `json:"name" binding:"required"`
	IconURL          *string   `json:"icon_url" binding:"omitempty,url"`
	IconThumbnailURL *string   `json:"icon_thumbnail_url" binding:"omitempty,url"`
	BannerColor      *string   `json:"banner_color" binding:"omitempty"`
	Description      *string   `json:"description" binding:"omitempty,max=500"`
	CreatedAt        time.Time `json:"created_at" binding:"required"`
	UpdatedAt        time.Time `json:"updated_at" binding:"required"`
	Owner            uuid.UUID `json:"owner" binding:"required"`
}
