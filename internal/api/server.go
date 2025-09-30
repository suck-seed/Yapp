package api

import (
	"fmt"

	"github.com/suck-seed/yapp/internal/ws"

	"github.com/gin-gonic/gin"
	"github.com/suck-seed/yapp/config"
	"github.com/suck-seed/yapp/internal/api/rest"
	"github.com/suck-seed/yapp/internal/auth"
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
	hallRepository := repositories.NewHallRepository(cfg.Postgres)
	floorRepository := repositories.NewFloorRepository(cfg.Postgres)
	roomRepository := repositories.NewRoomRepository(cfg.Postgres)
	messageRepository := repositories.NewMessageRepository(cfg.Postgres)

	// Service & Dependency Injection for services
	userService := services.NewUserService(userRepository)
	hallService := services.NewHallService(hallRepository)
	floorService := services.NewFloorService(hallRepository, floorRepository)
	roomService := services.NewRoomService(hallRepository, floorRepository, roomRepository)
	messageService := services.NewMessageService(roomRepository, messageRepository, userRepository)

	//	presist function
	presistFunction := ws.MakePresistFunction(messageService, userService)
	hub := ws.NewHub(presistFunction)
	go hub.Run()

	// Public Router ( Do not pass AuthMiddleware here pls )
	rest.RegisterAuthRoutes(router, userService)

	// Protected API routed (JWT required, AuthMiddleware Passed)
	api := router.Group("/api")
	api.Use(auth.AuthMiddleware())
	{
		rest.RegisterUserRoutes(api, userService)
		rest.RegisterHallRoutes(api, hallService)
		rest.RegisterFloorRoutes(api, floorService)
		rest.RegisterRoomRoutes(api, roomService)
		rest.RegisterMessageRoutes(api, messageService)

	}

	wsRouter := router.Group("/ws")
	wsRouter.Use(auth.AuthMiddleware())
	{

		rest.RegisterWebSocketRoutes(wsRouter, &hub, messageService, hallService, roomService, userService)
	}

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
