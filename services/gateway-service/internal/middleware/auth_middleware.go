package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/rinat074/chat-go/services/gateway-service/internal/clients"
)

type AuthMiddleware struct {
	authClient *clients.AuthClient
}

func NewAuthMiddleware(authClient *clients.AuthClient) *AuthMiddleware {
	return &AuthMiddleware{
		authClient: authClient,
	}
}

func (m *AuthMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Исключаем открытые маршруты
		if isPublicRoute(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		// Получаем токен из заголовка
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Требуется авторизация", http.StatusUnauthorized)
			return
		}

		// Извлекаем токен из формата "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Неверный формат токена", http.StatusUnauthorized)
			return
		}
		token := parts[1]

		// Проверяем токен через auth-service
		resp, err := m.authClient.ValidateToken(r.Context(), token)
		if err != nil || !resp.Valid {
			http.Error(w, "Недействительный токен", http.StatusUnauthorized)
			return
		}

		// Добавляем информацию о пользователе в контекст
		ctx := context.WithValue(r.Context(), userContextKey, userData{
			UserID:   resp.UserId,
			Username: resp.Username,
		})

		// Продолжаем запрос с обновленным контекстом
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Определяем открытые маршруты, не требующие авторизации
func isPublicRoute(path string) bool {
	publicRoutes := []string{
		"/api/auth/register",
		"/api/auth/login",
		"/api/auth/refresh",
		"/api/csrf-token",
	}

	for _, route := range publicRoutes {
		if path == route {
			return true
		}
	}
	return false
}

// Ключ и структура для данных пользователя в контексте
type contextKey string

const userContextKey contextKey = "user"

type userData struct {
	UserID   int64
	Username string
}
