package ws

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/suck-seed/yapp/internal/auth"
	dto "github.com/suck-seed/yapp/internal/dto/message"
	"github.com/suck-seed/yapp/internal/services"
	"github.com/suck-seed/yapp/internal/utils"
)

type WebsocketHandler struct {
	hub *Hub
	services.IMessageService
	services.IHallService
	services.IRoomService
	services.IUserService
}

func NewWebsocketHandler(h *Hub, messageService services.IMessageService, hallService services.IHallService, roomService services.IRoomService, userService services.IUserService) *WebsocketHandler {
	return &WebsocketHandler{
		h,
		messageService,
		hallService,
		roomService,
		userService,
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {

		// Restrict by Origin here if needed; browser CORS for HTTP is set in Nginx.
		return true
	},
}

// Connect godoc
// @Summary      Join a room via WebSocket
// @Description  Upgrades the HTTP connection to WebSocket and subscribes the caller to the room's message stream.
//
//	The client must send **inbound** messages as JSON matching `dto.InboundMessage` and will
//	receive **outbound** messages matching `dto.OutboundMessage`.
//
//	**Inbound message shape**
//	```json
//	{ "type": "send" | "edit" | "delete" | "typing", "content": "...", ... }
//	```
//
//	**Outbound message shape**
//	```json
//	{ "type": "send", "id": "uuid", "room_id": "uuid", "author_id": "uuid", "content": "...", "sent_at": "..." }
//	```
//
//	**Connection lifecycle**: the server sends a WebSocket ping every ~54 s and expects
//	a pong reply within 60 s, otherwise the connection is closed.
//
// @Tags         websocket
// @Produce      json
// @Security     CookieAuth
// @Param        room_id  path  string  true  "Room ID (UUID)"
// @Success      101      "Switching Protocols — WebSocket handshake successful"
// @Failure      400      {object}  map[string]interface{}  "Bad room / hall ID"
// @Failure      401      {object}  map[string]interface{}  "Not authenticated"
// @Failure      403      {object}  map[string]interface{}  "Not a hall/room member"
// @Failure      404      {object}  map[string]interface{}  "Room or hall not found"
// @Router       /ws/rooms/{room_id} [get]
func (h *WebsocketHandler) Connect(c *gin.Context) {
	// cant trust user with sending their userID, so we fetch it from context

	userInfo, err := auth.CurrentUserFromGinContext(c)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	// User Exists?
	user, err := h.IUserService.GetUserById(c.Request.Context(), userInfo.ID)
	if err != nil {
		utils.WriteError(c, utils.ErrorUserNotFound)
		return
	}

	// room_id -> hall_id
	// This is the core of your new specification: subscribe the app/device websocket to all rooms.
	subscribedRooms, err := h.IRoomService.GetAccessibleRoomsForUser(c.Request.Context(), &auth.UserInfo{ID: user.ID})
	if err != nil {
		utils.WriteError(c, utils.ErrorConnectingWebsocket)
		return
	}

	// Upgrade HTTP to websocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		utils.WriteError(c, utils.ErrorFailedUpgrade)
		return
	}

	// register client

	clientID, err := uuid.NewUUID()
	if err != nil {
		// close the connection already
		_ = conn.Close()
		utils.WriteError(c, utils.ErrorFailedUpgrade)
		return
	}

	const sendBuf = 1024
	client := &Client{
		ID:     clientID,
		Conn:   conn,
		Send:   make(chan *dto.OutboundMessage, sendBuf),
		UserID: user.ID,

		// INCLUDE THISSS ASAP
		SubscribedRooms: subscribedRooms,

		ConnectedAt: time.Now(),
		LastPing:    time.Now(),
	}

	log.Printf("Client Information: \n %v\n", client)

	h.hub.Register <- client

	// write message & read for message (new thread to stop blocking)
	go client.writePump()
	go client.readPump(h.hub)

	//handler can return if any error passed from writePump / readPump
	// return

}

type GetClientRes struct {
	ID string `json:"id"`
}

// GetClients godoc
// @Summary      List connected clients in a room
// @Description  Returns unique user IDs currently connected/subscribed to the room. Useful for @-mention autocomplete or debugging.
// @Tags         websocket
// @Produce      json
// @Security     CookieAuth
// @Param        room_id  path      string  true  "Room ID (UUID)"
// @Success      200      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}
// @Router       /ws/clients/{room_id} [get]
func (h *WebsocketHandler) GetClients(c *gin.Context) {
	roomID, err := uuid.Parse(c.Param("room_id"))
	if err != nil {
		utils.WriteError(c, utils.ErrorInvalidRoomIDFormat)
		return
	}

	clients := make([]GetClientRes, 0)
	seenUsers := make(map[uuid.UUID]bool)

	h.hub.mu.RLock()
	room, exists := h.hub.Rooms[roomID]
	if exists {
		for _, client := range room.Clients {
			if seenUsers[client.UserID] {
				continue
			}
			seenUsers[client.UserID] = true
			clients = append(clients, GetClientRes{ID: client.UserID.String()})
		}
	}
	h.hub.mu.RUnlock()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Room clients retrieved successfully",
		"data":    clients,
	})
}
