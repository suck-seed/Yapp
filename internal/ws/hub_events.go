package ws

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/suck-seed/yapp/internal/realtime"
)

func (h *Hub) handleEvents() {
	if h.EventBus == nil {
		return
	}

	for event := range h.EventBus.Events {
		h.handleHubEvent(event)
	}
}

func (h *Hub) handleHubEvent(event realtime.HubEvent) {
	switch event.Type {

	case realtime.HubEventUserJoinedHall:
		h.resyncUserAccess(context.Background(), event.UserID)

	case realtime.HubEventUserLeftHall,
		realtime.HubEventUserKickedFromHall,
		realtime.HubEventUserBannedFromHall:
		h.unsubscribeUserFromHall(event.UserID, event.HallID)

	case realtime.HubEventHallDeleted:
		h.unsubscribeAllClientsFromHall(event.HallID)

	case realtime.HubEventRoomCreated:
		if event.RoomID == uuid.Nil || event.HallID == uuid.Nil {
			return
		}

		if event.IsPrivate {
			// Private room access depends on room_members/floor sync.
			// Safest: resync creator/user if provided.
			if event.UserID != uuid.Nil {
				h.resyncUserAccess(context.Background(), event.UserID)
			}
			return
		}

		// Public room: every currently connected client already subscribed
		// to any room in this hall should get this room too.
		h.subscribeHallClientsToRoom(event.HallID, event.RoomID)

	case realtime.HubEventRoomDeleted:
		h.unsubscribeAllClientsFromRoom(event.RoomID)

	case realtime.HubEventRoomPrivacyChanged,
		realtime.HubEventRoomMoved:
		h.resyncHallAccess(context.Background(), event.HallID)

	case realtime.HubEventRoomMemberAdded:
		if event.UserID != uuid.Nil && event.RoomID != uuid.Nil && event.HallID != uuid.Nil {
			h.subscribeUserClientsToRoom(event.UserID, event.HallID, event.RoomID)
		}

	case realtime.HubEventRoomMemberRemoved:
		if event.UserID != uuid.Nil && event.RoomID != uuid.Nil {
			h.unsubscribeUserClientsFromRoom(event.UserID, event.RoomID)
		}

	case realtime.HubEventFloorMemberAdded,
		realtime.HubEventFloorMemberRemoved:
		// Floor member changes can sync many rooms in the floor.
		// Resync only that user if we know them.
		if event.UserID != uuid.Nil {
			h.resyncUserAccess(context.Background(), event.UserID)
		} else {
			h.resyncHallAccess(context.Background(), event.HallID)
		}

	case realtime.HubEventFloorPrivacyChanged,
		realtime.HubEventFloorDeleted,
		realtime.HubEventHallAccessResync:
		h.resyncHallAccess(context.Background(), event.HallID)

	case realtime.HubEventUserAccessResync:
		h.resyncUserAccess(context.Background(), event.UserID)

	default:
		log.Printf("unknown hub event type: %+v", event)
	}
}

// Subscription Mutation Helpers
func (h *Hub) subscribeHallClientsToRoom(hallID uuid.UUID, roomID uuid.UUID) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for _, client := range h.Clients {
		if !clientHasHall(client, hallID) {
			continue
		}

		h.subscribeClientToRoomLocked(client, hallID, roomID)
	}
}

func (h *Hub) subscribeUserClientsToRoom(userID uuid.UUID, hallID uuid.UUID, roomID uuid.UUID) {
	h.mu.Lock()
	defer h.mu.Unlock()

	clientsByUser := h.UserClients[userID]
	for _, client := range clientsByUser {
		h.subscribeClientToRoomLocked(client, hallID, roomID)
	}
}

func (h *Hub) unsubscribeUserClientsFromRoom(userID uuid.UUID, roomID uuid.UUID) {
	h.mu.Lock()
	defer h.mu.Unlock()

	clientsByUser := h.UserClients[userID]
	for _, client := range clientsByUser {
		h.unsubscribeClientFromRoomLocked(client, roomID)
	}
}

func (h *Hub) unsubscribeAllClientsFromRoom(roomID uuid.UUID) {
	h.mu.Lock()
	defer h.mu.Unlock()

	room, exists := h.Rooms[roomID]
	if exists {
		for _, client := range room.Clients {
			client.mu.Lock()
			delete(client.SubscribedRooms, roomID)
			client.mu.Unlock()
		}
	}

	delete(h.Rooms, roomID)
}

func (h *Hub) unsubscribeUserFromHall(userID uuid.UUID, hallID uuid.UUID) {
	h.mu.Lock()
	defer h.mu.Unlock()

	clientsByUser := h.UserClients[userID]
	for _, client := range clientsByUser {
		removeRooms := make([]uuid.UUID, 0)

		client.mu.Lock()
		for roomID, currentHallID := range client.SubscribedRooms {
			if currentHallID == hallID {
				removeRooms = append(removeRooms, roomID)
			}
		}
		client.mu.Unlock()

		for _, roomID := range removeRooms {
			h.unsubscribeClientFromRoomLocked(client, roomID)
		}
	}
}

func (h *Hub) unsubscribeAllClientsFromHall(hallID uuid.UUID) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for _, client := range h.Clients {
		removeRooms := make([]uuid.UUID, 0)

		client.mu.Lock()
		for roomID, currentHallID := range client.SubscribedRooms {
			if currentHallID == hallID {
				removeRooms = append(removeRooms, roomID)
			}
		}
		client.mu.Unlock()

		for _, roomID := range removeRooms {
			h.unsubscribeClientFromRoomLocked(client, roomID)
		}
	}
}

func (h *Hub) resyncUserAccess(ctx context.Context, userID uuid.UUID) {
	if h.AccessResolver == nil || userID == uuid.Nil {
		return
	}

	newRooms, err := h.AccessResolver(ctx, userID)
	if err != nil {
		log.Printf("failed to resync user access %s: %v", userID, err)
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	clientsByUser := h.UserClients[userID]
	for _, client := range clientsByUser {
		oldRooms := client.SubscribedRoomsSnapshot()

		for oldRoomID := range oldRooms {
			if _, stillAllowed := newRooms[oldRoomID]; !stillAllowed {
				h.unsubscribeClientFromRoomLocked(client, oldRoomID)
			}
		}

		for newRoomID, hallID := range newRooms {
			h.subscribeClientToRoomLocked(client, hallID, newRoomID)
		}
	}
}

func (h *Hub) resyncHallAccess(ctx context.Context, hallID uuid.UUID) {
	if hallID == uuid.Nil {
		return
	}

	userIDs := h.connectedUserIDsInHall(hallID)

	for _, userID := range userIDs {
		h.resyncUserAccess(ctx, userID)
	}
}

func (h *Hub) connectedUserIDsInHall(hallID uuid.UUID) []uuid.UUID {
	seen := make(map[uuid.UUID]bool)
	out := make([]uuid.UUID, 0)

	h.mu.RLock()
	defer h.mu.RUnlock()

	for userID, clientsByUser := range h.UserClients {
		for _, client := range clientsByUser {
			if clientHasHall(client, hallID) {
				if !seen[userID] {
					seen[userID] = true
					out = append(out, userID)
				}
				break
			}
		}
	}

	return out
}

// h.mu must already be write-locked.
func (h *Hub) subscribeClientToRoomLocked(client *Client, hallID uuid.UUID, roomID uuid.UUID) {
	if client == nil || roomID == uuid.Nil || hallID == uuid.Nil {
		return
	}

	client.mu.Lock()
	if client.SubscribedRooms == nil {
		client.SubscribedRooms = make(map[uuid.UUID]uuid.UUID)
	}
	client.SubscribedRooms[roomID] = hallID
	client.mu.Unlock()

	room, exists := h.Rooms[roomID]
	if !exists {
		room = &Room{
			ID:      roomID,
			HallID:  hallID,
			Clients: make(map[uuid.UUID]*Client),
		}
		h.Rooms[roomID] = room
	}

	room.HallID = hallID
	room.Clients[client.ID] = client
}

// h.mu must already be write-locked.
func (h *Hub) unsubscribeClientFromRoomLocked(client *Client, roomID uuid.UUID) {
	if client == nil || roomID == uuid.Nil {
		return
	}

	client.mu.Lock()
	delete(client.SubscribedRooms, roomID)
	client.mu.Unlock()

	room, exists := h.Rooms[roomID]
	if !exists {
		return
	}

	delete(room.Clients, client.ID)

	if len(room.Clients) == 0 {
		delete(h.Rooms, roomID)
	}
}

func clientHasHall(client *Client, hallID uuid.UUID) bool {
	if client == nil || hallID == uuid.Nil {
		return false
	}

	client.mu.Lock()
	defer client.mu.Unlock()

	for _, currentHallID := range client.SubscribedRooms {
		if currentHallID == hallID {
			return true
		}
	}

	return false
}
