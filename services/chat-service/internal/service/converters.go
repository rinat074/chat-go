package service

import (
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Временная заглушка для proto типов
type ChatMessage struct {
	Id         int64
	Type       int
	Content    string
	UserId     int64
	Username   string
	ReceiverId *int64
	GroupId    *int64
	CreatedAt  *timestamppb.Timestamp
}

// Константы для типов сообщений
const (
	MessageTypePublic  = 0
	MessageTypePrivate = 1
	MessageTypeGroup   = 2
)

// Конвертация из proto в модель домена
func MessageFromProto(protoMsg *ChatMessage) Message {
	msg := Message{
		Type:      messageTypeFromProto(protoMsg.Type),
		Content:   protoMsg.Content,
		UserID:    protoMsg.UserId,
		Username:  protoMsg.Username,
		CreatedAt: protoMsg.CreatedAt.AsTime(),
	}

	if protoMsg.ReceiverId != nil {
		msg.ReceiverID = protoMsg.ReceiverId
	}

	if protoMsg.GroupId != nil {
		msg.GroupID = protoMsg.GroupId
	}

	return msg
}

// Конвертация из модели домена в proto
func MessageToProto(msg *Message) *ChatMessage {
	protoMsg := &ChatMessage{
		Id:        msg.ID,
		Type:      messageTypeToProto(msg.Type),
		Content:   msg.Content,
		UserId:    msg.UserID,
		Username:  msg.Username,
		CreatedAt: timestamppb.New(msg.CreatedAt),
	}

	if msg.ReceiverID != nil {
		protoMsg.ReceiverId = msg.ReceiverID
	}

	if msg.GroupID != nil {
		protoMsg.GroupId = msg.GroupID
	}

	return protoMsg
}

// Конвертация типа сообщения из proto
func messageTypeFromProto(protoType int) MessageType {
	switch protoType {
	case MessageTypePrivate:
		return PrivateMessage
	case MessageTypeGroup:
		return GroupMessage
	default:
		return PublicMessage
	}
}

// Конвертация типа сообщения в proto
func messageTypeToProto(msgType MessageType) int {
	switch msgType {
	case PrivateMessage:
		return MessageTypePrivate
	case GroupMessage:
		return MessageTypeGroup
	default:
		return MessageTypePublic
	}
}
