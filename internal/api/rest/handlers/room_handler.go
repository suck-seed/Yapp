package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/suck-seed/yapp/internal/dto"
	"github.com/suck-seed/yapp/internal/services"
	"github.com/suck-seed/yapp/internal/utils"
	"net/http"
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
