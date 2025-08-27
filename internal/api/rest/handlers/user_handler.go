package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/suck-seed/yapp/internal/dto"
	"github.com/suck-seed/yapp/internal/services"
)

// UserHandler : CLASS
type UserHandler struct {
	// inject IUserService
	userService services.IUserService
}

func NewUserHandler(uSvc services.IUserService) *UserHandler {
	return &UserHandler{
		userService: uSvc,
	}
}

func (h *UserHandler) Hello(c *gin.Context) {
	if c.Request.Method != http.MethodGet {
		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"error": "Method not allowed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Hello Everynyan",
	})
}

// CreateUser : Functions / Methods accessed by UserHandler
func (h *UserHandler) Register(c *gin.Context) {

	if c.Request.Method != http.MethodPost {
		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"error": "Method not allowed",
		})
		return
	}

	user := dto.UserSignup{}
	err := c.BindJSON(&user)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Bad Request",
			"message": "Invalid JSON",
		})

		return
	}

	// Proceed to service Handlers
	token, err := h.userService.RegisterUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": token,
	})
}

func (h *UserHandler) GetUser(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{"id": c.Param("id")})
}
