package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/suck-seed/yapp/internal/auth"
	dto "github.com/suck-seed/yapp/internal/dto/user"
	"github.com/suck-seed/yapp/internal/services"
	"github.com/suck-seed/yapp/internal/utils"
)

// UserHandler : CLASS
type UserHandler struct {
	// inject IUserService
	services.IUserService
}

func NewUserHandler(userService services.IUserService) *UserHandler {
	return &UserHandler{
		userService,
	}
}

func (h *UserHandler) Ping(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Ping successful",
		"data":    nil,
	})
}

func (h *UserHandler) GetUserMe(c *gin.Context) {

	userInfo, err := auth.CurrentUserFromGinContext(c)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	user, err := h.IUserService.GetUserMe(c.Request.Context(), userInfo)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	res := dto.ToUserMe(*user)

	// Return user information
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Profile retrieved successfully",
		"data":    res,
	})
}

func (h *UserHandler) UpdateUserMe(c *gin.Context) {
	u := &dto.UpdateUserMeReq{}

	// Bind JSON payload: expects { "display_name": "...", "avatar_url": "..." }
	if err := c.ShouldBindJSON(u); err != nil {
		utils.WriteError(c, utils.ErrorInvalidInput)
		return
	}

	userInfo, err := auth.CurrentUserFromGinContext(c)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	// Call service with values from payload
	res, err := h.IUserService.UpdateUserMe(c.Request.Context(), userInfo, u)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Profile updated successfully",
		"data":    res,
	})
}
