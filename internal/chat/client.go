package chat

import (
	"bytes"
	"encoding/json"
	"log"
	"time"

	"chat-app/internal/models"

	"github.com/gorilla/websocket"
)

const (
	// Время ожидания для записи сообщения клиенту
	writeWait = 10 * time.Second

	// Время ожидания для чтения следующего сообщения от клиента
	pongWait = 60 * time.Second

	// Интервал отправки ping-сообщений
	pingPeriod = (pongWait * 9) / 10

	// Максимальный размер сообщения
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

// Client представляет собой соединение клиента по веб-сокету
type Client struct {
	hub      *Hub
	conn     *websocket.Conn
	send     chan models.Message
	userID   int64
	username string
}

// readPump запускает цикл чтения сообщений от клиента
func (c *Client) readPump() {
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
				log.Printf("Ошибка: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

		// Разбираем сообщение
		var msg struct {
			Content string `json:"content"`
		}
		if err := json.Unmarshal(message, &msg); err != nil {
			continue
		}

		// Формируем модель сообщения
		chatMessage := models.Message{
			Content:   msg.Content,
			UserID:    c.userID,
			Username:  c.username,
			CreatedAt: time.Now(),
		}

		// Отправляем сообщение в хаб
		c.hub.broadcast <- chatMessage
	}
}

// writePump запускает цикл записи сообщений клиенту
func (c *Client) writePump() {
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
				// Хаб закрыл канал
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			// Отправляем сообщение в формате JSON
			messageJSON, _ := json.Marshal(message)
			w.Write(messageJSON)

			// Добавляем в очередь ожидающие сообщения
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				nextMsg := <-c.send
				nextMsgJSON, _ := json.Marshal(nextMsg)
				w.Write(nextMsgJSON)
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
