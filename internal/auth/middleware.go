package auth

import (
	"context"
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

		// Everything's alright, place the context for ws_handler and other handler
		c.Set(CtxUSerIDKey, claims.ID)
		c.Set(CtxUsernameKey, claims.Username)

		// Also add them to context.Context, to be accessed from service and repository layer if we have to
		ctx := context.WithValue(c.Request.Context(), CtxUSerIDKey, claims.ID)
		ctx = context.WithValue(ctx, CtxUsernameKey, claims.Username)
		c.Request = c.Request.WithContext(ctx)

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

func CurrentUserFromGinContext(c *gin.Context) (string, string, error) {
	rawId, ok := c.Get(CtxUSerIDKey)
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

func CurrentUserFromContext(c context.Context) (id string, username string, err error) {

	rawId := c.Value(CtxUSerIDKey)
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
