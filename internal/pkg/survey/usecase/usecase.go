package usecase

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	"go.uber.org/zap"
)

type Repository interface {
	AddSurvey(ctx context.Context, survey models.Survey) (int, error)
	GetSurveyInfo(ctx context.Context) ([]models.Survey, error)
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

func (u *UseCase) AddSurvey(ctx context.Context, survey models.Survey) (int, error) {
	surveyID, err := u.repo.AddSurvey(ctx, survey)
	if err != nil {
		u.logger.Error("bad add survey", zap.Error(err))
		return -1, err
	}
	return surveyID, nil
}

func (u *UseCase) GetSurveyInfo(ctx context.Context) (map[string]models.SurveyStat, error) {
	stats := make(map[string]models.SurveyStat)
	surveys, err := u.repo.GetSurveyInfo(ctx)
	if err != nil {
		u.logger.Error("bad get survey info", zap.Error(err))
		return stats, fmt.Errorf("get survey info: %w", err)
	}
	for _, survey := range surveys {
		if _, ok := stats[survey.Question]; ok == false {
			_stat := models.SurveyStat{
				Question: survey.Question,
				Grade:    survey.Grade,
				Count:    1,
			}
			u.logger.Info("test")
			stats[survey.Question] = _stat
		}
		ss := stats[survey.Question]
		ss.Count = ss.Count + 1
		ss.Sum = ss.Sum + survey.Rating
		ss.Rating = float32(ss.Sum) / float32(ss.Count)
		stats[survey.Question] = ss
	}
	return stats, nil
}
