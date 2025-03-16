package service

import (
	"chat-app/proto/chat"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// MessageFromProto конвертирует proto-сообщение в модель домена
func MessageFromProto(protoMsg *chat.Message) Message {
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

// MessageToProto конвертирует модель домена в proto-сообщение
func MessageToProto(msg *Message) *chat.Message {
	protoMsg := &chat.Message{
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

// messageTypeFromProto конвертирует тип сообщения из proto
func messageTypeFromProto(protoType chat.MessageType) MessageType {
	switch protoType {
	case chat.MessageType_PRIVATE:
		return PrivateMessage
	case chat.MessageType_GROUP:
		return GroupMessage
	default:
		return PublicMessage
	}
}

// messageTypeToProto конвертирует тип сообщения в proto
func messageTypeToProto(msgType MessageType) chat.MessageType {
	switch msgType {
	case PrivateMessage:
		return chat.MessageType_PRIVATE
	case GroupMessage:
		return chat.MessageType_GROUP
	default:
		return chat.MessageType_PUBLIC
	}
}
