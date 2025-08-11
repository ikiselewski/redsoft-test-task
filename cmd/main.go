package main

import (
	"fmt"
	"log"
	"redsoft-test-task/config"
	"redsoft-test-task/internal/srv"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Get()
	if err != nil {
		log.Fatalf("%s", fmt.Errorf("unable to start server due to %w", err).Error())
	}

	// db, err := database.New(cfg.DSN)
	// if err != nil {
	// 	log.Fatalf("%s", fmt.Errorf("unable to start server due to %w", err).Error())
	// }

	server, err := srv.New(&cfg.SrvConfig, &srv.Dependencies{
		Database: nil,
	})
	if err != nil {
		log.Fatalf("%s", fmt.Errorf("unable to start server due to %w", err).Error())
	}

	r := gin.Default()

	srv.RegisterHandlersWithOptions(r, server, srv.GinServerOptions{})

	log.Fatal(r.Run(":8080"))
}
