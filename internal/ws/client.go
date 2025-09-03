package ws

import (
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/suck-seed/yapp/internal/models"
)

type Client struct {
	Conn    *websocket.Conn
	Message chan *models.Message

	RoomId uuid.UUID `json:"room_id"`

	UserId      uuid.UUID `json:"user_id"`
	Username    string    `json:"username"`
	DisplayName *string   `json:"display_name,omitempty" db:"display_name"`
	AvatarURL   *string   `json:"avatar_url,omitempty" db:"avatar_url"`
	Description *string   `json:"description,omitempty" db:"description"`
	Active      bool      `json:"active" db:"active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

func (c *Client) writeMessage() {

	defer c.Conn.Close()

	for {

		message, ok := <-c.Message
		if !ok {
			return
		}

		// TODO check profanity
		// TODO send to messageService to
		//

		c.Conn.WriteJSON(message)

	}

}

func (c *Client) readMessage(hub *Hub) {
	defer func() {
		hub.Unregister <- c
		c.Conn.Close()
	}()

	for {

		msg := &models.Message{}

		err := c.Conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break

		}

		msg.AuthorId = c.UserId
		msg.RoomId = c.RoomId

		hub.Broadcast <- msg
	}

}
