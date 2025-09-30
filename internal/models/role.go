package models

import (
	"time"

	"github.com/google/uuid"
)

type Role struct {
	ID        uuid.UUID `db:"id"`
	HallID    uuid.UUID `db:"hall_id"`
	Name      string    `db:"name"`
	Color     string    `db:"color"`
	IconURL   string    `db:"icon_url"`
	IsDefault bool      `db:"is_default"`
	IsAdmin   bool      `db:"is_admin"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
