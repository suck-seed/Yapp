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
	hallRepository := repositories.NewHallReposiroty(cfg.Postgres)
	floorRepository := repositories.NewFloorRepository(cfg.Postgres)
	roomRepository := repositories.NewRoomReposiroty(cfg.Postgres)
	messageRepository := repositories.NewMessageReposiroty(cfg.Postgres)

	// Service & Dependency Injection for services
	userService := services.NewUserService(userRepository)
	hallService := services.NewHallService(hallRepository)
	floorService := services.NewFloorService(hallRepository, floorRepository)
	roomService := services.NewRoomService(hallRepository, floorRepository, roomRepository)
	messageService := services.NewMessageService(roomRepository, messageRepository)

	// Routes Handler
	rest.RegisterUserRoutes(router, userService)
	rest.RegisterAuthRoutes(router, userService)
	rest.RegisterHallRoutes(router, hallService)
	rest.RegisterFloorRoutes(router, floorService)
	rest.RegisterRoomRoutes(router, roomService)
	rest.RegisterMessageRoutes(router, messageService)
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
