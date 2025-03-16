package config

import "os"

// Config содержит конфигурацию сервиса
type Config struct {
	DatabaseURL       string
	RedisURL          string
	GrpcServerAddress string
}

// Load загружает конфигурацию из переменных окружения
func Load() (*Config, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@postgres:5432/chatapp"
	}

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "redis:6379"
	}

	grpcAddr := os.Getenv("GRPC_SERVER_ADDRESS")
	if grpcAddr == "" {
		grpcAddr = ":50052"
	}

	return &Config{
		DatabaseURL:       dbURL,
		RedisURL:          redisURL,
		GrpcServerAddress: grpcAddr,
	}, nil
}
