package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	dto "github.com/suck-seed/yapp/internal/dto/hall"
	"github.com/suck-seed/yapp/internal/services"
	"github.com/suck-seed/yapp/internal/utils"
)

type HallHandler struct {
	services.IHallService
}

func NewHallHandler(hallService services.IHallService) *HallHandler {
	return &HallHandler{
		hallService,
	}
}

// TOP LEVEL HALL OPERATIONS
func (h *HallHandler) CreateHall(c *gin.Context) {

	u := &dto.CreateHallReq{}

	if err := c.ShouldBindJSON(u); err != nil {
		utils.WriteError(c, utils.ErrorInvalidInput)
		return
	}

	res, err := h.IHallService.CreateHall(c.Request.Context(), u)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	c.JSON(http.StatusCreated, res)

}

func (h *HallHandler) GetUserHalls(c *gin.Context) {

	res, err := h.IHallService.GetUserHalls(c.Request.Context())
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

// SINGLE HALL RUD
func (h *HallHandler) GetCurrentHall(c *gin.Context) {

}

func (h *HallHandler) UpdateCurrentHall(c *gin.Context) {

}

func (h *HallHandler) DeleteCurrentHall(c *gin.Context) {

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

func (h *HallHandler) UpdateHallMember(c *gin.Context) {

}

func (h *HallHandler) RemoveHallMember(c *gin.Context) {

}

// ROLE MANAGEMENT
func (h *HallHandler) GetHallRoles(c *gin.Context) {

}

func (h *HallHandler) CreateHallRoles(c *gin.Context) {

}

func (h *HallHandler) UpdateHallRoles(c *gin.Context) {

}

func (h *HallHandler) DeleteHallRoles(c *gin.Context) {

}

// ROLE PERMISSIONS

func (h *HallHandler) GetRolesPermissions(c *gin.Context) {

}

func (h *HallHandler) UpdateRolesPermissions(c *gin.Context) {

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

	// Parse room_id
	hallIdStr := c.Param("hall_id")
	hallID, err := uuid.Parse(hallIdStr)
	if err != nil {
		utils.WriteError(c, utils.ErrorInvalidRoomIDFormat)
		return
	}

}

func (h *HallHandler) BanAnUser(c *gin.Context) {

}

func (h *HallHandler) UnbanUser(c *gin.Context) {

}
