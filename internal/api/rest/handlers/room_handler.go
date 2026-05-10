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

// CreateRoom godoc
// @Summary      Create a room
// @Description  Adds a new text or audio room to the hall. Requires ManageChannels permission.
// @Tags         rooms
// @Accept       json
// @Produce      json
// @Security     CookieAuth
// @Param        hallID  path      string             true  "Hall ID (UUID)"
// @Param        body    body      dto.CreateRoomReq  true  "Room details"
// @Success      201     {object}  map[string]interface{}
// @Failure      400     {object}  map[string]interface{}
// @Failure      401     {object}  map[string]interface{}
// @Failure      403     {object}  map[string]interface{}
// @Router       /halls/{hallID}/rooms [post]
func (h *RoomHandler) CreateRoom(c *gin.Context) {
	userInfo, err := auth.CurrentUserFromGinContext(c)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	hallID, err := uuid.Parse(c.Param("hallID"))
	if err != nil {
		utils.WriteError(c, utils.ErrorInvalidIDFormart)
		return
	}

	var req dto.CreateRoomReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.WriteError(c, utils.ErrorInvalidInput)
		return
	}

	res, err := h.IRoomService.CreateRoom(c.Request.Context(), userInfo, hallID, &req)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code": http.StatusCreated, "success": true,
		"message": "Room created successfully", "data": res,
	})
}

// GetHallRooms godoc
// @Summary      List rooms in a hall
// @Description  Returns all rooms the authenticated user can see within the hall.
// @Tags         rooms
// @Produce      json
// @Security     CookieAuth
// @Param        hallID  path      string  true  "Hall ID (UUID)"
// @Success      200     {object}  map[string]interface{}
// @Failure      400     {object}  map[string]interface{}
// @Failure      401     {object}  map[string]interface{}
// @Router       /halls/{hallID}/rooms [get]
func (h *RoomHandler) GetHallRooms(c *gin.Context) {
	userInfo, err := auth.CurrentUserFromGinContext(c)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	hallID, err := uuid.Parse(c.Param("hallID"))
	if err != nil {
		utils.WriteError(c, utils.ErrorInvalidIDFormart)
		return
	}

	res, err := h.IRoomService.GetHallRooms(c.Request.Context(), userInfo, hallID)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK, "success": true,
		"message": "Rooms fetched successfully", "data": res,
	})
}

// GetRoom godoc
// @Summary      Get a single room
// @Description  Returns details for one room.
// @Tags         rooms
// @Produce      json
// @Security     CookieAuth
// @Param        hallID  path      string  true  "Hall ID (UUID)"
// @Param        roomID  path      string  true  "Room ID (UUID)"
// @Success      200     {object}  map[string]interface{}
// @Failure      400     {object}  map[string]interface{}
// @Failure      401     {object}  map[string]interface{}
// @Failure      404     {object}  map[string]interface{}
// @Router       /halls/{hallID}/rooms/{roomID} [get]
func (h *RoomHandler) GetRoom(c *gin.Context) {
	userInfo, err := auth.CurrentUserFromGinContext(c)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	hallID, err := uuid.Parse(c.Param("hallID"))
	if err != nil {
		utils.WriteError(c, utils.ErrorInvalidIDFormart)
		return
	}

	roomID, err := uuid.Parse(c.Param("roomID"))
	if err != nil {
		utils.WriteError(c, utils.ErrorInvalidIDFormart)
		return
	}

	res, err := h.IRoomService.GetRoom(c.Request.Context(), userInfo, hallID, roomID)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK, "success": true,
		"message": "Room fetched successfully", "data": res,
	})
}

// UpdateRoom godoc
// @Summary      Update a room
// @Description  Partially updates a room's name, type, floor, or privacy. Requires ManageChannels permission.
// @Tags         rooms
// @Accept       json
// @Produce      json
// @Security     CookieAuth
// @Param        hallID  path      string             true  "Hall ID (UUID)"
// @Param        roomID  path      string             true  "Room ID (UUID)"
// @Param        body    body      dto.UpdateRoomReq  true  "Fields to update"
// @Success      200     {object}  map[string]interface{}
// @Failure      400     {object}  map[string]interface{}
// @Failure      401     {object}  map[string]interface{}
// @Failure      403     {object}  map[string]interface{}
// @Router       /halls/{hallID}/rooms/{roomID} [patch]
func (h *RoomHandler) UpdateRoom(c *gin.Context) {
	userInfo, err := auth.CurrentUserFromGinContext(c)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	hallID, err := uuid.Parse(c.Param("hallID"))
	if err != nil {
		utils.WriteError(c, utils.ErrorInvalidIDFormart)
		return
	}

	roomID, err := uuid.Parse(c.Param("roomID"))
	if err != nil {
		utils.WriteError(c, utils.ErrorInvalidIDFormart)
		return
	}

	var req dto.UpdateRoomReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.WriteError(c, utils.ErrorInvalidInput)
		return
	}

	res, err := h.IRoomService.UpdateRoom(c.Request.Context(), userInfo, hallID, roomID, &req)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK, "success": true,
		"message": "Room updated successfully", "data": res,
	})
}

// DeleteRoom godoc
// @Summary      Delete a room
// @Description  Permanently deletes a room and all its messages. Requires ManageChannels permission.
// @Tags         rooms
// @Produce      json
// @Security     CookieAuth
// @Param        hallID  path      string  true  "Hall ID (UUID)"
// @Param        roomID  path      string  true  "Room ID (UUID)"
// @Success      200     {object}  map[string]interface{}
// @Failure      400     {object}  map[string]interface{}
// @Failure      401     {object}  map[string]interface{}
// @Failure      403     {object}  map[string]interface{}
// @Router       /halls/{hallID}/rooms/{roomID} [delete]
func (h *RoomHandler) DeleteRoom(c *gin.Context) {
	userInfo, err := auth.CurrentUserFromGinContext(c)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	hallID, err := uuid.Parse(c.Param("hallID"))
	if err != nil {
		utils.WriteError(c, utils.ErrorInvalidIDFormart)
		return
	}

	roomID, err := uuid.Parse(c.Param("roomID"))
	if err != nil {
		utils.WriteError(c, utils.ErrorInvalidIDFormart)
		return
	}

	if err := h.IRoomService.DeleteRoom(c.Request.Context(), userInfo, hallID, roomID); err != nil {
		utils.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK, "success": true,
		"message": "Room deleted successfully", "data": nil,
	})
}

// MoveRoom godoc
// @Summary      Reorder a room
// @Description  Moves a room to a different position or floor. Set `after_id` to null to place it at the top.
// @Tags         rooms
// @Accept       json
// @Produce      json
// @Security     CookieAuth
// @Param        hallID  path      string           true  "Hall ID (UUID)"
// @Param        roomID  path      string           true  "Room ID (UUID)"
// @Param        body    body      dto.MoveRoomReq  true  "Target position"
// @Success      200     {object}  map[string]interface{}
// @Failure      400     {object}  map[string]interface{}
// @Failure      401     {object}  map[string]interface{}
// @Failure      403     {object}  map[string]interface{}
// @Router       /halls/{hallID}/rooms/{roomID}/move [put]
func (h *RoomHandler) MoveRoom(c *gin.Context) {
	userInfo, err := auth.CurrentUserFromGinContext(c)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	hallID, err := uuid.Parse(c.Param("hallID"))
	if err != nil {
		utils.WriteError(c, utils.ErrorInvalidIDFormart)
		return
	}

	roomID, err := uuid.Parse(c.Param("roomID"))
	if err != nil {
		utils.WriteError(c, utils.ErrorInvalidIDFormart)
		return
	}

	var req dto.MoveRoomReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.WriteError(c, utils.ErrorInvalidInput)
		return
	}

	res, err := h.IRoomService.MoveRoom(c.Request.Context(), userInfo, hallID, roomID, &req)
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

// AddRoomMember godoc
// @Summary      Add a member to a private room
// @Description  Gives a hall member access to a private room. Requires ManageChannels permission.
// @Tags         rooms
// @Produce      json
// @Security     CookieAuth
// @Param        hallID    path      string  true  "Hall ID (UUID)"
// @Param        roomID    path      string  true  "Room ID (UUID)"
// @Param        memberID  path      string  true  "Hall Member ID (UUID)"
// @Success      200       {object}  map[string]interface{}
// @Failure      400       {object}  map[string]interface{}
// @Failure      401       {object}  map[string]interface{}
// @Failure      403       {object}  map[string]interface{}
// @Router       /halls/{hallID}/rooms/{roomID}/members/{memberID} [put]
func (h *RoomHandler) AddRoomMember(c *gin.Context) {
	userInfo, err := auth.CurrentUserFromGinContext(c)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	hallID, err := uuid.Parse(c.Param("hallID"))
	if err != nil {
		utils.WriteError(c, utils.ErrorInvalidIDFormart)
		return
	}

	roomID, err := uuid.Parse(c.Param("roomID"))
	if err != nil {
		utils.WriteError(c, utils.ErrorInvalidIDFormart)
		return
	}

	memberID, err := uuid.Parse(c.Param("memberID"))
	if err != nil {
		utils.WriteError(c, utils.ErrorInvalidIDFormart)
		return
	}

	res, err := h.IRoomService.AddRoomMember(c.Request.Context(), userInfo, hallID, roomID, memberID)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"success": true,
		"message": "Room member added successfully",
		"data":    res,
	})
}

// RemoveRoomMember godoc
// @Summary      Remove a member from a private room
// @Description  Removes a hall member's access from a private room. Requires ManageChannels permission.
// @Tags         rooms
// @Produce      json
// @Security     CookieAuth
// @Param        hallID    path      string  true  "Hall ID (UUID)"
// @Param        roomID    path      string  true  "Room ID (UUID)"
// @Param        memberID  path      string  true  "Hall Member ID (UUID)"
// @Success      200       {object}  map[string]interface{}
// @Failure      400       {object}  map[string]interface{}
// @Failure      401       {object}  map[string]interface{}
// @Failure      403       {object}  map[string]interface{}
// @Router       /halls/{hallID}/rooms/{roomID}/members/{memberID} [delete]
func (h *RoomHandler) RemoveRoomMember(c *gin.Context) {
	userInfo, err := auth.CurrentUserFromGinContext(c)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	hallID, err := uuid.Parse(c.Param("hallID"))
	if err != nil {
		utils.WriteError(c, utils.ErrorInvalidIDFormart)
		return
	}

	roomID, err := uuid.Parse(c.Param("roomID"))
	if err != nil {
		utils.WriteError(c, utils.ErrorInvalidIDFormart)
		return
	}

	memberID, err := uuid.Parse(c.Param("memberID"))
	if err != nil {
		utils.WriteError(c, utils.ErrorInvalidIDFormart)
		return
	}

	res, err := h.IRoomService.RemoveRoomMember(c.Request.Context(), userInfo, hallID, roomID, memberID)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"success": true,
		"message": "Room member removed successfully",
		"data":    res,
	})
}
