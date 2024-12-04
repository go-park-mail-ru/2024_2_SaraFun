package report

import (
	"context"
	stderr "errors"
	"fmt"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	"go.uber.org/zap"
)

//go:generate mockgen -destination=./mocks/mock_repository.go -package=mocks . Repository
type Repository interface {
	AddReport(ctx context.Context, report models.Report) (int, error)
	GetReportIfExists(ctx context.Context, firstUserID int, secondUserID int) (models.Report, error)
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
	//req_id := ctx.Value(consts.RequestIDKey).(string)
	//s.logger.Info("Usecase AddReport request_id", zap.String("req_id", req_id))
	reportId, err := s.repo.AddReport(ctx, report)
	if err != nil {
		s.logger.Error("Usecase AddReport error", zap.Error(err))
		return -1, fmt.Errorf("AddReport error: %w", err)
	}
	return reportId, nil
}

func (s *Usecase) GetReportIfExists(ctx context.Context, firstUserID int, secondUserID int) (models.Report, error) {
	rep, err := s.repo.GetReportIfExists(ctx, firstUserID, secondUserID)
	if err != nil {
		if err.Error() == "this report dont exists" {
			return models.Report{}, err
		}
		s.logger.Error("Usecase GetReport error", zap.Error(err))
		return models.Report{}, fmt.Errorf("GetReport error: %w", err)
	}
	return rep, nil
}

func (s *Usecase) CheckUsersBlockNotExists(ctx context.Context, firstUserID int, secondUserID int) (string, error) {
	rep, err := s.GetReportIfExists(ctx, firstUserID, secondUserID)
	if err != nil {
		if err.Error() == "this report dont exists" {
			return "", nil
		} else {
			s.logger.Error("Usecase GetReport error", zap.Error(err))
			return "", fmt.Errorf("GetReport error: %w", err)
		}
	}
	if rep.Author == firstUserID {
		return "Вы заблокировали данного пользователя", stderr.New("block exists")
	} else if rep.Author == secondUserID {
		return "Вас заблокировал данный пользователь", stderr.New("block exists")
	}
	return "", nil
}
