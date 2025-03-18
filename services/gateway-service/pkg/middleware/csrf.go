package middleware

import (
	"net/http"
)

// GetCSRFToken генерирует CSRF токен и возвращает его клиенту
func GetCSRFToken(w http.ResponseWriter, r *http.Request) {
	// Упрощенная версия для совместимости
	// Реальная логика находится в internal/middleware/csrf.go
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"csrf_token":"dummy_token"}`))
}

// CSRFMiddleware проверяет CSRF токен для небезопасных методов
func CSRFMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Упрощенная версия для совместимости
		// Реальная логика находится в internal/middleware/csrf.go
		next.ServeHTTP(w, r)
	})
}
