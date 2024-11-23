package grpc

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	generatedMessage "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/message/delivery/grpc/gen"
	"go.uber.org/zap"
)

type ReportUsecase interface {
	AddReport(ctx context.Context, report models.Report) (int, error)
}

type MessageUsecase interface {
	AddMessage(ctx context.Context, message *models.Message) (int, error)
	GetLastMessage(ctx context.Context, authorID int, receiverID int) (models.Message, bool, error)
	GetChatMessages(ctx context.Context, firstUserID int, secondUserID int) ([]models.Message, error)
}

type GRPCMessageHandler struct {
	generatedMessage.MessageServer
	reportUsecase  ReportUsecase
	messageUsecase MessageUsecase
	logger         *zap.Logger
}

func NewGRPCHandler(reportUsecase ReportUsecase, messageUsecase MessageUsecase, logger *zap.Logger) *GRPCMessageHandler {
	return &GRPCMessageHandler{
		reportUsecase:  reportUsecase,
		messageUsecase: messageUsecase,
		logger:         logger,
	}
}

func (h *GRPCMessageHandler) AddReport(ctx context.Context, in *generatedMessage.AddReportRequest) (*generatedMessage.AddReportResponse, error) {
	report := models.Report{
		ID:       int(in.Report.ID),
		Author:   int(in.Report.Author),
		Receiver: int(in.Report.Receiver),
		Body:     in.Report.Body,
	}
	id, err := h.reportUsecase.AddReport(ctx, report)
	if err != nil {
		return nil, fmt.Errorf("grpc AddReport error: %w", err)
	}
	response := &generatedMessage.AddReportResponse{
		ReportID: int32(id),
	}
	return response, nil
}

func (h *GRPCMessageHandler) AddMessage(ctx context.Context, in *generatedMessage.AddMessageRequest) (*generatedMessage.AddMessageResponse, error) {
	msg := models.Message{
		ID:       int(in.Message.ID),
		Author:   int(in.Message.Author),
		Receiver: int(in.Message.Receiver),
		Body:     in.Message.Body,
	}
	id, err := h.messageUsecase.AddMessage(ctx, &msg)
	if err != nil {
		return nil, fmt.Errorf("grpc AddMessage error: %w", err)
	}
	response := &generatedMessage.AddMessageResponse{
		MessageID: int32(id),
	}
	return response, nil
}

func (h *GRPCMessageHandler) GetLastMessage(ctx context.Context, in *generatedMessage.GetLastMessageRequest) (*generatedMessage.GetLastMessageResponse, error) {
	authorID := int(in.AuthorID)
	receiverID := int(in.ReceiverID)
	msg, self, err := h.messageUsecase.GetLastMessage(ctx, authorID, receiverID)
	if err != nil {
		h.logger.Error("grpc GetLastMessage error", zap.Error(err))
		return nil, fmt.Errorf("grpc GetLastMessage error: %w", err)
	}
	response := &generatedMessage.GetLastMessageResponse{Message: msg.Body, Self: self, Time: msg.Time}
	return response, nil
}

func (h *GRPCMessageHandler) GetChatMessages(ctx context.Context,
	in *generatedMessage.GetChatMessagesRequest) (*generatedMessage.GetChatMessagesResponse, error) {
	firstUserID := int(in.FirstUserID)
	secondUserID := int(in.SecondUserID)
	msgs, err := h.messageUsecase.GetChatMessages(ctx, firstUserID, secondUserID)
	if err != nil {
		h.logger.Error("grpc GetChatMessages error", zap.Error(err))
		return nil, fmt.Errorf("grpc GetChatMessages error: %w", err)
	}
	var Messages []*generatedMessage.ChatMessage
	for _, msg := range msgs {
		Message := &generatedMessage.ChatMessage{
			Author:   int32(msg.Author),
			Receiver: int32(msg.Receiver),
			Body:     msg.Body,
			Time:     msg.Time,
		}
		Messages = append(Messages, Message)
	}
	response := &generatedMessage.GetChatMessagesResponse{
		Messages: Messages,
	}
	return response, nil
}
