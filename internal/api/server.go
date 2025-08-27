package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/suck-seed/yapp/config"
	"github.com/suck-seed/yapp/internal/api/rest"
	"github.com/suck-seed/yapp/internal/services"
)

// StartServer : Runs up an instance of server, Pass dependency to individual Handlers, Middle ware implemented here
func StartServer(cfg config.AppConfig) {

	// new app handler
	router := gin.Default()

	// use CORS setting
	router.Use(cfg.CORS)

	// get services
	userService := services.NewUserService()

	// rest routes with services injected, can pass cfg too
	rest.RegisterUserRoutes(router, userService)
	//TODO Add similar routers for other too

	// runs on 8080 on default
	err := router.Run(":" + cfg.ServerPort)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
}
