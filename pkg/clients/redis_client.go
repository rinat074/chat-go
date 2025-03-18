package clients

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisClient обертка для Redis клиента
type RedisClient struct {
	client *redis.Client
}

// NewRedisClient создает новый клиент Redis
func NewRedisClient(addr string) *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // по умолчанию нет пароля
		DB:       0,  // использовать базу данных по умолчанию
	})

	return &RedisClient{
		client: client,
	}
}

// Get получает значение по ключу
func (c *RedisClient) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}

// Set устанавливает значение по ключу
func (c *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return c.client.Set(ctx, key, value, expiration).Err()
}

// Del удаляет ключ
func (c *RedisClient) Del(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

// Incr инкрементирует значение ключа
func (c *RedisClient) Incr(ctx context.Context, key string) (int64, error) {
	return c.client.Incr(ctx, key).Result()
}

// Exists проверяет существование ключа
func (c *RedisClient) Exists(ctx context.Context, key string) (bool, error) {
	result, err := c.client.Exists(ctx, key).Result()
	return result > 0, err
}

// Close закрывает соединение
func (c *RedisClient) Close() error {
	return c.client.Close()
}
