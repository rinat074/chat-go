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

// GetPublicMessages получает публичные сообщения
func (c *ChatClient) GetPublicMessages(ctx context.Context, limit, offset int) (*chat.MessagesResponse, error) {
	return c.client.GetPublicMessages(ctx, &chat.GetMessagesRequest{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
}

// GetPrivateMessages получает личные сообщения между двумя пользователями
func (c *ChatClient) GetPrivateMessages(ctx context.Context, userID, otherUserID int64, limit, offset int) (*chat.MessagesResponse, error) {
	return c.client.GetPrivateMessages(ctx, &chat.GetPrivateMessagesRequest{
		UserId:      userID,
		OtherUserId: otherUserID,
		Limit:       int32(limit),
		Offset:      int32(offset),
	})
}

// GetGroupMessages получает сообщения группы
func (c *ChatClient) GetGroupMessages(ctx context.Context, groupID, userID int64, limit, offset int) (*chat.MessagesResponse, error) {
	return c.client.GetGroupMessages(ctx, &chat.GetGroupMessagesRequest{
		GroupId: groupID,
		UserId:  userID,
		Limit:   int32(limit),
		Offset:  int32(offset),
	})
}

// CreateGroup создает новую группу
func (c *ChatClient) CreateGroup(ctx context.Context, name, description string, ownerID int64) (*chat.Group, error) {
	return c.client.CreateGroup(ctx, &chat.CreateGroupRequest{
		Name:        name,
		Description: description,
		OwnerId:     ownerID,
	})
}

// AddUserToGroup добавляет пользователя в группу
func (c *ChatClient) AddUserToGroup(ctx context.Context, groupID, userID, adminID int64) (*chat.AddUserToGroupResponse, error) {
	return c.client.AddUserToGroup(ctx, &chat.AddUserToGroupRequest{
		GroupId: groupID,
		UserId:  userID,
		AdminId: adminID,
	})
}
