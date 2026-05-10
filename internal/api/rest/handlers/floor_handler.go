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

// CreateFloor godoc
// @Summary      Create a floor
// @Description  Adds a new floor (category) to the given hall.
// @Tags         floors
// @Accept       json
// @Produce      json
// @Security     CookieAuth
// @Param        hallID  path      string               true  "Hall ID (UUID)"
// @Param        body    body      dto.CreateFloorReq   true  "Floor details"
// @Success      201     {object}  map[string]interface{}
// @Failure      400     {object}  map[string]interface{}
// @Failure      401     {object}  map[string]interface{}
// @Failure      403     {object}  map[string]interface{}
// @Router       /halls/{hallID}/floors [post]
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

// GetFloors godoc
// @Summary      List all floors in a hall
// @Description  Returns every floor that belongs to the specified hall.
// @Tags         floors
// @Produce      json
// @Security     CookieAuth
// @Param        hallID  path      string  true  "Hall ID (UUID)"
// @Success      200     {object}  map[string]interface{}
// @Failure      400     {object}  map[string]interface{}
// @Failure      401     {object}  map[string]interface{}
// @Router       /halls/{hallID}/floors [get]
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

// GetFloor godoc
// @Summary      Get a single floor
// @Description  Returns the floor identified by {id} inside the given hall.
// @Tags         floors
// @Produce      json
// @Security     CookieAuth
// @Param        hallID  path      string  true  "Hall ID (UUID)"
// @Param        id path      string  true  "Floor ID (UUID)"
// @Success      200     {object}  map[string]interface{}
// @Failure      400     {object}  map[string]interface{}
// @Failure      401     {object}  map[string]interface{}
// @Failure      404     {object}  map[string]interface{}
// @Router       /halls/{hallID}/floors/{floorID} [get]
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

// UpdateFloor godoc
// @Summary      Update a floor
// @Description  Partially updates a floor's name or privacy setting.
// @Tags         floors
// @Accept       json
// @Produce      json
// @Security     CookieAuth
// @Param        hallID  path      string               true  "Hall ID (UUID)"
// @Param        id      path      string               true  "Floor ID (UUID)"
// @Param        body    body      dto.UpdateFloorReq   true  "Fields to update (all optional)"
// @Success      200     {object}  map[string]interface{}
// @Failure      400     {object}  map[string]interface{}
// @Failure      401     {object}  map[string]interface{}
// @Failure      403     {object}  map[string]interface{}
// @Router       /halls/{hallID}/floors/{id} [patch]
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

// DeleteFloor godoc
// @Summary      Delete a floor
// @Description  Permanently removes a floor and all rooms within it.
// @Tags         floors
// @Produce      json
// @Security     CookieAuth
// @Param        hallID  path      string  true  "Hall ID (UUID)"
// @Param        id      path      string  true  "Floor ID (UUID)"
// @Success      200     {object}  map[string]interface{}
// @Failure      400     {object}  map[string]interface{}
// @Failure      401     {object}  map[string]interface{}
// @Failure      403     {object}  map[string]interface{}
// @Router       /halls/{hallID}/floors/{id} [delete]
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

// MoveFloor godoc
// @Summary      Reorder a floor
// @Description  Moves a floor to a new position. Set `after_id` to null to place it at the top.
// @Tags         floors
// @Accept       json
// @Produce      json
// @Security     CookieAuth
// @Param        hallID  path      string             true  "Hall ID (UUID)"
// @Param        id      path      string             true  "Floor ID (UUID)"
// @Param        body    body      dto.MoveFloorReq   true  "Target position"
// @Success      200     {object}  map[string]interface{}
// @Failure      400     {object}  map[string]interface{}
// @Failure      401     {object}  map[string]interface{}
// @Failure      403     {object}  map[string]interface{}
// @Router       /halls/{hallID}/floors/{id}/move [put]
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

// AddFloorMember godoc
// @Summary      Add a member to a private floor
// @Description  Gives a hall member access to a private floor and syncs all synced rooms inside it.
// @Tags         floors
// @Produce      json
// @Security     CookieAuth
// @Param        hallID    path      string  true  "Hall ID (UUID)"
// @Param        id        path      string  true  "Floor ID (UUID)"
// @Param        memberID  path      string  true  "Hall Member ID (UUID)"
// @Success      200       {object}  map[string]interface{}
// @Failure      400       {object}  map[string]interface{}
// @Failure      401       {object}  map[string]interface{}
// @Failure      403       {object}  map[string]interface{}
// @Router       /halls/{hallID}/floors/{id}/members/{memberID} [put]
func (h *FloorHandler) AddFloorMember(c *gin.Context) {
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

	memberID, err := uuid.Parse(c.Param("memberID"))
	if err != nil {
		utils.WriteError(c, utils.ErrorInvalidIDFormart)
		return
	}

	res, err := h.IFloorService.AddFloorMember(c.Request.Context(), userInfo, hallID, floorID, memberID)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Floor member added successfully",
		"success": true,
		"data":    res,
	})
}

// RemoveFloorMember godoc
// @Summary      Remove a member from a private floor
// @Description  Removes a hall member's access from a private floor and syncs all synced rooms inside it.
// @Tags         floors
// @Produce      json
// @Security     CookieAuth
// @Param        hallID    path      string  true  "Hall ID (UUID)"
// @Param        id        path      string  true  "Floor ID (UUID)"
// @Param        memberID  path      string  true  "Hall Member ID (UUID)"
// @Success      200       {object}  map[string]interface{}
// @Failure      400       {object}  map[string]interface{}
// @Failure      401       {object}  map[string]interface{}
// @Failure      403       {object}  map[string]interface{}
// @Router       /halls/{hallID}/floors/{id}/members/{memberID} [delete]
func (h *FloorHandler) RemoveFloorMember(c *gin.Context) {
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

	memberID, err := uuid.Parse(c.Param("memberID"))
	if err != nil {
		utils.WriteError(c, utils.ErrorInvalidIDFormart)
		return
	}

	res, err := h.IFloorService.RemoveFloorMember(c.Request.Context(), userInfo, hallID, floorID, memberID)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Floor member removed successfully",
		"success": true,
		"data":    res,
	})
}
