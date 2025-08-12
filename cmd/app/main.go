package main

import (
	"fmt"
	"log"
	"redsoft-test-task/config"
	"redsoft-test-task/internal/database"
	"redsoft-test-task/internal/srv"

	"github.com/gin-gonic/gin"
	"github.com/masonkmeyer/agify"
	"github.com/masonkmeyer/genderize"
	"github.com/masonkmeyer/nationalize"
)

func main() {
	cfg, err := config.Get()
	if err != nil {
		log.Fatalf("%s", fmt.Errorf("unable to start server due to %w", err).Error())
	}

	db, err := database.New(cfg.DSN)
	if err != nil {
		log.Fatalf("%s", fmt.Errorf("unable to start server due to %w", err).Error())
	}

	nationsClient := nationalize.NewClient()
	agesClient := agify.NewClient()
	gendersClient := genderize.NewClient()

	server, err := srv.New(&cfg.SrvConfig, &srv.Dependencies{
		Database:    db,
		Nationalize: nationsClient,
		Agify:       agesClient,
		Genderize:   gendersClient,
	})
	if err != nil {
		log.Fatalf("%s", fmt.Errorf("unable to start server due to %w", err).Error())
	}

	r := gin.Default()

	srv.RegisterHandlersWithOptions(r, server, srv.GinServerOptions{})

	log.Fatal(r.Run(":8080"))
}
