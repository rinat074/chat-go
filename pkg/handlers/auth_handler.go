package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/rinat074/chat-go/pkg/logger"
	"github.com/rinat074/chat-go/services/gateway-service/pkg/clients"
	"go.uber.org/zap"
)

// AuthHandler обрабатывает запросы аутентификации
type AuthHandler struct {
	authClient *clients.AuthClient
}

// NewAuthHandler создает новый экземпляр AuthHandler
func NewAuthHandler(authClient *clients.AuthClient) *AuthHandler {
	return &AuthHandler{
		authClient: authClient,
	}
}

// RegisterRequest структура для запроса регистрации
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest структура для запроса входа
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// RefreshTokenRequest структура для запроса обновления токена
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// LogoutRequest структура для запроса выхода
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// Register обрабатывает запрос на регистрацию нового пользователя
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	// Декодирование запроса
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Ошибка разбора запроса", http.StatusBadRequest)
		return
	}

	// Валидация данных
	if req.Username == "" || req.Email == "" || req.Password == "" {
		http.Error(w, "Все поля (username, email, password) должны быть заполнены", http.StatusBadRequest)
		return
	}

	// Вызов сервиса для регистрации
	resp, err := h.authClient.Register(r.Context(), req.Username, req.Email, req.Password)
	if err != nil {
		logger.Error("Ошибка регистрации", zap.Error(err))
		http.Error(w, "Ошибка регистрации", http.StatusInternalServerError)
		return
	}

	// Отправка ответа
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// Login обрабатывает запрос на вход в систему
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// Декодирование запроса
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Ошибка разбора запроса", http.StatusBadRequest)
		return
	}

	// Валидация данных
	if req.Username == "" || req.Password == "" {
		http.Error(w, "Все поля (username, password) должны быть заполнены", http.StatusBadRequest)
		return
	}

	// Получаем информацию о клиенте
	userAgent := r.UserAgent()
	ip := r.RemoteAddr

	// Вызов сервиса для входа
	resp, err := h.authClient.Login(r.Context(), req.Username, req.Password, userAgent, ip)
	if err != nil {
		logger.Error("Ошибка входа", zap.Error(err))
		http.Error(w, "Ошибка входа", http.StatusUnauthorized)
		return
	}

	// Отправка ответа
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// RefreshToken обрабатывает запрос на обновление токена
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	// Декодирование запроса
	var req RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Ошибка разбора запроса", http.StatusBadRequest)
		return
	}

	// Валидация данных
	if req.RefreshToken == "" {
		http.Error(w, "Отсутствует refresh token", http.StatusBadRequest)
		return
	}

	// Получаем информацию о клиенте
	userAgent := r.UserAgent()
	ip := r.RemoteAddr

	// Вызов сервиса для обновления токена
	resp, err := h.authClient.RefreshToken(r.Context(), req.RefreshToken, userAgent, ip)
	if err != nil {
		logger.Error("Ошибка обновления токена", zap.Error(err))
		http.Error(w, "Ошибка обновления токена", http.StatusUnauthorized)
		return
	}

	// Отправка ответа
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// Logout обрабатывает запрос на выход из системы
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Декодирование запроса
	var req LogoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Ошибка разбора запроса", http.StatusBadRequest)
		return
	}

	// Валидация данных
	if req.RefreshToken == "" {
		http.Error(w, "Отсутствует refresh token", http.StatusBadRequest)
		return
	}

	// Вызов сервиса для выхода
	resp, err := h.authClient.Logout(r.Context(), req.RefreshToken)
	if err != nil {
		logger.Error("Ошибка выхода", zap.Error(err))
		http.Error(w, "Ошибка выхода", http.StatusInternalServerError)
		return
	}

	// Отправка ответа
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
