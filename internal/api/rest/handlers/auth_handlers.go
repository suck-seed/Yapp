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
	services.IUserService
}

func NewAuthHandler(userService services.IUserService) *AuthHandler {
	return &AuthHandler{
		userService,
	}
}

func (h *AuthHandler) CreateUser(c *gin.Context) {

	var u dto.CreateUserReq

	if err := c.ShouldBindJSON(&u); err != nil {

		utils.WriteError(c, utils.ErrorInvalidInput)
		return

	}

	res, err := h.IUserService.CreateUser(c.Request.Context(), &u)
	if err != nil {

		utils.WriteError(c, err)
		return
	}

	c.JSON(http.StatusCreated, res)

}

func (h *AuthHandler) Login(c *gin.Context) {

	var user dto.LoginUserReq

	if err := c.ShouldBindJSON(&user); err != nil {

		utils.WriteError(c, utils.ErrorInvalidInput)
		return
	}

	u, err := h.IUserService.Login(c.Request.Context(), &user)
	if err != nil {

		utils.WriteError(c, err)
		return

	}

	// setcookie
	const cookieSecond = 24 * 60 * 60

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("jwt", u.AccessToken, cookieSecond, "/", "localhost", false, true)

	// a filtered req as we do not want to implicitely pass accessToken to client
	res := &dto.LoginUserRes{
		UserId:   u.UserId,
		Username: u.Username,
	}

	c.JSON(http.StatusOK, res)

}

func (h *AuthHandler) Logout(c *gin.Context) {

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("jwt", "", -1, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "Logout successful",
	})
}
