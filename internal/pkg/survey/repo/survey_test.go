package repo

import (
	"context"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
	"time"
)

func TestAddSurvey(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctx = context.WithValue(ctx, consts.RequestIDKey, "sparkit")
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock new error: %s", err)
	}
	defer db.Close()

	logger := zap.NewNop()
	repo := New(db, logger)

	tests := []struct {
		name       string
		survey     models.Survey
		queryID    int
		queryError error
		expectedID int
	}{
		{
			name: "successfull test",
			survey: models.Survey{
				Author:   1,
				Question: "test?",
				Comment:  "test comment",
				Rating:   5,
				Grade:    5,
			},
			queryID:    1,
			queryError: nil,
			expectedID: 1,
		},
		{
			name: "bad test",
			survey: models.Survey{
				Author:   1,
				Question: "test?",
				Comment:  "test comment",
				Rating:   3,
				Grade:    5,
			},
			queryID:    0,
			queryError: errors.New("error"),
			expectedID: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.queryError == nil {
				mock.ExpectQuery("INSERT INTO survey").
					WithArgs(tt.survey.Author, tt.survey.Question, tt.survey.Rating, tt.survey.Grade, tt.survey.Comment).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(tt.expectedID))
			} else {
				mock.ExpectQuery("INSERT INTO survey").
					WithArgs(tt.survey.Author, tt.survey.Question, tt.survey.Rating, tt.survey.Grade, tt.survey.Comment).
					WillReturnError(tt.queryError)
			}
			id, err := repo.AddSurvey(ctx, tt.survey)
			require.ErrorIs(t, tt.queryError, err)
			require.Equal(t, tt.expectedID, id)
		})
	}

}

func TestGetSurveyInfo(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctx = context.WithValue(ctx, consts.RequestIDKey, "sparkit")
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock new error: %s", err)
	}
	defer db.Close()
	logger := zap.NewNop()
	repo := New(db, logger)
	tests := []struct {
		name            string
		queryRows       *sqlmock.Rows
		queryError      error
		expectedSurveys []models.Survey
	}{
		{
			name: "successfull test",
			queryRows: mock.NewRows([]string{"author", "question", "rating", "grade", "comment"}).
				AddRow(1, "test?", 5, 5, "test comment"),
			queryError: nil,
			expectedSurveys: []models.Survey{
				{
					Author:   1,
					Question: "test?",
					Comment:  "test comment",
					Rating:   5,
					Grade:    5,
				},
			},
		},
		{
			name:            "bad test",
			queryRows:       nil,
			queryError:      errors.New("error"),
			expectedSurveys: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.queryError == nil {
				mock.ExpectQuery("SELECT").WillReturnRows(tt.queryRows)
			} else {
				mock.ExpectQuery("SELECT").WillReturnError(tt.queryError)
			}

			surveys, err := repo.GetSurveyInfo(ctx)
			require.ErrorIs(t, err, tt.queryError)
			require.Equal(t, tt.expectedSurveys, surveys)
		})
	}
}

func TestAddQuestion(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, consts.RequestIDKey, "sparkit")
	logger := zap.NewNop()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock new error: %s", err)
	}
	defer db.Close()
	repo := New(db, logger)

	tests := []struct {
		name       string
		question   models.AdminQuestion
		queryID    int
		queryError error
		expectedID int
	}{
		{
			name: "successfull test",
			question: models.AdminQuestion{
				Content: "test content",
				Grade:   5,
			},
			queryID:    1,
			queryError: nil,
			expectedID: 1,
		},
		{
			name:       "bad test",
			question:   models.AdminQuestion{},
			queryID:    0,
			queryError: errors.New("error"),
			expectedID: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.queryError == nil {
				mock.ExpectQuery("INSERT INTO question").
					WithArgs(tt.question.Content, tt.question.Grade).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(tt.expectedID))
			} else {
				mock.ExpectQuery("INSERT INTO question").
					WithArgs(tt.question.Content, tt.question.Grade).
					WillReturnError(tt.queryError)
			}
			id, err := repo.AddQuestion(ctx, tt.question)
			require.ErrorIs(t, err, tt.queryError)
			require.Equal(t, tt.expectedID, id)
		})
	}
}

func TestUpdateQuestion(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, consts.RequestIDKey, "sparkit")
	logger := zap.NewNop()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock new error: %s", err)
	}
	defer db.Close()
	repo := New(db, logger)
	tests := []struct {
		name       string
		question   models.AdminQuestion
		content    string
		queryID    int
		queryError error
		expectedID int
	}{
		{
			name: "successfull test",
			question: models.AdminQuestion{
				Content: "test content",
				Grade:   5,
			},
			content:    "test",
			queryID:    1,
			queryError: nil,
			expectedID: 1,
		},
		{
			name:       "bad test",
			question:   models.AdminQuestion{},
			content:    "test",
			queryID:    0,
			queryError: errors.New("error"),
			expectedID: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.queryError == nil {
				mock.ExpectQuery("UPDATE question").
					WithArgs(tt.question.Content, tt.question.Grade, tt.content).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(tt.expectedID))
			} else {
				mock.ExpectQuery("UPDATE question").
					WithArgs(tt.question.Content, tt.question.Grade, tt.content).
					WillReturnError(tt.queryError)
			}

			id, err := repo.UpdateQuestion(ctx, tt.question, tt.content)
			require.ErrorIs(t, err, tt.queryError)
			require.Equal(t, tt.expectedID, id)
		})
	}
}

func TestDeleteQuestion(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, consts.RequestIDKey, "sparkit")
	logger := zap.NewNop()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock new error: %s", err)
	}
	defer db.Close()
	repo := New(db, logger)
	tests := []struct {
		name       string
		content    string
		queryError error
	}{
		{
			name:       "successfull test",
			content:    "test",
			queryError: nil,
		},
		{
			name:       "bad test",
			content:    "test",
			queryError: errors.New("error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectExec("DELETE FROM question").
				WithArgs(tt.content).
				WillReturnResult(sqlmock.NewResult(0, 0)).
				WillReturnError(tt.queryError)

			err := repo.DeleteQuestion(ctx, tt.content)
			require.ErrorIs(t, err, tt.queryError)
		})
	}
}

func TestGetQuestions(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, consts.RequestIDKey, "sparkit")
	logger := zap.NewNop()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock new error: %s", err)
	}
	defer db.Close()
	repo := New(db, logger)
	tests := []struct {
		name              string
		queryRows         *sqlmock.Rows
		queryError        error
		expectedQuestions []models.AdminQuestion
	}{
		{
			name:       "successfull test",
			queryRows:  sqlmock.NewRows([]string{"content", "grade"}).AddRow("test", 5),
			queryError: nil,
			expectedQuestions: []models.AdminQuestion{
				{
					Content: "test",
					Grade:   5,
				},
			},
		},
		{
			name:              "bad test",
			queryRows:         nil,
			queryError:        errors.New("error"),
			expectedQuestions: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.queryError == nil {
				mock.ExpectQuery("SELECT content, grade FROM question").
					WillReturnRows(tt.queryRows)
			} else {
				mock.ExpectQuery("SELECT content, grade FROM question").
					WillReturnError(tt.queryError)
			}
			questions, err := repo.GetQuestions(ctx)
			require.ErrorIs(t, err, tt.queryError)
			require.Equal(t, tt.expectedQuestions, questions)
		})
	}
}
