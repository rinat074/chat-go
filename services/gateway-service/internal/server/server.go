package server

import (
	"github.com/rinat074/chat-go/services/gateway-service/internal/handlers"
	"github.com/rinat074/chat-go/services/gateway-service/internal/middleware"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	Router *chi.Mux
}

func NewServer(
	authHandler *handlers.AuthHandler,
	chatHandler *handlers.ChatHandler,
	wsHandler *handlers.WebSocketHandler,
	authMiddleware *middleware.AuthMiddleware,
	rateLimiter *middleware.RateLimiter,
) *Server {
	router := chi.NewRouter()

	// Глобальные middleware
	router.Use(chiMiddleware.RequestID)
	router.Use(chiMiddleware.RealIP)
	router.Use(chiMiddleware.Logger)
	router.Use(chiMiddleware.Recoverer)
	router.Use(rateLimiter.Handler)

	// CSRF-токен
	router.Get("/api/csrf-token", middleware.GetCSRFToken)

	// Маршруты аутентификации (без JWT-защиты)
	router.Route("/api/auth", func(r chi.Router) {
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)
		r.Post("/refresh", authHandler.RefreshToken)
		r.Post("/logout", authHandler.Logout)
	})

	// Защищенные маршруты
	router.Group(func(r chi.Router) {
		r.Use(authMiddleware.Handler)
		r.Use(middleware.CSRFMiddleware)

		// WebSocket
		r.Get("/ws", wsHandler.HandleWebSocket)

		// Маршруты чата
		r.Route("/api/chat", func(r chi.Router) {
			r.Get("/messages/public", chatHandler.GetPublicMessages)
			r.Get("/messages/private/{userID}", chatHandler.GetPrivateMessages)
			r.Get("/groups/{groupID}/messages", chatHandler.GetGroupMessages)
			r.Post("/groups", chatHandler.CreateGroup)
			r.Post("/groups/{groupID}/members", chatHandler.AddUserToGroup)
		})
	})

	return &Server{
		Router: router,
	}
}
