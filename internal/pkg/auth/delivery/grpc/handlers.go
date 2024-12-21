package authgrpc

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	generatedAuth "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc/gen"
	"go.uber.org/zap"
	"time"
)

//go:generate mockgen -destination=./mocks/mock_usecase.go -package=mocks github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/auth/delivery/grpc UseCase

type UseCase interface {
	CreateSession(ctx context.Context, user models.User) (models.Session, error)
	DeleteSession(ctx context.Context, sessionID string) error
	CheckSession(ctx context.Context, sessionID string) error
	GetUserIDBySessionID(ctx context.Context, sessionID string) (int, error)
}

type GrpcAuthHandler struct {
	generatedAuth.AuthServer
	uc     UseCase
	logger *zap.Logger
}

func NewGRPCAuthHandler(usecase UseCase, logger *zap.Logger) *GrpcAuthHandler {
	return &GrpcAuthHandler{uc: usecase, logger: logger}
}

func (h *GrpcAuthHandler) CreateSession(ctx context.Context, in *generatedAuth.CreateSessionRequest) (*generatedAuth.CreateSessionResponse, error) {
	h.logger.Info("CreateSession grpc started")
	payload := models.User{ID: int(in.User.ID)}
	h.logger.Info("test")
	session, err := h.uc.CreateSession(ctx, payload)
	if err != nil {
		return nil, err
	}
	return &generatedAuth.CreateSessionResponse{
		Session: &generatedAuth.Session{
			SessionID: session.SessionID,
			UserID:    int32(session.UserID),
			CreatedAt: session.CreatedAt.Format(time.RFC3339),
			ExpiresAt: session.ExpiresAt.Format(time.RFC3339),
		},
	}, nil
}

func (h *GrpcAuthHandler) DeleteSession(ctx context.Context, in *generatedAuth.DeleteSessionRequest) (*generatedAuth.DeleteSessionResponse, error) {
	h.logger.Info("DeleteSession grpc started")
	err := h.uc.DeleteSession(ctx, in.SessionID)
	if err != nil {
		return nil, fmt.Errorf("Grpc delete session error : %w", err)
	}
	return &generatedAuth.DeleteSessionResponse{}, nil
}

func (h *GrpcAuthHandler) CheckSession(ctx context.Context, in *generatedAuth.CheckSessionRequest) (*generatedAuth.CheckSessionResponse, error) {
	h.logger.Info("CheckSession grpc started")
	sessionId := in.SessionID
	err := h.uc.CheckSession(ctx, sessionId)
	if err != nil {
		return nil, fmt.Errorf("Grpc check session error : %w", err)
	}
	h.logger.Info("CheckSession grpc finished")
	return &generatedAuth.CheckSessionResponse{}, nil
}

func (h *GrpcAuthHandler) GetUserIDBySessionID(ctx context.Context, in *generatedAuth.GetUserIDBySessionIDRequest) (*generatedAuth.GetUserIDBYSessionIDResponse, error) {
	h.logger.Info("GetUserIdBySessionId grpc started")
	sessionId := in.SessionID
	userId, err := h.uc.GetUserIDBySessionID(ctx, sessionId)
	if err != nil {
		return nil, fmt.Errorf("Grpc get user id by session id error : %w", err)
	}
	return &generatedAuth.GetUserIDBYSessionIDResponse{UserId: int32(userId)}, nil
}
