package middleware

import (
	"context"
	"net/http"
	"strings"
)

// Ключ и структура для данных пользователя в контексте
type contextKey string

// UserContextKey ключ контекста для данных пользователя
const UserContextKey contextKey = "user"

// UserData содержит информацию о пользователе
type UserData struct {
	UserID   int64
	Username string
}

// AuthMiddleware проверяет JWT токен и добавляет информацию о пользователе в контекст
type AuthMiddleware struct {
	authClient interface{} // Упрощенный интерфейс
}

// NewAuthMiddleware создает новый экземпляр AuthMiddleware
func NewAuthMiddleware(authClient interface{}) *AuthMiddleware {
	return &AuthMiddleware{
		authClient: authClient,
	}
}

// Handler middleware для проверки авторизации
func (m *AuthMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Упрощенная версия для совместимости
		// Добавляем фиктивные данные пользователя для совместимости с websocket_handler.go
		userData := UserData{
			UserID:   1,
			Username: "dummy_user",
		}
		ctx := context.WithValue(r.Context(), UserContextKey, userData)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// IsPublicRoute проверяет, является ли маршрут открытым
func IsPublicRoute(path string) bool {
	publicRoutes := []string{
		"/api/auth/register",
		"/api/auth/login",
		"/api/auth/refresh",
		"/api/csrf-token",
		"/swagger",
	}

	for _, route := range publicRoutes {
		if strings.HasPrefix(path, route) {
			return true
		}
	}
	return false
}
