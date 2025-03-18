package config

import "time"

// Config содержит настройки для сервера
type Config struct {
	HTTPServerAddress  string
	AuthServiceAddress string
	ChatServiceAddress string
	RedisURL           string
	RateLimitDuration  time.Duration
}

// Load загружает конфигурацию из переменных окружения
func Load() (*Config, error) {
	return &Config{
		HTTPServerAddress:  ":8080",
		AuthServiceAddress: "localhost:50051",
		ChatServiceAddress: "localhost:50052",
		RedisURL:           "localhost:6379",
		RateLimitDuration:  time.Minute,
	}, nil
}
