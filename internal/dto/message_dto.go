package dto

import (
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

type InboundMessage struct {
	Type            MessageType      `json:"type"`
	Content         *string          `json:"content,omitempty" binding:"min=1,max=8000"`
	MentionEveryone *bool            `json:"mention_everyone,omitempty"`
	Mentions        *[]uuid.UUID     `json:"mentions,omitempty"` // array of user IDs
	Attachments     *[]AttachmentReq `json:"attachments,omitempty"`

	// These are set by server, not client, BUT KEEP THEM HERE for simplicity
	UserID uuid.UUID `json:"-"`
	RoomID uuid.UUID `json:"-"`
}

type OutboundMessage struct {
	Type MessageType `json:"type"`

	ID               uuid.UUID                   `json:"id"`
	RoomID           uuid.UUID                   `json:"room_id"`
	AuthorID         uuid.UUID                   `json:"author_id"`
	Content          *string                     `json:"content"`
	SentAt           time.Time                   `json:"sent_at"`
	MentionsEveryone bool                        `json:"mentions_everyone"`
	Mentions         []MentionResponseMinimal    `json:"mentions"`
	Attachments      []AttachmentResponseMinimal `json:"attachments"`

	// Optional fields for specific message types
	TypingUser uuid.UUID `json:"typing_user"` // for typing indicators
	Error      string    `json:"error"`       // for error messages

}

type AttachmentReq struct {
	FileName string  `json:"file_name"`
	URL      string  `json:"url"`
	FileType *string `json:"file_type,omitempty"`
	FileSize *int64  `json:"file_size,omitempty"` // in bytes

}

// RESPONSE TYPES FOR MESSAGE
type AttachmentResponseMinimal struct {
	ID        uuid.UUID `json:"id"`
	MessageID uuid.UUID `json:"message_id"`
	URL       string    `json:"URL"`
	FileName  string    `json:"fileName"`
	FileType  *string   `json:"fileType,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type MentionResponseMinimal struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
}
