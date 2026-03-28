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

// TOP LEVEL HALL OPERATIONS
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

	c.JSON(http.StatusCreated, res)

}

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

	c.JSON(http.StatusOK, res)
}

// SINGLE HALL RUD
func (h *HallHandler) GetCurrentHall(c *gin.Context) {

	hallID, err := uuid.Parse(c.Param("hallID"))
	if err != nil {
		utils.WriteError(c, utils.ErrorInternal)
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

	c.JSON(http.StatusOK, res)

}

// Depriciated not needed
// func (h *HallHandler) UpdateCurrentHall(c *gin.Context) {

// }

func (h *HallHandler) DeleteCurrentHall(c *gin.Context) {

	hallID, err := uuid.Parse(c.Param("hallID"))
	if err != nil {
		utils.WriteError(c, utils.ErrorInternal)
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

	c.JSON(http.StatusOK, res)

}

//----------------SETTING SCOPE

// PROFILE MANAGEMENT
func (h *HallHandler) GetHallProfile(c *gin.Context) {

}

func (h *HallHandler) UpdateHallProfile(c *gin.Context) {

}

// MEMBERS MANAGEMENT
func (h *HallHandler) GetHallMembers(c *gin.Context) {

}

func (h *HallHandler) GetHallMember(c *gin.Context) {

}

func (h *HallHandler) UpdateHallMember(c *gin.Context) {

}

func (h *HallHandler) RemoveHallMember(c *gin.Context) {

}

// ROLE MANAGEMENT
func (h *HallHandler) GetHallRoles(c *gin.Context) {

}

func (h *HallHandler) GetHallRole(c *gin.Context) {

}

func (h *HallHandler) CreateHallRoles(c *gin.Context) {

}

func (h *HallHandler) UpdateHallRoles(c *gin.Context) {

}

func (h *HallHandler) DeleteHallRoles(c *gin.Context) {

}

// ROLE PERMISSIONS

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

	c.JSON(http.StatusOK, res)

}

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

	c.JSON(http.StatusOK, res)
}

// INVITES MANAGEMENT

func (h *HallHandler) GetCurrentInviteLinks(c *gin.Context) {

}

func (h *HallHandler) CreateNewInviteLink(c *gin.Context) {

}

func (h *HallHandler) InvokeInviteLink(c *gin.Context) {

}

// JOIN REQUEST MANAGEMENT

func (h *HallHandler) GetCurrentRequests(c *gin.Context) {

}

func (h *HallHandler) CreateJoinRequest(c *gin.Context) {

}

func (h *HallHandler) AcceptJoinRequest(c *gin.Context) {

}

func (h *HallHandler) DeclineJoinRequest(c *gin.Context) {

}

// BANS
func (h *HallHandler) GetBannedUsers(c *gin.Context) {

}

func (h *HallHandler) BanAnUser(c *gin.Context) {

}

func (h *HallHandler) UnbanUser(c *gin.Context) {

}
