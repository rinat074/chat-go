package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/rinat074/chat-go/services/gateway-service/internal/clients"
	"github.com/rinat074/chat-go/services/gateway-service/pkg/logger"
)

type ChatHandler struct {
	clients *clients.ServiceClients
	log     logger.Logger
}

func NewChatHandler(clients *clients.ServiceClients, log logger.Logger) *ChatHandler {
	return &ChatHandler{
		clients: clients,
		log:     log,
	}
}

func (h *ChatHandler) GetPublicMessages(w http.ResponseWriter, r *http.Request) {
	limit := 50
	offset := 0

	if limitParam := r.URL.Query().Get("limit"); limitParam != "" {
		if l, err := strconv.Atoi(limitParam); err == nil && l > 0 {
			limit = l
		}
	}

	if offsetParam := r.URL.Query().Get("offset"); offsetParam != "" {
		if o, err := strconv.Atoi(offsetParam); err == nil && o >= 0 {
			offset = o
		}
	}

	messages, err := h.clients.ChatClient.GetPublicMessages(r.Context(), limit, offset)
	if err != nil {
		h.log.Error("ошибка получения публичных сообщений", "error", err)
		http.Error(w, "Ошибка получения сообщений: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

func (h *ChatHandler) GetPrivateMessages(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	if userID == "" {
		http.Error(w, "Отсутствует ID пользователя", http.StatusBadRequest)
		return
	}

	limit := 50
	offset := 0

	if limitParam := r.URL.Query().Get("limit"); limitParam != "" {
		if l, err := strconv.Atoi(limitParam); err == nil && l > 0 {
			limit = l
		}
	}

	if offsetParam := r.URL.Query().Get("offset"); offsetParam != "" {
		if o, err := strconv.Atoi(offsetParam); err == nil && o >= 0 {
			offset = o
		}
	}

	userIDInt, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		h.log.Error("неверный формат userID", "error", err, "userID", userID)
		http.Error(w, "Неверный ID пользователя", http.StatusBadRequest)
		return
	}

	// Получаем userID текущего пользователя из контекста
	currentUserID := r.Context().Value("userID").(int64)

	messages, err := h.clients.ChatClient.GetPrivateMessages(r.Context(), currentUserID, userIDInt, limit, offset)
	if err != nil {
		h.log.Error("ошибка получения приватных сообщений", "error", err, "userID", userID)
		http.Error(w, "Ошибка получения сообщений: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

func (h *ChatHandler) GetGroupMessages(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "groupID")
	if groupID == "" {
		http.Error(w, "Отсутствует ID группы", http.StatusBadRequest)
		return
	}

	limit := 50
	offset := 0

	if limitParam := r.URL.Query().Get("limit"); limitParam != "" {
		if l, err := strconv.Atoi(limitParam); err == nil && l > 0 {
			limit = l
		}
	}

	if offsetParam := r.URL.Query().Get("offset"); offsetParam != "" {
		if o, err := strconv.Atoi(offsetParam); err == nil && o >= 0 {
			offset = o
		}
	}

	groupIDInt, err := strconv.ParseInt(groupID, 10, 64)
	if err != nil {
		h.log.Error("неверный формат groupID", "error", err, "groupID", groupID)
		http.Error(w, "Неверный ID группы", http.StatusBadRequest)
		return
	}

	messages, err := h.clients.ChatClient.GetGroupMessages(r.Context(), groupIDInt, limit, offset)
	if err != nil {
		h.log.Error("ошибка получения групповых сообщений", "error", err, "groupID", groupID)
		http.Error(w, "Ошибка получения сообщений: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

func (h *ChatHandler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		MemberIDs   []int64 `json:"member_ids"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("ошибка декодирования запроса создания группы", "error", err)
		http.Error(w, "Неверный запрос", http.StatusBadRequest)
		return
	}

	// Получаем userID текущего пользователя из контекста
	userID := r.Context().Value("userID").(int64)

	group, err := h.clients.ChatClient.CreateGroup(r.Context(), req.Name, req.Description, userID, req.MemberIDs)
	if err != nil {
		h.log.Error("ошибка создания группы", "error", err, "name", req.Name)
		http.Error(w, "Ошибка создания группы: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(group)
}

func (h *ChatHandler) AddUserToGroup(w http.ResponseWriter, r *http.Request) {
	groupID := chi.URLParam(r, "groupID")
	if groupID == "" {
		http.Error(w, "Отсутствует ID группы", http.StatusBadRequest)
		return
	}

	var req struct {
		UserID int64 `json:"user_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("ошибка декодирования запроса добавления пользователя в группу", "error", err)
		http.Error(w, "Неверный запрос", http.StatusBadRequest)
		return
	}

	groupIDInt, err := strconv.ParseInt(groupID, 10, 64)
	if err != nil {
		h.log.Error("неверный формат groupID", "error", err, "groupID", groupID)
		http.Error(w, "Неверный ID группы", http.StatusBadRequest)
		return
	}

	// Получаем userID текущего пользователя из контекста для проверки прав
	currentUserID := r.Context().Value("userID").(int64)

	err = h.clients.ChatClient.AddUserToGroup(r.Context(), groupIDInt, req.UserID, currentUserID)
	if err != nil {
		h.log.Error("ошибка добавления пользователя в группу", "error", err, "groupID", groupID, "userID", req.UserID)
		http.Error(w, "Ошибка добавления пользователя: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
