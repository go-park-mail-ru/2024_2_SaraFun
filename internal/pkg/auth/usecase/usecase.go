package usecase

import (
	"context"
	"fmt"
	sparkiterrors "github.com/go-park-mail-ru/2024_2_SaraFun/internal/errors"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"time"
)

//go:generate mockgen -destination=./mocks/mock_repository.go -package=mocks . Repository
type Repository interface {
	AddSession(ctx context.Context, session models.Session) error
	DeleteSession(ctx context.Context, sessionID string) error
	GetUserIDBySessionID(ctx context.Context, sessionID string) (int, error)
	CheckSession(ctx context.Context, sessionID string) error
}

type UseCase struct {
	repo   Repository
	logger *zap.Logger
}

func New(repo Repository, logger *zap.Logger) *UseCase {
	return &UseCase{repo: repo, logger: logger}
}

func (s *UseCase) CreateSession(ctx context.Context, user models.User) (models.Session, error) {
	//req_id := ctx.Value(consts.RequestIDKey).(string)
	//s.logger.Info("usecase request-id", zap.String("request_id", req_id))
	session := models.Session{
		SessionID: uuid.New().String(),
		UserID:    user.ID,
		CreatedAt: time.Now(),
	}
	err := s.repo.AddSession(ctx, session)
	if err != nil {
		s.logger.Error("bad add session", zap.Error(err))
		return models.Session{}, fmt.Errorf("failed to create session: %v", err)
	}
	return session, nil
}

func (s *UseCase) GetUserIDBySessionID(ctx context.Context, sessionID string) (int, error) {
	//req_id := ctx.Value(consts.RequestIDKey).(string)
	//s.logger.Info("usecase request-id", zap.String("request_id", req_id))
	userID, err := s.repo.GetUserIDBySessionID(ctx, sessionID)
	if err != nil {
		s.logger.Error("failed to get user id by session id", zap.Error(err))
		return 0, sparkiterrors.ErrInvalidSession
	}
	return userID, nil
}

func (s *UseCase) CheckSession(ctx context.Context, sessionID string) error {
	//req_id := ctx.Value(consts.RequestIDKey).(string)
	//s.logger.Info("usecase request-id", zap.String("request_id", req_id))
	err := s.repo.CheckSession(ctx, sessionID)
	if err != nil {
		s.logger.Error("failed to check session", zap.Error(err))
		return sparkiterrors.ErrInvalidSession
	}
	return nil
}

func (s *UseCase) DeleteSession(ctx context.Context, sessionID string) error {
	//req_id := ctx.Value(consts.RequestIDKey).(string)
	//s.logger.Info("usecase request-id", zap.String("request_id", req_id))
	err := s.repo.DeleteSession(ctx, sessionID)
	if err != nil {
		s.logger.Error("failed to delete session", zap.Error(err))
		return sparkiterrors.ErrInvalidSession
	}
	return nil
}
