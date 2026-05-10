package ws

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	dto "github.com/suck-seed/yapp/internal/dto/message"
)

// Client represents one WebSocket connection from one browser tab / mobile app / device.
type Client struct {

	// Represents the current browser ID or instance ID of a user
	// Chrome -> Client 1st
	// Firefox -> New client created for it
	ID uuid.UUID

	// Connection essentials
	Conn *websocket.Conn
	Send chan *dto.OutboundMessage

	UserID uuid.UUID

	// Map RoomID -> HallID
	// The gateway subscribes this one client connection
	// to every room the user can access

	SubscribedRooms map[uuid.UUID]uuid.UUID

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

func (c *Client) IsSubscribedToRoom(roomID uuid.UUID) bool {
	if c == nil || c.SubscribedRooms == nil {
		return false
	}
	_, ok := c.SubscribedRooms[roomID]
	return ok
}

func (c *Client) RoomIDs() []uuid.UUID {
	ids := make([]uuid.UUID, 0, len(c.SubscribedRooms))
	for roomID := range c.SubscribedRooms {
		ids = append(ids, roomID)
	}
	return ids
}

// readPump continuously reads JSON events from this client.
// The client must now send room_id in the JSON payload because /ws is a global gateway.
func (c *Client) readPump(hub *Hub) {

	log.Print("Read Pump Hit")
	defer func() {
		hub.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	_ = c.Conn.SetReadDeadline(time.Now().Add(pongWait))

	c.Conn.SetPongHandler(func(string) error {

		c.LastPing = time.Now()
		_ = c.Conn.SetReadDeadline(time.Now().Add(pongWait))

		if hub.PresenceService != nil {
			_ = hub.PresenceService.RefreshConnection(context.Background(), c.UserID, c.ID)
		}
		return nil
	})

	for {
		inboundMessage := &dto.InboundMessage{}

		if err := c.Conn.ReadJSON(inboundMessage); err != nil {

			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("websocket read error: %v", err)
			}
			break

		}
		log.Printf("Inbound Message readPump: \n %v\n", inboundMessage)
		// Server-owned identity. Never trust these from frontend.
		inboundMessage.UserID = c.UserID
		inboundMessage.ClientID = c.ID

		// Validate if RoomID field is empty or not (should never be empty)
		if inboundMessage.RoomID == uuid.Nil {
			log.Print("Issue Here 1")
			hub.sendErrorToClient(c.ID, uuid.Nil, c.UserID, "room_id is required")
			continue
		}

		if !c.IsSubscribedToRoom(inboundMessage.RoomID) {
			log.Print("Issue Here 2")
			hub.sendErrorToClient(c.ID, inboundMessage.RoomID, c.UserID, "you are not subscribed to this room")
			continue
		}

		if hub.PresenceService != nil {
			_ = hub.PresenceService.RefreshConnection(context.Background(), c.UserID, c.ID)
		}

		log.Printf("Inbound Message: \n %v\n", inboundMessage)

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

			log.Printf("Out Message writePump() : \n %v\n", out)

			if err := c.Conn.WriteJSON(out); err != nil {
				log.Printf("Websocket write wrror %v", err)
				return
			}

		case <-ticker.C:
			_ = c.Conn.SetWriteDeadline(time.Now().Add(writeWait))

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
