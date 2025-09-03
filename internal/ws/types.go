package ws

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
	RoomId    uuid.UUID  `json:"room_id" db:"room_id"`
	HallID    uuid.UUID  `json:"hall_id" db:"hall_id"`
	FloorID   *uuid.UUID `json:"floor_id,omitempty" db:"floor_id"`
	Name      string     `json:"name" db:"name"`
	RoomType  RoomType   `json:"room_type" db:"room_type"`
	IsPrivate bool       `json:"is_private" db:"is_private"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`

	Clients map[uuid.UUID]*Client `json:"clients"`
}

type InboundMessage struct {
	Content  string    `json:"content"`
	RoomId   uuid.UUID `json:"-"`
	AuthorId uuid.UUID `json:"-"`
}

type OutboundMessage struct {
	MessageId       uuid.UUID `json:"message_id"`
	Content         string    `json:"content"`
	RoomId          uuid.UUID `json:"room_id"`
	AuthorId        uuid.UUID `json:"author_id"`
	CreatedAt       time.Time `json:"created_at"`
	MentionEveryone bool      `json:"mention_everyone" db:"mention_everyone"`
	SentAt          time.Time `json:"sent_at" db:"sent_at"`
}
