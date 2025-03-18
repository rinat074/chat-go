package clients

import (
	"context"

	"github.com/rinat074/chat-go/proto/chat"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ChatClient struct {
	conn   *grpc.ClientConn
	client chat.ChatServiceClient
}

func NewChatClient(address string) (*ChatClient, error) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := chat.NewChatServiceClient(conn)
	return &ChatClient{
		conn:   conn,
		client: client,
	}, nil
}

func (c *ChatClient) Close() error {
	return c.conn.Close()
}

// Добавляем методы для взаимодействия с чат-сервисом
func (c *ChatClient) SaveMessage(ctx context.Context, message *chat.Message) (*chat.Message, error) {
	return c.client.SaveMessage(ctx, message)
}

func (c *ChatClient) GetPublicMessages(ctx context.Context, limit, offset int) (*chat.MessagesResponse, error) {
	return c.client.GetPublicMessages(ctx, &chat.GetMessagesRequest{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
}

func (c *ChatClient) GetPrivateMessages(ctx context.Context, userId, receiverId int64, limit, offset int) (*chat.MessagesResponse, error) {
	return c.client.GetPrivateMessages(ctx, &chat.GetPrivateMessagesRequest{
		UserId:      userId,
		OtherUserId: receiverId,
		Limit:       int32(limit),
		Offset:      int32(offset),
	})
}

func (c *ChatClient) GetGroupMessages(ctx context.Context, groupId int64, limit, offset int) (*chat.MessagesResponse, error) {
	return c.client.GetGroupMessages(ctx, &chat.GetGroupMessagesRequest{
		GroupId: groupId,
		UserId:  0,
		Limit:   int32(limit),
		Offset:  int32(offset),
	})
}

func (c *ChatClient) CreateGroup(ctx context.Context, name, description string, creatorId int64, memberIds []int64) (*chat.Group, error) {
	return c.client.CreateGroup(ctx, &chat.CreateGroupRequest{
		Name:        name,
		Description: description,
		OwnerId:     creatorId,
	})
}

func (c *ChatClient) AddUserToGroup(ctx context.Context, groupId, userId, requesterId int64) error {
	_, err := c.client.AddUserToGroup(ctx, &chat.AddUserToGroupRequest{
		GroupId: groupId,
		UserId:  userId,
		AdminId: requesterId,
	})
	return err
}

// Дополнительные методы для других функций чата
