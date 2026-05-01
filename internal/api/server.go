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

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func StartServer(cfg config.AppConfig) {

	// Engine instance with the Logger and Recovery middleware already attached.
	router := gin.Default()

	router.Use(auth.CSRFCookieMiddleware())
	router.Use(cfg.CORS)

	// --- Swagger UI -------
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
	readRecieptFunction := ws.MakeReadReceiptFunction(messageService)
	hub := ws.NewHub(presistFunction, readRecieptFunction)
	go hub.Run()

	// Routes

	apiv1 := router.Group("/api/v1")

	// ---- PUBLIC ROUTES ,  NO AUTHENTICATION
	{
		rest.RegisterAuthRoutes(apiv1, userService)
		rest.RegisterInvitePublicRoutes(apiv1, inviteService)
	}

	// For endpoint with authentication required
	protectedv1 := apiv1.Group("", auth.AuthMiddleware())
	{
		rest.RegisterUserRoutes(protectedv1, userService)
		rest.RegisterHallRoutes(protectedv1, hallService, roleService, banService, inviteService, floorService, roomService, messageService)
		rest.RegisterMessageRoutes(protectedv1, messageService)
		rest.RegisterInvitePrivateRoutes(protectedv1, inviteService)
	}

	wsHandler := router.Group("/ws", auth.AuthMiddleware())
	{
		rest.RegisterWebSocketRoutes(wsHandler, &hub, messageService, hallService, roomService, userService)
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
