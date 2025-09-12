package models

import (
	"time"

	"github.com/google/uuid"
)

type RoomType string

const (
	AudioRoom RoomType = "audio"
	TextRoom  RoomType = "text"
)

type Room struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	HallId    uuid.UUID  `json:"hall_id" db:"hall_id"`
	FloorId   *uuid.UUID `json:"floor_id,omitempty" db:"floor_id"`
	Name      string     `json:"name" db:"name"`
	RoomType  RoomType   `json:"room_type" db:"room_type"`
	IsPrivate bool       `json:"is_private" db:"is_private"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
}
