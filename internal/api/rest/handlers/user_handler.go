package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/suck-seed/yapp/internal/services"
)

// UserHandler : CLASS
type UserHandler struct {
	// inject IUserService
	userService services.IUserService
}

func NewUserHandler(userService services.IUserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (userHandler *UserHandler) Ping(c *gin.Context) {

	c.JSON(200, gin.H{
		"message": "ping",
	})
}
