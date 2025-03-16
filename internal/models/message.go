package models

import (
	"time"
)

type MessageType string

const (
	PublicMessage  MessageType = "public"
	PrivateMessage MessageType = "private"
	GroupMessage   MessageType = "group"
)

type Message struct {
	ID         int64       `json:"id"`
	Type       MessageType `json:"type"`
	Content    string      `json:"content"`
	UserID     int64       `json:"user_id"`
	Username   string      `json:"username"`
	ReceiverID *int64      `json:"receiver_id,omitempty"` // Для личных сообщений
	GroupID    *int64      `json:"group_id,omitempty"`    // Для групповых сообщений
	CreatedAt  time.Time   `json:"created_at"`
}
