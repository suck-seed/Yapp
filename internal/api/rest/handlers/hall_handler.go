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

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Hall created successfully",
		"data":    res,
	})

}
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

// SINGLE HALL RUD
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

// Depriciated not needed
// func (h *HallHandler) UpdateCurrentHall(c *gin.Context) {

// }

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

//----------------SETTING SCOPE

// PROFILE MANAGEMENT
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

// MEMBERS MANAGEMENT
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

// ROLE MANAGEMENT
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

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Role permissions retrieved successfully",
		"data":    res,
	})

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

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Role permissions updated successfully",
		"data":    res,
	})
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

// BANS
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
