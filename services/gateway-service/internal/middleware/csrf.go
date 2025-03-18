package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"
)

// Использовать безопасное хранилище в production (Redis, БД и т.д.)
var csrfTokens = make(map[string]time.Time)

// GetCSRFToken генерирует CSRF токен и возвращает его клиенту
func GetCSRFToken(w http.ResponseWriter, r *http.Request) {
	// Генерация случайного токена
	token := generateToken()

	// Сохранение токена (в production лучше использовать Redis)
	csrfTokens[token] = time.Now().Add(24 * time.Hour)

	// Очистка устаревших токенов
	cleanExpiredTokens()

	// Возвращаем токен клиенту
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"csrf_token":"` + token + `"}`))
}

// CSRFMiddleware проверяет CSRF токен для небезопасных методов
func CSRFMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем только для небезопасных методов
		if r.Method != "GET" && r.Method != "HEAD" && r.Method != "OPTIONS" {
			token := r.Header.Get("X-CSRF-Token")
			if token == "" {
				http.Error(w, "CSRF токен отсутствует", http.StatusForbidden)
				return
			}

			// Проверяем существование и действительность токена
			expiry, exists := csrfTokens[token]
			if !exists || time.Now().After(expiry) {
				http.Error(w, "Недействительный или устаревший CSRF токен", http.StatusForbidden)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

// generateToken генерирует случайный токен
func generateToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

// cleanExpiredTokens удаляет устаревшие токены
func cleanExpiredTokens() {
	now := time.Now()
	for token, expiry := range csrfTokens {
		if now.After(expiry) {
			delete(csrfTokens, token)
		}
	}
}
