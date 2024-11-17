package report

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"
	"go.uber.org/zap"
)

type Repository interface {
	AddReport(ctx context.Context, report models.Report) (int, error)
}

type Usecase struct {
	repo   Repository
	logger *zap.Logger
}

func New(repo Repository, logger *zap.Logger) *Usecase {
	return &Usecase{
		repo:   repo,
		logger: logger,
	}
}

func (s *Usecase) AddReport(ctx context.Context, report models.Report) (int, error) {
	req_id := ctx.Value(consts.RequestIDKey).(string)
	s.logger.Info("Usecase AddReport request_id", zap.String("req_id", req_id))

	reportId, err := s.repo.AddReport(ctx, report)
	if err != nil {
		s.logger.Error("Usecase AddReport error", zap.Error(err))
		return -1, fmt.Errorf("AddReport error: %w", err)
	}
	return reportId, nil
}
