package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

// ── POST /floors ──────────────────────────────────────────────────────────────

func (h *FloorHandler) CreateFloor(c *gin.Context) {
	var req dto.CreateFloorReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.WriteError(c, utils.ErrorInvalidInput)
		return
	}

	res, err := h.IFloorService.CreateFloor(c.Request.Context(), &req)
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

// ── GET /floors/:id ───────────────────────────────────────────────────────────

func (h *FloorHandler) GetFloor(c *gin.Context) {
	floorID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.WriteError(c, utils.ErrorInvalidIDFormart)
		return
	}

	res, err := h.IFloorService.GetFloor(c.Request.Context(), floorID)
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

// ── GET /floors?hall_id= ──────────────────────────────────────────────────────

func (h *FloorHandler) GetFloors(c *gin.Context) {
	var req dto.GetFloorsReq
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.WriteError(c, utils.ErrorInvalidInput)
		return
	}

	res, err := h.IFloorService.GetFloors(c.Request.Context(), req.HallID)
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

// ── PATCH /floors/:id ─────────────────────────────────────────────────────────

func (h *FloorHandler) UpdateFloor(c *gin.Context) {
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

	res, err := h.IFloorService.UpdateFloor(c.Request.Context(), floorID, &req)
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

// ── DELETE /floors/:id ────────────────────────────────────────────────────────

func (h *FloorHandler) DeleteFloor(c *gin.Context) {
	floorID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.WriteError(c, utils.ErrorInvalidIDFormart)
		return
	}

	if err := h.IFloorService.DeleteFloor(c.Request.Context(), floorID); err != nil {
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

// ── PUT /floors/reorder ───────────────────────────────────────────────────────

func (h *FloorHandler) ReorderFloors(c *gin.Context) {
	var req dto.ReorderFloorsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.WriteError(c, utils.ErrorInvalidInput)
		return
	}

	if err := h.IFloorService.ReorderFloors(c.Request.Context(), &req); err != nil {
		utils.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Floors reordered successfully",
		"success": true,
		"data":    nil,
	})
}
