package auth

import (
	"context"
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

		ctx := context.WithValue(c.Request.Context(), CtxUserIDKey, claims.ID)
		ctx = context.WithValue(ctx, CtxUsernameKey, claims.Username)
		c.Request = c.Request.WithContext(ctx)

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

// ---- UserID and Username from gin context ----
func CurrentUserFromGinContext(c *gin.Context) (string, string, error) {

	rawId, ok := c.Get(CtxUserIDKey)
	if !ok {
		return "", "", utils.ErrorNoUserIdInContext
	}

	idString, _ := rawId.(string)
	if idString == "" {
		return "", "", utils.ErrorEmptyUserIdInContext
	}

	rawUsername, _ := c.Get(CtxUsernameKey)
	usernameString, _ := rawUsername.(string)

	return idString, usernameString, nil

}

// ---- UserID and Username from context ----
func CurrentUserFromContext(c context.Context) (id string, username string, err error) {

	rawId := c.Value(CtxUserIDKey)
	if rawId == nil {
		return "", "", utils.ErrorNoUserIdInContext
	}

	idString, _ := rawId.(string)
	if idString == "" {
		return "", "", utils.ErrorEmptyUserIdInContext

	}

	usernameString, _ := c.Value(CtxUsernameKey).(string)

	return idString, usernameString, nil

}
