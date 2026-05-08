package ws

import (
	"context"
	"log"

	"sync"
	"time"

	"github.com/google/uuid"
	dto "github.com/suck-seed/yapp/internal/dto/message"
	"github.com/suck-seed/yapp/internal/services"
)

type Hub struct {
	// Room management
	Rooms map[uuid.UUID]*Room

	// Connection lifecycle
	Register   chan *Client
	Unregister chan *Client

	// Message processing
	Inbound  chan *dto.InboundMessage
	Outbound chan *dto.OutboundMessage

	// Persistence callback
	PersistFunc     PersistFunction
	ReadReceiptFunc ReadReceiptFunction

	// Presence Service
	PresenceService services.IPresenceService

	mu sync.RWMutex
}

// Creates a New Presistent working Hub
func NewHub(
	p PersistFunction,
	readFunc ReadReceiptFunction,
	presenceService services.IPresenceService,
) Hub {
	return Hub{
		Rooms:           make(map[uuid.UUID]*Room),
		Register:        make(chan *Client, 1024),
		Unregister:      make(chan *Client, 1024),
		Inbound:         make(chan *dto.InboundMessage, 1024),
		Outbound:        make(chan *dto.OutboundMessage, 1024),
		PersistFunc:     p,
		ReadReceiptFunc: readFunc,
		PresenceService: presenceService,
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

	// passed from client/readPump()
	// contains all info about the inbound message
	for inboundMessage := range h.Inbound {

		switch inboundMessage.Type {

		case dto.MessageTypeText:
			h.processTextMessage(inboundMessage)

		case dto.MessageTypeTyping:
			h.processTypingIndicator(inboundMessage)

		case dto.MessageTypeStopTyping:
			h.processStopTypingIndicator(inboundMessage)

		case dto.MessageTypeRead:
			h.processReadReciept(inboundMessage)

		default:
			log.Printf("Unknown message type %s", inboundMessage.Type)
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
			HallID:  client.HallID,
			Clients: make(map[uuid.UUID]*Client),
		}

		h.Rooms[client.RoomID] = room
	}

	room.Clients[client.ID] = client

	presence, err := h.PresenceService.MarkConnected(context.Background(), client.UserID, client.ID)
	if err == nil && presence != nil {
		userID := client.UserID

		presenceMsg := &dto.OutboundMessage{
			Type:           dto.MessageTypePresence,
			RoomID:         client.RoomID,
			AuthorID:       client.UserID,
			PresenceUserID: &userID,
			PresenceStatus: string(presence.Status),
			LastSeenAt:     presence.LastSeenAt,
			SentAt:         time.Now(),
		}

		select {
		case h.Outbound <- presenceMsg:
		default:
		}
	}

	joinMsg := &dto.OutboundMessage{
		Type:     dto.MessageTypeJoin,
		RoomID:   client.RoomID,
		AuthorID: client.UserID,
		SentAt:   time.Now(),
	}

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

	delete(room.Clients, client.ID)
	client.SafeClose()

	_ = h.PresenceService.StopTyping(context.Background(), client.RoomID, client.UserID)

	presence, err := h.PresenceService.MarkDisconnected(context.Background(), client.UserID, client.ID)
	if err == nil && presence != nil {
		userID := client.UserID

		presenceMsg := &dto.OutboundMessage{
			Type:           dto.MessageTypePresence,
			RoomID:         client.RoomID,
			AuthorID:       client.UserID,
			PresenceUserID: &userID,
			PresenceStatus: string(presence.Status),
			LastSeenAt:     presence.LastSeenAt,
			SentAt:         time.Now(),
		}

		select {
		case h.Outbound <- presenceMsg:
		default:
		}
	}

	leavingMsg := &dto.OutboundMessage{
		Type:     dto.MessageTypeLeave,
		RoomID:   client.RoomID,
		AuthorID: client.UserID,
		SentAt:   time.Now(),
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

// MESSAGE PROCESSING
func (h *Hub) processTextMessage(msg *dto.InboundMessage) {

	outboundingMsg, err := h.PersistFunc(context.Background(), msg)

	// handle any error occured during message creation
	if err != nil {
		errMsg := &dto.OutboundMessage{
			Type:     dto.MessageTypeError,
			RoomID:   msg.RoomID,
			AuthorID: msg.UserID,
			Error:    err.Error(),
			SentAt:   time.Now(),
		}

		// send error the the concerning user only
		h.sendToUser(msg.UserID, msg.RoomID, errMsg)
		return
	}

	// add to broadcast queue (non blocking)
	select {
	case h.Outbound <- outboundingMsg:

	default:
		//		log.Printf("Broadcast channel full, dropping message %s", outboundingMsg.)

	}

}

func (h *Hub) processTypingIndicator(msg *dto.InboundMessage) {

	_ = h.PresenceService.SetTyping(context.Background(), msg.RoomID, msg.UserID)

	typingUser := msg.UserID

	typingMsg := &dto.OutboundMessage{
		Type:       dto.MessageTypeTyping,
		RoomID:     msg.RoomID,
		AuthorID:   msg.UserID,
		SentAt:     time.Now(),
		TypingUser: &typingUser,
	}

	select {
	case h.Outbound <- typingMsg:
	default:
	}

}

func (h *Hub) processStopTypingIndicator(msg *dto.InboundMessage) {
	_ = h.PresenceService.StopTyping(context.Background(), msg.RoomID, msg.UserID)

	typingUser := msg.UserID

	stopTypingMsg := &dto.OutboundMessage{
		Type:       dto.MessageTypeStopTyping,
		RoomID:     msg.RoomID,
		AuthorID:   msg.UserID,
		SentAt:     time.Now(),
		TypingUser: &typingUser,
	}

	select {
	case h.Outbound <- stopTypingMsg:
	default:
	}
}

func (h *Hub) processReadReciept(msg *dto.InboundMessage) {
	if msg.MessageID == nil {
		errMsg := &dto.OutboundMessage{
			Type:     dto.MessageTypeError,
			RoomID:   msg.RoomID,
			AuthorID: msg.UserID,
			Error:    "message_id is required for read receipt",
			SentAt:   time.Now(),
		}
		h.sendToUser(msg.UserID, msg.RoomID, errMsg)
		return
	}

	out, err := h.ReadReceiptFunc(context.Background(), msg)
	if err != nil {
		errMsg := &dto.OutboundMessage{
			Type:     dto.MessageTypeError,
			RoomID:   msg.RoomID,
			AuthorID: msg.UserID,
			Error:    err.Error(),
			SentAt:   time.Now(),
		}
		h.sendToUser(msg.UserID, msg.RoomID, errMsg)
		return
	}

	if out == nil {
		return
	}

	select {
	case h.Outbound <- out:
	default:
		log.Printf("Could not broadcast read receipt for user %s", msg.UserID)
	}

}

// Todo: Finish this function to close hub connection
func (h *Hub) Close() error {
	return nil
}

func (h *Hub) deliverToRoom(roomId uuid.UUID, msg *dto.OutboundMessage) {

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
			client.SafeClose()
			disconnectedClients = append(disconnectedClients, userId)

		}
	}

	room.mu.RUnlock()

	// remove any disconnected clients and empty Rooms
	if len(disconnectedClients) == 0 {
		return
	}

	// If length of disconnected clients is > 0
	room.mu.Lock()

	for _, clientID := range disconnectedClients {
		delete(room.Clients, clientID)
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

func (h *Hub) sendToUser(userID uuid.UUID, roomID uuid.UUID, msg *dto.OutboundMessage) {

	// Read lock to main idemptoancy
	h.mu.RLock()
	room, exists := h.Rooms[roomID]
	h.mu.RUnlock()

	if !exists {
		return
	}

	h.mu.RLock()
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

		client.SafeClose()

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
