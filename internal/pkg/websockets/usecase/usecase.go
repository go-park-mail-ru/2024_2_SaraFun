package usecase

import (
	"context"
	"fmt"
	ws "github.com/gorilla/websocket"
	"go.uber.org/zap"
)

//go:generate mockgen -destination=./mocks/mock_repository.go -package=mocks github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/websockets/usecase Repository

type Repository interface {
	AddConnection(ctx context.Context, conn *ws.Conn, userId int) error
	DeleteConnection(ctx context.Context, userId int) error
	WriteMessage(ctx context.Context, authorID int, receiverID int, message string, username string) error
	SendNotification(ctx context.Context, receiverID int, receiverImageLink string, authorUsername string) error
}

type UseCase struct {
	repo   Repository
	logger *zap.Logger
}

func New(repo Repository, logger *zap.Logger) *UseCase {
	return &UseCase{
		repo:   repo,
		logger: logger,
	}
}

func (u *UseCase) AddConnection(ctx context.Context, conn *ws.Conn, userId int) error {
	u.logger.Info("Usecase Add connection start", zap.Int("user_id", userId))
	err := u.repo.AddConnection(ctx, conn, userId)
	if err != nil {
		u.logger.Error("repo Add connection call in Usecase failed", zap.Error(err))
		return fmt.Errorf("repo Add connection call in Usecase failed: %w", err)
	}
	return nil
}

func (u *UseCase) DeleteConnection(ctx context.Context, userId int) error {
	u.logger.Info("Usecase Delete connection start", zap.Int("user_id", userId))
	err := u.repo.DeleteConnection(ctx, userId)
	if err != nil {
		u.logger.Error("repo Delete connection call in Usecase failed", zap.Error(err))
		return fmt.Errorf("repo Delete connection call in Usecase failed: %w", err)
	}
	return nil
}

func (u *UseCase) WriteMessage(ctx context.Context, authorID int, receiverID int, message string, username string) error {
	u.logger.Info("Usecase WriteMessage start", zap.Int("user_id", receiverID))
	err := u.repo.WriteMessage(ctx, authorID, receiverID, message, username)
	if err != nil {
		u.logger.Error("repo WriteMessage call in Usecase failed", zap.Error(err))
		return fmt.Errorf("repo WriteMessage call in Usecase failed: %w", err)
	}
	return nil
}

func (u *UseCase) SendNotification(ctx context.Context, receiverID int, authorUsername string, authorImageLink string) error {
	u.logger.Info("Usecase SendNotification start", zap.Int("user_id", receiverID))
	err := u.repo.SendNotification(ctx, receiverID, authorImageLink, authorUsername)
	if err != nil {
		u.logger.Error("repo SendNotification call in Usecase failed", zap.Error(err))
		return fmt.Errorf("repo SendNotification call in Usecase failed: %w", err)
	}
	return nil
}
