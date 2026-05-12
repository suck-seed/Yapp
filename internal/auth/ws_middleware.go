// internal/auth/ws_middleware.go
package auth

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// WebSocketAuthMiddleware allows browser WebSocket auth using ?token=
// while still reusing your normal AuthMiddleware logic.
func WebSocketAuthMiddleware() gin.HandlerFunc {
	baseAuth := AuthMiddleware()

	return func(c *gin.Context) {
		token := strings.TrimSpace(c.Query("token"))

		// Optional alias if you ever use ?access_token=
		if token == "" {
			token = strings.TrimSpace(c.Query("access_token"))
		}

		// Browser WebSocket cannot send Authorization header,
		// so we convert ?token= into Authorization header
		// before your existing AuthMiddleware runs.
		if token != "" && c.GetHeader("Authorization") == "" {
			c.Request.Header.Set("Authorization", "Bearer "+token)
		}

		baseAuth(c)
	}
}
