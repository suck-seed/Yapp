package ws

import (
	"context"
	"sync"

	"github.com/google/uuid"
)

type PersistFunc func(ctx context.Context, in *InboundMessage) (*OutboundMessage, error)

type Hub struct {
	// maps roomID to roomStruct
	Rooms map[uuid.UUID]*Room

	// WebSocket channels for managing clients and messages
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *OutboundMessage
	Presist    chan *InboundMessage

	mu sync.RWMutex
}

func NewHub() Hub {
	return Hub{
		Rooms:      make(map[uuid.UUID]*Room),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *OutboundMessage, 64),
		Presist:    make(chan *InboundMessage, 256),
	}
}

func (h *Hub) Run(presist PersistFunc) {

}
