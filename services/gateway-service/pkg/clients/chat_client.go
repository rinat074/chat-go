package clients

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

// Добавьте необходимые методы-заглушки
