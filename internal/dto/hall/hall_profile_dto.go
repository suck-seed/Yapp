package dto

import (
	"time"

	"github.com/google/uuid"
)

// Contains basic information about hall
// - NAME, ICON, DESCRIPTION
// CAN UPDATE USING PROFILE DTO

// GET
type GetHallProfileRes struct {
	ID               uuid.UUID `json:"id"`
	Name             string    `json:"name"`
	IsPrivate        bool      `json:"is_private"`
	IconURL          *string   `json:"icon_url"`
	IconThumbnailURL *string   `json:"icon_thumbnail_url"`
	BannerColor      *string   `json:"banner_color"`
	Description      *string   `json:"description"`
	OwnerID          uuid.UUID `json:"owner_id"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// PATCH
type HallProfileUpdateReq struct {
	Name        *string `json:"name"         binding:"omitempty,min=1,max=100"`
	IsPrivate   *bool   `json:"is_private" binding:"omitempty"`
	Description *string `json:"description"  binding:"omitempty,max=500"`
	BannerColor *string `json:"banner_color" binding:"omitempty"`
}

type HallProfileUpdateRes struct {
	ID               uuid.UUID `json:"id"`
	Name             string    `json:"name"`
	IsPrivate        bool      `json:"is_private"`
	IconURL          *string   `json:"icon_url"`
	IconThumbnailURL *string   `json:"icon_thumbnail_url"`
	BannerColor      *string   `json:"banner_color"`
	Description      *string   `json:"description"`
	OwnerID          uuid.UUID `json:"owner_id"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// currentHalldto
type GetCurrentHallRes struct {
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
}
