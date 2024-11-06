package reaction

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"sparkit/internal/models"
	"sparkit/internal/utils/consts"
)

//go:generate mockgen -destination=./mocks/mock_repository.go -package=mocks . Repository
type Repository interface {
	AddReaction(ctx context.Context, reaction models.Reaction) error
	GetMatchList(ctx context.Context, userId int) ([]int, error)
}

type UseCase struct {
	repo   Repository
	logger *zap.Logger
}

func New(repo Repository, logger *zap.Logger) *UseCase {
	return &UseCase{repo: repo, logger: logger}
}

func (u *UseCase) AddReaction(ctx context.Context, reaction models.Reaction) error {
	req_id := ctx.Value(consts.RequestIDKey).(string)
	u.logger.Info("usecase request-id", zap.String("request_id", req_id))
	err := u.repo.AddReaction(ctx, reaction)
	if err != nil {
		u.logger.Error("UseCase AddReaction: failed to add reaction", zap.Error(err))
		return fmt.Errorf("failed to AddReaction: %w", err)
	}
	return nil
}

func (u *UseCase) GetMatchList(ctx context.Context, userId int) ([]int, error) {
	req_id := ctx.Value(consts.RequestIDKey).(string)
	u.logger.Info("usecase request-id", zap.String("request_id", req_id))
	authors, err := u.repo.GetMatchList(ctx, userId)
	if err != nil {
		u.logger.Error("UseCase GetMatchList: failed to GetMatchList", zap.Error(err))
		return nil, fmt.Errorf("failed to GetMatchList: %w", err)
	}
	return authors, nil
}
