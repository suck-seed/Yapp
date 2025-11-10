package dto

import (
	"github.com/google/uuid"
)

type CreateMessageReq struct {

	// we will get author id from headers
	RoomID      uuid.UUID        `json:"room_id" binding:"required"`
	AuthorID    uuid.UUID        `json:"author_id" binding:"required"`
	Content     *string          `json:"content,omitempty" binding:"min=1,max=8000"`
	Attachments *[]AttachmentReq `json:"attachments,omitempty"`

	MentionEveryone *bool        `json:"mention_everyone,omitempty"`
	Mentions        *[]uuid.UUID `json:"mentions,omitempty"`
}

type MessageQueryParams struct {
	RoomID uuid.UUID
	Limit  int        // default: 50, max: 100
	Before *uuid.UUID // Get messages before this message ID
	After  *uuid.UUID // Get messages after this message ID
	Around *uuid.UUID // Get messages around this message ID
}
