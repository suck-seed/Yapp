package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/suck-seed/yapp/internal/models"
)

type CreateMessageRes struct {
	ID               uuid.UUID           `json:"id"`
	RoomID           uuid.UUID           `json:"room_id"`
	AuthorID         uuid.UUID           `json:"author_id"`
	Content          *string             `json:"content"`
	SentAt           time.Time           `json:"sent_at"`
	MentionsEveryone bool                `json:"mentions_everyone"`
	Mentions         []UserBasic         `json:"mentions"`
	Attachments      []models.Attachment `json:"attachments"`

	EditedAt  *time.Time `json:"edited_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type AttachmentResponseMinimal struct {
	ID        uuid.UUID `json:"id"`
	MessageID uuid.UUID `json:"message_id"`
	URL       string    `json:"url"`
	FileName  string    `json:"file_name"`
	FileType  *string   `json:"file_type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type MentionResponseMinimal struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
}

type UserBasic struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	Email     string    `json:"email" db:"email"`
	AvatarURL *string   `json:"avatar_url" db:"avatar_url"`
}

type ReactionGroup struct {
	Emoji    string      `json:"emoji"`
	Count    int         `json:"count"`
	Reactors []UserBasic `json:"user_ids"`
}

type ReactionRes struct {
	MessageID uuid.UUID `json:"message_id"`
	UserID    uuid.UUID `json:"user_id"`
	Emoji     string    `json:"emoji"`
}

type MessageDetailed struct {
	models.Message
	Author      UserBasic                   `json:"author"`
	Attachments []AttachmentResponseMinimal `json:"attachments"`
	Reactions   []ReactionGroup             `json:"reactions"`
	Mentions    []UserBasic                 `json:"mentions"`
}

type MessageListResponse struct {
	Messages []*MessageDetailed `json:"messages"`
	HasMore  bool               `json:"has_more"`
}
