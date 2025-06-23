package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {

	// new mux handler
	r := gin.Default()

	// CORS middleware use
	r.Use(CORSMiddleware())

	// listen to /ping
	r.GET("/ping", func(c *gin.Context) {

		// return JSON
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// runs on 8080 on default
	r.Run()
}

// CORS middleware
func CORSMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Max")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}
