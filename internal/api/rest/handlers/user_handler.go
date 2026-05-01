package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/suck-seed/yapp/internal/auth"
	dto "github.com/suck-seed/yapp/internal/dto/user"
	"github.com/suck-seed/yapp/internal/models"
	"github.com/suck-seed/yapp/internal/services"
	"github.com/suck-seed/yapp/internal/utils"
)

type UserHandler struct {
	services.IUserService
}

func NewUserHandler(userService services.IUserService) *UserHandler {
	return &UserHandler{
		userService,
	}
}

func (h *UserHandler) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Ping successful",
		"data":    nil,
	})
}

func (h *UserHandler) GetUserMe(c *gin.Context) {
	userInfo, err := auth.CurrentUserFromGinContext(c)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	user, err := h.IUserService.GetUserMe(c.Request.Context(), userInfo)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Profile retrieved successfully",
		"data":    user,
	})
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	targetUserID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		utils.WriteError(c, utils.ErrorInvalidIDFormart)
		return
	}

	currentUserID := uuid.Nil
	userInfo, err := auth.CurrentUserFromGinContext(c)
	if err == nil {
		currentUserID = userInfo.ID
	}

	res, err := h.IUserService.GetUserPublic(c.Request.Context(), currentUserID, targetUserID)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User retrieved successfully",
		"data":    res,
	})
}

func (h *UserHandler) GetMutualFriends(c *gin.Context) {
	userInfo, err := auth.CurrentUserFromGinContext(c)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	targetUserID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		utils.WriteError(c, utils.ErrorInvalidIDFormart)
		return
	}

	res, err := h.IUserService.GetMutualFriends(c.Request.Context(), userInfo.ID, targetUserID)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Mutual friends retrieved successfully",
		"data":    res,
	})
}

func (h *UserHandler) UpdateUserMe(c *gin.Context) {
	u := &dto.UpdateUserMeReq{}

	if err := c.ShouldBindJSON(u); err != nil {
		utils.WriteError(c, utils.ErrorInvalidInput)
		return
	}

	userInfo, err := auth.CurrentUserFromGinContext(c)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	res, err := h.IUserService.UpdateUserMe(c.Request.Context(), userInfo, u)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Profile updated successfully",
		"data":    res,
	})
}

func (h *UserHandler) UpdateUsername(c *gin.Context) {
	var req dto.UpdateUsernameReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.WriteError(c, utils.ErrorInvalidInput)
		return
	}

	userInfo, err := auth.CurrentUserFromGinContext(c)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	res, err := h.IUserService.UpdateUsername(c.Request.Context(), userInfo, &req)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Username updated successfully",
		"data":    res,
	})
}

func (h *UserHandler) UpdateEmail(c *gin.Context) {
	var req dto.UpdateEmailReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.WriteError(c, utils.ErrorInvalidInput)
		return
	}

	userInfo, err := auth.CurrentUserFromGinContext(c)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	res, err := h.IUserService.UpdateEmail(c.Request.Context(), userInfo, &req)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Email updated successfully",
		"data":    res,
	})
}

func (h *UserHandler) DeleteMe(c *gin.Context) {
	userInfo, err := auth.CurrentUserFromGinContext(c)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	if err := h.IUserService.DeleteMe(c.Request.Context(), userInfo); err != nil {
		utils.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User deleted successfully",
		"data":    nil,
	})
}

func (h *UserHandler) SendFriendRequest(c *gin.Context) {
	var req dto.SendFriendRequestReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.WriteError(c, utils.ErrorInvalidInput)
		return
	}

	userInfo, err := auth.CurrentUserFromGinContext(c)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	res, err := h.IUserService.SendFriendRequest(c.Request.Context(), userInfo, &req)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Friend request sent successfully",
		"data":    res,
	})
}

func (h *UserHandler) RespondFriendRequest(c *gin.Context) {
	userInfo, err := auth.CurrentUserFromGinContext(c)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	requestID, err := uuid.Parse(c.Param("request_id"))
	if err != nil {
		utils.WriteError(c, utils.ErrorInvalidIDFormart)
		return
	}

	var req dto.RespondFriendRequestReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.WriteError(c, utils.ErrorInvalidInput)
		return
	}

	if err := h.IUserService.RespondFriendRequest(c.Request.Context(), userInfo, requestID, &req); err != nil {
		utils.WriteError(c, err)
		return
	}

	message := "Friend request declined successfully"
	if req.Action == "accept" {
		message = "Friend request accepted successfully"
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": message,
		"data":    nil,
	})
}

func (h *UserHandler) Unfriend(c *gin.Context) {
	userInfo, err := auth.CurrentUserFromGinContext(c)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	targetUserID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		utils.WriteError(c, utils.ErrorInvalidIDFormart)
		return
	}

	if err := h.IUserService.Unfriend(c.Request.Context(), userInfo, targetUserID); err != nil {
		utils.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User unfriended successfully",
		"data":    nil,
	})
}

func (h *UserHandler) GetMyFriends(c *gin.Context) {
	userInfo, err := auth.CurrentUserFromGinContext(c)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	res, err := h.IUserService.GetMyFriends(c.Request.Context(), userInfo)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Friends retrieved successfully",
		"data":    res,
	})
}

func (h *UserHandler) UpsertMyAppLink(c *gin.Context) {
	userInfo, err := auth.CurrentUserFromGinContext(c)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	var req dto.UpsertAppLinkReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.WriteError(c, utils.ErrorInvalidInput)
		return
	}

	res, err := h.IUserService.UpsertMyAppLink(c.Request.Context(), userInfo, &req)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "App link saved successfully",
		"data":    res,
	})
}

func (h *UserHandler) DeleteMyAppLink(c *gin.Context) {
	userInfo, err := auth.CurrentUserFromGinContext(c)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	provider := models.AppProvider(c.Param("provider"))
	if provider == "" {
		utils.WriteError(c, utils.ErrorInvalidInput)
		return
	}

	if err := h.IUserService.DeleteMyAppLink(c.Request.Context(), userInfo, provider); err != nil {
		utils.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "App link deleted successfully",
		"data":    nil,
	})
}
