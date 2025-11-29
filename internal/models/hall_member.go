package models

import (
	"time"

	"github.com/google/uuid"
)

type HallMember struct {
	ID        uuid.UUID `db:"id"`
	HallID    uuid.UUID `db:"hall_id"`
	UserID    uuid.UUID `db:"user_id"`
	RoleID    uuid.UUID `db:"role_id"`
	JoinedAt  time.Time `db:"joined_at"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
