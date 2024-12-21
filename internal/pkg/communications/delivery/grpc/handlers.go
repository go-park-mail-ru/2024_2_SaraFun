package communicationsgrpc

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	generatedCommunications "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/communications/delivery/grpc/gen"
	"go.uber.org/zap"
)

//go:generate mockgen -destination=./mocks/mock_ReactionService.go -package=communications_mocks . ReactionUseCase
type ReactionUseCase interface {
	AddReaction(ctx context.Context, reaction models.Reaction) error
	GetMatchList(ctx context.Context, userId int) ([]int, error)
	GetReactionList(ctx context.Context, userId int) ([]int, error)
	GetMatchTime(ctx context.Context, firstUser int, secondUser int) (string, error)
	GetMatchesBySearch(ctx context.Context, userID int, search string) ([]int, error)
	UpdateOrCreateReaction(ctx context.Context, reaction models.Reaction) error
	CheckMatchExists(ctx context.Context, firstUser int, secondUser int) (bool, error)
}

type GrpcCommunicationsHandler struct {
	generatedCommunications.CommunicationsServer
	reactionUC ReactionUseCase
	logger     *zap.Logger
}

func NewGrpcCommunicationHandler(uc ReactionUseCase, logger *zap.Logger) *GrpcCommunicationsHandler {
	return &GrpcCommunicationsHandler{reactionUC: uc, logger: logger}
}

func (h *GrpcCommunicationsHandler) AddReaction(ctx context.Context,
	in *generatedCommunications.AddReactionRequest) (*generatedCommunications.AddReactionResponse, error) {
	reaction := models.Reaction{
		Id:       int(in.Reaction.ID),
		Author:   int(in.Reaction.Author),
		Receiver: int(in.Reaction.Receiver),
		Type:     in.Reaction.Type,
	}
	err := h.reactionUC.AddReaction(ctx, reaction)
	if err != nil {
		return nil, fmt.Errorf("grpc add reaction error: %w", err)
	}
	return &generatedCommunications.AddReactionResponse{}, nil
}

func (h *GrpcCommunicationsHandler) GetMatchList(ctx context.Context,
	in *generatedCommunications.GetMatchListRequest) (*generatedCommunications.GetMatchListResponse, error) {
	userId := int(in.UserID)
	users, err := h.reactionUC.GetMatchList(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("grpc get match list error: %w", err)
	}
	resUsers := make([]int32, len(users))
	for i, v := range users {
		resUsers[i] = int32(v)
	}
	res := &generatedCommunications.GetMatchListResponse{
		Authors: resUsers,
	}
	return res, nil
}

func (h *GrpcCommunicationsHandler) GetReactionList(ctx context.Context,
	in *generatedCommunications.GetReactionListRequest) (*generatedCommunications.GetReactionListResponse, error) {
	userId := int(in.UserId)
	users, err := h.reactionUC.GetReactionList(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("grpc get reaction list error: %w", err)
	}
	resUsers := make([]int32, len(users))

	for i, v := range users {
		resUsers[i] = int32(v)
	}

	res := &generatedCommunications.GetReactionListResponse{Receivers: resUsers}

	return res, nil
}

func (h *GrpcCommunicationsHandler) GetMatchTime(ctx context.Context,
	in *generatedCommunications.GetMatchTimeRequest) (*generatedCommunications.GetMatchTimeResponse, error) {
	firstUser := int(in.FirstUser)
	secondUser := int(in.SecondUser)
	time, err := h.reactionUC.GetMatchTime(ctx, firstUser, secondUser)
	if err != nil {
		return nil, fmt.Errorf("grpc get match time error: %w", err)
	}
	response := &generatedCommunications.GetMatchTimeResponse{
		Time: time,
	}
	return response, nil
}

func (h *GrpcCommunicationsHandler) GetMatchesBySearch(ctx context.Context,
	in *generatedCommunications.GetMatchesBySearchRequest) (*generatedCommunications.GetMatchesBySearchResponse, error) {
	userId := int(in.UserID)
	search := in.Search

	authors, err := h.reactionUC.GetMatchesBySearch(ctx, userId, search)
	if err != nil {
		return nil, fmt.Errorf("grpc get matches by search error: %w", err)
	}

	var respAuthors []int32
	for _, v := range authors {
		respAuthors = append(respAuthors, int32(v))
	}

	response := &generatedCommunications.GetMatchesBySearchResponse{
		Authors: respAuthors,
	}
	return response, nil
}

func (h *GrpcCommunicationsHandler) UpdateOrCreateReaction(ctx context.Context,
	in *generatedCommunications.UpdateOrCreateReactionRequest) (*generatedCommunications.UpdateOrCreateReactionResponse, error) {
	reaction := models.Reaction{
		Author:   int(in.Reaction.Author),
		Receiver: int(in.Reaction.Receiver),
		Type:     in.Reaction.Type,
	}
	h.logger.Info("reaction", zap.Any("reaction", reaction))
	err := h.reactionUC.UpdateOrCreateReaction(ctx, reaction)
	if err != nil {
		return nil, fmt.Errorf("grpc update reaction error: %w", err)
	}
	return &generatedCommunications.UpdateOrCreateReactionResponse{}, nil
}

func (h *GrpcCommunicationsHandler) CheckMatchExists(ctx context.Context,
	in *generatedCommunications.CheckMatchExistsRequest) (*generatedCommunications.CheckMatchExistsResponse, error) {
	firstUser := int(in.FirstUser)
	secondUser := int(in.SecondUser)
	h.logger.Info("check match exists grpc")
	response, err := h.reactionUC.CheckMatchExists(ctx, firstUser, secondUser)
	if err != nil {
		return nil, fmt.Errorf("grpc check match exists error: %w", err)
	}
	return &generatedCommunications.CheckMatchExistsResponse{Exists: response}, nil
}
