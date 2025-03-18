package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"github.com/rinat074/chat-go/services/gateway-service/internal/clients"
	"github.com/rinat074/chat-go/services/gateway-service/pkg/logger"
)

const (
	// Максимальное время ожидания для записи сообщения клиенту
	writeWait = 10 * time.Second

	// Максимальное время между pong от клиента
	pongWait = 60 * time.Second

	// Интервал отправки ping сообщений
	pingPeriod = (pongWait * 9) / 10

	// Максимальный размер входящего сообщения
	maxMessageSize = 10 * 1024 // 10KB
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Разрешаем соединения с любых источников (в продакшене стоит ограничить)
	},
}

// WebSocketClient представляет WebSocket клиента
type WebSocketClient struct {
	conn   *websocket.Conn
	userID int64
	send   chan []byte
	hub    *WebSocketHub
	log    logger.Logger
}

// WebSocketHub управляет набором активных WebSocket клиентов
type WebSocketHub struct {
	clients    map[*WebSocketClient]bool
	broadcast  chan []byte
	register   chan *WebSocketClient
	unregister chan *WebSocketClient
	redis      *redis.Client
	log        logger.Logger
	mu         sync.RWMutex
}

// WebSocketHandler обрабатывает WebSocket соединения
type WebSocketHandler struct {
	clients *clients.ServiceClients
	hub     *WebSocketHub
	log     logger.Logger
}

// NewWebSocketHandler создает новый обработчик WebSocket
func NewWebSocketHandler(clients *clients.ServiceClients, log logger.Logger) *WebSocketHandler {
	hub := &WebSocketHub{
		clients:    make(map[*WebSocketClient]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *WebSocketClient),
		unregister: make(chan *WebSocketClient),
		redis:      clients.RedisClient,
		log:        log,
	}

	go hub.run()

	// Подписка на канал сообщений Redis
	go hub.subscribeToRedis()

	return &WebSocketHandler{
		clients: clients,
		hub:     hub,
		log:     log,
	}
}

// HandleWebSocket обрабатывает WebSocket соединение
func (h *WebSocketHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.log.Error("ошибка обновления соединения до WebSocket", "error", err)
		return
	}

	// Получение userID из контекста (установлен middleware аутентификации)
	userID := r.Context().Value("userID").(int64)

	client := &WebSocketClient{
		conn:   conn,
		userID: userID,
		send:   make(chan []byte, 256),
		hub:    h.hub,
		log:    h.log,
	}

	// Регистрация клиента в хабе
	h.hub.register <- client

	// Запуск горутин для чтения и записи сообщений
	go client.readPump()
	go client.writePump()
}

// run запускает цикл обработки событий WebSocket хаба
func (h *WebSocketHub) run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			h.log.Info("WebSocket клиент подключен", "userID", client.userID)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mu.Unlock()
			h.log.Info("WebSocket клиент отключен", "userID", client.userID)

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					h.mu.RUnlock()
					h.mu.Lock()
					delete(h.clients, client)
					close(client.send)
					h.mu.Unlock()
					h.mu.RLock()
				}
			}
			h.mu.RUnlock()
		}
	}
}

// subscribeToRedis подписывается на сообщения Redis
func (h *WebSocketHub) subscribeToRedis() {
	pubsub := h.redis.Subscribe(context.Background(), "chat_messages")
	defer pubsub.Close()

	ch := pubsub.Channel()
	for msg := range ch {
		h.broadcast <- []byte(msg.Payload)
	}
}

// readPump читает сообщения от WebSocket клиента
func (c *WebSocketClient) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.log.Error("ошибка чтения WebSocket", "error", err)
			}
			break
		}

		// Обработка входящего сообщения
		var msg struct {
			Type    string          `json:"type"`
			Content json.RawMessage `json:"content"`
		}

		if err := json.Unmarshal(message, &msg); err != nil {
			c.log.Error("ошибка разбора WebSocket сообщения", "error", err)
			continue
		}

		// Здесь можно добавить логику обработки различных типов сообщений
	}
}

// writePump отправляет сообщения WebSocket клиенту
func (c *WebSocketClient) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// Канал закрыт
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Добавляем в очередь все ждущие сообщения
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
