package models

import (
	"time"

	"github.com/google/uuid"
)

type HallBan struct {
	ID     uuid.UUID `db:"id" json:"id"`
	Reason string    `db:"reason" json:"reason"`

	UserID uuid.UUID `db:"user_id" json:"user_id"`
	HallID uuid.UUID `db:"hall_id" json:"hall_id"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
