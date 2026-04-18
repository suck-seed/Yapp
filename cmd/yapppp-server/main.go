// Package main is the entry point for the Yapp API server.
//
// @title           Yapp API
// @version         1.0
// @description     Real-time chat + voice application API (REST & WebSocket).
// @termsOfService  http://swagger.io/terms/
//
// @contact.name   Yapp API Support
// @contact.email  sandeshburner@gmail.com
//
// @license.name  MIT
//
// @host      yappserver.onrender.com
// @Schemes   https
// @BasePath  /api/v1
//
// @securityDefinitions.apikey  CookieAuth
// @in                          cookie
// @name                        jwt
// @description                 JWT token stored in the `jwt` HttpOnly cookie. Set automatically on sign-in.
package main

import (
	"fmt"
	"log"

	"github.com/suck-seed/yapp/config"
	"github.com/suck-seed/yapp/internal/api"
	"github.com/suck-seed/yapp/internal/database"

	_ "github.com/suck-seed/yapp/docs" // ← this must be present
)

func main() {

	//	get config
	cfg, err := config.SetupEnvironment()
	if err != nil {
		log.Fatalf("Config File is not loaded properly: %v\n", err)
	}

	if err := database.RunProductionMigrations(); err != nil {
		log.Fatalf("Migration failed: %v\n", err)
	}

	// start server
	fmt.Println(cfg.ServerPort)
	api.StartServer(cfg)

}
