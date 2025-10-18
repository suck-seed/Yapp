package dto

import (
	"time"

	"github.com/google/uuid"
)

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

type GetRoomMessagesResponse struct {
	Messages   []EnrichedMessage `json:"messages"`
	HasMore    bool              `json:"has_more"`
	NextCursor *time.Time        `json:"next_cursor,omitempty"`
}

type EnrichedMessage struct {
	ID               uuid.UUID                   `json:"id"`
	RoomID           uuid.UUID                   `json:"room_id"`
	AuthorID         uuid.UUID                   `json:"author_id"`
	AuthorUsername   string                      `json:"author_username"`
	AuthorAvatar     *string                     `json:"author_avatar"`
	Content          *string                     `json:"content"`
	SentAt           time.Time                   `json:"sent_at"`
	EditedAt         *time.Time                  `json:"edited_at"`
	MentionsEveryone bool                        `json:"mentions_everyone"`
	Mentions         []MentionResponseMinimal    `json:"mentions"`
	Attachments      []AttachmentResponseMinimal `json:"attachments"`
}
