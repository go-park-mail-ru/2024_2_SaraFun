package repo

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	"go.uber.org/zap"
)

type Storage struct {
	DB     *sql.DB
	logger *zap.Logger
}

func New(db *sql.DB, logger *zap.Logger) *Storage {
	return &Storage{
		DB:     db,
		logger: logger,
	}
}

func (repo *Storage) AddSurvey(ctx context.Context, survey models.Survey) (int, error) {
	_, err := repo.DB.Exec("INSERT INTO survey (author, question, rating, grade, comment) VALUES ($1, $2, $3, $4, $5)",
		survey.Author, survey.Question, survey.Rating, survey.Grade, survey.Comment)
	if err != nil {
		repo.logger.Error("bad insert survey", zap.Error(err))
		return -1, err
	}
	repo.logger.Info("success added survey")
	return survey.ID, nil
}

func (repo *Storage) GetSurveyInfo(ctx context.Context) ([]models.Survey, error) {
	rows, err := repo.DB.Query("SELECT author, question, rating, grade, comment FROM survey")
	if err != nil {
		repo.logger.Error("bad insert survey", zap.Error(err))
		return nil, fmt.Errorf("bad select survey: %v", err)
	}
	defer rows.Close()
	var surveys []models.Survey

	for rows.Next() {
		survey := models.Survey{}
		err = rows.Scan(&survey.Author, &survey.Question, &survey.Rating, &survey.Grade, &survey.Comment)
		if err != nil {
			repo.logger.Error("bad scan survey", zap.Error(err))
			return nil, fmt.Errorf("bad scan survey: %v", err)
		}
		surveys = append(surveys, survey)
	}
	return surveys, nil
}
