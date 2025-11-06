package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

type ServerConfig struct {
	Port string `env:"PORT" envDefault:"8080"`
	Env  string `env:"ENV" envDefault:"development"`
}

type DatabaseConfig struct {
	URL string `env:"DB_CONN_URL,required"`
}

func Load() (*Config, error) {
	_ = godotenv.Load() // Ignore .env file loading error in case we have our envs set
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return cfg, nil
}
