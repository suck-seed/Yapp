package main

import (
	"fmt"
	"log"

	"github.com/suck-seed/yapp/config"
	"github.com/suck-seed/yapp/internal/api"
	"github.com/suck-seed/yapp/internal/database"
)

func main() {

	// setting gin SETMODE
	// gin.SetMode(gin.DebugMode)

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
