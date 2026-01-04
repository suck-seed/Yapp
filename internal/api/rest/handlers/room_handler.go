package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	dto "github.com/suck-seed/yapp/internal/dto/room"
	"github.com/suck-seed/yapp/internal/services"
	"github.com/suck-seed/yapp/internal/utils"
)

type RoomHandler struct {
	services.IRoomService
}

func NewRoomHandler(roomService services.IRoomService) *RoomHandler {
	return &RoomHandler{
		roomService,
	}
}

func (h *RoomHandler) CreateRoom(c *gin.Context) {

	u := &dto.CreateRoomReq{}

	if err := c.ShouldBindJSON(u); err != nil {
		utils.WriteError(c, utils.ErrorInvalidInput)
		return
	}

	res, err := h.IRoomService.CreateRoom(c.Request.Context(), u)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)

}
