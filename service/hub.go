package service

import (
	"context"
)

// Hub управляет всеми активными клиентами и сообщениями
type Hub struct {
	// Канал входящих сообщений от клиентов
	Broadcast chan Message
	service   *ChatService
}

// NewHub создает новый экземпляр хаба
func NewHub(service *ChatService) *Hub {
	return &Hub{
		Broadcast: make(chan Message, 256),
		service:   service,
	}
}

// Run запускает цикл обработки сообщений
func (h *Hub) Run() {
	for {
		select {
		case message := <-h.Broadcast:
			// Здесь будет логика обработки сообщений
			// Например, сохранение в базу данных и отправка клиентам
			_, _ = h.service.SaveMessage(context.Background(), message)
		}
	}
}
