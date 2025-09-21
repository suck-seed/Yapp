package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/suck-seed/yapp/internal/auth"
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

func (h *UserHandler) GetUser(c *gin.Context) {
	// Extract user info from context (already validated by middleware)
	userId, err := auth.GetUserIDFromContext(c)
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

	res := dto.ToUserMe(*user)

	// Return user information
	c.JSON(http.StatusOK, res)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	var u dto.UpdateProfileReq

	// Extract username from context (already validated by middleware)
	username, err := auth.GetUsernameFromContext(c)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	// Bind JSON payload: expects { "display_name": "...", "avatar_url": "..." }
	if err := c.ShouldBindJSON(&u); err != nil {
		utils.WriteError(c, utils.ErrorInvalidInput)
		return
	}

	// Call service with values from payload
	user, err := h.IUserService.UpdateUser(c.Request.Context(), username, u.DisplayName, u.AvatarURL)
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
