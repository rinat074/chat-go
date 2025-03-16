package service

import (
	"context"
	"time"

	"chat-app/services/chat-service/internal/cache"
	"chat-app/services/chat-service/internal/db"
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
