package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/rinat074/chat-go/services/gateway-service/pkg/handlers"
	"github.com/rinat074/chat-go/services/gateway-service/pkg/middleware"
)

// Server представляет HTTP-сервер приложения
type Server struct {
	Router *chi.Mux
}

// NewServer создает новый экземпляр Server
func NewServer(
	authHandler *handlers.AuthHandler,
	chatHandler *handlers.ChatHandler,
	wsHandler *handlers.WebSocketHandler,
	authMiddleware *middleware.AuthMiddleware,
	rateLimiter *middleware.RateLimiter,
) *Server {
	r := chi.NewRouter()

	// Заглушка для совместимости
	return &Server{
		Router: r,
	}
}
