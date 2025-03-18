package clients

import (
	"github.com/redis/go-redis/v9"
)

// ServiceClients содержит всех клиентов для сервисов
type ServiceClients struct {
	AuthClient  *AuthClient
	ChatClient  *ChatClient
	RedisClient *redis.Client
}

// NewServiceClients создает новый экземпляр клиентов сервисов
func NewServiceClients(
	authClient *AuthClient,
	chatClient *ChatClient,
	redisClient *redis.Client,
) *ServiceClients {
	return &ServiceClients{
		AuthClient:  authClient,
		ChatClient:  chatClient,
		RedisClient: redisClient,
	}
}
