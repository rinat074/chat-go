package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rinat074/chat-go/pkg/logger"
	"github.com/rinat074/chat-go/services/gateway-service/internal/clients"
	"github.com/rinat074/chat-go/services/gateway-service/internal/config"
	"github.com/rinat074/chat-go/services/gateway-service/internal/handlers"
	"github.com/rinat074/chat-go/services/gateway-service/internal/middleware"
	"github.com/rinat074/chat-go/services/gateway-service/internal/server"
	"go.uber.org/zap"
)

func main() {
	// Инициализация логгера
	logger.Init("gateway-service", "info")
	defer func() {
		// Записываем логи перед выходом
		logger.Info("Сервис Gateway завершил работу")
	}()

	// Загрузка конфигурации
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("Ошибка загрузки конфигурации", logger.WithError(err))
	}

	// Создание gRPC клиентов для микросервисов
	authClient, err := clients.NewAuthClient(cfg.AuthServiceAddress)
	if err != nil {
		logger.Fatal("Ошибка создания клиента auth-service", logger.WithError(err))
	}
	defer authClient.Close()

	chatClient, err := clients.NewChatClient(cfg.ChatServiceAddress)
	if err != nil {
		logger.Fatal("Ошибка создания клиента chat-service", logger.WithError(err))
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
	srv := server.NewServer(
		authHandler,
		chatHandler,
		wsHandler,
		authMiddleware,
		rateLimiter,
	)

	// Настройка HTTP-сервера
	httpServer := &http.Server{
		Addr:         cfg.HTTPServerAddress,
		Handler:      srv.Router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Запуск сервера в горутине
	go func() {
		logger.Info("API Gateway запущен", zap.String("address", cfg.HTTPServerAddress))
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Ошибка HTTP-сервера", logger.WithError(err))
		}
	}()

	// Обработка сигналов
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Graceful shutdown
	logger.Info("Завершение работы сервера...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Fatal("Ошибка при остановке сервера", logger.WithError(err))
	}
}
