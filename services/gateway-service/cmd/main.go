package main

import (
	"log"
	"net/http"

	"chat-app/services/gateway-service/internal/clients"
	"chat-app/services/gateway-service/internal/config"
	"chat-app/services/gateway-service/internal/handlers"
	"chat-app/services/gateway-service/internal/middleware"
	"chat-app/services/gateway-service/internal/server"
)

func main() {
	// Загрузка конфигурации
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	// Создание gRPC клиентов для микросервисов
	authClient, err := clients.NewAuthClient(cfg.AuthServiceAddress)
	if err != nil {
		log.Fatalf("Ошибка создания клиента auth-service: %v", err)
	}
	defer authClient.Close()

	chatClient, err := clients.NewChatClient(cfg.ChatServiceAddress)
	if err != nil {
		log.Fatalf("Ошибка создания клиента chat-service: %v", err)
	}
	defer chatClient.Close()

	// Инициализация Redis для rate limiting и кэширования
	redisClient := clients.NewRedisClient(cfg.RedisURL)

	// Создание обработчиков API
	authHandler := handlers.NewAuthHandler(authClient)
	chatHandler := handlers.NewChatHandler(chatClient)
	wsHandler := handlers.NewWebSocketHandler(chatClient)

	// Создание middleware
	rateLimiter := middleware.NewRateLimiter(redisClient, 100, cfg.RateLimitDuration)
	authMiddleware := middleware.NewAuthMiddleware(authClient)

	// Создание и запуск сервера
	server := server.NewServer(
		authHandler,
		chatHandler,
		wsHandler,
		authMiddleware,
		rateLimiter,
	)

	log.Printf("API Gateway запущен на %s", cfg.HTTPServerAddress)
	if err := http.ListenAndServe(cfg.HTTPServerAddress, server.Router); err != nil {
		log.Fatalf("Ошибка HTTP сервера: %v", err)
	}
}
