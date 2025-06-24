package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/suck-seed/yapp/internal/services/user"
	"net/http"
)

// UserHandler : CLASS
type UserHandler struct {
	// inject IUserService
	userService user.IUserService
}

func NewUserHandler(uSvc user.IUserService) *UserHandler {
	return &UserHandler{
		userService: uSvc,
	}
}

// CreateUser : Functions / Methods accessed by UserHandler
func (h *UserHandler) CreateUser(c *gin.Context) {

	id := c.Param("id")
	returnUser, err := h.userService.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})

	}
	c.JSON(http.StatusOK, returnUser)
}

func (h *UserHandler) GetUser(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{"id": c.Param("id")})
}
