package models

import (
	"github.com/google/uuid"
	"golang.org/x/net/websocket"
)

type Client struct {
	Conn     *websocket.Conn
	Message  chan *Message
	ID       uuid.UUID `json:"id"`
	RoomID   uuid.UUID `json:"room_id"`
	Username string    `json:"username"`
}
