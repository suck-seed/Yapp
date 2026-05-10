package ws

import (
	"context"
	"log"

	"sync"
	"time"

	"github.com/google/uuid"
	dto "github.com/suck-seed/yapp/internal/dto/message"
	"github.com/suck-seed/yapp/internal/services"
	"github.com/suck-seed/yapp/internal/utils"
)

type Hub struct {
	// room_id -> room subscription bucket
	Rooms map[uuid.UUID]*Room

	// client_id -> client
	Clients map[uuid.UUID]*Client

	// user_id -> client_id -> client
	UserClients map[uuid.UUID]map[uuid.UUID]*Client

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

	// One lock protects Rooms, Clients, and UserClients.
	mu sync.RWMutex
}

// Creates a New Presistent working Hub
func NewHub(
	p PersistFunction,
	readFunc ReadReceiptFunction,
	presenceService services.IPresenceService,
) Hub {
	return Hub{
		Rooms:           make(map[uuid.UUID]*Room),   // room_id -> room subscription bucket
		Clients:         make(map[uuid.UUID]*Client), // client_id -> client
		UserClients:     make(map[uuid.UUID]map[uuid.UUID]*Client),
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
			log.Printf("Inbound Message: handleInboundMessage() \n %v\n", inboundMessage)
			h.processTextMessage(inboundMessage)
		case dto.MessageTypeTyping:
			h.processTypingIndicator(inboundMessage)
		case dto.MessageTypeStopTyping:
			h.processStopTypingIndicator(inboundMessage)
		case dto.MessageTypeRead:
			h.processReadReciept(inboundMessage)
		default:
			log.Printf("unknown websocket message type %s", inboundMessage.Type)
			h.sendErrorToClient(inboundMessage.ClientID, inboundMessage.RoomID, inboundMessage.UserID, "unknown websocket message type")
		}

	}

}

func (h *Hub) handleOutbound() {

	for msg := range h.Outbound {
		log.Printf("Outbound Message handleOutbound() : \n %v\n", msg)
		h.deliverToRoom(msg.RoomID, msg)
	}
}

func (h *Hub) registerClient(client *Client) {

	//  just in case if Client's Subscribe Room isnt made
	if client.SubscribedRooms == nil {
		client.SubscribedRooms = make(map[uuid.UUID]uuid.UUID)
	}

	h.mu.Lock()

	h.Clients[client.ID] = client

	// If client with this UserID doesnt exist
	if _, exists := h.UserClients[client.UserID]; !exists {

		// Make a new UserClient
		// With this new UserID as client.UserID
		h.UserClients[client.UserID] = make(map[uuid.UUID]*Client)
	}

	// Then map the current client as one of the client of this user
	h.UserClients[client.UserID][client.ID] = client

	// userID ->->->-> ClientID1 ->->->-> Client1 (Chrome)
	// 			\ ->->->-> ClientID2 ->->->-> Client2 (Firefox)

	// Now create a suscribtion
	for roomID, hallID := range client.SubscribedRooms {

		// Check if the subscribed room of client
		// exist or not in the hub.Rooms
		room, exists := h.Rooms[roomID]

		// if doesnt exist, create the room and add in h.Rooms
		if !exists {

			room = &Room{
				ID:      roomID,
				HallID:  hallID,
				Clients: make(map[uuid.UUID]*Client),
			}

			h.Rooms[roomID] = room
		}

		// Now assign this client in room;s client
		room.Clients[client.ID] = client

	}
	h.mu.Unlock()

	if h.PresenceService != nil {

		presence, err := h.PresenceService.MarkConnected(context.Background(), client.UserID, client.ID)

		if err == nil && presence != nil {
			h.broadcastPresenceToRooms(client.SubscribedRooms, client.UserID, string(presence.Status), presence.LastSeenAt)
		}
	}

	log.Printf("Client Registered: \n %v\n", client)

}

func (h *Hub) unregisterClient(client *Client) {

	roomIDs := client.RoomIDs()
	roomsForPresence := cloneRoomMap(client.SubscribedRooms)

	h.mu.Lock()

	current, exists := h.Clients[client.ID]
	if exists {
		h.removeClientLocked(current)
	}
	// unlock the client
	h.mu.Unlock()

	client.SafeClose()

	if h.PresenceService != nil {
		for _, roomID := range roomIDs {
			_ = h.PresenceService.StopTyping(context.Background(), roomID, client.UserID)
		}

		presence, err := h.PresenceService.MarkDisconnected(context.Background(), client.UserID, client.ID)
		if err == nil && presence != nil {
			h.broadcastPresenceToRooms(roomsForPresence, client.UserID, string(presence.Status), presence.LastSeenAt)
		}
	}

	log.Printf("Client Unregistered: \n %v\n", client)

}

// removeClientLocked must only be called while h.mu is write-locked.
func (h *Hub) removeClientLocked(client *Client) {
	delete(h.Clients, client.ID)

	if clientsByUser, exists := h.UserClients[client.UserID]; exists {
		delete(clientsByUser, client.ID)
		if len(clientsByUser) == 0 {
			delete(h.UserClients, client.UserID)
		}
	}

	for roomID := range client.SubscribedRooms {
		room, exists := h.Rooms[roomID]
		if !exists {
			continue
		}

		delete(room.Clients, client.ID)
		if len(room.Clients) == 0 {
			delete(h.Rooms, roomID)
		}
	}
}

func cloneRoomMap(in map[uuid.UUID]uuid.UUID) map[uuid.UUID]uuid.UUID {
	out := make(map[uuid.UUID]uuid.UUID, len(in))
	for roomID, hallID := range in {
		out[roomID] = hallID
	}
	return out
}

// MESSAGE PROCESSING
func (h *Hub) processTextMessage(msg *dto.InboundMessage) {

	if msg.RoomID == uuid.Nil {
		h.sendErrorToClient(msg.ClientID, uuid.Nil, msg.UserID, "room_id is required")
		return
	}

	// is client even subscribed into this room
	if !h.isClientSubscribedToRoom(msg.ClientID, msg.RoomID) {
		h.sendErrorToClient(msg.ClientID, msg.RoomID, msg.UserID, "you are not subscribed to this room")
		return
	}

	outboundingMsg, err := h.PersistFunc(context.Background(), msg)
	if err != nil {
		h.sendErrorToClient(msg.ClientID, msg.RoomID, msg.UserID, err.Error())
		return
	}

	select {
	case h.Outbound <- outboundingMsg:
	default:
		log.Printf("outbound channel full, dropping message %s", outboundingMsg.ID)
	}

}

func (h *Hub) processTypingIndicator(msg *dto.InboundMessage) {

	if !h.isClientSubscribedToRoom(msg.ClientID, msg.RoomID) {
		h.sendErrorToClient(msg.ClientID, msg.RoomID, msg.UserID, "you are not subscribed to this room")
		return
	}

	if h.PresenceService != nil {
		_ = h.PresenceService.SetTyping(context.Background(), msg.RoomID, msg.UserID)
	}

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

	if !h.isClientSubscribedToRoom(msg.ClientID, msg.RoomID) {
		h.sendErrorToClient(msg.ClientID, msg.RoomID, msg.UserID, "you are not subscribed to this room")
		return
	}

	if h.PresenceService != nil {
		_ = h.PresenceService.StopTyping(context.Background(), msg.RoomID, msg.UserID)
	}

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

	if !h.isClientSubscribedToRoom(msg.ClientID, msg.RoomID) {
		h.sendErrorToClient(msg.ClientID, msg.RoomID, msg.UserID, "you are not subscribed to this room")
		return
	}

	out, err := h.ReadReceiptFunc(context.Background(), msg)
	if err != nil {
		h.sendErrorToClient(msg.ClientID, msg.RoomID, msg.UserID, err.Error())
		return
	}
	if out == nil {
		return
	}

	select {
	case h.Outbound <- out:
	default:
		log.Printf("could not broadcast read receipt for user %s", msg.UserID)
	}

}

func (h *Hub) deliverToRoom(roomID uuid.UUID, msg *dto.OutboundMessage) {
	var disconnectedClients []uuid.UUID

	log.Print(" Hit: deliverToRoom")
	h.mu.RLock()
	room, exists := h.Rooms[roomID]
	if !exists {
		h.mu.RUnlock()
		return
	}

	if msg.HallID == uuid.Nil && room.HallID != uuid.Nil {
		msg.HallID = room.HallID
	}

	for clientID, client := range room.Clients {
		select {
		case client.Send <- msg:
		default:
			log.Printf("client %s buffer full, disconnecting", clientID)
			disconnectedClients = append(disconnectedClients, clientID)
		}
	}
	h.mu.RUnlock()

	if len(disconnectedClients) == 0 {
		return
	}

	h.mu.Lock()
	for _, clientID := range disconnectedClients {
		if client, exists := h.Clients[clientID]; exists {
			h.removeClientLocked(client)
			client.SafeClose()
		}
	}
	h.mu.Unlock()

}

func (h *Hub) sendErrorToClient(clientID uuid.UUID, roomID uuid.UUID, userID uuid.UUID, message string) {
	errMsg := &dto.OutboundMessage{
		Type:     dto.MessageTypeError,
		RoomID:   roomID,
		AuthorID: userID,
		Error:    utils.StringToPointer(message),
		SentAt:   time.Now(),
	}
	h.sendToClientID(clientID, errMsg)
}

func (h *Hub) sendToClientID(clientID uuid.UUID, msg *dto.OutboundMessage) {
	var disconnected *Client

	h.mu.RLock()
	client, exists := h.Clients[clientID]
	if exists {
		select {
		case client.Send <- msg:
		default:
			disconnected = client
		}
	}
	h.mu.RUnlock()

	if disconnected == nil {
		return
	}

	h.mu.Lock()
	if current, exists := h.Clients[clientID]; exists {
		h.removeClientLocked(current)
		current.SafeClose()
	}
	h.mu.Unlock()
}

// sendToUser sends to all currently connected browser tabs/devices for the user.
// Useful later for friend requests, direct notifications, etc.
func (h *Hub) sendToUser(userID uuid.UUID, msg *dto.OutboundMessage) {
	var disconnectedClients []uuid.UUID

	h.mu.RLock()
	clientsByUser := h.UserClients[userID]
	for clientID, client := range clientsByUser {
		select {
		case client.Send <- msg:
		default:
			disconnectedClients = append(disconnectedClients, clientID)
		}
	}
	h.mu.RUnlock()

	if len(disconnectedClients) == 0 {
		return
	}

	h.mu.Lock()
	for _, clientID := range disconnectedClients {
		if client, exists := h.Clients[clientID]; exists {
			h.removeClientLocked(client)
			client.SafeClose()
		}
	}
	h.mu.Unlock()
}

func (h *Hub) Close() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	for _, client := range h.Clients {
		client.SafeClose()
		_ = client.Conn.Close()
	}

	close(h.Inbound)
	close(h.Outbound)
	return nil
}

// ->->->->->->->->->- OTHER HELPER FUNCTIONS

func (h *Hub) isClientSubscribedToRoom(clientID uuid.UUID, roomID uuid.UUID) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	client, exists := h.Clients[clientID]
	if !exists {
		return false
	}
	return client.IsSubscribedToRoom(roomID)
}

func (h *Hub) broadcastPresenceToRooms(rooms map[uuid.UUID]uuid.UUID, userID uuid.UUID, status string, lastSeenAt *time.Time) {
	for roomID, hallID := range rooms {
		uid := userID
		hid := hallID
		msg := &dto.OutboundMessage{
			Type:           dto.MessageTypePresence,
			RoomID:         roomID,
			HallID:         hid,
			AuthorID:       userID,
			PresenceUserID: &uid,
			PresenceStatus: utils.StringToPointer(status),
			LastSeenAt:     lastSeenAt,
			SentAt:         time.Now(),
		}
		h.deliverToRoom(roomID, msg)
	}
}
