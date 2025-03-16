package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"chat-app/services/gateway-service/internal/clients"
)

type ChatHandler struct {
	chatClient *clients.ChatClient
}

func NewChatHandler(chatClient *clients.ChatClient) *ChatHandler {
	return &ChatHandler{
		chatClient: chatClient,
	}
}

func (h *ChatHandler) GetPublicMessages(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := int32(50)
	offset := int32(0)

	if limitStr != "" {
		if l, err := strconv.ParseInt(limitStr, 10, 32); err == nil && l > 0 {
			limit = int32(l)
		}
	}

	if offsetStr != "" {
		if o, err := strconv.ParseInt(offsetStr, 10, 32); err == nil && o >= 0 {
			offset = int32(o)
		}
	}

	resp, err := h.chatClient.GetPublicMessages(r.Context(), limit, offset)
	if err != nil {
		http.Error(w, "Ошибка при получении сообщений: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// Остальные методы обработчика чата: GetPrivateMessages, GetGroupMessages, etc.
