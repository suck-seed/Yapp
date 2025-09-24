package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/suck-seed/yapp/internal/dto"
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

func (userHandler *UserHandler) Ping(c *gin.Context) {

	c.JSON(200, gin.H{
		"message": "ping",
	})
}

func (h *UserHandler) GetUserMe(c *gin.Context) {

	user, err := h.IUserService.GetUserMe(c.Request.Context())
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	res := dto.ToUserMe(*user)

	// Return user information
	c.JSON(http.StatusOK, res)
}

func (h *UserHandler) UpdateUserMe(c *gin.Context) {
	var u dto.UpdateUserMeReq

	// Bind JSON payload: expects { "display_name": "...", "avatar_url": "..." }
	if err := c.ShouldBindJSON(&u); err != nil {
		utils.WriteError(c, utils.ErrorInvalidInput)
		return
	}

	// Call service with values from payload
	user, err := h.IUserService.UpdateUserMe(c.Request.Context(), &u)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	// Return updated user info
	c.JSON(http.StatusOK, gin.H{
		"username":     user.Username,
		"display_name": user.DisplayName,
		"email":        user.Email,
		"avatar_url":   user.AvatarURL,
		"active":       user.Active,
	})
}
