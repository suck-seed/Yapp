// internal/models/message_read.go
package models

import (
	"time"

	"github.com/google/uuid"
)

type MessageRead struct {
	RoomID    uuid.UUID `json:"room_id" db:"room_id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	MessageID uuid.UUID `json:"message_id" db:"message_id"`

	ReadAt    time.Time `json:"read_at" db:"read_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
