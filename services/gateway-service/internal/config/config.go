package config

import "os"

// Config содержит конфигурацию сервиса
type Config struct {
	ServerAddress      string
	AuthServiceAddress string
	ChatServiceAddress string
	RedisAddress       string
	JwtSecret          string
}

// Load загружает конфигурацию из переменных окружения
func Load() (*Config, error) {
	serverAddr := os.Getenv("SERVER_ADDRESS")
	if serverAddr == "" {
		serverAddr = ":8080"
	}

	authServiceAddr := os.Getenv("AUTH_SERVICE_ADDRESS")
	if authServiceAddr == "" {
		authServiceAddr = "auth-service:50051"
	}

	chatServiceAddr := os.Getenv("CHAT_SERVICE_ADDRESS")
	if chatServiceAddr == "" {
		chatServiceAddr = "chat-service:50052"
	}

	redisAddr := os.Getenv("REDIS_ADDRESS")
	if redisAddr == "" {
		redisAddr = "redis:6379"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "ваш-секретный-ключ" // В продакшене должен быть безопасно сохранен
	}

	return &Config{
		ServerAddress:      serverAddr,
		AuthServiceAddress: authServiceAddr,
		ChatServiceAddress: chatServiceAddr,
		RedisAddress:       redisAddr,
		JwtSecret:          jwtSecret,
	}, nil
}
