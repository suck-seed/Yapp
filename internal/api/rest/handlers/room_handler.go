package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/suck-seed/yapp/internal/auth"
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
	log.Printf(">>> CREATE ROOM HANDLER HIT")
	u := &dto.CreateRoomReq{}

	if err := c.ShouldBindJSON(u); err != nil {
		log.Printf("BIND ERROR: %v", err)
		utils.WriteError(c, utils.ErrorInvalidInput)
		return
	}
	log.Printf(">>> AFTER BINDING")

	userInfo, err := auth.CurrentUserFromGinContext(c)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	res, err := h.IRoomService.CreateRoom(c.Request.Context(), userInfo, u)
	if err != nil {
		log.Printf(">>> SERVICE ERROR: %v", err) // add this
		utils.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Room created successfully",
		"data":    res,
	})
}
