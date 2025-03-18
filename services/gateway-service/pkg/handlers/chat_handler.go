package handlers

import (
	"github.com/rinat074/chat-go/services/gateway-service/pkg/clients"
)

// ChatHandler обрабатывает запросы к чату
type ChatHandler struct {
	chatClient *clients.ChatClient
}

// NewChatHandler создает новый экземпляр ChatHandler
func NewChatHandler(chatClient *clients.ChatClient) *ChatHandler {
	return &ChatHandler{
		chatClient: chatClient,
	}
}

// Добавьте необходимые методы-заглушки
