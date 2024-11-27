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
	GetReportIfExists(ctx context.Context, firstUserID int, secondUserID int) (models.Report, error)
	CheckUsersBlockNotExists(ctx context.Context, firstUserID int, secondUserID int) (string, error)
}

type MessageUsecase interface {
	AddMessage(ctx context.Context, message *models.Message) (int, error)
	GetLastMessage(ctx context.Context, authorID int, receiverID int) (models.Message, bool, error)
	GetChatMessages(ctx context.Context, firstUserID int, secondUserID int) ([]models.Message, error)
	GetMessagesBySearch(ctx context.Context, userID int, page int, search string) ([]models.Message, error)
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
		Reason:   in.Report.Reason,
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

func (h *GRPCMessageHandler) GetMessagesBySearch(ctx context.Context,
	in *generatedMessage.GetMessagesBySearchRequest) (*generatedMessage.GetMessagesBySearchResponse, error) {
	userID := int(in.UserID)
	page := int(in.Page)
	search := in.Search
	msgs, err := h.messageUsecase.GetMessagesBySearch(ctx, userID, page, search)
	if err != nil {
		h.logger.Error("grpc GetMessagesBySearch error", zap.Error(err))
		return nil, fmt.Errorf("grpc GetMessagesBySearch error: %w", err)
	}
	h.logger.Info("msgs", zap.Any("msgs", msgs))
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
	response := &generatedMessage.GetMessagesBySearchResponse{
		Messages: Messages,
	}
	return response, nil
}

func (h *GRPCMessageHandler) GetReportIfExists(ctx context.Context,
	in *generatedMessage.GetReportIfExistsRequest) (*generatedMessage.GetReportIfExistsResponse, error) {
	report := models.Report{
		ID:       int(in.Report.ID),
		Author:   int(in.Report.Author),
		Receiver: int(in.Report.Receiver),
		Body:     in.Report.Body,
	}
	response, err := h.reportUsecase.GetReportIfExists(ctx, report.Author, report.Receiver)
	if err != nil {
		if err.Error() == "this report dont exists" {
			return &generatedMessage.GetReportIfExistsResponse{}, err
		}
		h.logger.Error("grpc GetReportIfExists error", zap.Error(err))
		return &generatedMessage.GetReportIfExistsResponse{}, fmt.Errorf("grpc GetReportIfExists error: %w", err)
	}
	resp := &generatedMessage.GetReportIfExistsResponse{
		Report: &generatedMessage.Report{
			Author:   int32(response.Author),
			Receiver: int32(response.Receiver),
			Body:     response.Body,
		},
	}
	return resp, nil
}

func (h *GRPCMessageHandler) CheckUsersBlockNotExists(ctx context.Context,
	in *generatedMessage.CheckUsersBlockNotExistsRequest) (*generatedMessage.CheckUsersBlockNotExistsResponse, error) {
	firstUserID := int(in.FirstUserID)
	secondUserID := int(in.SecondUserID)

	status, err := h.reportUsecase.CheckUsersBlockNotExists(ctx, firstUserID, secondUserID)
	if err != nil {
		if err.Error() == "block exists" {
			resp := &generatedMessage.CheckUsersBlockNotExistsResponse{Status: status}
			h.logger.Info("resp", zap.Any("resp", resp))
			return resp, nil
		} else {
			h.logger.Error("grpc CheckUsersBlockNotExists error", zap.Error(err))
			return &generatedMessage.CheckUsersBlockNotExistsResponse{Status: ""}, fmt.Errorf("grpc CheckUsersBlockNotExists error: %w", err)
		}
	}
	response := &generatedMessage.CheckUsersBlockNotExistsResponse{Status: ""}
	h.logger.Info("grpc check block success")
	h.logger.Info("response", zap.Any("response", status))
	return response, nil
}
