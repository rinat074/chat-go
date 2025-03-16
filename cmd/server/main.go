package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"chat-app/internal/auth"
	"chat-app/internal/cache"
	"chat-app/internal/chat"
	"chat-app/internal/db"
	"chat-app/internal/middleware"
	"chat-app/internal/server"
)

func main() {
	// Получаем настройки из переменных окружения
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/chatapp"
	}

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "localhost:6379"
	}

	// Инициализация подключения к базе данных
	database, err := db.NewConnection(dbURL)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer database.Close()

	// Инициализация Redis кэша
	redisCache := cache.NewRedisCache(redisURL)

	// Инициализация сервисов
	authService := auth.NewService(database)
	chatService := chat.NewService(database, redisCache)

	// Инициализация хаба веб-сокетов
	hub := chat.NewHub(chatService)
	go hub.Run()

	// Инициализация обработчиков
	authHandler := auth.NewHandler(authService)
	chatHandler := chat.NewHandler(hub, chatService)

	// Создание Rate Limiter
	rateLimiter := middleware.NewRateLimiter(redisCache.Client, 100, time.Minute)

	// Создание и запуск сервера
	srv := server.NewServer(authHandler, chatHandler, authService, rateLimiter)

	log.Println("Сервер запущен на http://localhost:8080")
	if err := http.ListenAndServe(":8080", srv.Router); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
