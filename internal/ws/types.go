package ws

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
)

type MessageType string

const (
	MessageTypeText       MessageType = "text"
	MessageTypeTyping     MessageType = "typing"
	MessageTypeStopTyping MessageType = "stop_typing"
	MessageTypeRead       MessageType = "read"
	MessageTypeEdit       MessageType = "edit"
	MessageTypeDelete     MessageType = "delete"
	MessageTypeReact      MessageType = "react"

	// System messages (sent by server only)
	MessageTypeJoin  MessageType = "join"
	MessageTypeLeave MessageType = "leave"
	MessageTypeError MessageType = "error"
)

type RoomType string

const (
	AudioRoom RoomType = "audio"
	TextRoom  RoomType = "text"
)

// type Room struct {
// 	RoomId    uuid.UUID  `json:"room_id" db:"room_id"`
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
	RoomId uuid.UUID `json:"room_id"`

	// Active connections - use map for O(1) lookup
	Clients map[uuid.UUID]*Client `json:"-"`

	// Room broadcast channel
	Broadcast chan *OutboundMessage `json:"-"`

	// Metadata (fetched once, cached)
	IsPrivate bool      `json:"is_private"`
	CreatedAt time.Time `json:"created_at"`
	mu        sync.RWMutex
}

// What clients Send to server
type InboundMessage struct {
	Type    MessageType `json:"type"`
	Content string      `json:"content,omitempty"`

	// Mention fields - parsed by frontend
	MentionEveryone bool     `json:"mention_everyone,omitempty"`
	Mentions        []string `json:"mentions,omitempty"` // array of user IDs

	// These are set by server, not client
	UserId uuid.UUID `json:"-"`
	RoomId uuid.UUID `json:"-"`
}

// What server sends to clients
type OutboundMessage struct {
	Type      MessageType `json:"type"`
	MessageId uuid.UUID   `json:"message_id,omitempty"`
	RoomId    uuid.UUID   `json:"room_id"`
	AuthorId  uuid.UUID   `json:"author_id"`
	Content   string      `json:"content,omitempty"`
	Timestamp time.Time   `json:"timestamp"`

	// Optional fields for specific message types
	TypingUser *uuid.UUID `json:"typing_user,omitempty"` // for typing indicators
	Error      string     `json:"error,omitempty"`       // for error messages
}

type PersistFunc func(ctx context.Context, in *InboundMessage) (*OutboundMessage, error)

// TypingIndicator represents typing state
type TypingIndicator struct {
	UserId    uuid.UUID `json:"user_id"`
	RoomId    uuid.UUID `json:"room_id"`
	IsTyping  bool      `json:"is_typing"`
	Timestamp time.Time `json:"timestamp"`
}
