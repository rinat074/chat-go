package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/rinat074/chat-go/proto/chat"
	"github.com/rinat074/chat-go/services/gateway-service/internal/clients"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Для разработки разрешаем все источники
	},
}

type WebSocketHandler struct {
	chatClient *clients.ChatClient
	// Можно добавить hub для локального управления подключениями
}

func NewWebSocketHandler(chatClient *clients.ChatClient) *WebSocketHandler {
	return &WebSocketHandler{
		chatClient: chatClient,
	}
}

func (h *WebSocketHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Получаем данные пользователя из контекста
	user, ok := r.Context().Value(userContextKey).(userData)
	if !ok {
		http.Error(w, "Требуется авторизация", http.StatusUnauthorized)
		return
	}

	// Обновляем соединение до WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Ошибка обновления до WebSocket: %v", err)
		return
	}
	defer conn.Close()

	// Обработка сообщений от клиента
	for {
		_, msgBytes, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Ошибка WebSocket: %v", err)
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
			log.Printf("Ошибка разбора сообщения: %v", err)
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
			log.Printf("Ошибка сохранения сообщения: %v", err)
			continue
		}

		// Отправляем подтверждение клиенту
		if err := conn.WriteJSON(savedMsg); err != nil {
			log.Printf("Ошибка отправки подтверждения: %v", err)
		}
	}
}
