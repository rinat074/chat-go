package clients

import (
	"context"

	"github.com/rinat074/chat-go/proto/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// AuthClient клиент для взаимодействия с auth-service
type AuthClient struct {
	conn   *grpc.ClientConn
	client auth.AuthServiceClient
}

// NewAuthClient создает новый клиент для auth-service
func NewAuthClient(addr string) (*AuthClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := auth.NewAuthServiceClient(conn)
	return &AuthClient{
		conn:   conn,
		client: client,
	}, nil
}

// Register регистрирует нового пользователя
func (c *AuthClient) Register(ctx context.Context, username, email, password string) (*auth.AuthResponse, error) {
	return c.client.Register(ctx, &auth.RegisterRequest{
		Username: username,
		Email:    email,
		Password: password,
	})
}

// Login выполняет вход пользователя
func (c *AuthClient) Login(ctx context.Context, username, password, userAgent, ip string) (*auth.AuthResponse, error) {
	return c.client.Login(ctx, &auth.LoginRequest{
		Username:  username,
		Password:  password,
		UserAgent: userAgent,
		Ip:        ip,
	})
}

// RefreshToken обновляет токен
func (c *AuthClient) RefreshToken(ctx context.Context, refreshToken, userAgent, ip string) (*auth.TokenPair, error) {
	return c.client.RefreshToken(ctx, &auth.RefreshTokenRequest{
		RefreshToken: refreshToken,
		UserAgent:    userAgent,
		Ip:           ip,
	})
}

// Logout выполняет выход пользователя
func (c *AuthClient) Logout(ctx context.Context, refreshToken string) (*auth.LogoutResponse, error) {
	return c.client.Logout(ctx, &auth.LogoutRequest{
		RefreshToken: refreshToken,
	})
}

// ValidateToken проверяет валидность токена
func (c *AuthClient) ValidateToken(ctx context.Context, token string) (*auth.ValidateTokenResponse, error) {
	return c.client.ValidateToken(ctx, &auth.ValidateTokenRequest{
		Token: token,
	})
}

// Close закрывает соединение
func (c *AuthClient) Close() error {
	return c.conn.Close()
}
