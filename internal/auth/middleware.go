package auth

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/suck-seed/yapp/internal/utils"
)

const (
	CtxUserIDKey   = "user_id"
	CtxUsernameKey = "username"
)

// UserInfo hold authenticated user information
type UserInfo struct {
	ID       uuid.UUID
	Username string
}

// Verifies JWT from cookie "jwt" or "Authorization : Bearer <token>"
// and injects userId/username into gin.Context
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		token, ok := GetTokenFromRequest(c)
		if !ok {
			utils.WriteError(c, utils.ErrorMissingToken)
			c.Abort()
			return
		}

		// Parse claims from token
		claims, err := ParseAndVerify(token)
		if err != nil {
			utils.WriteError(c, utils.ErrorInvalidToken)
			c.Abort()
			return
		}

		// Parsing UUID for the userID
		userID, err := uuid.Parse(claims.ID)
		if err != nil {
			utils.WriteError(c, utils.ErrorInvalidUserUUID)
			c.Abort()
			return
		}

		userInfo := &UserInfo{
			ID:       userID,
			Username: claims.Username,
		}

		// store in gin.Context
		c.Set(CtxUserIDKey, userInfo.ID)
		c.Set(CtxUsernameKey, userInfo.Username)

		// context.context
		ctx := context.WithValue(c.Request.Context(), CtxUserIDKey, userInfo.ID)
		ctx = context.WithValue(ctx, CtxUsernameKey, userInfo.Username)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}

}

func GetTokenFromRequest(c *gin.Context) (string, bool) {
	// Trying cookie
	if cookie, err := c.Cookie("jwt"); err == nil && cookie != "" {
		return cookie, true
	}

	// Trying Authorization header
	if token, ok := strings.CutPrefix(c.GetHeader("Authorization"), "Bearer "); ok {
		token := strings.TrimSpace(token)
		if token != "" {
			return token, true
		}
	}

	return "", false
}

func CurrentUserFromGinContext(c *gin.Context) (*UserInfo, error) {
	rawId, ok := c.Get(CtxUserIDKey)
	if !ok {
		return nil, utils.ErrorNoUserIdInContext
	}

	userID, ok := rawId.(uuid.UUID)
	if !ok {
		return nil, utils.ErrorInvalidUserUUID
	}

	rawUsername, _ := c.Get(CtxUsernameKey)
	username, _ := rawUsername.(string)

	return &UserInfo{
		ID:       userID,
		Username: username,
	}, nil

}
