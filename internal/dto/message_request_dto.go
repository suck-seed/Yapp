package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateMessageReq struct {

	// we will get author id from headers
	RoomID      uuid.UUID        `json:"room_id" binding:"required"`
	AuthorID    uuid.UUID        `json:"author_id" binding:"required"`
	Content     *string          `json:"content,omitempty" binding:"min=1,max=8000"`
	SentAt      time.Time        `json:"sent_at" binding:"required"`
	Attachments *[]AttachmentReq `json:"attachments,omitempty"`

	MentionEveryone *bool        `json:"mention_everyone,omitempty"`
	Mentions        *[]uuid.UUID `json:"mentions,omitempty"`
}

type AttachmentReq struct {
	FileName string  `json:"file_name"`
	URL      string  `json:"url"`
	FileType *string `json:"file_type,omitempty"`
	FileSize *int64  `json:"file_size,omitempty"` // in bytes

}

// TO REQUEST MESSAGE (PAGINATION)

type MessageQueryParams struct {
	RoomID uuid.UUID
	Limit  int        // default: 50, max: 100
	Before *uuid.UUID // Get messages before this message ID
	After  *uuid.UUID // Get messages after this message ID
	Around *uuid.UUID // Get messages around this message ID
}
