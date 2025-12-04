package ws

import (
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	dto "github.com/suck-seed/yapp/internal/dto/message"
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

	// Channel state tracking
	closed bool
	mu     sync.Mutex
}

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10 // send pings a bit before the pong deadline
	maxMessageSize = 512 << 10           // 512KB (set what you want)
)

// readPump : Continously Reads the upcoming json stream from the websocket connection and passes incoming/inboundMessage to Hub to be handled
func (c *Client) readPump(hub *Hub) {
	defer func() {
		hub.Unregister <- c
		c.Conn.Close()
	}()

	// c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {

		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
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

func (c *Client) writePump() {

	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	c.Conn.EnableWriteCompression(true)

	for {

		select {

		case out, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))

			if !ok {
				_ = c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.Conn.WriteJSON(out); err != nil {
				log.Printf("Websocket write wrror %v", err)
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))

			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {

				return
			}

		}

	}

}

func (c *Client) SafeClose() {

	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.closed {
		close(c.Send)
		c.closed = true
	}

}
