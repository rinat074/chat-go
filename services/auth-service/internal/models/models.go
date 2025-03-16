package models

import (
	"time"
)

// User представляет пользователя
type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // Не возвращаем пароль в JSON
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// LoginRequest запрос на вход
type LoginRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	UserAgent string `json:"user_agent"`
	IP        string `json:"ip"`
}

// RegisterRequest запрос на регистрацию
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// TokenPair пара токенов (access и refresh)
type TokenPair struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// RefreshSession сессия refresh токена
type RefreshSession struct {
	ID           int64     `json:"id"`
	UserID       int64     `json:"user_id"`
	RefreshToken string    `json:"refresh_token"`
	UserAgent    string    `json:"user_agent"`
	IP           string    `json:"ip"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
}

// AuthResponse ответ на аутентификацию
type AuthResponse struct {
	Tokens *TokenPair `json:"tokens,omitempty"`
	Token  string     `json:"token,omitempty"` // Для обратной совместимости
	User   User       `json:"user"`
}
