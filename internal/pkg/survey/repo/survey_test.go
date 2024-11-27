package repo

import (
	"context"
	_ "database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestAddSurvey(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	logger := zap.NewNop()
	repo := New(db, logger)

	ctx := context.Background()
	survey := models.Survey{
		Author:   1,
		Question: "Sample Question",
		Rating:   5,
		Comment:  "Great!",
	}

	mock.ExpectExec("INSERT INTO survey").
		WithArgs(survey.Author, survey.Question, survey.Rating, survey.Grade, survey.Comment).
		WillReturnResult(sqlmock.NewResult(1, 1))

	id, err := repo.AddSurvey(ctx, survey)
	require.NoError(t, err)
	require.Equal(t, survey.ID, id)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetSurveyInfo(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	logger := zap.NewNop()
	repo := New(db, logger)

	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"author", "question", "rating", "grade", "comment"}).
		AddRow(1, "Q1", 5, 5, "Comment1").
		AddRow(2, "Q2", 4, 4, "Comment2")

	mock.ExpectQuery("SELECT author, question, rating, grade, comment FROM survey").
		WillReturnRows(rows)

	surveys, err := repo.GetSurveyInfo(ctx)
	require.NoError(t, err)
	require.Len(t, surveys, 2)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestAddQuestion(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	logger := zap.NewNop()
	repo := New(db, logger)

	ctx := context.Background()
	question := models.AdminQuestion{
		Content: "Question Content",
		Grade:   5,
	}

	mock.ExpectQuery("INSERT INTO question").
		WithArgs(question.Content, question.Grade).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	id, err := repo.AddQuestion(ctx, question)
	require.NoError(t, err)
	require.Equal(t, 1, id)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateQuestion(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	logger := zap.NewNop()
	repo := New(db, logger)

	ctx := context.Background()
	question := models.AdminQuestion{
		Content: "Updated Content",
		Grade:   4,
	}
	content := "Old Content"

	mock.ExpectQuery("UPDATE question SET content =").
		WithArgs(question.Content, question.Grade, content).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	id, err := repo.UpdateQuestion(ctx, question, content)
	require.NoError(t, err)
	require.Equal(t, 1, id)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteQuestion(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	logger := zap.NewNop()
	repo := New(db, logger)

	ctx := context.Background()
	content := "Content to Delete"

	mock.ExpectExec("DELETE FROM question WHERE content =").
		WithArgs(content).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.DeleteQuestion(ctx, content)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetQuestions(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	logger := zap.NewNop()
	repo := New(db, logger)

	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"content", "grade"}).
		AddRow("Question 1", 5).
		AddRow("Question 2", 4)

	mock.ExpectQuery("SELECT content, grade FROM question").
		WillReturnRows(rows)

	questions, err := repo.GetQuestions(ctx)
	require.NoError(t, err)
	require.Len(t, questions, 2)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestErrorCases(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	logger := zap.NewNop()
	repo := New(db, logger)

	ctx := context.Background()

	// AddSurvey error
	mock.ExpectExec("INSERT INTO survey").
		WillReturnError(errors.New("db error"))
	_, err = repo.AddSurvey(ctx, models.Survey{})
	require.Error(t, err)

	// GetSurveyInfo error
	mock.ExpectQuery("SELECT author, question, rating, grade, comment FROM survey").
		WillReturnError(errors.New("db error"))
	_, err = repo.GetSurveyInfo(ctx)
	require.Error(t, err)

	// AddQuestion error
	mock.ExpectQuery("INSERT INTO question").
		WillReturnError(errors.New("db error"))
	_, err = repo.AddQuestion(ctx, models.AdminQuestion{})
	require.Error(t, err)

	// UpdateQuestion error
	mock.ExpectQuery("UPDATE question SET content =").
		WillReturnError(errors.New("db error"))
	_, err = repo.UpdateQuestion(ctx, models.AdminQuestion{}, "Old Content")
	require.Error(t, err)

	// DeleteQuestion error
	mock.ExpectExec("DELETE FROM question WHERE content =").
		WillReturnError(errors.New("db error"))
	err = repo.DeleteQuestion(ctx, "Content")
	require.Error(t, err)

	// GetQuestions error
	mock.ExpectQuery("SELECT content, grade FROM question").
		WillReturnError(errors.New("db error"))
	_, err = repo.GetQuestions(ctx)
	require.Error(t, err)
}
