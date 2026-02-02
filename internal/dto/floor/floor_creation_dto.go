package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateFloorReq struct {
	HallID    uuid.UUID `json:"hall_id" binding:"required"`
	Name      string    `json:"name" binding:"required"`
	IsPrivate bool      `json:"is_private" binding:"required"`
}

type CreateFloorRes struct {
	ID        uuid.UUID `json:"id"`
	HallID    uuid.UUID `json:"hall_id"`
	Name      string    `json:"name"`
	IsPrivate bool      `json:"is_private"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
