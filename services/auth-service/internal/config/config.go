package config

import "os"

// Config содержит конфигурацию сервиса
type Config struct {
	DatabaseURL       string
	GrpcServerAddress string
	JwtSecret         string
}

// Load загружает конфигурацию из переменных окружения
func Load() (*Config, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@postgres:5432/chatapp"
	}

	grpcAddr := os.Getenv("GRPC_SERVER_ADDRESS")
	if grpcAddr == "" {
		grpcAddr = ":50051"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "ваш-секретный-ключ" // В продакшене должен быть безопасно сохранен
	}

	return &Config{
		DatabaseURL:       dbURL,
		GrpcServerAddress: grpcAddr,
		JwtSecret:         jwtSecret,
	}, nil
}
