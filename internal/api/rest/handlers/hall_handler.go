package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/suck-seed/yapp/internal/services"
)

type HallHandler struct {
	// inject IHallService
	services.IHallService
}

func NewHallHandler(hallService services.IHallService) *HallHandler {
	return &HallHandler{
		hallService,
	}
}

func (h *HallHandler) Ping(c *gin.Context) {

	c.JSON(200, gin.H{
		"message": "ping",
	})
}
