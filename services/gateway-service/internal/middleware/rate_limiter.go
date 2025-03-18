package middleware

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	defaultWindow    = 60 * time.Second
	defaultMaxTokens = 100
)

// RateLimiter реализует ограничение скорости запросов
type RateLimiter struct {
	redisClient *redis.Client
	window      time.Duration
	maxTokens   int
}

// NewRateLimiter создает новый ограничитель скорости
func NewRateLimiter(redisClient *redis.Client) *RateLimiter {
	return &RateLimiter{
		redisClient: redisClient,
		window:      defaultWindow,
		maxTokens:   defaultMaxTokens,
	}
}

// Handler middleware для ограничения скорости запросов
func (rl *RateLimiter) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Получаем IP адрес клиента
		clientIP := r.RemoteAddr

		// Создаем ключ для Redis
		key := "rate_limit:" + clientIP

		// Проверяем количество запросов
		allowed, err := rl.isAllowed(r.Context(), key)
		if err != nil {
			// В случае ошибки Redis продолжаем выполнение
			next.ServeHTTP(w, r)
			return
		}

		if !allowed {
			w.Header().Set("Retry-After", strconv.Itoa(int(rl.window.Seconds())))
			http.Error(w, "Слишком много запросов", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// isAllowed проверяет, разрешен ли запрос
func (rl *RateLimiter) isAllowed(ctx context.Context, key string) (bool, error) {
	// Получаем текущее количество запросов
	count, err := rl.redisClient.Get(ctx, key).Int()
	if err != nil && err != redis.Nil {
		return true, err
	}

	if err == redis.Nil {
		// Если ключа нет, создаем его
		err = rl.redisClient.Set(ctx, key, 1, rl.window).Err()
		return true, err
	}

	// Проверяем, не превышен ли лимит
	if count >= rl.maxTokens {
		return false, nil
	}

	// Увеличиваем счетчик
	_, err = rl.redisClient.Incr(ctx, key).Result()
	return true, err
}
