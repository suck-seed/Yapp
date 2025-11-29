package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	messageService := services.NewMessageService(hallRepository, roomRepository, messageRepository, userRepository)

	//	presist function
	presistFunction := ws.MakePresistFunction(messageService, userService)
	hub := ws.NewHub(presistFunction)
	go hub.Run()

	// Public Router ( Do not pass AuthMiddleware here pls )
	rest.RegisterAuthRoutes(router, userService)

	// Protected API routed (JWT required, AuthMiddleware Passed)
	api := router.Group("/api/v1")
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

	// start(router, cfg)
	startGracefully(router, cfg, &hub)
}

func startGracefully(router *gin.Engine, cfg config.AppConfig, hub *ws.Hub) {

	server := http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: router,

		// Should be longer than Nginx proxy_read_timeout (30s)
		ReadTimeout: 45 * time.Second,

		// Should be longer than Nginx proxy_send_timeout (30s)
		WriteTimeout: 45 * time.Second,

		// IdleTimeout: Maximum time to wait for the next request when keep-alives are enabled
		// CRITICAL: Must be LONGER than Nginx keepalive_timeout (65s default)
		// This ensures Nginx closes the connection first, preventing "broken pipe" errors
		IdleTimeout: 120 * time.Second,

		// ReadHeaderTimeout: Time allowed to read request headers
		ReadHeaderTimeout: 10 * time.Second,

		MaxHeaderBytes: 1 << 20, // 1 MB

	}

	// starting server in a go routine
	go func() {

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("error: %v\n", err)
		}

	}()

	// Shutdown handler
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("Shutdown signal received, starting graceful shutdown...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// CLosing websockets
	if err := hub.Close(); err != nil {
		log.Printf("Error closing WebSocket hub: %v", err)
	}

	// http server clousere
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Error closing http server: %v", err)
	}

	// close postgres db
	cfg.Postgres.Close()

	// log message
	log.Printf("Gracefully stopped server")

}

func start(router *gin.Engine, cfg config.AppConfig) {

	// runs on 8080 on default
	err := router.Run(":" + cfg.ServerPort)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
}
