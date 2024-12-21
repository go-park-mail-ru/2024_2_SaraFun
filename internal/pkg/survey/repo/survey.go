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
	var surveyID int
	err := repo.DB.QueryRow("INSERT INTO survey (author, question, rating, grade, comment) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		survey.Author, survey.Question, survey.Rating, survey.Grade, survey.Comment).Scan(&surveyID)
	if err != nil {
		repo.logger.Error("bad insert survey", zap.Error(err))
		return -1, err
	}
	repo.logger.Info("success added survey")
	return surveyID, nil
}

func (repo *Storage) GetSurveyInfo(ctx context.Context) ([]models.Survey, error) {
	rows, err := repo.DB.Query("SELECT author, question, rating, grade, comment FROM survey")
	if err != nil {
		repo.logger.Error("bad insert survey", zap.Error(err))
		return nil, fmt.Errorf("bad select survey: %w", err)
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

func (repo *Storage) AddQuestion(ctx context.Context, question models.AdminQuestion) (int, error) {
	var id int
	err := repo.DB.QueryRow("INSERT INTO question (content, grade) VALUES ($1, $2) RETURNING id", question.Content, question.Grade).
		Scan(&id)
	if err != nil {
		repo.logger.Error("bad insert question", zap.Error(err))
		return -1, err
	}
	return id, nil
}

func (repo *Storage) UpdateQuestion(ctx context.Context, question models.AdminQuestion, content string) (int, error) {
	var id int
	err := repo.DB.QueryRow("UPDATE question SET content = $1, grade = $2 WHERE content = $3 RETURNING id",
		question.Content, question.Grade, content).Scan(&id)
	if err != nil {
		repo.logger.Error("bad insert question", zap.Error(err))
		return -1, fmt.Errorf("bad update question: %w", err)
	}
	return id, nil
}

func (repo *Storage) DeleteQuestion(ctx context.Context, content string) error {
	_, err := repo.DB.Exec("DELETE FROM question WHERE content = $1", content)
	if err != nil {
		repo.logger.Error("bad insert question", zap.Error(err))
		return fmt.Errorf("bad delete question: %w", err)
	}
	return nil
}

func (repo *Storage) GetQuestions(ctx context.Context) ([]models.AdminQuestion, error) {
	rows, err := repo.DB.Query("SELECT content, grade FROM question")
	if err != nil {
		repo.logger.Error("bad insert question", zap.Error(err))
		return nil, fmt.Errorf("bad get questions: %w", err)
	}
	defer rows.Close()
	var questions []models.AdminQuestion
	for rows.Next() {
		question := models.AdminQuestion{}
		err = rows.Scan(&question.Content, &question.Grade)
		if err != nil {
			repo.logger.Error("bad get question", zap.Error(err))
			return nil, fmt.Errorf("bad get question: %w", err)
		}
		questions = append(questions, question)
	}
	return questions, nil
}
