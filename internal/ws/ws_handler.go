package ws

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/suck-seed/yapp/internal/models"
)

type WsHandler struct {
	hub *models.Hub
}

func NewWsHandler(h *models.Hub) *WsHandler {
	return &WsHandler{
		hub: h,
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {

		// get origin of frontend
		// origin := r.Header.Get("Origin")
		// return origin == config.FrontEndOrigin()
		return true
	},
}

// JoinRoom :
