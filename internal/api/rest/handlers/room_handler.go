package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/suck-seed/yapp/internal/auth"
	dto "github.com/suck-seed/yapp/internal/dto/room"
	"github.com/suck-seed/yapp/internal/services"
	"github.com/suck-seed/yapp/internal/utils"
)

type RoomHandler struct {
	services.IRoomService
}

func NewRoomHandler(roomService services.IRoomService) *RoomHandler {
	return &RoomHandler{roomService}
}

// ── POST /rooms ───────────────────────────────────────────────────────────────

func (h *RoomHandler) CreateRoom(c *gin.Context) {
	var req dto.CreateRoomReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.WriteError(c, utils.ErrorInvalidInput)
		return
	}

	userInfo, err := auth.CurrentUserFromGinContext(c)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	res, err := h.IRoomService.CreateRoom(c.Request.Context(), userInfo, &req)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code": http.StatusCreated, "success": true,
		"message": "Room created successfully", "data": res,
	})
}

// ── GET /rooms/:id ────────────────────────────────────────────────────────────

func (h *RoomHandler) GetRoom(c *gin.Context) {
	roomID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.WriteError(c, utils.ErrorInvalidIDFormart)
		return
	}

	res, err := h.IRoomService.GetRoom(c.Request.Context(), roomID)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK, "success": true,
		"message": "Room fetched successfully", "data": res,
	})
}

// ── GET /rooms?hall_id= ───────────────────────────────────────────────────────

func (h *RoomHandler) GetHallRooms(c *gin.Context) {
	var req dto.GetHallRoomsReq
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.WriteError(c, utils.ErrorInvalidInput)
		return
	}

	res, err := h.IRoomService.GetHallRooms(c.Request.Context(), req.HallID)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK, "success": true,
		"message": "Rooms fetched successfully", "data": res,
	})
}

// ── PATCH /rooms/:id ──────────────────────────────────────────────────────────

func (h *RoomHandler) UpdateRoom(c *gin.Context) {
	roomID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.WriteError(c, utils.ErrorInvalidIDFormart)
		return
	}

	var req dto.UpdateRoomReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.WriteError(c, utils.ErrorInvalidInput)
		return
	}

	res, err := h.IRoomService.UpdateRoom(c.Request.Context(), roomID, &req)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK, "success": true,
		"message": "Room updated successfully", "data": res,
	})
}

// ── DELETE /rooms/:id ─────────────────────────────────────────────────────────

func (h *RoomHandler) DeleteRoom(c *gin.Context) {
	roomID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.WriteError(c, utils.ErrorInvalidIDFormart)
		return
	}

	if err := h.IRoomService.DeleteRoom(c.Request.Context(), roomID); err != nil {
		utils.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK, "success": true,
		"message": "Room deleted successfully", "data": nil,
	})
}

// ── PUT /rooms/:id/move ───────────────────────────────────────────────────────

func (h *RoomHandler) MoveRoom(c *gin.Context) {
	roomID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.WriteError(c, utils.ErrorInvalidIDFormart)
		return
	}

	var req dto.MoveRoomReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.WriteError(c, utils.ErrorInvalidInput)
		return
	}

	res, err := h.IRoomService.MoveRoom(c.Request.Context(), roomID, &req)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"success": true,
		"message": "Room moved successfully",
		"data":    res,
	})
}
