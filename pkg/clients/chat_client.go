package clients

import (
	"context"

	"github.com/rinat074/chat-go/proto/chat"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ChatClient клиент для взаимодействия с chat-service
type ChatClient struct {
	conn   *grpc.ClientConn
	client chat.ChatServiceClient
}

// NewChatClient создает новый клиент для chat-service
func NewChatClient(addr string) (*ChatClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := chat.NewChatServiceClient(conn)
	return &ChatClient{
		conn:   conn,
		client: client,
	}, nil
}

// SaveMessage сохраняет сообщение
func (c *ChatClient) SaveMessage(ctx context.Context, msg *chat.Message) (*chat.Message, error) {
	return c.client.SaveMessage(ctx, msg)
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

// Close закрывает соединение
func (c *ChatClient) Close() error {
	return c.conn.Close()
}
