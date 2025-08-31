package ws

import "github.com/suck-seed/yapp/internal/models"

func NewHub() *models.Hub {
	return &models.Hub{
		Rooms: make(map[string]*models.Room),
	}
}
