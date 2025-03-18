package config

import (
	"time"

	"github.com/rinat074/chat-go/pkg/config"
)

// Config содержит настройки приложения.
type Config struct {
	HTTPServerAddress  string
	AuthServiceAddress string
	ChatServiceAddress string
	RedisURL           string
	RateLimitDuration  time.Duration
}

// Load загружает конфигурацию из переменных окружения.
func Load() (*Config, error) {
	// Загружаем переменные окружения из .env файла, если он существует
	if err := config.Load(); err != nil {
		return nil, err
	}

	// Получаем настройки из переменных окружения
	return &Config{
		HTTPServerAddress:  config.GetString("HTTP_SERVER_ADDRESS", ":8080"),
		AuthServiceAddress: config.GetString("AUTH_SERVICE_ADDRESS", "localhost:50051"),
		ChatServiceAddress: config.GetString("CHAT_SERVICE_ADDRESS", "localhost:50052"),
		RedisURL:           config.GetString("REDIS_URL", "localhost:6379"),
		RateLimitDuration:  config.GetDuration("RATE_LIMIT_DURATION", 1*time.Minute),
	}, nil
}
