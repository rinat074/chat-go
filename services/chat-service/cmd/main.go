package main

import (
	"log"
	"net"

	"github.com/rinat074/chat-go/services/chat-service/internal/cache"
	"github.com/rinat074/chat-go/services/chat-service/internal/config"
	"github.com/rinat074/chat-go/services/chat-service/internal/db"
	"github.com/rinat074/chat-go/services/chat-service/internal/server"
	"github.com/rinat074/chat-go/services/chat-service/internal/service"

	"google.golang.org/grpc"
)

func main() {
	// Загрузка конфигурации
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	// Подключение к базе данных
	database, err := db.NewConnection(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer database.Close()

	// Подключение к Redis
	redisCache := cache.NewRedisCache(cfg.RedisURL)

	// Создание сервиса чата
	chatService := service.NewChatService(database, redisCache)

	// Создание и запуск WebSocket хаба
	hub := service.NewHub(chatService)
	go hub.Run()

	// Создание gRPC сервера
	grpcServer := grpc.NewServer()
	server.RegisterChatServer(grpcServer, chatService, hub)

	// Запуск gRPC сервера
	lis, err := net.Listen("tcp", cfg.GrpcServerAddress)
	if err != nil {
		log.Fatalf("Ошибка при запуске gRPC сервера: %v", err)
	}

	log.Printf("Chat-сервис запущен на %s", cfg.GrpcServerAddress)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Ошибка gRPC сервера: %v", err)
	}
}
