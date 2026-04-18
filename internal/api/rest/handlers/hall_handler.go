package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/suck-seed/yapp/internal/auth"
	dto "github.com/suck-seed/yapp/internal/dto/hall"
	"github.com/suck-seed/yapp/internal/services"
	"github.com/suck-seed/yapp/internal/utils"
)

type HallHandler struct {
	services.IHallService
	services.IRoleService
	services.IBanService
}

func NewHallHandler(hallService services.IHallService, roleServices services.IRoleService, banServices services.IBanService) *HallHandler {
	return &HallHandler{
		hallService,
		roleServices,
		banServices,
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// TOP-LEVEL HALL OPERATIONS
// ─────────────────────────────────────────────────────────────────────────────

// CreateHall godoc
// @Summary      Create a hall
// @Description  Creates a new hall owned by the authenticated user.
// @Tags         halls
// @Accept       json
// @Produce      json
// @Security     CookieAuth
// @Param        body  body      dto.CreateHallReq      true  "Hall details"
// @Success      201   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Failure      401   {object}  map[string]interface{}
// @Router       /halls [post]
func (h *HallHandler) CreateHall(c *gin.Context) {

	u := &dto.CreateHallReq{}

	if err := c.ShouldBindJSON(u); err != nil {
		utils.WriteError(c, utils.ErrorInvalidInput)
		return
	}

	userInfo, err := auth.CurrentUserFromGinContext(c)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	res, err := h.IHallService.CreateHall(c.Request.Context(), userInfo, u)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Hall created successfully",
		"data":    res,
	})

}

// JoinHall godoc
// @Summary      Join or request to join a hall
// @Description  For public halls the user joins immediately; for private halls a join-request is created.
// @Tags         halls
// @Produce      json
// @Security     CookieAuth
// @Param        hallID  path      string  true  "Hall ID (UUID)"
// @Success      200     {object}  map[string]interface{}  "joined or requested"
// @Failure      400     {object}  map[string]interface{}
// @Failure      401     {object}  map[string]interface{}
// @Failure      403     {object}  map[string]interface{}  "banned"
// @Router       /halls/{hallID}/join [post]
func (h *HallHandler) JoinHall(c *gin.Context) {
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

	res, err := h.IHallService.JoinHall(c.Request.Context(), userInfo, hallID)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	message := "Joined hall successfully"
	if res.Status == "requested" {
		message = "Join request created successfully"
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": message,
		"data":    res,
	})
}

// GetUserHalls godoc
// @Summary      List halls for the current user
// @Description  Returns all halls that the authenticated user is a member of.
// @Tags         halls
// @Produce      json
// @Security     CookieAuth
// @Success      200  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Router       /halls [get]
func (h *HallHandler) GetUserHalls(c *gin.Context) {

	userInfo, err := auth.CurrentUserFromGinContext(c)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	res, err := h.IHallService.GetUserHalls(c.Request.Context(), userInfo)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Halls retrieved successfully",
		"data":    res,
	})
}

// GetCurrentHall godoc
// @Summary      Get a hall
// @Description  Returns basic info for a single hall.
// @Tags         halls
// @Produce      json
// @Security     CookieAuth
// @Param        hallID  path      string  true  "Hall ID (UUID)"
// @Success      200     {object}  map[string]interface{}
// @Failure      400     {object}  map[string]interface{}
// @Failure      401     {object}  map[string]interface{}
// @Failure      404     {object}  map[string]interface{}
// @Router       /halls/{hallID} [get]
func (h *HallHandler) GetCurrentHall(c *gin.Context) {

	hallID, err := uuid.Parse(c.Param("hallID"))
	if err != nil {
		utils.WriteError(c, utils.ErrorInvalidIDFormart)
		return
	}

	userInfo, err := auth.CurrentUserFromGinContext(c)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	res, err := h.IHallService.GetCurrentHall(c.Request.Context(), userInfo, hallID)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Hall retrieved successfully",
		"data":    res,
	})

}

// DeleteCurrentHall godoc
// @Summary      Delete a hall
// @Description  Permanently deletes the hall. Only the owner may perform this action.
// @Tags         halls
// @Produce      json
// @Security     CookieAuth
// @Param        hallID  path      string  true  "Hall ID (UUID)"
// @Success      200     {object}  map[string]interface{}
// @Failure      400     {object}  map[string]interface{}
// @Failure      401     {object}  map[string]interface{}
// @Failure      403     {object}  map[string]interface{}
// @Router       /halls/{hallID} [delete]
func (h *HallHandler) DeleteCurrentHall(c *gin.Context) {

	hallID, err := uuid.Parse(c.Param("hallID"))
	if err != nil {
		utils.WriteError(c, utils.ErrorInvalidIDFormart)
		return
	}

	userInfo, err := auth.CurrentUserFromGinContext(c)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	res, err := h.IHallService.DeleteHall(c.Request.Context(), userInfo, hallID)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Hall deleted successfully",
		"data":    res,
	})

}

// ─────────────────────────────────────────────────────────────────────────────
// SETTINGS — PROFILE
// ─────────────────────────────────────────────────────────────────────────────

// GetHallProfile godoc
// @Summary      Get hall profile / settings
// @Description  Returns the full editable profile for the hall (name, icon, banner, description).
// @Tags         hall-settings
// @Produce      json
// @Security     CookieAuth
// @Param        hallID  path      string  true  "Hall ID (UUID)"
// @Success      200     {object}  map[string]interface{}
// @Failure      400     {object}  map[string]interface{}
// @Failure      401     {object}  map[string]interface{}
// @Failure      403     {object}  map[string]interface{}
// @Router       /halls/{hallID}/settings/profile [get]
func (h *HallHandler) GetHallProfile(c *gin.Context) {

	hallID, err := uuid.Parse(c.Param("hallID"))
	if err != nil {
		utils.WriteError(c, utils.ErrorInvalidIDFormart)
		return
	}

	log.Println(hallID)

	userInfo, err := auth.CurrentUserFromGinContext(c)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	log.Println(userInfo)

	res, err := h.IHallService.GetHallProfile(c.Request.Context(), userInfo, hallID)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	log.Println(res)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Hall profile retrieved successfully",
		"data":    res,
	})

}

// UpdateHallProfile godoc
// @Summary      Update hall profile
// @Description  Updates the hall name, icon, banner colour, or description.
// @Tags         hall-settings
// @Accept       json
// @Produce      json
// @Security     CookieAuth
// @Param        hallID  path      string                    true  "Hall ID (UUID)"
// @Param        body    body      dto.HallProfileUpdateReq  true  "Profile fields to update"
// @Success      200     {object}  map[string]interface{}
// @Failure      400     {object}  map[string]interface{}
// @Failure      401     {object}  map[string]interface{}
// @Failure      403     {object}  map[string]interface{}
// @Router       /halls/{hallID}/settings/profile [patch]
func (h *HallHandler) UpdateHallProfile(c *gin.Context) {

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

	var req dto.HallProfileUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.WriteError(c, err)
		return
	}

	res, err := h.IHallService.UpdateHallProfile(c.Request.Context(), userInfo, hallID, &req)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Hall profile updated successfully",
		"data":    res,
	})

}

// ─────────────────────────────────────────────────────────────────────────────
// SETTINGS — MEMBERS
// ─────────────────────────────────────────────────────────────────────────────

// GetHallMembers godoc
// @Summary      List hall members
// @Description  Returns all members of a hall with their roles and nicknames.
// @Tags         hall-settings
// @Produce      json
// @Security     CookieAuth
// @Param        hallID  path      string  true  "Hall ID (UUID)"
// @Success      200     {object}  map[string]interface{}
// @Failure      400     {object}  map[string]interface{}
// @Failure      401     {object}  map[string]interface{}
// @Router       /halls/{hallID}/settings/members [get]
func (h *HallHandler) GetHallMembers(c *gin.Context) {
	hallID, err := uuid.Parse(c.Param("hallID"))
	if err != nil {
		utils.WriteError(c, utils.ErrorInvalidIDFormart)
		return
	}

	userInfo, err := auth.CurrentUserFromGinContext(c)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	res, err := h.IHallService.GetHallMembers(c.Request.Context(), userInfo, hallID)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Hall members retrieved successfully",
		"data":    res,
	})
}

// GetHallMember godoc
// @Summary      Get a hall member
// @Description  Returns a single member's details within a hall.
// @Tags         hall-settings
// @Produce      json
// @Security     CookieAuth
// @Param        hallID    path      string  true  "Hall ID (UUID)"
// @Param        memberID  path      string  true  "Member ID (UUID)"
// @Success      200       {object}  map[string]interface{}
// @Failure      400       {object}  map[string]interface{}
// @Failure      401       {object}  map[string]interface{}
// @Failure      404       {object}  map[string]interface{}
// @Router       /halls/{hallID}/settings/members/{memberID} [get]
func (h *HallHandler) GetHallMember(c *gin.Context) {
	hallID, err := uuid.Parse(c.Param("hallID"))
	if err != nil {
		utils.WriteError(c, utils.ErrorInvalidIDFormart)
		return
	}

	memberID, err := uuid.Parse(c.Param("memberID"))
	if err != nil {
		utils.WriteError(c, utils.ErrorInvalidIDFormart)
		return
	}

	userInfo, err := auth.CurrentUserFromGinContext(c)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	res, err := h.IHallService.GetHallMember(c.Request.Context(), userInfo, hallID, memberID)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Hall member retrieved successfully",
		"data":    res,
	})
}

// UpdateHallMemberRole godoc
// @Summary      Change a member's role
// @Description  Assigns a different role to a hall member. Requires ManageRoles permission.
// @Tags         hall-settings
// @Accept       json
// @Produce      json
// @Security     CookieAuth
// @Param        hallID    path      string                       true  "Hall ID (UUID)"
// @Param        memberID  path      string                       true  "Member ID (UUID)"
// @Param        body      body      dto.UpdateHallMemberRoleReq  true  "New role"
// @Success      200       {object}  map[string]interface{}
// @Failure      400       {object}  map[string]interface{}
// @Failure      401       {object}  map[string]interface{}
// @Failure      403       {object}  map[string]interface{}
// @Router       /halls/{hallID}/settings/members/{memberID}/role [patch]
func (h *HallHandler) UpdateHallMemberRole(c *gin.Context) {
	hallID, err := uuid.Parse(c.Param("hallID"))
	if err != nil {
		utils.WriteError(c, utils.ErrorInvalidIDFormart)
		return
	}

	memberID, err := uuid.Parse(c.Param("memberID"))
	if err != nil {
		utils.WriteError(c, utils.ErrorInvalidIDFormart)
		return
	}

	var req dto.UpdateHallMemberRoleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.WriteError(c, utils.ErrorInvalidInput)
		return
	}

	userInfo, err := auth.CurrentUserFromGinContext(c)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	res, err := h.IHallService.UpdateHallMemberRole(c.Request.Context(), userInfo, hallID, memberID, &req)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Hall member role updated successfully",
		"data":    res,
	})
}

// UpdateHallMemberNickname godoc
// @Summary      Change a member's nickname
// @Description  Sets or clears a member's display nickname within the hall.
// @Tags         hall-settings
// @Accept       json
// @Produce      json
// @Security     CookieAuth
// @Param        hallID    path      string                           true  "Hall ID (UUID)"
// @Param        memberID  path      string                           true  "Member ID (UUID)"
// @Param        body      body      dto.UpdateHallMemberNicknameReq  true  "Nickname payload"
// @Success      200       {object}  map[string]interface{}
// @Failure      400       {object}  map[string]interface{}
// @Failure      401       {object}  map[string]interface{}
// @Failure      403       {object}  map[string]interface{}
// @Router       /halls/{hallID}/settings/members/{memberID}/nickname [patch]
func (h *HallHandler) UpdateHallMemberNickname(c *gin.Context) {
	hallID, err := uuid.Parse(c.Param("hallID"))
	if err != nil {
		utils.WriteError(c, utils.ErrorInvalidIDFormart)
		return
	}

	memberID, err := uuid.Parse(c.Param("memberID"))
	if err != nil {
		utils.WriteError(c, utils.ErrorInvalidIDFormart)
		return
	}

	var req dto.UpdateHallMemberNicknameReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.WriteError(c, utils.ErrorInvalidInput)
		return
	}

	userInfo, err := auth.CurrentUserFromGinContext(c)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	res, err := h.IHallService.UpdateHallMemberNickname(c.Request.Context(), userInfo, hallID, memberID, &req)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Hall member nickname updated successfully",
		"data":    res,
	})
}

// KickHallMember godoc
// @Summary      Kick a member from the hall
// @Description  Removes a member from the hall without banning them. Requires KickMembers permission.
// @Tags         hall-settings
// @Produce      json
// @Security     CookieAuth
// @Param        hallID    path      string  true  "Hall ID (UUID)"
// @Param        memberID  path      string  true  "Member ID (UUID)"
// @Success      200       {object}  map[string]interface{}
// @Failure      400       {object}  map[string]interface{}
// @Failure      401       {object}  map[string]interface{}
// @Failure      403       {object}  map[string]interface{}
// @Router       /halls/{hallID}/settings/members/{memberID}/kick [delete]
func (h *HallHandler) KickHallMember(c *gin.Context) {
	hallID, err := uuid.Parse(c.Param("hallID"))
	if err != nil {
		utils.WriteError(c, utils.ErrorInvalidIDFormart)
		return
	}

	memberID, err := uuid.Parse(c.Param("memberID"))
	if err != nil {
		utils.WriteError(c, utils.ErrorInvalidIDFormart)
		return
	}

	userInfo, err := auth.CurrentUserFromGinContext(c)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	res, err := h.IHallService.KickHallMember(c.Request.Context(), userInfo, hallID, memberID)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Member removed from hall successfully",
		"data":    res,
	})
}

// ─────────────────────────────────────────────────────────────────────────────
// SETTINGS — ROLES
// ─────────────────────────────────────────────────────────────────────────────

// GetHallRoles godoc
// @Summary      List roles in a hall
// @Description  Returns all roles defined for the hall.
// @Tags         hall-settings
// @Produce      json
// @Security     CookieAuth
// @Param        hallID  path      string  true  "Hall ID (UUID)"
// @Success      200     {object}  map[string]interface{}
// @Failure      400     {object}  map[string]interface{}
// @Failure      401     {object}  map[string]interface{}
// @Router       /halls/{hallID}/settings/roles [get]
func (h *HallHandler) GetHallRoles(c *gin.Context) {
	hallID, err := uuid.Parse(c.Param("hallID"))
	if err != nil {
		utils.WriteError(c, utils.ErrorInvalidIDFormart)
		return
	}

	userInfo, err := auth.CurrentUserFromGinContext(c)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	res, err := h.IRoleService.ListHallRoles(c.Request.Context(), userInfo, hallID)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Hall roles retrieved successfully",
		"data":    res,
	})
}

// GetHallRole godoc
// @Summary      Get a single role
// @Description  Returns a specific role by ID.
// @Tags         hall-settings
// @Produce      json
// @Security     CookieAuth
// @Param        hallID  path      string  true  "Hall ID (UUID)"
// @Param        roleID  path      string  true  "Role ID (UUID)"
// @Success      200     {object}  map[string]interface{}
// @Failure      400     {object}  map[string]interface{}
// @Failure      401     {object}  map[string]interface{}
// @Failure      404     {object}  map[string]interface{}
// @Router       /halls/{hallID}/settings/roles/{roleID} [get]
func (h *HallHandler) GetHallRole(c *gin.Context) {
	hallID, err := uuid.Parse(c.Param("hallID"))
	if err != nil {
		utils.WriteError(c, utils.ErrorInvalidIDFormart)
		return
	}

	roleID, err := uuid.Parse(c.Param("roleID"))
	if err != nil {
		utils.WriteError(c, utils.ErrorInvalidIDFormart)
		return
	}

	userInfo, err := auth.CurrentUserFromGinContext(c)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	res, err := h.IRoleService.GetHallRole(c.Request.Context(), userInfo, hallID, roleID)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Hall role retrieved successfully",
		"data":    res,
	})
}

// CreateHallRoles godoc
// @Summary      Create a role
// @Description  Creates a new role in the hall. Requires ManageRoles permission.
// @Tags         hall-settings
// @Accept       json
// @Produce      json
// @Security     CookieAuth
// @Param        hallID  path      string                  true  "Hall ID (UUID)"
// @Param        body    body      dto.CreateHallRoleReq   true  "Role details"
// @Success      201     {object}  map[string]interface{}
// @Failure      400     {object}  map[string]interface{}
// @Failure      401     {object}  map[string]interface{}
// @Failure      403     {object}  map[string]interface{}
// @Router       /halls/{hallID}/settings/roles [post]
func (h *HallHandler) CreateHallRoles(c *gin.Context) {
	hallID, err := uuid.Parse(c.Param("hallID"))
	if err != nil {
		utils.WriteError(c, utils.ErrorInvalidIDFormart)
		return
	}

	var req dto.CreateHallRoleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.WriteError(c, utils.ErrorInvalidInput)
		return
	}

	userInfo, err := auth.CurrentUserFromGinContext(c)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	res, err := h.IRoleService.CreateHallRole(c.Request.Context(), userInfo, hallID, &req)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Hall role created successfully",
		"data":    res,
	})
}

// UpdateHallRoles godoc
// @Summary      Update a role
// @Description  Renames a role or changes its colour / icon. Requires ManageRoles permission.
// @Tags         hall-settings
// @Accept       json
// @Produce      json
// @Security     CookieAuth
// @Param        hallID  path      string                  true  "Hall ID (UUID)"
// @Param        roleID  path      string                  true  "Role ID (UUID)"
// @Param        body    body      dto.UpdateHallRoleReq   true  "Fields to update"
// @Success      200     {object}  map[string]interface{}
// @Failure      400     {object}  map[string]interface{}
// @Failure      401     {object}  map[string]interface{}
// @Failure      403     {object}  map[string]interface{}
// @Router       /halls/{hallID}/settings/roles/{roleID} [patch]
func (h *HallHandler) UpdateHallRoles(c *gin.Context) {
	hallID, err := uuid.Parse(c.Param("hallID"))
	if err != nil {
		utils.WriteError(c, utils.ErrorInvalidIDFormart)
		return
	}

	roleID, err := uuid.Parse(c.Param("roleID"))
	if err != nil {
		utils.WriteError(c, utils.ErrorInvalidIDFormart)
		return
	}

	var req dto.UpdateHallRoleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.WriteError(c, utils.ErrorInvalidInput)
		return
	}

	userInfo, err := auth.CurrentUserFromGinContext(c)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	res, err := h.IRoleService.UpdateHallRole(c.Request.Context(), userInfo, hallID, roleID, &req)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Hall role updated successfully",
		"data":    res,
	})
}

// DeleteHallRoles godoc
// @Summary      Delete a role
// @Description  Permanently removes a role from the hall. Requires ManageRoles permission.
// @Tags         hall-settings
// @Produce      json
// @Security     CookieAuth
// @Param        hallID  path      string  true  "Hall ID (UUID)"
// @Param        roleID  path      string  true  "Role ID (UUID)"
// @Success      200     {object}  map[string]interface{}
// @Failure      400     {object}  map[string]interface{}
// @Failure      401     {object}  map[string]interface{}
// @Failure      403     {object}  map[string]interface{}
// @Router       /halls/{hallID}/settings/roles/{roleID} [delete]
func (h *HallHandler) DeleteHallRoles(c *gin.Context) {
	hallID, err := uuid.Parse(c.Param("hallID"))
	if err != nil {
		utils.WriteError(c, utils.ErrorInvalidIDFormart)
		return
	}

	roleID, err := uuid.Parse(c.Param("roleID"))
	if err != nil {
		utils.WriteError(c, utils.ErrorInvalidIDFormart)
		return
	}

	userInfo, err := auth.CurrentUserFromGinContext(c)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	res, err := h.IRoleService.DeleteHallRole(c.Request.Context(), userInfo, hallID, roleID)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Hall role deleted successfully",
		"data":    res,
	})
}

// ─────────────────────────────────────────────────────────────────────────────
// SETTINGS — ROLE PERMISSIONS
// ─────────────────────────────────────────────────────────────────────────────

// GetRolesPermissions godoc
// @Summary      Get role permissions
// @Description  Returns the full permission set for a specific role.
// @Tags         hall-settings
// @Produce      json
// @Security     CookieAuth
// @Param        hallID  path      string  true  "Hall ID (UUID)"
// @Param        roleID  path      string  true  "Role ID (UUID)"
// @Success      200     {object}  map[string]interface{}
// @Failure      400     {object}  map[string]interface{}
// @Failure      401     {object}  map[string]interface{}
// @Failure      403     {object}  map[string]interface{}
// @Router       /halls/{hallID}/settings/roles/{roleID}/permissions [get]
func (h *HallHandler) GetRolesPermissions(c *gin.Context) {

	// fetch hallID and roleID from the params
	hallID, err := uuid.Parse(c.Param("hallID"))
	if err != nil {
		utils.WriteError(c, utils.ErrorInternal)
		return
	}

	roleID, err := uuid.Parse(c.Param("roleID"))
	if err != nil {
		utils.WriteError(c, utils.ErrorInternal)
		return
	}

	// fetch userinformation from the current instance of gin context
	userInfo, err := auth.CurrentUserFromGinContext(c)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	res, err := h.IRoleService.GetRolePermissions(c.Request.Context(), userInfo, hallID, roleID)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Role permissions retrieved successfully",
		"data":    res,
	})

}

// UpdateRolesPermissions godoc
// @Summary      Update role permissions
// @Description  Overwrites the permission flags for a role. Requires ManageRoles permission.
// @Tags         hall-settings
// @Accept       json
// @Produce      json
// @Security     CookieAuth
// @Param        hallID  path      string                       true  "Hall ID (UUID)"
// @Param        roleID  path      string                       true  "Role ID (UUID)"
// @Param        body    body      dto.UpdateRolePermissionReq  true  "Permission flags"
// @Success      200     {object}  map[string]interface{}
// @Failure      400     {object}  map[string]interface{}
// @Failure      401     {object}  map[string]interface{}
// @Failure      403     {object}  map[string]interface{}
// @Router       /halls/{hallID}/settings/roles/{roleID}/permissions [patch]
func (h *HallHandler) UpdateRolesPermissions(c *gin.Context) {

	// fetch hallID and roleID from the params
	hallID, err := uuid.Parse(c.Param("hallID"))
	if err != nil {
		utils.WriteError(c, utils.ErrorInternal)
		return
	}

	roleID, err := uuid.Parse(c.Param("roleID"))
	if err != nil {
		utils.WriteError(c, utils.ErrorInternal)
		return
	}

	// RolePermissionUpdate struct binding to the requesting json payload
	u := &dto.UpdateRolePermissionReq{}

	if err := c.ShouldBindJSON(u); err != nil {
		utils.WriteError(c, utils.ErrorInvalidInput)
		return
	}

	userInfo, err := auth.CurrentUserFromGinContext(c)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	// call the service layer implementation responsible for updating the role's permissions
	res, err := h.IRoleService.UpdateRolePermissions(c.Request.Context(), userInfo, hallID, roleID, u)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Role permissions updated successfully",
		"data":    res,
	})
}

// ─────────────────────────────────────────────────────────────────────────────
// SETTINGS — JOIN REQUESTS
// ─────────────────────────────────────────────────────────────────────────────

// GetCurrentRequests godoc
// @Summary      List pending join requests
// @Description  Returns all outstanding join requests for a private hall.
// @Tags         hall-settings
// @Produce      json
// @Security     CookieAuth
// @Param        hallID  path      string  true  "Hall ID (UUID)"
// @Success      200     {object}  map[string]interface{}
// @Failure      400     {object}  map[string]interface{}
// @Failure      401     {object}  map[string]interface{}
// @Failure      403     {object}  map[string]interface{}
// @Router       /halls/{hallID}/settings/requests [get]
func (h *HallHandler) GetCurrentRequests(c *gin.Context) {
	hallID, err := uuid.Parse(c.Param("hallID"))
	if err != nil {
		utils.WriteError(c, utils.ErrorInvalidIDFormart)
		return
	}

	userInfo, err := auth.CurrentUserFromGinContext(c)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	res, err := h.IHallService.GetCurrentRequests(c.Request.Context(), userInfo, hallID)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Join requests retrieved successfully",
		"data":    res,
	})
}

// AcceptJoinRequest godoc
// @Summary      Accept a join request
// @Description  Approves a pending join request and adds the user to the hall. Requires ManageRequests permission.
// @Tags         hall-settings
// @Produce      json
// @Security     CookieAuth
// @Param        hallID     path      string  true  "Hall ID (UUID)"
// @Param        requestID  path      string  true  "Request ID (UUID)"
// @Success      200        {object}  map[string]interface{}
// @Failure      400        {object}  map[string]interface{}
// @Failure      401        {object}  map[string]interface{}
// @Failure      403        {object}  map[string]interface{}
// @Router       /halls/{hallID}/settings/requests/{requestID}/accept [patch]
func (h *HallHandler) AcceptJoinRequest(c *gin.Context) {
	hallID, err := uuid.Parse(c.Param("hallID"))
	if err != nil {
		utils.WriteError(c, utils.ErrorInvalidIDFormart)
		return
	}

	requestID, err := uuid.Parse(c.Param("requestID"))
	if err != nil {
		utils.WriteError(c, utils.ErrorInvalidIDFormart)
		return
	}

	userInfo, err := auth.CurrentUserFromGinContext(c)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	res, err := h.IHallService.AcceptJoinRequest(c.Request.Context(), userInfo, hallID, requestID)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Join request accepted successfully",
		"data":    res,
	})
}

// DeclineJoinRequest godoc
// @Summary      Decline a join request
// @Description  Rejects and removes a pending join request. Requires ManageRequests permission.
// @Tags         hall-settings
// @Produce      json
// @Security     CookieAuth
// @Param        hallID     path      string  true  "Hall ID (UUID)"
// @Param        requestID  path      string  true  "Request ID (UUID)"
// @Success      200        {object}  map[string]interface{}
// @Failure      400        {object}  map[string]interface{}
// @Failure      401        {object}  map[string]interface{}
// @Failure      403        {object}  map[string]interface{}
// @Router       /halls/{hallID}/settings/requests/{requestID} [delete]
func (h *HallHandler) DeclineJoinRequest(c *gin.Context) {
	hallID, err := uuid.Parse(c.Param("hallID"))
	if err != nil {
		utils.WriteError(c, utils.ErrorInvalidIDFormart)
		return
	}

	requestID, err := uuid.Parse(c.Param("requestID"))
	if err != nil {
		utils.WriteError(c, utils.ErrorInvalidIDFormart)
		return
	}

	userInfo, err := auth.CurrentUserFromGinContext(c)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	res, err := h.IHallService.DeclineJoinRequest(c.Request.Context(), userInfo, hallID, requestID)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Join request declined successfully",
		"data":    res,
	})
}

// ─────────────────────────────────────────────────────────────────────────────
// SETTINGS — BANS
// ─────────────────────────────────────────────────────────────────────────────

// GetBannedUsers godoc
// @Summary      List banned users
// @Description  Returns all active bans for the hall. Requires BanMembers permission.
// @Tags         hall-settings
// @Produce      json
// @Security     CookieAuth
// @Param        hallID  path      string  true  "Hall ID (UUID)"
// @Success      200     {object}  map[string]interface{}
// @Failure      400     {object}  map[string]interface{}
// @Failure      401     {object}  map[string]interface{}
// @Failure      403     {object}  map[string]interface{}
// @Router       /halls/{hallID}/settings/bans [get]
func (h *HallHandler) GetBannedUsers(c *gin.Context) {
	hallID, err := uuid.Parse(c.Param("hallID"))
	if err != nil {
		utils.WriteError(c, utils.ErrorInvalidIDFormart)
		return
	}

	userInfo, err := auth.CurrentUserFromGinContext(c)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	res, err := h.IBanService.GetAllHallBans(c.Request.Context(), userInfo, hallID)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Banned users retrieved successfully",
		"data":    res,
	})
}

// BanAnUser godoc
// @Summary      Ban a user
// @Description  Kicks and permanently bans a user from the hall. Requires BanMembers permission.
// @Tags         hall-settings
// @Accept       json
// @Produce      json
// @Security     CookieAuth
// @Param        hallID  path      string           true  "Hall ID (UUID)"
// @Param        body    body      dto.BanUserReq   true  "Ban payload (user + reason)"
// @Success      201     {object}  map[string]interface{}
// @Failure      400     {object}  map[string]interface{}
// @Failure      401     {object}  map[string]interface{}
// @Failure      403     {object}  map[string]interface{}
// @Router       /halls/{hallID}/settings/bans [post]
func (h *HallHandler) BanAnUser(c *gin.Context) {
	hallID, err := uuid.Parse(c.Param("hallID"))
	if err != nil {
		utils.WriteError(c, utils.ErrorInvalidIDFormart)
		return
	}

	var req dto.BanUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.WriteError(c, utils.ErrorInvalidInput)
		return
	}

	userInfo, err := auth.CurrentUserFromGinContext(c)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	res, err := h.IBanService.BanUser(c.Request.Context(), userInfo, hallID, &req)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "User banned successfully",
		"data":    res,
	})
}

// UnbanUser godoc
// @Summary      Unban a user
// @Description  Lifts an existing ban so the user may rejoin. Requires BanMembers permission.
// @Tags         hall-settings
// @Produce      json
// @Security     CookieAuth
// @Param        hallID  path      string  true  "Hall ID (UUID)"
// @Param        banID   path      string  true  "Ban ID (UUID)"
// @Success      200     {object}  map[string]interface{}
// @Failure      400     {object}  map[string]interface{}
// @Failure      401     {object}  map[string]interface{}
// @Failure      403     {object}  map[string]interface{}
// @Router       /halls/{hallID}/settings/bans/{banID} [delete]
func (h *HallHandler) UnbanUser(c *gin.Context) {
	hallID, err := uuid.Parse(c.Param("hallID"))
	if err != nil {
		utils.WriteError(c, utils.ErrorInvalidIDFormart)
		return
	}

	banID, err := uuid.Parse(c.Param("banID"))
	if err != nil {
		utils.WriteError(c, utils.ErrorInvalidIDFormart)
		return
	}

	userInfo, err := auth.CurrentUserFromGinContext(c)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	res, err := h.IBanService.UnbanUser(c.Request.Context(), userInfo, hallID, banID)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User unbanned successfully",
		"data":    res,
	})
}
