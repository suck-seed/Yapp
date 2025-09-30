package ws

import (
	"github.com/suck-seed/yapp/internal/dto"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Represents a Websocket connection
type Client struct {
	// Connection essentials
	Conn *websocket.Conn
	Send chan *dto.OutboundMessage

	// Identity - only what's needed for routing
	UserID uuid.UUID
	RoomID uuid.UUID

	// Connection metadata
	ConnectedAt time.Time
	LastPing    time.Time
}

func (c *Client) writePump() {

	ticker := time.NewTicker(60 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {

		select {

		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.Conn.WriteJSON(message); err != nil {
				log.Printf("Websocket write wrror %v", err)
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {

				return
			}

		}

	}

}

// readPump : Continously Reads the upcoming json stream from the websocket connection and passes incoming/inboundMessage to Hub to be handled
func (c *Client) readPump(hub *Hub) {
	defer func() {
		hub.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error {

		c.LastPing = time.Now()
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil

	})

	for {

		inboundMessage := &dto.InboundMessage{}

		if err := c.Conn.ReadJSON(&inboundMessage); err != nil {

			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break

		}

		// We added these values on client in JoinRoom Func
		inboundMessage.UserID = c.UserID
		inboundMessage.RoomID = c.RoomID

		hub.Inbound <- inboundMessage
	}

}
