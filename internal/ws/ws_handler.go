package ws

import (
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

// JoinRoom godoc
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
func (h *WebsocketHandler) JoinRoom(c *gin.Context) {
	// cant trust user with sending their userID, so we fetch it from context

	userInfo, err := auth.CurrentUserFromGinContext(c)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	// User Exists?
	user, err := h.IUserService.GetUserById(c, userInfo.ID)
	if err != nil {
		utils.WriteError(c, utils.ErrorUserNotFound)
		return
	}

	// Parse room_id
	roomIDStr := c.Param("room_id")
	roomID, err := uuid.Parse(roomIDStr)
	if err != nil {
		utils.WriteError(c, utils.ErrorInvalidRoomIDFormat)
		return
	}

	// Room exists? && Fetch Room
	room, err := h.IRoomService.GetRoomByID(c, roomID)
	if err != nil {
		utils.WriteError(c, utils.ErrorRoomNotFound)
		return
	}

	// Hall exists?
	hallExists, err := h.IHallService.DoesHallExist(c, room.HallID)
	if err != nil {
		utils.WriteError(c, err)
		return
	}
	if !hallExists {
		utils.WriteError(c, utils.ErrorHallNotFound)
		return
	}

	// Hall Member ?
	belongs, err := h.IHallService.IsUserHallMember(c, room.HallID, user.ID)
	if err != nil {
		utils.WriteError(c, err)
		return
	}
	if !belongs {
		utils.WriteError(c, utils.ErrorUserDoesntBelongHall)
		return
	}

	// Room Member ? ( ON PRIVATE ROOMS )
	if room.IsPrivate {

		// check on room_member table
		belongs, err := h.IRoomService.IsUserRoomMember(c, room.ID, user.ID)
		if err != nil {
			utils.WriteError(c, err)
			return
		}
		if !belongs {
			utils.WriteError(c, utils.ErrorUserDoesntBelongRoom)
			return
		}

	}

	// get all the info and validate themmm, ani ball upgrading them

	// Upgrade HTTP to websocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		utils.WriteError(c, utils.ErrorFailedUpgrade)
		return
	}

	// Do this during registiring not here, race condition aauxa
	// if _, exists := h.hub.Rooms[room.RoomID]; !exists {

	// 	// add in the collection of rooms in hub
	// 	h.hub.Rooms[room.RoomID] = &Room{
	// 		RoomID:    room.RoomID,
	// 		HallID:    room.HallId,
	// 		FloorID:   room.FloorId,
	// 		Name:      room.Name,
	// 		RoomType:  RoomType(room.RoomType),
	// 		IsPrivate: room.IsPrivate,
	// 		CreatedAt: room.CreatedAt,
	// 		UpdatedAt: room.UpdatedAt,

	// 		// initialize a client list
	// 		Clients: make(map[uuid.UUID]*Client),
	// 	}
	// }

	// register client
	const sendBuf = 1024

	client := &Client{
		Conn:        conn,
		Send:        make(chan *dto.OutboundMessage, sendBuf),
		UserID:      user.ID,
		RoomID:      room.ID,
		ConnectedAt: time.Now(),
	}

	h.hub.Register <- client

	// write message & read for message (new thread to stop blocking)
	go client.writePump()
	go client.readPump(h.hub)

	//handler can return if any error passed from writePump / readPump
	// return

}

type GetClientRes struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

// GetClients godoc
// @Summary      List connected clients in a room
// @Description  Returns the user IDs of every WebSocket client currently connected to the room. Useful for @-mention autocomplete.
// @Tags         websocket
// @Produce      json
// @Security     CookieAuth
// @Param        room_id  path      string  true  "Room ID (UUID)"
// @Success      200      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}
// @Router       /ws/clients/{room_id} [get]
func (h *WebsocketHandler) GetClients(c *gin.Context) {

	var clients []GetClientRes
	roomIdString := c.Param("room_id")

	roomId, err := uuid.Parse(roomIdString)
	if err != nil {
		utils.WriteError(c, utils.ErrorInvalidRoomIDFormat)
		return
	}

	if _, ok := h.hub.Rooms[roomId]; !ok {
		clients = make([]GetClientRes, 0)
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Room clients retrieved successfully",
			"data":    clients,
		})
		return
	}

	for _, client := range h.hub.Rooms[roomId].Clients {
		clients = append(clients, GetClientRes{
			ID: client.UserID.String(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Room clients retrieved successfully",
		"data":    clients,
	})
}
