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

// ListInviteLinks godoc
// @Summary      List invite links for a hall
// @Description  Returns all active invite links for the hall. Requires ManageInvites permission.
// @Tags         invites
// @Produce      json
// @Security     CookieAuth
// @Param        hallID  path      string  true  "Hall ID (UUID)"
// @Success      200     {object}  map[string]interface{}
// @Failure      400     {object}  map[string]interface{}
// @Failure      401     {object}  map[string]interface{}
// @Failure      403     {object}  map[string]interface{}
// @Router       /halls/{hallID}/settings/invites [get]
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
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Current Invite Links",
		"data":    res,
	})
}

// CreateInviteLink godoc
// @Summary      Create an invite link
// @Description  Generates a new invite link with optional expiry and max-use limits. Requires ManageInvites permission.
// @Tags         invites
// @Accept       json
// @Produce      json
// @Security     CookieAuth
// @Param        hallID  path      string                    true  "Hall ID (UUID)"
// @Param        body    body      dto.CreateInviteLinkReq   true  "Invite options"
// @Success      201     {object}  map[string]interface{}
// @Failure      400     {object}  map[string]interface{}
// @Failure      401     {object}  map[string]interface{}
// @Failure      403     {object}  map[string]interface{}
// @Router       /halls/{hallID}/settings/invites [post]
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

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Invite Link created sucessfully",
		"data":    res,
	})
}

// RevokeInviteLink godoc
// @Summary      Revoke an invite link
// @Description  Permanently deletes an invite link so it can no longer be used. Requires ManageInvites permission.
// @Tags         invites
// @Produce      json
// @Security     CookieAuth
// @Param        hallID    path      string  true  "Hall ID (UUID)"
// @Param        inviteID  path      string  true  "Invite ID (UUID)"
// @Success      200       {object}  map[string]interface{}
// @Failure      400       {object}  map[string]interface{}
// @Failure      401       {object}  map[string]interface{}
// @Failure      403       {object}  map[string]interface{}
// @Router       /halls/{hallID}/settings/invites/{inviteID} [delete]
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
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Invite link revoked sucessfully",
		"data":    res,
	})
}

// GetInviteLinkInfo godoc
// @Summary      Get invite info (public)
// @Description  Returns hall preview info for an invite code without requiring authentication.
// @Tags         invites
// @Produce      json
// @Param        code  path      string  true  "Invite code"
// @Success      200   {object}  map[string]interface{}
// @Failure      404   {object}  map[string]interface{}  "Invalid or expired code"
// @Router       /invites/{code} [get]
func (h *InviteHandler) GetInviteLinkInfo(c *gin.Context) {

	// GET /invites/:code  — public, no auth required

	code := c.Param("code")
	res, err := h.inviteService.GetInviteLinkInfo(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Invite Links fetched successfully",
		"data":    res,
	})
}

// AcceptInviteLink godoc
// @Summary      Accept an invite link
// @Description  Joins the hall associated with the given invite code. Consumes one use if the link has a max-use limit.
// @Tags         invites
// @Produce      json
// @Security     CookieAuth
// @Param        code  path      string  true  "Invite code"
// @Success      202   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}  "Code expired, exhausted, or already a member"
// @Failure      401   {object}  map[string]interface{}
// @Router       /invites/{code}/accept [post]
func (h *InviteHandler) AcceptInviteLink(c *gin.Context) {
	userInfo, err := auth.CurrentUserFromGinContext(c)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	code := c.Param("code")

	res, err := h.inviteService.AcceptInviteLink(c.Request.Context(), userInfo, code)
	if err != nil {
		utils.WriteError(c, err)
		return
	}
	c.JSON(http.StatusAccepted, gin.H{
		"code":    http.StatusCreated,
		"success": true,
		"message": "Invite accepted successfully",
		"data":    res,
	})
}
