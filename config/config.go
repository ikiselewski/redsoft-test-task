package config

import (
	"redsoft-test-task/internal/database"
	"redsoft-test-task/internal/srv"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	database.DBConfig
	srv.SrvConfig
}

func Get() (*Config, error) {
	var cfg Config
	err := env.Parse(&cfg)
	return &cfg, err
}
