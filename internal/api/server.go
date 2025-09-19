package api

import (
	"fmt"

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
	messageService := services.NewMessageService(roomRepository, messageRepository)

	// Routes Handler
	rest.RegisterAuthRoutes(router, userService)

	// Protected API routed (JWT required)
	api := router.Group("/api")
	api.Use(auth.AuthMiddleware())
	{
		rest.RegisterUserRoutes(api, userService)
		rest.RegisterHallRoutes(api, hallService)
		rest.RegisterFloorRoutes(api, floorService)
		rest.RegisterRoomRoutes(api, roomService)
		rest.RegisterMessageRoutes(api, messageService)

	}
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
