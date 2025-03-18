package clients

import (
	"context"

	"github.com/rinat074/chat-go/proto/auth"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthClient struct {
	conn   *grpc.ClientConn
	client auth.AuthServiceClient
}

func NewAuthClient(address string) (*AuthClient, error) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := auth.NewAuthServiceClient(conn)
	return &AuthClient{
		conn:   conn,
		client: client,
	}, nil
}

func (c *AuthClient) Close() error {
	return c.conn.Close()
}

func (c *AuthClient) Register(ctx context.Context, username, email, password string) (*auth.AuthResponse, error) {
	return c.client.Register(ctx, &auth.RegisterRequest{
		Username: username,
		Email:    email,
		Password: password,
	})
}

func (c *AuthClient) Login(ctx context.Context, username, password, userAgent, ip string) (*auth.AuthResponse, error) {
	return c.client.Login(ctx, &auth.LoginRequest{
		Username:  username,
		Password:  password,
		UserAgent: userAgent,
		Ip:        ip,
	})
}

func (c *AuthClient) ValidateToken(ctx context.Context, token string) (*auth.ValidateTokenResponse, error) {
	return c.client.ValidateToken(ctx, &auth.ValidateTokenRequest{
		Token: token,
	})
}

// И другие методы для взаимодействия с auth-service
