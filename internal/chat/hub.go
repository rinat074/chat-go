package chat

import (
	"context"

	"chat-app/internal/models"
)

// Hub управляет всеми активными клиентами и сообщениями
type Hub struct {
	// Зарегистрированные клиенты
	clients map[*Client]bool

	// Канал входящих сообщений от клиентов
	broadcast chan models.Message

	// Канал регистрации клиента
	register chan *Client

	// Канал отмены регистрации клиента
	unregister chan *Client

	// Сервис для работы с сообщениями
	service *Service
}

func NewHub(service *Service) *Hub {
	return &Hub{
		broadcast:  make(chan models.Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		service:    service,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			// Сохраняем сообщение в базе данных
			savedMsg, err := h.service.SaveMessage(context.Background(), message)
			if err != nil {
				continue
			}

			// Отправляем сообщение всем клиентам
			for client := range h.clients {
				select {
				case client.send <- *savedMsg:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
