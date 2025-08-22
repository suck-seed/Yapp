package websocket

import (
	"github.com/gorilla/websocket"
	"github.com/suck-seed/yapp/internal/websocket"
)

var (
	websocketUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)
