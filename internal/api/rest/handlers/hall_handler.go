package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/suck-seed/yapp/internal/dto"
	"github.com/suck-seed/yapp/internal/services"
	"github.com/suck-seed/yapp/internal/utils"
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

func (h *HallHandler) CreateHall(c *gin.Context) {

	var u dto.CreateHallReq

	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	res, err := h.IHallService.CreateHall(c.Request.Context(), &u)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	c.JSON(http.StatusCreated, res)
}
