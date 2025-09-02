package ws

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/suck-seed/yapp/internal/dto"
	"github.com/suck-seed/yapp/internal/services"
	"github.com/suck-seed/yapp/internal/utils"
)

type WsHandler struct {
	hub *Hub
	services.IMessageService
	services.IRoomService
	services.IUserService
}

func NewWsHandler(h *Hub, messageService services.IMessageService, roomService services.IRoomService, userService services.IUserService) *WsHandler {
	return &WsHandler{
		h,
		messageService,
		roomService,
		userService,
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {

		// get origin of frontend
		// origin := r.Header.Get("Origin")
		// return origin == config.FrontEndOrigin()
		return true
	},
}

// JoinRoom :
func (h *WsHandler) CreateRoom(c *gin.Context) {

	req := dto.CreateRoomReq{}

	if err := c.ShouldBindJSON(&req); err != nil {

	}
}

// Join Room
func (h *WsHandler) JoinRoom(c *gin.Context) {

	// Upgrade HTTP to websocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		utils.WriteError(c, utils.ErrorFailedUpgrade)
		return
	}

	// Parse parameters
	// As we will be sending /room_id to join the room
	roomIdStr := c.Query("room-id")
	rooomId, err := uuid.Parse(roomIdStr)
	if err != nil {
		utils.WriteError(c, utils.ErrorInvalidRoomId)

	}

	// cant trust user with sending their userID, so we fetch it from context

}
