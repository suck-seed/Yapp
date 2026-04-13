package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/suck-seed/yapp/internal/auth"
	dto "github.com/suck-seed/yapp/internal/dto/floor"
	"github.com/suck-seed/yapp/internal/services"
	"github.com/suck-seed/yapp/internal/utils"
)

type FloorHandler struct {
	services.IFloorService
}

func NewFloorHandler(floorService services.IFloorService) *FloorHandler {
	return &FloorHandler{floorService}
}

// ── POST /halls/:hallID/floors ────────────────────────────────────────────────

func (h *FloorHandler) CreateFloor(c *gin.Context) {
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

	var req dto.CreateFloorReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.WriteError(c, utils.ErrorInvalidInput)
		return
	}

	res, err := h.IFloorService.CreateFloor(c.Request.Context(), userInfo, hallID, &req)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code":    http.StatusCreated,
		"message": "Floor created successfully",
		"success": true,
		"data":    res,
	})
}

// ── GET /halls/:hallID/floors ─────────────────────────────────────────────────

func (h *FloorHandler) GetFloors(c *gin.Context) {
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

	res, err := h.IFloorService.GetFloors(c.Request.Context(), userInfo, hallID)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Floors fetched successfully",
		"success": true,
		"data":    res,
	})
}

// ── GET /halls/:hallID/floors/:id ─────────────────────────────────────────────

func (h *FloorHandler) GetFloor(c *gin.Context) {
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

	floorID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.WriteError(c, utils.ErrorInvalidIDFormart)
		return
	}

	res, err := h.IFloorService.GetFloor(c.Request.Context(), userInfo, hallID, floorID)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Floor fetched successfully",
		"success": true,
		"data":    res,
	})
}

// ── PATCH /halls/:hallID/floors/:id ──────────────────────────────────────────

func (h *FloorHandler) UpdateFloor(c *gin.Context) {
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

	floorID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.WriteError(c, utils.ErrorInvalidIDFormart)
		return
	}

	var req dto.UpdateFloorReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.WriteError(c, utils.ErrorInvalidInput)
		return
	}

	res, err := h.IFloorService.UpdateFloor(c.Request.Context(), userInfo, hallID, floorID, &req)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Floor updated successfully",
		"success": true,
		"data":    res,
	})
}

// ── DELETE /halls/:hallID/floors/:id ─────────────────────────────────────────

func (h *FloorHandler) DeleteFloor(c *gin.Context) {
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

	floorID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.WriteError(c, utils.ErrorInvalidIDFormart)
		return
	}

	if err := h.IFloorService.DeleteFloor(c.Request.Context(), userInfo, hallID, floorID); err != nil {
		utils.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Floor deleted successfully",
		"success": true,
		"data":    nil,
	})
}

// ── PUT /halls/:hallID/floors/:id/move ────────────────────────────────────────

func (h *FloorHandler) MoveFloor(c *gin.Context) {
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

	floorID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.WriteError(c, utils.ErrorInvalidIDFormart)
		return
	}

	var req dto.MoveFloorReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.WriteError(c, utils.ErrorInvalidInput)
		return
	}

	res, err := h.IFloorService.MoveFloor(c.Request.Context(), userInfo, hallID, floorID, &req)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"success": true,
		"message": "Floor moved successfully",
		"data":    res,
	})
}
