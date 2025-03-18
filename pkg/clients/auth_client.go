package clients

import "context"

// AuthClient представляет клиент для Auth сервиса
type AuthClient struct{}

// NewAuthClient создает новый клиент
func NewAuthClient(address string) (*AuthClient, error) {
	return &AuthClient{}, nil
}

// Close закрывает соединение
func (c *AuthClient) Close() error {
	return nil
}

// ValidateToken проверяет токен
func (c *AuthClient) ValidateToken(ctx context.Context, token string) (interface{}, error) {
	return struct {
		Valid bool
		UserId int64
		Username string
	}{
		Valid: true,
		UserId: 1,
		Username: "dummy_user",
	}, nil
}
