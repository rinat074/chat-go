package clients

import "context"

// RedisClient - клиент для работы с Redis
type RedisClient struct{}

// NewRedisClient создает новый клиент Redis
func NewRedisClient(url string) *RedisClient {
	return &RedisClient{}
}

// Exists проверяет существование ключа
func (r *RedisClient) Exists(ctx context.Context, key string) (bool, error) {
	// Заглушка для совместимости
	return false, nil
}

// Set устанавливает значение ключа
func (r *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration interface{}) error {
	// Заглушка для совместимости
	return nil
}

// Get получает значение по ключу
func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	// Заглушка для совместимости
	return "0", nil
}

// Incr увеличивает значение ключа на 1
func (r *RedisClient) Incr(ctx context.Context, key string) (int64, error) {
	// Заглушка для совместимости
	return 1, nil
}
