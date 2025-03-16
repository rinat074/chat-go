package clients

import (
	"context"

	"chat-app/proto/chat"

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

func (c *ChatClient) GetPublicMessages(ctx context.Context, limit, offset int32) (*chat.MessagesResponse, error) {
	return c.client.GetPublicMessages(ctx, &chat.GetMessagesRequest{
		Limit:  limit,
		Offset: offset,
	})
}

// Дополнительные методы для других функций чата
