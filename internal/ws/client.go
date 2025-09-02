package ws

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/suck-seed/yapp/internal/models"
)

type Client struct {
	Conn             *websocket.Conn
	Message          chan *models.Message
	WebsocketMessage chan *WebsocketMessage
	ID               uuid.UUID `json:"id"`
	RoomID           uuid.UUID `json:"room_id"`
	Username         string    `json:"username"`
}

type WebsocketMessage struct {
	Type      string    `json:"type"` // "message", "user_joined", "user_left", etc.
	Content   string    `json:"content"`
	RoomID    uuid.UUID `json:"room_id"`
	AuthorID  uuid.UUID `json:"author_id"`
	Username  string    `json:"username"`
	MessageID uuid.UUID `json:"message_id,omitempty"` // Only for actual messages
	Timestamp string    `json:"timestamp"`
}
