package usecase

import (
	"context"
	"errors"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/survey/usecase/mocks"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
	"time"
)

func TestAddSurvey(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, consts.RequestIDKey, "sparkit")
	logger := zap.NewNop()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	repo := mocks.NewMockRepository(mockCtrl)

	tests := []struct {
		name             string
		survey           models.Survey
		repoSurveyID     int
		repoError        error
		repoCount        int
		expectedSurveyID int
	}{
		{
			name: "successfull test",
			survey: models.Survey{
				Author:   1,
				Question: "test?",
				Comment:  "test",
				Rating:   5,
				Grade:    5,
			},
			repoSurveyID:     1,
			repoError:        nil,
			repoCount:        1,
			expectedSurveyID: 1,
		},
		{
			name: "bad test",
			survey: models.Survey{
				Author: 1,
			},
			repoSurveyID:     -1,
			repoError:        errors.New("error"),
			repoCount:        1,
			expectedSurveyID: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo.EXPECT().AddSurvey(ctx, tt.survey).Return(tt.repoSurveyID, tt.repoError).Times(tt.repoCount)
			usecase := New(repo, logger)
			id, err := usecase.AddSurvey(ctx, tt.survey)
			require.ErrorIs(t, err, tt.repoError)
			require.Equal(t, tt.expectedSurveyID, id)
		})
	}

}

func TestGetSurveyInfo(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, consts.RequestIDKey, "sparkit")
	logger := zap.NewNop()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	repo := mocks.NewMockRepository(mockCtrl)
	tests := []struct {
		name          string
		repoSurveys   []models.Survey
		repoError     error
		repoCount     int
		expectedStats map[string]models.SurveyStat
	}{
		{
			name: "successfull test",
			repoSurveys: []models.Survey{
				{
					Author:   1,
					Question: "test?",
					Comment:  "test",
					Rating:   5,
					Grade:    5,
				},
				{
					Author:   2,
					Question: "test?",
					Comment:  "test",
					Rating:   3,
					Grade:    5,
				},
			},
			repoError: nil,
			repoCount: 1,
			expectedStats: map[string]models.SurveyStat{
				"test?": {
					Question: "test?",
					Grade:    5,
					Rating:   4,
					Sum:      8,
					Count:    2,
				},
			},
		},
		{
			name:          "bad test",
			repoSurveys:   nil,
			repoError:     errors.New("error"),
			repoCount:     1,
			expectedStats: map[string]models.SurveyStat{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo.EXPECT().GetSurveyInfo(ctx).Return(tt.repoSurveys, tt.repoError).Times(tt.repoCount)
			usecase := New(repo, logger)
			stats, err := usecase.GetSurveyInfo(ctx)
			require.ErrorIs(t, err, tt.repoError)
			require.Equal(t, tt.expectedStats, stats)
		})
	}
}

func TestAddQuestion(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, consts.RequestIDKey, "sparkit")
	logger := zap.NewNop()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	repo := mocks.NewMockRepository(mockCtrl)
	tests := []struct {
		name           string
		question       models.AdminQuestion
		repoQuestionID int
		repoError      error
		repoCount      int
		expectedID     int
	}{
		{
			name: "successfull test",
			question: models.AdminQuestion{
				Content: "Насколько вам нравится наш сервис?",
				Grade:   5,
			},
			repoQuestionID: 1,
			repoError:      nil,
			repoCount:      1,
			expectedID:     1,
		},
		{
			name: "bad test",
			question: models.AdminQuestion{
				Content: "",
				Grade:   0,
			},
			repoQuestionID: -1,
			repoError:      errors.New("error"),
			repoCount:      1,
			expectedID:     -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo.EXPECT().AddQuestion(ctx, tt.question).Return(tt.repoQuestionID, tt.repoError).Times(tt.repoCount)
			usecase := New(repo, logger)
			id, err := usecase.AddQuestion(ctx, tt.question)
			require.ErrorIs(t, err, tt.repoError)
			require.Equal(t, tt.expectedID, id)
		})
	}
}

func TestUpdateQuestion(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, consts.RequestIDKey, "sparkit")
	logger := zap.NewNop()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	repo := mocks.NewMockRepository(mockCtrl)
	tests := []struct {
		name           string
		question       models.AdminQuestion
		content        string
		repoQuestionID int
		repoError      error
		repoCount      int
		expectedID     int
	}{
		{
			name: "successfull test",
			question: models.AdminQuestion{
				Content: "Насколько вам нравится наш сервис?",
				Grade:   5,
			},
			content:        "Насколько?",
			repoQuestionID: 1,
			repoError:      nil,
			repoCount:      1,
			expectedID:     1,
		},
		{
			name: "bad test",
			question: models.AdminQuestion{
				Content: "",
				Grade:   0,
			},
			content:        "",
			repoQuestionID: -1,
			repoError:      errors.New("error"),
			repoCount:      1,
			expectedID:     -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo.EXPECT().UpdateQuestion(ctx, tt.question, tt.content).Return(tt.repoQuestionID, tt.repoError).Times(tt.repoCount)
			usecase := New(repo, logger)
			id, err := usecase.UpdateQuestion(ctx, tt.question, tt.content)
			require.ErrorIs(t, err, tt.repoError)
			require.Equal(t, tt.expectedID, id)
		})
	}
}

func TestDeleteQuestion(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, consts.RequestIDKey, "sparkit")
	logger := zap.NewNop()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	repo := mocks.NewMockRepository(mockCtrl)
	tests := []struct {
		name      string
		content   string
		repoError error
		repoCount int
	}{
		{
			name:      "successfull test",
			content:   "Насколько вам нравятся свайпы?",
			repoError: nil,
			repoCount: 1,
		},
		{
			name:      "bad test",
			content:   "",
			repoError: errors.New("error"),
			repoCount: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo.EXPECT().DeleteQuestion(ctx, tt.content).Return(tt.repoError).Times(tt.repoCount)
			usecase := New(repo, logger)
			err := usecase.DeleteQuestion(ctx, tt.content)
			require.ErrorIs(t, err, tt.repoError)
		})
	}
}

func TestGetQuestions(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ctx = context.WithValue(ctx, consts.RequestIDKey, "sparkit")
	logger := zap.NewNop()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	repo := mocks.NewMockRepository(mockCtrl)
	tests := []struct {
		name              string
		repoQuestions     []models.AdminQuestion
		repoError         error
		repoCount         int
		expectedQuestions []models.AdminQuestion
	}{
		{
			name: "successfull test",
			repoQuestions: []models.AdminQuestion{
				{
					Content: "Насколько вам нравятся свайпы?",
					Grade:   5,
				},
			},
			repoError: nil,
			repoCount: 1,
			expectedQuestions: []models.AdminQuestion{
				{
					Content: "Насколько вам нравятся свайпы?",
					Grade:   5,
				},
			},
		},
		{
			name: "bad test",
			repoQuestions: []models.AdminQuestion{
				{
					Content: "",
					Grade:   0,
				},
			},
			repoError:         errors.New("error"),
			repoCount:         1,
			expectedQuestions: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo.EXPECT().GetQuestions(ctx).Return(tt.repoQuestions, tt.repoError).Times(tt.repoCount)
			usecase := New(repo, logger)
			questions, err := usecase.GetQuestions(ctx)
			require.ErrorIs(t, err, tt.repoError)
			require.Equal(t, tt.expectedQuestions, questions)
		})
	}
}
