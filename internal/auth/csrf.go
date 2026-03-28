package auth

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"
)

func generateCSRFToken() string {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}

// Sets csrf_token cookie if missing
func CSRFCookieMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := c.Cookie("csrf_token")
		if err != nil {
			isHTTPS := c.GetHeader("X-Forwarded-Proto") == "https"

			token := generateCSRFToken()

			if isHTTPS {
				c.SetSameSite(http.SameSiteNoneMode)
				c.SetCookie("csrf_token", token, 24*60*60, "/", "", true, false)
			} else {
				c.SetSameSite(http.SameSiteLaxMode)
				c.SetCookie("csrf_token", token, 24*60*60, "/", "", false, false)
			}
		}

		c.Next()
	}
}

// Enforce CSRF on unsafe methods
func CSRFMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		switch c.Request.Method {
		case http.MethodGet, http.MethodHead, http.MethodOptions:
			c.Next()
			return
		}

		cookieToken, err := c.Cookie("csrf_token")
		if err != nil || cookieToken == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "missing csrf cookie",
			})
			return
		}

		headerToken := c.GetHeader("X-CSRF-Token")
		if headerToken == "" || headerToken != cookieToken {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "invalid csrf token",
			})
			return
		}

		c.Next()
	}
}
