package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/suck-seed/yapp/internal/dto"
	"github.com/suck-seed/yapp/internal/services"
	"github.com/suck-seed/yapp/internal/utils"
)

type HallHandler struct {
	services.IHallService
}

func NewHallhandler(hallService services.IHallService) *HallHandler {
	return &HallHandler{
		hallService,
	}
}

func (h *HallHandler) CreateHall(c *gin.Context) {

	u := &dto.CreateHallReq{}

	if err := c.ShouldBindJSON(u); err != nil {
		utils.WriteError(c, utils.ErrorInvalidInput)
		return
	}

	res, err := h.IHallService.CreateHall(c.Request.Context(), u)
	if err != nil {
		utils.WriteError(c, err)
	}

	c.JSON(http.StatusCreated, res)

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
