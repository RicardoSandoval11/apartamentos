package conf

import (
	"fmt"
	"log/slog"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type Config struct {
	ServiceName        string `env:"SERVICE_NAME" envDefault:"local"`
	ServiceVersion     string `env:"SERVICE_VERSION,required"`
	ApplicationPort    string `env:"APPLICATION_PORT" envDefault:":8080"`
	Environment        string `env:"ENVIRONMENT" envDefault:"LOCAL"`
	DbConnectionString string `env:"DATABASE_CONN_STRING,required"`
}

func GetEnv() (Config, error) {
	var cfg Config

	if err := godotenv.Load("conf/.env.local"); err != nil {
		slog.Debug("file not found")
	}

	if err := env.Parse(&cfg); err != nil {
		return Config{}, fmt.Errorf("could not load environment variables: %w", err)
	}

	return cfg, nil
}
