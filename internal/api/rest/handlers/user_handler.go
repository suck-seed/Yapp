package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/suck-seed/yapp/internal/auth"
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

func (h *UserHandler) GetUser(c *gin.Context) {
	// Extract user info from context (already validated by middleware)
	userId, username, err := auth.GetUsernameAndIdFromContext(c)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	userIdParsed, err := utils.ParseUUID(userId)
	if err != nil {
		utils.WriteError(c, utils.ErrorInvalidUserIdInContext)
		return
	}

	user, err := h.IUserService.GetUserByID(c.Request.Context(), userIdParsed)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	// Return user information
	c.JSON(http.StatusOK, gin.H{
		"username":    username,
		"displayName": user.DisplayName,
		"email":       user.Email,
		"avatarUrl":   user.AvatarURL,
		"active":      user.Active,
		"success":     true,
	})
}
