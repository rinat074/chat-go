package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/rinat074/chat-go/pkg/clients"
	"github.com/rinat074/chat-go/pkg/logger"
	"go.uber.org/zap"
)

// RateLimiter middleware для ограничения частоты запросов
type RateLimiter struct {
	redisClient *clients.RedisClient
	limit       int
	window      time.Duration
}

// NewRateLimiter создает новый экземпляр RateLimiter
func NewRateLimiter(redisClient *clients.RedisClient, limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		redisClient: redisClient,
		limit:       limit,
		window:      window,
	}
}

// Handler обрабатывает ограничение частоты запросов
func (m *RateLimiter) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Используем IP клиента как идентификатор
		ip := r.RemoteAddr
		key := fmt.Sprintf("ratelimit:%s", ip)

		// Проверяем, существует ли ключ
		exists, err := m.redisClient.Exists(r.Context(), key)
		if err != nil {
			logger.Error("Ошибка при проверке существования ключа", zap.Error(err))
			// В случае ошибки пропускаем запрос
			next.ServeHTTP(w, r)
			return
		}

		if !exists {
			// Создаем счетчик
			err = m.redisClient.Set(r.Context(), key, 1, m.window)
			if err != nil {
				logger.Error("Ошибка при создании счетчика", zap.Error(err))
				next.ServeHTTP(w, r)
				return
			}
		} else {
			// Получаем текущее значение счетчика
			val, err := m.redisClient.Get(r.Context(), key)
			if err != nil {
				logger.Error("Ошибка при получении счетчика", zap.Error(err))
				next.ServeHTTP(w, r)
				return
			}

			// Преобразуем строку в число
			count, err := strconv.Atoi(val)
			if err != nil {
				logger.Error("Ошибка при преобразовании счетчика", zap.Error(err))
				next.ServeHTTP(w, r)
				return
			}

			// Проверяем, не превышен ли лимит
			if count >= m.limit {
				logger.Warn("Превышен лимит запросов", zap.String("ip", ip), zap.Int("count", count), zap.Int("limit", m.limit))
				w.Header().Set("Retry-After", fmt.Sprintf("%d", int(m.window.Seconds())))
				http.Error(w, "Слишком много запросов", http.StatusTooManyRequests)
				return
			}

			// Увеличиваем счетчик
			_, err = m.redisClient.Incr(r.Context(), key)
			if err != nil {
				logger.Error("Ошибка при увеличении счетчика", zap.Error(err))
				next.ServeHTTP(w, r)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

// GetCSRFToken генерирует и возвращает CSRF-токен
func GetCSRFToken(w http.ResponseWriter, r *http.Request) {
	// Генерация случайного токена
	token := fmt.Sprintf("%d", time.Now().UnixNano())

	// Установка токена в куки
	cookie := http.Cookie{
		Name:     "X-CSRF-Token",
		Value:    token,
		Path:     "/",
		HttpOnly: false, // Чтобы был доступен из JavaScript
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   3600, // 1 час
	}
	http.SetCookie(w, &cookie)

	// Отправка токена в ответе
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"token":"` + token + `"}`))
}

// CSRFMiddleware обеспечивает защиту от CSRF-атак
func CSRFMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Только для небезопасных методов
		if r.Method != http.MethodGet && r.Method != http.MethodHead && r.Method != http.MethodOptions {
			// Получаем токен из cookie
			cookie, err := r.Cookie("X-CSRF-Token")
			if err != nil {
				http.Error(w, "CSRF-токен отсутствует", http.StatusForbidden)
				return
			}

			// Получаем токен из заголовка
			headerToken := r.Header.Get("X-CSRF-Token")
			if headerToken == "" {
				http.Error(w, "CSRF-токен отсутствует в заголовке", http.StatusForbidden)
				return
			}

			// Сравниваем токены
			if cookie.Value != headerToken {
				http.Error(w, "CSRF-токен недействителен", http.StatusForbidden)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
