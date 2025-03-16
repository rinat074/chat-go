package cache

import (
	"context"
	"encoding/json"
	"time"

	"chat-app/internal/models"

	"github.com/go-redis/redis/v8"
)

// Cache представляет кэш Redis
type Cache struct {
	Client *redis.Client
}

// NewRedisCache создает новый экземпляр кэша Redis
func NewRedisCache(addr string) *Cache {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // в продакшене используйте пароль
		DB:       0,
	})

	return &Cache{
		Client: client,
	}
}

// SetMessages кэширует сообщения для определенного ключа
func (c *Cache) SetMessages(ctx context.Context, key string, messages []models.Message, ttl time.Duration) error {
	data, err := json.Marshal(messages)
	if err != nil {
		return err
	}

	return c.Client.Set(ctx, key, data, ttl).Err()
}

// GetMessages получает кэшированные сообщения
func (c *Cache) GetMessages(ctx context.Context, key string) ([]models.Message, error) {
	data, err := c.Client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	var messages []models.Message
	if err := json.Unmarshal(data, &messages); err != nil {
		return nil, err
	}

	return messages, nil
}

// DeleteByPattern удаляет ключи по паттерну
func (c *Cache) DeleteByPattern(ctx context.Context, pattern string) error {
	keys, err := c.Client.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		return c.Client.Del(ctx, keys...).Err()
	}

	return nil
}
