package conf

import (
	"errors"

	"github.com/caarlos0/env"
)

type Config struct {
	ServiceName        string `env:"SERVICE_NAME" envDefault:"local"`
	ServiceVersion     string `env:"SERVICE_VERSION,required"`
	ApplicationPort    string `env:"APPLICATION_PORT" envDefautl:":8080"`
	Environment        string `env:"ENVIRONMENT" envDefault:"LOCAL"`
	DbConnectionString string `env:"DATABASE_CONN_STRING,required"`
}

func GetEnv() (Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return Config{}, errors.New("could not load environment variables")
	}

	return cfg, nil
}
