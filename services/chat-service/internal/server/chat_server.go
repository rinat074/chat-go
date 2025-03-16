package server

import (
	"context"

	"chat-app/services/chat-service/internal/service"

	"github.com/gochat/proto/chat"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ChatServer struct {
	chat.UnimplementedChatServiceServer
	service *service.ChatService
	hub     *service.Hub
}

func RegisterChatServer(s *grpc.Server, svc *service.ChatService, hub *service.Hub) {
	chat.RegisterChatServiceServer(s, &ChatServer{
		service: svc,
		hub:     hub,
	})
}

// Преобразование из protobuf в модель домена
func messageFromProto(protoMsg *chat.Message) service.Message {
	msgType := service.PublicMessage
	switch protoMsg.Type {
	case chat.MessageType_PRIVATE:
		msgType = service.PrivateMessage
	case chat.MessageType_GROUP:
		msgType = service.GroupMessage
	}

	msg := service.Message{
		Content:   protoMsg.Content,
		UserID:    protoMsg.UserId,
		Username:  protoMsg.Username,
		CreatedAt: protoMsg.CreatedAt.AsTime(),
		Type:      msgType,
	}

	if protoMsg.ReceiverId != nil && *protoMsg.ReceiverId != 0 {
		receiverID := *protoMsg.ReceiverId
		msg.ReceiverID = &receiverID
	}

	if protoMsg.GroupId != nil && *protoMsg.GroupId != 0 {
		groupID := *protoMsg.GroupId
		msg.GroupID = &groupID
	}

	return msg
}

// Преобразование из модели домена в protobuf
func messageToProto(msg *service.Message) *chat.Message {
	msgType := chat.MessageType_PUBLIC
	switch msg.Type {
	case service.PrivateMessage:
		msgType = chat.MessageType_PRIVATE
	case service.GroupMessage:
		msgType = chat.MessageType_GROUP
	}

	protoMsg := &chat.Message{
		Id:        msg.ID,
		Type:      msgType,
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

func (s *ChatServer) SaveMessage(ctx context.Context, req *chat.Message) (*chat.Message, error) {
	// Преобразуем protobuf сообщение в модель домена
	message := messageFromProto(req)

	// Сохраняем сообщение через сервис
	savedMsg, err := s.service.SaveMessage(ctx, message)
	if err != nil {
		return nil, err
	}

	// Отправляем сообщение в хаб для рассылки клиентам
	s.hub.Broadcast <- *savedMsg

	// Преобразуем результат обратно в protobuf и возвращаем
	return messageToProto(savedMsg), nil
}

func (s *ChatServer) GetPublicMessages(ctx context.Context, req *chat.GetMessagesRequest) (*chat.MessagesResponse, error) {
	messages, err := s.service.GetPublicMessages(ctx, int(req.Limit), int(req.Offset))
	if err != nil {
		return nil, err
	}

	// Преобразуем сообщения в формат protobuf
	protoMessages := make([]*chat.Message, len(messages))
	for i, msg := range messages {
		msgCopy := msg // Создаем копию, чтобы избежать проблем с указателями
		protoMessages[i] = messageToProto(&msgCopy)
	}

	return &chat.MessagesResponse{
		Messages: protoMessages,
	}, nil
}

// Аналогично реализуем остальные методы: GetPrivateMessages, GetGroupMessages, CreateGroup, AddUserToGroup
