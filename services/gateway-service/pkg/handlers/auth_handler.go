package handlers

import (
	"github.com/rinat074/chat-go/services/gateway-service/pkg/clients"
)

// AuthHandler обрабатывает запросы аутентификации
type AuthHandler struct {
	authClient *clients.AuthClient
}

// NewAuthHandler создает новый экземпляр AuthHandler
func NewAuthHandler(authClient *clients.AuthClient) *AuthHandler {
	return &AuthHandler{
		authClient: authClient,
	}
}

// Добавьте необходимые методы-заглушки
