package server

import (
	"chat-app/internal/auth"
	"chat-app/internal/chat"
	customMiddleware "chat-app/internal/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	Router      *chi.Mux
	authHandler *auth.Handler
	chatHandler *chat.Handler
}

func NewServer(authHandler *auth.Handler, chatHandler *chat.Handler, authService *auth.Service, rateLimiter *customMiddleware.RateLimiter) *Server {
	s := &Server{
		Router:      chi.NewRouter(),
		authHandler: authHandler,
		chatHandler: chatHandler,
	}

	// Применяем общие middleware
	s.Router.Use(middleware.RequestID)
	s.Router.Use(middleware.RealIP)
	s.Router.Use(middleware.Logger)
	s.Router.Use(middleware.Recoverer)
	s.Router.Use(rateLimiter.Middleware) // Добавляем rate limiting

	// CSRF токен
	s.Router.Get("/api/csrf-token", customMiddleware.GetCSRFToken)

	// Маршруты аутентификации (без JWT-защиты)
	s.Router.Route("/api/auth", func(r chi.Router) {
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)
		r.Post("/refresh", authHandler.RefreshToken)
		r.Post("/logout", authHandler.Logout)
	})

	// Защищенные маршруты (с JWT-защитой и CSRF)
	s.Router.Group(func(r chi.Router) {
		r.Use(auth.Middleware(authService))
		r.Use(customMiddleware.CSRFMiddleware) // Добавляем CSRF защиту

		// WebSocket маршрут (для WebSocket не нужна CSRF защита)
		r.Get("/ws", chatHandler.WebSocketHandler)

		// API для чата
		r.Route("/api/chat", func(r chi.Router) {
			// Публичные сообщения
			r.Get("/messages/public", chatHandler.GetPublicMessages)

			// Личные сообщения
			r.Get("/messages/private/{userID}", chatHandler.GetPrivateMessages)

			// Групповые сообщения
			r.Get("/groups/{groupID}/messages", chatHandler.GetGroupMessages)

			// Управление группами
			r.Post("/groups", chatHandler.CreateGroup)
			r.Post("/groups/{groupID}/members", chatHandler.AddUserToGroup)
		})
	})

	return s
}
