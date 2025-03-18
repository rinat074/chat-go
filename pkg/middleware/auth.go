package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/rinat074/chat-go/pkg/clients"
	"github.com/rinat074/chat-go/pkg/logger"
	"go.uber.org/zap"
)

// ContextKey - тип для ключей контекста.
type ContextKey string

// UserContextKey - ключ для хранения данных пользователя в контексте.
const UserContextKey ContextKey = "user"

// UserData - данные пользователя.
type UserData struct {
	UserID   int64
	Username string
}

// AuthMiddleware - middleware для аутентификации запросов.
type AuthMiddleware struct {
	authClient *clients.AuthClient
}

// NewAuthMiddleware создает новый экземпляр AuthMiddleware.
func NewAuthMiddleware(authClient *clients.AuthClient) *AuthMiddleware {
	return &AuthMiddleware{
		authClient: authClient,
	}
}

// Handler обрабатывает аутентификацию HTTP запросов.
func (m *AuthMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Получаем токен из заголовка Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Требуется авторизация", http.StatusUnauthorized)
			return
		}

		// Проверяем формат Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Неверный формат токена авторизации", http.StatusUnauthorized)
			return
		}

		token := parts[1]

		// Валидируем токен через auth-service
		resp, err := m.authClient.ValidateToken(r.Context(), token)
		if err != nil {
			logger.Error("Ошибка валидации токена", zap.Error(err))
			http.Error(w, "Ошибка авторизации", http.StatusUnauthorized)
			return
		}

		if !resp.Valid {
			http.Error(w, "Невалидный токен", http.StatusUnauthorized)
			return
		}

		// Токен валиден, добавляем информацию о пользователе в контекст
		userData := UserData{
			UserID:   resp.UserId,
			Username: resp.Username,
		}

		ctx := context.WithValue(r.Context(), UserContextKey, userData)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
