package usecase

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	"go.uber.org/zap"
)

//go:generate mockgen -destination=./mocks/mock_repository.go -package=mocks . Repository
type Repository interface {
	AddSurvey(ctx context.Context, survey models.Survey) (int, error)
	GetSurveyInfo(ctx context.Context) ([]models.Survey, error)
	AddQuestion(ctx context.Context, question models.AdminQuestion) (int, error)
	DeleteQuestion(ctx context.Context, content string) error
	UpdateQuestion(ct context.Context, question models.AdminQuestion, content string) (int, error)
	GetQuestions(ctx context.Context) ([]models.AdminQuestion, error)
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
		if _, ok := stats[survey.Question]; !ok {
			_stat := models.SurveyStat{
				Question: survey.Question,
				Grade:    survey.Grade,
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

func (u *UseCase) AddQuestion(ctx context.Context, question models.AdminQuestion) (int, error) {
	id, err := u.repo.AddQuestion(ctx, question)
	if err != nil {
		u.logger.Error("bad add question", zap.Error(err))
		return -1, fmt.Errorf("add question usecase: %w", err)
	}
	return id, nil
}

func (u *UseCase) UpdateQuestion(ctx context.Context, question models.AdminQuestion, content string) (int, error) {
	id, err := u.repo.UpdateQuestion(ctx, question, content)
	if err != nil {
		u.logger.Error("bad update question", zap.Error(err))
		return -1, fmt.Errorf("update question usecase: %w", err)
	}
	return id, nil
}

func (u *UseCase) DeleteQuestion(ctx context.Context, content string) error {
	err := u.repo.DeleteQuestion(ctx, content)
	if err != nil {
		u.logger.Error("bad delete question", zap.Error(err))
		return fmt.Errorf("delete question usecase: %w", err)
	}
	return nil
}

func (u *UseCase) GetQuestions(ctx context.Context) ([]models.AdminQuestion, error) {
	var questions []models.AdminQuestion
	questions, err := u.repo.GetQuestions(ctx)
	if err != nil {
		u.logger.Error("bad get questions", zap.Error(err))
		return nil, fmt.Errorf("get questions usecase: %w", err)
	}
	return questions, nil
}
