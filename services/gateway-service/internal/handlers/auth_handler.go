package handlers

import (
	"encoding/json"
	"net/http"

	"chat-app/services/gateway-service/internal/clients"
)

type AuthHandler struct {
	authClient *clients.AuthClient
}

func NewAuthHandler(authClient *clients.AuthClient) *AuthHandler {
	return &AuthHandler{
		authClient: authClient,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный запрос", http.StatusBadRequest)
		return
	}

	response, err := h.authClient.Register(r.Context(), req.Username, req.Email, req.Password)
	if err != nil {
		http.Error(w, "Ошибка регистрации: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный запрос", http.StatusBadRequest)
		return
	}

	userAgent := r.UserAgent()
	ip := r.RemoteAddr

	response, err := h.authClient.Login(r.Context(), req.Username, req.Password, userAgent, ip)
	if err != nil {
		http.Error(w, "Ошибка входа: "+err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Дополнительные методы обработчика (RefreshToken, Logout, etc.)
