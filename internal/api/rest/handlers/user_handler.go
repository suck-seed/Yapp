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

func NewUserHandler(userService services.IUserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) CreateUser(c *gin.Context) {

	// fetch json and bind into CreateUserReq

	var u dto.CreateUserReq

	if err := c.ShouldBindJSON(&u); err != nil {

		// bad request
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return

	}

	res, err := h.userService.CreateUser(c.Request.Context(), &u)
	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// no errors, so return
	c.JSON(http.StatusOK, res)

}
