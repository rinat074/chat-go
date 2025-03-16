package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

type contextKey string

const UserContextKey contextKey = "user"

type UserClaims struct {
	UserID   int64
	Username string
}

func Middleware(service *Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Пропускаем пути авторизации
			if r.URL.Path == "/api/auth/register" || r.URL.Path == "/api/auth/login" {
				next.ServeHTTP(w, r)
				return
			}

			// Получаем токен из заголовка
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Требуется авторизация", http.StatusUnauthorized)
				return
			}

			// Обычно токен в формате "Bearer <token>"
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Неверный формат токена авторизации", http.StatusUnauthorized)
				return
			}

			tokenString := parts[1]
			token, err := service.ValidateToken(tokenString)
			if err != nil || !token.Valid {
				http.Error(w, "Недействительный токен", http.StatusUnauthorized)
				return
			}

			// Извлекаем данные пользователя из токена
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, "Недействительный токен", http.StatusUnauthorized)
				return
			}

			userID := int64(claims["user_id"].(float64))
			username := claims["username"].(string)

			// Добавляем информацию о пользователе в контекст
			ctx := context.WithValue(r.Context(), UserContextKey, UserClaims{
				UserID:   userID,
				Username: username,
			})

			// Продолжаем обработку запроса
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
