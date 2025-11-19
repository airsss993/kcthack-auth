package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	HTTP struct {
		Port string
	}

	Database struct {
		DSN string
	}

	JWT struct {
		JWTSecret string
	}

	Auth struct {
		AccessTTL  time.Duration
		RefreshTTL time.Duration
	}
}

func Init() (*Config, error) {
	_ = godotenv.Load()

	var cfg Config

	err := loadFromEnv(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to load env vars from .env: %w", err)
	}

	viper.SetConfigName("main")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath("../configs")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to find yaml config file: %w", err)
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to set vars from .yml: %w", err)
	}

	return &cfg, nil
}

func loadFromEnv(cfg *Config) error {
	cfg.Database.DSN = os.Getenv("DATABASE_DSN")
	if cfg.Database.DSN == "" {
		return errors.New("DATABASE_DSN env var is required")
	}

	cfg.JWT.JWTSecret = os.Getenv("JWT_SECRET")
	if cfg.JWT.JWTSecret == "" {
		return errors.New("JWT_SECRET env var is required")
	}

	return nil
}
