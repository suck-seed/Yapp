package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateMessageReq struct {

	// we will get author id from headers
	RoomId          string   `json:"room_id" binding:"required"`
	Content         string   `json:"content" binding:"required,min=1,max=8000"`
	MentionEveryone bool     `json:"mention_everyone"`
	Mentions        []string `json:"mentions" binding:"dive,uuid4"`
}

type CreateMessageRes struct {
}

type FetchMessageByIdReq struct {
	MessageId uuid.UUID `json:"message_id" binding:"required"`
}

type FetchMessageByRoomIDReq struct {
	RoomId string `json:"room_id" binding:"required"`
	Limit  int    `json:"limit" binding:"required"`
	Offset int    `json:"offset" binding:"required"`
}

type FetchRoomMessageReq struct {
	RoomId string    `json:"room_id" binding:"required"`
	Time   time.Time `json:"time" binding:"required"`
	Limit  int       `json:"limit" binding:"required"`
}

// To update
// Fetch 1. Message , check if AuthorId == userId

type UpdateMessageReq struct {
	MessageId       string `json:"message_id" binding:"required"`
	Content         string `json:"content" binding:"required,min=1,max=8000"`
	MentionEveryone bool   `json:"mention_everyone" binding:"omitempty"`
}

// To Delete
// Fetch 1. Message , check if AuthorId == userId
// Or 2. Message, Find user in hall_members, fetch role_id from it, ani check role_permissions ma if text_manage_message is available or not

type DeleteMessageReq struct {
	MessageId string `json:"message_id" binding:"required"`
}
