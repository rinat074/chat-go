package service

import (
	"context"
	"time"

	"github.com/rinat074/chat-go/services/chat-service/internal/cache"
	"github.com/rinat074/chat-go/services/chat-service/internal/db"
)

// ChatService представляет сервис чата
type ChatService struct {
	db    *db.Database
	cache *cache.Cache
}

// NewChatService создает новый экземпляр сервиса чата
func NewChatService(db *db.Database, cache *cache.Cache) *ChatService {
	return &ChatService{
		db:    db,
		cache: cache,
	}
}

// Hub управляет всеми активными клиентами и сообщениями
type Hub struct {
	// Канал входящих сообщений от клиентов
	Broadcast chan Message
	service   *ChatService
}

// NewHub создает новый экземпляр хаба
func NewHub(service *ChatService) *Hub {
	return &Hub{
		Broadcast: make(chan Message, 256),
		service:   service,
	}
}

// Run запускает цикл обработки сообщений
func (h *Hub) Run() {
	for {
		select {
		case message := <-h.Broadcast:
			// Здесь будет логика обработки сообщений
			// Например, сохранение в базу данных и отправка клиентам
			_, _ = h.service.SaveMessage(context.Background(), message)
		}
	}
}

// SaveMessage сохраняет сообщение
func (s *ChatService) SaveMessage(ctx context.Context, msg Message) (*Message, error) {
	// Временная реализация
	msg.ID = time.Now().UnixNano()
	msg.CreatedAt = time.Now()
	return &msg, nil
}

// GetPublicMessages возвращает публичные сообщения
func (s *ChatService) GetPublicMessages(ctx context.Context, limit, offset int) ([]Message, error) {
	// Временная реализация
	return []Message{
		{
			ID:        1,
			Type:      PublicMessage,
			Content:   "Тестовое сообщение",
			UserID:    1,
			Username:  "user1",
			CreatedAt: time.Now(),
		},
	}, nil
}
