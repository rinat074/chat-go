package clients

import (
	"context"

	"github.com/rinat074/chat-go/proto/chat"
)

// ChatClient представляет клиент для Chat сервиса
type ChatClient struct{}

// NewChatClient создает новый клиент для Chat сервиса
func NewChatClient(address string) (*ChatClient, error) {
	return &ChatClient{}, nil
}

// Close закрывает соединение с сервером
func (c *ChatClient) Close() error {
	return nil
}

// SaveMessage сохраняет сообщение в чат-сервисе
func (c *ChatClient) SaveMessage(ctx context.Context, message *chat.Message) (*chat.Message, error) {
	// Возвращаем то же сообщение как заглушка
	return message, nil
}

// Добавьте необходимые методы-заглушки
