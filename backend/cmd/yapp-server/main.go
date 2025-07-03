package main

import (
	"log"

	"github.com/suck-seed/yapp/config"
	"github.com/suck-seed/yapp/internal/api"
)

func main() {

	//	get config
	cfg, err := config.SetupEnvironment()
	if err != nil {
		log.Fatalf("Config File is not loaded properly: %v\n", err)
	}

	// start server
	api.StartServer(cfg)

}
