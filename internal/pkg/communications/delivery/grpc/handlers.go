package communicationsgrpc

//go:generate mockgen -source=handlers.go -destination=mocks/reaction_usecase_mock.go -package=mocks

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2024_2_SaraFun/internal/models"
	generatedCommunications "github.com/go-park-mail-ru/2024_2_SaraFun/internal/pkg/communications/delivery/grpc/gen"
)

type ReactionUseCase interface {
	AddReaction(ctx context.Context, reaction models.Reaction) error
	GetMatchList(ctx context.Context, userId int) ([]int, error)
	GetReactionList(ctx context.Context, userId int) ([]int, error)
}

type GrpcCommunicationsHandler struct {
	generatedCommunications.CommunicationsServer
	reactionUC ReactionUseCase
}

func NewGrpcCommunicationHandler(uc ReactionUseCase) *GrpcCommunicationsHandler {
	return &GrpcCommunicationsHandler{reactionUC: uc}
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
