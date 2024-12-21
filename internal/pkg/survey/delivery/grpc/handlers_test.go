package grpc_test

import (
	"context"
	"errors"

	"testing"

	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/survey/delivery/grpc"
	generatedSurvey "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/survey/delivery/grpc/gen"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/survey/delivery/grpc/mocks"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/utils/consts"

	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
)

func TestGRPCHandler_AddSurvey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uc := mocks.NewMockSurveyUsecase(ctrl)
	logger := zap.NewNop()
	h := grpc.New(logger, uc)

	ctx := context.WithValue(context.Background(), consts.RequestIDKey, "test_req_id")
	in := &generatedSurvey.AddSurveyRequest{
		Survey: &generatedSurvey.SSurvey{
			Author:   10,
			Question: "Q?",
			Comment:  "C",
			Rating:   5,
			Grade:    2,
		},
	}

	t.Run("success", func(t *testing.T) {
		uc.EXPECT().AddSurvey(ctx, models.Survey{
			Author:   10,
			Question: "Q?",
			Comment:  "C",
			Rating:   5,
			Grade:    2,
		}).Return(123, nil)

		resp, err := h.AddSurvey(ctx, in)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if resp.SurveyID != 123 {
			t.Errorf("got %v, want %v", resp.SurveyID, 123)
		}
	})

	t.Run("error", func(t *testing.T) {
		uc.EXPECT().AddSurvey(ctx, gomock.Any()).Return(0, errors.New("add error"))

		_, err := h.AddSurvey(ctx, in)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
		if !contains(err.Error(), "bad add survey grpc: add error") {
			t.Errorf("error mismatch: got %v", err)
		}
	})
}

func TestGRPCHandler_GetSurveyInfo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uc := mocks.NewMockSurveyUsecase(ctrl)
	logger := zap.NewNop()
	h := grpc.New(logger, uc)

	ctx := context.WithValue(context.Background(), consts.RequestIDKey, "test_req_id")

	t.Run("success", func(t *testing.T) {
		stats := map[string]models.SurveyStat{
			"Q1": {Question: "Q1", Rating: 4.5, Grade: 2},
			"Q2": {Question: "Q2", Rating: 3.0, Grade: 1},
		}
		uc.EXPECT().GetSurveyInfo(ctx).Return(stats, nil)

		in := &generatedSurvey.GetSurveyInfoRequest{}
		resp, err := h.GetSurveyInfo(ctx, in)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if len(resp.Stats) != 2 {
			t.Errorf("got %d stats, want 2", len(resp.Stats))
		}
	})

	t.Run("error", func(t *testing.T) {
		uc.EXPECT().GetSurveyInfo(ctx).Return(nil, errors.New("info error"))

		in := &generatedSurvey.GetSurveyInfoRequest{}
		_, err := h.GetSurveyInfo(ctx, in)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
		if !contains(err.Error(), "get survey info: info error") {
			t.Errorf("error mismatch: got %v", err)
		}
	})
}

func TestGRPCHandler_AddQuestion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uc := mocks.NewMockSurveyUsecase(ctrl)
	logger := zap.NewNop()
	h := grpc.New(logger, uc)

	ctx := context.WithValue(context.Background(), consts.RequestIDKey, "test_req_id")
	in := &generatedSurvey.AddQuestionRequest{
		Question: &generatedSurvey.AdminQuestion{
			Content: "Q?",
			Grade:   2,
		},
	}

	t.Run("success", func(t *testing.T) {
		uc.EXPECT().AddQuestion(ctx, models.AdminQuestion{Content: "Q?", Grade: 2}).Return(456, nil)
		resp, err := h.AddQuestion(ctx, in)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if resp.QuestionID != 456 {
			t.Errorf("got %v, want %v", resp.QuestionID, 456)
		}
	})

	t.Run("error", func(t *testing.T) {
		uc.EXPECT().AddQuestion(ctx, gomock.Any()).Return(0, errors.New("q error"))
		_, err := h.AddQuestion(ctx, in)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
		if !contains(err.Error(), "add question grpc: q error") {
			t.Errorf("error mismatch: got %v", err)
		}
	})
}

func TestGRPCHandler_UpdateQuestion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uc := mocks.NewMockSurveyUsecase(ctrl)
	logger := zap.NewNop()
	h := grpc.New(logger, uc)

	ctx := context.WithValue(context.Background(), consts.RequestIDKey, "test_req_id")
	in := &generatedSurvey.UpdateQuestionRequest{
		Question: &generatedSurvey.AdminQuestion{
			Content: "Q new",
			Grade:   3,
		},
		Content: "Q old",
	}

	t.Run("success", func(t *testing.T) {
		uc.EXPECT().UpdateQuestion(ctx, models.AdminQuestion{Content: "Q new", Grade: 3}, "Q old").Return(789, nil)
		resp, err := h.UpdateQuestion(ctx, in)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if resp.Id != 789 {
			t.Errorf("got %v, want %v", resp.Id, 789)
		}
	})

	t.Run("error", func(t *testing.T) {
		uc.EXPECT().UpdateQuestion(ctx, gomock.Any(), "Q old").Return(0, errors.New("update error"))
		_, err := h.UpdateQuestion(ctx, in)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
		if !contains(err.Error(), "update question grpc: update error") {
			t.Errorf("error mismatch: got %v", err)
		}
	})
}

func TestGRPCHandler_DeleteQuestion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uc := mocks.NewMockSurveyUsecase(ctrl)
	logger := zap.NewNop()
	h := grpc.New(logger, uc)

	ctx := context.WithValue(context.Background(), consts.RequestIDKey, "test_req_id")
	in := &generatedSurvey.DeleteQuestionRequest{
		Content: "Q?",
	}

	t.Run("success", func(t *testing.T) {
		uc.EXPECT().DeleteQuestion(ctx, "Q?").Return(nil)
		resp, err := h.DeleteQuestion(ctx, in)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if resp == nil {
			t.Errorf("expected response, got nil")
		}
	})

	t.Run("error", func(t *testing.T) {
		uc.EXPECT().DeleteQuestion(ctx, "Q?").Return(errors.New("delete error"))
		_, err := h.DeleteQuestion(ctx, in)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
		if !contains(err.Error(), "delete question grpc: delete error") {
			t.Errorf("error mismatch: got %v", err)
		}
	})
}

func TestGRPCHandler_GetQuestions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uc := mocks.NewMockSurveyUsecase(ctrl)
	logger := zap.NewNop()
	h := grpc.New(logger, uc)

	ctx := context.WithValue(context.Background(), consts.RequestIDKey, "test_req_id")
	in := &generatedSurvey.GetQuestionsRequest{}

	t.Run("success", func(t *testing.T) {
		questions := []models.AdminQuestion{
			{Content: "Q1", Grade: 1},
			{Content: "Q2", Grade: 2},
		}
		uc.EXPECT().GetQuestions(ctx).Return(questions, nil)
		resp, err := h.GetQuestions(ctx, in)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if len(resp.Questions) != 2 {
			t.Errorf("got %d questions, want 2", len(resp.Questions))
		}
	})

	t.Run("error", func(t *testing.T) {
		uc.EXPECT().GetQuestions(ctx).Return(nil, errors.New("get error"))
		_, err := h.GetQuestions(ctx, in)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
		if !contains(err.Error(), "get questions grpc: get error") {
			t.Errorf("error mismatch: got %v", err)
		}
	})
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchSubstring(s, substr)
}

func searchSubstring(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
