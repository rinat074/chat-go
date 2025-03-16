package main

import (
	"log"
	"net"

	"chat-app/services/auth-service/internal/config"
	"chat-app/services/auth-service/internal/db"
	"chat-app/services/auth-service/internal/server"
	"chat-app/services/auth-service/internal/service"

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

	// Создание сервиса аутентификации
	authService := service.NewAuthService(database)

	// Создание gRPC сервера
	grpcServer := grpc.NewServer()
	server.RegisterAuthServer(grpcServer, authService)

	// Запуск gRPC сервера
	lis, err := net.Listen("tcp", cfg.GrpcServerAddress)
	if err != nil {
		log.Fatalf("Ошибка при запуске gRPC сервера: %v", err)
	}

	log.Printf("Auth-сервис запущен на %s", cfg.GrpcServerAddress)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Ошибка gRPC сервера: %v", err)
	}
}
