package models

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID              uuid.UUID `json:"id" db:"id"`
	RoomID          uuid.UUID `json:"room_id" db:"room_id"`
	AuthorID        uuid.UUID `json:"author_id" db:"author_id"`
	Content         *string   `json:"content,omitempty" db:"content"`
	MentionEveryone bool      `json:"mention_everyone" db:"mention_everyone"`

	SentAt    time.Time  `json:"sent_at" db:"sent_at"`
	EditedAt  *time.Time `json:"edited_at,omitempty" db:"edited_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
}

type Attachment struct {
	ID        uuid.UUID `json:"id" db:"id"`
	MessageID uuid.UUID `json:"message_id"`
	FileName  string    `json:"file_name"`
	URL       string    `json:"url"`
	FileType  *string   `json:"file_type,omitempty"`
	FileSize  *int64    `json:"file_size,omitempty"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Reaction struct {
	ID        uuid.UUID `json:"id"`
	MessageID uuid.UUID `json:"message_id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	Emoji     string    `json:"emoji"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
