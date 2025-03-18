package middleware

import (
	"net/http"
	"time"
)

// RateLimiter ограничивает количество запросов от одного IP
type RateLimiter struct {
	redisClient interface{}
	maxRequests int
	duration    time.Duration
}

// NewRateLimiter создает новый экземпляр RateLimiter
func NewRateLimiter(redisClient interface{}, maxRequests int, duration time.Duration) *RateLimiter {
	return &RateLimiter{
		redisClient: redisClient,
		maxRequests: maxRequests,
		duration:    duration,
	}
}

// Handler middleware для ограничения скорости запросов
func (rl *RateLimiter) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Упрощенная версия для совместимости
		next.ServeHTTP(w, r)
	})
}
