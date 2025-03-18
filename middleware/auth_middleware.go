package middleware

import (
	"github.com/rinat074/chat-go/services/gateway-service/internal/clients"
	"github.com/rinat074/chat-go/services/gateway-service/internal/handlers"
)

// Ключ контекста
var userContextKey = handlers.UserContextKey()

type AuthMiddleware struct {
	authClient *clients.AuthClient
}

func NewAuthMiddleware(authClient *clients.AuthClient) *AuthMiddleware {
	return &AuthMiddleware{
		authClient: authClient,
	}
}
