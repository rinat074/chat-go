package server

import (
	"context"

	"chat-app/proto/chat"
	"chat-app/services/chat-service/internal/service"

	"google.golang.org/grpc"
)

// ChatServer представляет gRPC сервер чата
type ChatServer struct {
	chat.UnimplementedChatServiceServer
	service *service.ChatService
	hub     *service.Hub
}

// RegisterChatServer регистрирует сервер чата в gRPC сервере
func RegisterChatServer(s *grpc.Server, svc *service.ChatService, hub *service.Hub) {
	chat.RegisterChatServiceServer(s, &ChatServer{
		service: svc,
		hub:     hub,
	})
}

// SaveMessage сохраняет сообщение
func (s *ChatServer) SaveMessage(ctx context.Context, req *chat.Message) (*chat.Message, error) {
	// Преобразуем protobuf сообщение в модель домена
	message := service.MessageFromProto(req)

	// Сохраняем сообщение через сервис
	savedMsg, err := s.service.SaveMessage(ctx, message)
	if err != nil {
		return nil, err
	}

	// Отправляем сообщение в хаб для рассылки клиентам
	s.hub.Broadcast <- *savedMsg

	// Преобразуем результат обратно в protobuf и возвращаем
	return service.MessageToProto(savedMsg), nil
}

// GetPublicMessages возвращает публичные сообщения
func (s *ChatServer) GetPublicMessages(ctx context.Context, req *chat.GetMessagesRequest) (*chat.MessagesResponse, error) {
	messages, err := s.service.GetPublicMessages(ctx, int(req.Limit), int(req.Offset))
	if err != nil {
		return nil, err
	}

	// Преобразуем сообщения в формат protobuf
	protoMessages := make([]*chat.Message, len(messages))
	for i, msg := range messages {
		protoMessages[i] = service.MessageToProto(&msg)
	}

	return &chat.MessagesResponse{
		Messages: protoMessages,
	}, nil
}
