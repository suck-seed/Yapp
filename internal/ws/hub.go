package ws

type Hub struct {
	// maps roomID to roomStruct
	Rooms map[string]*Room

	// WebSocket channels for managing clients and messages
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *WebsocketMessage
}

func NewHub() Hub {
	return Hub{
		Rooms: make(map[string]*Room),
	}
}
