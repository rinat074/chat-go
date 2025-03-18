package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/rinat074/chat-go/pkg/clients"
	"github.com/rinat074/chat-go/pkg/logger"
	"github.com/rinat074/chat-go/pkg/middleware"
	"github.com/rinat074/chat-go/proto/chat"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Для разработки разрешаем все источники
	},
}

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
	// Получаем данные пользователя из контекста
	user, ok := r.Context().Value(middleware.UserContextKey).(middleware.UserData)
	if !ok {
		http.Error(w, "Требуется авторизация", http.StatusUnauthorized)
		return
	}

	// Обновляем соединение до WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error("Ошибка обновления до WebSocket", zap.Error(err))
		return
	}
	defer conn.Close()

	// Обработка сообщений от клиента
	for {
		_, msgBytes, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.Error("Ошибка WebSocket", zap.Error(err))
			}
			break
		}

		// Разбор сообщения
		var msgData struct {
			Type    string `json:"type"`
			Content string `json:"content"`
			GroupID *int64 `json:"group_id,omitempty"`
			UserID  *int64 `json:"user_id,omitempty"` // ID получателя для личных сообщений
		}

		if err := json.Unmarshal(msgBytes, &msgData); err != nil {
			logger.Error("Ошибка разбора сообщения", zap.Error(err))
			continue
		}

		// Создаем proto-сообщение
		msgType := chat.MessageType_PUBLIC
		switch msgData.Type {
		case "private":
			msgType = chat.MessageType_PRIVATE
		case "group":
			msgType = chat.MessageType_GROUP
		}

		protoMsg := &chat.Message{
			Type:      msgType,
			Content:   msgData.Content,
			UserId:    user.UserID,
			Username:  user.Username,
			CreatedAt: timestamppb.New(time.Now()),
		}

		if msgData.GroupID != nil {
			protoMsg.GroupId = msgData.GroupID
		}

		if msgData.UserID != nil {
			protoMsg.ReceiverId = msgData.UserID
		}

		// Отправляем сообщение в чат-сервис
		savedMsg, err := h.chatClient.SaveMessage(r.Context(), protoMsg)
		if err != nil {
			logger.Error("Ошибка сохранения сообщения", zap.Error(err))
			continue
		}

		// Отправляем подтверждение клиенту
		if err := conn.WriteJSON(savedMsg); err != nil {
			logger.Error("Ошибка отправки подтверждения", zap.Error(err))
		}
	}
}
