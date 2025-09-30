package ws

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/suck-seed/yapp/internal/dto"
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
	ID uuid.UUID `json:"id"`

	Clients   map[uuid.UUID]*Client `json:"-"`
	Broadcast chan *OutboundMessage `json:"-"`

	// Metadata (fetched once, cached)
	IsPrivate bool      `json:"is_private"`
	CreatedAt time.Time `json:"created_at"`
	mu        sync.RWMutex
}

type InboundMessage struct {
	Type        MessageType           `json:"type"`
	Content     *string               `json:"content,omitempty" binding:"min=1,max=8000"`
	Attachments *[]dto.AttachmentType `json:"attachments,omitempty"`

	// Mention fields - parsed by frontend
	MentionEveryone *bool     `json:"mention_everyone,omitempty"`
	Mentions        *[]string `json:"mentions,omitempty"` // array of user IDs

	// These are set by server, not client, BUT KEEP THEM HERE for simplicity
	UserID uuid.UUID `json:"-"`
	RoomID uuid.UUID `json:"-"`
}

// What server sends to clients
type OutboundMessage struct {
	Type      MessageType `json:"type"`
	ID        uuid.UUID   `json:"id,omitempty"`
	RoomID    uuid.UUID   `json:"room_id"`
	AuthorID  uuid.UUID   `json:"author_id"`
	Content   string      `json:"content,omitempty"`
	Timestamp time.Time   `json:"timestamp"`

	// Optional fields for specific message types
	TypingUser *uuid.UUID `json:"typing_user,omitempty"` // for typing indicators
	Error      string     `json:"error,omitempty"`       // for error messages
}

// TypingIndicator represents typing state
type TypingIndicator struct {
	UserID    uuid.UUID `json:"user_id"`
	RoomID    uuid.UUID `json:"room_id"`
	IsTyping  bool      `json:"is_typing"`
	Timestamp time.Time `json:"timestamp"`
}
