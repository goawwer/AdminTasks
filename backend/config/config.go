package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	// database
	DatabaseName     string `env:"DB_NAME"`
	DatabaseHost     string `env:"DB_HOST"`
	DatabasePort     string `env:"DB_PORT"`
	DatabaseUser     string `env:"DB_USER"`
	DatabasePassword string `env:"DB_PASSWORD"`

	// api
	BindAddr    string `env:"BIND_ADDR"`
	LoggerLevel string `env:"LOGGER_LEVEL"`
}

func (c *Config) DatabaseURI() string {
	return fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		c.DatabaseUser,
		c.DatabasePassword,
		c.DatabaseHost,
		c.DatabasePort,
		c.DatabaseName,
	)
}

func New() (*Config, error) {
	cfg, err := env.ParseAs[Config]()
	if err != nil {
		return nil, fmt.Errorf("cannot load config: %w", err)
	}

	return &cfg, nil
}
