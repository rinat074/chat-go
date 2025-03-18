package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/rinat074/chat-go/services/gateway-service/internal/clients"
	"github.com/rinat074/chat-go/services/gateway-service/pkg/logger"
)

type AuthMiddleware struct {
	clients   *clients.ServiceClients
	log       logger.Logger
	jwtSecret string
}

func NewAuthMiddleware(clients *clients.ServiceClients, log logger.Logger, jwtSecret string) *AuthMiddleware {
	return &AuthMiddleware{
		clients:   clients,
		log:       log,
		jwtSecret: jwtSecret,
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
			m.log.Warn("отсутствует заголовок авторизации", "path", r.URL.Path)
			http.Error(w, "Требуется авторизация", http.StatusUnauthorized)
			return
		}

		// Извлекаем токен из формата "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			m.log.Warn("неверный формат токена", "auth_header", authHeader)
			http.Error(w, "Неверный формат токена", http.StatusUnauthorized)
			return
		}
		token := parts[1]

		// Проверяем токен через auth-service
		resp, err := m.clients.AuthClient.ValidateToken(r.Context(), token)
		if err != nil {
			m.log.Error("ошибка проверки токена", "error", err)
			http.Error(w, "Ошибка проверки токена", http.StatusUnauthorized)
			return
		}

		if !resp.Valid {
			m.log.Warn("недействительный токен", "token", token[:10]+"...")
			http.Error(w, "Недействительный токен", http.StatusUnauthorized)
			return
		}

		// Добавляем информацию о пользователе в контекст
		ctx := context.WithValue(r.Context(), "userID", resp.UserId)
		ctx = context.WithValue(ctx, "username", resp.Username)

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
		"/swagger",
	}

	for _, route := range publicRoutes {
		if strings.HasPrefix(path, route) {
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
