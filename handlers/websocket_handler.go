package handlers

import (
	"net/http"

	"github.com/gorilla/websocket"
)

// Ключ для контекста пользователя
type contextKey string

const userContextKey contextKey = "user"

// UserContextKey возвращает ключ для контекста пользователя
func UserContextKey() contextKey {
	return userContextKey
}

// Структура данных пользователя
type userData struct {
	UserID   int64
	Username string
}

// UserData возвращает структуру данных пользователя
func UserData(userID int64, username string) userData {
	return userData{
		UserID:   userID,
		Username: username,
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Для разработки разрешаем все источники
	},
}

// ... остальной код не меняется
