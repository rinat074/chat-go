package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config содержит настройки приложения
type Config struct {
	HTTPServerAddress  string
	AuthServiceAddress string
	ChatServiceAddress string
	RedisURL           string
	RateLimitDuration  time.Duration
}

// Load загружает конфигурацию из переменных окружения
func Load() (*Config, error) {
	// Загружаем переменные окружения из .env файла, если он существует
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			return nil, fmt.Errorf("error loading .env file: %w", err)
		}
	}

	// Получаем настройки из переменных окружения
	return &Config{
		HTTPServerAddress:  GetString("HTTP_SERVER_ADDRESS", ":8080"),
		AuthServiceAddress: GetString("AUTH_SERVICE_ADDRESS", "localhost:50051"),
		ChatServiceAddress: GetString("CHAT_SERVICE_ADDRESS", "localhost:50052"),
		RedisURL:           GetString("REDIS_URL", "localhost:6379"),
		RateLimitDuration:  GetDuration("RATE_LIMIT_DURATION", 1*time.Minute),
	}, nil
}

// GetString возвращает значение переменной окружения как string
func GetString(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

// GetInt возвращает значение переменной окружения как int
func GetInt(key string, defaultValue int) int {
	valueStr, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

// GetFloat возвращает значение переменной окружения как float64
func GetFloat(key string, defaultValue float64) float64 {
	valueStr, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return defaultValue
	}
	return value
}

// GetBool возвращает значение переменной окружения как bool
func GetBool(key string, defaultValue bool) bool {
	valueStr, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

// GetDuration возвращает значение переменной окружения как time.Duration
func GetDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	value, err := time.ParseDuration(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

// RequireString возвращает значение переменной окружения как string или паникует если оно не установлено.
func RequireString(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		panic(fmt.Sprintf("required environment variable %s not set", key))
	}
	return value
}
