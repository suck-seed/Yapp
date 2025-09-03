package models

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	MessageId       uuid.UUID  `json:"message_id" db:"message_id"`
	RoomId          uuid.UUID  `json:"room_id" db:"room_id"`
	AuthorId        uuid.UUID  `json:"author_id" db:"author_id"`
	Content         string     `json:"content" db:"content"`
	SentAt          time.Time  `json:"sent_at" db:"sent_at"`
	EditedAt        *time.Time `json:"edited_at" db:"edited_at"`
	DeletedAt       *time.Time `json:"deleted_at" db:"deleted_at"`
	MentionEveryone bool       `json:"mention_everyone" db:"mention_everyone"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`
}
