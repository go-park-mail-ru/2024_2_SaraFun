package message

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"
	"go.uber.org/zap"
)

//go:generate mockgen -destination=./mocks/mock_repository.go -package=mocks . Repository

type Repository interface {
	AddMessage(ctx context.Context, message *models.Message) (int, error)
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

func (s *UseCase) AddMessage(ctx context.Context, message *models.Message) (int, error) {
	req_id := ctx.Value(consts.RequestIDKey).(string)
	s.logger.Info("AddMessage usecase", zap.String("req_id", req_id))
	messageId, err := s.repo.AddMessage(ctx, message)
	if err != nil {
		s.logger.Error("AddMessage error", zap.Error(err))
		return -1, fmt.Errorf("Usecase AddMessage error: %w", err)
	}
	return messageId, nil
}
