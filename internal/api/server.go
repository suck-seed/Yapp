package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/suck-seed/yapp/config"
	"github.com/suck-seed/yapp/internal/api/rest"
	"github.com/suck-seed/yapp/internal/repositories"
	"github.com/suck-seed/yapp/internal/services"
)

// StartServer : Runs up an instance of server, Pass dependency to individual Handlers, Middle ware implemented here
func StartServer(cfg config.AppConfig) {

	// new app handler
	router := gin.Default()

	// use CORS setting
	router.Use(cfg.CORS)

	// repositories
	userRepository := repositories.NewUserRepository(cfg.Postgres)

	// get services
	userService := services.NewUserService(userRepository)

	// rest routes with services injected, can pass cfg too
	rest.RegisterUserRoutes(router, userService)
	//TODO Add similar routers for other too

	start(router, cfg)
}

func start(router *gin.Engine, cfg config.AppConfig) {

	// runs on 8080 on default
	err := router.Run(":" + cfg.ServerPort)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
}
