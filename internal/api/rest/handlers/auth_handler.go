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

// Signup godoc
// @Summary      Register a new account
// @Description  Creates a new user account and returns the created user.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      dto.SignupUserReq          true  "Signup payload"
// @Success      201   {object}  map[string]interface{}  "Account created"
// @Failure      400   {object}  map[string]interface{}  "Invalid input"
// @Failure      409   {object}  map[string]interface{}  "Username / e-mail already taken"
// @Router       /auth/signup [post]
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
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Account created successfully",
		"data":    res,
	})
}

// Signin godoc
// @Summary      Sign in
// @Description  Authenticates a user and sets an HttpOnly `jwt` cookie on success.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      dto.SigninUserReq          true  "Signin payload"
// @Success      200   {object}  map[string]interface{}  "Signed in — jwt cookie is set"
// @Failure      400   {object}  map[string]interface{}  "Invalid input"
// @Failure      401   {object}  map[string]interface{}  "Wrong credentials"
// @Router       /auth/signin [post]
func (h *AuthHandler) Signin(c *gin.Context) {
	userSignIn := &dto.SigninUserReq{}
	if err := c.ShouldBindJSON(userSignIn); err != nil {
		utils.WriteError(c, utils.ErrorInvalidInput)
		return
	}

	signInRes, err := h.IUserService.Signin(c.Request.Context(), userSignIn)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	// set cookie
	const cookieSeconds = 24 * 60 * 60

	isHTTPS := c.GetHeader("X-Forwarded-Proto") == "https"
	if isHTTPS {
		c.SetSameSite(http.SameSiteNoneMode)
		c.SetCookie("jwt", signInRes.AccessToken, cookieSeconds, "/", "", true, true)
	} else {
		// local dev over plain http
		c.SetSameSite(http.SameSiteLaxMode)
		c.SetCookie("jwt", signInRes.AccessToken, cookieSeconds, "/", "", false, true)
	}

	// filtered response (not sending accesstoken over https, so removed it)
	res := &dto.SigninUserRes{
		UserMe:      signInRes.UserMe,
		Success:     signInRes.Success,
		AccessToken: signInRes.AccessToken,
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Signed in successfully",
		"data":    res,
	})
}

// Signout godoc
// @Summary      Sign out
// @Description  Clears the `jwt` cookie, effectively ending the session.
// @Tags         auth
// @Produce      json
// @Security     CookieAuth
// @Success      200  {object}  map[string]interface{}  "Signed out"
// @Router       /auth/signout [post]
func (h *AuthHandler) Signout(c *gin.Context) {

	isHTTPS := c.GetHeader("X-Forwarded-Proto") == "https"

	if isHTTPS {
		c.SetSameSite(http.SameSiteNoneMode)
		c.SetCookie("jwt", "", -1, "/", "", true, true)
	} else {
		c.SetSameSite(http.SameSiteLaxMode)
		c.SetCookie("jwt", "", -1, "/", "", false, true)
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Signed out successfully",
		"data":    nil,
	})
}
