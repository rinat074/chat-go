package cache

import (
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
