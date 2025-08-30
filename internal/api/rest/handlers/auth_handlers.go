package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/suck-seed/yapp/internal/dto"
	"github.com/suck-seed/yapp/internal/services"
	"github.com/suck-seed/yapp/internal/utils"
)

type AuthHandler struct {
	// inject IUserService
	userService services.IUserService
}

func NewAuthHandler(userService services.IUserService) *AuthHandler {
	return &AuthHandler{
		userService: userService,
	}
}

func (h *AuthHandler) CreateUser(c *gin.Context) {

	var u dto.CreateUserReq

	if err := c.ShouldBindJSON(&u); err != nil {

		utils.WriteError(c, utils.ErrorInvalidInput)
		return

	}

	res, err := h.userService.CreateUser(c.Request.Context(), &u)
	if err != nil {

		utils.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)

}

func (h *AuthHandler) Login(c *gin.Context) {

	var user dto.LoginUserReq

	if err := c.ShouldBindJSON(&user); err != nil {

		utils.WriteError(c, utils.ErrorInvalidInput)
		return
	}

	u, err := h.userService.Login(c.Request.Context(), &user)
	if err != nil {

		utils.WriteError(c, err)
		return

	}

	// setcookie
	c.SetCookie("jwt", u.AccessToken, 3600, "/", "localhost", false, true)

	// a filtered req as we do not want to implicitely pass accessToken to client
	res := &dto.LoginUserRes{
		ID:       u.ID,
		Username: u.Username,
	}

	c.JSON(http.StatusOK, res)

}

func (h *AuthHandler) Logout(c *gin.Context) {
	c.SetCookie("jwt", "", -1, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"message": "Logout successful",
	})
}
