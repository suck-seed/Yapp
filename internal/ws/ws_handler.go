package ws

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/suck-seed/yapp/config"
	"github.com/suck-seed/yapp/internal/auth"
	"github.com/suck-seed/yapp/internal/dto"
	"github.com/suck-seed/yapp/internal/models"
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

		// get origin of frontend
		origin := r.Header.Get("Origin")
		return origin == config.FrontEndOrigin()
		// return true
	},
}

// JoinRoom :
func (h *WebsocketHandler) CreateRoom(c *gin.Context) {

	req := dto.CreateRoomReq{}

	if err := c.ShouldBindJSON(&req); err != nil {

	}
}

// Join Room
// /ws/JoinRoom/:roomID
func (h *WebsocketHandler) JoinRoom(c *gin.Context) {
	// cant trust user with sending their userID, so we fetch it from context

	userIdString, _, err := auth.CurrentUserFromContext(c)
	if err != nil {
		utils.WriteError(c, utils.ErrorInvalidToken)
		return
	}

	// parse uuid
	userId, err := uuid.Parse(userIdString)
	if err != nil {
		utils.WriteError(c, utils.ErrorInvalidUserIdInContext)
		return
	}

	// verify if user exists
	user, err := h.IUserService.GetUserByID(c, &userId)
	if err != nil {
		utils.WriteError(c, utils.ErrorUserNotFound)
		return
	}

	// Parse room_id
	roomIdStr := c.Param("room_id")
	roomId, err := uuid.Parse(roomIdStr)
	if err != nil {
		utils.WriteError(c, utils.ErrorInvalidRoomId)
		return
	}

	// check if room exists
	room, err := h.IRoomService.GetRoomByID(c, &roomId)
	if err != nil {
		utils.WriteError(c, utils.ErrorRoomDoesntExist)
		return
	}

	// check if user is in the hall or not
	belongs, err := h.IHallService.IsMember(c, &room.HallId, &user.UserId)
	if err != nil || !belongs {
		utils.WriteError(c, utils.ErrorUserDoesntBelongHall)
		return
	}

	// if private, check if belongs or not
	if room.IsPrivate {

		// check on room_member table
		belongs, err := h.IRoomService.IsMember(c, &room.RoomId, &user.UserId)
		if err != nil || !belongs {
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

	// now we duplicate the room inMemory and add current cllient in that struct
	// check if it already exists or not, if it doesnt add it in room

	if _, exists := h.hub.Rooms[room.RoomId]; !exists {

		// add in the collection of rooms in hub
		h.hub.Rooms[room.RoomId] = &Room{
			RoomId:    room.RoomId,
			HallID:    room.HallId,
			FloorID:   room.FloorId,
			Name:      room.Name,
			RoomType:  RoomType(room.RoomType),
			IsPrivate: room.IsPrivate,
			CreatedAt: room.CreatedAt,
			UpdatedAt: room.UpdatedAt,

			// initialize a client list
			Clients: make(map[uuid.UUID]*Client),
		}
	}

	// register client
	client := &Client{
		Conn:        conn,
		Message:     make(chan *models.Message, 50),
		UserId:      user.UserId,
		RoomId:      room.RoomId,
		Username:    user.Username,
		DisplayName: user.DisplayName,
		AvatarURL:   user.AvatarURL,
		Description: user.Description,
		Active:      user.Active,
		CreatedAt:   user.CreatedAt,
	}

	h.hub.Register <- client

	// write message & read for message
	go client.writeMessage()
	client.readMessage(h.hub)

}

type GetClientRes struct {
	ID       string `json:"user_id"`
	Username string `json:"username"`
}

func (h *WebsocketHandler) GetClients(c *gin.Context) {

	var clients []GetClientRes
	roomIdString := c.Param("room_id")

	roomId, err := uuid.Parse(roomIdString)
	if err != nil {
		utils.WriteError(c, utils.ErrorInvalidRoomId)
	}

	if _, ok := h.hub.Rooms[roomId]; !ok {

		clients = make([]GetClientRes, 0)
		c.JSON(http.StatusOK, clients)
	}

	for _, client := range h.hub.Rooms[roomId].Clients {
		clients = append(clients, GetClientRes{
			ID:       client.UserId.String(),
			Username: client.Username,
		})
	}

	// return the clients
	c.JSON(http.StatusOK, clients)
}
