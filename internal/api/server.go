package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/suck-seed/yapp/config"
)

// StartServer
func StartServer(cfg config.AppConfig) {

	// new app handler
	router := gin.Default()

	// use CORS setting
	router.Use(cfg.CORS)

	// listen to /ping
	router.GET("/ping", func(c *gin.Context) {

		// return JSON
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// runs on 8080 on default
	err := router.Run("localhost:" + cfg.ServerPort)
	if err != nil {
		return
	}
}
