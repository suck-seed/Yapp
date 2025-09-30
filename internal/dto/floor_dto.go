package dto

import (
	"github.com/google/uuid"
	"time"
)

type CreateFloorReq struct {
	HallID    uuid.UUID `json:"hall_id"`
	Name      string    `json:"name"`
	IsPrivate bool      `json:"is_private"`
}

type CreateFloorRes struct {
	ID        uuid.UUID `json:"id"`
	HallID    uuid.UUID `json:"hall_id"`
	Name      string    `json:"name"`
	IsPrivate bool      `json:"is_private"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
