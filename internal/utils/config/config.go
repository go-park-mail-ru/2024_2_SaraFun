package config

import (
	"fmt"
	"github.com/caarlos0/env/v11"
	"go.uber.org/zap"
	"os"
)

type EnvConfig struct {
	RedisUser     string `env:"REDIS_USER"`
	RedisPassword string `env:"REDIS_PASSWORD"`
	DbHost        string `env:"DB_HOST"`
	DbPort        string `env:"DB_PORT"`
	DbUser        string `env:"DB_USER"`
	DbPassword    string `env:"DB_PASSWORD"`
	DbName        string `env:"DB_NAME"`
	DbSSLMode     string `env:"DB_SSLMODE"`
}

func NewConfig(logger *zap.Logger) (EnvConfig, error) {
	var config EnvConfig
	db_host := os.Getenv("DB_HOST")
	logger.Info("db_host", zap.String("db_host", db_host))
	err := env.Parse(&config)
	if err != nil {
		logger.Error("error parsing config", zap.Error(err))
		return EnvConfig{}, fmt.Errorf("error parsing config: %w", err)
	}
	logger.Info("config parse result", zap.Any("config", config))
	return config, nil
}
