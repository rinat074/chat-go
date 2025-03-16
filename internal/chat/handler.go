package chat

import (
	"encoding/json"
	"net/http"
	"strconv"

	"chat-app/internal/auth"
	"chat-app/internal/models"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Разрешить все источники в целях разработки
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Handler struct {
	hub     *Hub
	service *Service
}

func NewHandler(hub *Hub, service *Service) *Handler {
	return &Handler{
		hub:     hub,
		service: service,
	}
}

// WebSocketHandler обрабатывает WebSocket соединения
func (h *Handler) WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем данные пользователя из контекста
	userClaims, ok := r.Context().Value(auth.UserContextKey).(auth.UserClaims)
	if !ok {
		http.Error(w, "Требуется авторизация", http.StatusUnauthorized)
		return
	}

	// Обновляем соединение до WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Не удалось обновить до WebSocket", http.StatusInternalServerError)
		return
	}

	// Создаем нового клиента
	client := &Client{
		hub:      h.hub,
		conn:     conn,
		send:     make(chan models.Message, 256),
		userID:   userClaims.UserID,
		username: userClaims.Username,
	}

	// Регистрируем клиента в хабе
	client.hub.register <- client

	// Запускаем горутины для чтения и записи сообщений
	go client.writePump()
	go client.readPump()
}

// GetChatHistory возвращает историю сообщений чата (исправление ошибки)
func (h *Handler) GetChatHistory(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 50
	offset := 0

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	// Используем GetPublicMessages вместо неопределенного GetMessages
	messages, err := h.service.GetPublicMessages(r.Context(), limit, offset)
	if err != nil {
		http.Error(w, "Ошибка при получении истории сообщений", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}
