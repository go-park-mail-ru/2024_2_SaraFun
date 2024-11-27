package grpc

import (
	"context"
	"testing"

	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	generatedSurvey "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/survey/delivery/grpc/gen"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/survey/delivery/grpc/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestGetSurveyInfo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockSurveyUsecase(ctrl)
	logger := zap.NewNop()
	handler := New(logger, mockUsecase)

	ctx := context.Background()
	stats := map[string]models.SurveyStat{
		"Question1": {Question: "Question1", Grade: 3, Rating: 4.5},
		"Question2": {Question: "Question2", Grade: 4, Rating: 3.8},
	}

	mockUsecase.EXPECT().GetSurveyInfo(ctx).Return(stats, nil)

	req := &generatedSurvey.GetSurveyInfoRequest{}
	resp, err := handler.GetSurveyInfo(ctx, req)
	require.NoError(t, err)
	require.Len(t, resp.Stats, 2)
}

func TestAddQuestion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockSurveyUsecase(ctrl)
	logger := zap.NewNop()
	handler := New(logger, mockUsecase)

	ctx := context.Background()
	question := models.AdminQuestion{
		Content: "Sample Question",
		Grade:   3,
	}

	mockUsecase.EXPECT().AddQuestion(ctx, question).Return(1, nil)

	req := &generatedSurvey.AddQuestionRequest{
		Question: &generatedSurvey.AdminQuestion{
			Content: "Sample Question",
			Grade:   3,
		},
	}
	resp, err := handler.AddQuestion(ctx, req)
	require.NoError(t, err)
	require.Equal(t, int32(1), resp.QuestionID)
}

func TestUpdateQuestion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockSurveyUsecase(ctrl)
	logger := zap.NewNop()
	handler := New(logger, mockUsecase)

	ctx := context.Background()
	question := models.AdminQuestion{
		Content: "Updated Question",
		Grade:   3,
	}

	mockUsecase.EXPECT().UpdateQuestion(ctx, question, "Old Question").Return(1, nil)

	req := &generatedSurvey.UpdateQuestionRequest{
		Question: &generatedSurvey.AdminQuestion{
			Content: "Updated Question",
			Grade:   3,
		},
		Content: "Old Question",
	}
	resp, err := handler.UpdateQuestion(ctx, req)
	require.NoError(t, err)
	require.Equal(t, int32(1), resp.Id)
}

func TestDeleteQuestion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockSurveyUsecase(ctrl)
	logger := zap.NewNop()
	handler := New(logger, mockUsecase)

	ctx := context.Background()
	content := "Question to delete"

	mockUsecase.EXPECT().DeleteQuestion(ctx, content).Return(nil)

	req := &generatedSurvey.DeleteQuestionRequest{
		Content: content,
	}
	_, err := handler.DeleteQuestion(ctx, req)
	require.NoError(t, err)
}

func TestGetQuestions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockSurveyUsecase(ctrl)
	logger := zap.NewNop()
	handler := New(logger, mockUsecase)

	ctx := context.Background()
	questions := []models.AdminQuestion{
		{Content: "Question1", Grade: 3},
		{Content: "Question2", Grade: 4},
	}

	mockUsecase.EXPECT().GetQuestions(ctx).Return(questions, nil)

	req := &generatedSurvey.GetQuestionsRequest{}
	resp, err := handler.GetQuestions(ctx, req)
	require.NoError(t, err)
	require.Len(t, resp.Questions, 2)
}
