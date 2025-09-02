package models

import (
	"time"

	"github.com/google/uuid"
)

type Floor struct {
	ID        uuid.UUID `json:"floor_id" db:"floor_id"`
	HallID    uuid.UUID `json:"hall_id" db:"hall_id"`
	Name      string    `json:"name" db:"name"`
	IsPrivate bool      `json:"is_private" db:"is_private"`

	// Position  int       `json:"position" db:"position"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
