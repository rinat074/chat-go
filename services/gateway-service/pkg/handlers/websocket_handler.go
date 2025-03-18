package handlers

import (
	"net/http"

	"github.com/rinat074/chat-go/services/gateway-service/pkg/clients"
)

// WebSocketHandler обрабатывает WebSocket соединения
type WebSocketHandler struct {
	chatClient *clients.ChatClient
}

// NewWebSocketHandler создает новый экземпляр WebSocketHandler
func NewWebSocketHandler(chatClient *clients.ChatClient) *WebSocketHandler {
	return &WebSocketHandler{
		chatClient: chatClient,
	}
}

// HandleWebSocket обрабатывает WebSocket соединения
func (h *WebSocketHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Заглушка для совместимости
}
