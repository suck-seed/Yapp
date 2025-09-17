package auth

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/suck-seed/yapp/internal/utils"
)

const (
	CtxUserIDKey   = "user_id"
	CtxUsernameKey = "username"
)

// Verifies JWT from cookie "jwt" or "Authorization : Bearer <token>"
// and injects userId/username into gin.Context
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		tokenString := GetTokenFromRequest(c)
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
		c.Set(CtxUserIDKey, claims.ID)
		c.Set(CtxUsernameKey, claims.Username)
		c.Next()
	}

}

func GetTokenFromRequest(c *gin.Context) string {

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

func GetUsernameAndIdFromContext(c *gin.Context) (userId string, username string, err error) {
	rawId, ok := c.Get(CtxUserIDKey) // Fix typo here
	if !ok {
		return "", "", utils.ErrorNoUserIdInContext
	}
	rawUsername, _ := c.Get(CtxUsernameKey)

	userId, _ = rawId.(string)
	username, _ = rawUsername.(string)

	if userId == "" {
		return "", "", utils.ErrorEmptyUserIdInContext
	}
	return userId, username, nil
}
