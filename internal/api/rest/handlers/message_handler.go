package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/suck-seed/yapp/internal/auth"
	dto "github.com/suck-seed/yapp/internal/dto/message"
	"github.com/suck-seed/yapp/internal/services"
	"github.com/suck-seed/yapp/internal/utils"
)

type MessageHandler struct {
	services.IMessageService
}

func NewMessageHandler(messageService services.IMessageService) *MessageHandler {
	return &MessageHandler{messageService}
}

// ── GET /rooms/:roomID/messages ───────────────────────────────────────────────

func (h *MessageHandler) FetchMessages(c *gin.Context) {
	userInfo, err := auth.CurrentUserFromGinContext(c)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	roomID, err := uuid.Parse(c.Param("roomID"))
	if err != nil {
		utils.WriteError(c, utils.ErrorInvalidIDFormart)
		return
	}

	query := &dto.FetchMessagesQuery{}

	limitStr := c.DefaultQuery("limit", "50")
	query.Limit, err = strconv.Atoi(limitStr)
	if err != nil || query.Limit <= 0 {
		query.Limit = 50
	}
	if query.Limit > 100 {
		query.Limit = 100
	}

	if beforeStr := c.Query("before"); beforeStr != "" {
		id, err := uuid.Parse(beforeStr)
		if err != nil {
			utils.WriteError(c, utils.ErrorInvalidInput)
			return
		}
		query.Before = &id
	}

	if afterStr := c.Query("after"); afterStr != "" {
		id, err := uuid.Parse(afterStr)
		if err != nil {
			utils.WriteError(c, utils.ErrorInvalidInput)
			return
		}
		query.After = &id
	}

	if aroundStr := c.Query("around"); aroundStr != "" {
		id, err := uuid.Parse(aroundStr)
		if err != nil {
			utils.WriteError(c, utils.ErrorInvalidInput)
			return
		}
		query.Around = &id
	}

	res, err := h.IMessageService.FetchMessages(c.Request.Context(), userInfo, roomID, query)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Messages retrieved successfully",
		"data":    res,
	})
}

// ── GET /rooms/:roomID/messages/:messageID ────────────────────────────────────

func (h *MessageHandler) GetMessage(c *gin.Context) {
	userInfo, err := auth.CurrentUserFromGinContext(c)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	roomID, err := uuid.Parse(c.Param("roomID"))
	if err != nil {
		utils.WriteError(c, utils.ErrorInvalidIDFormart)
		return
	}

	messageID, err := uuid.Parse(c.Param("messageID"))
	if err != nil {
		utils.WriteError(c, utils.ErrorInvalidIDFormart)
		return
	}

	res, err := h.IMessageService.GetMessage(c.Request.Context(), userInfo, roomID, messageID)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Message retrieved successfully",
		"data":    res,
	})
}

// ── PATCH /rooms/:roomID/messages/:messageID ──────────────────────────────────

func (h *MessageHandler) UpdateMessage(c *gin.Context) {
	userInfo, err := auth.CurrentUserFromGinContext(c)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	roomID, err := uuid.Parse(c.Param("roomID"))
	if err != nil {
		utils.WriteError(c, utils.ErrorInvalidIDFormart)
		return
	}

	messageID, err := uuid.Parse(c.Param("messageID"))
	if err != nil {
		utils.WriteError(c, utils.ErrorInvalidIDFormart)
		return
	}

	var req dto.UpdateMessageReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.WriteError(c, utils.ErrorInvalidInput)
		return
	}

	res, err := h.IMessageService.UpdateMessage(c.Request.Context(), userInfo, roomID, messageID, &req)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Message updated successfully",
		"data":    res,
	})
}

// ── DELETE /rooms/:roomID/messages/:messageID ─────────────────────────────────

func (h *MessageHandler) DeleteMessage(c *gin.Context) {
	userInfo, err := auth.CurrentUserFromGinContext(c)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	roomID, err := uuid.Parse(c.Param("roomID"))
	if err != nil {
		utils.WriteError(c, utils.ErrorInvalidIDFormart)
		return
	}

	messageID, err := uuid.Parse(c.Param("messageID"))
	if err != nil {
		utils.WriteError(c, utils.ErrorInvalidIDFormart)
		return
	}

	if err := h.IMessageService.DeleteMessage(c.Request.Context(), userInfo, roomID, messageID); err != nil {
		utils.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Message deleted successfully",
		"data":    nil,
	})
}

// ── PUT /rooms/:roomID/messages/:messageID/reactions/:emoji ───────────────────
// Adds a reaction. If the same user+emoji already exists, returns 200 with no change (idempotent).
// PUT /rooms/{roomID}/messages/{messageID}/reactions/%F0%9F%98%8C
// That %F0%9F%98%8C is the URL-encoded form of 😌.
func (h *MessageHandler) AddReaction(c *gin.Context) {
	userInfo, err := auth.CurrentUserFromGinContext(c)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	roomID, err := uuid.Parse(c.Param("roomID"))
	if err != nil {
		utils.WriteError(c, utils.ErrorInvalidIDFormart)
		return
	}

	messageID, err := uuid.Parse(c.Param("messageID"))
	if err != nil {
		utils.WriteError(c, utils.ErrorInvalidIDFormart)
		return
	}

	// Gin automatically URL-decodes path params, so emoji arrives as the actual character
	emoji := c.Param("emoji")
	if emoji == "" {
		utils.WriteError(c, utils.ErrorInvalidInput)
		return
	}

	res, err := h.IMessageService.AddReaction(c.Request.Context(), userInfo, roomID, messageID, emoji)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Reaction added",
		"data":    res,
	})
}

// ── DELETE /rooms/:roomID/messages/:messageID/reactions/:emoji ────────────────

func (h *MessageHandler) RemoveReaction(c *gin.Context) {
	userInfo, err := auth.CurrentUserFromGinContext(c)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	roomID, err := uuid.Parse(c.Param("roomID"))
	if err != nil {
		utils.WriteError(c, utils.ErrorInvalidIDFormart)
		return
	}

	messageID, err := uuid.Parse(c.Param("messageID"))
	if err != nil {
		utils.WriteError(c, utils.ErrorInvalidIDFormart)
		return
	}

	emoji := c.Param("emoji")
	if emoji == "" {
		utils.WriteError(c, utils.ErrorInvalidInput)
		return
	}

	if err := h.IMessageService.RemoveReaction(c.Request.Context(), userInfo, roomID, messageID, emoji); err != nil {
		utils.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Reaction removed",
		"data":    nil,
	})
}
