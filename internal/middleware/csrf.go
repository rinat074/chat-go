package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"
)

const (
	csrfTokenHeader = "X-CSRF-Token"
	csrfTokenCookie = "csrf_token"
	csrfTokenExpiry = 24 * time.Hour
)

// GenerateCSRFToken создает новый CSRF токен
func GenerateCSRFToken() (string, error) {
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(tokenBytes), nil
}

// CSRFMiddleware добавляет CSRF защиту
func CSRFMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Пропускаем GET, HEAD, OPTIONS, TRACE
		if r.Method == http.MethodGet || r.Method == http.MethodHead ||
			r.Method == http.MethodOptions || r.Method == http.MethodTrace {
			next.ServeHTTP(w, r)
			return
		}

		// Проверяем наличие токена в заголовке
		csrfToken := r.Header.Get(csrfTokenHeader)
		if csrfToken == "" {
			http.Error(w, "CSRF токен отсутствует", http.StatusForbidden)
			return
		}

		// Получаем токен из куки
		cookie, err := r.Cookie(csrfTokenCookie)
		if err != nil || cookie.Value == "" {
			http.Error(w, "CSRF токен недействителен", http.StatusForbidden)
			return
		}

		// Сравниваем токены
		if csrfToken != cookie.Value {
			http.Error(w, "CSRF токен недействителен", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// GetCSRFToken устанавливает CSRF токен в куки если его нет и возвращает токен
func GetCSRFToken(w http.ResponseWriter, r *http.Request) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var token string
		cookie, err := r.Cookie(csrfTokenCookie)

		if err != nil || cookie.Value == "" {
			newToken, err := GenerateCSRFToken()
			if err != nil {
				http.Error(w, "Не удалось создать CSRF токен", http.StatusInternalServerError)
				return
			}

			token = newToken
			http.SetCookie(w, &http.Cookie{
				Name:     csrfTokenCookie,
				Value:    token,
				Path:     "/",
				HttpOnly: true,
				SameSite: http.SameSiteStrictMode,
				Secure:   true, // в продакшене
				MaxAge:   int(csrfTokenExpiry.Seconds()),
			})
		} else {
			token = cookie.Value
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"token": "` + token + `"}`))
	}
}
