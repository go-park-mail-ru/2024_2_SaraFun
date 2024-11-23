package grpc

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	generatedSurvey "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/survey/delivery/grpc/gen"
	"go.uber.org/zap"
)

type SurveyUsecase interface {
	AddSurvey(ctx context.Context, survey models.Survey) (int, error)
	GetSurveyInfo(ctx context.Context) (map[string]models.SurveyStat, error)
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
