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

// FetchMessages godoc
// @Summary      Fetch messages in a room
// @Description  Returns a paginated list of messages. Use `before`, `after`, or `around` (message UUID) to paginate. Maximum 100 per request.
// @Tags         messages
// @Produce      json
// @Security     CookieAuth
// @Param        roomID  path      string  true   "Room ID (UUID)"
// @Param        limit   query     int     false  "Number of messages (1-100, default 50)"
// @Param        before  query     string  false  "Return messages older than this message ID"
// @Param        after   query     string  false  "Return messages newer than this message ID"
// @Param        around  query     string  false  "Return messages around this message ID"
// @Success      200     {object}  map[string]interface{}
// @Failure      400     {object}  map[string]interface{}
// @Failure      401     {object}  map[string]interface{}
// @Failure      403     {object}  map[string]interface{}
// @Router       /rooms/{roomID}/messages [get]
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

// GetMessage godoc
// @Summary      Get a single message
// @Description  Returns one message by ID.
// @Tags         messages
// @Produce      json
// @Security     CookieAuth
// @Param        roomID     path      string  true  "Room ID (UUID)"
// @Param        messageID  path      string  true  "Message ID (UUID)"
// @Success      200        {object}  map[string]interface{}
// @Failure      400        {object}  map[string]interface{}
// @Failure      401        {object}  map[string]interface{}
// @Failure      404        {object}  map[string]interface{}
// @Router       /rooms/{roomID}/messages/{messageID} [get]
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

// UpdateMessage godoc
// @Summary      Edit a message
// @Description  Updates the content of a message. Only the original author may edit.
// @Tags         messages
// @Accept       json
// @Produce      json
// @Security     CookieAuth
// @Param        roomID     path      string                 true  "Room ID (UUID)"
// @Param        messageID  path      string                 true  "Message ID (UUID)"
// @Param        body       body      dto.UpdateMessageReq   true  "New content"
// @Success      200        {object}  map[string]interface{}
// @Failure      400        {object}  map[string]interface{}
// @Failure      401        {object}  map[string]interface{}
// @Failure      403        {object}  map[string]interface{}
// @Router       /rooms/{roomID}/messages/{messageID} [patch]
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

// DeleteMessage godoc
// @Summary      Delete a message
// @Description  Soft-deletes a message. The author or a member with ManageMessages permission may delete.
// @Tags         messages
// @Produce      json
// @Security     CookieAuth
// @Param        roomID     path      string  true  "Room ID (UUID)"
// @Param        messageID  path      string  true  "Message ID (UUID)"
// @Success      200        {object}  map[string]interface{}
// @Failure      400        {object}  map[string]interface{}
// @Failure      401        {object}  map[string]interface{}
// @Failure      403        {object}  map[string]interface{}
// @Router       /rooms/{roomID}/messages/{messageID} [delete]
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

// AddReaction godoc
// @Summary      Add a reaction
// @Description  Adds an emoji reaction to a message. Idempotent — adding the same emoji twice is a no-op.
// @Tags         messages
// @Produce      json
// @Security     CookieAuth
// @Param        roomID     path      string  true  "Room ID (UUID)"
// @Param        messageID  path      string  true  "Message ID (UUID)"
// @Param        emoji      path      string  true  "URL-encoded emoji character, e.g. %F0%9F%98%8C for 😌"
// @Success      200        {object}  map[string]interface{}
// @Failure      400        {object}  map[string]interface{}
// @Failure      401        {object}  map[string]interface{}
// @Router       /rooms/{roomID}/messages/{messageID}/reactions/{emoji} [put]
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

// RemoveReaction godoc
// @Summary      Remove a reaction
// @Description  Removes the caller's emoji reaction from a message.
// @Tags         messages
// @Produce      json
// @Security     CookieAuth
// @Param        roomID     path      string  true  "Room ID (UUID)"
// @Param        messageID  path      string  true  "Message ID (UUID)"
// @Param        emoji      path      string  true  "URL-encoded emoji character"
// @Success      200        {object}  map[string]interface{}
// @Failure      400        {object}  map[string]interface{}
// @Failure      401        {object}  map[string]interface{}
// @Router       /rooms/{roomID}/messages/{messageID}/reactions/{emoji} [delete]
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
