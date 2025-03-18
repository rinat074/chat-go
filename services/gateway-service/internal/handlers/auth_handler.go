package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/rinat074/chat-go/services/gateway-service/internal/clients"
	"github.com/rinat074/chat-go/services/gateway-service/pkg/logger"
)

type AuthHandler struct {
	clients *clients.ServiceClients
	log     logger.Logger
}

func NewAuthHandler(clients *clients.ServiceClients, log logger.Logger) *AuthHandler {
	return &AuthHandler{
		clients: clients,
		log:     log,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("ошибка декодирования запроса регистрации", "error", err)
		http.Error(w, "Неверный запрос", http.StatusBadRequest)
		return
	}

	response, err := h.clients.AuthClient.Register(r.Context(), req.Username, req.Email, req.Password)
	if err != nil {
		h.log.Error("ошибка регистрации пользователя", "error", err)
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
		h.log.Error("ошибка декодирования запроса входа", "error", err)
		http.Error(w, "Неверный запрос", http.StatusBadRequest)
		return
	}

	userAgent := r.UserAgent()
	ip := r.RemoteAddr

	response, err := h.clients.AuthClient.Login(r.Context(), req.Username, req.Password, userAgent, ip)
	if err != nil {
		h.log.Error("ошибка входа пользователя", "error", err, "username", req.Username)
		http.Error(w, "Ошибка входа: "+err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("ошибка декодирования запроса обновления токена", "error", err)
		http.Error(w, "Неверный запрос", http.StatusBadRequest)
		return
	}

	userAgent := r.UserAgent()
	ip := r.RemoteAddr

	tokenPair, err := h.clients.AuthClient.RefreshToken(r.Context(), req.RefreshToken, userAgent, ip)
	if err != nil {
		h.log.Error("ошибка обновления токена", "error", err)
		http.Error(w, "Ошибка обновления токена: "+err.Error(), http.StatusUnauthorized)
		return
	}

	// Формируем ответ
	response := struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("ошибка декодирования запроса выхода", "error", err)
		http.Error(w, "Неверный запрос", http.StatusBadRequest)
		return
	}

	err := h.clients.AuthClient.Logout(r.Context(), req.RefreshToken)
	if err != nil {
		h.log.Error("ошибка выхода", "error", err)
		http.Error(w, "Ошибка выхода: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
