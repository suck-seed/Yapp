package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateMessageReq struct {
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
	FileSize *int64  `json:"file_size" binding:"omitempty"`
}

type FetchMessagesQuery struct {
	Limit  int        `form:"limit"`
	Before *uuid.UUID `form:"before" binding:"omitempty"`
	After  *uuid.UUID `form:"after" binding:"omitempty"`
	Around *uuid.UUID `form:"around" binding:"omitempty"`
}

type MessageQueryParams struct {
	RoomID uuid.UUID  `json:"room_id"`
	Limit  int        `json:"limit"`
	Before *uuid.UUID `json:"before" binding:"omitempty"`
	After  *uuid.UUID `json:"after" binding:"omitempty"`
	Around *uuid.UUID `json:"around" binding:"omitempty"`
}

type UpdateMessageReq struct {
	Content string `json:"content" binding:"required,min=1,max=8000"`
}
