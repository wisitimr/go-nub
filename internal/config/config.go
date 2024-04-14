package config

import (
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Configuration struct {
	HTTPServer
	Database
}

type HTTPServer struct {
	IdleTimeout  time.Duration `envconfig:"HTTP_SERVER_IDLE_TIMEOUT" default:"60s"`
	Port         int           `envconfig:"PORT" default:"8080"`
	ReadTimeout  time.Duration `envconfig:"HTTP_SERVER_READ_TIMEOUT" default:"1s"`
	WriteTimeout time.Duration `envconfig:"HTTP_SERVER_WRITE_TIMEOUT" default:"2s"`
}

type Database struct {
	URI        string `envconfig:"DATABASE_URL" required:"true"`
	Name       string `envconfig:"DATABASE_NAME" default:"nub"`
	Collection string
}

func Load() (Configuration, error) {
	if err := godotenv.Load(); err != nil {
		log.Fatal("load .env error")
	}
	var cfg Configuration
	if err := envconfig.Process("", &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}
