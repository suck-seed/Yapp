package dto

import "github.com/google/uuid"

type CreateHallReq struct {
	ID          uuid.UUID `json:"id" binding:"required"`
	Name        string    `json:"name" binding:"required"`
	IconURL     *string   `json:"icon_url" binding:"omitempty,url"`
	BannerColor *string   `json:"banner_color" binding:"omitempty"`
	Description *string   `json:"description" binding:"omitempty,max=500"`
	CreatedBy   uuid.UUID `json:"created_by" binding:"required"`
}

type CreateHallRes struct {
	ID          string  `json:"id" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	IconURL     *string `json:"icon_url" binding:"omitempty,url"`
	BannerColor *string `json:"banner_color" binding:"omitempty"`
	Description *string `json:"description" binding:"omitempty,max=500"`
	CreatedAt   string  `json:"created_at" binding:"required"`
	UpdatedAt   string  `json:"updated_at" binding:"required"`
	CreatedBy   string  `json:"created_by" binding:"required"`
}
