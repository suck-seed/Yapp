package auth

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/suck-seed/yapp/internal/utils"
)

const (
	CtxUSerIDKey   = "user_id"
	CtxUsernameKey = "username"
)

// Verifies JWT from cookie "jwt" or "Authorization : Bearer <token>"
// and injects userId/username into gin.Context
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		tokenString := getTokenFromRequest(c)
		if tokenString == "" {
			utils.WriteError(c, utils.ErrorMissingToken)
			return
		}

		// Parse claims from token
		claims, err := ParseAndVerify(tokenString)
		if err != nil {
			utils.WriteError(c, utils.ErrorInvalidToken)
			return

		}

		// Everything's alright, place the context for handlers
		c.Set(CtxUSerIDKey, claims.ID)
		c.Set(CtxUsernameKey, claims.Username)
		c.Next()
	}

}

func getTokenFromRequest(c *gin.Context) string {

	// Trying cookie
	if cookie, err := c.Cookie("jwt"); err == nil && cookie != "" {
		return cookie
	}

	// Trying Authorization header
	authHeader := c.GetHeader("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
	}

	return ""
}

func CurrentUserFromContext(c *gin.Context) (userId string, username string, err error) {
	rawId, ok := c.Get(CtxUSerIDKey)
	if !ok {
		return "", "", utils.ErrorNoUserIdInContext
	}

	rawUsername, _ := c.Get(CtxUsernameKey)

	idString, _ := rawId.(string)
	usernameString, _ := rawUsername.(string)

	if idString == "" {
		return "", "", utils.ErrorEmptyUserIdInContext
	}

	return idString, usernameString, nil
}
