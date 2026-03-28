package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	dto "github.com/suck-seed/yapp/internal/dto/user"
	"github.com/suck-seed/yapp/internal/services"
	"github.com/suck-seed/yapp/internal/utils"
)

type AuthHandler struct {
	services.IUserService
}

func NewAuthHandler(userService services.IUserService) *AuthHandler {
	return &AuthHandler{userService}
}

func (h *AuthHandler) Signup(c *gin.Context) {
	u := &dto.SignupUserReq{}
	if err := c.ShouldBindJSON(u); err != nil {
		utils.WriteError(c, utils.ErrorInvalidInput)
		return
	}
	res, err := h.IUserService.Signup(c.Request.Context(), u)
	if err != nil {
		utils.WriteError(c, err)
		return
	}
	c.JSON(http.StatusCreated, res)
}

func (h *AuthHandler) Signin(c *gin.Context) {
	userSignIn := &dto.SigninUserReq{}
	if err := c.ShouldBindJSON(userSignIn); err != nil {
		utils.WriteError(c, utils.ErrorInvalidInput)
		return
	}

	SignInRes, err := h.IUserService.Signin(c.Request.Context(), userSignIn)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	const cookieSecond = 24 * 60 * 60
	// CRITICAL: SetSameSite MUST be called before SetCookie
	c.SetSameSite(http.SameSiteNoneMode)
	// Secure=true and HttpOnly=true are required for cross-site cookies
	c.SetCookie("jwt", SignInRes.AccessToken, cookieSecond, "/", "", true, true)

	res := &dto.SigninUserRes{
		UserMe:  SignInRes.UserMe,
		Success: SignInRes.Success,
	}
	c.JSON(http.StatusOK, res)
}

func (h *AuthHandler) Signout(c *gin.Context) {
	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie("jwt", "", -1, "/", "", true, true)
	c.JSON(http.StatusOK, gin.H{"message": "Signed out successfully"})
}
