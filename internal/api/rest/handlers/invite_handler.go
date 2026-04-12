// internal/api/rest/handlers/invite_handler.go
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/suck-seed/yapp/internal/auth"
	dto "github.com/suck-seed/yapp/internal/dto/hall"
	"github.com/suck-seed/yapp/internal/services"
	"github.com/suck-seed/yapp/internal/utils"
)

type InviteHandler struct {
	inviteService services.IInviteService
}

func NewInviteHandler(inviteService services.IInviteService) *InviteHandler {
	return &InviteHandler{inviteService}
}

// GET /halls/:hallID/settings/invites
func (h *InviteHandler) ListInviteLinks(c *gin.Context) {
	userInfo, err := auth.CurrentUserFromGinContext(c)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	hallID, err := uuid.Parse(c.Param("hallID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid hall id"})
		return
	}

	res, err := h.inviteService.ListInviteLinks(c.Request.Context(), userInfo, hallID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

// POST /halls/:hallID/settings/invites
func (h *InviteHandler) CreateInviteLink(c *gin.Context) {
	userInfo, err := auth.CurrentUserFromGinContext(c)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	hallID, err := uuid.Parse(c.Param("hallID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid hall id"})
		return
	}

	var req dto.CreateInviteLinkReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.inviteService.CreateInviteLink(c.Request.Context(), userInfo, hallID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, res)
}

// DELETE /halls/:hallID/settings/invites/:inviteID
func (h *InviteHandler) RevokeInviteLink(c *gin.Context) {
	userInfo, err := auth.CurrentUserFromGinContext(c)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	hallID, err := uuid.Parse(c.Param("hallID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid hall id"})
		return
	}
	inviteID, err := uuid.Parse(c.Param("inviteID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid invite id"})
		return
	}

	res, err := h.inviteService.RevokeInviteLink(c.Request.Context(), userInfo, hallID, inviteID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

// GET /invites/:code  — public, no auth required
func (h *InviteHandler) GetInviteLinkInfo(c *gin.Context) {
	code := c.Param("code")
	res, err := h.inviteService.GetInviteLinkInfo(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

// POST /invites/:code/accept  — requires auth
func (h *InviteHandler) AcceptInviteLink(c *gin.Context) {
	userInfo, err := auth.CurrentUserFromGinContext(c)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	code := c.Param("code")

	res, err := h.inviteService.AcceptInviteLink(c.Request.Context(), userInfo, code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}
