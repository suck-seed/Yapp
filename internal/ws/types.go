package ws

import (
	"github.com/suck-seed/yapp/internal/dto"
	"sync"
	"time"

	"github.com/google/uuid"
)

type RoomType string

const (
	AudioRoom RoomType = "audio"
	TextRoom  RoomType = "text"
)

// type Room struct {
// 	RoomID    uuid.UUID  `json:"room_id" db:"room_id"`
// 	HallID    uuid.UUID  `json:"hall_id" db:"hall_id"`
// 	FloorID   *uuid.UUID `json:"floor_id,omitempty" db:"floor_id"`
// 	Name      string     `json:"name" db:"name"`
// 	RoomType  RoomType   `json:"room_type" db:"room_type"`
// 	IsPrivate bool       `json:"is_private" db:"is_private"`
// 	CreatedAt time.Time  `json:"created_at" db:"created_at"`
// 	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`

// 	Clients map[uuid.UUID]*Client `json:"clients"`
// }

// Room represents a chat room in memory
type Room struct {
	ID uuid.UUID `json:"id"`

	Clients   map[uuid.UUID]*Client     `json:"-"`
	Broadcast chan *dto.OutboundMessage `json:"-"`

	// Metadata (fetched once, cached)
	IsPrivate bool      `json:"is_private"`
	CreatedAt time.Time `json:"created_at"`
	mu        sync.RWMutex
}

// What server sends to clients

// TypingIndicator represents typing state
type TypingIndicator struct {
	UserID    uuid.UUID `json:"user_id"`
	RoomID    uuid.UUID `json:"room_id"`
	IsTyping  bool      `json:"is_typing"`
	Timestamp time.Time `json:"timestamp"`
}
