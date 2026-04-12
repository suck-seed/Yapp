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
	ID        uuid.UUID  `json:"id"                  db:"id"`
	HallID    uuid.UUID  `json:"hall_id"              db:"hall_id"`
	FloorID   *uuid.UUID `json:"floor_id,omitempty"   db:"floor_id"`
	Name      string     `json:"name"                 db:"name"`
	RoomType  string     `json:"room_type"            db:"room_type"`
	Position  float64    `json:"position"             db:"position"`
	IsPrivate bool       `json:"is_private"           db:"is_private"`
	CreatedAt time.Time  `json:"created_at"           db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"           db:"updated_at"`
}
