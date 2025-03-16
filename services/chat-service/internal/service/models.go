package service

import (
	"time"
)

// MessageType определяет тип сообщения
type MessageType string

const (
	PublicMessage  MessageType = "public"
	PrivateMessage MessageType = "private"
	GroupMessage   MessageType = "group"
)

// Message представляет сообщение чата
type Message struct {
	ID         int64       `json:"id"`
	Type       MessageType `json:"type"`
	Content    string      `json:"content"`
	UserID     int64       `json:"user_id"`
	Username   string      `json:"username"`
	ReceiverID *int64      `json:"receiver_id,omitempty"`
	GroupID    *int64      `json:"group_id,omitempty"`
	CreatedAt  time.Time   `json:"created_at"`
}
