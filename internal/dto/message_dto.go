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
	Type            MessageType       `json:"type"`
	Content         *string           `json:"content,omitempty" binding:"min=1,max=8000"`
	MentionEveryone *bool             `json:"mention_everyone,omitempty"`
	Mentions        *[]uuid.UUID      `json:"mentions,omitempty"` // array of user IDs
	Attachments     *[]AttachmentType `json:"attachments,omitempty"`

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
	SentAt           time.Time                   `json:"sentAt"`
	MentionsEveryone bool                        `json:"mentionsEveryone"`
	Mentions         []MentionResponseMinimal    `json:"mentions"`
	Attachments      []AttachmentResponseMinimal `json:"attachments"`

	// Optional fields for specific message types
	TypingUser uuid.UUID `json:"typing_user"` // for typing indicators
	Error      string    `json:"error"`       // for error messages

}

type CreateMessageReq struct {

	// we will get author id from headers
	RoomID      uuid.UUID `json:"room_id" binding:"required"`
	AuthorID    uuid.UUID
	Content     *string           `json:"content,omitempty" binding:"min=1,max=8000"`
	Attachments *[]AttachmentType `json:"attachments,omitempty"`

	MentionEveryone *bool        `json:"mention_everyone,omitempty"`
	Mentions        *[]uuid.UUID `json:"mentions,omitempty"`
}

type AttachmentType struct {
	FileName string  `json:"file_name"`
	URL      string  `json:"url"`
	FileType *string `json:"file_type,omitempty"`
	FileSize *int64  `json:"file_size,omitempty"` // in bytes

}

type MentionType struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
}

// RESPONSE TO MESSAGE CREATION
type AttachmentResponseMinimal struct {
	ID       uuid.UUID `json:"id"`
	URL      string    `json:"URL"`
	FileName string    `json:"fileName"`
	FileType *string   `json:"fileType,omitempty"`
}

type MentionResponseMinimal struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
}

type CreateMessageRes struct {
	ID               uuid.UUID                   `json:"id"`
	RoomID           uuid.UUID                   `json:"roomID"`
	AuthorID         uuid.UUID                   `json:"authorID"`
	Content          *string                     `json:"content"`
	SentAt           time.Time                   `json:"sentAt"`
	MentionsEveryone bool                        `json:"mentionsEveryone"`
	Mentions         []MentionResponseMinimal    `json:"mentions"`
	Attachments      []AttachmentResponseMinimal `json:"attachments"`
}

type FetchMessageByIdReq struct {
	ID uuid.UUID `json:"id" binding:"required"`
}

type FetchMessageByRoomIDReq struct {
	RoomId uuid.UUID `json:"room_id" binding:"required"`
	Limit  int       `json:"limit" binding:"required"`
	Offset int       `json:"offset" binding:"required"`
}

type FetchRoomMessageReq struct {
	RoomId uuid.UUID `json:"room_id" binding:"required"`
	Time   time.Time `json:"time" binding:"required"`
	Limit  int       `json:"limit" binding:"required"`
}

// To update
// Fetch 1. Message , check if AuthorId == userId

type UpdateMessageReq struct {
	ID              uuid.UUID `json:"id" binding:"required"`
	Content         string    `json:"content" binding:"required,min=1,max=8000"`
	MentionEveryone bool      `json:"mention_everyone" binding:"omitempty"`
}

// To Delete
// Fetch 1. Message , check if AuthorId == userId
// Or 2. Message, Find user in hall_members, fetch role_id from it, ani check role_permissions ma if text_manage_message is available or not

type DeleteMessageReq struct {
	ID uuid.UUID `json:"id" binding:"required"`
}
