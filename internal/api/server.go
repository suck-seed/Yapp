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

	"github.com/gin-gonic/gin"
	"github.com/suck-seed/yapp/config"
	"github.com/suck-seed/yapp/internal/api/rest"
	"github.com/suck-seed/yapp/internal/auth"
	"github.com/suck-seed/yapp/internal/repositories"
	"github.com/suck-seed/yapp/internal/services"
	"github.com/suck-seed/yapp/internal/ws"
)

func StartServer(cfg config.AppConfig) {
	router := gin.Default()

	router.Use(auth.CSRFCookieMiddleware())
	router.Use(cfg.CORS)

	// health check
	router.GET("/health", func(c *gin.Context) {
		if cfg.PostgresPool != nil {
			if err := cfg.PostgresPool.Ping(context.Background()); err != nil {
				c.JSON(http.StatusServiceUnavailable, gin.H{
					"status": "unhealthy",
					"error":  "database unavailable",
				})
				return
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// Dependency Injection
	userRepository := repositories.NewUserRepository()
	hallRepository := repositories.NewHallRepository()
	roleRepository := repositories.NewRoleRepository()
	banRepository := repositories.NewBanRepository()
	floorRepository := repositories.NewFloorRepository()
	roomRepository := repositories.NewRoomRepository()
	messageRepository := repositories.NewMessageRepository()
	inviteRepository := repositories.NewInviteRepository()

	permissionCheckerService := services.NewPermissionCheckerService(roleRepository, userRepository, hallRepository, banRepository, cfg.PostgresPool)

	userService := services.NewUserService(userRepository, cfg.PostgresPool)
	hallService := services.NewHallService(hallRepository, userRepository, roleRepository, banRepository, permissionCheckerService, cfg.PostgresPool)
	floorService := services.NewFloorService(hallRepository, floorRepository, roomRepository, banRepository, permissionCheckerService, cfg.PostgresPool)
	roomService := services.NewRoomService(hallRepository, floorRepository, roomRepository, banRepository, permissionCheckerService, cfg.PostgresPool)
	roleService := services.NewRoleService(roleRepository, userRepository, hallRepository, banRepository, permissionCheckerService, cfg.PostgresPool)
	banService := services.NewBanService(banRepository, userRepository, hallRepository, permissionCheckerService, cfg.PostgresPool)
	messageService := services.NewMessageService(hallRepository, roomRepository, messageRepository, userRepository, permissionCheckerService, cfg.PostgresPool)
	inviteService := services.NewInviteService(inviteRepository, hallRepository, roleRepository, permissionCheckerService, cfg.PostgresPool)

	presistFunction := ws.MakePresistFunction(messageService, userService)
	hub := ws.NewHub(presistFunction)
	go hub.Run()

	// Routes
	rest.RegisterAuthRoutes(router, userService)

	api := router.Group("/api/v1")
	api.Use(func(c *gin.Context) {
		log.Printf(">>> MIDDLEWARE HIT: %s %s", c.Request.Method, c.Request.URL.Path)
		c.Next()
		log.Printf(">>> RESPONSE STATUS: %d", c.Writer.Status())
	})
	api.Use(auth.AuthMiddleware())
	{
		rest.RegisterUserRoutes(api, userService)
		rest.RegisterHallRoutes(api, hallService, roleService, banService, inviteService, floorService, roomService, messageService)
		// rest.RegisterFloorRoutes(api, floorService)
		// rest.RegisterRoomRoutes(api, roomService)
		rest.RegisterMessageRoutes(api, messageService)
		rest.RegisterInviteRoutes(api, inviteService)
	}

	wsRouter := router.Group("/ws")
	wsRouter.Use(auth.AuthMiddleware())
	{
		rest.RegisterWebSocketRoutes(wsRouter, &hub, messageService, hallService, roomService, userService)
	}

	startGracefully(router, cfg, &hub)
}

func startGracefully(router *gin.Engine, cfg config.AppConfig, hub *ws.Hub) {
	server := http.Server{
		Addr:         "0.0.0.0:" + cfg.ServerPort,
		Handler:      router,
		ReadTimeout:  45 * time.Second,
		WriteTimeout: 45 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("error: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutdown signal received, starting graceful shutdown...")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	hub.Close()
	server.Shutdown(ctx)
	cfg.PostgresPool.Close()
	log.Printf("Gracefully stopped server")
}
