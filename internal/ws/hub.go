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

	// Handles Presistent handling on incoming messages, and returns a full fiedged cononical outbound message
	// go func() {
	// 	for inboundedMessage := range h.Presist {

	// 		outboundedMessage, err := presist(context.Background(), inboundedMessage)
	// 		if err != nil {

	// 			// Send error to the user
	// 			continue
	// 		}

	// 		h.Broadcast <- outboundedMessage
	// 	}
	// }()

	//
	for {
		select {
		case cl := <-h.Register:
			h.registerClient(cl)

		case cl := <-h.Unregister:
			h.unregisterClient(cl)

			// case m := <-h.Broadcast:

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
		h.deliverToRoom(msg.RoomId, msg)

	}

}

func (h *Hub) registerClient(client *Client) {

	h.mu.Lock()
	defer h.mu.Unlock()

	room, exists := h.Rooms[client.RoomId]
	if !exists {
		room = &Room{
			RoomId:  client.RoomId,
			Clients: make(map[uuid.UUID]*Client),
		}

		h.Rooms[client.RoomId] = room
	}

	room.Clients[client.UserId] = client

	// Notift room of user joining via broadcast channel
	joinMsg := &OutboundMessage{
		Type:      MessageTypeJoin,
		RoomId:    client.RoomId,
		AuthorId:  client.UserId,
		Timestamp: time.Now(),
	}

	// broadcast
	select {
	case h.Outbound <- joinMsg:

	default:
		log.Printf("Could not broadcast join message for user %s", client.UserId)
	}

}

func (h *Hub) unregisterClient(client *Client) {

	h.mu.Lock()
	defer h.mu.Unlock()

	room, exists := h.Rooms[client.RoomId]
	if !exists {
		return
	}

	// if exists
	delete(room.Clients, client.UserId)
	close(client.Send)

	leavingMsg := &OutboundMessage{
		Type:      MessageTypeLeave,
		RoomId:    client.RoomId,
		AuthorId:  client.UserId,
		Timestamp: time.Now(),
	}

	// remove room from memory if no user
	if len(room.Clients) == 0 {
		delete(h.Rooms, client.RoomId)
	} else {

		select {
		case h.Outbound <- leavingMsg:

		default:
			log.Printf("Could not broadcast leave message for user %s", client.UserId)

		}
	}
}

func (h *Hub) processTextMessage(msg *InboundMessage) {

	outboundingMsg, err := h.PersistFunc(context.Background(), msg)
	if err != nil {
		errMsg := &OutboundMessage{
			Type:      MessageTypeError,
			RoomId:    msg.RoomId,
			AuthorId:  msg.UserId,
			Error:     "Failed to send message",
			Timestamp: time.Now(),
		}

		// send error the the concerning user only
		h.sendToUser(msg.UserId, msg.RoomId, errMsg)
		return
	}

	// add to broadcast queue (non blocking)
	select {
	case h.Outbound <- outboundingMsg:

	default:
		log.Printf("Broadcast channel full, dropping message %s", outboundingMsg.MessageId)

	}

}

func (h *Hub) processTypingIndicator(msg *InboundMessage) {

	typingMsg := &OutboundMessage{
		Type:       MessageTypeTyping,
		RoomId:     msg.RoomId,
		AuthorId:   msg.UserId,
		Timestamp:  time.Now(),
		TypingUser: &msg.UserId,
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
			RoomId:    msg.RoomId,
			AuthorId:  msg.UserId,
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
			log.Printf("Client %s buffer full, disconnecting", msg.AuthorId)
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
