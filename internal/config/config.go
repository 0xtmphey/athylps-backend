package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	Server     ServerConfig
	Database   DatabaseConfig
	RevenueCat RevenueCatConfig
}

type ServerConfig struct {
	Port string `env:"PORT" envDefault:"8080"`
	Env  string `env:"ENV" envDefault:"development"`
}

type DatabaseConfig struct {
	Name     string `env:"DB_NAME,required"`
	Port     string `env:"DB_PORT,required"`
	Host     string `env:"DB_HOST,required"`
	User     string `env:"DB_USER,required"`
	Password string `env:"DB_PASSWORD,required"`
}

type RevenueCatConfig struct {
	BearerToken string `env:"RC_BEARER,required"`
}

type TelegramConfig struct {
	BotToken     string `env:"BOT_TOKEN,required"`
	NotifyChatID string `env:"NOTIFY_CHAT_ID,required"`
}

func Load() (*Config, error) {
	_ = godotenv.Load() // Ignore .env file loading error in case we have our envs set
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return cfg, nil
}

func (dbConfig *DatabaseConfig) ConnectionUrl() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Name,
	)
}
