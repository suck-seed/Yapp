package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateMessageReq struct {

	// we will get author id from headers
	RoomID      uuid.UUID         `json:"room_id" binding:"required"`
	AuthorID    uuid.UUID         `json:"author_id" binding:"required"`
	Content     *string           `json:"content,omitempty" binding:"min=1,max=8000"`
	Attachments *[]AttachmentType `json:"attachments,omitempty"`

	MentionEveryone *bool        `json:"mention_everyone,omitempty"`
	Mentions        *[]uuid.UUID `json:"mentions,omitempty"`
}

type FetchMessageByIdReq struct {
	ID uuid.UUID `json:"id" binding:"required"`
}

type FetchMessageByRoomIDReq struct {

}


type FetchRoomMessageReq struct {
	RoomID uuid.UUID  `json:"room_id" binding:"required"`
	UserID uuid.UUID  `json:"user_id" binding:"required"` // For authorization (for private room and type shi)
	Before *time.Time `json:"before,omitempty"`           // For pagination
	Limit  int        `json:"limit" binding:"required,min=1,max=100"`
}

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
