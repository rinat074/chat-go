package main

import (
	"context"
	"log"
	"net/http"

	"github.com/redis/go-redis/v9"
	"github.com/rinat074/chat-go/services/gateway-service/internal/clients"
	"github.com/rinat074/chat-go/services/gateway-service/internal/config"
	"github.com/rinat074/chat-go/services/gateway-service/internal/handlers"
	"github.com/rinat074/chat-go/services/gateway-service/internal/middleware"
	"github.com/rinat074/chat-go/services/gateway-service/internal/server"
	"github.com/rinat074/chat-go/services/gateway-service/pkg/logger"

	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	// Логгер
	logger := logger.NewMockLogger()

	// Загрузка конфигурации
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("ошибка загрузки конфигурации: %v", err)
	}

	// Создание клиентов сервисов
	authClient, err := clients.NewAuthClient(cfg.AuthServiceAddress)
	if err != nil {
		log.Fatalf("ошибка подключения к auth-service: %v", err)
	}
	defer authClient.Close()

	chatClient, err := clients.NewChatClient(cfg.ChatServiceAddress)
	if err != nil {
		log.Fatalf("ошибка подключения к chat-service: %v", err)
	}
	defer chatClient.Close()

	// Подключение к Redis для управления сессиями и кэширования
	redisClient := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddress,
	})
	defer redisClient.Close()

	// Проверка подключения к Redis
	if _, err := redisClient.Ping(context.Background()).Result(); err != nil {
		log.Fatalf("ошибка подключения к redis: %v", err)
	}

	// Создание клиентов для сервисов
	serviceClients := clients.NewServiceClients(authClient, chatClient, redisClient)

	// Создание обработчиков
	authHandler := handlers.NewAuthHandler(serviceClients, logger)
	chatHandler := handlers.NewChatHandler(serviceClients, logger)
	wsHandler := handlers.NewWebSocketHandler(serviceClients, logger)

	// Создание middleware
	authMiddleware := middleware.NewAuthMiddleware(serviceClients, logger, cfg.JwtSecret)
	rateLimiter := middleware.NewRateLimiter(redisClient)

	// Создание HTTP сервера
	srv := server.NewServer(authHandler, chatHandler, wsHandler, authMiddleware, rateLimiter)

	// Swagger
	srv.Router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	// Запуск сервера
	log.Printf("Gateway-сервис запущен на %s", cfg.ServerAddress)
	if err := http.ListenAndServe(cfg.ServerAddress, srv.Router); err != nil {
		log.Fatalf("ошибка запуска сервера: %v", err)
	}
}
