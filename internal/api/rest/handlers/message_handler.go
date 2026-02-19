package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/suck-seed/yapp/internal/auth"
	dto "github.com/suck-seed/yapp/internal/dto/message"
	"github.com/suck-seed/yapp/internal/services"
	"github.com/suck-seed/yapp/internal/utils"
)

type MessageHandler struct {
	services.IMessageService
}

func NewMessageHandler(messageService services.IMessageService) *MessageHandler {

	return &MessageHandler{
		messageService,
	}
}

func (h *MessageHandler) FetchMessage(c *gin.Context) {

	// a dto
	u := &dto.MessageQueryParams{}

	if err := c.ShouldBindJSON(u); err != nil {
		utils.WriteError(c, utils.ErrorInvalidInput)
		return
	}

	userInfo, err := auth.CurrentUserFromGinContext(c)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	// calling server layer
	res, err := h.IMessageService.FetchMessages(c.Request.Context(), userInfo, u)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)

}
