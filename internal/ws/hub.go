package ws

import (
	"context"
	"log"

	"sync"
	"time"

	"github.com/google/uuid"
)

type Hub struct {
	// Room management
	Rooms map[uuid.UUID]*Room

	// Connection lifecycle
	Register   chan *Client
	Unregister chan *Client

	// Message processing
	Inbound  chan *InboundMessage
	Outbound chan *OutboundMessage

	// Persistence callback
	PersistFunc PersistFunc

	mu sync.RWMutex
}

// Creates a New Presistent working Hub
func NewHub(presistFunc PersistFunc) Hub {
	return Hub{
		Rooms:       make(map[uuid.UUID]*Room),
		Register:    make(chan *Client, 256),
		Unregister:  make(chan *Client, 256),
		Inbound:     make(chan *InboundMessage, 1024),
		Outbound:    make(chan *OutboundMessage, 1024),
		PersistFunc: presistFunc,
	}
}

// Handles Register, Unregister, Broadcast, Presist
func (h *Hub) Run() {

	go h.handleInboundMessage()
	go h.handleOutbound()

	for {
		select {
		case cl := <-h.Register:
			h.registerClient(cl)

		case cl := <-h.Unregister:
			h.unregisterClient(cl)
		}
	}
}

func (h *Hub) handleInboundMessage() {

	for msg := range h.Inbound {

		switch msg.Type {

		case MessageTypeText:
			h.processTextMessage(msg)

		case MessageTypeTyping:
			h.processTypingIndicator(msg)

		case MessageTypeRead:
			h.processReadReciept(msg)

		default:
			log.Printf("Unknown message type %s", msg.Type)
		}

	}

}

func (h *Hub) handleOutbound() {

	for msg := range h.Outbound {
		h.deliverToRoom(msg.RoomID, msg)

	}

}

func (h *Hub) registerClient(client *Client) {

	h.mu.Lock()
	defer h.mu.Unlock()

	room, exists := h.Rooms[client.RoomID]
	if !exists {
		room = &Room{
			ID:      client.RoomID,
			Clients: make(map[uuid.UUID]*Client),
		}

		h.Rooms[client.RoomID] = room
	}

	room.Clients[client.UserID] = client

	// Notift room of user joining via broadcast channel
	joinMsg := &OutboundMessage{
		Type:      MessageTypeJoin,
		RoomID:    client.RoomID,
		AuthorID:  client.UserID,
		Timestamp: time.Now(),
	}

	// broadcast
	select {
	case h.Outbound <- joinMsg:

	default:
		log.Printf("Could not broadcast join message for user %s", client.UserID)
	}

}

func (h *Hub) unregisterClient(client *Client) {

	h.mu.Lock()
	defer h.mu.Unlock()

	room, exists := h.Rooms[client.RoomID]
	if !exists {
		return
	}

	// if exists
	delete(room.Clients, client.UserID)
	close(client.Send)

	leavingMsg := &OutboundMessage{
		Type:      MessageTypeLeave,
		RoomID:    client.RoomID,
		AuthorID:  client.UserID,
		Timestamp: time.Now(),
	}

	// remove room from memory if no user
	if len(room.Clients) == 0 {
		delete(h.Rooms, client.RoomID)
	} else {

		select {
		case h.Outbound <- leavingMsg:

		default:
			log.Printf("Could not broadcast leave message for user %s", client.UserID)

		}
	}
}

func (h *Hub) processTextMessage(msg *InboundMessage) {

	outboundingMsg, err := h.PersistFunc(context.Background(), msg)
	if err != nil {
		errMsg := &OutboundMessage{
			Type:      MessageTypeError,
			RoomID:    msg.RoomID,
			AuthorID:  msg.UserID,
			Error:     "Failed to send message",
			Timestamp: time.Now(),
		}

		// send error the the concerning user only
		h.sendToUser(msg.UserID, msg.RoomID, errMsg)
		return
	}

	// add to broadcast queue (non blocking)
	select {
	case h.Outbound <- outboundingMsg:

	default:
		log.Printf("Broadcast channel full, dropping message %s", outboundingMsg.ID)

	}

}

func (h *Hub) processTypingIndicator(msg *InboundMessage) {

	typingMsg := &OutboundMessage{
		Type:       MessageTypeTyping,
		RoomID:     msg.RoomID,
		AuthorID:   msg.UserID,
		Timestamp:  time.Now(),
		TypingUser: &msg.UserID,
	}

	select {
	case h.Outbound <- typingMsg:

	default:
		// no need to have anything, typing message can be dropped
	}

	go func() {
		time.Sleep(5 * time.Second)
		stopTyping := &OutboundMessage{
			Type:      MessageTypeStopTyping,
			RoomID:    msg.RoomID,
			AuthorID:  msg.UserID,
			Timestamp: time.Now(),
			// Typing User is nil, thus stopped typing
		}

		select {
		case h.Outbound <- stopTyping:

		default:

		}
	}()

}

func (h *Hub) processReadReciept(msg *InboundMessage) {

}

func (h *Hub) deliverToRoom(roomId uuid.UUID, msg *OutboundMessage) {

	h.mu.RLock()
	room, exists := h.Rooms[roomId]
	defer h.mu.RUnlock()

	if !exists {
		return
	}

	// collect clients to disconnect
	room.mu.RLock()
	var disconnectedClients []uuid.UUID

	// Deliver to all client in room
	for userId, client := range room.Clients {

		select {
		case client.Send <- msg:

		default:
			// Client send buffer is full -
			// Disconnect
			log.Printf("Client %s buffer full, disconnecting", msg.AuthorID)
			close(client.Send)
			disconnectedClients = append(disconnectedClients, userId)

		}
	}

	room.mu.RUnlock()

	// remove any disconnected clients and empty Rooms
	if len(disconnectedClients) > 0 {

		room.mu.Lock()

		if room, exists := h.Rooms[roomId]; exists {

			for _, userId := range disconnectedClients {
				delete(room.Clients, userId)
			}

			isEmpty := len(room.Clients) == 0
			room.mu.Unlock()

			// If no of clients in a room became while removing disconnected clients , remove room from mem
			if isEmpty {
				h.mu.Lock()
				delete(h.Rooms, roomId)
				h.mu.Unlock()
				log.Printf("Cleaned up empty room %s", roomId)
			}
		}

	}

}

func (h *Hub) sendToUser(userID uuid.UUID, roomID uuid.UUID, msg *OutboundMessage) {
	h.mu.RLock()
	room, exists := h.Rooms[roomID]
	if !exists {
		h.mu.RUnlock()
		return
	}

	client, exists := room.Clients[userID]
	h.mu.RUnlock()

	if !exists {
		return
	}

	// sendToUser
	select {

	case client.Send <- msg:

	default:
		// buffer full, cleanup

		close(client.Send)

		h.mu.Lock()
		if room, exists := h.Rooms[roomID]; exists {

			delete(room.Clients, userID)

			// clean up potential empty room
			if len(room.Clients) == 0 {
				delete(h.Rooms, roomID)
			}
		}

		h.mu.Unlock()
	}

}
