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

	// System messages (sent by server only)
	MessageTypeJoin  MessageType = "join"
	MessageTypeLeave MessageType = "leave"
	MessageTypeError MessageType = "error"
)

// InboundMessage : InboundMessage is mapped to CreateMessageReq for MessageTypeText
type InboundMessage struct {
	Type            MessageType      `json:"type"`
	Content         *string          `json:"content,omitempty" binding:"min=1,max=8000"`
	SentAt          time.Time        `json:"sent_at" binding:"required"`
	MentionEveryone *bool            `json:"mention_everyone,omitempty"`
	Mentions        *[]uuid.UUID     `json:"mentions,omitempty"` // array of user IDs
	Attachments     *[]AttachmentReq `json:"attachments,omitempty"`

	// These are set by server, not client, BUT KEEP THEM HERE for simplicity
	UserID uuid.UUID `json:"-"`
	RoomID uuid.UUID `json:"-"`
}

type OutboundMessage struct {
	Type MessageType `json:"type"`

	ID               uuid.UUID           `json:"id"`
	RoomID           uuid.UUID           `json:"room_id"`
	AuthorID         uuid.UUID           `json:"author_id"`
	Content          *string             `json:"content"`
	SentAt           time.Time           `json:"sent_at"`
	MentionsEveryone bool                `json:"mentions_everyone"`
	Mentions         []UserBasic         `json:"mentions"`
	Attachments      []models.Attachment `json:"attachments"`

	EditedAt  *time.Time `json:"edited_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`

	// Optional fields for specific message types
	TypingUser uuid.UUID `json:"typing_user"` // for typing indicators
	Error      string    `json:"error"`       // for error messages

}
