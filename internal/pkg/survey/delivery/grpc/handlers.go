package grpc

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	generatedSurvey "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/survey/delivery/grpc/gen"
	"go.uber.org/zap"
)

//go:generate mockgen -destination=./mocks/mock_surveyusecase.go -package=mocks github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/survey/delivery/grpc SurveyUsecase

type SurveyUsecase interface {
	AddSurvey(ctx context.Context, survey models.Survey) (int, error)
	GetSurveyInfo(ctx context.Context) (map[string]models.SurveyStat, error)
	AddQuestion(ctx context.Context, question models.AdminQuestion) (int, error)
	DeleteQuestion(ctx context.Context, content string) error
	UpdateQuestion(ctx context.Context, question models.AdminQuestion, content string) (int, error)
	GetQuestions(ctx context.Context) ([]models.AdminQuestion, error)
}

type GRPCHandler struct {
	generatedSurvey.SurveyServer
	uc     SurveyUsecase
	logger *zap.Logger
}

func New(logger *zap.Logger, uc SurveyUsecase) *GRPCHandler {
	return &GRPCHandler{
		uc:     uc,
		logger: logger,
	}
}

func (h *GRPCHandler) AddSurvey(ctx context.Context, in *generatedSurvey.AddSurveyRequest) (*generatedSurvey.AddSurveyResponse, error) {
	survey := models.Survey{
		Author:   int(in.Survey.Author),
		Question: in.Survey.Question,
		Comment:  in.Survey.Comment,
		Rating:   int(in.Survey.Rating),
		Grade:    int(in.Survey.Grade),
	}
	id, err := h.uc.AddSurvey(ctx, survey)
	if err != nil {
		return nil, fmt.Errorf("bad add survey grpc: %w", err)
	}
	response := &generatedSurvey.AddSurveyResponse{
		SurveyID: int32(id),
	}
	return response, nil
}

func (h *GRPCHandler) GetSurveyInfo(ctx context.Context, in *generatedSurvey.GetSurveyInfoRequest) (*generatedSurvey.GetSurveyInfoResponse, error) {
	stats, err := h.uc.GetSurveyInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("get survey info: %w", err)
	}
	var respStats []*generatedSurvey.Stat
	for _, stat := range stats {
		respStats = append(respStats, &generatedSurvey.Stat{
			Question:  stat.Question,
			AvgRating: stat.Rating,
			Grade:     int32(stat.Grade),
		})
	}
	response := &generatedSurvey.GetSurveyInfoResponse{
		Stats: respStats,
	}
	return response, nil
}

func (h *GRPCHandler) AddQuestion(ctx context.Context, in *generatedSurvey.AddQuestionRequest) (*generatedSurvey.AddQuestionResponse, error) {
	question := models.AdminQuestion{
		Content: in.Question.Content,
		Grade:   int(in.Question.Grade),
	}
	id, err := h.uc.AddQuestion(ctx, question)
	if err != nil {
		return nil, fmt.Errorf("add question grpc: %w", err)
	}
	response := &generatedSurvey.AddQuestionResponse{
		QuestionID: int32(id),
	}
	return response, nil
}

func (h *GRPCHandler) UpdateQuestion(ctx context.Context, in *generatedSurvey.UpdateQuestionRequest) (*generatedSurvey.UpdateQuestionResponse, error) {
	question := models.AdminQuestion{
		Content: in.Question.Content,
		Grade:   int(in.Question.Grade),
	}
	h.logger.Info("content", zap.String("content", in.Content))
	id, err := h.uc.UpdateQuestion(ctx, question, in.Content)
	if err != nil {
		return nil, fmt.Errorf("update question grpc: %w", err)
	}
	response := &generatedSurvey.UpdateQuestionResponse{
		Id: int32(id),
	}
	return response, nil
}

func (h *GRPCHandler) DeleteQuestion(ctx context.Context, in *generatedSurvey.DeleteQuestionRequest) (*generatedSurvey.DeleteQuestionResponse, error) {
	content := in.Content
	err := h.uc.DeleteQuestion(ctx, content)
	if err != nil {
		return nil, fmt.Errorf("delete question grpc: %w", err)
	}
	response := &generatedSurvey.DeleteQuestionResponse{}
	return response, nil
}

func (h *GRPCHandler) GetQuestions(ctx context.Context, in *generatedSurvey.GetQuestionsRequest) (*generatedSurvey.GetQuestionResponse, error) {
	questions, err := h.uc.GetQuestions(ctx)
	if err != nil {
		return nil, fmt.Errorf("get questions grpc: %w", err)
	}
	var respQuestions []*generatedSurvey.AdminQuestion
	for _, question := range questions {
		respQuestions = append(respQuestions, &generatedSurvey.AdminQuestion{
			Content: question.Content,
			Grade:   int32(question.Grade),
		})
	}
	response := &generatedSurvey.GetQuestionResponse{
		Questions: respQuestions,
	}
	return response, nil
}
