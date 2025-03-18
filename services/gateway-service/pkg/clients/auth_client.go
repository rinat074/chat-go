package clients

// AuthClient представляет клиент для Auth сервиса
type AuthClient struct{}

// NewAuthClient создает новый клиент для Auth сервиса
func NewAuthClient(address string) (*AuthClient, error) {
	return &AuthClient{}, nil
}

// Close закрывает соединение с сервером
func (c *AuthClient) Close() error {
	return nil
}

// Добавьте необходимые методы-заглушки
