package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/suck-seed/yapp/internal/models"
)

// DTO HERE DEALS WITH I/O STREAM

type MessageType string

const (
	MessageTypeText       MessageType = "text"
	MessageTypeTyping     MessageType = "typing"
	MessageTypeStopTyping MessageType = "stop_typing"
	MessageTypeRead       MessageType = "read"
	MessageTypeEdit       MessageType = "edit"
	MessageTypeDelete     MessageType = "delete"
	MessageTypeReact      MessageType = "react"

	MessageTypePresence MessageType = "presence"

	// Client asks server to refresh this WS client's room subscriptions from DB.
	MessageTypeSyncSubscriptions MessageType = "sync_subscriptions"

	// Server confirms subscriptions were refreshed.
	MessageTypeSubscriptionsSynced MessageType = "subscriptions_synced"

	// System messages (sent by server only)
	MessageTypeJoin  MessageType = "join"
	MessageTypeLeave MessageType = "leave"
	MessageTypeError MessageType = "error"
)

// InboundMessage : InboundMessage is mapped to CreateMessageReq for MessageTypeText
type InboundMessage struct {
	Type MessageType `json:"type"`

	// In new /ws gateway design, frontend must explicitly send RoomID
	RoomID          uuid.UUID        `json:"room_id" binding:"required"`
	Content         *string          `json:"content,omitempty" binding:"min=1,max=8000"`
	SentAt          time.Time        `json:"sent_at" binding:"required"`
	MentionEveryone *bool            `json:"mention_everyone,omitempty"`
	Mentions        *[]uuid.UUID     `json:"mentions,omitempty"` // array of user IDs
	Attachments     *[]AttachmentReq `json:"attachments,omitempty"`

	// For read receipt
	MessageID *uuid.UUID `json:"message_id,omitempty"`

	// Server-owned fields. Never accept these from frontend.
	UserID   uuid.UUID `json:"-"`
	ClientID uuid.UUID `json:"-"`
}

type OutboundMessage struct {

	// To keep the outbounding json minimal as possible
	// Some field are * and omitted if field are empty

	Type MessageType `json:"type"`

	ID       uuid.UUID `json:"id"`
	RoomID   uuid.UUID `json:"room_id"`
	HallID   uuid.UUID `json:"hall_id"`
	AuthorID uuid.UUID `json:"author_id"` // UserID for denoting the user who created the message

	Content          *string             `json:"content,omitempty"`
	SentAt           time.Time           `json:"sent_at"`
	MentionsEveryone bool                `json:"mentions_everyone"`
	Mentions         []UserBasic         `json:"mentions"`
	Attachments      []models.Attachment `json:"attachments"`

	EditedAt  *time.Time `json:"edited_at,omitempty"`  // opt
	DeletedAt *time.Time `json:"deleted_at,omitempty"` // opt
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`

	// Typing
	TypingUser *uuid.UUID `json:"typing_user,omitempty"` // opt

	// Read receipt
	MessageID *uuid.UUID `json:"message_id,omitempty"` // opt
	ReadBy    *uuid.UUID `json:"read_by,omitempty"`    // opt
	ReadAt    *time.Time `json:"read_at,omitempty"`    // opt

	// Pressence
	PresenceUserID *uuid.UUID `json:"presence_user_id,omitempty"`
	PresenceStatus *string    `json:"presence_status,omitempty"`
	LastSeenAt     *time.Time `json:"last_seen_at,omitempty"`

	Error *string `json:"error,omitempty"` // opt

	// New subscription sync response
	SubscribedRoomCount *int                 `json:"subscribed_room_count,omitempty"`
	SubscribedRooms     []SubscribedRoomInfo `json:"subscribed_rooms,omitempty"`
	SyncedAt            *time.Time           `json:"synced_at,omitempty"`
}

type SubscribedRoomInfo struct {
	RoomID uuid.UUID `json:"room_id"`
	HallID uuid.UUID `json:"hall_id"`
}
