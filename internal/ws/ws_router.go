package ws

import (
	"github.com/gin-gonic/gin"
	"github.com/suck-seed/yapp/internal/services"
)

func RegisterWebSocketRoutes(r *gin.RouterGroup, hub *Hub, messageService services.IMessageService, hallService services.IHallService, roomService services.IRoomService, userService services.IUserService) {
	// inject dependency to wsHandler
	wsHandler := NewWebsocketHandler(hub, messageService, hallService, roomService, userService)

	// Single gateway socket per browser/app/device.
	r.GET("/", wsHandler.Connect)

	// Used for debugging
	// TODO : remove before pushing
	r.GET("/clients/:room_id")

}
