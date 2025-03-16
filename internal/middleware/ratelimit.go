package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

// RateLimiter ограничивает количество запросов
type RateLimiter struct {
	redisClient *redis.Client
	requests    int           // Количество запросов
	duration    time.Duration // Период ограничения
}

func NewRateLimiter(redisClient *redis.Client, requests int, duration time.Duration) *RateLimiter {
	return &RateLimiter{
		redisClient: redisClient,
		requests:    requests,
		duration:    duration,
	}
}

// Middleware возвращает middleware для ограничения запросов
func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Получаем IP клиента для идентификации
		ip := r.RemoteAddr

		ctx := r.Context()
		key := "rate_limit:" + ip

		// Увеличиваем счетчик для данного IP
		count, err := rl.redisClient.Incr(ctx, key).Result()
		if err != nil {
			http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
			return
		}

		// Если это первый запрос, устанавливаем время жизни ключа
		if count == 1 {
			rl.redisClient.Expire(ctx, key, rl.duration)
		}

		// Получаем оставшееся время для этого ключа
		ttl, err := rl.redisClient.TTL(ctx, key).Result()
		if err != nil {
			http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
			return
		}

		// Устанавливаем заголовки с информацией о Rate Limiting
		w.Header().Set("X-RateLimit-Limit", strconv.Itoa(rl.requests))
		w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(rl.requests-int(count)))
		w.Header().Set("X-RateLimit-Reset", strconv.Itoa(int(ttl.Seconds())))

		// Если превышен лимит запросов
		if count > int64(rl.requests) {
			w.Header().Set("Retry-After", strconv.Itoa(int(ttl.Seconds())))
			http.Error(w, "Слишком много запросов", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
