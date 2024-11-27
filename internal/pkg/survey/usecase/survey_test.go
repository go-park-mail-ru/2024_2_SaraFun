//go:generate mockgen -source=usecase.go -destination=./mocks/mock_repository.go -package=mocks

package usecase

import (
	"context"
	"errors"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/survey/usecase/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
)

func TestUseCase_AddSurvey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	logger := zap.NewNop()

	ctx := context.Background()
	useCase := New(mockRepo, logger)

	survey := models.Survey{ID: 1, Question: "Test Question", Rating: 5}
	mockRepo.EXPECT().AddSurvey(ctx, survey).Return(1, nil).Times(1)

	id, err := useCase.AddSurvey(ctx, survey)
	require.NoError(t, err)
	require.Equal(t, 1, id)
}

func TestUseCase_AddQuestion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	logger := zap.NewNop()

	ctx := context.Background()
	useCase := New(mockRepo, logger)

	question := models.AdminQuestion{Content: "Test Question"}
	mockRepo.EXPECT().AddQuestion(ctx, question).Return(1, nil).Times(1)

	id, err := useCase.AddQuestion(ctx, question)
	require.NoError(t, err)
	require.Equal(t, 1, id)
}

func TestUseCase_UpdateQuestion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	logger := zap.NewNop()

	ctx := context.Background()
	useCase := New(mockRepo, logger)

	question := models.AdminQuestion{Content: "Updated Question"}
	mockRepo.EXPECT().UpdateQuestion(ctx, question, "Old Content").Return(1, nil).Times(1)

	id, err := useCase.UpdateQuestion(ctx, question, "Old Content")
	require.NoError(t, err)
	require.Equal(t, 1, id)
}

func TestUseCase_DeleteQuestion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	logger := zap.NewNop()

	ctx := context.Background()
	useCase := New(mockRepo, logger)

	mockRepo.EXPECT().DeleteQuestion(ctx, "Test Question").Return(nil).Times(1)

	err := useCase.DeleteQuestion(ctx, "Test Question")
	require.NoError(t, err)
}

func TestUseCase_GetQuestions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	logger := zap.NewNop()

	ctx := context.Background()
	useCase := New(mockRepo, logger)

	questions := []models.AdminQuestion{
		{Content: "Question 1"},
		{Content: "Question 2"},
	}
	mockRepo.EXPECT().GetQuestions(ctx).Return(questions, nil).Times(1)

	result, err := useCase.GetQuestions(ctx)
	require.NoError(t, err)
	require.Equal(t, questions, result)
}

func TestUseCase_ErrorCases(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepository(ctrl)
	logger := zap.NewNop()

	ctx := context.Background()
	useCase := New(mockRepo, logger)

	// AddSurvey error
	mockRepo.EXPECT().AddSurvey(ctx, gomock.Any()).Return(-1, errors.New("db error")).Times(1)
	_, err := useCase.AddSurvey(ctx, models.Survey{})
	require.Error(t, err)

	// GetSurveyInfo error
	mockRepo.EXPECT().GetSurveyInfo(ctx).Return(nil, errors.New("db error")).Times(1)
	_, err = useCase.GetSurveyInfo(ctx)
	require.Error(t, err)

	// AddQuestion error
	mockRepo.EXPECT().AddQuestion(ctx, gomock.Any()).Return(-1, errors.New("db error")).Times(1)
	_, err = useCase.AddQuestion(ctx, models.AdminQuestion{})
	require.Error(t, err)

	// UpdateQuestion error
	mockRepo.EXPECT().UpdateQuestion(ctx, gomock.Any(), gomock.Any()).Return(-1, errors.New("db error")).Times(1)
	_, err = useCase.UpdateQuestion(ctx, models.AdminQuestion{}, "Old Content")
	require.Error(t, err)

	// DeleteQuestion error
	mockRepo.EXPECT().DeleteQuestion(ctx, gomock.Any()).Return(errors.New("db error")).Times(1)
	err = useCase.DeleteQuestion(ctx, "Test Question")
	require.Error(t, err)

	// GetQuestions error
	mockRepo.EXPECT().GetQuestions(ctx).Return(nil, errors.New("db error")).Times(1)
	_, err = useCase.GetQuestions(ctx)
	require.Error(t, err)
}
