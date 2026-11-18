package conf

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

const defaultEnv = "LOCAL"

type Config struct {
	ServiceName         string `env:"SERVICE_NAME" envDefault:"local"`
	ServiceVersion      string `env:"SERVICE_VERSION,required"`
	ApplicationPort     string `env:"APPLICATION_PORT" envDefault:":8080"`
	Environment         string `env:"ENVIRONMENT" envDefault:"LOCAL"`
	DbConnectionString  string `env:"DATABASE_CONN_STRING,required"`
	PrivateKeyAuthToken string `env:"PRIVATE_KEY_AUTH_TOKEN,required"`
	PublicKeyAuthToken  string `env:"PUBLIC_KEY_AUTH_TOKEN,required"`
}

func getEnvironment() string {
	env := os.Getenv("ENVIRONMENT")

	if env == "" {
		return defaultEnv
	}

	return env
}

func getEnvFile() string {
	env := getEnvironment()

	switch env {
	case "LOCAL":
		return "conf/.env.local"
	case "DEVELOPMENT":
		return "conf/.env.development"
	case "PRODUCTION":
		return "conf/.env.production"
	default:
		return "conf/.env.local"
	}
}

func GetEnv() (Config, error) {
	var cfg Config

	if err := godotenv.Load(getEnvFile()); err != nil {
		slog.Debug("file not found")
	}

	if err := env.Parse(&cfg); err != nil {
		return Config{}, fmt.Errorf("could not load environment variables: %w", err)
	}

	return cfg, nil
}
