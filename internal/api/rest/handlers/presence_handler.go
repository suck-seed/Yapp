package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/suck-seed/yapp/internal/auth"
	dto "github.com/suck-seed/yapp/internal/dto/user"
	"github.com/suck-seed/yapp/internal/services"
	"github.com/suck-seed/yapp/internal/utils"
)

type PresenceHandler struct {
	services.IPresenceService
}

func NewPresenceHandler(presenceService services.IPresenceService) *PresenceHandler {
	return &PresenceHandler{presenceService}
}

func (h *PresenceHandler) GetMyPresence(c *gin.Context) {
	userInfo, err := auth.CurrentUserFromGinContext(c)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	res, err := h.IPresenceService.GetUserPresence(c.Request.Context(), userInfo.ID)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Presence fetched successfully",
		"data":    res,
	})
}

func (h *PresenceHandler) UpdateMyPresence(c *gin.Context) {
	userInfo, err := auth.CurrentUserFromGinContext(c)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	var req dto.UpdatePresenceReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.WriteError(c, utils.ErrorInvalidInput)
		return
	}

	res, err := h.IPresenceService.SetManualStatus(c.Request.Context(), userInfo.ID, req.Status)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Presence updated successfully",
		"data":    res,
	})
}

func (h *PresenceHandler) GetManyPresence(c *gin.Context) {
	rawIDs := c.Query("ids")
	if rawIDs == "" {
		utils.WriteError(c, utils.ErrorInvalidInput)
		return
	}

	parts := strings.Split(rawIDs, ",")
	userIDs := make([]uuid.UUID, 0, len(parts))

	for _, part := range parts {
		id, err := uuid.Parse(strings.TrimSpace(part))
		if err != nil {
			utils.WriteError(c, utils.ErrorInvalidIDFormart)
			return
		}
		userIDs = append(userIDs, id)
	}

	res, err := h.IPresenceService.GetManyPresences(c.Request.Context(), userIDs)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Presence fetched successfully",
		"data": gin.H{
			"presences": res,
		},
	})
}
