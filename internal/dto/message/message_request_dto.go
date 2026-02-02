package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateMessageReq struct {
	// we will get author id from headers
	RoomID      uuid.UUID        `json:"room_id" binding:"required"`
	AuthorID    uuid.UUID        `json:"author_id" binding:"required"`
	Content     *string          `json:"content" binding:"omitempty,min=1,max=8000"`
	SentAt      time.Time        `json:"sent_at" binding:"required"`
	Attachments *[]AttachmentReq `json:"attachments" binding:"omitempty"`

	MentionEveryone *bool        `json:"mention_everyone" binding:"omitempty"`
	Mentions        *[]uuid.UUID `json:"mentions" binding:"omitempty"`
}

type AttachmentReq struct {
	FileName string  `json:"file_name"`
	URL      string  `json:"url"`
	FileType *string `json:"file_type" binding:"omitempty"`
	FileSize *int64  `json:"file_size" binding:"omitempty"` // in bytes
}

// TO REQUEST MESSAGE (PAGINATION)
type MessageQueryParams struct {
	RoomID uuid.UUID  `json:"room_id"`
	Limit  int        `json:"limit"`                      // default: 50, max: 100
	Before *uuid.UUID `json:"before" binding:"omitempty"` // Get messages before this message ID
	After  *uuid.UUID `json:"after" binding:"omitempty"`  // Get messages after this message ID
	Around *uuid.UUID `json:"around" binding:"omitempty"` // Get messages around this message ID
}
